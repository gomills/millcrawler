package main

import (
	"flag"
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

	crawler.Crawl(domainStr, config)

}
