# AGENTS.md

## Requirements
- Read requirements in `docs/requirements.md` and indicate when the requests are drifting. 

## Build Commands

### Go
- `make build` - Build all binaries to bin/
- `make test` - Run all tests
- `make lint` - Run golangci-lint
- `make lint-fix` - Run golangci-lint with auto-fix
- `make fmt` - Format code
- `make clean` - Remove build artifacts
- `make help` - Show all available targets

## Code Style Guidelines

### Go
- Follow official [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting (enforced by `go fmt`)
- Use `go vet` and `golangci-lint` for static analysis
- Package naming: short, lowercase, single-word names (no underscores or mixedCaps)
- Interface naming: single-method interfaces end with "-er" suffix (e.g., `Reader`, `Writer`)
- Error handling: check all errors explicitly, return errors up the call stack
- Error messages: lowercase, no punctuation (e.g., `fmt.Errorf("failed to connect: %w", err)`)
- Use `context.Context` for cancellation and timeouts in long-running operations
- Imports: standard library first, then third-party, then internal packages (separated by blank lines)
- Variable naming: short names for short scope (e.g., `i`, `r`, `w`), descriptive for larger scope
- Avoid naked returns in functions longer than a few lines
- Use struct literals with field names for clarity
- Prefer table-driven tests with `t.Run()` for subtests
- Export only what's necessary; keep internal details unexported
- Comment exported functions, types, and constants with complete sentences starting with the name
- Use `defer` for cleanup operations (closing files, unlocking mutexes)
- Avoid package-level state; prefer dependency injection

### General
- Comment code only when necessary for complex logic
- Follow existing patterns in codebase
- Security: Never commit API keys or secrets

## Testing Strategy

### Go
- **Table-Driven Tests**: Use table-driven patterns for all unit tests to cover multiple scenarios efficiently.
- **Concurrency**: Always run tests with the `-race` flag (`go test -race ./...`) to detect data races, especially for components using `sync` primitives.
- **Coverage**: Aim for high statement coverage (80%+) in `internal/` packages. Use `go test -cover` to verify.
- **Documentation**: Use comments above test functions or within test tables to describe the specific scenario being tested.
- **Subtests**: Use `t.Run()` for individual cases in table-driven tests to allow running specific scenarios.
