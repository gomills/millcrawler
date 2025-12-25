# Stealth refactor

General TODOS


- Define URL extraction perfectly. URLs must come from:

    - .js: already working correctly with tree-sitter. maybe revise some settings, but not too much, it already works, don't over engineer.

    - .html: same as with .js.

    - other scripts: they just will add too much overhead because will have to be regexed. We'll ignore them.
    
    - HTTPs responses: we'll inspect headers for [debugging endpoints] and regex error payloads.

- Stealthness: requests (TLS fingerprint) and CloudFare (antibot capacities).

- Improve outcomes: we have to know exactly how the crawl went and the info we want to store from it. we should store in a dataframe: URLs crawled, time spent, reason for stop (in code format: timeout, anti_bot, etc.). this will improve predictability.

- Stop signals should only be: 429 / timeout

- Bundle config file in compile time

Side quests TODOS for refactoring:

- Simplify config, just accept .json. DONE
- Discriminate the /crawler file. Separate into crawl_urls and crawler DONE
- Merge urlextraction with the corresponding crawl_urls DONE
- Deprecate kw_interrupt_signal DONE
