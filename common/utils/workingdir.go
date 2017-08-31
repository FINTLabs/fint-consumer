package utils

import "fmt"

func GetWorkingDir(name string) string {
	return fmt.Sprintf("./%s", name)
}

func GetDotGitDir(name string) string {
	return fmt.Sprintf("%s/.git", GetWorkingDir(name))
}

func GetGradleFile(name string) string {
	return fmt.Sprintf("%s/build.gradle", GetWorkingDir(name))
}
