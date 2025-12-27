package crawler

type CrawlingOutcome struct {
	Domain      string
	N_urls      int
	S_crawling  float64
	Stop_reason string
}

func GetCrawlingOutcome(domain string, n_urls int, s_crawling float64, stop_reason string) *CrawlingOutcome {
	return &CrawlingOutcome{
		Domain:      domain,
		N_urls:      n_urls,
		S_crawling:  s_crawling,
		Stop_reason: stop_reason,
	}
}
