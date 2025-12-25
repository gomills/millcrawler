package crawler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
	"github.com/gomills/gofocusedcrawler/pkg/utils"
	"github.com/google/uuid"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/sync/errgroup"
)

func StartCrawling(domain string, config *config.Config) {

	// Get registered domain
	registeredDomain, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		log.Print(err)
		return
	}

	// Gather initial urls (domain, robots and some other brute forced urls)
	initialUrls, err := utils.GetBruteForcedUrls(domain, registeredDomain)
	if err != nil {
		log.Print(err)
		return
	}

	// Instantiate queue
	qp := queue.NewQueue(config.Workers)

	// Instantiate context. It will be cancelled on:
	//	- Timeout
	//	- Empty queue
	//	- On first 429 status code of any worker
	//	- On Ctrl+C signal
	timeOutDuration := time.Duration(config.TimeOutDuration) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeOutDuration)
	g, ctx := errgroup.WithContext(ctx)

	// Enqueue the initial URLs
	for _, urlInfo := range initialUrls {
		qp.AddUrl(ctx, urlInfo)
	}

	// Set the relay for the keyboard interrupt signal
	// g.Go(func() error {
	// 	utils.SetKeyboardInterruptSignal(ctx, cancel)
	// 	return nil
	// })

	// Set the pool of workers
	for range config.Workers {

		g.Go(func() error {
			return worker(qp, ctx, cancel, config, domain, registeredDomain)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println("Error:", err)
	}

}

func worker(qp *queue.Queue, ctx context.Context, cancel context.CancelFunc, config *config.Config, domain string, registeredDomain string) error {

	workerId := uuid.New().String()

	for {

		url := qp.TakeUrl(ctx)
		if url == nil {
			cancel()
			return nil
		}

		log.Println(url.String())

		stopCrawling := crawlUrl(ctx, config, url, qp, domain, registeredDomain, workerId)
		if stopCrawling {
			cancel()
			return fmt.Errorf("429 hit")
		}

	}
}
