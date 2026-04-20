# Rules for AI Agents (for Go projects)

## Repository Scope

| Path | Policy |
| - | - |
| Module packages | Modifiable. Agents own quality, testing, and documentation for all changes. |
| `vendor/` | **Read-only.** Assume correct; report issues to the user without attempting fixes. |
| `build/` | Not under version control. Use as a scratch area for any temporary files needed during a session. |

For standard library and external packages, always consult the relevant specifications or documentation rather than guessing behavior.

## Process

- Follow the existing project coding style and conventions.
- Run all shell commands with `timeout` to prevent the agent from hanging on unresponsive processes.
- Never use `sudo` or install software automatically. If a needed tool is unavailable, suggest the user install it.
- Do not create new files (source, tests, docs, or otherwise) unless explicitly instructed.
- Continuously refactor: if a more readable and provably correct alternative exists, apply it.
- Be thorough—before adding or changing production code, verify that every code path functions correctly.

### Exit Criteria

Before finishing any editing session:
1. Proofread all modifications for correctness and intent.
2. Build, test (`go test`), vet (`go vet`), and format (`go fmt`) the code; fix any issues found.

## Go Language Conventions

- Maintain compatibility with Go 1.21.
- Use the standard library instead of third-party packages whenever applicable.

## Code Organization

### Imports

Group `import` directives as follows, sorted alphabetically within each group:

1. Standard library imports
2. Third-party imports

## Comments and Documentation

- Write in English: accurate, professional, concise. Do not state the obvious.
- Do not add comments for struct fields; use self-documenting field names instead.
- Document behavior in code comments; update `README.md` when documented content changes.

## Logging

Always use the `github.com/hexian000/gosnippets/slog` package for logging. Log errors in `"(%T) %s"` format whenever available.

## Testing

Use industry best practices for unit testing; aim for high coverage.
