package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

var (
	colorEnabled = true
	isTTY        = true
)

func initColor(noColorFlag bool) {
	isTTY = term.IsTerminal(int(os.Stdout.Fd()))
	if noColorFlag || os.Getenv("NO_COLOR") != "" || !isTTY {
		colorEnabled = false
	}
}

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	dim    = "\033[2m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
)

func colored(color, text string) string {
	if !colorEnabled {
		return text
	}
	return color + text + reset
}

func ciSymbol(status CIStatus) string {
	switch status {
	case CISuccess:
		return colored(green, "✓")
	case CIFailure:
		return colored(red, "✗")
	case CIPending:
		return colored(yellow, "◌")
	case CIError:
		return colored(red, "!")
	default:
		return colored(dim, "–")
	}
}

func reviewSymbol(status ReviewStatus) string {
	switch status {
	case ReviewApproved:
		return colored(green, "●")
	case ReviewChangesRequested:
		return colored(red, "●")
	case ReviewRequired:
		return colored(yellow, "●")
	default:
		return colored(dim, "○")
	}
}

func mergeSymbol(status MergeStatus) string {
	switch status {
	case Mergeable:
		return colored(green, "◯")
	case Conflicting:
		return colored(red, "⊘")
	default:
		return colored(dim, "?")
	}
}

func ciLabel(status CIStatus) string {
	switch status {
	case CISuccess:
		return "passing"
	case CIFailure:
		return "failing"
	case CIPending:
		return "pending"
	case CIError:
		return "error"
	default:
		return "–"
	}
}

func reviewLabel(status ReviewStatus) string {
	switch status {
	case ReviewApproved:
		return "approved"
	case ReviewChangesRequested:
		return "changes"
	case ReviewRequired:
		return "review needed"
	default:
		return "–"
	}
}

func mergeLabel(status MergeStatus) string {
	switch status {
	case Mergeable:
		return "clean"
	case Conflicting:
		return "conflict"
	default:
		return "?"
	}
}

func hyperlink(url, text string) string {
	if !colorEnabled || !isTTY {
		return text
	}
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", url, text)
}
