package run

import (
	"reflect"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/setup"
)

// The prose session of the step setup (docs/spec/scaffold.md, The prose rows;
// task T-eahr, #277): its stop table, the paths that it writes, its prompt, and
// the records that its files become.

func proseRows() []setup.StopRow {
	return []setup.StopRow{
		{Step: "S07", Question: "F-docs/onboarding-for-engineers.md", Ask: "Write the text of docs/onboarding-for-engineers.md for the target, and give it as inputs/files/docs/onboarding-for-engineers.md.", Where: "docs/onboarding-for-engineers.md"},
		{Step: "S07", Question: "F-docs/glossary.md", Ask: "Write the text of docs/glossary.md for the target, and give it as inputs/files/docs/glossary.md.", Where: "docs/glossary.md"},
	}
}

// The files that the session writes are the paths of the F- rows, in the order
// of the stop table; a row that is no F- row of its own path is a defect of the
// runner, and the step is fail.
func TestTheProseNamesAreThePathsOfTheFRows(t *testing.T) {
	names, err := proseNames(proseRows())
	if want := []string{"docs/onboarding-for-engineers.md", "docs/glossary.md"}; err != nil || !reflect.DeepEqual(names, want) {
		t.Errorf("the names %q, %v; want %q", names, err, want)
	}
	for name, rows := range map[string][]setup.StopRow{
		"no row":                       nil,
		"an M- row":                    {{Step: "S10", Question: "M-aaaaaaaa", Ask: "q", Where: "NOTES.md:2 <x>"}},
		"an M- row that names itself":  {{Step: "S10", Question: "M-aaaaaaaa", Ask: "q", Where: "M-aaaaaaaa"}},
		"the root as a path":           {{Step: "S07", Question: "F-.", Ask: "q", Where: "."}},
		"an F- row of another path":    {{Step: "S07", Question: "F-docs/a.md", Ask: "q", Where: "docs/b.md"}},
		"a path that leaves the tree":  {{Step: "S07", Question: "F-../a.md", Ask: "q", Where: "../a.md"}},
		"an absolute path":             {{Step: "S07", Question: "F-/a.md", Ask: "q", Where: "/a.md"}},
		"the same path twice":          {proseRows()[0], proseRows()[0]},
		"an F- row with no path":       {{Step: "S07", Question: "F-", Ask: "q", Where: ""}},
		"an F- row of an unclean path": {{Step: "S07", Question: "F-docs//a.md", Ask: "q", Where: "docs//a.md"}},
	} {
		if names, err := proseNames(rows); err == nil {
			t.Errorf("%s: %q; want an error", name, names)
		}
	}
}

// The files go to the records at setup/inputs/files/<path>, and only the files
// that the stop named: a file that no name covers, or that is missing, is a
// refusal of the output (the session reader has checked it; this is the second line).
func TestTheProseKeepGivesTheFilesAtTheirRecordsPaths(t *testing.T) {
	keep := proseKeep([]string{"docs/glossary.md", "README.md"})
	got, err := keep("S-1a2b3c4d", map[string][]byte{"docs/glossary.md": []byte("# G\n"), "README.md": []byte("# R\n")})
	want := map[string][]byte{"setup/inputs/files/docs/glossary.md": []byte("# G\n"), "setup/inputs/files/README.md": []byte("# R\n")}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("the records %q, %v; want %q", got, err, want)
	}
	for name, files := range map[string]map[string][]byte{
		"a file that no name covers": {"docs/glossary.md": []byte("g"), "README.md": []byte("r"), "extra.md": []byte("x")},
		"a missing file":             {"README.md": []byte("r")},
		"a file in place of a named": {"docs/glossary.md": []byte("g"), "other.md": []byte("x")},
		"no file":                    {},
	} {
		if got, err := keep("S-1a2b3c4d", files); err == nil || got != nil {
			t.Errorf("%s: %q, %v; want an error and no record", name, got, err)
		}
	}
}

// The prompt names, one on each line, the problem statement, the vision brief
// when there is one, the stop table and the directory of the output files.
func TestTheProsePromptNamesItsPaths(t *testing.T) {
	input := map[string]string{"problem-statement.md": "/s/input/problem-statement.md", "vision.md": "/s/input/vision.md", "stop.tsv": "/s/input/stop.tsv"}
	text := string(proseText(input, "/s/result/files"))
	for _, line := range []string{"The problem statement: /s/input/problem-statement.md", "The vision brief: /s/input/vision.md", "The stop table: /s/input/stop.tsv", "The output directory: /s/result/files"} {
		if !strings.Contains("\n"+text, "\n"+line+"\n") {
			t.Errorf("the prompt has no line %q:\n%s", line, text)
		}
	}
	delete(input, "vision.md")
	if text := string(proseText(input, "/s/result/files")); strings.Contains(text, "The vision brief:") || strings.Contains(text, "vision.md") {
		t.Errorf("the prompt of a target with no vision brief names its file:\n%s", text)
	}
}

// A file that the runner asks to adapt keeps each marker that it holds: S14 refuses an input that
// loses a place of one (docs/spec/setup.md), so the ban of angle quotes has that exception.
func TestTheProsePromptLetsAnAdaptedFileKeepItsMarkers(t *testing.T) {
	text := strings.Join(strings.Fields(string(proseText(map[string]string{"stop.tsv": "/s/input/stop.tsv"}, "/s/result/files"))), " ")
	for _, want := range []string{"A row whose ask says Adapt", "keep each marker that it holds", "U+2039 and U+203A", "in every place, byte for byte, with its line ends"} {
		if !strings.Contains(text, want) {
			t.Errorf("the prompt does not say %q:\n%s", want, text)
		}
	}
}

// The stop table of the session is the table that the runner printed, the block
// setup-stop, byte for byte.
func TestTheStopTableOfTheSessionIsTheRunnersTable(t *testing.T) {
	got, err := stopTable(proseRows())
	var want strings.Builder
	if werr := (setup.Result{Stops: proseRows()}).Write(&want); err != nil || werr != nil || string(got) != want.String() || !strings.HasPrefix(string(got), "step\tquestion\task\twhere\n") {
		t.Errorf("the stop table %q, %v; want %q, %v", got, err, want.String(), werr)
	}
}
