.PHONY: help all build build-release clean test 

help: #Display this help message.
	@echo ""
	@echo "Note, a templating tool for note-taking."
	@echo ""
	@grep '^[#[a-zA-Z].*:' Makefile | sed 's/:.*#/:/' | column -s ':' -t | sort -h

all: build

build:
	go build -tags=debug -o build/note ./cmd/note

build-release:
	go build -o build/note ./cmd/note

clean:
	rm -r build
	go clean ./...

test:
	if [[ -x $$(command -v gotest) ]]; then gotest -tags=debug ./... ; else go test -tags=debug ./... ; fi
