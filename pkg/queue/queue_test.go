package queue

import (
	"context"
	"net/url"
	"strconv"
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
	if count := q.GetNFoundUrls(); count != 1 {
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

	if count := q.GetNFoundUrls(); count != 0 {
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
	if count := q.GetNFoundUrls(); count != 1 {
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

	takenUrl, _ := q.TakeUrl(ctx)

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
	takenUrl, err := q.TakeUrl(ctx)

	if takenUrl != nil {
		t.Errorf("Expected nil after context cancellation, got %v", takenUrl)
	}
	if err != ctx.Err() {
		t.Errorf("Expected error '%s', got '%s'", context.Canceled, err)
	}

}

// TestQueue_EmptyForGood tests that the queue returns a nil url when the
// queue is empty (workers waiting for urls = total workers)
func TestQueue_EmptyForGood(t *testing.T) {
	nWorkers := 5

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := NewQueue(nWorkers)

	for i := 0; i < nWorkers-1; i++ {

		go func() {
			u, _ := q.TakeUrl(ctx)
			if u != nil {
				panic("received url")
			}
		}()

	}

	// Wait until all workers - 1 are waiting
	time.Sleep(100 * time.Millisecond)

	// Now queue is empty, all workers - 1 waiting, now should be nil
	sentinel, _ := q.TakeUrl(ctx)
	if sentinel != nil {
		t.Fatal("last url wasn't nil")
	}

}

// TestAddSecretStoresExactly tests that AddSecret stores the secret and URL exactly as provided
func TestAddSecretStoresExactly(t *testing.T) {
	q := NewQueue(1)

	secret := "aws_key_123456781241"
	sourceUrl := "https://example.com/config.json"

	q.AddSecret(secret, sourceUrl)

	secretsCopy := q.GetSecretsCopy()

	if len(secretsCopy) != 1 {
		t.Errorf("Expected 1 secret, got %d", len(secretsCopy))
	}

	if storedUrl, ok := secretsCopy[secret]; !ok {
		t.Errorf("Expected secret '%q' to be stored, but not found", secret)
	} else if storedUrl != sourceUrl {
		t.Errorf("Expected URL '%q', got '%q'", sourceUrl, storedUrl)
	}
}

// TestAddSecretDeduplication tests that adding the same secret twice results in only one entry
func TestAddSecretDeduplication(t *testing.T) {
	q := NewQueue(1)

	secret := "api_token_abcdef23c23x2r3x22"
	sourceUrl := "https://example.com/credentials"

	q.AddSecret(secret, sourceUrl)
	q.AddSecret(secret, sourceUrl)

	if count := q.GetNSecrets(); count != 1 {
		t.Errorf("Expected 1 secret after deduplication, got %d", count)
	}
}

// TestGetNSecrets_OneSecret tests that GetNSecrets returns 1 when exactly one secret is stored
func TestGetNSecrets_OneSecret(t *testing.T) {
	q := NewQueue(1)

	q.AddSecret("secret1", "https://example.com/source1")

	if count := q.GetNSecrets(); count != 1 {
		t.Errorf("Expected 1 secret, got %d", count)
	}
}

// TestGetNSecrets_ThreeSecrets tests that GetNSecrets returns 3 when three distinct secrets are stored
func TestGetNSecrets_ThreeSecrets(t *testing.T) {
	q := NewQueue(1)

	q.AddSecret("secret1", "https://example.com/source1")
	q.AddSecret("secret2", "https://example.com/source2")
	q.AddSecret("secret3", "https://example.com/source3")

	if count := q.GetNSecrets(); count != 3 {
		t.Errorf("Expected 3 secrets, got %d", count)
	}
}

// TestGetSecretsCopy_Isolation tests that modifying the returned copy does not affect internal state
func TestGetSecretsCopy_Isolation(t *testing.T) {
	q := NewQueue(1)

	secret1 := "secret_original"
	url1 := "https://example.com/original"

	q.AddSecret(secret1, url1)

	// Get a copy and modify it
	secretsCopy := q.GetSecretsCopy()
	secretsCopy["modified_secret"] = "https://example.com/modified"
	delete(secretsCopy, secret1)

	// Verify internal state is unchanged
	secondCopy := q.GetSecretsCopy()

	if len(secondCopy) != 1 {
		t.Errorf("Expected internal state to have 1 secret, got '%d'", len(secondCopy))
	}

	if url, ok := secondCopy[secret1]; !ok {
		t.Errorf("Expected secret '%q' to still exist in internal state", secret1)
	} else if url != url1 {
		t.Errorf("Expected URL '%q', got '%q'", url1, url)
	}

	if _, exists := secondCopy["modified_secret"]; exists {
		t.Error("Expected modified_secret to not exist in internal state")
	}
}

// TestConcurrentSecretAddition tests that multiple workers can safely and concurrently add secrets
func TestConcurrentSecretAddition(t *testing.T) {
	nWorkers := 2
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := NewQueue(nWorkers)

	// Define secrets to be added
	secrets := map[string]string{
		"secret_1": "https://example.com/source1",
		"secret_2": "https://example.com/source2",
	}

	// Channel to synchronize worker startup
	startSignal := make(chan struct{})
	// Channel to signal all workers completed
	done := make(chan struct{}, nWorkers)

	// Spawn workers to add secrets concurrently
	workerIdx := 0
	for secret, url := range secrets {
		go func(s, u string) {
			<-startSignal // Wait for signal
			q.AddSecret(s, u)
			done <- struct{}{}
		}(secret, url)
		workerIdx++
	}

	// Give all goroutines time to reach the signal
	time.Sleep(50 * time.Millisecond)

	// Close the channel to trigger all workers simultaneously
	close(startSignal)

	// Wait for all workers to complete
	for i := 0; i < nWorkers; i++ {
		<-done
	}

	// Verify both secrets were stored correctly
	storedSecrets := q.GetSecretsCopy()

	if len(storedSecrets) != 2 {
		t.Errorf("Expected 2 secrets, got %d", len(storedSecrets))
	}

	for secret, expectedUrl := range secrets {
		if storedUrl, ok := storedSecrets[secret]; !ok {
			t.Errorf("Expected secret %q to be stored", secret)
		} else if storedUrl != expectedUrl {
			t.Errorf("Expected URL %q for secret %q, got %q", expectedUrl, secret, storedUrl)
		}
	}
}

// TestTakeUrl_ConcurrentConsumption tests concurrent consumption of 10 known URLs from the queue
func TestTakeUrl_ConcurrentConsumption(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := NewQueue(10)

	// Define 10 known URLs
	testUrls := make([]*url.URL, 10)
	for i := 0; i < 10; i++ {
		testUrls[i], _ = url.Parse("https://example.com/path" + strconv.Itoa(i))
	}

	// Add 10 URLs to the queue one by one
	for _, testUrl := range testUrls {
		q.AddUrl(ctx, testUrl)
	}

	// Give AddUrl goroutines time to queue the URLs
	time.Sleep(50 * time.Millisecond)

	// Semaphore channel to synchronize consumer goroutines
	startSignal := make(chan struct{})
	// Channel to collect consumed results
	results := make(chan *url.URL, 10)

	// Spawn 10 consumer goroutines that will take URLs simultaneously
	for i := 0; i < 10; i++ {
		go func() {
			<-startSignal // Wait for signal
			u, _ := q.TakeUrl(ctx)
			results <- u
		}()
	}

	// Give all goroutines time to reach the signal
	time.Sleep(50 * time.Millisecond)

	// Close the channel to trigger all consumers simultaneously
	close(startSignal)

	// Collect all consumed URLs
	consumedUrls := make(map[string]bool)
	for i := 0; i < 10; i++ {
		url := <-results
		if url == nil {
			t.Errorf("Consumer %d got nil URL", i)
			continue
		}
		consumedUrls[url.String()] = true
	}

	// Verify we consumed all unique URLs
	if len(consumedUrls) != 10 {
		t.Errorf("Expected 10 unique URLs consumed, got %d", len(consumedUrls))
	}

	// Verify all consumed URLs match the original ones
	for _, testUrl := range testUrls {
		if !consumedUrls[testUrl.String()] {
			t.Errorf("Expected URL %q to be consumed", testUrl.String())
		}
	}
}
