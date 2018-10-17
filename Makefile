all: build

build:
	go build ./cmd/note

clean:
	go clean ./...
	rm note
