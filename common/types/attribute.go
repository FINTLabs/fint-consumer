package types

type Attribute struct {
	Name       string
	Type       string
	Package    string
	List       bool
	Optional   bool
	Writable   bool
	Deprecated bool
}

type InheritedAttribute struct {
	Owner string
	Attribute
}
