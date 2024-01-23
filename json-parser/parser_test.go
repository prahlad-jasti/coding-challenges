package parser

import "testing"

func TestValidJSON(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{`{"key": "value", "numbers": [1, 2, 3], "nested": {"foo": "bar"}}`, true},
		{"{}", true},
		{`{"key1": "value", "key2": "value"}`, true},
		{`{"key": "value", "numbers": [1, 2, 3], "nested": {"foo": "bar"}`, false},
		{`{"key": "value", "numbers": 1, 2, 3], "nested": {"foo": "bar"}}`, false},
		{"invalid json", false},
	}
    for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := IsValidJSON(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}
