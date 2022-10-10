# Contributing to the UOR Client and Libraries

Thank you for investing your time in contributing to our project!

When contributing to this repository, please first discuss the change you wish to make via issue or GitHub Discussion
with the maintainers to ensure the change lines up with the long term plans for the project.

## Pull Request Process

### We Use [GitHub Flow](https://docs.github.com/en/get-started/quickstart/github-flow)

Please use the following workflow to make changes to the UOR Client codebase:

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes (Run `make test-unit`)
5. Make sure your code lints (Run `make sanity`)
6. Create a pull request against the `uor-client-go` main branch.


We encourage draft pull requests, when applicable, for early feedback and better communication.

*Note: In the Makefile there are targets for code and file generation. If any changes are made to the gRPC API in the `api/services`
directory, run `make generate-protobuf`. This operation assumes the protobuf compiler and Go plugins are installed. If any changes are made
to the client CLI under `cmd/client`, run `make generate-usage-docs` to update the documentation under `docs/usage`.*


## Report bugs and feature ideas using GitHub's [issues](https://github.com/uor-framework/uor-client-go/issues/new/choose)
Each issue type has a template attached to guide the submission.


## Code Styling

- Run `go fmt` 
- Run `golangci-lint`
- Use `go-imports` (This should be configured to group the standard library, third-party, and uor-client-go module imports separately)


## License
By contributing, you agree that your contributions will be licensed under its [Apache 2.0 License](https://choosealicense.com/licenses/apache-2.0/).