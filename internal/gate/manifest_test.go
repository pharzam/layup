package gate

import (
	"strings"
	"testing"
)

const manifestHeader = "kind\tstate\ttool\tcommand\tscope\tconfig\n"

func TestReadManifestGivesTheKindsInTheirOrder(t *testing.T) {
	kinds, err := ReadManifest([]byte(manifestHeader +
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
		"no scope pattern":       manifestHeader + strings.Replace(row, "./*.go", "—", 1),
		"a config path .git":     manifestHeader + strings.Replace(row, "\t—\n", "\t.git\n", 1),
		"a config path in .git":  manifestHeader + strings.Replace(row, "\t—\n", "\tlayout .git/hooks/x\n", 1),
		"a config path .GIT":     manifestHeader + strings.Replace(row, "\t—\n", "\t.GIT\n", 1),
		"a config path in .Git":  manifestHeader + strings.Replace(row, "\t—\n", "\t.Git/x\n", 1),
	} {
		if _, err := ReadManifest([]byte(text)); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
}

// A config path that only starts with .git is a path like any other.
func TestReadManifestTakesAConfigPathThatOnlyStartsWithGit(t *testing.T) {
	text := manifestHeader + "static\tactive\tgo\tgo vet ./...\t./*.go\t.gitignore .github/x.yml\n"
	if kinds, err := ReadManifest([]byte(text)); err != nil || strings.Join(kinds[0].Config, " ") != ".gitignore .github/x.yml" {
		t.Fatalf("%+v, %v", kinds, err)
	}
}
