// Package git provides git helper functions
package git

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// GetChangeIDFromCommit gets Change-Id from current git commit
func GetChangeIDFromCommit(cwd string) (string, error) {
	cmd := exec.Command("git", "log", "-1", "--format=%B")
	cmd.Dir = cwd

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git commit message: %w", err)
	}

	message := string(output)
	re := regexp.MustCompile(`Change-Id: (I[a-f0-9]{40})`)
	match := re.FindStringSubmatch(message)

	if match == nil {
		return "", nil
	}

	return match[1], nil
}

// GetHostFromRemote extracts host from git remote URL
func GetHostFromRemote(cwd string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = cwd

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git remote URL: %w", err)
	}

	remoteURL := strings.TrimSpace(string(output))

	// Parse SSH URL: ssh://user@host:port/project or git@host:project
	// Parse HTTPS URL: https://host/project
	var host string

	if strings.HasPrefix(remoteURL, "ssh://") {
		re := regexp.MustCompile(`ssh://(?:[^@]+@)?([^:/]+)`)
		if match := re.FindStringSubmatch(remoteURL); match != nil {
			host = match[1]
		}
	} else if strings.HasPrefix(remoteURL, "https://") {
		re := regexp.MustCompile(`https://([^/]+)`)
		if match := re.FindStringSubmatch(remoteURL); match != nil {
			host = match[1]
		}
	} else if strings.Contains(remoteURL, "@") {
		// git@host:project format
		re := regexp.MustCompile(`[^@]+@([^:]+):`)
		if match := re.FindStringSubmatch(remoteURL); match != nil {
			host = match[1]
		}
	}

	if host == "" {
		return "", fmt.Errorf("could not parse host from git remote URL: %s", remoteURL)
	}

	return host, nil
}
