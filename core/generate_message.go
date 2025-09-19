//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"
)

type Names []string
type Tags []string

type FieldDecl struct {
	Name       string
	Type       string
	Tag        string
	ShouldWrap bool
}

type StructDecl struct {
	Name       string
	Fields     []FieldDecl
	ShouldWrap bool
}

type StructDecls []StructDecl

func goTypeToTypeScript(goType string) string {
	// Handle common Go to TypeScript type mappings
	switch {
	case goType == "string":
		return "string"
	case goType == "bool":
		return "boolean"
	case goType == "int" || goType == "int8" || goType == "int16" || goType == "int32" || goType == "int64":
		return "number"
	case goType == "uint" || goType == "uint8" || goType == "uint16" || goType == "uint32" || goType == "uint64":
		return "number"
	case goType == "float32" || goType == "float64":
		return "number"
	case goType == "time.Time":
		return "string" // JSON-serialized time becomes string
	case goType == "uuid.UUID":
		return "string" // UUID becomes string in JSON
	case strings.HasPrefix(goType, "[]"):
		// Slice/array type
		innerType := strings.TrimPrefix(goType, "[]")
		return goTypeToTypeScript(innerType) + "[]"
	case strings.Contains(goType, "Unpacker"):
		// Handle unpacker types - strip the Unpacker suffix
		return strings.Replace(goType, "Unpacker", "", 1)
	default:
		// For custom types, assume they match interface names
		return goType
	}
}

func extractJsonTag(tag string) string {
	// Extract the json tag value, e.g., `json:"name"` -> "name"
	if tag == "" {
		return ""
	}
	re := regexp.MustCompile(`json:"([^"]+)"`)
	matches := re.FindStringSubmatch(tag)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func parseStruct(structType *ast.StructType, fset *token.FileSet) (decls []FieldDecl) {
	for _, field := range structType.Fields.List {
		if len(field.Names) != 1 {
			fmt.Println(field)
			fmt.Println(field.Names)
			panic("Unexpected")
		}

		var buf bytes.Buffer
		printer.Fprint(&buf, fset, field.Type)
		fmt.Println("Field type:", buf.String())

		re := regexp.MustCompile(`\bwrap\b`)
		matches := re.FindAllString(field.Comment.Text(), -1)
		shouldWrap := false
		if len(matches) > 0 {
			shouldWrap = true
		}

		tag := ""
		if field.Tag != nil {
			tag = field.Tag.Value
		}

		decls = append(decls, FieldDecl{
			Name:       field.Names[0].Name,
			Type:       buf.String(),
			Tag:        tag,
			ShouldWrap: shouldWrap,
		})
	}
	return decls
}

func main() {
	fset := token.NewFileSet()
	files := []string{os.Getenv("GOFILE")}

	typeIdentName := os.Args[1]
	fmt.Println(typeIdentName)

	// Collect type names
	var interfaceName string
	var interfaceImplementor string
	structDecls := make([]StructDecl, 0)

	var imports []string

	for _, fname := range files {
		f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}

		for _, decl := range f.Decls {
			// A generic declaration that is a type
			if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok.String() == "type" {
				// List of items under a type block (there can be more than one)
				for _, spec := range gd.Specs {
					typeSpec := spec.(*ast.TypeSpec)
					fmt.Println("Type: ", reflect.TypeOf(typeSpec.Type))

					switch v := typeSpec.Type.(type) {
					case *ast.ArrayType:
					case *ast.BadExpr:
					case *ast.BasicLit:
					case *ast.BinaryExpr:
					case *ast.CallExpr:
					case *ast.ChanType:
					case *ast.CompositeLit:
					case *ast.Ellipsis:
					case *ast.FuncLit:
					case *ast.FuncType:
					case *ast.Ident:
						// fmt.Println("Ident: ", v.Name)
						// typeIndicator.Name = typeSpec.Name.Name
						// typeIndicator.UnderlyingType = v
					case *ast.IndexExpr:
					case *ast.IndexListExpr:
					case *ast.InterfaceType:
						fmt.Println("Interface:", typeSpec.Name.Name)
						interfaceName = typeSpec.Name.Name
						for _, method := range v.Methods.List {
							interfaceImplementor = method.Names[0].Name
						}

					case *ast.KeyValueExpr:
					case *ast.MapType:
					case *ast.ParenExpr:
					case *ast.SelectorExpr:
					case *ast.SliceExpr:
					case *ast.StarExpr:
						break

					case *ast.StructType:
						fmt.Println("Parsing ", typeSpec.Name.Name)
						decls := parseStruct(v, fset)
						shouldWrap := false
						for _, decl := range decls {
							if decl.ShouldWrap {
								shouldWrap = true
							}
						}
						structDecls = append(structDecls, StructDecl{
							Name:       typeSpec.Name.Name,
							Fields:     decls,
							ShouldWrap: shouldWrap,
						})

					case *ast.TypeAssertExpr:
					case *ast.UnaryExpr:
						fmt.Println("something else")
					default:
						panic(fmt.Sprintf("unexpected ast.Expr: %#v", typeSpec.Type))
					}
				}
			} else if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok.String() == "import" {
				for _, spec := range gd.Specs {
					typeSpec := spec.(*ast.ImportSpec)
					if typeSpec.Name == nil {
						imports = append(imports, typeSpec.Path.Value)
					} else {
						imports = append(imports, typeSpec.Name.Name+" "+typeSpec.Path.Value)
					}
				}
			}
		}
		fmt.Println("list: ", structDecls)

	}

	// v--- nvm, this is what wrap does
	// we also need to unmarshal the structs themselves, in case
	// they have interfaces within them

	tmpl := template.Must(template.New("action").Funcs(template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}).Parse(registryTemplate))

	out, err := os.Create(strings.Split(os.Getenv("GOFILE"), ".")[0] + "_generated.go")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	err = tmpl.Execute(out, struct {
		InterfaceName string
		InterfaceImpl string
		Decls         StructDecls
		Imports       []string
	}{
		interfaceName,
		interfaceImplementor,
		structDecls,
		imports,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate TypeScript file
	tsTemplate := template.Must(template.New("typescript").Funcs(template.FuncMap{
		"upper":          strings.ToUpper,
		"lower":          strings.ToLower,
		"tsType":         goTypeToTypeScript,
		"extractJsonTag": extractJsonTag,
	}).Parse(typescriptTemplate))

	// Create TypeScript output file in the app/messaging directory
	baseName := strings.Split(filepath.Base(os.Getenv("GOFILE")), ".")[0]
	tsOutputPath := filepath.Join("..", "..", "app", "messaging", baseName+"_generated.ts")
	
	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(tsOutputPath), 0755); err != nil {
		fmt.Printf("Warning: Could not create TypeScript directory: %v\n", err)
		return
	}
	
	tsOut, err := os.Create(tsOutputPath)
	if err != nil {
		fmt.Printf("Warning: Could not create TypeScript file: %v\n", err)
		return
	}
	defer tsOut.Close()

	err = tsTemplate.Execute(tsOut, struct {
		InterfaceName string
		InterfaceImpl string
		Decls         StructDecls
		Imports       []string
	}{
		interfaceName,
		interfaceImplementor,
		structDecls,
		imports,
	})

	if err != nil {
		fmt.Println("TypeScript generation error:", err)
	} else {
		fmt.Printf("Generated TypeScript file: %s\n", tsOutputPath)
	}
}

const registryTemplate = `// Code generated by go generate; DO NOT EDIT.
package core

import (
    "encoding/json"
    "fmt"
    {{- range .Imports }}
    {{.}}
    {{- end }}
)

type {{.InterfaceName}}Type uint8

type {{.InterfaceName}}Unpacker struct {
    {{.InterfaceName}}
}

const (
{{- range .Decls }}
	{{upper .Name }} {{$.InterfaceName}}Type = iota
{{- end }}
)

{{- range .Decls }}
func ({{ .Name }}) {{ $.InterfaceImpl }}() {}
func (obj {{ .Name}}) MarshalJSON() ([]byte, error) {
    var raw struct {
		{{$.InterfaceName}}Type {{$.InterfaceName}}Type ` + "`" + `json:"{{lower $.InterfaceName}}_type"` + "`" + `
		{{- range .Fields}}
        {{.Name}} {{.Type}} {{.Tag}}
		{{- end}}
	}

    raw.{{$.InterfaceName}}Type = {{upper .Name}}
	{{- range .Fields}}
    raw.{{.Name}} = obj.{{.Name}}
	{{- end}}

    return json.Marshal(raw)
}
{{if .ShouldWrap }}
func (obj *{{ .Name}}) UnmarshalJSON(rawData []byte) error {
    var raw struct {
        {{$.InterfaceName}}Type {{$.InterfaceName}}Type ` + "`" + `json:"{{lower $.InterfaceName}}_type"` + "`" + `
		{{- range .Fields}}
        {{- if .ShouldWrap}}
        {{.Name}} {{.Type}}Unpacker {{.Tag}}
        {{- else}}
        {{.Name}} {{.Type}} {{.Tag}} {{- end}}
		{{- end}}
	}

	err := json.Unmarshal(rawData, &raw)

    {{- range .Fields}}
    {{- if .ShouldWrap}}
    obj.{{.Name}} = raw.{{.Name}}.{{.Type}}
    {{- else}}
    obj.{{.Name}} = raw.{{.Name}}
    {{- end}}
    {{- end}}

    return err
}
{{end}}
{{- end }}

func (msg *{{.InterfaceName}}Unpacker) Uncover() {{.InterfaceName}} {
	return msg.{{.InterfaceName}}
}

func (msg *{{.InterfaceName}}Unpacker) UnmarshalJSON(rawData []byte) error {
	var raw struct {
		{{$.InterfaceName}}Type {{$.InterfaceName}}Type ` + "`" + `json:"{{lower $.InterfaceName}}_type"` + "`" + `
	}

	if err := json.Unmarshal(rawData, &raw); err != nil {
		return err
	}

	switch raw.{{$.InterfaceName}}Type {
{{- range .Decls }}
	case {{upper .Name}}:
		message := {{.Name}}{}
		err := json.Unmarshal(rawData, &message)
		if err != nil {
			return err
		}
		msg.{{$.InterfaceName}} = message
{{- end }}
	default:
		return fmt.Errorf("unexpected type: %#v", raw.{{$.InterfaceName}}Type)
	}
	return nil
}

type {{ .InterfaceName -}}Handler[T any, E any] interface {
{{- range .Decls}}
    Handle{{.Name}}({{.Name}}, E) (T, error)
{{- end }}
}

func {{.InterfaceName}}Decode[T any, E any](handler {{.InterfaceName}}Handler[T, E], data {{.InterfaceName}}, extraData E) (ret T, err error) {
	switch v := data.(type) {
{{- range .Decls}}
    case {{.Name}}:
        return handler.Handle{{.Name}}(v, extraData)
{{- end }}
	default:
		return ret, fmt.Errorf("unexpected type: %#v", data)
	}
}
`

const typescriptTemplate = `// Code generated by go generate; DO NOT EDIT.

// {{.InterfaceName}} Union Type and Enum
export type {{.InterfaceName}} = {{range $i, $decl := .Decls}}{{if $i}} | {{end}}{{$decl.Name}}{{end}};

export enum {{.InterfaceName}}Type {
{{- range $i, $decl := .Decls}}
    {{$decl.Name}} = {{$i}},
{{- end}}
}

// Individual struct interfaces
{{- range .Decls}}

export interface {{.Name}} {
{{- range .Fields}}
    {{- $jsonTag := extractJsonTag .Tag}}
    {{- if $jsonTag}}
    {{$jsonTag}}: {{tsType .Type}};
    {{- else}}
    {{lower .Name}}: {{tsType .Type}};
    {{- end}}
{{- end}}
}
{{- end}}
`
