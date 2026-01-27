# Claude Code Context - sm-ssh-add

## Project Context Files

**IMPORTANT:**
- Always read `.claude/TODO.md` first to understand current development priorities and progress
- Always read `.claude/project.md` and `README.md` first to understand project purpose, specification and structure

## Tech Stack

- **Language**: Go
- **Agent**: Use `golang-expert` subagent for Go-specific tasks
- **Documentation**: Follow `go-doc-comments` skill (official Go conventions from go.dev/doc/comment)
- **Testing**: Follow `golang-testing` skill (test files in same directory as source)

## Development Workflow

**CRITICAL: Use `superpowers:test-driven-development` skill for ALL development**
- Invoke TDD skill before implementing any feature
- Tests drive implementation - always write tests first

## Quick Commands

```bash
go build -o sm-ssh-add .
go test ./...
go run main.go
```
