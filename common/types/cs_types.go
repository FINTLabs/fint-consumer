package types

import "strings"

var CS_TYPE_MAP = map[string]string{
	"string":   "string",
	"boolean":  "bool",
	"date":     "DateTime",
	"datetime": "DateTime",
	"double":   "double",
}

var CS_VALUE_TYPES = []string{
	"bool",
	"byte",
	"char",
	"decimal",
	"double",
	"float",
	"int",
	"long",
	"DateTime" }

func GetCSType(t string) string {

	value, ok := CS_TYPE_MAP[strings.ToLower(t)]
	if ok {
		return value
	} else {
		return t
	}
}

func IsValueType(t string) bool {
	for _, value := range CS_VALUE_TYPES {
		if t == value {
			return true
		}
	}
	return false
}
