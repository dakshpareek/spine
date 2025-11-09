package cmd

import (
	"testing"
	"time"
)

func TestHumanizeDuration(t *testing.T) {
	tests := []struct {
		input    time.Duration
		expected string
	}{
		{10 * time.Second, "just now"},
		{2 * time.Minute, "2 minutes ago"},
		{3 * time.Hour, "3 hours ago"},
		{49 * time.Hour, "2 days ago"},
	}

	for _, tc := range tests {
		if got := humanizeDuration(tc.input); got != tc.expected {
			t.Fatalf("humanizeDuration(%v) = %q, expected %q", tc.input, got, tc.expected)
		}
	}
}
