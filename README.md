# asyncer

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/dmitrymomot/asyncer)](https://github.com/dmitrymomot/asyncer/tags)
[![Go Reference](https://pkg.go.dev/badge/github.com/dmitrymomot/asyncer.svg)](https://pkg.go.dev/github.com/dmitrymomot/asyncer)
[![License](https://img.shields.io/github/license/dmitrymomot/asyncer)](https://github.com/dmitrymomot/asyncer/blob/main/LICENSE)


[![Tests](https://github.com/dmitrymomot/asyncer/actions/workflows/tests.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/tests.yml)
[![CodeQL Analysis](https://github.com/dmitrymomot/asyncer/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/codeql-analysis.yml)
[![GolangCI Lint](https://github.com/dmitrymomot/asyncer/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/golangci-lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dmitrymomot/asyncer)](https://goreportcard.com/report/github.com/dmitrymomot/asyncer)

This is a simple, reliable, and efficient distributed task queue in Go.
The asyncer just wrapps [hibiken/asynq](https://github.com/hibiken/asynq) package with some predefined settings. So, if you need more flexibility, you can use [hibiken/asynq](https://github.com/hibiken/asynq) directly.

## Usage

See [_example](https://github.com/dmitrymomot/asyncer/tree/main/_example) directory for usage examples.

## Todo

- [ ] Add tests

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/dmitrymomot/asyncer/tree/main/LICENSE) file for details. This project contains some code from [hibiken/asynq](https://github.com/hibiken/asynq) package, which is also licensed under the MIT License.