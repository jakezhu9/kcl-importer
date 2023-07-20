package kclschema

import (
	_ "embed"
	"io"
	"text/template"
)

type Generator struct {
	template  *template.Template
	functions *template.FuncMap
}

var (
	//go:embed templates/schema.gotmpl
	schemaTmpl string
)

func NewGenerator() *Generator {
	g := &Generator{
		template:  &template.Template{},
		functions: defaultFunctions(),
	}
	g.addDefaultTemplates()
	return g
}

func (g *Generator) addDefaultTemplates() {
	g.addTemplate("schema", schemaTmpl)
	g.template = g.template.Funcs(*g.functions)
}

func (g *Generator) addTemplate(name, data string) {
	tmpl := template.Must(template.New(name).Funcs(*g.functions).Parse(data))
	g.template = template.Must(g.template.AddParseTree(name, tmpl.Tree))
}

func defaultFunctions() *template.FuncMap {
	return &template.FuncMap{
		"formatTypes": func(t []Type) string {
			var result string
			for i, v := range t {
				if i > 0 {
					result += " | "
				}
				result += string(v)
			}
			return result
		},
	}
}

func (g *Generator) Generate(data Schema, writer io.Writer) error {
	return g.template.Execute(writer, data)
}
