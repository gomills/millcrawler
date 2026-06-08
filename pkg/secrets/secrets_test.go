package secrets

import (
	"testing"
)

func TestFindSecrets(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name: "finds key2",
			content: `
Some text here and then the secret key:
mywebsecret_AbCdEfGhIjKlMnOpQrSt
More text after the key.
`,
			expected: []string{"mywebsecret_AbCdEfGhIjKlMnOpQrSt"},
		},
		{
			name: "finds key3",
			content: `
Configuration file with:
mywebsecret_XyZaBcDeFgHiJkLmNoPqRs
embedded in the middle.
`,
			expected: []string{"mywebsecret_XyZaBcDeFgHiJkLmNoPqRs"},
		},
		{
			name: "finds all three keys",
			content: `
Database config:
mywebsecret_AbCdEfGhIjKlMnOpQrSt

API settings:
mywebsecret_XyZaBcDeFgHiJkLmNoPqRs
mywebsecret_1234567890AbCdEfGhIj

End of file.
`,
			expected: []string{
				"mywebsecret_AbCdEfGhIjKlMnOpQrSt",
				"mywebsecret_XyZaBcDeFgHiJkLmNoPqRs",
				"mywebsecret_1234567890AbCdEfGhIj",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FindSecrets(tt.content)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("expected %d secrets, got %d", len(tt.expected), len(result))
			}

			for i, exp := range tt.expected {
				if i >= len(result) {
					t.Errorf("missing expected secret: %s", exp)
					continue
				}
				if result[i] != exp {
					t.Errorf("secret %d: expected %s, got %s", i, exp, result[i])
				}
			}
		})
	}
}
