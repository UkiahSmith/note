all: build

build:
	go build -o note ./cmd/note

clean:
	go clean ./...
	rm note
