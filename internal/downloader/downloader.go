// Package downloader provides functionality for fetching repository contents
// using git sparse-checkout and partial clones to minimize download size.
package downloader

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Fetch downloads repository contents for a given subpath to a local destination.
// It uses git's sparse-checkout and partial-clone capabilities to minimize
// network traffic and disk space usage.
//
// Arguments:
//   - host: the Git host (e.g., github.com)
//   - user: the repository's account or organization name
//   - repo: the repository's name
//   - subpath: the specific sub-folder or file to fetch; if empty, clones the root
//   - destination: the local path where contents should be placed
//
// It returns an error if git is not installed, if network operations fail,
// or if the destination cannot be written to.
func Fetch(host, user, repo, subpath, destination string) error {
	// 1. Check if git is installed
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed or not in PATH")
	}

	// 2. Create a temporary directory
	tempDir, err := os.MkdirTemp("", "yoink-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	repoURL := fmt.Sprintf("https://%s/%s/%s", host, user, repo)

	// 3. Clone with sparse-checkout
	// git clone --depth 1 --filter=blob:none --sparse <repoURL> <tempDir>
	cloneCmd := exec.Command("git", "clone", "--depth", "1", "--filter=blob:none", "--sparse", repoURL, tempDir)
	if output, err := cloneCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clone repository: %w, output: %s", err, string(output))
	}

	// 4. Set sparse-checkout for the subpath if provided
	if subpath != "" {
		// Use --no-cone to allow individual files to be matched.
		// Cone mode (the default) is optimized for directories and often fails on files.
		checkoutCmd := exec.Command("git", "-C", tempDir, "sparse-checkout", "set", "--no-cone", subpath)
		if output, err := checkoutCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to set sparse-checkout: %w, output: %s", err, string(output))
		}
	} else {
		// If no subpath, we want everything
		checkoutCmd := exec.Command("git", "-C", tempDir, "sparse-checkout", "disable")
		if output, err := checkoutCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to disable sparse-checkout: %w, output: %s", err, string(output))
		}
	}

	// 5. Check if the subpath is a file or directory
	src := filepath.Join(tempDir, subpath)
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source path: %w", err)
	}

	// 6. Ensure destination directory exists if moving multiple files or if destination is clearly a directory
	if info.IsDir() {
		if err := os.MkdirAll(destination, 0755); err != nil {
			return fmt.Errorf("failed to create destination directory: %w", err)
		}
	} else {
		// If it's a file, ensure the parent directory of the destination exists
		destDir := destination
		destInfo, err := os.Stat(destination)
		if err == nil && destInfo.IsDir() {
			// destination is an existing directory, we'll put the file inside it
		} else {
			// destination is a file path or doesn't exist
			destDir = filepath.Dir(destination)
		}

		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create destination parent directory: %w", err)
		}
	}

	// 7. Move the contents or the file
	return movePath(src, destination, info.IsDir())
}

// movePath handles moving either a single file or directory contents from source
// to destination. If src is a directory, it moves all files inside it (except .git)
// into dst.
func movePath(src, dst string, isDir bool) error {
	if isDir {
		entries, err := os.ReadDir(src)
		if err != nil {
			return fmt.Errorf("failed to read source directory: %w", err)
		}

		for _, entry := range entries {
			if entry.Name() == ".git" {
				continue
			}
			oldPath := filepath.Join(src, entry.Name())
			newPath := filepath.Join(dst, entry.Name())
			if err := move(oldPath, newPath); err != nil {
				return err
			}
		}
		return nil
	}

	// For a single file
	destPath := dst
	if info, err := os.Stat(dst); err == nil && info.IsDir() {
		destPath = filepath.Join(dst, filepath.Base(src))
	}
	return move(src, destPath)
}

// move attempts to rename a file/directory from oldPath to newPath. If renaming
// fails (e.g., across different filesystems), it falls back to a recursive copy.
func move(oldPath, newPath string) error {
	// Try rename first (fastest)
	if err := os.Rename(oldPath, newPath); err != nil {
		// If rename fails (e.g. cross-device link), we should copy and remove.
		// Using 'cp -r' to handle directories if called for them
		cpCmd := exec.Command("cp", "-r", oldPath, newPath)
		if output, err := cpCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to copy %s to %s: %w, output: %s", oldPath, newPath, err, string(output))
		}
	}
	return nil
}
