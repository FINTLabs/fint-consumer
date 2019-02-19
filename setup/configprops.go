package setup

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/FINTLabs/fint-consumer/common/types"
)

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
