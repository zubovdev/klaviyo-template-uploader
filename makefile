build:
	set GOOS=darwin;
	set GOARCH=amd64;
	go build -v ./cmd/app

.DEFAULT_GOAL := build

.PHONY: build