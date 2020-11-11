package document

import (
	"fmt"
	"os"

	"github.com/FINTLabs/fint-consumer/common/github"
	xmlquery "github.com/antchfx/xquery/xml"
)

func Get(owner string, repo string, tag string, filename string, force bool) (*xmlquery.Node, error) {

	fileName := github.GetXMIFile(owner, repo, tag, filename, force)

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	doc, err := xmlquery.Parse(f)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return doc, nil

}
