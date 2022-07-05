.PHONY: deps test publish clean

GOPATH ?= /go
GOBIN  := $(GOPATH)/bin
PATH   := $(GOPATH)/bin:$(PATH)
PROJ   := kubectl-ssm_secret
DOCKER_USERNAME ?= Monkey
DOCKER_PASSWORD ?= Magic
PKG_DIR := pkg
VERSION		:= $(shell git describe --tags --match "v[0-9]*" --exact HEAD 2>/dev/null || echo snapshot)

#LDFLAGS := -ldflags "-X main.commit=`git rev-parse HEAD`"
LDFLAGS := -ldflags "-s -w -X github.com/pr8kerl/kubectl-ssm-secret/pkg/cmd.version=$(VERSION)"

all: deps fmt test $(PROJ) 

deps:
	@echo "--- collecting ingredients :bento:"
	go mod tidy -v

fmt:
	go fmt ./...
	go vet ./...

test: deps 
	@echo "--- Is this thing working? :hammer_and_wrench:"
	GOOS=linux go test -cover -v ./pkg/...

$(PROJ): deps
	@echo "--- make it so :cooking:"
	go build $(LDFLAGS) -o $@ -v ./cmd/$(PROJ)
	touch $@ && chmod 755 $@

build: deps 
	@echo "--- goreleaser build :cooking:"
	goreleaser build --single-target

release: deps 
	@echo "--- package it up! :box:"
	goreleaser --skip-validate --skip-publish --rm-dist
	sha256sum dist/kubectl-ssm-secret*.gz

publish: deps
	@echo "--- release :octocat:"
	goreleaser --skip-validate --rm-dist
	sha256sum dist/kubectl-ssm-secret*.gz

clean:
	rm -rf $(PROJ) $(PROJ)-* dist
