run: test
	./reinze

test: fmt
	go test

fmt: build
	go fmt

build:
	go build

.PHONY: fmt test build run