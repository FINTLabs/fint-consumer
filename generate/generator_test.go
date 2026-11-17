package generate

import (
	"testing"

	"github.com/FINTLabs/fint-consumer/common/config"
)

// The tests use config.JAVA_PACKAGE_BASE, so they work for the v3 base
// (no.fint.model) and the v4 base (no.novari.fint.model).

func TestGetMainPackage(t *testing.T) {
	base := config.JAVA_PACKAGE_BASE
	cases := map[string]string{
		base + ".arkiv.noark":             "arkiv/noark",
		base + ".felles":                  "felles",
		base + ".felles.kodeverk.iso":     "felles/kodeverk/iso",
		base + ".utdanning.elev":          "utdanning/elev",
		base + ".administrasjon.personal": "administrasjon/personal",
	}
	for in, want := range cases {
		if got := GetMainPackage(in); got != want {
			t.Errorf("GetMainPackage(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestModelPkg(t *testing.T) {
	modelPkg := funcMap["modelPkg"].(func(string) string)
	base := config.JAVA_PACKAGE_BASE
	cases := map[string]string{
		base + ".arkiv.noark":         "",
		base + ".felles.kodeverk":     "",
		base + ".felles.kodeverk.iso": "iso.",
		base + ".felles":              "",
	}
	for in, want := range cases {
		if got := modelPkg(in); got != want {
			t.Errorf("modelPkg(%q) = %q, want %q", in, got, want)
		}
	}
}
