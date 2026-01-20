package server

// This is an example routes file for demonstration purposes.
// Based on ares-layout project structure using github.com/xushuhui/ares framework.
//
// Example usage:
// import (
// 	"github.com/xushuhui/ares"
// )
//
// func NewHTTPServer(userHandler *handler.UserHandler) *ares.Ares {
// 	app := ares.Default()
//
// 	api := app.Group("/api")
// 	// 创建用户 @tag Users
// 	api.POST("/users", userHandler.Create)
// 	// 获取用户列表 @tag Users
// 	api.GET("/users", userHandler.List)
// 	// 获取用户详情 @tag Users
// 	api.GET("/users/{id}", userHandler.Get)
// 	// 更新用户 @tag Users
// 	api.PUT("/users/{id}", userHandler.Update)
// 	// 删除用户 @tag Users
// 	api.DELETE("/users/{id}", userHandler.Delete)
//
// 	return app
// }

// Example routes for OpenAPI generation
type Router interface {
	GET(path string, handler any)
	POST(path string, handler any)
	PUT(path string, handler any)
	DELETE(path string, handler any)
	Group(prefix string) Router
}

type UserHandler struct{}

func (h *UserHandler) Create(ctx any) error  { return nil }
func (h *UserHandler) List(ctx any) error    { return nil }
func (h *UserHandler) Get(ctx any) error     { return nil }
func (h *UserHandler) Update(ctx any) error  { return nil }
func (h *UserHandler) Delete(ctx any) error  { return nil }

func NewHTTPServer(userHandler *UserHandler) Router {
	var app Router // In real app: app := ares.Default()

	api := app.Group("/api")
	// 创建用户 @tag Users
	api.POST("/users", userHandler.Create)
	// 获取用户列表 @tag Users
	api.GET("/users", userHandler.List)
	// 获取用户详情 @tag Users
	api.GET("/users/{id}", userHandler.Get)
	// 更新用户 @tag Users
	api.PUT("/users/{id}", userHandler.Update)
	// 删除用户 @tag Users
	api.DELETE("/users/{id}", userHandler.Delete)

	return app
}
