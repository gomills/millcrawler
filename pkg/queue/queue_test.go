package queue

import (
	"context"
	"net/url"
	"testing"
	"time"
)

// TestAddUrlDeduplication tests that AddUrl deduplicates identical URLs
func TestAddUrlDeduplication(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := NewQueue(1)

	testUrl, _ := url.Parse("https://example.com/test")

	// Add same URL twice
	q.AddUrl(ctx, testUrl)
	q.AddUrl(ctx, testUrl)

	time.Sleep(50 * time.Millisecond)

	// Should only be 1 URL in the set despite adding twice
	if count := q.GetNCrawledUrls(); count != 1 {
		t.Errorf("Expected 1 URL after deduplication, got %d", count)
	}
}

// TestAddUrlNilSafety tests that AddUrl handles nil URLs gracefully
func TestAddUrlNilSafety(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := NewQueue(1)

	// Add nil URL - should not panic and not be added
	q.AddUrl(ctx, nil)

	time.Sleep(50 * time.Millisecond)

	if count := q.GetNCrawledUrls(); count != 0 {
		t.Errorf("Expected 0 URLs after adding nil, got %d", count)
	}
}

// TestAddUrlContextClosure tests that AddUrl respects context cancellation
func TestAddUrlContextClosure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	q := NewQueue(1)

	testUrl, _ := url.Parse("https://example.com/test")

	// Add URL before context cancellation
	q.AddUrl(ctx, testUrl)
	time.Sleep(10 * time.Millisecond)

	// Cancel context
	cancel()
	time.Sleep(10 * time.Millisecond)

	// Try to add another URL after context is cancelled - should not be added
	testUrl2, _ := url.Parse("https://example.com/test2")
	q.AddUrl(ctx, testUrl2)

	time.Sleep(50 * time.Millisecond)

	// Should only have the first URL (added before cancellation)
	if count := q.GetNCrawledUrls(); count != 1 {
		t.Errorf("Expected 1 URL, got %d", count)
	}

	urls := q.GetFoundUrls()
	if len(urls) > 0 && urls[0] != "https://example.com/test" {
		t.Errorf("Expected first URL to be test, got %s", urls[0])
	}
}

// TestTakeUrlReturnsExactUrl tests that TakeUrl returns the URL exactly as it was added
func TestTakeUrlReturnsExactUrl(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := NewQueue(2)

	testUrlStr := "https://example.com/path/to/resource?query=value"
	testUrl, _ := url.Parse(testUrlStr)

	q.AddUrl(ctx, testUrl)
	time.Sleep(10 * time.Millisecond)

	takenUrl := q.TakeUrl(ctx)

	if takenUrl == nil {
		t.Fatal("Expected to take a URL, got nil")
	}

	if takenUrl.String() != testUrlStr {
		t.Errorf("Expected URL %q, got %q", testUrlStr, takenUrl.String())
	}
}

// TestTakeUrlContextClosure tests that TakeUrl returns nil when context is cancelled
func TestTakeUrlContextClosure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	q := NewQueue(2)

	testUrl, _ := url.Parse("https://example.com/test")
	q.AddUrl(ctx, testUrl)
	time.Sleep(10 * time.Millisecond)

	// Cancel context
	cancel()
	time.Sleep(10 * time.Millisecond)

	// TakeUrl should return nil because context is cancelled
	takenUrl := q.TakeUrl(ctx)

	if takenUrl != nil {
		t.Errorf("Expected nil after context cancellation, got %v", takenUrl)
	}
}

func TestQueue_EmptyForGood(t *testing.T) {
	nWorkers := 5

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := NewQueue(nWorkers)

	// Add exactly nWorkers URLs
	// url1, _ := url.Parse("https://example.com/1")
	// url2, _ := url.Parse("https://example.com/2")
	// url3, _ := url.Parse("https://example.com/3")
	// q.AddUrl(ctx, url1)
	// q.AddUrl(ctx, url2)
	// q.AddUrl(ctx, url3)

	for i := 0; i < nWorkers-1; i++ {

		go func() {
			u := q.TakeUrl(ctx)
			if u != nil {
				panic("received url")
			}
		}()

	}

	// Wait until all workers - 1 are waiting
	time.Sleep(100 * time.Millisecond)

	// Now queue is empty, all workers - 1 waiting, now should be nil
	sentinel := q.TakeUrl(ctx)
	if sentinel != nil {
		t.Fatal("last url wasn't nil")
	}

}
