package setup

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"text/template"

	"github.com/FINTLabs/fint-consumer/common/types"
)

var funcMap = template.FuncMap{
	"ToLower": strings.ToLower,
	"ToUpper": strings.ToUpper,
	"GetInitialRate": func(i int) string {
		rate := (i * 10000) + 900000
		return strconv.Itoa(rate)
	},
}

func getConstantsClass(name string, models []types.Model) string {
	tpl := template.New("class").Funcs(funcMap)

	parse, err := tpl.Parse(CONSTANTS_TEMPLATE)

	if err != nil {
		panic(err)
	}

	m := struct {
		Name   string
		Models []types.Model
	}{
		name,
		models,
	}

	var b bytes.Buffer
	err = parse.Execute(&b, m)
	if err != nil {
		panic(err)
	}
	return b.String()
}

func writeConstantsFile(content string, name string) {
	fmt.Println("  > Setup Constants.java")
	file := fmt.Sprintf("%s/src/main/java/no/fint/consumer/config/Constants.java", getConsumerName(name))
	err := ioutil.WriteFile(file, []byte(content), 0777)
	if err != nil {
		fmt.Printf("Unable to write file: %s", err)
	}
}
