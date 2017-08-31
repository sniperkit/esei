.PHONY: build

build:
	go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)"

run: build
	./bado

release: *.go es/*.go 
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)" -a -o esei_linux
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)" -a -o esei_darwin
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)" -a -o esei_windows
