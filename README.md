# gbx — Git Branch Explorer

A terminal UI for exploring git branches and commits. Navigate branches, view commit history, and switch branches — all from your terminal.

![TUI Preview](https://img.shields.io/badge/TUI-bubbletea_v2-ff69b4?style=flat-square)
![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-square&logo=go)
![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)

## Install

### One-liner (Linux / macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/Pixie2468/git-branch-explorer/main/install.sh | bash
```

Install to a custom directory:

```bash
curl -fsSL https://raw.githubusercontent.com/Pixie2468/git-branch-explorer/main/install.sh | bash -s -- --dir ~/.local/bin
```

### Go install

```bash
go install github.com/Pixie2468/git-branch-explorer@latest
```

### From source

```bash
git clone https://github.com/Pixie2468/git-branch-explorer.git
cd git-branch-explorer
make build
./gbx run
```

## Usage

```bash
# Run in current directory
gbx run

# Run in a specific repo
gbx run --dir /path/to/repo

# Set a timeout for git commands
gbx run --timeout 5s

# Check version
gbx version
```

## Keybindings

| Key | Action |
|-----|--------|
| `Tab` / `→` | Switch focus to commits pane |
| `Shift+Tab` / `←` | Switch focus to branches pane |
| `↑` / `k` | Move cursor up |
| `↓` / `j` | Move cursor down |
| `Enter` | Switch to selected branch |
| `q` / `Ctrl+C` | Quit |

## Build

```bash
make build         # Build for current platform
make dist          # Cross-compile all platforms
make checksums     # Generate SHA256 checksums
make test          # Run tests
make vet           # Run go vet
make clean         # Remove build artifacts
```

## Release

Releases are automated via GitHub Actions. Tag a commit and push:

```bash
git tag v0.1.0
git push origin v0.1.0
```

This triggers the [release workflow](.github/workflows/release.yml), which cross-compiles, generates checksums, and publishes a GitHub Release with all binaries.

## License

MIT
