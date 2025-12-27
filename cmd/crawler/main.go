package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/crawler"
)

func main() {

	domain := flag.String("domain", "", "domain to crawl (e.g. example.com www.example.com)")

	flag.Parse()
	if domain == nil || *domain == "" {
		log.Fatalln("need_valid_domain")
	}

	domainStr := strings.TrimSpace(*domain)

	// load configuration
	config := config.LoadConfig()

	log.Printf(">>> Start crawling %s\n\n", domainStr)

	results := crawler.Crawl(domainStr, config)

	printResults(results)
}

func printResults(outcome *crawler.CrawlingOutcome) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("CRAWL SUMMARY")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Domain:           %s\n", outcome.Domain)
	fmt.Printf("URLs Crawled:     %d\n", outcome.N_urls)
	fmt.Printf("Time Elapsed:     %.2f seconds\n", outcome.S_crawling)
	fmt.Printf("Stop Reason:      %s\n", outcome.Stop_reason)
	fmt.Println(strings.Repeat("=", 60) + "\n")
}
