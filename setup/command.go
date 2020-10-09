package setup

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/FINTLabs/fint-consumer/common/config"
	"github.com/FINTLabs/fint-consumer/common/github"
	"github.com/FINTLabs/fint-consumer/common/types"
	"github.com/FINTLabs/fint-consumer/common/utils"
	"github.com/FINTLabs/fint-consumer/generate"
	"github.com/urfave/cli"
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

	var ref string
	if c.String("tag") != "" {
		ref = "refs/tags/" + c.String("tag")
	} else {
		ref = "refs/heads/" + c.String("branch")
	}

	version := c.String("version")

	setupSkeleton(name, ref)
	resources := generate.Generate(c.GlobalString("owner"), c.GlobalString("repo"), tag, c.GlobalString("filename"), force, component, pkg)
	sort.Sort(types.ByName(resources))

	addModels(component, pkg, name)

	includePerson := c.Bool("includePerson")
	addPerson(includePerson, name)

	updateConfigFiles(component, pkg, name, resources)

	reportNeedOfChanges(name)

	addModelToGradle(component, name)

	createGradleProperties(tag, name, version)

	createGradleSettings(name)

	createReadme(c.App.Name, c.App.Version, tag, pkg, component, name, ref)

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
	var models = []types.Model{}

	walkModels := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, "CacheService.java") {
			name := strings.ToLower(strings.TrimSuffix(filepath.Base(path), "CacheService.java"))
			l := strings.Split(filepath.Dir(path), "/")
			pkg := strings.Join(l[4:], ".")
			models = append(models, types.Model{Name: name, Package: pkg})
			return filepath.SkipDir
		}
		return nil
	}
	err := filepath.Walk(fmt.Sprintf("%s/src/main/java/no/fint/consumer/models", getConsumerName(name)), walkModels)
	if err != nil {
		fmt.Println(err)
	}
	return models
}
func getModelsFromResources(resources []*types.Class) []types.Model {
	var models = []types.Model{}
	for _, c := range resources {
		models = append(models, types.Model{Name: c.Name, Package: c.Package})
	}
	return models
}
func getAssociationsFromResources(resources []*types.Class) []types.Association {
	var assocs = []types.Association{}
	var exists = make(map[string]struct{})
	for _, c := range resources {
		for _, rel := range c.Relations {
			_, ok := exists[rel.Target]
			if !ok {
				assocs = append(assocs, rel)
				exists[rel.Target] = struct{}{}
			}
		}
	}
	return assocs
}
func updateConfigFiles(component string, pkg string, name string, resources []*types.Class) {
	models := getModelsFromResources(resources)
	assocs := getAssociationsFromResources(resources)
	/*writeConsumerPropsFile(getConsumerPropsClass(models), name)*/
	writeConstantsFile(getConstantsClass(name, models), name)
	writeLinkMapperFile(getLinkMapperClass(component, pkg, models, assocs), name)
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

func createReadme(app string, version string, tag string, pkg string, component string, name string, ref string) {
	content := fmt.Sprintf("# %s\n\n"+
		"Generated by %s `%s` from tag `%s` on package `%s` and component `%s`.\n\n"+
		"Based on fint-consumer-skeleton %s.\n",
		getConsumerName(name), app, version, tag, pkg, component, ref)
	err := ioutil.WriteFile(utils.GetReadmeFile(getConsumerName(name)), []byte(content), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func createGradleProperties(tag string, name string, version string) {
	apiVersion := strings.TrimPrefix(tag, "v")
	content := fmt.Sprintf("version=%s\napiVersion=%s\n", version, apiVersion)
	gradleProperties := fmt.Sprintf("%s/gradle.properties", utils.GetWorkingDir(getConsumerName(name)))
	err := ioutil.WriteFile(gradleProperties, []byte(content), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func createGradleSettings(name string) {
	content := fmt.Sprintf("rootProject.name = 'fint-consumer-%s'\n", name)
	gradleSettings := fmt.Sprintf("%s/settings.gradle", utils.GetWorkingDir(getConsumerName(name)))
	err := ioutil.WriteFile(gradleSettings, []byte(content), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
