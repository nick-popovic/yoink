// Package validator provides functions for validating local destination paths
// and parsing/validating remote repository URLs.
package validator

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

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
	_, _, _, _, err := Parse(rawURL)
	return err
}

// Parse extracts host, user, repo, and subpath from a clean URL string.
// It supports common Git hosting formats (e.g., github.com/user/repo/path/to/item).
//
// The returned host is the domain (e.g., github.com).
// The user is the account or organization name.
// The repo is the repository name.
// The subpath is the relative path within the repository; it will be empty if the
// URL points to the repository root.
//
// An error is returned if the URL is malformed or missing required components.
func Parse(rawURL string) (host, user, repo, subpath string, err error) {
	// Prepend https:// if no scheme is provided
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	// Remove trailing slashes for consistency
	rawURL = strings.TrimSuffix(rawURL, "/")

	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to parse URL: %w", err)
	}

	if u.Host == "" {
		return "", "", "", "", fmt.Errorf("URL must have a host (e.g., github.com)")
	}
	host = u.Host

	// Split path parts
	pathParts := strings.FieldsFunc(u.Path, func(r rune) bool { return r == '/' })
	if len(pathParts) < 2 {
		return "", "", "", "", fmt.Errorf("URL must include at least a user and a repository (e.g., github.com/user/repo)")
	}

	user = pathParts[0]
	repo = pathParts[1]
	if len(pathParts) > 2 {
		subpath = strings.Join(pathParts[2:], "/")
	}

	return host, user, repo, subpath, nil
}
