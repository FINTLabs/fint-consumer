package types

type Class struct {
	Name             string
	Abstract         bool
	Extends          string
	Package          string
	Imports          []string
	Namespace        string
	Using            []string
	DocumentationUrl string
	Attributes       []Attribute
	Relations        []string
	Identifiable     bool
}
