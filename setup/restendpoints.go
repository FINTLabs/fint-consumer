package setup

import (
	"strings"
	"bytes"
	"text/template"
	"fmt"
	"io/ioutil"
	"github.com/FINTprosjektet/fint-consumer/common/types"
)


func getRestEndpointsClass(m []types.Model) string {
	var funcMap = template.FuncMap{
		"ToLower": strings.ToLower,
		"ToUpper": strings.ToUpper,
		"ToTitle": strings.Title,
	}
	tpl := template.New("class").Funcs(funcMap)

	parse, err := tpl.Parse(RESTENDPOINTS_TEMPLATE)

	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	err = parse.Execute(&b, m)
	if err != nil {
		panic(err)
	}
	return b.String()
}

func writeRestEndpointsFile(content string, name string) {
	fmt.Println("  > Setup RestEndpoints.java")
	file := fmt.Sprintf("%s/src/main/java/no/fint/consumer/utils/RestEndpoints.java", getConsumerName(name))
	err := ioutil.WriteFile(file, []byte(content), 0777)
	if err != nil {
		fmt.Printf("Unable to write file: %s", err)
	}
}