# .PHONY: all
# all: container

ARCH?=amd64
OS?=darwin
PKG=github.com/hiraken-w/event-cwl-exporter
REPO_INFO=$(shell git config --get remote.origin.url)
GO111MODULE=on
GOPROXY=direct
GOBIN:=$(shell pwd)/.bin

.EXPORT_ALL_VARIABLES:

ifndef GIT_COMMIT
  GIT_COMMIT := git-$(shell git rev-parse --short HEAD)
endif

LDFLAGS=-X $(PKG)/version.COMMIT=$(GIT_COMMIT) -X $(PKG)/version.RELEASE=$(TAG) -X $(PKG)/version.REPO=$(REPO_INFO)

app: cmd/*go internal/*/*.go
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -a -installsuffix cgo -ldflags '-s -w $(LDFLAGS)' -o app ./cmd

.PHONY: clean
clean:
	rm -rf app

.PHONY: container
container:
	docker build .

.PHONY: lint
lint:
	go install -v github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOBIN)/golangci-lint run --deadline=10m
