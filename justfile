# To add commands, see https://cheatography.com/linux-china/cheat-sheets/justfile/

[doc('List all available commands')]
help:
  @just -l --list-heading $'Available commands:\n'

[doc('Run golangci-lint')]
lint:
  docker run --rm -v $(pwd):/app -v ~/.cache/golangci-lint/:/root/.cache -w /app golangci/golangci-lint:v2.2 golangci-lint run --timeout 3m --verbose

[doc('Run all unit tests')]
unit-tests:
  @go test ./...
