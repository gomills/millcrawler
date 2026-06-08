# TODOs

1. crawlUrl() gets passed the queue instead of returning URLs and having the worker() append them. Not too idiomatic but reduces allocations. DONE
2. unmarshal crawling outcome into .json. DONE.
3. add bool flag to respoect robots.txt or attack it. PASS
4. aggressive bruteforcing: bool for a config, where we don't bruteforce subdomains DONE
5. secret detection should be optional DONE
6. cookies cumulator DONE
7. add a listener to OS signals: SIGTERM, SIGINT, SIGSTP DONE
8. change PrintUrls flag to interactivelogging PASS
9. revisit the extraction from .js; let it be deeper DONE
10. sane defaults with sentinel values (-1 etc) DONE
11. improve external urls handling and clarify what happens with them. DONE
12. optimization: if secret scanning is not enabled, why would we request external urls???
13. improve headers and user agent
14. it needs general optimizations, some stealth improvements, redundant testing, deep analysis of where we're taking urls right now from .js and .html, for being even more robust (it's very functional right now as is...)