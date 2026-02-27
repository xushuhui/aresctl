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
	Name        string `yaml:"name"`
	In          string `yaml:"in"`
	Required    bool   `yaml:"required"`
	Description string `yaml:"description,omitempty"`
	Schema      Schema `yaml:"schema"`
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
	Type        string     `yaml:"type,omitempty"`
	Description string     `yaml:"description,omitempty"`
	Items       *SchemaRef `yaml:"items,omitempty"`
	Format      string     `yaml:"format,omitempty"`
	Ref         string     `yaml:"$ref,omitempty"`
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
		Paths:      buildPaths(routes, schemas),
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
						} else if strings.HasPrefix(part, "query:") {
							jsonTag = strings.Trim(strings.TrimPrefix(part, "query:"), `"`)
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
							prop.Ref = "#/components/schemas/" + ident.Name
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

func buildPaths(routes []Route, schemas map[string]SchemaObject) map[string]PathItem {
	paths := make(map[string]PathItem)

	for _, route := range routes {
		path := strings.ReplaceAll(route.Path, "{id}", "{id}")

		item := paths[path]
		op := buildOperation(route, schemas)

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

func buildOperation(route Route, schemas map[string]SchemaObject) *Operation {
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

	if route.Method == "GET" {
		queryParams := getQueryParametersForHandler(route.Handler, schemas)
		op.Parameters = append(op.Parameters, queryParams...)
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

func getQueryParametersForHandler(handler string, schemas map[string]SchemaObject) []Parameter {
	requestStructName := handler + "Request"

	if schema, exists := schemas[requestStructName]; exists {
		var params []Parameter

		for propName, propDef := range schema.Properties {
			param := Parameter{
				Name:     propName,
				In:       "query",
				Required: false,
				Schema:   Schema{Type: propDef.Type},
			}

			if propDef.Description != "" {
				param.Description = propDef.Description
			}

			params = append(params, param)
		}

		return params
	}

	return []Parameter{}
}

func inferResponseSchema(handler, method string) string {
	return handler + "Response"
}

func inferRequestSchema(handler, method string) string {
	if method != "POST" && method != "PUT" {
		return ""
	}

	return handler + "Request"
}
