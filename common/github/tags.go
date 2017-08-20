package github

import (
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
)

func GetTagList() []string {
	client := github.NewClient(nil)
	ctx := context.Background()
	var tagList []string

	opt := &github.ListOptions{}
	tags, _, err := client.Repositories.ListTags(ctx, GITHUB_OWNER, GITHUB_REPO, opt)

	if err != nil {
		fmt.Printf("Unable to get tag list from GitHub: %s", err)
	}

	for _, tag := range tags {
		tagList = append(tagList, tag.GetName())
	}

	return tagList
}

func GetLatest() string {
	client := github.NewClient(nil)
	ctx := context.Background()
	release, _, err := client.Repositories.GetLatestRelease(ctx, GITHUB_OWNER, GITHUB_REPO)

	if err != nil {
		fmt.Printf("Unable to get latest release from GitHub: %s", err)
	}

	return release.GetTagName()
}
