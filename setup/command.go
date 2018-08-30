package setup

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/FINTprosjektet/fint-consumer/common/config"
	"github.com/FINTprosjektet/fint-consumer/common/github"
	"github.com/FINTprosjektet/fint-consumer/common/types"
	"github.com/FINTprosjektet/fint-consumer/common/utils"
	"github.com/FINTprosjektet/fint-consumer/generate"
	"github.com/codegangsta/cli"
)

func CmdSetupConsumer(c *cli.Context) {

	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		tag = github.GetLatest(c.GlobalString("owner"), c.GlobalString("repo"))
	} else {
		tag = c.GlobalString("tag")
	}
	force := c.GlobalBool("force")

	name := c.String("name")
	verfifyParameter(name, "Name parameter missing!")

	pkg := c.String("package")

	component := c.String("component")
	verfifyParameter(component, "Component parameter missing!")

	setupSkeleton(name)
	generate.Generate(c.GlobalString("owner"), c.GlobalString("repo"), tag, c.GlobalString("filename"), force)

	addModels(component, pkg, name)

	includePerson := c.Bool("includePerson")
	addPerson(includePerson, name)

	updateConfigFiles(name)

	reportNeedOfChanges(name)

	addModelToGradle(component, name)

	createGradleProperties(tag, name)

	createGradleSettings(name)

	createReadme(tag, pkg, component, name)

	/*
		r, _ := git.PlainInit(getConsumerName(name), false)
		w, _ := r.Worktree()

		w.Add(".gitignore")
		commit, _ := w.Commit("Initial commit", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "fint-provider cli",
				Email: "post@fintlabs.no",
				When:  time.Now(),
			},
		})

		obj, _ := r.CommitObject(commit)
		fmt.Println(obj)
	*/
}

func verfifyParameter(name string, message string) {
	if len(name) < 1 {
		fmt.Println(message)
		os.Exit(-1)
	}
}
func getModels(name string) []types.Model {
	files, _ := ioutil.ReadDir(fmt.Sprintf("%s/src/main/java/no/fint/consumer/models", getConsumerName(name)))

	var models = []types.Model{}
	for _, f := range files {
		if f.IsDir() {
			models = append(models, types.Model{Name: f.Name()})
		}
	}

	return models
}
func updateConfigFiles(name string) {
	models := getModels(name)
	writeConsumerPropsFile(getConsumerPropsClass(models), name)
	writeConstantsFile(getConstantsClass(name), name)
	writeLinkMapperFile(getLinkMapperClass(models), name)
	writeRestEndpointsFile(getRestEndpointsClass(models), name)
}
func addModels(component string, pkg string, name string) {
	src := fmt.Sprintf("%s/%s/%s/%s", utils.GetTempDirectory(), config.BASE_PATH, component, pkg)
	dest := fmt.Sprintf("./%s/src/main/java/no/fint/consumer/models/", getConsumerName(name))
	fmt.Printf("  > Copying models from %s to %s\n", src, dest)
	os.RemoveAll(dest)
	err := utils.CopyDir(src, dest)
	if err != nil {
		fmt.Println(err)
	}
}
func addPerson(includePerson bool, name string) {
	if includePerson {
		src := fmt.Sprintf("%s/%s/felles/person", utils.GetTempDirectory(), config.BASE_PATH)
		dest := fmt.Sprintf("./%s/src/main/java/no/fint/consumer/models/person/", getConsumerName(name))
		err := utils.CopyDir(src, dest)
		if err != nil {
			fmt.Println(err)
		}

		src = fmt.Sprintf("%s/%s/felles/kontaktperson", utils.GetTempDirectory(), config.BASE_PATH)
		dest = fmt.Sprintf("./%s/src/main/java/no/fint/consumer/models/kontaktperson/", getConsumerName(name))
		err = utils.CopyDir(src, dest)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func addModelToGradle(model string, name string) {
	m := fmt.Sprintf("    compile(\"no.fint:fint-%s-resource-model-java:${apiVersion}\")", model)
	gradleFile := utils.GetGradleFile(getConsumerName(name))
	input, err := ioutil.ReadFile(gradleFile)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "/* --> Models <-- */") {
			lines[i] = m
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(gradleFile, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func createReadme(tag string, pkg string, component string, name string) {
	content := fmt.Sprintf("# %s\n\nGenerated from tag `%s` on package `%s` and component `%s`.\n",
		getConsumerName(name), tag, pkg, component)
	err := ioutil.WriteFile(utils.GetReadmeFile(getConsumerName(name)), []byte(content), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func createGradleProperties(tag string, name string) {
	apiVersion := strings.TrimPrefix(tag, "v")
	content := fmt.Sprintf("version=0.0.0\napiVersion=%s\n", apiVersion)
	gradleProperties := fmt.Sprintf("%s/gradle.properties", utils.GetWorkingDir(getConsumerName(name)))
	err := ioutil.WriteFile(gradleProperties, []byte(content), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func createGradleSettings(name string) {
	content := fmt.Sprintf("rootProject.name = '%s'\n", name)
	gradleSettings := fmt.Sprintf("%s/settings.gradle", utils.GetWorkingDir(getConsumerName(name)))
	err := ioutil.WriteFile(gradleSettings, []byte(content), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
