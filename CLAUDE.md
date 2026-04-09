# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Test

```bash
# Build binary
go build -o aresctl

# Run all tests
go test ./...

# Run specific test
go test -run TestName ./cmd/
```

## Architecture

Aresctl is a CLI tool for the [Ares](https://github.com/xushuhui/ares) Go web framework. It uses [cobra](https://github.com/spf13/cobra) for command parsing.

### Command Structure (`cmd/`)

- `root.go` - Root command and `openapi` subcommand (OpenAPI 3.0 spec generation)
- `new.go` - `new` subcommand: scaffolds new projects from [ares-layout](https://github.com/xushuhui/ares-layout)
- `generate.go` - `generate` subcommand: handler and CRUD code generation from templates
- `gen.go` - `gen` subcommand: GORM model/query code generation from database schema (MySQL only)
- `openapi.go` - OpenAPI spec generation logic using Go AST parsing

### Key Dependencies

- `github.com/spf13/cobra` - CLI framework
- `gorm.io/gen` - GORM code generation
- `gopkg.in/yaml.v3` - YAML config parsing

### Configuration

- `gorm.yaml` - Configuration file for `gen` command (DSN, tables, output path)
