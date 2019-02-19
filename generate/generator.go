package generate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/FINTLabs/fint-consumer/common/config"
	"github.com/FINTLabs/fint-consumer/common/parser"
	"github.com/FINTLabs/fint-consumer/common/types"
	"github.com/FINTLabs/fint-consumer/common/utils"
)

var funcMap = template.FuncMap{
	"ToLower": strings.ToLower,
	"ToUpper": strings.ToUpper,
	"ToTitle": strings.Title,
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
	"resourcePkg": func(s string) string {
		return strings.Replace(s, "model", "model.resource", -1)
	},
	"modelPkg": func(s string) string {
		l := strings.Split(s, ".")
		if len(l) <= 5 {
			return ""
		}
		return strings.Join(l[5:], ".") + "."
	},
}

func Generate(owner string, repo string, tag string, filename string, force bool) {

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

	classes, _, _, _ := parser.GetClasses(owner, repo, tag, filename, force)
	for _, c := range classes {

		if !c.Abstract && c.Identifiable {
			fmt.Printf("  > Creating consumer package and classes for: %s\n", fmt.Sprintf("%s.%s", c.Package, c.Name))

			setupPackagePath(c)

			writeClassFile(getLinkerClass(c), getMainPackage(c.Package), c.Name, getLinkerClassFileName(c.Name))
			writeClassFile(getCacheServiceClass(c), getMainPackage(c.Package), c.Name, getCacheServiceClassFileName(c.Name))
			writeClassFile(getControllerClass(c), getMainPackage(c.Package), c.Name, getControllerClassFileName(c.Name))

		}

	}

	fmt.Println("Finished generating Java code!")

}

func setupPackagePath(c types.Class) {
	path := fmt.Sprintf("%s/%s/%s/%s", utils.GetTempDirectory(), config.BASE_PATH, getMainPackage(c.Package), strings.ToLower(c.Name))
	fmt.Printf("    > Creating directory: %s\n", path)
	err := os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Println("Unable to create packages structure")
		fmt.Println(err)
	}
}
func getMainPackage(path string) string {
	a := strings.Split(path, ".")
	return strings.Join(a[3:], "/")
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

func getLinkerClassFileName(name string) string {
	return fmt.Sprintf("%sLinker.java", name)
}

func getCacheServiceClassFileName(name string) string {
	return fmt.Sprintf("%sCacheService.java", name)
}

func getControllerClassFileName(name string) string {
	return fmt.Sprintf("%sController.java", name)
}

func getLinkerClass(c types.Class) string {
	return getClass(c, LINKER_TEMPLATE)
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
