// Package validator provides functions for validating local destination paths
// and parsing/validating remote repository URLs.
package validator

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var listRemoteRefs = func(repoURL string) (map[string]struct{}, error) {
	cmd := exec.Command("git", "ls-remote", "--heads", "--tags", repoURL)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	refs := make(map[string]struct{})
	for line := range bytes.SplitSeq(bytes.TrimSpace(output), []byte{'\n'}) {
		fields := strings.Fields(string(line))
		if len(fields) != 2 {
			continue
		}

		ref := fields[1]
		ref = strings.TrimPrefix(ref, "refs/heads/")
		ref = strings.TrimPrefix(ref, "refs/tags/")
		refs[ref] = struct{}{}
	}

	return refs, nil
}

// Destination checks if the given path's parent directory exists and is accessible.
// It returns nil if the path exists or if its parent exists, and an error otherwise.
// This ensures that the downloader can safely write to the target location.
func Destination(path string) error {
	cleanPath := filepath.Clean(path)

	_, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			parent := filepath.Dir(cleanPath)
			// If parent is "." and cleanPath is just a filename, parent exists.
			if parent == "." {
				return nil
			}
			if _, err := os.Stat(parent); err != nil {
				return fmt.Errorf("destination parent directory does not exist: %s", parent)
			}
			return nil
		}
		return err
	}

	return nil
}

// URL performs a quick validation to determine if the raw string matches a
// supported repository URL format. It returns an error if the format is invalid.
func URL(rawURL string) error {
	_, _, _, _, _, err := Parse(rawURL)
	return err
}

// Parse extracts host, user, repo, branch, and subpath from a clean URL string.
// It supports common Git hosting formats (e.g., github.com/user/repo/path/to/item).
//
// The returned host is the domain (e.g., github.com).
// The user is the account or organization name.
// The repo is the repository name.
// The branch is the specified branch, if present in a tree or blob URL.
// The subpath is the relative path within the repository; it will be empty if the
// URL points to the repository root.
//
// An error is returned if the URL is malformed or missing required components.
func Parse(rawURL string) (host, user, repo, branch, subpath string, err error) {
	// Prepend https:// if no scheme is provided
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	// Remove trailing slashes for consistency
	rawURL = strings.TrimSuffix(rawURL, "/")

	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to parse URL: %w", err)
	}

	if u.Host == "" {
		return "", "", "", "", "", fmt.Errorf("URL must have a host (e.g., github.com)")
	}
	host = u.Host

	// Split path parts
	pathParts := strings.FieldsFunc(u.Path, func(r rune) bool { return r == '/' })
	if len(pathParts) < 2 {
		return "", "", "", "", "", fmt.Errorf("URL must include at least a user and a repository (e.g., github.com/user/repo)")
	}

	user = pathParts[0]
	repo = pathParts[1]
	if len(pathParts) > 2 {
		if (pathParts[2] == "tree" || pathParts[2] == "blob") && len(pathParts) > 3 {
			branch = pathParts[3]
			if len(pathParts) > 4 {
				subpath = strings.Join(pathParts[4:], "/")
				repoURL := fmt.Sprintf("https://%s/%s/%s", host, user, repo)
				if resolvedBranch, resolvedSubpath, resolveErr := resolveRefAndSubpath(repoURL, pathParts[3:]); resolveErr == nil {
					branch = resolvedBranch
					subpath = resolvedSubpath
				}
			}
		} else {
			subpath = strings.Join(pathParts[2:], "/")
		}
	}

	return host, user, repo, branch, subpath, nil
}

func resolveRefAndSubpath(repoURL string, parts []string) (branch, subpath string, err error) {
	refs, err := listRemoteRefs(repoURL)
	if err != nil {
		return "", "", err
	}

	for i := len(parts); i > 0; i-- {
		candidate := strings.Join(parts[:i], "/")
		if _, ok := refs[candidate]; ok {
			return candidate, strings.Join(parts[i:], "/"), nil
		}
	}

	return "", "", fmt.Errorf("no matching branch or tag found")
}
