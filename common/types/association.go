package types

type Association struct {
	Name          string
	Target        string
	TargetPackage string
	Deprecated    bool
	Optional      bool
	List          bool
	Stereotype    string
}
