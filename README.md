# Aresctl

[中文文档](README.zh.md)

Command-line tool for the [Ares](https://github.com/xushuhui/ares) framework. Aresctl helps you develop Ares applications faster with code generation, OpenAPI documentation, and more.

## Features

- 🚀 **OpenAPI Generation** - Automatically generate OpenAPI 3.0 specifications from your Go code
- 📝 **Code Analysis** - Parse route definitions and API structures
- 🎯 **Convention-based** - Works seamlessly with [ares-layout](https://github.com/xushuhui/ares-layout) project structure
- ⚡ **Fast & Lightweight** - Minimal dependencies, quick execution

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

**Example output:**
```
Cloning ares-layout template...
✓ Successfully created project 'myproject'

Next steps:
  cd myproject
  docker-compose -f deploy/docker-compose.yml up -d
  go run main.go
```

### Generate Code

Generate handler or CRUD code from templates:

**Generate a handler:**
```bash
aresctl generate handler product
```

**Output:**
```
Created: internal/handler/product.go
✓ Successfully generated handler for 'product'
```

**Generate complete CRUD (handler + API types):**
```bash
aresctl generate crud product
```

**Output:**
```
Created: internal/handler/product.go
Created: api/product.go
✓ Successfully generated crud for 'product'

Next steps:
1. Add routes in internal/server/http.go:
   api.POST("/products", productHandler.Create)
   api.GET("/products", productHandler.List)
   api.GET("/products/{id}", productHandler.Get)
   api.PUT("/products/{id}", productHandler.Update)
   api.DELETE("/products/{id}", productHandler.Delete)
2. Implement business logic in internal/biz/
3. Implement data access in internal/data/
```

### Generate OpenAPI Specification

Navigate to your Ares project directory and run:

```bash
aresctl openapi
```

This will generate an `openapi.yaml` file in your project root by analyzing:
- Route definitions in `internal/server/`
- API structures in `api/`

**Example output:**
```
✓ Generated openapi.yaml successfully
```

### Generate GORM Code

Run in your project root:

```bash
aresctl gen
```

`gen` reads `gorm.yaml` from current directory by default.

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

### View Help

Get general help:

```bash
aresctl --help
```

**Output:**
```
Aresctl is a powerful CLI tool that helps you develop applications with the Ares framework. It provides code generation, OpenAPI documentation, and more.

Usage:
  aresctl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  gen         Generate GORM models and query code
  generate    Generate code from templates
  help        Help about any command
  new         Create a new Ares project
  openapi     Generate OpenAPI specification

Flags:
  -h, --help   help for aresctl

Use "aresctl [command] --help" for more information about a command.
```

Get command-specific help:

```bash
# New project help
aresctl new --help

# Generate code help
aresctl generate --help

# OpenAPI help
aresctl openapi --help
```

### Shell Autocompletion

Aresctl supports shell autocompletion for bash, zsh, fish, and PowerShell. This allows you to use Tab to autocomplete commands and flags.

**Setup for Bash:**
```bash
# Load completions in current session
source <(aresctl completion bash)

# Load completions for every new session (Linux)
aresctl completion bash > /etc/bash_completion.d/aresctl

# Load completions for every new session (macOS)
aresctl completion bash > $(brew --prefix)/etc/bash_completion.d/aresctl
```

**Setup for Zsh:**
```bash
# Load completions in current session
source <(aresctl completion zsh)

# Load completions for every new session (Linux)
aresctl completion zsh > "${fpath[1]}/_aresctl"

# Load completions for every new session (macOS)
aresctl completion zsh > $(brew --prefix)/share/zsh/site-functions/_aresctl
```

**Setup for Fish:**
```bash
aresctl completion fish | source

# To load completions for each session
aresctl completion fish > ~/.config/fish/completions/aresctl.fish
```

**Setup for PowerShell:**
```powershell
aresctl completion powershell | Out-String | Invoke-Expression

# To load completions for every new session
aresctl completion powershell > aresctl.ps1
# Then source this file from your PowerShell profile
```

After setup, you can use Tab to autocomplete:
```bash
aresctl <Tab>          # Shows: completion, generate, help, new, openapi
aresctl generate <Tab> # Shows: handler, crud
```

## Project Structure Requirements

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

Routes should be defined with comments for documentation:

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

## OpenAPI Generation Features

Aresctl automatically:

- ✅ Extracts route paths and HTTP methods
- ✅ Parses comments for API descriptions
- ✅ Generates request/response schemas from Go structs
- ✅ Infers request bodies for POST/PUT operations
- ✅ Infers response schemas based on handler names
- ✅ Groups endpoints by tags (from `@tag` annotations)
- ✅ Adds standard error responses (400, 404, 500)
- ✅ Supports path parameters (e.g., `/users/{id}`)

## Examples

### Complete Workflow

Create a new project and add features:

```bash
# 1. Create new project
aresctl new myshop
cd myshop

# 2. Generate CRUD for products
aresctl generate crud product

# 3. Add routes (manually edit internal/server/http.go)
# Add the routes shown in the generate output

# 4. Generate OpenAPI documentation
aresctl openapi

# 5. Start development
docker-compose -f deploy/docker-compose.yml up -d
go run main.go
```

### Quick Start with Template

```bash
# Use aresctl to create from template
aresctl new myproject
cd myproject

# Generate OpenAPI documentation
aresctl openapi

# View the generated specification
cat openapi.yaml
```

### Generate Multiple Resources

```bash
# Generate handlers for different resources
aresctl generate crud user
aresctl generate crud product
aresctl generate crud order

# Generate OpenAPI for all routes
aresctl openapi
```

## Ecosystem

Aresctl is part of the Ares ecosystem:

- **[Ares](https://github.com/xushuhui/ares)** - Lightweight Go web framework
- **[ares-layout](https://github.com/xushuhui/ares-layout)** - Production-ready application template
- **[aresctl](https://github.com/xushuhui/aresctl)** - Command-line tool (this project)

## Development

### Build from Source

```bash
git clone https://github.com/xushuhui/aresctl.git
cd aresctl
go build -o aresctl
```

### Run Tests

```bash
go test ./...
```

## Roadmap

Completed features:

- ✅ Project scaffolding (`aresctl new project-name`)
- ✅ Handler generation (`aresctl generate handler user`)
- ✅ CRUD generation (`aresctl generate crud user`)
- ✅ OpenAPI specification generation

Future features planned:

- [ ] Database migration tools
- [ ] Hot reload for development
- [ ] Configuration file support
- [ ] Interactive project setup wizard
- [ ] Custom template support

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License

## Links

- [Ares Framework](https://github.com/xushuhui/ares)
- [Ares Layout Template](https://github.com/xushuhui/ares-layout)
- [Documentation](https://github.com/xushuhui/ares)
- [Issues](https://github.com/xushuhui/aresctl/issues)
