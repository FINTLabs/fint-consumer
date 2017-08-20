package generate

import (
	"bytes"
	"github.com/FINTprosjektet/fint-consumer/common/types"
	"text/template"
	"github.com/FINTprosjektet/fint-consumer/generate/java"
	"github.com/FINTprosjektet/fint-consumer/generate/cs"
)

var funcMap = template.FuncMap{
	//"add": func(i int, ii int) int { return i + ii },
	"sub":      func(i int, ii int) int { return i - ii },
	"javaType": types.GetJavaType,
	"csType":   types.GetCSType,
}

func GetJavaClass(c types.Class) string {
	return getClass(c, java.CLASS_TEMPLATE)
}

func GetJavaActionEnum(a types.Action) string  {
	return getActionEnum(a, java.ACTION_ENUM_TEMPLATE)
}

func GetCSActionEnum(a types.Action) string  {
	return getActionEnum(a, cs.ACTION_ENUM_TEMPLATE)
}

func GetCSClass(c types.Class) string {
	return getClass(c, cs.CLASS_TEMPLATE)
}

func getClass(c types.Class, t string) string {
	tpl := template.New("class").Funcs(funcMap)

	parse, err := tpl.Parse(t)

	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = parse.Execute(&b, c)
	if err != nil {
		panic(err)
	}
	return b.String()
}

func getActionEnum(a types.Action, t string) string {
	tpl := template.New("action_enum").Funcs(funcMap)

	parse, err := tpl.Parse(t)

	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = parse.Execute(&b, a)
	if err != nil {
		panic(err)
	}
	return b.String()
}
