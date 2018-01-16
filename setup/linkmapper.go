package setup

import (
	"strings"
	"bytes"
	"text/template"
	"fmt"
	"io/ioutil"
	"github.com/FINTprosjektet/fint-consumer/common/types"
)


func getLinkMapperClass(m []types.Model) string {
	var funcMap = template.FuncMap{
		"ToLower": strings.ToLower,
		"ToUpper": strings.ToUpper,
		"ToTitle": strings.Title,
	}
	tpl := template.New("class").Funcs(funcMap)

	parse, err := tpl.Parse(LINKMAPPER_TEMPLATE)

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

func writeLinkMapperFile(content string, name string) {
	fmt.Println("  > Setup LinkMapper.java")
	file := fmt.Sprintf("%s/src/main/java/no/fint/consumer/config/LinkMapper.java", getConsumerName(name))
	err := ioutil.WriteFile(file, []byte(content), 0777)
	if err != nil {
		fmt.Printf("Unable to write file: %s", err)
	}
}