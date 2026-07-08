package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

func displayCompact(result fetchResult) {
	hasOutput := false

	if result.ReviewReqErr == nil && len(result.ReviewRequests) > 0 {
		fmt.Printf("\n %s\n\n", colored(bold, fmt.Sprintf("Pending Reviews (%d)", len(result.ReviewRequests))))
		for _, rr := range result.ReviewRequests {
			displayReviewRequest(rr)
		}
		hasOutput = true
	}

	if result.MyPRsErr == nil && len(result.MyPRs) > 0 {
		if hasOutput {
			fmt.Println()
		}
		fmt.Printf("\n %s\n\n", colored(bold, fmt.Sprintf("My PRs (%d)", len(result.MyPRs))))
		for _, pr := range result.MyPRs {
			displayMyPR(pr)
		}
		hasOutput = true
	}

	if !hasOutput && result.MyPRsErr == nil && result.ReviewReqErr == nil {
		fmt.Printf("\n %s\n\n", colored(dim, "No open PRs or pending reviews."))
	}

	if result.MyPRsErr != nil {
		fmt.Fprintf(os.Stderr, " %s %s\n", colored(yellow, "⚠"), result.MyPRsErr)
	}
	if result.ReviewReqErr != nil {
		fmt.Fprintf(os.Stderr, " %s %s\n", colored(yellow, "⚠"), result.ReviewReqErr)
	}

	fmt.Println()
}

func displayReviewRequest(rr ReviewRequest) {
	draft := ""
	if rr.IsDraft {
		draft = colored(dim, " [draft]")
	}
	ref := fmt.Sprintf("%s#%d", repoName(rr.Repo), rr.Number)
	link := hyperlink(rr.URL, colored(cyan, ref))
	fmt.Printf(" %s  %s  %s  %s  %s%s\n",
		colored(yellow, "●"),
		link,
		truncate(rr.Title, 50),
		colored(dim, "@"+rr.Author),
		relativeTime(rr.UpdatedAt),
		draft,
	)
}

func displayMyPR(pr PR) {
	draft := ""
	if pr.IsDraft {
		draft = colored(dim, " [draft]")
	}
	ref := fmt.Sprintf("%s#%d", repoName(pr.Repo), pr.Number)
	link := hyperlink(pr.URL, colored(cyan, ref))
	fmt.Printf(" %s %s %s  %s  %s  %s%s\n",
		ciSymbol(pr.CI),
		reviewSymbol(pr.Review),
		mergeSymbol(pr.Merge),
		link,
		truncate(pr.Title, 45),
		relativeTime(pr.UpdatedAt),
		draft,
	)
}

func displayTable(result fetchResult) {
	if result.ReviewReqErr == nil && len(result.ReviewRequests) > 0 {
		fmt.Printf("\n %s\n\n", colored(bold, fmt.Sprintf("Pending Reviews (%d)", len(result.ReviewRequests))))
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, " %s\t%s\t%s\t%s\t%s\n",
			colored(dim, "REPO"),
			colored(dim, "#"),
			colored(dim, "TITLE"),
			colored(dim, "AUTHOR"),
			colored(dim, "AGE"),
		)
		for _, rr := range result.ReviewRequests {
			fmt.Fprintf(w, " %s\t%d\t%s\t@%s\t%s\n",
				repoName(rr.Repo),
				rr.Number,
				truncate(rr.Title, 40),
				rr.Author,
				relativeTime(rr.UpdatedAt),
			)
		}
		w.Flush()
	}

	if result.MyPRsErr == nil && len(result.MyPRs) > 0 {
		fmt.Printf("\n %s\n\n", colored(bold, fmt.Sprintf("My PRs (%d)", len(result.MyPRs))))
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, " %s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			colored(dim, "REPO"),
			colored(dim, "#"),
			colored(dim, "TITLE"),
			colored(dim, "CI"),
			colored(dim, "REVIEW"),
			colored(dim, "MERGE"),
			colored(dim, "AGE"),
		)
		for _, pr := range result.MyPRs {
			fmt.Fprintf(w, " %s\t%d\t%s\t%s\t%s\t%s\t%s\n",
				repoName(pr.Repo),
				pr.Number,
				truncate(pr.Title, 35),
				ciLabel(pr.CI),
				reviewLabel(pr.Review),
				mergeLabel(pr.Merge),
				relativeTime(pr.UpdatedAt),
			)
		}
		w.Flush()
	}

	if result.MyPRsErr == nil && result.ReviewReqErr == nil &&
		len(result.MyPRs) == 0 && len(result.ReviewRequests) == 0 {
		fmt.Printf("\n %s\n", colored(dim, "No open PRs or pending reviews."))
	}

	if result.MyPRsErr != nil {
		fmt.Fprintf(os.Stderr, " %s %s\n", colored(yellow, "⚠"), result.MyPRsErr)
	}
	if result.ReviewReqErr != nil {
		fmt.Fprintf(os.Stderr, " %s %s\n", colored(yellow, "⚠"), result.ReviewReqErr)
	}

	fmt.Println()
}

func repoName(repo string) string {
	parts := strings.SplitN(repo, "/", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return repo
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func relativeTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return colored(dim, "now")
	case d < time.Hour:
		return colored(dim, fmt.Sprintf("%dm ago", int(d.Minutes())))
	case d < 24*time.Hour:
		return colored(dim, fmt.Sprintf("%dh ago", int(d.Hours())))
	case d < 30*24*time.Hour:
		return colored(dim, fmt.Sprintf("%dd ago", int(d.Hours()/24)))
	default:
		return colored(dim, fmt.Sprintf("%dw ago", int(d.Hours()/(24*7))))
	}
}
