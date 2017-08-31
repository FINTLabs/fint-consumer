package github

import (
	"fmt"
	"os"
	"gopkg.in/src-d/go-git.v4"
	"errors"
	"github.com/FINTprosjektet/fint-consumer/common/utils"
)

func Clone(name string, cloneUrl string) error  {

	r, err := git.PlainClone(utils.GetWorkingDir(name), false, &git.CloneOptions{
		URL:      cloneUrl,
		Progress: os.Stdout,
	})

	_, err = r.Head()

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to clone %s", cloneUrl))
	}

	return nil
}
