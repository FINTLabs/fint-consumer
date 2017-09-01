package generate

import (
	"strings"
	"text/template"
	"fmt"
	"bytes"
	"os"
	"github.com/FINTprosjektet/fint-consumer/common/types"
	"github.com/FINTprosjektet/fint-consumer/common/config"
	"github.com/FINTprosjektet/fint-consumer/common/parser"
	"github.com/FINTprosjektet/fint-consumer/common/utils"
	"io/ioutil"
)

var funcMap = template.FuncMap{
	"ToLower": strings.ToLower,
	"ToUpper": strings.ToUpper,
	"GetActionPackage": func(p string) string {
		a := strings.Split(p, ".")
		action := a[len(a)-1]
		action = strings.Title(action) + "Actions"
		return fmt.Sprintf("%s.%s", p, action)
	},
	"GetAction": func(p string) string {
		a := strings.Split(p, ".")
		action := a[len(a)-1]
		action = strings.Title(action) + "Actions"
		return action
	},
}

func Generate(tag string, force bool) {

	//document.Get(tag, force)
	fmt.Println("Generating Java code:")

	fmt.Println("  > Setup directory structure.")
	tempGenerateDir := fmt.Sprintf("%s/%s", utils.GetTempDirectory(), config.BASE_PATH)
	os.RemoveAll(tempGenerateDir)
	err := os.MkdirAll(tempGenerateDir, 0777)
	if err != nil {
		fmt.Println("Unable to create base structure")
		fmt.Println(err)
	}

	classes, packageMap, _, _ := parser.GetClasses(tag, force)
	for _, c := range classes {

		if !c.Abstract && c.Identifiable {
			fmt.Printf("  > Creating consumer package and classes for: %s\n", fmt.Sprintf("%s.%s", c.Package, c.Name))

			setupPackagePath(packageMap, c)

			writeClassFile(getAssemblerClass(c), getMainPackage(packageMap[c.Name].Java), c.Name, getAssemblerClassFileName(c.Name))
			writeClassFile(getCacheServiceClass(c), getMainPackage(packageMap[c.Name].Java), c.Name, getCacheServiceClassFileName(c.Name))
			writeClassFile(getControllerClass(c), getMainPackage(packageMap[c.Name].Java), c.Name, getControllerClassFileName(c.Name))

		}

	}

	fmt.Println("Finish generating Java code!")

}

func setupPackagePath(packageMap map[string]types.Import, c types.Class) {
	path := fmt.Sprintf("%s/%s/%s/%s", utils.GetTempDirectory(), config.BASE_PATH, getMainPackage(packageMap[c.Name].Java), strings.ToLower(c.Name))
	fmt.Printf("    > Creating directory: %s\n", path)
	err := os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Println("Unable to create packages structure")
		fmt.Println(err)
	}
}
func getMainPackage(path string) string {
	a := strings.Split(path, ".")
	pkg := ""
	if len(a) == 6 {
		pkg = fmt.Sprintf("%s/%s", a[3], a[4])
	}
	if len(a) == 5 {
		pkg = a[3]
	}
	return pkg
}

func writeClassFile(content string, pkg string, name string, className string) {
	fmt.Printf("    > Creating class: %s\n", className)
	file := ""
	if len(pkg) > 0 {
		file = fmt.Sprintf("%s/%s/%s/%s/%s", utils.GetTempDirectory(), config.BASE_PATH, pkg, strings.ToLower(name), className)
	} else {
		file = fmt.Sprintf("%s/%s/%s/%s", utils.GetTempDirectory(), config.BASE_PATH, strings.ToLower(name), className)
	}
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
