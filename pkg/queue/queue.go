package queue

import (
	"context"
	"log"
	"net/url"
	"sync"
)

type Queue struct {
	mu             sync.Mutex          // internal locker
	toCrawl        chan *url.URL       // FIFO queue of ulrs to crawl
	foundUrlsSet   map[string]struct{} // Set to track found URLs with O(1) lookup time
	workers        int                 // number of concurrent workers
	waitingWorkers int                 // current number of workers dequeueing a url to crawl
}

// Initialize a Queue
func NewQueue(nWorkers int) *Queue {
	return &Queue{
		toCrawl:        make(chan *url.URL, nWorkers*2),
		foundUrlsSet:   make(map[string]struct{}, 300),
		workers:        nWorkers,
		waitingWorkers: 0,
	}
}

// Add a URL to queue with deduplication and emptiness safety
func (q *Queue) AddUrl(ctx context.Context, newUrl *url.URL) {

	if newUrl == nil {
		// log.Println("<Q: found empty url")
		return
	}

	if ctx.Err() != nil {
		// log.Printf("<Q: %s, context cancelled %s", newUrl.String(), ctx.Err())
		return
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	// Deduplicate
	if _, ok := q.foundUrlsSet[newUrl.String()]; ok {
		// log.Printf("<Q: duplicate!: %s", newUrl.String())
		return
	}
	q.foundUrlsSet[newUrl.String()] = struct{}{}
	// log.Println("<Q: added url")

	// Hopefully this goroutine doesn't perform blocking sending. If it does, it exits by listening to the context anyway,
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

func (q *Queue) TakeUrl(ctx context.Context) *url.URL {

	if ctx.Err() != nil {
		return nil
	}

	// Check that queue is not empty for good. We know it's empty when all the workers are simultaneously waiting a URL
	q.mu.Lock()
	q.waitingWorkers++
	// log.Printf("Worker entered queue: %d waiter and %d workers", q.waitingWorkers, q.workers)
	isQueueEmpty := q.waitingWorkers == q.workers
	if isQueueEmpty {
		log.Println("Worker found the other one waiting as well!")
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
