package main

import (
	"testing"
	"time"
)

func TestRelativeTime(t *testing.T) {
	colorEnabled = false

	tests := []struct {
		input time.Time
		want  string
	}{
		{time.Now().Add(-30 * time.Second), "now"},
		{time.Now().Add(-5 * time.Minute), "5m ago"},
		{time.Now().Add(-3 * time.Hour), "3h ago"},
		{time.Now().Add(-2 * 24 * time.Hour), "2d ago"},
		{time.Now().Add(-10 * 7 * 24 * time.Hour), "10w ago"},
		{time.Time{}, ""},
	}

	for _, tt := range tests {
		got := relativeTime(tt.input)
		if got != tt.want {
			t.Errorf("relativeTime(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input string
		max   int
		want  string
	}{
		{"hello", 10, "hello"},
		{"hello world", 5, "hell…"},
		{"abc", 3, "abc"},
	}

	for _, tt := range tests {
		got := truncate(tt.input, tt.max)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.max, got, tt.want)
		}
	}
}

func TestRepoName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"org/repo", "repo"},
		{"very-long-org-name/my-service", "my-service"},
		{"short/repo-name", "repo-name"},
		{"solo", "solo"},
	}

	for _, tt := range tests {
		got := repoName(tt.input)
		if got != tt.want {
			t.Errorf("repoName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
