# Aresctl

[English](README.md)

[Ares](https://github.com/xushuhui/ares) 框架的命令行工具。Aresctl 通过代码生成、OpenAPI 文档生成等功能，帮助你更快地开发 Ares 应用。

## 特性

- 🚀 **OpenAPI 生成** - 自动从 Go 代码生成 OpenAPI 3.0 规范
- 📝 **代码分析** - 解析路由定义和 API 结构
- 🎯 **约定优于配置** - 与 [ares-layout](https://github.com/xushuhui/ares-layout) 项目结构无缝配合
- ⚡ **快速轻量** - 最小依赖，快速执行

## 安装

```bash
go install github.com/xushuhui/aresctl@latest
```

确保 `$GOPATH/bin` 在你的 `PATH` 中：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## 使用

### 创建新项目

从 ares-layout 模板创建新的 Ares 项目：

```bash
aresctl new myproject
```

这将：
- 克隆 ares-layout 模板
- 删除 .git 目录
- 使用你的项目名称更新 go.mod
- 初始化新的 git 仓库

**示例输出：**
```
Cloning ares-layout template...
✓ Successfully created project 'myproject'

Next steps:
  cd myproject
  docker-compose -f deploy/docker-compose.yml up -d
  go run main.go
```

### 生成代码

从模板生成 handler 或 CRUD 代码：

**生成 handler：**
```bash
aresctl generate handler product
```

**输出：**
```
Created: internal/handler/product.go
✓ Successfully generated handler for 'product'
```

**生成完整 CRUD（handler + API 类型）：**
```bash
aresctl generate crud product
```

**输出：**
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

### 生成 OpenAPI 规范

进入你的 Ares 项目目录并运行：

```bash
aresctl openapi
```

这将通过分析以下内容在项目根目录生成 `openapi.yaml` 文件：
- `internal/server/` 中的路由定义
- `api/` 中的 API 结构

**示例输出：**
```
✓ Generated openapi.yaml successfully
```

### 查看帮助

获取通用帮助：

```bash
aresctl --help
```

**输出：**
```
Aresctl is a powerful CLI tool that helps you develop applications with the Ares framework. It provides code generation, OpenAPI documentation, and more.

Usage:
  aresctl [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    Generate code from templates
  help        Help about any command
  new         Create a new Ares project
  openapi     Generate OpenAPI specification

Flags:
  -h, --help   help for aresctl

Use "aresctl [command] --help" for more information about a command.
```

获取命令特定帮助：

```bash
# 新项目帮助
aresctl new --help

# 生成代码帮助
aresctl generate --help

# OpenAPI 帮助
aresctl openapi --help
```

### Shell 自动补全

Aresctl 支持 bash、zsh、fish 和 PowerShell 的 shell 自动补全。这允许你使用 Tab 键自动补全命令和标志。

**Bash 设置：**
```bash
# 在当前会话中加载补全
source <(aresctl completion bash)

# 为每个新会话加载补全（Linux）
aresctl completion bash > /etc/bash_completion.d/aresctl

# 为每个新会话加载补全（macOS）
aresctl completion bash > $(brew --prefix)/etc/bash_completion.d/aresctl
```

**Zsh 设置：**
```bash
# 在当前会话中加载补全
source <(aresctl completion zsh)

# 为每个新会话加载补全（Linux）
aresctl completion zsh > "${fpath[1]}/_aresctl"

# 为每个新会话加载补全（macOS）
aresctl completion zsh > $(brew --prefix)/share/zsh/site-functions/_aresctl
```

**Fish 设置：**
```bash
aresctl completion fish | source

# 为每个会话加载补全
aresctl completion fish > ~/.config/fish/completions/aresctl.fish
```

**PowerShell 设置：**
```powershell
aresctl completion powershell | Out-String | Invoke-Expression

# 为每个新会话加载补全
aresctl completion powershell > aresctl.ps1
# 然后从你的 PowerShell 配置文件中引用此文件
```

设置后，你可以使用 Tab 键自动补全：
```bash
aresctl <Tab>          # 显示: completion, generate, help, new, openapi
aresctl generate <Tab> # 显示: handler, crud
```

## 项目结构要求

Aresctl 期望你的项目遵循 [ares-layout](https://github.com/xushuhui/ares-layout) 结构：

```
your-project/
├── internal/
│   └── server/
│       └── http.go          # 路由定义
├── api/
│   ├── user.go              # 请求/响应结构
│   └── response.go          # 通用响应类型
└── openapi.yaml             # 生成的文件
```

### 路由定义格式

路由应该使用注释来提供文档：

```go
func NewHTTPServer(userHandler *handler.UserHandler) *ares.Ares {
    app := ares.Default()

    api := app.Group("/api")
    // 创建用户 @tag Users
    api.POST("/users", userHandler.Create)
    // 获取用户列表 @tag Users
    api.GET("/users", userHandler.List)
    // 获取用户详情 @tag Users
    api.GET("/users/{id}", userHandler.Get)

    return app
}
```

### API 结构格式

使用注释定义你的 API 结构：

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

## OpenAPI 生成特性

Aresctl 自动：

- ✅ 提取路由路径和 HTTP 方法
- ✅ 解析注释作为 API 描述
- ✅ 从 Go 结构体生成请求/响应 schema
- ✅ 为 POST/PUT 操作推断请求体
- ✅ 根据处理器名称推断响应 schema
- ✅ 按标签分组端点（来自 `@tag` 注解）
- ✅ 添加标准错误响应（400、404、500）
- ✅ 支持路径参数（例如 `/users/{id}`）

## 示例

### 完整工作流

创建新项目并添加功能：

```bash
# 1. 创建新项目
aresctl new myshop
cd myshop

# 2. 为产品生成 CRUD
aresctl generate crud product

# 3. 添加路由（手动编辑 internal/server/http.go）
# 添加生成输出中显示的路由

# 4. 生成 OpenAPI 文档
aresctl openapi

# 5. 开始开发
docker-compose -f deploy/docker-compose.yml up -d
go run main.go
```

### 使用模板快速开始

```bash
# 使用 aresctl 从模板创建
aresctl new myproject
cd myproject

# 生成 OpenAPI 文档
aresctl openapi

# 查看生成的规范
cat openapi.yaml
```

### 生成多个资源

```bash
# 为不同资源生成 handler
aresctl generate crud user
aresctl generate crud product
aresctl generate crud order

# 为所有路由生成 OpenAPI
aresctl openapi
```

## 生态系统

Aresctl 是 Ares 生态系统的一部分：

- **[Ares](https://github.com/xushuhui/ares)** - 轻量级 Go Web 框架
- **[ares-layout](https://github.com/xushuhui/ares-layout)** - 生产级应用模板
- **[aresctl](https://github.com/xushuhui/aresctl)** - 命令行工具（本项目）

## 开发

### 从源码构建

```bash
git clone https://github.com/xushuhui/aresctl.git
cd aresctl
go build -o aresctl
```

### 运行测试

```bash
go test ./...
```

## 路线图

已完成的特性：

- ✅ 项目脚手架（`aresctl new project-name`）
- ✅ Handler 生成（`aresctl generate handler user`）
- ✅ CRUD 生成（`aresctl generate crud user`）
- ✅ OpenAPI 规范生成

计划中的未来特性：

- [ ] 数据库迁移工具
- [ ] 开发热重载
- [ ] 配置文件支持
- [ ] 交互式项目设置向导
- [ ] 自定义模板支持

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

MIT License

## 链接

- [Ares 框架](https://github.com/xushuhui/ares)
- [Ares Layout 模板](https://github.com/xushuhui/ares-layout)
- [文档](https://github.com/xushuhui/ares)
- [问题反馈](https://github.com/xushuhui/aresctl/issues)
