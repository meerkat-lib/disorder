package generator

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/meerkat-lib/disorder/internal/schema"
	"github.com/meerkat-lib/disorder/internal/utils/strcase"
)

var goTypes = map[schema.Type]string{
	schema.TypeBool:      "bool",
	schema.TypeI8:        "int8",
	schema.TypeU8:        "uint8",
	schema.TypeI16:       "int16",
	schema.TypeU16:       "uint16",
	schema.TypeI32:       "int32",
	schema.TypeU32:       "uint32",
	schema.TypeI64:       "int64",
	schema.TypeU64:       "uint64",
	schema.TypeF32:       "float32",
	schema.TypeF64:       "float64",
	schema.TypeString:    "string",
	schema.TypeTimestamp: "int64",
	schema.TypeBytes:     "[]byte",
}

func goType(typ schema.Type, ref string) string {
	if typ.IsPrimary() {
		return goTypes[typ]
	}
	if strings.Contains(ref, ".") {
		names := strings.Split(ref, ".")
		pkg := strcase.SnakeCase(names[len(names)-2])
		obj := strcase.PascalCase(names[len(names)-1])
		if typ == schema.TypeEnum {
			return fmt.Sprintf("%s.%s", pkg, obj)
		}
		return fmt.Sprintf("*%s.%s", pkg, obj)
	}
	if typ == schema.TypeEnum {
		return strcase.PascalCase(ref)
	}
	return fmt.Sprintf("*%s", strcase.PascalCase(ref))
}

func NewGoGenerator() Generator {
	return newGeneratorImpl(&goLanguage{})
}

type goLanguage struct {
	goTemplate *template.Template
}

func (g *goLanguage) folder(pkg string) string {
	folders := strings.Split(pkg, ".")
	for i := range folders {
		folders[i] = strcase.SnakeCase(folders[i])
	}
	return strings.Join(folders, "/")
}

func (g *goLanguage) template() *template.Template {
	if g.goTemplate == nil {
		t := `// Code generated by https://github.com/meerkat-lib/disorder; DO NOT EDIT.
package {{PackageName .Package}}
{{- range $index, $enum := .Enums}}

type {{PascalCase $enum.Name}} string
const(
{{- range .Values}}
	{{PascalCase $enum.Name}}{{PascalCase .}} {{PascalCase $enum.Name}} = "{{.}}"
{{- end}}
)
{{- end}}
{{- range .Messages}}

type {{PascalCase .Name}} struct {
	{{- range .Fields}}
	{{PascalCase .Name}} {{Type .Type}}` +
			" `disorder:\"{{.Name}}\"`" + `{{- end}}
}
{{- end}}
`
		funcMap := template.FuncMap{
			"PascalCase": func(name string) string {
				return strcase.PascalCase(name)
			},
			"SnakeCase": func(name string) string {
				return strcase.SnakeCase(name)
			},
			"PackageName": func(pkg string) string {
				names := strings.Split(pkg, ".")
				return strcase.SnakeCase(names[len(names)-1])
			},
			"Type": func(typ *schema.TypeInfo) string {
				switch typ.Type {
				case schema.TypeArray:
					return fmt.Sprintf("[]%s", goType(typ.SubType, typ.TypeRef))
				case schema.TypeMap:
					return fmt.Sprintf("map[string]%s", goType(typ.SubType, typ.TypeRef))
				default:
					return goType(typ.Type, typ.TypeRef)
				}
			},
		}
		g.goTemplate = template.New("go").Funcs(funcMap)
		template.Must(g.goTemplate.Parse(t))
	}
	return g.goTemplate
}
