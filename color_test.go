package main

import "testing"

func TestCISymbol(t *testing.T) {
	colorEnabled = false

	tests := []struct {
		status CIStatus
		want   string
	}{
		{CISuccess, "✓"},
		{CIFailure, "✗"},
		{CIPending, "◌"},
		{CIError, "!"},
		{CINone, "–"},
	}

	for _, tt := range tests {
		got := ciSymbol(tt.status)
		if got != tt.want {
			t.Errorf("ciSymbol(%d) = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestReviewSymbol(t *testing.T) {
	colorEnabled = false

	tests := []struct {
		status ReviewStatus
		want   string
	}{
		{ReviewApproved, "●"},
		{ReviewChangesRequested, "●"},
		{ReviewRequired, "●"},
		{ReviewNone, "○"},
	}

	for _, tt := range tests {
		got := reviewSymbol(tt.status)
		if got != tt.want {
			t.Errorf("reviewSymbol(%d) = %q, want %q", tt.status, got, tt.want)
		}
	}
}

func TestMergeSymbol(t *testing.T) {
	colorEnabled = false

	tests := []struct {
		status MergeStatus
		want   string
	}{
		{Mergeable, "◯"},
		{Conflicting, "⊘"},
		{MergeUnknown, "?"},
	}

	for _, tt := range tests {
		got := mergeSymbol(tt.status)
		if got != tt.want {
			t.Errorf("mergeSymbol(%d) = %q, want %q", tt.status, got, tt.want)
		}
	}
}
