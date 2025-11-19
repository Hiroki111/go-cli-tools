# goci

Continuous Integration tool for Go programs.

- Building the program using `go build` to verify if the program structure is valid.
- Executing tests using `go test` to ensure the program does what it’s intended to do.
- Executing `gofmt` to ensure the program’s format conforms to the standards.
- Executing `git push` to push the code to the remote shared Git repository that hosts the program code.

## Pre-requisite

Install the following:

- [golangci-lint](https://golangci-lint.run/docs/welcome/install/)
- [gocyclo](https://github.com/fzipp/gocyclo)

## Example usage

```bash
go build

./goci -p <path-to-a-go-project>
```