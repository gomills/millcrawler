package crawler

type CrawlingOutcome struct {
	domain      string
	n_urls      int
	s_crawling  float64
	stop_reason string
}

func GetCrawlingOutcome(domain string, n_urls int, s_crawling float64, stop_reason string) *CrawlingOutcome {
	return &CrawlingOutcome{
		domain:      domain,
		n_urls:      n_urls,
		s_crawling:  s_crawling,
		stop_reason: stop_reason,
	}
}
