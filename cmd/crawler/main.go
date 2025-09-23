package main

import (
	"flag"
	"log"
	"strings"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/crawler"
)

func main() {

	domain := flag.String("domain", "", "domain to crawl (e.g. example.com or https://example.com)")

	flag.Parse() // Parse CLI flags before using them
	if domain == nil || *domain == "" {
		log.Fatalln("A domain is necessary to call gofocusedcrawler!")
	}
	domainStr := strings.TrimSpace(*domain)

	// Load configuration
	config, err := config.LoadConfig("path/to/your/config_example.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Start crawling %s\n\n", domainStr)

	crawler.StartCrawling(domainStr, config)

}
