.PHONY: deps test publish clean

GOPATH ?= /go
GOBIN  := $(GOPATH)/bin
PATH   := $(GOPATH)/bin:$(PATH)
PROJ   := kubectl-ssm_secret
DOCKER_USERNAME ?= Monkey
DOCKER_PASSWORD ?= Magic
PKG_DIR := pkg

#LDFLAGS := -ldflags "-X main.commit=`git rev-parse HEAD`"
LDFLAGS := -ldflags '-extldflags "-static"'

all: deps fmt test $(PROJ) publish

deps:
	@echo "--- collecting ingredients :bento:"
	go mod download

fmt:
	$(foreach dir, $(wildcard $(PKG_DIR)/*), go fmt $(dir)/*.go;)
	go fmt cmd/*.go
	$(foreach dir, $(wildcard $(PKG_DIR)/*), go vet $(dir)/*.go;)
	go vet cmd/*.go

test: deps 
	@echo "--- Is this thing working? :hammer_and_wrench:"
	GOOS=linux go test -cover -v ./pkg/...

$(PROJ): deps
	@echo "--- make it so :cooking:"
	go build $(LDFLAGS) -o $@ -v ./cmd/
	touch $@ && chmod 755 $@

ifdef TRAVIS_TAG
publish: deps
	@echo "--- release :octocat:"
	echo docker login -u "$(DOCKER_USERNAME)" -p "$(DOCKER_PASSWORD)"
	echo goreleaser --skip-validate --rm-dist
endif

clean:
	rm -rf $(PROJ) $(PROJ)-windows-amd64.exe $(PROJ)-linux-amd64 $(PROJ)-darwin-amd64 dist
