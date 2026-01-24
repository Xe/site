# Repository Guidelines

Go-based static site generator with JavaScript tooling for content automation.

## Project Structure

```
cmd/        # Binaries: fabricate-generation, github-sponsor-webhook, hydrate, no-way-to-prevent-this, patreon-saasproxy, xesite (main), xesitectl
internal/   # Private packages
lume/       # Static site generator configuration and pages
```

## Commands

| Command | Description |
|---------|-------------|
| `npm test` | Run `go test ./...` |
| `npm run dev` | Development server |
| `npm run deploy` | Deploy to Kubernetes |
| `npm run extract-meta` | Extract metadata from content |
| `npm run validate-content-dates` | Validate blog post dates |
| `go build ./...` | Compile all packages |
| `go run ./cmd/<app>` | Run specific binary |

## Conventions

- [Go Conventions](.claude/go-conventions.md) - Formatting, style, testing
- [JavaScript Conventions](.claude/javascript-conventions.md) - Prettier, HTML/CSS
- [Git Workflow](.claude/git-workflow.md) - Commits, PRs, security
