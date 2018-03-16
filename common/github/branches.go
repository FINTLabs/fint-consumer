package github

import (
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/net/context"
)

func GetBranchList(owner string, repo string) []string {
	client := github.NewClient(nil)
	ctx := context.Background()
	var branchList []string

	opt := &github.ListOptions{}
	branches, _, err := client.Repositories.ListBranches(ctx, owner, repo, opt)

	if err != nil {
		fmt.Printf("Unable to get branch list from GitHub: %s", err)
	}

	for _, b := range branches {
		branchList = append(branchList, b.GetName())
	}

	return branchList
}
