package crawler

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/gomills/millcrawler/pkg/config"
	"github.com/gomills/millcrawler/pkg/queue"
	"github.com/gomills/millcrawler/pkg/utilities"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/sync/errgroup"
)

// Crawl is the entry point of crawling to a given domain. It consumes a domain and a crawling configuration and
// returns the crawling outcome in a predictable fashion.
func Crawl(ctx context.Context, domain string, config *config.Config) *CrawlingOutcome {

	// log.Printf(">>> %s\n", domain)

	// 0. get registered domain for bruteforcing subdomains (e.g domain: foo.bar.golang.com => registered domain: golang.com)
	registeredDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		// log.Println(err.Error())
		return CreateCrawlingOutcome(domain, 0, 0, err.Error(), 0, nil)
	}

	// 1. get the initial urls, which are brute forced url paths and subdomains
	initialUrls, err := utilities.GetBruteForcedUrls(domain, registeredDomain, config)
	if err != nil {
		// log.Println(err.Error())
		return CreateCrawlingOutcome(domain, 0, 0, err.Error(), 0, nil)
	}

	// 2. instantiate queue
	q := queue.NewQueue(config.Workers)

	// 3. instantiate context. It will be cancelled on either:
	//	- timeout, using the duration coming from the config
	//	- first 429 status code of any worker, which signals an anti-bot response code
	//	- empty queue
	timeOutDuration := time.Duration(config.TimeOutDuration) * time.Second

	ctx, cancel := context.WithTimeout(ctx, timeOutDuration)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// 4. seed the queue
	var wg sync.WaitGroup
	q.SeedQueue(ctx, &wg, initialUrls)
	wg.Wait()

	startTime := time.Now()

	// 5. get a stealth client
	client, err := utilities.GetStealthHttps(config)
	if err != nil {
		return CreateCrawlingOutcome(domain, 0, 0, err.Error(), 0, nil)
	}

	// 6. spam pool of workers inside the error group
	for range config.Workers {

		g.Go(func() error {
			return worker(q, ctx, config, client, domain, registeredDomain)
		})

	}

	// 7. await the error group cancellation. It will safely always happen due to the timeout.
	err = g.Wait()

	// 8. return the crawling outcome
	return returnCrawlingOutcome(q, domain, startTime, err)

}

// worker is a consumer&producer of the urls queue.
// It returns error for: empty queue, timeout or a single 429 status code
func worker(q *queue.Queue, ctx context.Context, config *config.Config, client tls_client.HttpClient, domain string, registeredDomain string) error {

	// main loop [crawl -> populate queue -> crawl]
	for {

		// 1. check if we hit maximum crawled urls
		if config.MaxNCrawledUrls > 0 {
			if q.GetNCrawledUrls() >= config.MaxNCrawledUrls {
				return errors.New("hit crawled urls limit")
			}
		}

		// 2. consume a url from the queue.
		url, err := q.TakeUrl(ctx)
		if err != nil {
			return err
		}
		if url == nil {
			return errors.New("empty queue")
		}

		if config.PrintUrls {
			log.Println(url.String())
		}

		// 3. crawl it
		err = crawlUrl(ctx, client, config, url, q, domain, registeredDomain)
		if err != nil {
			return err
		}

	}
}

// returnCrawlingOutcome creates a crawling outcome and returns it
func returnCrawlingOutcome(q *queue.Queue, domain string, startTime time.Time, err error) *CrawlingOutcome {
	numURLs := q.GetNCrawledUrls()
	durationSeconds := time.Since(startTime).Seconds()
	stopReason := err.Error()
	secretsFound := q.GetNSecrets()
	secretsCopy := q.GetSecretsCopy()
	return CreateCrawlingOutcome(domain, numURLs, durationSeconds, stopReason, secretsFound, secretsCopy)
}
