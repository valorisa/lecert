BINARY := lecert
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: build clean test release

build:
	go build $(LDFLAGS) -o bin/$(BINARY) ./cmd/lecert/

test:
	go test ./... -v -count=1

clean:
	rm -rf bin/ dist/

release: clean
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} \
		go build $(LDFLAGS) -o dist/$(BINARY)-$${platform%/*}-$${platform#*/}$$([ "$${platform%/*}" = "windows" ] && echo ".exe") ./cmd/lecert/ ; \
		echo "Built: $${platform}"; \
	done
	@echo "Release binaries in dist/"
	@ls -la dist/
