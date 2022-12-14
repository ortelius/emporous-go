# Contributing to the UOR Client and Libraries

Thank you for investing your time in contributing to our project!

When contributing to this repository, please first discuss the change you wish to make via GitHub Issues or Discussions
as to ensure the change aligns with the project's long-term plans.

## Pull Request Process

### We Use [GitHub Flow](https://docs.github.com/en/get-started/quickstart/github-flow)

Please use the following workflow to make changes to the UOR Client codebase:

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes (Run `make test-unit`)
5. Make sure your code lints (Run `make sanity`)
6. Create a pull request against the `uor-client-go` `main` branch.


When applicable, we encourage [draft pull requests](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/changing-the-stage-of-a-pull-request) for early feedback and better communication.

*Note: In the Makefile, there are code and file generation targets. If any changes are made to the gRPC API in the `api/services`
directory, run `make generate-protobuf`. This operation assumes the protobuf compiler and Go plugins are installed. If any changes are made
to the client CLI under `cmd/client`, run `make generate-usage-docs` to update the documentation under `docs/usage`.*

## Report bugs and feature ideas using GitHub's [issues](https://github.com/uor-framework/uor-client-go/issues/new/choose)
Each issue type has a template attached to guide the submission.

## Required Tools

To make changes to the gRPC API, please install the following tools:
  
- [protoc](https://github.com/protocolbuffers/protobuf/releases)

Install the plugins need to generated Go code with protoc:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@vlatest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```


## Code Styling

- Run `go fmt`
- Run [golangci-lint](https://github.com/golangci/golangci-lint)
- Use [go-imports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)
  - (This should be configured to group the standard library, third-party, and uor-client-go module imports separately)


## License
By contributing, you agree that your contributions will be licensed under its [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/).
