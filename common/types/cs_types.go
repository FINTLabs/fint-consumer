package types

var CS_TYPE_MAP = map[string]string{
	"string":   "string",
	"boolean":  "bool",
	"date":     "DateTime",
	"dateTime": "DateTime",
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

	value, ok := CS_TYPE_MAP[t]
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