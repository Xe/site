# Repository Guidelines

## Project Structure & Module Organization
- `cmd/` – entry‑point binaries (e.g., `xesite`, `patreon-saasproxy`). Each subdirectory contains a `main.go`.
- `internal/` – reusable Go packages for core logic.
- `dhall/` – configuration files written in Dhall, used by the site generator.
- `lume/` – static content (blog posts, videos) in Markdown/MDX.
- `scripts/` – helper scripts (e.g., `fabricate-generation`, `imgoptimize`).
- `manifest/` – Kubernetes manifests for deployment.
- `var/` – generated artifacts (compiled binaries, assets).

## Build, Test, and Development Commands
- `go build ./cmd/xesite` – compile the main site binary.
- `go run ./cmd/xesite --site-url https://preview.xeiaso.net --devel` – start the dev server (see `package.json` script `dev`).
- `npm run dev` – alias for the above Go command.
- `npm run deploy` – deploy the site to production.
- `go test ./...` – run any Go tests (currently none).

## Coding Style & Naming Conventions
- Go code follows `gofmt`; run `go fmt ./...` before committing.
- Use `camelCase` for variables/functions, `PascalCase` for exported types.
- Indentation: tabs (default `go fmt`).
- Dhall files use kebab‑case filenames (e.g., `my-config.dhall`).
- Bash/Node scripts use `snake_case` for variable names.

## Testing Guidelines
- Add tests using the `testing` package; name files `*_test.go` and place them alongside the package.
- Run `go test ./...` to execute all tests.

## Commit & Pull Request Guidelines
- Commit messages follow the conventional format:
  - `type(scope): short description`
  - Example: `feat(cli): add --site-url flag`.
- Keep commits atomic and self‑contained.
- PR description must include:
  - Summary of changes.
  - Related issue number (if any).
  - Verification steps (e.g., run `go run ./cmd/xesite`).

## Security & Configuration Tips
- Secrets live in `.env` and must never be committed; `.gitignore` already excludes it.
- When adding new environment variables, document them in `README.md` or a dedicated config file.

