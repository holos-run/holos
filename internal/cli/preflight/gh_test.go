package preflight

import "testing"

func TestGhTokenAllowsRepoCreation(t *testing.T) {
	testCases := []struct {
		name     string
		status   ghAuthStatusResponse
		expected bool
	}{
		{
			name:     "token has necessary scopes",
			status:   "- Token: gho_************************************\n- Token scopes: 'gist', 'read:org', 'repo'",
			expected: true,
		},
		{
			name:     "token has necessary scopes",
			status:   "  - Token scopes: 'gist', 'read:org', 'repo'",
			expected: true,
		},
		{
			name:     "token does not have necessary scopes",
			status:   "- Token: gho_************************************\n- Token scopes: 'gist', 'read:org'",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ghTokenAllowsRepoCreation(tc.status)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
