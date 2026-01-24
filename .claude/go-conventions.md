# Go Conventions

## Formatting

- Use `go fmt` or `goimports` for formatting
- Indent with tabs
- Variables: `camelCase`
- Exported identifiers: `PascalCase`
- Files: `snake_case`
- Packages: `lower-case` module names

## Testing

- Write tests using the standard `testing` package (`*_test.go`)
- Keep test files next to the code they cover
- Prefer table-driven tests
- Aim for high coverage on new modules (existing coverage not enforced)
