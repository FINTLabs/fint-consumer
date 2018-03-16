package packages

import (
	"github.com/FINTprosjektet/fint-consumer/common/parser"
	"github.com/FINTprosjektet/fint-consumer/common/utils"
)

func DistinctPackageList(owner string, repo string, tag string, filename string, force bool) []string {
	classes, _, _, _ := parser.GetClasses(owner, repo, tag, filename, force)

	var p []string
	for _, c := range classes {
		p = append(p, c.Package)
	}

	return utils.Distinct(p)
}
