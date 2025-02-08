# INTERNAL.md

This document contains notes for development.

## Running Tests

- To run all tests:

  ```bash
  go test ./...
  ```

  To get coverage:

  ```bash
  go test -cover ./... -coverprofile=coverage.out
  ```

  To view coverage:

  ```bash
  go tool cover -html=coverage.out
  ```
