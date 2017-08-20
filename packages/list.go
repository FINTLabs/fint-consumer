package packages

import (
	"github.com/FINTprosjektet/fint-consumer/common/parser"
	"github.com/FINTprosjektet/fint-consumer/common/utils"
)

func DistinctPackageList(tag string, force bool) []string {
	classes, _, _, _ := parser.GetClasses(tag, force)

	var p []string
	for _, c := range classes {
		p = append(p, c.Package)
	}

	return utils.Distinct(p)
}
