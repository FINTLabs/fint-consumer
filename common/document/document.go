package document

import (
	"fmt"
	"os"

	"github.com/FINTprosjektet/fint-consumer/common/github"
	"github.com/antchfx/xquery/xml"
)

func Get(owner string, repo string, tag string, filename string, force bool) *xmlquery.Node {

	fileName := github.GetXMIFile(owner, repo, tag, filename, force)

	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	doc, err := xmlquery.Parse(f)
	if err != nil {
		fmt.Println(err)
	}
	return doc

}
