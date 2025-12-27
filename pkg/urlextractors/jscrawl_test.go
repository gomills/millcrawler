package urlsextractors

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	http "github.com/bogdanfinn/fhttp"
	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
)

var jsTestConfig = &config.Config{
	TimeOutDuration:        60,
	MaxPathDepth:           1,
	AllowedExternalDomains: []string{"github.com", "cdn.example.com"},
	SensitivePatterns:      []string{"dashboard", "test", "repo", "private"},
	AllowedExtensions:      []string{".git", ".txt", ".js", ".json", ""},
	Workers:                3,
}

const (
	jsDomain           = "www.example.com"
	jsRegisteredDomain = "example.com"
)

// TestCrawlJsCallExpression tests extracting URL from a call_expression (function call)
func TestCrawlJsCallExpression(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Simple call_expression: fetch('//example.com/api')
	jsCode := `fetch('//example.com/api');`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsCode)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlJs(ctx, jsTestConfig, response, q, jsDomain, jsRegisteredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 1 // One URL in fetch call

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URL from call_expression, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlJsImportStatement tests extracting URL from inside an object (relevant parent type)
func TestCrawlJsImportStatement(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Object with URL string property
	jsCode := `const config = { apiUrl: 'https://example.com/api' };`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsCode)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlJs(ctx, jsTestConfig, response, q, jsDomain, jsRegisteredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 1 // One URL in object property

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URL from object, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlJsAssignmentExpression tests extracting URL from an assignment_expression
func TestCrawlJsAssignmentExpression(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Simple assignment_expression: const apiUrl = 'https://example.com/api'
	jsCode := `const apiUrl = 'https://example.com/api';`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsCode)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlJs(ctx, jsTestConfig, response, q, jsDomain, jsRegisteredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 1 // One URL in assignment

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URL from assignment_expression, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlJsVariableDeclarator tests extracting URL from a variable_declarator
func TestCrawlJsVariableDeclarator(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Simple variable_declarator: let endpoint = 'https://example.com/endpoint'
	jsCode := `let endpoint = 'https://example.com/endpoint';`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsCode)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlJs(ctx, jsTestConfig, response, q, jsDomain, jsRegisteredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 1 // One URL in variable declarator

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URL from variable_declarator, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlJsCombinedParentTypes tests extracting URLs from multiple parent types in one file
func TestCrawlJsCombinedParentTypes(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// JS code with multiple parent types: call_expression, assignment_expression, variable_declarator, object
	jsCode := `
const config = { apiUrl: 'https://example.com/api' };
const baseUrl = 'https://example.com/base';
let endpoint = 'https://example.com/endpoint';
fetch('//example.com/data');
`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsCode)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlJs(ctx, jsTestConfig, response, q, jsDomain, jsRegisteredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 4 // One from each: object property, const (assignment), let (variable), fetch (call)

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URLs from combined parent types, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlJsNoValidUrls tests JS with strings that don't pass validation
func TestCrawlJsNoValidUrls(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// JS code with non-URL strings
	jsCode := `
const greeting = 'hello world';
const version = '1.0.0';
fetch('some local function');
`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsCode)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlJs(ctx, jsTestConfig, response, q, jsDomain, jsRegisteredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 0 // No valid URLs

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URLs from invalid strings, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlJsEmptyFile tests JS with empty or whitespace-only content
func TestCrawlJsEmptyFile(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	jsCode := `// Empty file with just comments`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsCode)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlJs(ctx, jsTestConfig, response, q, jsDomain, jsRegisteredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 0 // No URLs

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URLs from empty file, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlJsPairType tests extracting URL from a pair (object property)
func TestCrawlJsPairType(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Pair in object literal: { endpoint: 'https://example.com/api' }
	jsCode := `const settings = { endpoint: 'https://example.com/api' };`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsCode)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlJs(ctx, jsTestConfig, response, q, jsDomain, jsRegisteredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 1 // One URL in pair

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URL from pair, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlJsAssignmentPatternType tests extracting URL from assignment_pattern (default parameter)
func TestCrawlJsAssignmentPatternType(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Assignment pattern in function parameter: function(url = 'https://example.com/default')
	jsCode := `function connect(url = 'https://example.com/default') { return url; }`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsCode)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlJs(ctx, jsTestConfig, response, q, jsDomain, jsRegisteredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 1 // One URL in assignment pattern

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URL from assignment_pattern, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlJsArgumentsType tests extracting URL from arguments (function call arguments)
func TestCrawlJsArgumentsType(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Arguments in function call: fetch('//example.com/resource')
	jsCode := `fetch('//example.com/resource');`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsCode)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlJs(ctx, jsTestConfig, response, q, jsDomain, jsRegisteredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 1 // One URL in arguments

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URL from arguments, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlJsBinaryExpressionType tests extracting URL from binary_expression (concatenation)
func TestCrawlJsBinaryExpressionType(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	// Binary expression with string concatenation: 'https://example.com' + '/api'
	jsCode := `const url = 'https://example.com' + '/api';`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsCode)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlJs(ctx, jsTestConfig, response, q, jsDomain, jsRegisteredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 2 // Both operands extracted (the URL and '/api')

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URLs from binary_expression, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}

// TestCrawlJsAllNewParentTypes tests all four new parent types in a single file
func TestCrawlJsAllNewParentTypes(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	jsCode := `
const settings = { endpoint: 'https://example.com/api' };
function connect(url = 'https://example.com/default') { return url; }
fetch('//example.com/resource');
const fullUrl = 'https://example.com' + '/api';
`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsCode)),
	}

	initialCount := q.GetNCrawledUrls()
	CrawlJs(ctx, jsTestConfig, response, q, jsDomain, jsRegisteredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNCrawledUrls()
	extractedCount := finalCount - initialCount
	expectedCount := 5 // pair (1) + assignment_pattern (1) + arguments (1) + binary_expression (2 operands)

	if extractedCount != expectedCount {
		t.Errorf("Expected %d URLs from all new parent types, got %d (initial: %d, final: %d)", expectedCount, extractedCount, initialCount, finalCount)
	}
}
