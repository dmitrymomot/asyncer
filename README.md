# asyncer


[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/dmitrymomot/asyncer)](https://github.com/dmitrymomot/asyncer)
[![Tests](https://github.com/dmitrymomot/asyncer/actions/workflows/tests.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/tests.yml)
[![CodeQL Analysis](https://github.com/dmitrymomot/asyncer/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/codeql-analysis.yml)
[![GolangCI Lint](https://github.com/dmitrymomot/asyncer/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/dmitrymomot/asyncer/actions/workflows/golangci-lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dmitrymomot/asyncer)](https://goreportcard.com/report/github.com/dmitrymomot/asyncer)
[![Go Reference](https://pkg.go.dev/badge/github.com/dmitrymomot/asyncer.svg)](https://pkg.go.dev/github.com/dmitrymomot/asyncer)
[![License](https://img.shields.io/github/license/dmitrymomot/asyncer)](https://github.com/dmitrymomot/asyncer/blob/main/LICENSE)

This is a simple, reliable, and efficient distributed task queue in Go.
The asyncer just wrapps [hibiken/asynq](https://github.com/hibiken/asynq) package with some predefined settings. So, if you need more flexibility, you can use [hibiken/asynq](https://github.com/hibiken/asynq) directly.