VERSION ?= latest

.PHONY: build
build:
	CGO_ENABLED=0 go build -o build/repo-audit main.go

.PHONY: docker-build
docker-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/repo-audit *.go
	docker build -t quay.io/helmpack/repo-audit:$(VERSION) .

.PHONY: docker-push
docker-push:
	docker push quay.io/helmpack/repo-audit