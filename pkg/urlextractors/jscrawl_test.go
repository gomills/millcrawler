package urlsextractors

import (
	"context"
	"io"
	"slices"
	"strings"
	"testing"
	"time"

	http "github.com/bogdanfinn/fhttp"
	"github.com/gomills/millcrawler/pkg/config"
	"github.com/gomills/millcrawler/pkg/queue"
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

type jsTestCase struct {
	name          string
	jsCode        string
	expectedCount int
}

// These are the tested parent types as of now:
// call_expression
// import_statement
// assignment_expression
// variable_declarator
// pair
// assignment_pattern
// arguments
// binary_expression

func runJsTest(t *testing.T, tc jsTestCase) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(tc.jsCode)),
	}

	initialCount := q.GetNFoundUrls()
	CrawlJs(ctx, jsTestConfig, response.Body, q, jsDomain, jsRegisteredDomain)

	time.Sleep(10 * time.Millisecond)

	finalCount := q.GetNFoundUrls()
	extractedCount := finalCount - initialCount

	if extractedCount != tc.expectedCount {
		t.Errorf("%s: expected %d URLs, got %d (initial: %d, final: %d)",
			tc.name, tc.expectedCount, extractedCount, initialCount, finalCount)
	}
}

func TestCrawlJsParentTypes(t *testing.T) {
	tests := []jsTestCase{
		{
			name:          "CallExpression",
			jsCode:        `fetch('//example.com/api');`,
			expectedCount: 1,
		},
		{
			name:          "ImportStatement/Object",
			jsCode:        `const config = { apiUrl: 'https://example.com/api' };`,
			expectedCount: 1,
		},
		{
			name:          "AssignmentExpression",
			jsCode:        `const apiUrl = 'https://example.com/api';`,
			expectedCount: 1,
		},
		{
			name:          "VariableDeclarator",
			jsCode:        `let endpoint = 'https://example.com/endpoint';`,
			expectedCount: 1,
		},
		{
			name:          "PairType",
			jsCode:        `const settings = { endpoint: 'https://example.com/api' };`,
			expectedCount: 1,
		},
		{
			name:          "AssignmentPatternType",
			jsCode:        `function connect(url = 'https://example.com/default') { return url; }`,
			expectedCount: 1,
		},
		{
			name:          "ArgumentsType",
			jsCode:        `fetch('//example.com/resource');`,
			expectedCount: 1,
		},
		{
			name:          "BinaryExpressionType",
			jsCode:        `const url = 'https://example.com' + '/api';`,
			expectedCount: 2,
		},
		{
			name: "CombinedParentTypes",
			jsCode: `
const config = { apiUrl: 'https://example.com/api' };
const baseUrl = 'https://example.com/base';
let endpoint = 'https://example.com/endpoint';
function connect(url = 'https://example.com/default') { return url; }
fetch('//example.com/resource');
const fullUrl = 'https://example.com' + '/api';
`,
			expectedCount: 7,
		},
		{
			name: "NoValidUrls",
			jsCode: `
const greeting = 'hello world';
const version = '1.0.0';
fetch('some local function');
`,
			expectedCount: 0,
		},
		{
			name:          "EmptyFile",
			jsCode:        `// Empty file with just comments`,
			expectedCount: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			runJsTest(t, tc)
		})
	}
}

func TestCrawlJsExactUrlExtraction(t *testing.T) {
	ctx := context.Background()
	q := queue.NewQueue(1)

	jsCode := `
const config = { apiUrl: 'https://www.example.com/api' };
const baseUrl = 'https://www.example.com/base';
let endpoint = 'https://www.example.com/endpoint';
function connect(url = 'https://www.example.com/default') { return url; }
fetch('//www.example.com/resource');
const fullUrl = 'https://example.com/se66' + '/api';
const fullUrl = 'https://example.com/se66' + '/2api';
`

	response := &http.Response{
		Body: io.NopCloser(strings.NewReader(jsCode)),
	}

	CrawlJs(ctx, jsTestConfig, response.Body, q, jsDomain, jsRegisteredDomain)
	time.Sleep(10 * time.Millisecond)

	extractedUrls := q.GetFoundUrls()

	// Expected URLs after sanitization (relative to absolute)
	expectedUrls := []string{
		"https://www.example.com/api",
		"https://www.example.com/base",
		"https://www.example.com/endpoint",
		"https://www.example.com/default",
		"https://www.example.com/resource",
		"https://example.com/se66",
		"https://www.example.com/2api",
	}

	// Sort for consistent comparison
	slices.Sort(extractedUrls)
	slices.Sort(expectedUrls)

	if len(extractedUrls) != len(expectedUrls) {
		t.Logf("Extracted URLs (%d): %v", len(extractedUrls), extractedUrls)
		t.Logf("Expected URLs (%d): %v", len(expectedUrls), expectedUrls)
		t.Fatalf("Expected %d URLs, got %d", len(expectedUrls), len(extractedUrls))
	}

	for i, expected := range expectedUrls {
		if extractedUrls[i] != expected {
			t.Errorf("URL mismatch at index %d: expected %q, got %q", i, expected, extractedUrls[i])
		}
	}
}
