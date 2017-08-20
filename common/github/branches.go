package github

import (
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
)

func GetBranchList() []string {
	client := github.NewClient(nil)
	ctx := context.Background()
	var branchList []string

	opt := &github.ListOptions{}
	branches, _, err := client.Repositories.ListBranches(ctx, GITHUB_OWNER, GITHUB_REPO, opt)

	if err != nil {
		fmt.Printf("Unable to get branch list from GitHub: %s", err)
	}

	for _, b := range branches {
		branchList = append(branchList, b.GetName())
	}

	return branchList
}
