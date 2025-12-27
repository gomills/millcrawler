# Stealth refactor

General TODOS


- URLs must come from, and only from:

    - .js: improve AST parent types. DONE

    - .html: same as with .js. DONE

    - other scripts: ignore all except robots.txt and sitemap.yaml DONE
    
    - HTTPs responses: headers for [debugging endpoints] and error payloads.

- Stealthness: 
    - TLS fingerprinting: https://github.com/refraction-networking/utls
    - CloudFare: we'll model cloudfare with TLS fingerprinting as well.

- Improve outcome: to improve predictability we should store in a dataframe: URLs crawled, time spent, reason for stop (in code format: timeout, anti_bot, etc.). DONE

- Take a look at more status codes for bot blocking

- Optimize

- Improve external url handling. DONE

- Stop signals should only be: 429 / timeout DONE

Side quests TODOS for refactoring:

- Simplify config, just accept .json. DONE
- Discriminate the /crawler file. Separate into crawl_urls and crawler DONE
- Merge urlextraction with the corresponding crawl_urls DONE
- Deprecate kw_interrupt_signal DONE
