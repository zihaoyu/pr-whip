SHELL := /bin/bash

build: # compile source code
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/whip  cmd/main.go
.PHONY: build

docker: # build docker image
	docker build -t zihaoyu/pr-whip .
.PHONY: docker

publish: # publish docker image
	docker push zihaoyu/pr-whip
.PHONY: publish
