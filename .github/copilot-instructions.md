# Go Project Setup Instructions

This is a basic Go project for a todo authentication application.

- [x] Create .github/copilot-instructions.md file
- [x] Get project setup information
- [x] Scaffold the Go project structure
- [x] Compile and verify the project
- [x] Create VS Code run task
- [x] Finalize documentation

## Quick Start

**Build the project:**
```bash
go build -o bin/server ./cmd/server
```

**Run the project:**
```bash
./bin/server
```

Or use the VS Code Build task (Shift+Cmd+B) to build and run.

## Project Structure

- `cmd/server/` - Main application entry point
- `internal/models/` - Data models (User, Task)
- `internal/handlers/` - HTTP handlers (coming soon)
- `go.mod` - Go module definition
