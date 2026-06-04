package main

import (
	"context"
	"encoding/json"
	"log"
	"os/signal"

	"github.com/gomills/millcrawler/pkg/crawler"
	"golang.org/x/sys/unix"
)

func main() {

	log.SetFlags(0)

	// load config from flags
	domain, configuration, err := parseFlagsIntoConfig()
	if err != nil {
		log.Fatalf("Error while parsing flags: %s", err)
	}

	// print banner if interactive logging enabled
	if configuration.PrintUrls {
		banner := loadBanner()
		log.Println(banner)
	}

	// handle gracefully the following OS signals: SIGTERM, SIGINT, SIGSTP.
	ctx, stop := signal.NotifyContext(context.Background(), unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)
	defer stop()

	// crawl
	crawlingOutcome := crawler.Crawl(ctx, domain, configuration)

	// return json to stdout
	j, err := json.Marshal(crawlingOutcome)
	if err != nil {
		log.Printf("Error while serializing results to JSON: %s\n", err)
	}

	log.Println(string(j))
}
