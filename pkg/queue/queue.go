package queue

import (
	"context"
	"net/url"
	"sync"
)

// it passes 'go test -race -count=1 ./pkg/queue -v'
type Queue struct {
	mu             sync.Mutex          // internal locker
	toCrawl        chan *url.URL       // buffered FIFO queue of uncrawled ulrs (before sending to this channel they immediately get into foundUrlsSet)
	foundUrlsSet   map[string]struct{} // set to store all found URLs, either used and unused
	secretsMap     map[string]string   // key-value map of [secret]-[source url]
	workers        int                 // number of concurrent workers
	waitingWorkers int                 // current number of workers dequeueing a url to crawl
}

// instantiate a Queue
func NewQueue(nWorkers int) *Queue {
	return &Queue{
		toCrawl:        make(chan *url.URL, 2000),
		foundUrlsSet:   make(map[string]struct{}, 2000),
		secretsMap:     make(map[string]string, 3),
		workers:        nWorkers,
		waitingWorkers: 0,
	}
}

// SeedQueue adds the initial urls to the queue with a wait group primitive avoiding
// race conditions when initializing the crawler. This function relies on that
// the toCrawl urls channel is buffered so sending to it won't be blocking.
func (q *Queue) SeedQueue(ctx context.Context, wg *sync.WaitGroup, urls []*url.URL) {
	for _, seedUrl := range urls {

		if _, ok := q.foundUrlsSet[seedUrl.String()]; ok {
			continue
		} else {
			q.foundUrlsSet[seedUrl.String()] = struct{}{}
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			case q.toCrawl <- seedUrl:
				return
			}
		}()
	}
}

// add a URL to queue with deduplication and emptiness (nil) safety
func (q *Queue) AddUrl(ctx context.Context, newUrl *url.URL) {
	if newUrl == nil {
		// log.Println("<Q: found empty url")
		return
	}

	// monitor context closure
	if ctx.Err() != nil {
		// log.Printf("<Q: %s, context cancelled %s", newUrl.String(), ctx.Err())
		return
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	// avoid duplicates
	if _, ok := q.foundUrlsSet[newUrl.String()]; ok {
		// log.Printf("<Q: duplicate!: %s", newUrl.String())
		return
	} else {
		q.foundUrlsSet[newUrl.String()] = struct{}{}
	}
	// log.Println("<Q: added url")

	// hopefully this goroutine doesn't perform blocking sending.
	// If it does, it exits by listening to the context anyway,
	// so there is your memory safety!
	go func() {
		select {
		case <-ctx.Done():
			return
		case q.toCrawl <- newUrl:
			return
		}
	}()

}

// take a URL from the queue. Returns an error for empty-for-good queue and context cancellation.
func (q *Queue) TakeUrl(ctx context.Context) (*url.URL, error) {

	// monitor context closure
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	q.mu.Lock()
	q.waitingWorkers++
	q.mu.Unlock()

	if err := ctx.Err(); err != nil {
		q.mu.Lock()
		q.waitingWorkers--
		q.mu.Unlock()
		return nil, err
	}

	select {
	// regular case, take a url successfully instantly
	case url := <-q.toCrawl:
		q.mu.Lock()
		q.waitingWorkers--
		q.mu.Unlock()
		return url, nil
	default:

		// check that queue is not empty-for-good. We know it's empty-for-good when
		// all the workers are simultaneously waiting a URL
		q.mu.Lock()
		queueIsEmpty := q.waitingWorkers == q.workers
		q.mu.Unlock()
		if queueIsEmpty {

			q.mu.Lock()
			q.waitingWorkers--
			q.mu.Unlock()
			return nil, nil

		} else {
			// second regular case, take a url after waiting
			select {
			case url := <-q.toCrawl:
				q.mu.Lock()
				q.waitingWorkers--
				q.mu.Unlock()
				return url, nil
			case <-ctx.Done():
				q.mu.Lock()
				q.waitingWorkers--
				q.mu.Unlock()
				return nil, ctx.Err()
			}
		}
	}

}

// add a secret to the secrets vault with deduplication
func (q *Queue) AddSecret(secret string, url string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// deduplicate it
	if _, ok := q.secretsMap[secret]; ok {
		return
	} else {
		q.secretsMap[secret] = url
	}
}

// get number of crawled urls (it's computed from found minus to crawl)
func (q *Queue) GetNCrawledUrls() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return (len(q.foundUrlsSet) - len(q.toCrawl))
}

// get number of found urls
func (q *Queue) GetNFoundUrls() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return (len(q.foundUrlsSet))
}

// get number of secrets in the secrets vault
func (q *Queue) GetNSecrets() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.secretsMap)
}

// get a copy of the secrets-source map
func (q *Queue) GetSecretsCopy() map[string]string {
	q.mu.Lock()
	defer q.mu.Unlock()

	secretsCopy := make(map[string]string, len(q.secretsMap))
	for secret, url := range q.secretsMap {
		secretsCopy[secret] = url
	}
	return secretsCopy
}

// get all found URLs as a slice.
func (q *Queue) GetFoundUrls() []string {
	q.mu.Lock()
	defer q.mu.Unlock()

	urls := make([]string, 0, len(q.foundUrlsSet))
	for url := range q.foundUrlsSet {
		urls = append(urls, url)
	}
	return urls
}
