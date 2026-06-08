# Millcrawler

```text
~/$ ./millcrawler -domain=crawler-test.com -printurls=true

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
- smart, custom URL validation heuristics
- optional secrets detection, although manual implementation of secrets is necessary in `pkg/secrets/secrets.go` 
- stops abruptly on first 429 antibot response status code
- extensive testing coverage for most packages

## Paradigm

Instead of betting on a massive crawl following anything but mostly on HTML's `<a href="">`, `millcrawler` invests on parsing and targetting specific static attributes of *JavaScript* and *HTML* while filtering out most of useless urls.

## Notes

By design it will include false positives urls, strings that because of .js parsing they were taken into the queue as relative urls when they were just anything else.

If interactive urls printing is enabled, all urls coming out of the queue fresh will be printed. This doesn't mean they've been successfully crawled or if they will even be crawled. For example, external urls won't be crawled, just scanned for secrets if enabled.

It performs static analysis of `.js` so constructed and formatted urls won't be crawled.

## Options

```text
~/$ ./millcrawler -h
Usage of ./millcrawler:
  -allowedextdomains string
        Comma-separated list of allowed external domains. Example: -allowedextdomains=github.com,gitlab.com (default "github.com,gitlab.com,docker.com,pastebin.com,nuget.org,bitbucket.org,s3.amazonaws.com,telegram.org,slack.com,drive.google.com,docs.google.com,codeberg.org")
  -allowedextensions string
        Comma-separated list of allowed file extensions ('' for subpages...). Example: allowedextensions=,.html,.js,.json (default ",.js,.ts,.map,.properties,.log,.cfg,.pem,.crt,.npmrc,.yarnrc,.html,.htm,.json,.yaml,.yml,.env,.conf,.xml,.txt,.py,.php,.git,.sh,.key,.go,.ini,.example,.md,.rb,.java,.cpp,.c,.pl,.zsh,.bak,.old,.sql,.db,.tar,.gz,.zip,.ovpn,.toml")
  -bruteforcesubdomains
        Enable subdomain bruteforcing. Example: -bruteforcesubdomains=true
  -cookies
        Enable cookie jar for use during crawling. Example: -cookies=true
  -domain string
        Target domain to crawl. Example: -domain=example.com
  -maxnumurls int
        Maximum number of URLs to crawl before stopping, 0 for unlimited. Example: -maxnumurls=400
  -maxpathdepth int
        Maximum URL path depth to crawl, 0 for unlimited. Example: -maxpathdepth=2 allows /path1/path2 but not /path1/path2/path3
  -printurls
        Interactive logging: banner and URLs to stdout. If disabled, it outputs .json results directly for scripting integration. Example: -printurls=true
  -scansecrets
        Decide if scan responses for secrets. Example: -scansecrets=true
  -sensitivepatterns string
        Comma-separated URL patterns that bypass max path depth restrictions. Example: sensitivepatterns=api,admin,debug
  -timeoutseconds int
        Maximum crawl duration in seconds, never unlimited. Example: -timeoutseconds=60 (default 3600)
  -workers int
        Number of concurrent crawling workers. Example: -workers=10 (default 1)
```

## Building from source

```bash
git clone https://github.com/gomills/millcrawler
cd millcrawler
go test ./... # optional: run tests
go build
./millcrawler -h
```