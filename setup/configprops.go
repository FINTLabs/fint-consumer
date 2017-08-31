package setup

import (
	"io/ioutil"
	"github.com/FINTprosjektet/fint-consumer/common/types"
	"strings"
	"text/template"
	"bytes"
	"strconv"
	"fmt"
)

var funcMap = template.FuncMap{
	"ToLower": strings.ToLower,
	"ToUpper": strings.ToUpper,
	"GetInitialRate": func(i int) string {
		rate := (i * 10000) + 60000
		return strconv.Itoa(rate)
	},
}

func getModels(name string) []types.Model {
	files, _ := ioutil.ReadDir(fmt.Sprintf("%s/src/main/java/no/fint/consumer/models", getConsumerName(name)))

	var models = []types.Model{}
	for _, f := range files {
		models = append(models, types.Model{Name: f.Name()})
	}

	return models
}

func getConsumerPropsClass(m []types.Model) string {
	tpl := template.New("class").Funcs(funcMap)

	parse, err := tpl.Parse(CONSUMER_PROPS_TEMPLATE)

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

func writeConsumerPropsFile(content string, name string) {
	fmt.Println("  > Setup ConsumerProps.java")

	file := fmt.Sprintf("%s/src/main/java/no/fint/consumer/config/ConsumerProps.java", getConsumerName(name))
	err := ioutil.WriteFile(file, []byte(content), 0777)
	if err != nil {
		fmt.Printf("Unable to write file: %s", err)
	}
}
