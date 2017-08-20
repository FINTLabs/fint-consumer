package generate

import (
	"fmt"
	"github.com/FINTprosjektet/fint-consumer/common/config"
	"github.com/FINTprosjektet/fint-consumer/common/document"
	"github.com/FINTprosjektet/fint-consumer/common/github"
	"github.com/FINTprosjektet/fint-consumer/common/parser"
	"github.com/FINTprosjektet/fint-consumer/namespaces"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"strings"
	"github.com/FINTprosjektet/fint-consumer/common/types"
	"github.com/FINTprosjektet/fint-consumer/packages"
)

func CmdGenerate(c *cli.Context) {

	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		tag = github.GetLatest()
	} else {
		tag = c.GlobalString("tag")
	}
	force := c.GlobalBool("force")

	if c.String("lang") == "JAVA" {
		generateJavaCode(tag, force)
	}

	if c.String("lang") == "CS" {
		generateCSCode(tag, force)
	}

	if c.String("lang") == "ALL" {
		generateCSCode(tag, force)
		generateJavaCode(tag, force)
	}
}

func generateJavaCode(tag string, force bool) {

	document.Get(tag, force)
	fmt.Println("Generating Java code:")
	setupJavaDirStructure(tag, force)
	classes, _, packageClassMap, _ := parser.GetClasses(tag, force)
	for _, c := range classes {
		fmt.Printf("  > Creating class: %s.java\n", c.Name)
		class := GetJavaClass(c)

		path := fmt.Sprintf("%s/%s/%s.java", config.JAVA_BASE_PATH, strings.Replace(c.Package, ".", "/", -1), c.Name)
		err := ioutil.WriteFile(removeJavaPackagePathFromFilePath(path), []byte(class), 0777)
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}

	}

	for p, cl := range packageClassMap {
		action := getAction(p, cl)
		fmt.Printf("  > Creating action: %s.java\n", action.Name)
		actionEnum := GetJavaActionEnum(action)
		path := fmt.Sprintf("%s/%s/%s.java", config.JAVA_BASE_PATH, strings.Replace(p, ".", "/", -1), action.Name)
		err := ioutil.WriteFile(removeJavaPackagePathFromFilePath(path), []byte(actionEnum), 0777)
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}

	}

	fmt.Println("Finish generating Java code!")
}

func removeJavaPackagePathFromFilePath(path string) string {
	return strings.Replace(path, "no/fint/model/", "", -1)
}

func getAction(p string, cl []types.Class) types.Action {
	var action types.Action

	packageList := strings.Split(p, ".")
	pkg := packageList[len(packageList)-1]
	action.Name = fmt.Sprintf("%sActions", strings.Title(pkg))

	action.Package = p
	action.Namespace = p

	for _, c := range cl {
		if c.Identifiable && !c.Abstract {
			action.Classes = append(action.Classes, strings.ToUpper(c.Name))
		}
	}
	return action
}

func generateCSCode(tag string, force bool) {

	document.Get(tag, force)
	fmt.Println("Generating CSharp code:")
	setupCSDirStructure(tag, force)
	classes, _, _, packageClassMap := parser.GetClasses(tag, force)
	for _, c := range classes {
		fmt.Printf("  > Creating class: %s.cs\n", c.Name)

		class := GetCSClass(c)

		path := fmt.Sprintf("%s/%s.cs", getCSPath(c.Namespace), c.Name)
		err := ioutil.WriteFile(path, []byte(class), 0777)
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}

	}

	for p, cl := range packageClassMap {
		action := getAction(p, cl)
		fmt.Printf("  > Creating action: %s.cs\n", action.Name)
		actionEnum := GetCSActionEnum(action)
		path := fmt.Sprintf("%s/%s.cs", getCSPath(p), action.Name)
		err := ioutil.WriteFile(path, []byte(actionEnum), 0777)
		if err != nil {
			fmt.Printf("Unable to write file: %s", err)
		}

	}

	fmt.Println("Finish generating CSharp code!")

}

func setupCSDirStructure(tag string, force bool) {
	fmt.Println("  > Setup directory structure.")
	os.RemoveAll("net")
	err := os.MkdirAll(config.CS_BASE_PATH, 0777)
	if err != nil {
		fmt.Println("Unable to create base structure")
		fmt.Println(err)
	}
	for _, ns := range namespaces.DistinctNamespaceList(tag, force) {
		path := getCSPath(ns)
		err := os.MkdirAll(path, 0777)
		if err != nil {
			fmt.Println("Unable to create namespace structure")
			fmt.Println(err)
		}
	}
}
func getCSPath(ns string) string {
	nsList := strings.Split(ns, ".")
	projectDir := fmt.Sprintf("%s.%s.%s", nsList[0], nsList[1], nsList[2])
	subDirs := strings.Replace(ns, projectDir, "", -1)
	subDirs = strings.Replace(subDirs, ".", "/", -1)
	path := fmt.Sprintf("%s/%s/%s", config.CS_BASE_PATH, projectDir, subDirs)
	return path
}

func setupJavaDirStructure(tag string, force bool) {
	fmt.Println("  > Setup directory structure.")
	os.RemoveAll("java")
	err := os.MkdirAll(config.JAVA_BASE_PATH, 0777)
	if err != nil {
		fmt.Println("Unable to create base structure")
		fmt.Println(err)
	}

	for _, pkg := range packages.DistinctPackageList(tag, force) {
		path := fmt.Sprintf("%s/%s", config.JAVA_BASE_PATH, strings.Replace(pkg, ".", "/", -1))
		err := os.MkdirAll(removeJavaPackagePathFromFilePath(path), 0777)
		if err != nil {
			fmt.Println("Unable to create packages structure")
			fmt.Println(err)
		}
	}

}
