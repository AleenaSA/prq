package main

import (
	"encoding/json"
	"testing"
)

func TestConvertPR(t *testing.T) {
	raw := `{
		"number": 42,
		"title": "Fix auth flow",
		"url": "https://github.com/org/repo/pull/42",
		"isDraft": false,
		"mergeable": "CONFLICTING",
		"reviewDecision": "APPROVED",
		"baseRepository": {"nameWithOwner": "org/repo"},
		"author": {"login": "testuser"},
		"createdAt": "2024-01-01T00:00:00Z",
		"updatedAt": "2024-01-02T00:00:00Z",
		"additions": 50,
		"deletions": 10,
		"commits": {
			"nodes": [{
				"commit": {
					"statusCheckRollup": {"state": "FAILURE"}
				}
			}]
		}
	}`

	var node apiPullRequest
	if err := json.Unmarshal([]byte(raw), &node); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	pr := convertPR(node)

	if pr.Repo != "org/repo" {
		t.Errorf("Repo = %q, want %q", pr.Repo, "org/repo")
	}
	if pr.Number != 42 {
		t.Errorf("Number = %d, want 42", pr.Number)
	}
	if pr.CI != CIFailure {
		t.Errorf("CI = %d, want CIFailure", pr.CI)
	}
	if pr.Review != ReviewApproved {
		t.Errorf("Review = %d, want ReviewApproved", pr.Review)
	}
	if pr.Merge != Conflicting {
		t.Errorf("Merge = %d, want Conflicting", pr.Merge)
	}
	if pr.Additions != 50 {
		t.Errorf("Additions = %d, want 50", pr.Additions)
	}
}

func TestConvertPRNoChecks(t *testing.T) {
	raw := `{
		"number": 1,
		"title": "WIP",
		"url": "https://github.com/org/repo/pull/1",
		"isDraft": true,
		"mergeable": "UNKNOWN",
		"reviewDecision": "",
		"baseRepository": {"nameWithOwner": "org/repo"},
		"author": {"login": "dev"},
		"createdAt": "2024-01-01T00:00:00Z",
		"updatedAt": "2024-01-01T00:00:00Z",
		"additions": 0,
		"deletions": 0,
		"commits": {
			"nodes": [{
				"commit": {
					"statusCheckRollup": null
				}
			}]
		}
	}`

	var node apiPullRequest
	if err := json.Unmarshal([]byte(raw), &node); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	pr := convertPR(node)

	if pr.CI != CINone {
		t.Errorf("CI = %d, want CINone", pr.CI)
	}
	if pr.IsDraft != true {
		t.Errorf("IsDraft = %v, want true", pr.IsDraft)
	}
	if pr.Merge != MergeUnknown {
		t.Errorf("Merge = %d, want MergeUnknown", pr.Merge)
	}
}

func TestConvertReviewRequest(t *testing.T) {
	raw := `{
		"number": 99,
		"title": "Big feature",
		"url": "https://github.com/org/repo/pull/99",
		"isDraft": false,
		"author": {"login": "colleague"},
		"baseRepository": {"nameWithOwner": "org/repo"},
		"createdAt": "2024-01-01T00:00:00Z",
		"updatedAt": "2024-01-01T12:00:00Z",
		"additions": 200,
		"deletions": 50
	}`

	var node apiSearchNode
	if err := json.Unmarshal([]byte(raw), &node); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	rr := convertReviewRequest(node)

	if rr.Author != "colleague" {
		t.Errorf("Author = %q, want %q", rr.Author, "colleague")
	}
	if rr.Additions != 200 {
		t.Errorf("Additions = %d, want 200", rr.Additions)
	}
}
