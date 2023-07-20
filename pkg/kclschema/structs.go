package kclschema

type Schema struct {
	Name        string
	Description string

	Types      []Type
	Properties []Schema
}

type Type string

const (
	Str   Type = "str"
	Bool  Type = "bool"
	Int   Type = "int"
	Float Type = "float"
)
