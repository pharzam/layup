package gate

import (
	"strings"
	"testing"
)

const manifestHeader = "kind\tstate\ttool\tcommand\tscope\tconfig\n"

func TestReadManifestGivesTheKindsInTheirOrder(t *testing.T) {
	kinds, err := readManifest([]byte(manifestHeader +
		"static\tactive\tgo\ttest -z \"$(gofmt -l .)\"\t./*.go\t—\n" +
		"layout\tpending\tgo\tgo test ./layout/\t./*.go internal\tlayout layout/rules.txt\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(kinds) != 2 || kinds[0].Name != "static" || kinds[1].Name != "layout" || kinds[1].State != "pending" ||
		len(kinds[1].Scope) != 2 || strings.Join(kinds[1].Config, " ") != "layout layout/rules.txt" || len(kinds[0].Config) != 0 {
		t.Fatalf("the kinds %+v", kinds)
	}
}

func TestReadManifestRefusesEachMalformedForm(t *testing.T) {
	row := "static\tactive\tgo\tgo vet ./...\t./*.go\t—\n"
	for name, text := range map[string]string{
		"a wrong header":         "kind\tstate\ttool\tcommand\tscope\n" + row,
		"a wrong field count":    manifestHeader + "static\tactive\tgo\tgo vet ./...\t./*.go\n",
		"a bad state":            manifestHeader + strings.Replace(row, "active", "on", 1),
		"a kind with a digit":    manifestHeader + strings.Replace(row, "static", "static2", 1),
		"a repeated kind":        manifestHeader + row + row,
		"a .. in a config path":  manifestHeader + strings.Replace(row, "\t—\n", "\t../x\n", 1),
		"an empty field":         manifestHeader + strings.Replace(row, "\tgo\t", "\t\t", 1),
		"no row":                 manifestHeader,
		"a pattern of no form":   manifestHeader + strings.Replace(row, "./*.go", "*.go", 1),
		"a pattern that is ./ a": manifestHeader + strings.Replace(row, "./*.go", "./internal", 1),
	} {
		if _, err := readManifest([]byte(text)); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
}
