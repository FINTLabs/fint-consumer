package setup

import (
	"github.com/codegangsta/cli"
	"github.com/FINTprosjektet/fint-consumer/generate"
	"github.com/FINTprosjektet/fint-consumer/common/github"
	"github.com/FINTprosjektet/fint-consumer/common/config"
	"fmt"
	"os"
	"github.com/FINTprosjektet/fint-consumer/common/utils"
	"log"
	"strings"
	"io/ioutil"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"time"
	"gopkg.in/src-d/go-git.v4"
)

func CmdSetupConsumer(c *cli.Context) {

	var tag string
	if c.GlobalString("tag") == config.DEFAULT_TAG {
		tag = github.GetLatest()
	} else {
		tag = c.GlobalString("tag")
	}
	force := c.GlobalBool("force")

	name := c.String("name")
	verfifyParameter(name, "Name parameter missing!")

	pkg := c.String("package")
	verfifyParameter(pkg, "Package parameter missing!")

	component := c.String("component")
	verfifyParameter(component, "Component parameter missing!")

	setupSkeleton(name)
	generate.Generate(tag, force)

	addModels(component, pkg, name)

	includePerson := c.Bool("includePerson")
	addPerson(includePerson, name)

	updateConfigFiles(name)

	reportNeedOfChanges(name)

	addModelToGradle(component, name)

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
}

func verfifyParameter(name string, message string) {
	if len(name) < 1 {
		fmt.Println(message)
		os.Exit(-1)
	}
}
func updateConfigFiles(name string) {
	models := getModels(name)
	writeConsumerPropsFile(getConsumerPropsClass(models), name)
	writeConstantsFile(getConstantsClass(name), name)
}
func addModels(component string, pkg string, name string) {
	src := fmt.Sprintf("%s/%s/%s/%s", utils.GetTempDirectory(), config.BASE_PATH, component, pkg)
	dest := fmt.Sprintf("./%s/src/main/java/no/fint/consumer/models/", getConsumerName(name))
	os.RemoveAll(dest)
	err := utils.CopyDir(src, dest)
	if err != nil {
		fmt.Println(err)
	}
}
func addPerson(includePerson bool, name string) {
	if includePerson {
		src := fmt.Sprintf("%s/%s/felles/person", utils.GetTempDirectory(), config.BASE_PATH)
		dest := fmt.Sprintf("./%s/src/main/java/no/fint/consumer/models/", getConsumerName(name))
		err := utils.CopyDir(src, dest)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func addModelToGradle(model string, name string) {
	m := fmt.Sprintf("compile('no.fint:fint-%s-model-java:+')", model)
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
