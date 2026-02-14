# Contributing to P9 Labs

We are building P9 without heavy frameworks to keep the binary light and the code readable.

## Project Structure
```
p9/
├── cmd/
│   └── main.go          ← Entry point, CLI routing only
├── internal/
│   ├── ports/           ← Port operations (local/remote)
│   │   ├── local.go     ← List listening ports
│   │   └── remote.go    ← Check remote connectivity
│   ├── dns/             ← DNS/domain operations
│   │   └── lookup.go    ← DNS resolution, WHOIS
│   └── cli/             ← Output formatting ONLY
│       └── output.go    ← Display logic, no business logic
```

**Guidelines:**
- Business logic goes in `internal/ports/` or `internal/dns/`
- Output formatting goes in `internal/cli/`
- Keep `main.go` thin - just routing

## Workflow
1. **Fork the repo**: This creates your own copy where you have permission to work.
2. **Create a branch**: Name it after your feature (e.g., `feat-dns-trace`).
3. **No Frameworks**: Use the Go standard library (`os`, `flag`, `net`) wherever possible.
4. **Submit a PR**: Propose merging your fork's branch into our `main` branch.

## Pull Request Guidelines
- **One feature per PR** - keeps reviews manageable
- **Test your changes** - actually run the tool
- **Update README** - if adding user-facing features
- **Clear description** - explain what and why

## Code Quality

Before submitting:
- Run `go fmt ./...` to format code
- Run `go vet ./...` to check for issues
- Test your changes: `go build -o p9 ./cmd && ./p9 [your-feature]`