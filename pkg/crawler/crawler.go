package crawler

import (
	"context"
	"fmt"
	"log"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
	"github.com/gomills/gofocusedcrawler/pkg/utils"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/sync/errgroup"
)

// Crawl is the entry point of crawling to a given domain. Returns the crawling outcome.
func Crawl(domain string, config *config.Config) *CrawlingOutcome {

	// 0. get registered domain
	registeredDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		log.Print(err)
		return GetCrawlingOutcome(domain, 0, 0, "reg_domain_failed")
	}

	// 0.1 get client and headers
	client, headers, err := getStealthHttps()
	if err != nil {
		log.Print(err)
		return GetCrawlingOutcome(domain, 0, 0, err.Error())
	}

	// 1. get seed urls, which are brute forced url paths and subdomains
	seedUrls, err := utils.GetBruteForcedUrls(domain, registeredDomain)
	if err != nil {
		log.Print(err)
		return GetCrawlingOutcome(domain, 0, 0, err.Error())
	}

	// 2. instantiate queue
	qp := queue.NewQueue(config.Workers)

	// 3. instantiate context. It will be cancelled on either:
	//	- timeout
	//	- empty queue
	//	- first 429 status code of any worker
	timeOutDuration := time.Duration(config.TimeOutDuration) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeOutDuration)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// 4. seed the queue
	for _, seedUrl := range seedUrls {
		qp.AddUrl(ctx, seedUrl)
	}

	start := time.Now()

	// 5. spam pool of workers
	for range config.Workers {

		g.Go(func() error {
			return worker(qp, ctx, client, headers, config, domain, registeredDomain)
		})

	}

	// 6. await the error group cancellation
	err = g.Wait()

	s_crawling := time.Since(start).Seconds()

	return GetCrawlingOutcome(domain, qp.GetNCrawledUrls(), s_crawling, err.Error())

}

// worker is a consumer&producer of the urls queue.
// It returns error for: empty queue, timeout or a single 429 status code
func worker(qp *queue.Queue, ctx context.Context, client tls_client.HttpClient, headers http.Header, config *config.Config, domain string, registeredDomain string) error {

	for {

		// 1. consume a url from the queue. If nil it's empty queue.
		url := qp.TakeUrl(ctx)
		if url == nil {
			return fmt.Errorf("empty_queue")
		}

		log.Println(url.String())

		// 2. crawl it. Error is only for 429 status code hit.
		err := crawlUrl(ctx, client, headers, config, url, qp, domain, registeredDomain)
		if err != nil {
			return err
		}

	}
}
