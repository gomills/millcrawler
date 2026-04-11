# Stealth refactor

This codebase has been updated with:

- testing files for most of packages coverage
- fixed issues and bugs
- improved documentation with README.md's on each package and inline comments
- improved predictability, now the crawler outputs same results in a structured fashion
- simplified url extraction, now they only come from .js, .html, robots.txt, sitemap.xml and error payloads

## What GOFOCUSEDCRAWLER brings new to the table:

- **Intuitive, granular control**  
  Through config, the crawler is controlled by these parameters:
  ```go
  TimeOutDuration
  AllowedExternalDomains
  MaxPathDepth
  SensitivePatterns
  AllowedExtensions
  Workers
  ```

- **Tree-Sitter parsing for `.js`.**  
  Fast, lightweight, yet precise parsing allows selective extraction of URLs from JavaScript files. It won’t miss them as long as they are within the heuristics. Thorough parsing is especially important because most sensitive endpoints are found in JavaScript.

- **Extensive URL extraction from HTML.**  
  Not limited to a `"follow a[href]"` algorithm. Like `.js` files, heuristics were extended to avoid missing important URLs. Elements and attributes were carefully selected.

- **Smart, custom URL validation heuristics.**  
  Beyond extraction, URLs are validated per the user's config:  
  - Accept URLs with allowed file extensions.  
  - Accept subpages (`.htm`, `.html`, or no extension`) only if the path depth is within limits (bypassed by sensitive patterns).  
  - Same logic applies to **any and all subdomains**.  
  - External URLs are kept only if they match allowed domains (e.g., `github.com`), but they won’t be crawled.

- **Exits on first 429 status code.**  
  Once flagged by a server, further requests are usually futile.

- **Bruteforces certain URLs**, such as `sitemap` and some subdomains like `dev.` or `staging.`

- **URL regexing on error payloads** for debugging sensitive endpoint leaks. This could be further expanded with headers in the future.

- **Fully synchronized goroutines.** Once a stop signal is fired, all close automatically thorugh shared context immediately, and no further requests proceed. 

- **Fully modularized, encapsulated codebase.**  
  Well documented, intuitive, easily modifiable and refactorable, and can be deployed in distributed architectures with minor changes.
