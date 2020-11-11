package types

import "strings"

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

type ByName []*Class

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return strings.ToLower(n[i].Name) < strings.ToLower(n[j].Name) }
