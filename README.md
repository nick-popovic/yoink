# Yoink

Yoink is a simple CLI tool that allows you to download specific subdirectories or individual files from a Git repository without needing to clone the entire repository. It uses `git sparse-checkout` under the hood to keep things fast and lightweight.

## Features

- **Selective Download:** Only download what you need (files or folders).
- **Fast:** Uses shallow clones and blob filtering to minimize data transfer.
- **Cross-Platform:** Works on Windows, macOS, and Linux.

## Installation

If you have Go installed:

```bash
go install github.com/nick-popovic/yoink@latest
```

Alternatively, you can download the pre-built binaries from the [Releases](https://github.com/nick-popovic/yoink/releases) page.

## Usage

```bash
yoink [destination] <url>
```

- If `destination` is omitted, it defaults to the current directory.
- `url` should be the full URL to the file or directory you want to "yoink".

### Examples

**Download a subdirectory:**
```bash
yoink https://github.com/google/go-github/tree/master/github
```

**Download a single file:**
```bash
yoink https://github.com/google/go-github/blob/master/README.md
```

**Download and rename a file:**
```bash
yoink my-readme.md https://github.com/google/go-github/blob/master/README.md
```

## How it works

Yoink creates a temporary directory, initializes a shallow Git repository with `sparse-checkout` enabled, pulls only the requested path, and then moves it to your specified destination before cleaning up.
