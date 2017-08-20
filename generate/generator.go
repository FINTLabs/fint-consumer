package generate

import (
	"strings"
	"text/template"
	"fmt"
)

var funcMap = template.FuncMap{
	"ToLower": strings.ToLower,
	"ToUpper": strings.ToUpper,
	"GetActionPackage": func(p string) string {
		a := strings.Split(p, ".")
		action := a[len(a) - 1]
		action = strings.Title(action) + "Actions"
		return fmt.Sprintf("%s.%s", p, action)
	},
	"GetAction": func(p string) string {
		a := strings.Split(p, ".")
		action := a[len(a) - 1]
		action = strings.Title(action) + "Actions"
		return action
	},
}
