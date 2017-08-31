.PHONY: build

build:
	go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)"

run: build
	./bado

release: *.go es/*.go 
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d)" -a
	docker build -t vikings/bado .
	docker push vikings/bado
