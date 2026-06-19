APP_NAME   := gbx
MODULE     := github.com/Pixie2468/git-branch-explorer
VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT     := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS    := -s -w -X '$(MODULE)/cmd.Version=$(VERSION)' -X '$(MODULE)/cmd.Commit=$(COMMIT)' -X '$(MODULE)/cmd.Date=$(BUILD_DATE)'

PLATFORMS  := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# ──────────────────────────────────────────────

.PHONY: build install clean test vet dist checksums

## build: Compile for current OS/arch
build:
	go build -ldflags "$(LDFLAGS)" -o $(APP_NAME) .

## install: Build and move to $GOPATH/bin
install:
	go install -ldflags "$(LDFLAGS)" .

## test: Run all tests
test:
	go test -v -race ./...

## vet: Run go vet
vet:
	go vet ./...

## clean: Remove build artifacts
clean:
	rm -rf $(APP_NAME) dist/

## dist: Cross-compile for all release platforms
dist: clean
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d/ -f1); \
		ARCH=$$(echo $$platform | cut -d/ -f2); \
		EXT=""; \
		if [ "$$OS" = "windows" ]; then EXT=".exe"; fi; \
		OUTPUT="dist/$(APP_NAME)-$$OS-$$ARCH$$EXT"; \
		echo "→ Building $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH go build -ldflags "$(LDFLAGS)" -o $$OUTPUT . || exit 1; \
	done
	@echo "✓ All binaries in dist/"

## checksums: Generate SHA256 checksums for dist binaries
checksums: dist
	@cd dist && sha256sum * > checksums.txt
	@echo "✓ Checksums written to dist/checksums.txt"
