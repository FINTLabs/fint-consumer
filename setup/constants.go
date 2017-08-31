package setup

import (
	"bytes"
	"text/template"
	"fmt"
	"io/ioutil"
)

func getConstantsClass(name string) string {
	tpl := template.New("class")

	parse, err := tpl.Parse(CONSTANTS_TEMPLATE)

	if err != nil {
		panic(err)
	}
	m := map[string]string {
		"Name": name,
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