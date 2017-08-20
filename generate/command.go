package generate

import (
	"github.com/codegangsta/cli"
	"fmt"
	"github.com/FINTprosjektet/fint-consumer/common/document"
	"github.com/FINTprosjektet/fint-consumer/common/parser"
	"bytes"
	"github.com/FINTprosjektet/fint-consumer/common/types"
	"text/template"
	"github.com/FINTprosjektet/fint-consumer/common/config"
	"github.com/FINTprosjektet/fint-consumer/common/github"
	"os"
	"strings"
	"io/ioutil"
)

func CmdGenerate(c *cli.Context) {
	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		tag = github.GetLatest()
	} else {
		tag = c.GlobalString("tag")
	}
	force := c.GlobalBool("force")

	generate(tag, force)

}
func generate(tag string, force bool) {

	document.Get(tag, force)
	fmt.Println("Generating Java code:")

	fmt.Println("  > Setup directory structure.")
	os.RemoveAll("java")
	err := os.MkdirAll(config.BASE_PATH, 0777)
	if err != nil {
		fmt.Println("Unable to create base structure")
		fmt.Println(err)
	}


	classes, _, _, _ := parser.GetClasses(tag, force)
	for _, c := range classes {

		if !c.Abstract && c.Identifiable {
			fmt.Printf("  > Creating class: %s.java\n", c.Name)

			path := fmt.Sprintf("%s/%s", config.BASE_PATH, strings.ToLower(c.Name))
			err := os.MkdirAll(path, 0777)
			if err != nil {
				fmt.Println("Unable to create packages structure")
				fmt.Println(err)
			}


			writeClassFile(getAssemblerClass(c), c.Name, getAssemblerClassFileName(c.Name))
			writeClassFile(getCacheServiceClass(c), c.Name, getCacheServiceClassFileName(c.Name) )
			writeClassFile(getControllerClass(c), c.Name, getControllerClassFileName(c.Name))

		}

	}

}

func writeClassFile(content string, name string, className string) {
	file := fmt.Sprintf("%s/%s/%s", config.BASE_PATH, strings.ToLower(name), className)
	err := ioutil.WriteFile(file, []byte(content), 0777)
	if err != nil {
		fmt.Printf("Unable to write file: %s", err)
	}
}

func getAssemblerClassFileName(name string) string {
	return fmt.Sprintf("%sAssembler.java", name)
}

func getCacheServiceClassFileName(name string) string {
	return fmt.Sprintf("%sCacheService.java", name)
}

func getControllerClassFileName(name string) string {
	return fmt.Sprintf("%sController.java", name)
}

func getAssemblerClass(c types.Class) string {
	return getClass(c, RESOURCE_ASSEMBLER_TEMPLATE)
}

func getCacheServiceClass(c types.Class) string {
	return getClass(c, CACHE_SERVICE_TEMPLATE)
}

func getControllerClass(c types.Class) string {
	return getClass(c, CONTROLLER_TEMPLATE)
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
