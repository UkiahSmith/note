.PHONY: help all build build-release clean test 

VERSION := $(shell git tag | grep ^v | sort -V | tail -n 1)
SHORT_HASH := $(shell git rev-parse --short HEAD)
BUILD_TIMESTAMP := $(shell date -u +%FT%TZ)
DEV_TIMESTAMP := $(shell date -u +%s)
LDFLAGS = -ldflags "-X main.buildVersion=${VERSION} -X main.buildHash=${SHORT_HASH} -X main.buildTimestamp=${BUILD_TIMESTAMP}"
DEVLDFLAGS = -ldflags "-X main.buildVersion=dev-${DEV_TIMESTAMP}"

help: #Display this help message.
	@echo ""
	@echo "Note, a templating tool for note-taking."
	@echo ""
	@grep '^[#[a-z].*:' Makefile | sed 's/:.*#/:/' | column -s ':' -t | sort -h

build: #Compile note with extras needed during the developent process.
	go build ${DEVLDFLAGS} -tags=debug -o build/note ./cmd/note

build-release: #Compile note without development hooks and logging.
	go build ${LDFLAGS} -o build/note ./cmd/note

clean: #Remove all build artifacts.
	rm -r build
	go clean ./...

test: #Run the test suite.
	if [[ -x $$(command -v gotest) ]]; then gotest -tags=debug ./... ; else go test -tags=debug ./... ; fi
