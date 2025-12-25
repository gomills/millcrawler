# Stealth refactor

General TODOS

- Bundle config file in compile time

- Stealthness: TLS fingerprint and CloudFare

- Define URL extraction perfectly

- Improve outcomes: we have to know exactly how the crawl went and the info we want to store from it

- Define stop signals: 429 / timeout

Side quests TODOS for refactoring:

- Simplify config, just accept .json. DONE
- Discriminate the /crawler file. Separate into crawl_urls and crawler DONE
- Merge urlextraction with the corresponding crawl_urls DONE
- Deprecate kw_interrupt_signal
