package document

import (
	"fmt"
	"github.com/FINTprosjektet/fint-consumer/common/github"
	"github.com/antchfx/xquery/xml"
	"os"
)

func Get(tag string, force bool) *xmlquery.Node {

	fileName := github.GetXMIFile(tag, force)

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
