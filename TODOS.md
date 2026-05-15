# TODOs

1. crawlUrl() gets passed the queue instead of returning URLs and having the worker() append them. Not too idiomatic but reduces allocations. DONE
2. unmarshal crawling outcome into .json. DONE.
3. add bool flag to respoect robots.txt or attack it. PASS
4. aggressive bruteforcing: bool for a config, where we don't bruteforce subdomains DONE
5. secret detection should be optional DONE
6. cookies cumulator DONE
7. add a listener to OS signals: SIGTERM, SIGINT, SIGSTP
8. change PrintUrls flag to interactivelogging