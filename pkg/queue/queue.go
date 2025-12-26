package queue

import (
	"context"
	"net/url"
	"sync"
)

type Queue struct {
	mu             sync.Mutex          // internal locker
	toCrawl        chan *url.URL       // buffered FIFO queue of ulrs unused for crawling
	foundUrlsSet   map[string]struct{} // set to store all found URLs, either used and unused
	workers        int                 // number of concurrent workers
	waitingWorkers int                 // current number of workers dequeueing a url to crawl
}

// instantiate a Queue
func NewQueue(nWorkers int) *Queue {
	return &Queue{
		toCrawl:        make(chan *url.URL, nWorkers*2),
		foundUrlsSet:   make(map[string]struct{}, 400),
		workers:        nWorkers,
		waitingWorkers: 0,
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

// take a URL from the queue. Returns nil for context closure and empty-for-good queue.
func (q *Queue) TakeUrl(ctx context.Context) *url.URL {

	// monitor context closure
	if ctx.Err() != nil {
		return nil
	}

	// check that queue is not empty-for-good. We know it's empty-for-good when all the workers are simultaneously waiting a URL
	q.mu.Lock()
	q.waitingWorkers++
	// log.Printf("Worker entered queue: %d waiter and %d workers", q.waitingWorkers, q.workers)

	isQueueEmpty := q.waitingWorkers == q.workers
	if isQueueEmpty {
		// log.Println("worker found the other one waiting as well!")
		q.waitingWorkers--
		q.mu.Unlock()
		return nil
	}
	q.mu.Unlock()

	select {
	case giveUrl := <-q.toCrawl:
		q.mu.Lock()
		q.waitingWorkers--
		q.mu.Unlock()
		return giveUrl
	case <-ctx.Done():
		q.mu.Lock()
		q.waitingWorkers--
		q.mu.Unlock()
		return nil
	}

}
