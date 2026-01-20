package cmd

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type OpenAPI struct {
	OpenAPI    string              `yaml:"openapi"`
	Info       Info                `yaml:"info"`
	Paths      map[string]PathItem `yaml:"paths"`
	Components Components          `yaml:"components"`
}

type Info struct {
	Title   string `yaml:"title"`
	Version string `yaml:"version"`
}

type PathItem struct {
	Get    *Operation `yaml:"get,omitempty"`
	Post   *Operation `yaml:"post,omitempty"`
	Put    *Operation `yaml:"put,omitempty"`
	Delete *Operation `yaml:"delete,omitempty"`
}

type Operation struct {
	Summary     string              `yaml:"summary,omitempty"`
	Tags        []string            `yaml:"tags,omitempty"`
	Parameters  []Parameter         `yaml:"parameters,omitempty"`
	RequestBody *RequestBody        `yaml:"requestBody,omitempty"`
	Responses   map[string]Response `yaml:"responses"`
}

type Parameter struct {
	Name     string `yaml:"name"`
	In       string `yaml:"in"`
	Required bool   `yaml:"required"`
	Schema   Schema `yaml:"schema"`
}

type RequestBody struct {
	Required bool                       `yaml:"required"`
	Content  map[string]MediaTypeObject `yaml:"content"`
}

type MediaTypeObject struct {
	Schema SchemaRef `yaml:"schema"`
}

type Response struct {
	Description string                     `yaml:"description"`
	Content     map[string]MediaTypeObject `yaml:"content,omitempty"`
}

type Schema struct {
	Type string `yaml:"type"`
}

type SchemaRef struct {
	Ref string `yaml:"$ref,omitempty"`
}

type Components struct {
	Schemas map[string]SchemaObject `yaml:"schemas"`
}

type SchemaObject struct {
	Type        string                 `yaml:"type"`
	Description string                 `yaml:"description,omitempty"`
	Properties  map[string]PropertyDef `yaml:"properties,omitempty"`
}

type PropertyDef struct {
	Type        string     `yaml:"type"`
	Description string     `yaml:"description,omitempty"`
	Items       *SchemaRef `yaml:"items,omitempty"`
	Format      string     `yaml:"format,omitempty"`
}

type Route struct {
	Method  string
	Path    string
	Handler string
	Comment string
	Tag     string
}

func GenerateOpenAPI(routeDir, apiDir, outputFile string) {
	routes := parseRoutesFromDir(routeDir)
	schemas := parseSchemas(apiDir)

	openapi := OpenAPI{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:   "API Documentation",
			Version: "1.0.0",
		},
		Paths:      buildPaths(routes),
		Components: Components{Schemas: schemas},
	}

	data, _ := yaml.Marshal(openapi)
	os.WriteFile(outputFile, data, 0644)
}

func parseRoutesFromDir(dir string) []Route {
	var routes []Route

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return err
		}

		fileRoutes := parseRoutes(path)
		routes = append(routes, fileRoutes...)

		return nil
	})

	return routes
}

func parseRoutes(file string) []Route {
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, file, nil, parser.ParseComments)

	var routes []Route

	for _, decl := range node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		for _, stmt := range fn.Body.List {
			exprStmt, ok := stmt.(*ast.ExprStmt)
			if !ok {
				continue
			}

			call, ok := exprStmt.X.(*ast.CallExpr)
			if !ok {
				continue
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}

			method := strings.ToUpper(sel.Sel.Name)
			if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
				continue
			}

			if len(call.Args) >= 2 {
				if lit, ok := call.Args[0].(*ast.BasicLit); ok {
					path := strings.Trim(lit.Value, `"`)
					handler := ""
					if sel2, ok := call.Args[1].(*ast.SelectorExpr); ok {
						handler = sel2.Sel.Name
					}

					comment := ""
					tag := ""
					for _, cg := range node.Comments {
						if cg.End() < exprStmt.Pos() && fset.Position(cg.End()).Line == fset.Position(exprStmt.Pos()).Line-1 {
							text := strings.TrimSpace(strings.TrimPrefix(cg.Text(), "//"))
							if strings.Contains(text, "@tag") {
								parts := strings.SplitN(text, "@tag", 2)
								comment = strings.TrimSpace(parts[0])
								if len(parts) == 2 {
									tag = strings.TrimSpace(parts[1])
								}
							} else {
								comment = text
							}
							break
						}
					}

					routes = append(routes, Route{Method: method, Path: path, Handler: handler, Comment: comment, Tag: tag})
				}
			}
		}
	}
	return routes
}

func parseSchemas(dir string) map[string]SchemaObject {
	schemas := make(map[string]SchemaObject)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return err
		}

		fset := token.NewFileSet()
		node, _ := parser.ParseFile(fset, path, nil, parser.ParseComments)

		commentsByLine := make(map[int]string)
		for _, cg := range node.Comments {
			for _, comment := range cg.List {
				line := fset.Position(comment.Pos()).Line
				commentsByLine[line] = strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
			}
		}

		ast.Inspect(node, func(n ast.Node) bool {
			typeSpec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				return true
			}

			typeComment := ""

			if typeSpec.Doc != nil && len(typeSpec.Doc.List) > 0 {
				lastComment := typeSpec.Doc.List[len(typeSpec.Doc.List)-1]
				commentText := strings.TrimSpace(lastComment.Text)

				if strings.Contains(commentText, typeSpec.Name.Name) {
					parts := strings.SplitN(commentText, typeSpec.Name.Name, 2)
					if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
						typeComment = strings.TrimSpace(parts[1])
					}
				} else {
					typeComment = commentText
				}
			} else {
				typePos := fset.Position(typeSpec.Pos())
				var bestComment *ast.CommentGroup
				minDistance := 3

				for _, cg := range node.Comments {
					commentEnd := fset.Position(cg.End())
					distance := typePos.Line - commentEnd.Line

					if distance > 0 && distance <= minDistance {
						commentText := cg.Text()
						if strings.Contains(commentText, typeSpec.Name.Name) {
							bestComment = cg
							minDistance = distance
						}
					}
				}

				if bestComment != nil {
					commentText := strings.TrimSpace(bestComment.Text())
					if strings.Contains(commentText, typeSpec.Name.Name) {
						parts := strings.SplitN(commentText, typeSpec.Name.Name, 2)
						if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
							typeComment = strings.TrimSpace(parts[1])
						}
					}
				}
			}

			props := make(map[string]PropertyDef)
			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue
				}

				jsonTag := ""
				if field.Tag != nil {
					tag := strings.Trim(field.Tag.Value, "`")
					for _, part := range strings.Fields(tag) {
						if strings.HasPrefix(part, "json:") {
							jsonTag = strings.Trim(strings.TrimPrefix(part, "json:"), `"`)
							jsonTag = strings.Split(jsonTag, ",")[0]
						}
					}
				}

				if jsonTag == "" || jsonTag == "-" {
					continue
				}

				prop := PropertyDef{}

				if field.Comment != nil {
					prop.Description = strings.TrimSpace(field.Comment.Text())
				} else {
					fieldLine := fset.Position(field.End()).Line
					if comment, ok := commentsByLine[fieldLine]; ok {
						prop.Description = comment
					}
				}

				switch t := field.Type.(type) {
				case *ast.Ident:
					prop.Type = mapGoType(t.Name)
					if t.Name == "int64" {
						prop.Format = "int64"
					}
				case *ast.StarExpr:
					if ident, ok := t.X.(*ast.Ident); ok {
						if isBasicType(ident.Name) {
							prop.Type = mapGoType(ident.Name)
							if ident.Name == "int64" {
								prop.Format = "int64"
							}
						} else {
							prop.Type = "object"
							prop.Items = &SchemaRef{Ref: "#/components/schemas/" + ident.Name}
						}
					}
				case *ast.ArrayType:
					prop.Type = "array"
					if ident, ok := t.Elt.(*ast.Ident); ok {
						if ident.Name == "string" {
						} else {
							prop.Items = &SchemaRef{Ref: "#/components/schemas/" + ident.Name}
						}
					} else if starExpr, ok := t.Elt.(*ast.StarExpr); ok {
						if ident, ok := starExpr.X.(*ast.Ident); ok {
							prop.Items = &SchemaRef{Ref: "#/components/schemas/" + ident.Name}
						}
					}
				}

				props[jsonTag] = prop
			}

			if len(props) > 0 {
				schemas[typeSpec.Name.Name] = SchemaObject{
					Type:        "object",
					Description: typeComment,
					Properties:  props,
				}
			}
			return true
		})
		return nil
	})

	schemas["ErrorResponse"] = SchemaObject{
		Type: "object",
		Properties: map[string]PropertyDef{
			"code":    {Type: "integer", Description: "错误码"},
			"message": {Type: "string", Description: "错误信息"},
		},
	}

	return schemas
}

func mapGoType(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int32", "int64":
		return "integer"
	case "bool":
		return "boolean"
	default:
		return "string"
	}
}

func isBasicType(goType string) bool {
	switch goType {
	case "string", "int", "int32", "int64", "bool", "uint", "uint32", "uint64", "float32", "float64":
		return true
	default:
		return false
	}
}

func buildPaths(routes []Route) map[string]PathItem {
	paths := make(map[string]PathItem)

	for _, route := range routes {
		path := strings.ReplaceAll(route.Path, "{id}", "{id}")

		item := paths[path]
		op := buildOperation(route)

		switch route.Method {
		case "GET":
			item.Get = op
		case "POST":
			item.Post = op
		case "PUT":
			item.Put = op
		case "DELETE":
			item.Delete = op
		}

		paths[path] = item
	}

	return paths
}

func buildOperation(route Route) *Operation {
	op := &Operation{
		Summary: route.Comment,
		Responses: map[string]Response{
			"200": {Description: "Success"},
			"400": {Description: "Bad Request", Content: map[string]MediaTypeObject{
				"application/json": {Schema: SchemaRef{Ref: "#/components/schemas/ErrorResponse"}},
			}},
			"500": {Description: "Internal Error", Content: map[string]MediaTypeObject{
				"application/json": {Schema: SchemaRef{Ref: "#/components/schemas/ErrorResponse"}},
			}},
		},
	}

	if route.Comment == "" {
		op.Summary = route.Handler
	}

	if route.Tag != "" {
		op.Tags = []string{route.Tag}
	}

	if needsIdQueryParameter(route.Path, route.Handler) {
		op.Parameters = []Parameter{{
			Name:     "id",
			In:       "query",
			Required: true,
			Schema:   Schema{Type: "integer"},
		}}
	}

	requestSchema := inferRequestSchema(route.Handler, route.Method)
	if requestSchema != "" {
		op.RequestBody = &RequestBody{
			Required: true,
			Content: map[string]MediaTypeObject{
				"application/json": {Schema: SchemaRef{Ref: "#/components/schemas/" + requestSchema}},
			},
		}
	}

	responseSchema := inferResponseSchema(route.Handler, route.Method)
	if responseSchema != "" {
		op.Responses["200"] = Response{
			Description: "Success",
			Content: map[string]MediaTypeObject{
				"application/json": {Schema: SchemaRef{Ref: "#/components/schemas/" + responseSchema}},
			},
		}
	}

	if route.Method == "GET" && strings.HasPrefix(route.Handler, "Get") {
		op.Responses["404"] = Response{Description: "Not Found", Content: map[string]MediaTypeObject{
			"application/json": {Schema: SchemaRef{Ref: "#/components/schemas/ErrorResponse"}},
		}}
	}

	return op
}

func needsIdQueryParameter(path, handler string) bool {
	needsIdPaths := []string{
		"/help", "/expert", "/help/share", "/expert/share",
	}

	needsIdHandlers := []string{
		"GetHelp", "GetExpert", "GetHelpShare", "GetExpertShare",
	}

	for _, p := range needsIdPaths {
		if path == p {
			return true
		}
	}

	for _, h := range needsIdHandlers {
		if handler == h {
			return true
		}
	}

	return false
}

func inferResponseSchema(handler, method string) string {
	replyName := handler + "Reply"

	if strings.HasPrefix(handler, "List") {
		return replyName
	}

	if strings.HasPrefix(handler, "Get") || strings.HasPrefix(handler, "Create") {
		return replyName
	}

	if method == "POST" || method == "PUT" || method == "DELETE" {
		return replyName
	}

	return ""
}

func inferRequestSchema(handler, method string) string {
	if method != "POST" && method != "PUT" {
		return ""
	}

	switch handler {
	case "LoginPhone":
		return "LoginPhoneRequest"
	case "LoginMini":
		return "LoginMiniRequest"
	case "SendCode":
		return "SendCodeRequest"
	case "AuthorizePhoneMini":
		return "AuthorizePhoneMiniRequest"
	case "CreateHelp":
		return "CreateHelpRequest"
	default:
		return handler + "Request"
	}
}
