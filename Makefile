GO := go

GO_BUILD_PACKAGES := ./cmd/...
GO_BUILD_BINDIR :=./bin

all: clean vendor test-unit build
.PHONY: all

build:
	mkdir -p ${GO_BUILD_BINDIR}
	GOOS=linux GOARCH=amd64 $(GO) build -o $(GO_BUILD_BINDIR)/client $(GO_BUILD_PACKAGES)
.PHONY: build

vendor:
	$(GO) mod tidy
	$(GO) mod verify
	$(GO) mod vendor
.PHONY: vendor

clean:
	@rm -rf ./$(GO_BUILD_BINDIR)/*
.PHONY: clean

test-unit:
	$(GO) test $(GO_BUILD_FLAGS) -coverprofile=coverage.out -race -count=1 ./...
.PHONY: test-unit

sanity: vendor format vet
	git diff --exit-code
.PHONY: sanity

format: 
	$(GO) fmt ./...
.PHONY: format

vet: 
	$(GO) vet ./...
.PHONY: vet