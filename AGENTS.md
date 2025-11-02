# Repository Guidelines

## Project Structure & Module Organization

```text
├─ cmd/            # Main applications (each sub‑directory is a binary)
├─ docs/           # Documentation assets
├─ internal/       # Private packages used by the repo
├─ lume/           # Static site generator configuration
├─ lume/src        # Static site pages
├─ test files      # Go test files live alongside source (`*_test.go`)
├─ go.mod          # Go module definition
└─ package.json    # npm scripts and JS tooling
```

Source code is primarily Go; JavaScript tooling lives under `node_modules` and the root `package.json`.

## Development Workflow

### Build, Test & Development Commands

| Command                      | Description                                         |
| ---------------------------- | --------------------------------------------------- |
| `npm run generate`           | Regenerates protobuf, Go code, and runs formatters. |
| `npm test` or `npm run test` | Runs `generate` then executes `go test ./...`.      |
| `go build ./...`             | Compiles all Go packages.                           |
| `go run ./cmd/<app>`         | Runs a specific binary from `cmd/`.                 |
| `npm run format`             | Formats Go (`goimports`) and JS/HTML (`prettier`).  |

### Code Formatting & Style

- **Go** – use `go fmt`/`goimports`; tabs for indentation, `camelCase` for variables, `PascalCase` for exported identifiers.
- **JavaScript/HTML/CSS** – formatted with Prettier (2‑space tabs, trailing commas, double quotes).
- Files are snake_case; packages use lower‑case module names.
- Run `npm run format` before committing.

### Testing

- Tests are written in Go using the standard `testing` package (`*_test.go`).
- Keep test files next to the code they cover.
- Run the full suite with `npm test`.
- Aim for high coverage on new modules; existing coverage is not enforced.
- **Go** – follow the standard library style; prefer table‑driven tests.

## Code Quality & Security

### Commit Guidelines

Commit messages follow **Conventional Commits** format:

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`

- Add `!` after type/scope for breaking changes or include `BREAKING CHANGE:` in the footer.
- Keep descriptions concise, imperative, lowercase, and without a trailing period.
- Reference issues/PRs in the footer when applicable.

### Attribution Requirements

AI agents must disclose what tool and model they are using in the "Assisted-by" commit footer:

```text
Assisted-by: [Model Name] via [Tool Name]
```

Example:

```text
Assisted-by: GLM 4.6 via Claude Code
```

### Additional Guidelines

## Pull Request Requirements

- Include a clear description of changes.
- Reference any related issues.
- Pass CI (`npm test`).
- Optionally add screenshots for UI changes.

### Security Best Practices

- Secrets never belong in the repo; use environment variables via `.env` files.
- Run `npm audit` periodically and address reported vulnerabilities.

_This file is consulted by the repository's tooling. Keep it up‑to‑date as the project evolves._
