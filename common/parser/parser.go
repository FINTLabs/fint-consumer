package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/FINTLabs/fint-consumer/common/config"
	"github.com/FINTLabs/fint-consumer/common/document"
	"github.com/FINTLabs/fint-consumer/common/types"
	"github.com/FINTLabs/fint-consumer/common/utils"
	xmlquery "github.com/antchfx/xquery/xml"
)

func GetClasses(owner string, repo string, tag string, filename string, force bool) ([]*types.Class, map[string]types.Import, map[string][]*types.Class, map[string][]*types.Class) {
	fmt.Printf("Fetching document %s/%s/%s @ %s ...", owner, repo, filename, tag)
	doc, err := document.Get(owner, repo, tag, filename, force)
	if err != nil {
		return nil, nil, nil, nil
	}
	fmt.Println("ok")

	fmt.Print("Parsing ...")

	var classes []*types.Class
	// TODO BUG: packageMap and classMap fail for classes with the same name!
	packageMap := make(map[string]types.Import)
	classMap := make(map[string]*types.Class)
	javaPackageClassMap := make(map[string][]*types.Class)
	csPackageClassMap := make(map[string][]*types.Class)

	classElements := xmlquery.Find(doc, "//element[@type='Class']")

	fmt.Print(".")

	for i, c := range classElements {

		if i%10 == 0 {
			fmt.Print(".")
		}

		properties := c.SelectElement("properties")
		class := new(types.Class)

		class.Tag = tag
		class.Name = replaceNO(c.SelectAttr("name"))
		class.Abstract = toBool(properties.SelectAttr("isAbstract"))
		class.Extends = getExtends(doc, c)
		class.Attributes = getAttributes(c)
		class.Relations = getRelations(doc, c)
		class.Package = getPackagePath(c, doc)
		class.Namespace = getNamespacePath(c, doc)
		class.Identifiable = identifiable(class.Attributes)
		class.Writable = writable(class.Attributes)
		class.Stereotype = properties.SelectAttr("stereotype")
		class.Documentation = properties.SelectAttr("documentation")
		class.Deprecated = c.SelectElement("tags/tag[@name='DEPRECATED']") != nil

		if len(class.Stereotype) == 0 {
			if class.Abstract {
				class.Stereotype = "abstrakt"
			}
		}

		imp := types.Import{
			Java:   fmt.Sprintf("%s.%s", class.Package, class.Name),
			CSharp: class.Namespace,
		}
		packageMap[class.Name] = imp

		classes = append(classes, class)
		classMap[class.Name] = class
	}

	fmt.Print(".")

	packageMap["Date"] = types.Import{
		Java: "java.util.Date",
	}

	for _, class := range classes {
		class.Imports = getImports(class, packageMap)
		class.Using = getUsing(class, packageMap)
		class.Identifiable = identifiableFromExtends(class, classMap)
		class.Resource = isResource(class, classMap)
		class.Identifiers = getIdentifiers(class, classMap)
		javaPackageClassMap[class.Package] = append(javaPackageClassMap[class.Package], class)
		csPackageClassMap[class.Namespace] = append(csPackageClassMap[class.Namespace], class)
		if len(class.Stereotype) == 0 {
			if class.Identifiable {
				class.Stereotype = "hovedklasse"
			} else {
				class.Stereotype = "datatype"
			}
		}
		if !class.Writable {
			class.Writable = getWritableFromExtends(class, classMap)
		}
	}

	fmt.Print(".")

	for _, class := range classes {
		for i, a := range class.Attributes {
			if typ, found := classMap[a.Type]; found {
				class.Attributes[i].Package = typ.Package
				if typ.Resource {
					class.Resources = append(class.Resources, a)
				}
			}
		}
	}

	fmt.Print(".")

	for _, class := range classes {
		if len(class.Extends) > 0 {
			if typ, found := classMap[class.Extends]; found {
				class.ExtendsResource = typ.Resource || len(typ.Resources) > 0
			}
		}
		class.InheritedAttributes = getAttributesFromExtends(class, classMap)
	}

	fmt.Print(".")

	for _, class := range classes {
		for i, r := range class.Relations {
			class.Relations[i].Stereotype = classMap[r.Target].Stereotype
		}
	}

	fmt.Println(". done")
	return classes, packageMap, javaPackageClassMap, csPackageClassMap
}

func getWritableFromExtends(class *types.Class, classMap map[string]*types.Class) bool {
	if len(class.Extends) > 0 {
		extendedClass := classMap[class.Extends]
		if extendedClass.Writable {
			return true
		}
		for len(extendedClass.Extends) > 0 {
			extendedClass = classMap[extendedClass.Extends]
			if extendedClass.Writable {
				return true
			}
		}
	}
	return false
}

func identifiableFromExtends(class *types.Class, classMap map[string]*types.Class) bool {
	if class.Identifiable {
		return true
	}
	if len(class.Extends) > 0 {
		return identifiableFromExtends(classMap[class.Extends], classMap)
	}
	return false
}

func getClassByIdRef(idref string, doc *xmlquery.Node) *xmlquery.Node {
	result := xmlquery.Find(doc, fmt.Sprintf("//element[@idref='%s']", idref))
	return result[0]
}

func getAttributesFromExtends(class *types.Class, classMap map[string]*types.Class) []types.InheritedAttribute {

	var result []types.InheritedAttribute

	extendedClass := class

	for len(extendedClass.Extends) > 0 {
		extendedClass = classMap[extendedClass.Extends]
		for _, a := range extendedClass.Attributes {
			var att = types.InheritedAttribute{
				Owner:     extendedClass.Name,
				Attribute: a,
			}
			result = append(result, att)
		}
	}

	return result
}

func isResource(class *types.Class, classMap map[string]*types.Class) bool {
	if len(class.Relations) > 0 {
		return true
	}
	if len(class.Extends) > 0 {
		return isResource(classMap[class.Extends], classMap)
	}
	return false
}

func writable(attribs []types.Attribute) bool {

	for _, value := range attribs {
		if value.Writable {
			return true
		}
	}
	return false
}

func identifiable(attribs []types.Attribute) bool {

	for _, value := range attribs {
		if value.Type == "Identifikator" {
			return true
		}
	}

	return false
}

func getIdentifiers(class *types.Class, classMap map[string]*types.Class) []types.Identifier {
	var identifiers []types.Identifier

	for _, a := range class.Attributes {
		if a.Type == "Identifikator" {
			identifier := types.Identifier{}
			identifier.Name = a.Name
			identifier.Optional = a.Optional
			identifiers = append(identifiers, identifier)
		}
	}

	if len(class.Extends) > 0 {
		identifiers = append(identifiers, getIdentifiers(classMap[class.Extends], classMap)...)
	}

	return identifiers
}

func getImports(c *types.Class, imports map[string]types.Import) []string {

	attribs := c.Attributes
	var imps []string
	for _, att := range attribs {
		javaType := types.GetJavaType(att.Type)
		if imports[javaType].Java != c.Package && len(javaType) > 0 {
			imps = append(imps, imports[javaType].Java)
		}
	}

	if len(c.Extends) > 0 {
		imps = append(imps, imports[c.Extends].Java)
	}

	if c.Identifiable {
		imps = append(imps, imports["Identifikator"].Java)
	}

	return utils.Distinct(utils.TrimArray(imps))
}

func getUsing(c *types.Class, imports map[string]types.Import) []string {

	attribs := c.Attributes
	var imps []string
	for _, att := range attribs {
		csType := types.GetCSType(att.Type)
		if imports[csType].CSharp != c.Package && len(imports[csType].CSharp) > 0 {
			imps = append(imps, imports[csType].CSharp)
		}
	}

	if len(c.Extends) > 0 {
		imps = append(imps, imports[c.Extends].CSharp)
	}

	return utils.Distinct(utils.TrimArray(imps))

}

func getPackagePath(c *xmlquery.Node, doc *xmlquery.Node) string {

	var pkgs []string
	var parentPkg string
	classPkg := getPackage(c)
	pkgs = append(pkgs, getNameLower(classPkg, doc))

	parentPkg = getParentPackage(classPkg, doc)
	for parentPkg != "" {
		pkgs = append(pkgs, getNameLower(parentPkg, doc))
		parentPkg = getParentPackage(parentPkg, doc)
	}
	pkgs = utils.TrimArray(pkgs)
	pkgs = utils.Reverse(pkgs)
	return replaceNO(fmt.Sprintf("%s.%s", config.JAVA_PACKAGE_BASE, strings.Join(pkgs, ".")))

}

func getNamespacePath(c *xmlquery.Node, doc *xmlquery.Node) string {

	var pkgs []string
	var parentPkg string
	classPkg := getPackage(c)
	pkgs = append(pkgs, getName(classPkg, doc))

	parentPkg = getParentPackage(classPkg, doc)
	for parentPkg != "" {
		pkgs = append(pkgs, getName(parentPkg, doc))
		parentPkg = getParentPackage(parentPkg, doc)
	}
	pkgs = utils.TrimArray(pkgs)
	pkgs = utils.Reverse(pkgs)
	return replaceNO(fmt.Sprintf("%s.%s", config.NET_NAMESPACE_BASE, strings.Join(pkgs, ".")))

}

func getName(idref string, doc *xmlquery.Node) string {
	name := ""
	if len(idref) > 0 {
		xpath := fmt.Sprintf("//element[@idref='%s']", idref)
		parent := xmlquery.Find(doc, xpath)

		name = parent[0].SelectAttr("name")
		name = excludeName(name)
	}
	return strings.Replace(name, " ", "", -1)
}

func excludeName(name string) string {
	if name == "FINT" {
		name = strings.Replace(name, "FINT", "", -1)
	}
	if name == "Model" {
		name = strings.Replace(name, "Model", "", -1)
	}
	return name
}

func getNameLower(idref string, doc *xmlquery.Node) string {

	return strings.ToLower(getName(idref, doc))
}

func getParentPackage(idref string, doc *xmlquery.Node) string {
	xpath := fmt.Sprintf("//element[@idref='%s']", idref)

	parent := xmlquery.Find(doc, xpath)

	if len(parent) > 1 {
		fmt.Printf("More than one element with idref %s\n", idref)
		return ""
	}
	if len(parent) < 1 {
		return ""
	}

	model := parent[0].SelectElement("model")
	if model == nil {
		return ""
	}

	return model.SelectAttr("package")
}

func getPackage(c *xmlquery.Node) string {
	return c.SelectElement("model").SelectAttr("package")
}

func getExtends(doc *xmlquery.Node, c *xmlquery.Node) string {

	var extends []string
	for _, rr := range xmlquery.Find(doc, fmt.Sprintf("//connectors/connector/properties[@ea_type='Generalization']/../source[@idref='%s']/../target/model[@name]", c.SelectAttr("idref"))) {
		if len(rr.SelectAttr("name")) > 0 {
			extends = append(extends, replaceNO(rr.SelectAttr("name")))
		}
	}

	if len(extends) == 1 {
		return extends[0]
	}

	return ""
}

func getAttributes(c *xmlquery.Node) []types.Attribute {
	var attributes []types.Attribute
	for _, a := range xmlquery.Find(c, "//attributes/attribute") {

		attrib := types.Attribute{}
		attrib.Name = replaceNO(a.SelectAttr("name"))
		attrib.Deprecated = a.SelectElement("tags/tag[@name='DEPRECATED']") != nil
		attrib.List = strings.Compare(a.SelectElement("bounds").SelectAttr("upper"), "*") == 0
		attrib.Optional = !attrib.List && strings.Compare(a.SelectElement("bounds").SelectAttr("lower"), "0") == 0
		attrib.Type = replaceNO(a.SelectElement("properties").SelectAttr("type"))
		attrib.Writable = a.SelectElement("stereotype").SelectAttr("stereotype") == "writable"

		attributes = append(attributes, attrib)
	}

	return attributes
}

func getRelations(doc *xmlquery.Node, c *xmlquery.Node) []types.Association {
	var assocs []types.Association
	isAbstract := toBool(c.SelectElement("properties").SelectAttr("isAbstract"))
	if !isAbstract {
		assocs = getAssociations(doc, c)
		assocs = append(assocs, getRecursivelyAssociationsFromExtends(doc, c)...)
	}

	return assocs
}

// TODO: This is actually iterative and not recursive, and works only for linear inheritance.
// TODO: Possible to combine with getAssociationsFromExtends?
func getRecursivelyAssociationsFromExtends(doc *xmlquery.Node, c *xmlquery.Node) []types.Association {
	var assocs []types.Association
	extAssocs, extends := getAssociationsFromExtends(doc, c)
	assocs = append(assocs, extAssocs...)
	for {
		if extends == nil {
			break
		}
		extAssocs, extends = getAssociationsFromExtends(doc, extends)
		assocs = append(assocs, extAssocs...)
	}
	return assocs
}

func getAssociationsFromExtends(doc *xmlquery.Node, c *xmlquery.Node) ([]types.Association, *xmlquery.Node) {
	var assocs []types.Association
	extends := xmlquery.Find(doc, fmt.Sprintf("//connectors/connector/properties[@ea_type='Generalization']/../source[@idref='%s']/../target", c.SelectAttr("idref")))
	if len(extends) == 1 {
		assocs = append(assocs, getAssociations(doc, extends[0])...)
		return assocs, extends[0]
	}

	return assocs, nil
}

func getMultiplicity(multiplicity string) (bool, bool) {
	return strings.HasPrefix(multiplicity, "0"),
		strings.HasSuffix(multiplicity, "*")
}

func getAssociations(doc *xmlquery.Node, c *xmlquery.Node) []types.Association {
	var assocs []types.Association
	for _, rr := range xmlquery.Find(doc, fmt.Sprintf("//connectors/connector/properties[@ea_type='Association']/../source[@idref='%s']/../target/role", c.SelectAttr("idref"))) {
		if len(rr.SelectAttr("name")) > 0 {
			assoc := types.Association{}
			assoc.Name = replaceNO(rr.SelectAttr("name"))
			assoc.Target = replaceNO(rr.SelectElement("../model").SelectAttr("name"))
			assoc.Optional, assoc.List = getMultiplicity(rr.SelectElement("../type").SelectAttr("multiplicity"))
			assoc.Deprecated = rr.SelectElement("../../tags/tag[@name='DEPRECATED']") != nil
			assoc.TargetPackage = getPackagePath(getClassByIdRef(rr.SelectElement("../../target").SelectAttr("idref"), doc), doc)
			assocs = append(assocs, assoc)
		}
	}
	for _, rl := range xmlquery.Find(doc, fmt.Sprintf("//connectors/connector/properties[@ea_type='Association']/../target[@idref='%s']/../source/role", c.SelectAttr("idref"))) {
		if len(rl.SelectAttr("name")) > 0 {
			assoc := types.Association{}
			assoc.Name = replaceNO(rl.SelectAttr("name"))
			assoc.Target = replaceNO(rl.SelectElement("../model").SelectAttr("name"))
			assoc.Optional, assoc.List = getMultiplicity(rl.SelectElement("../type").SelectAttr("multiplicity"))
			assoc.Deprecated = rl.SelectElement("../../tags/tag[@name='DEPRECATED']") != nil
			assoc.TargetPackage = getPackagePath(getClassByIdRef(rl.SelectElement("../../source").SelectAttr("idref"), doc), doc)
			assocs = append(assocs, assoc)
		}
	}
	return assocs
}

func replaceNO(s string) string {

	r := strings.Replace(s, "æ", "a", -1)
	r = strings.Replace(r, "ø", "o", -1)
	r = strings.Replace(r, "å", "a", -1)
	r = strings.Replace(r, "Æ", "A", -1)
	r = strings.Replace(r, "Ø", "O", -1)
	r = strings.Replace(r, "Å", "A", -1)
	return r
}

func toBool(s string) bool {
	b, err := strconv.ParseBool(s)

	if err != nil {
		return false
	}

	return b
}
