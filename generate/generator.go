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
	"GetIdentifikatorPackage": func(attr []types.Attribute, inh []types.InheritedAttribute) string {
		for _, a := range attr {
			if a.Type == "Identifikator" {
				return a.Package + "." + a.Type
			}
		}
		for _, a := range inh {
			if a.Type == "Identifikator" {
				return a.Package + "." + a.Type
			}
		}
		return "java.util.Random"
	},
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
		return strings.Replace(s, ".model.", ".model.resource.", -1)
	},
	"modelPkg": func(s string) string {
		l := strings.Split(s, ".")
		if len(l) <= 5 {
			return ""
		}
		return strings.Join(l[5:], ".") + "."
	},
}

func Generate(owner string, repo string, tag string, filename string, force bool, component string, pkg string, includePerson bool) []*types.Class {

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

	var classMap = make(map[string]*types.Class)

	classes, _, _, _ := parser.GetClasses(owner, repo, tag, filename, force)
	for _, c := range classes {
		classMap[c.Package+"."+c.Name] = c
	}

	var pkgName = "." + component
	if len(pkg) > 0 {
		pkgName = pkgName + "." + pkg
	}
	var resources []*types.Class
	var attributes = make(map[string]bool)

	for _, c := range classes {

		if (strings.Contains(c.Package, pkgName) && !c.Abstract && c.Identifiable) ||
			(includePerson && (c.Name == "Person" || c.Name == "Kontaktperson")) {
			fmt.Printf("  > Creating consumer package and classes for: %s.%s\n", c.Package, c.Name)

			setupPackagePath(c)

			writeClassFile(getLinkerClass(c), GetMainPackage(c.Package), c.Name, getLinkerClassFileName(c.Name))
			writeClassFile(getCacheServiceClass(c), GetMainPackage(c.Package), c.Name, getCacheServiceClassFileName(c.Name))
			writeClassFile(getControllerClass(c), GetMainPackage(c.Package), c.Name, getControllerClassFileName(c.Name))

			resources = append(resources, c)
			resources = append(resources, getLinks(c, classMap, attributes)...)
		}
	}

	fmt.Println("Finished generating Java code!")

	for _, c := range classes {
		_, ok := attributes[c.Name]
		if ok {
			resources = append(resources, c)
		}
	}

	return resources
}

func getLinks(c *types.Class, classMap map[string]*types.Class, seen map[string]bool) []*types.Class {
	var result []*types.Class
	for _, r := range c.Resources {
		name := r.Package + "." + r.Type
		if _, ok := seen[name]; !ok {
			if target, found := classMap[name]; found {
				//fmt.Printf("Need to link to %s\n", name)
				seen[name] = true
				result = append(result, target)
				result = append(result, getLinks(target, classMap, seen)...)
			}
		}
	}
	return result
}

func setupPackagePath(c *types.Class) {
	path := fmt.Sprintf("%s/%s/%s/%s", utils.GetTempDirectory(), config.BASE_PATH, GetMainPackage(c.Package), strings.ToLower(c.Name))
	fmt.Printf("    > Creating directory: %s\n", path)
	err := os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Println("Unable to create packages structure")
		fmt.Println(err)
	}
}
func GetMainPackage(path string) string {
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

func getLinkerClass(c *types.Class) string {
	return getClass(c, LINKER_TEMPLATE)
}

func getCacheServiceClass(c *types.Class) string {
	return getClass(c, CACHE_SERVICE_TEMPLATE)
}

func getControllerClass(c *types.Class) string {
	return getClass(c, CONTROLLER_TEMPLATE)
}

func getClass(c *types.Class, t string) string {
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
