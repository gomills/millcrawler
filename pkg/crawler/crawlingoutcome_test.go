package crawler

import "testing"

// TestGetCrawlingOutcome_CorrectStorage verifies that GetCrawlingOutcome correctly stores all fields and secrets without data corruption
func TestGetCrawlingOutcome_CorrectStorage(t *testing.T) {

	// dummy data
	domain := "example.com"
	numURLs := 1893
	durationSeconds := 131.13
	stopReason := "empty_queue"
	secretsFound := 4

	secretsMap := map[string]string{
		"aws_key_abc123":      "https://example.com/config.json",
		"api_token_xyz":       "https://example.com/env",
		"github_token_def456": "https://example.com/secrets.yaml",
		"_key_ghi789":         "https://example.com/api/path1/path2/keys",
	}

	// craft outcome
	outcome := CreateCrawlingOutcome(domain, numURLs, durationSeconds, stopReason, secretsFound, secretsMap)

	// Verify all metrics are stored correctly and exactly as they were inputted
	if outcome.Domain != domain {
		t.Errorf("Expected Domain %q, got %q", domain, outcome.Domain)
	}

	if outcome.NumURLs != numURLs {
		t.Errorf("Expected NumURLs %d, got %d", numURLs, outcome.NumURLs)
	}

	if outcome.DurationSeconds != durationSeconds {
		t.Errorf("Expected DurationSeconds %f, got %f", durationSeconds, outcome.DurationSeconds)
	}

	if outcome.StopReason != stopReason {
		t.Errorf("Expected StopReason %q, got %q", stopReason, outcome.StopReason)
	}

	if outcome.SecretsFound != secretsFound {
		t.Errorf("Expected SecretsFound %d, got %d", secretsFound, outcome.SecretsFound)
	}

	// Verify all 4 secrets are stored correctly
	if len(outcome.SecretsMap) != 4 {
		t.Errorf("Expected 4 secrets, got %d", len(outcome.SecretsMap))
	}

	for secret, expectedURL := range secretsMap {
		if storedURL, ok := outcome.SecretsMap[secret]; !ok {
			t.Errorf("Expected secret %q to be stored", secret)
		} else if storedURL != expectedURL {
			t.Errorf("Expected URL %q for secret %q, got %q", expectedURL, secret, storedURL)
		}
	}
}
