# Aresctl

[中文文档](README.zh.md)

Command-line tool for the [Ares](https://github.com/xushuhui/ares) framework. Aresctl helps you develop Ares applications faster with code generation, OpenAPI documentation, and GORM model generation.

## Features

- **Project Scaffolding** - Create new Ares projects from template
- **Code Generation** - Generate handlers and CRUD code from templates
- **OpenAPI Generation** - Automatically generate OpenAPI 3.0 specifications from your Go code
- **GORM Code Generation** - Generate GORM models and query code from database schema

## Installation

```bash
go install github.com/xushuhui/aresctl@latest
```

Make sure `$GOPATH/bin` is in your `PATH`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage

### Create a New Project

Create a new Ares project from the ares-layout template:

```bash
aresctl new myproject
```

This will:
- Clone the ares-layout template
- Remove the .git directory
- Update go.mod with your project name
- Initialize a new git repository

### Generate Code

Generate handler or CRUD code from templates:

```bash
# Generate a handler
aresctl generate handler product

# Generate complete CRUD (handler + API types)
aresctl generate crud product
```

### Generate OpenAPI Specification

Navigate to your Ares project directory and run:

```bash
aresctl openapi
```

This will generate an `openapi.yaml` file in your project root by analyzing:
- Route definitions in `internal/server/`
- API structures in `api/`

### Generate GORM Code

Run in your project root:

```bash
aresctl gen
```

`gen` reads `gorm.yaml` from the current directory.

Minimum `gorm.yaml` example (MySQL only):

```yaml
version: "0.1"
database:
  dsn: "user:pass@tcp(127.0.0.1:3306)/your_db?charset=utf8mb4&parseTime=True&loc=Local"
  db: "mysql"
  tables: []
  outPath: "./dao/query"
```

If `tables` is empty, it will generate for all tables.

## Project Structure

Aresctl expects your project to follow the [ares-layout](https://github.com/xushuhui/ares-layout) structure:

```
your-project/
├── internal/
│   └── server/
│       └── http.go          # Route definitions
├── api/
│   ├── user.go              # Request/Response structures
│   └── response.go          # Common response types
└── openapi.yaml             # Generated file
```

### Route Definition Format

Routes should be defined with comments and `@tag` annotations:

```go
func NewHTTPServer(userHandler *handler.UserHandler) *ares.Ares {
    app := ares.Default()

    api := app.Group("/api")
    // Create user @tag Users
    api.POST("/users", userHandler.Create)
    // Get user list @tag Users
    api.GET("/users", userHandler.List)
    // Get user details @tag Users
    api.GET("/users/{id}", userHandler.Get)

    return app
}
```

### API Structure Format

Define your API structures with comments:

```go
// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
    Name  string `json:"name"`  // 用户名
    Email string `json:"email"` // 邮箱
}

// UserResponse 用户响应
type UserResponse struct {
    ID    int64  `json:"id"`    // 用户ID
    Name  string `json:"name"`  // 用户名
    Email string `json:"email"` // 邮箱
}
```

## Shell Autocompletion

Aresctl supports shell autocompletion for bash, zsh, fish, and PowerShell.

**Setup for Zsh:**
```bash
source <(aresctl completion zsh)
```

**Setup for Bash:**
```bash
source <(aresctl completion bash)
```

## Help

```bash
aresctl --help
aresctl new --help
aresctl generate --help
aresctl openapi --help
aresctl gen --help
```

## Development

```bash
# Build from source
git clone https://github.com/xushuhui/aresctl.git
cd aresctl
go build -o aresctl

# Run tests
go test ./...
```

## Ecosystem

- **[Ares](https://github.com/xushuhui/ares)** - Lightweight Go web framework
- **[ares-layout](https://github.com/xushuhui/ares-layout)** - Production-ready application template
- **[aresctl](https://github.com/xushuhui/aresctl)** - Command-line tool (this project)

## License

MIT License
