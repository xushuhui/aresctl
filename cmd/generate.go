package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [type] [name]",
	Short: "Generate code from templates",
	Long:  `Generate code from templates. Supported types: handler, crud`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		genType := args[0]
		name := args[1]

		if err := generateCode(genType, name); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Successfully generated %s for '%s'\n", genType, name)
	},
}

func generateCode(genType, name string) error {
	switch genType {
	case "handler":
		return generateHandler(name)
	case "crud":
		return generateCRUD(name)
	default:
		return fmt.Errorf("unsupported type: %s (supported: handler, crud)", genType)
	}
}

// capitalize capitalizes the first letter of a string
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func generateHandler(name string) error {
	// Capitalize first letter for type names
	typeName := capitalize(name)

	// Create handler file
	handlerDir := "internal/handler"
	if err := os.MkdirAll(handlerDir, 0755); err != nil {
		return fmt.Errorf("failed to create handler directory: %w", err)
	}

	handlerFile := filepath.Join(handlerDir, strings.ToLower(name)+".go")
	if _, err := os.Stat(handlerFile); !os.IsNotExist(err) {
		return fmt.Errorf("handler file already exists: %s", handlerFile)
	}

	tmpl := `package handler

import (
	"strconv"

	"github.com/xushuhui/ares"
)

type {{.TypeName}}Handler struct {
	// Add your dependencies here
}

func New{{.TypeName}}Handler() *{{.TypeName}}Handler {
	return &{{.TypeName}}Handler{}
}

func (h *{{.TypeName}}Handler) Create(ctx *ares.Context) error {
	// TODO: Implement create logic
	return ctx.JSON(200, map[string]string{"message": "{{.TypeName}} created"})
}

func (h *{{.TypeName}}Handler) Get(ctx *ares.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid id"})
	}

	// TODO: Implement get logic
	return ctx.JSON(200, map[string]interface{}{"id": id, "message": "{{.TypeName}} details"})
}

func (h *{{.TypeName}}Handler) List(ctx *ares.Context) error {
	// TODO: Implement list logic
	return ctx.JSON(200, map[string]interface{}{"items": []interface{}{}, "total": 0})
}

func (h *{{.TypeName}}Handler) Update(ctx *ares.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid id"})
	}

	// TODO: Implement update logic
	return ctx.JSON(200, map[string]interface{}{"id": id, "message": "{{.TypeName}} updated"})
}

func (h *{{.TypeName}}Handler) Delete(ctx *ares.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid id"})
	}

	// TODO: Implement delete logic
	return ctx.JSON(200, map[string]string{"message": "{{.TypeName}} deleted"})
}
`

	t, err := template.New("handler").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	f, err := os.Create(handlerFile)
	if err != nil {
		return fmt.Errorf("failed to create handler file: %w", err)
	}
	defer f.Close()

	data := map[string]string{
		"TypeName": typeName,
	}

	if err := t.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	fmt.Printf("Created: %s\n", handlerFile)
	return nil
}

func generateCRUD(name string) error {
	typeName := capitalize(name)
	lowerName := strings.ToLower(name)

	// Generate handler
	if err := generateHandler(name); err != nil {
		return err
	}

	// Generate API types
	apiDir := "api"
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return fmt.Errorf("failed to create api directory: %w", err)
	}

	apiFile := filepath.Join(apiDir, lowerName+".go")
	if _, err := os.Stat(apiFile); !os.IsNotExist(err) {
		return fmt.Errorf("api file already exists: %s", apiFile)
	}

	apiTmpl := `package api

// Create{{.TypeName}}Request 创建{{.TypeName}}请求
type Create{{.TypeName}}Request struct {
	Name string ` + "`json:\"name\"`" + ` // 名称
}

// Create{{.TypeName}}Response 创建{{.TypeName}}响应
type Create{{.TypeName}}Response struct {
	ID int64 ` + "`json:\"id\"`" + ` // ID
}

// Update{{.TypeName}}Request 更新{{.TypeName}}请求
type Update{{.TypeName}}Request struct {
	Name string ` + "`json:\"name\"`" + ` // 名称
}

// {{.TypeName}}Response {{.TypeName}}响应
type {{.TypeName}}Response struct {
	ID   int64  ` + "`json:\"id\"`" + `   // ID
	Name string ` + "`json:\"name\"`" + ` // 名称
}

// List{{.TypeName}}Response {{.TypeName}}列表响应
type List{{.TypeName}}Response struct {
	Total int64                ` + "`json:\"total\"`" + ` // 总数
	List  []*{{.TypeName}}Response ` + "`json:\"list\"`" + `  // 列表
}
`

	t, err := template.New("api").Parse(apiTmpl)
	if err != nil {
		return fmt.Errorf("failed to parse api template: %w", err)
	}

	f, err := os.Create(apiFile)
	if err != nil {
		return fmt.Errorf("failed to create api file: %w", err)
	}
	defer f.Close()

	data := map[string]string{
		"TypeName": typeName,
	}

	if err := t.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute api template: %w", err)
	}

	fmt.Printf("Created: %s\n", apiFile)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("1. Add routes in internal/server/http.go:\n")
	fmt.Printf("   api.POST(\"/%s\", %sHandler.Create)\n", lowerName+"s", lowerName)
	fmt.Printf("   api.GET(\"/%s\", %sHandler.List)\n", lowerName+"s", lowerName)
	fmt.Printf("   api.GET(\"/%s/{id}\", %sHandler.Get)\n", lowerName+"s", lowerName)
	fmt.Printf("   api.PUT(\"/%s/{id}\", %sHandler.Update)\n", lowerName+"s", lowerName)
	fmt.Printf("   api.DELETE(\"/%s/{id}\", %sHandler.Delete)\n", lowerName+"s", lowerName)
	fmt.Printf("2. Implement business logic in internal/biz/\n")
	fmt.Printf("3. Implement data access in internal/data/\n")

	return nil
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
