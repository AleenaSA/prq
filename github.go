package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type graphQLRequest struct {
	Query string `json:"query"`
}

type graphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type myPRsResponse struct {
	Viewer struct {
		Login        string `json:"login"`
		PullRequests struct {
			Nodes []apiPullRequest `json:"nodes"`
		} `json:"pullRequests"`
	} `json:"viewer"`
}

type apiPullRequest struct {
	Number         int    `json:"number"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	IsDraft        bool   `json:"isDraft"`
	Mergeable      string `json:"mergeable"`
	ReviewDecision string `json:"reviewDecision"`
	BaseRepository struct {
		NameWithOwner string `json:"nameWithOwner"`
	} `json:"baseRepository"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Commits   struct {
		Nodes []struct {
			Commit struct {
				StatusCheckRollup *struct {
					State string `json:"state"`
				} `json:"statusCheckRollup"`
			} `json:"commit"`
		} `json:"nodes"`
	} `json:"commits"`
}

type searchResponse struct {
	Search struct {
		Nodes []apiSearchNode `json:"nodes"`
	} `json:"search"`
}

type apiSearchNode struct {
	Number         int    `json:"number"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	IsDraft        bool   `json:"isDraft"`
	Author         *struct {
		Login string `json:"login"`
	} `json:"author"`
	BaseRepository *struct {
		NameWithOwner string `json:"nameWithOwner"`
	} `json:"baseRepository"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

type fetchResult struct {
	MyPRs          []PR
	ReviewRequests []ReviewRequest
	MyPRsErr       error
	ReviewReqErr   error
}

func fetchAll(token string) fetchResult {
	var result fetchResult
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		result.MyPRs, result.MyPRsErr = fetchMyPRs(token)
	}()

	go func() {
		defer wg.Done()
		result.ReviewRequests, result.ReviewReqErr = fetchReviewRequests(token)
	}()

	wg.Wait()
	return result
}

func fetchMyPRs(token string) ([]PR, error) {
	raw, err := executeQuery(token, myPRsQuery)
	if err != nil {
		return nil, err
	}

	var resp myPRsResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	prs := make([]PR, 0, len(resp.Viewer.PullRequests.Nodes))
	for _, node := range resp.Viewer.PullRequests.Nodes {
		prs = append(prs, convertPR(node))
	}
	return prs, nil
}

func fetchReviewRequests(token string) ([]ReviewRequest, error) {
	raw, err := executeQuery(token, reviewRequestsQuery)
	if err != nil {
		return nil, err
	}

	var resp searchResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	reqs := make([]ReviewRequest, 0, len(resp.Search.Nodes))
	for _, node := range resp.Search.Nodes {
		if node.Number == 0 {
			continue
		}
		reqs = append(reqs, convertReviewRequest(node))
	}
	return reqs, nil
}

func executeQuery(token, query string) (json.RawMessage, error) {
	body, err := json.Marshal(graphQLRequest{Query: query})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GitHub API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(respBody))
	}

	var gqlResp graphQLResponse
	if err := json.Unmarshal(respBody, &gqlResp); err != nil {
		return nil, fmt.Errorf("parsing GraphQL response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", gqlResp.Errors[0].Message)
	}

	return gqlResp.Data, nil
}

func convertPR(node apiPullRequest) PR {
	pr := PR{
		Number:    node.Number,
		Title:     node.Title,
		URL:       node.URL,
		Author:    node.Author.Login,
		IsDraft:   node.IsDraft,
		Additions: node.Additions,
		Deletions: node.Deletions,
	}

	if node.BaseRepository.NameWithOwner != "" {
		pr.Repo = node.BaseRepository.NameWithOwner
	}

	pr.CreatedAt, _ = time.Parse(time.RFC3339, node.CreatedAt)
	pr.UpdatedAt, _ = time.Parse(time.RFC3339, node.UpdatedAt)

	switch node.Mergeable {
	case "MERGEABLE":
		pr.Merge = Mergeable
	case "CONFLICTING":
		pr.Merge = Conflicting
	default:
		pr.Merge = MergeUnknown
	}

	switch node.ReviewDecision {
	case "APPROVED":
		pr.Review = ReviewApproved
	case "CHANGES_REQUESTED":
		pr.Review = ReviewChangesRequested
	case "REVIEW_REQUIRED":
		pr.Review = ReviewRequired
	default:
		pr.Review = ReviewNone
	}

	if len(node.Commits.Nodes) > 0 {
		rollup := node.Commits.Nodes[0].Commit.StatusCheckRollup
		if rollup == nil {
			pr.CI = CINone
		} else {
			switch rollup.State {
			case "SUCCESS":
				pr.CI = CISuccess
			case "FAILURE":
				pr.CI = CIFailure
			case "PENDING":
				pr.CI = CIPending
			case "ERROR":
				pr.CI = CIError
			default:
				pr.CI = CINone
			}
		}
	} else {
		pr.CI = CINone
	}

	return pr
}

func convertReviewRequest(node apiSearchNode) ReviewRequest {
	rr := ReviewRequest{
		Number:    node.Number,
		Title:     node.Title,
		URL:       node.URL,
		IsDraft:   node.IsDraft,
		Additions: node.Additions,
		Deletions: node.Deletions,
	}

	if node.Author != nil {
		rr.Author = node.Author.Login
	}
	if node.BaseRepository != nil {
		rr.Repo = node.BaseRepository.NameWithOwner
	}

	rr.CreatedAt, _ = time.Parse(time.RFC3339, node.CreatedAt)
	rr.UpdatedAt, _ = time.Parse(time.RFC3339, node.UpdatedAt)

	return rr
}
