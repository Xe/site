# Git Workflow

## Commit Messages

Follow **Conventional Commits** format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`

- Add `!` after type/scope for breaking changes
- Keep descriptions: concise, imperative, lowercase, no trailing period
- Reference issues/PRs in footers

## Attribution

AI agents must disclose tool and model in commit footer:

```
Assisted-by: [Model Name] via [Tool Name]
```

Example: `Assisted-by: GLM 4.6 via Claude Code`

## Pull Requests

- Clear description of changes
- Reference related issues
- Pass CI (`npm test`)
- Screenshots for UI changes (optional)

## Security

- Secrets in environment variables or `.env` files only (never in repo)
- Run `npm audit` periodically and address vulnerabilities
