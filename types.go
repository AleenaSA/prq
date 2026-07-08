package main

import "time"

type CIStatus int

const (
	CISuccess CIStatus = iota
	CIFailure
	CIPending
	CIError
	CINone
)

type ReviewStatus int

const (
	ReviewApproved ReviewStatus = iota
	ReviewChangesRequested
	ReviewRequired
	ReviewNone
)

type MergeStatus int

const (
	Mergeable MergeStatus = iota
	Conflicting
	MergeUnknown
)

type PR struct {
	Repo      string
	Number    int
	Title     string
	URL       string
	Author    string
	IsDraft   bool
	CI        CIStatus
	Review    ReviewStatus
	Merge     MergeStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ReviewRequest struct {
	Repo      string
	Number    int
	Title     string
	URL       string
	Author    string
	IsDraft   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
