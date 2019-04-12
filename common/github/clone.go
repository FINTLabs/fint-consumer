package github

import (
	"errors"
	"fmt"
	"os"

	"github.com/FINTLabs/fint-consumer/common/utils"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func Clone(name string, cloneUrl string, reference string) error {

	ref := plumbing.ReferenceName(reference)

	r, err := git.PlainClone(utils.GetWorkingDir(name), false, &git.CloneOptions{
		URL:           cloneUrl,
		Progress:      os.Stdout,
		Depth:         1,
		ReferenceName: ref,
	})

	_, err = r.Head()

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to clone %s", cloneUrl))
	}

	return nil
}
