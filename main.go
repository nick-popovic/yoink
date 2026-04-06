// Package main provides the entry point for the yoink CLI utility.
// yoink is a tool designed to efficiently download specific sub-directories
// or files from Git repositories without cloning the entire repository history.
//
// Usage:
//
//	yoink [destination] <url>
//
// Examples:
//
//	yoink . github.com/user/repo/path/to/dir
//	yoink my-local-dir github.com/user/repo/path/to/file.txt
package main

import (
	"fmt"
	"os"
	"yoink/internal/downloader"
	"yoink/internal/validator"
)

func main() {
	version := "1.0.0"
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--version" {
			fmt.Printf("yoink version %s\n", version)
			return
		}
	}

	var destination string
	var repoURL string

	switch len(os.Args) {
	case 2:
		destination = "."
		repoURL = os.Args[1]
	case 3:
		destination = os.Args[1]
		repoURL = os.Args[2]
	default:
		fmt.Println("Usage: yoink [destination] <url>")
		os.Exit(1)
	}

	if err := validator.Destination(destination); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid destination: %v\n", err)
		os.Exit(1)
	}

	host, user, repo, subpath, err := validator.Parse(repoURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid URL: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Yoinking %s/%s/%s [%s] to %s...\n", host, user, repo, subpath, destination)

	if err := downloader.Fetch(host, user, repo, subpath, destination); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Success!")
}
