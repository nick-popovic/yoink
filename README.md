# Yoink

Yoink is a simple CLI tool that allows you to download specific subdirectories or individual files from a Git repository without needing to clone the entire repository history. It uses `git sparse-checkout` under the hood to keep things fast and lightweight.

## Features

- **Selective Download:** Only download what you need (files or folders).
- **Fast:** Uses shallow clones and blob filtering to minimize data transfer.
- **Branch Support:** Fully supports downloading from specific branches by automatically parsing GitHub `/tree/<branch>` or `/blob/<branch>` URLs.
- **Cross-Platform:** Written with a pure Go native file mover, making it truly reliable and seamless across Windows, macOS, and Linux without relying on Unix shell commands.

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

**Download a subdirectory from a specific branch:**

```bash
yoink https://github.com/google/go-github/tree/master/github
```

**Download a single file:**

```bash
yoink https://github.com/google/go-github/blob/master/README.md
```

**Download and rename a file locally:**

```bash
yoink my-readme.md https://github.com/google/go-github/blob/master/README.md
```

## How it works

Yoink parses your URL to determine the repository and the specific branch you want. It creates a temporary directory, initializes a shallow Git repository (`--depth 1 --branch <branch>`), enables `sparse-checkout`, pulls only the requested path, and then uses a native Go recursive copy to move the files to your specified destination before cleaning up.

## Contributing

Contributions are welcome! To ensure code quality, all Pull Requests are automatically tested (via `go build`, `go test`, and `go vet`) by GitHub Actions. All commits must be signed and code cannot be merged into `main` unless all automated tests pass successfully.
