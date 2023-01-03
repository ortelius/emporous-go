GO := go

GO_BUILD_PACKAGES := ./cmd/client
GO_BUILD_BINDIR :=./bin
EXECUTABLE_NAME := "uor-client-go"
GIT_COMMIT := $(or $(SOURCE_GIT_COMMIT),$(shell git rev-parse --short HEAD))
GIT_TAG :="$(shell git tag | sort -V | tail -1)"

GO_LD_EXTRAFLAGS :=-X github.com/uor-framework/uor-client-go/cmd/client/commands.version="$(shell git tag | sort -V | tail -1)" \
				   -X github.com/uor-framework/uor-client-go/cmd/client/commands.buildData="dev" \
				   -X github.com/uor-framework/uor-client-go/cmd/client/commands.commit="$(GIT_COMMIT)" \
				   -X github.com/uor-framework/uor-client-go/cmd/client/commands.buildDate="$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')"

build: prep-build-dir
	$(GO) build -o $(GO_BUILD_BINDIR)/$(EXECUTABLE_NAME)  -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: build

cross-build-darwin-amd64:
	env GOOS=darwin  GOARCH=amd64	$(GO) build -o $(GO_BUILD_BINDIR)/$(EXECUTABLE_NAME)-darwin-amd64  -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-darwin-amd64

cross-build-darwin-arm64:
	env GOOS=darwin  GOARCH=arm64	$(GO) build -o $(GO_BUILD_BINDIR)/$(EXECUTABLE_NAME)-go-darwin-arm64  -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-darwin-arm64

cross-build-windows-amd64:
	env GOOS=windows GOARCH=amd64	$(GO) build -o $(GO_BUILD_BINDIR)/$(EXECUTABLE_NAME)-windows-amd64 -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-windows-amd64

cross-build-linux-amd64:
	env GOOS=linux   GOARCH=amd64	$(GO) build -o $(GO_BUILD_BINDIR)/$(EXECUTABLE_NAME)-linux-amd64   -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-linux-amd64

cross-build-linux-arm64:
	env GOOS=linux   GOARCH=arm64   $(GO) build -o $(GO_BUILD_BINDIR)/$(EXECUTABLE_NAME)-linux-arm64   -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-linux-arm64

cross-build-linux-ppc64le:
	env GOOS=linux   GOARCH=ppc64le $(GO) build -o $(GO_BUILD_BINDIR)/$(EXECUTABLE_NAME)-linux-ppc64le -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-linux-ppc64le

cross-build-linux-s390x:
	env GOOS=linux   GOARCH=s390x	$(GO) build -o $(GO_BUILD_BINDIR)/$(EXECUTABLE_NAME)-linux-s390x   -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-linux-s390x

cross-build-linux-riscv64:
	env GOOS=linux   GOARCH=riscv64 $(GO) build -o $(GO_BUILD_BINDIR)/$(EXECUTABLE_NAME)-linux-riscv64 -ldflags="$(GO_LD_EXTRAFLAGS)" $(GO_BUILD_PACKAGES)
.PHONY: cross-build-linux-riscv64

cross-build: prep-build-dir cross-build-darwin-amd64 cross-build-darwin-arm64 cross-build-windows-amd64 cross-build-linux-amd64 cross-build-linux-arm64 cross-build-linux-ppc64le cross-build-linux-s390x cross-build-linux-riscv64
.PHONY: cross-build

prep-build-dir:
	mkdir -p ${GO_BUILD_BINDIR}
.PHONY: prep-build-dir

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

sanity: vendor format vet generate-usage-docs
	git diff --exit-code
.PHONY: sanity

format: 
	$(GO) fmt ./...
.PHONY: format

vet: 
	$(GO) vet ./...
.PHONY: vet

generate-usage-docs: prep-build-dir
	# The gendoc executable is built with the name "uor-client-go" since the
    # root command name is built from the base path name of the cli
    # at runtime.
	$(GO) build -o $(GO_BUILD_BINDIR)/tmp/$(EXECUTABLE_NAME) "./cmd/gendoc"
	$(GO_BUILD_BINDIR)/tmp/$(EXECUTABLE_NAME) "docs/usage"
	@rm -rf ./$(GO_BUILD_BINDIR)/tmp/
.PHONY: generate-usage-docs

generate-protobuf:
	protoc api/services/*/*/*.proto --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_out=. --go_opt=paths=source_relative --proto_path=.
.PHONY: generate-protobuf

all: clean vendor test-unit build
.PHONY: all

cross-build-all: clean vendor test-unit cross-build
.PHONY: cross-build-all
