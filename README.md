# Millcrawler

```text
username@COMPUTERID:~/$ ./millcrawler -domain=crawler-test.com -printurls=true

                _ _ _                         _           
               (_) | |                       | |          
      _ __ ___  _| | | ___ _ __ __ ___      _| | ___ _ __ 
     | '_ ' _ \| | | |/ __| '__/ _' \ \ /\ / / |/ _ \ '__|
     | | | | | | | | | (__| | | (_| |\ V  V /| |  __/ |   
     |_| |_| |_|_|_|_|\___|_|  \__,_| \_/\_/ |_|\___|_|   
----------------------------------------------------------- by github.com/gomills
                             

https://internal.crawler-test.com/
https://crawler-test.com/
https://crawler-test.com/robots.txt
https://crawler-test.com/sitemap.xml
https://crawler-test.com/mobile/separate_desktop_with_different_links_out
https://crawler-test.com/mobile/separate_desktop_with_mobile_not_subdomain
https://crawler-test.com/mobile/dynamic
[···]
^C
{"Domain":"crawler-test.com","NumURLs":28,"DurationSeconds":5.691633071,"StopReason":"context canceled","SecretsFound":0,"SecretsMap":{}}
```

Configurable web crawler for static crawling (no java execution) that outputs a predictable `.json`, characterized by:

- browser-fingerprinted-TLS client
- targeted urls extraction from: `.js`, `.html`, `robots.txt`, `sitemap.xml` and HTTP error payloads
- *tree-sitter* parsing for javascript
- Smart, custom URL validation heuristics
- optional secrets detection, although manual implementation of secrets is necessary in `pkg/secrets/secrets.go` 
- stops abruptly on first 429 antibot response status code
- testing coverage for most packages

## Paradigm

Instead of betting on a massive crawl following anything but mostly on HTML's `<a href="">`, `millcrawler` invests on parsing and targetting specific static attributes of *JavaScript* and *HTML* while filtering out most of useless urls.

## Options

```text
username@COMPUTERID:~/$ ./millcrawler -h
Usage of ./millcrawler:
  -allowedextdomains string
        Comma-separated list of allowed external domains. Example: 'github.com,gitlab.com'
  -allowedextensions string
        Comma-separated list of allowed file extensions. Example: '.html,.js,.json,.git,.yaml' (default ",.htm,.html,.js")
  -bruteforcesubdomains
        Enable subdomain bruteforcing. E.g: 'true' for testing.example.com, admin.example.com, ... (default true)
  -cookies
        Enable cookie jar for use during crawling
  -domain string
        Target domain to crawl. Example: 'example.com'
  -maxnumurls int
        Maximum number of URLs to crawl before stopping (default 100)
  -maxpathdepth int
        Maximum URL path depth to crawl. Example: '2' allows '/path1/path2' (default 1)
  -printurls
        Interactive logging: banner, URLs, to stdout.
  -scansecrets
        Decide if scan responses for secrets
  -sensitivepatterns string
        Comma-separated URL patterns that bypass max path depth restrictions. Example: 'api,admin,debug'
  -timeoutseconds int
        Maximum crawl duration in seconds. Example: '60' (default 60)
  -workers int
        Number of concurrent crawling workers (default 1)
```

## Building from source

```bash
git clone https://github.com/gomills/millcrawler
cd millcrawler
go build
./millcrawler -h
```