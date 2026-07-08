package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func resolveToken() (string, error) {
	if token := ghAuthToken(); token != "" {
		return token, nil
	}
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}
	return "", fmt.Errorf("no GitHub token found.\n\n  Option 1: Install gh and run `gh auth login`\n  Option 2: Set the GITHUB_TOKEN environment variable\n")
}

func ghAuthToken() string {
	cmd := exec.Command("gh", "auth", "token")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
