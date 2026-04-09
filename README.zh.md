# Aresctl

[English](README.md)

[Ares](https://github.com/xushuhui/ares) 框架的命令行工具。Aresctl 通过代码生成、OpenAPI 文档生成和 GORM 模型生成等功能，帮助你更快地开发 Ares 应用。

## 特性

- **项目脚手架** - 从模板创建新的 Ares 项目
- **代码生成** - 从模板生成 handler 和 CRUD 代码
- **OpenAPI 生成** - 自动从 Go 代码生成 OpenAPI 3.0 规范
- **GORM 代码生成** - 从数据库 schema 生成 GORM 模型和查询代码

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

### 生成代码

从模板生成 handler 或 CRUD 代码：

```bash
# 生成 handler
aresctl generate handler product

# 生成完整 CRUD（handler + API 类型）
aresctl generate crud product
```

### 生成 OpenAPI 规范

进入你的 Ares 项目目录并运行：

```bash
aresctl openapi
```

这将通过分析以下内容在项目根目录生成 `openapi.yaml` 文件：
- `internal/server/` 中的路由定义
- `api/` 中的 API 结构

### 生成 GORM 代码

在项目根目录运行：

```bash
aresctl gen
```

`gen` 命令从当前目录读取 `gorm.yaml` 配置文件。

最小 `gorm.yaml` 配置示例（仅支持 MySQL）：

```yaml
version: "0.1"
database:
  dsn: "user:pass@tcp(127.0.0.1:3306)/your_db?charset=utf8mb4&parseTime=True&loc=Local"
  db: "mysql"
  tables: []
  outPath: "./dao/query"
```

如果 `tables` 为空，则会为所有表生成代码。

## 项目结构

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

路由应使用注释和 `@tag` 注解定义：

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

## 帮助

```bash
aresctl --help
aresctl new --help
aresctl generate --help
aresctl openapi --help
aresctl gen --help
```

## 开发

```bash
# 从源码构建
git clone https://github.com/xushuhui/aresctl.git
cd aresctl
go build -o aresctl

# 运行测试
go test ./...
```

## 生态系统

- **[Ares](https://github.com/xushuhui/ares)** - 轻量级 Go Web 框架
- **[ares-layout](https://github.com/xushuhui/ares-layout)** - 生产级应用模板
- **[aresctl](https://github.com/xushuhui/aresctl)** - 命令行工具（本项目）

## 许可证

MIT License
