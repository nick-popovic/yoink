# 🚀 Yoink v1.0.0 — The Initial Release!

We are excited to announce the first official release of **Yoink**! 🎉

Yoink is a fast, lightweight CLI tool designed to solve a common developer pain point: downloading specific subdirectories or individual files from a Git repository without the overhead of cloning the entire project history.

## ✨ Features

- **Selective Downloads:** Fetch only the specific folder or file you need.
- **Optimized Performance:** Uses `git sparse-checkout` and partial clones (`--filter=blob:none`) to minimize network traffic and disk usage.
- **Simple CLI:** A straightforward, easy-to-use interface.
- **Cross-Platform:** Works seamlessly on Windows, macOS, and Linux.

## 📦 How to get it?

### Using Go

If you have Go installed, you can install Yoink directly:

```bash
go install github.com/nick-popovic/yoink@latest
```

### From GitHub Releases

Pre-built binaries for all major platforms are available in the **Assets** section below.

## 🛠 Usage Example

```bash
# Yoink a specific directory from a repository
yoink . https://github.com/nick-popovic/yoink/tree/main/internal/downloader

# Yoink a single file and rename it locally
yoink my-validator.go https://github.com/nick-popovic/yoink/blob/main/internal/validator/validator.go
```

## 🔍 How it works

Yoink leverages modern Git features (`sparse-checkout` and shallow clones) to efficiently isolate and fetch only the requested paths, ensuring you get your files quickly and with minimal data transfer.

---

**Thank you for using Yoink!** If you find it useful, feel free to [star the project on GitHub](https://github.com/nick-popovic/yoink) and report any issues you encounter.
