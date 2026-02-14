# Contributing to P9 Labs

We are building P9 without heavy frameworks to keep the binary light and the code readable.

## Workflow
1. **Fork the repo**: This creates your own copy where you have permission to work.
2. **Create a branch**: Name it after your feature (e.g., `feat-dns-trace`).
3. **No Frameworks**: Use the Go standard library (`os`, `flag`, `net`) wherever possible.
4. **Submit a PR**: Propose merging your fork's branch into our `main` branch.

## Project Structure
- `cmd/p9/main.go`: The main switchboard.
- `internal/`: Where the actual logic for DNS and Network tools lives.