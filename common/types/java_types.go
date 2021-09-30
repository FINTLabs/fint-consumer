package types

import "strings"

var JAVA_TYPE_MAP = map[string]string{
	"string":      "String",
	"boolean":     "Boolean",
	"date":        "Date",
	"datetime":    "Date",
	"float":       "Float",
	"double":      "Double",
	"long":        "Long",
	"int":         "Integer",
	"hovedklasse": "FintMainObject",
	"referanse":   "FintReference",
	"abstrakt":    "FintAbstractObject",
	"datatype":    "FintComplexDatatypeObject",
}

func GetJavaType(t string) string {

	value, ok := JAVA_TYPE_MAP[strings.ToLower(t)]
	if ok {
		return value
	} else {
		return t
	}
}
