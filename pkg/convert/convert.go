package convert

import (
	"fmt"
	"kcl-importer/pkg/jsonschema"
	"kcl-importer/pkg/kclschema"
)

func JsonSchemaToKclSchema(s jsonschema.Schema) kclschema.Schema {
	if s.SchemaType != jsonschema.SchemaTypeObject {
		panic("schema type is not object")
	}

	kclSch := kclschema.Schema{}

	for _, k := range s.OrderedKeywords {
		switch v := s.Keywords[k].(type) {
		case *jsonschema.Title:
			kclSch.Name = string(*v)
		case *jsonschema.ID:
			kclSch.Description += string(*v)
		case *jsonschema.Description:
			kclSch.Description += string(*v)
		case *jsonschema.Comment:
			kclSch.Description += string(*v)
		case *jsonschema.Type:
			if !v.StrVal || len(v.Vals) == 0 {
				fmt.Printf("unknown type: %#v\n", v)
				break
			}
			if v.Vals[0] != "object" {
				kclSch.Types = JsonTypesToKclTypes(v.Vals)
			}
		case *jsonschema.Properties:
			for key, val := range *v {
				propSch := JsonSchemaToKclSchema(*val)
				propSch.Name = key
				kclSch.Properties = append(kclSch.Properties, propSch)
			}
		default:
			fmt.Printf("unknown Keyword: %s\n\tValue: %#v\n", k, v)
		}
	}

	return kclSch
}

func JsonTypesToKclTypes(t []string) []kclschema.Type {
	var kclTypes []kclschema.Type
	for _, v := range t {
		kclTypes = append(kclTypes, JsonTypeToKclType(v))
	}
	return kclTypes
}

func JsonTypeToKclType(t string) kclschema.Type {
	switch t {
	case "string":
		return kclschema.Str
	case "boolean":
		return kclschema.Bool
	case "integer":
		return kclschema.Int
	case "number":
		return kclschema.Float
	default:
		fmt.Printf("unknown type: %s\n", t)
		return ""
	}
}
