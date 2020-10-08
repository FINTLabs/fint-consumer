package types

type Class struct {
	Tag                 string
	Name                string
	Abstract            bool
	Deprecated          bool
	Extends             string
	Package             string
	Imports             []string
	Namespace           string
	Using               []string
	Documentation       string
	Attributes          []Attribute
	InheritedAttributes []InheritedAttribute
	Relations           []Association
	Resources           []Attribute
	Resource            bool
	ExtendsResource     bool
	Identifiable        bool
	Writable            bool
	Stereotype          string
	Identifiers         []Identifier
}
