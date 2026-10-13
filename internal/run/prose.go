package run

// This file holds the prose session of the step setup (docs/spec/scaffold.md,
// The prose rows; row 54 of the plan, task T-eahr, #277): one session of the role
// product-owner for each stop of the runner with F- rows, whose output is one
// file for each row, at result/files/<path>.

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/setup"
)

// proseOutput is the directory of result/ that holds the files of the session.
const proseOutput = "files"

// proseNames gives the paths that the session writes: the path of each F- row of
// the stop, in the order of the table. A row that is no F- row of a path of the
// target below its root (the runner gives the path as the question's name and as
// its where), or a path twice, is a defect of the runner, and the step is fail.
func proseNames(rows []setup.StopRow) ([]string, error) {
	if len(rows) == 0 {
		return nil, errors.New("the stop has no row for a file")
	}
	var names []string
	for _, r := range rows {
		path, ok := strings.CutPrefix(r.Question, "F-")
		if !ok || r.Where != path || !fs.ValidPath(path) || path == "." {
			return nil, fmt.Errorf("the row %s of %s is not the F- row of a path of the target", r.Question, r.Step)
		}
		if slices.Contains(names, path) {
			return nil, fmt.Errorf("the path %s is in the stop twice", path)
		}
		names = append(names, path)
	}
	return names, nil
}

// proseKeep gives the Keep of the session for the paths that the stop named: the
// records of the files at setup/inputs/files/<path>. The reader of the output has
// checked the set; a file that no name covers or a missing one is refused here
// too, so a records commit never holds a file that the stop did not name.
func proseKeep(names []string) func(id string, files map[string][]byte) (map[string][]byte, error) {
	return func(id string, files map[string][]byte) (map[string][]byte, error) {
		if len(files) != len(names) {
			return nil, fmt.Errorf("%d files; the stop named %d", len(files), len(names))
		}
		kept := map[string][]byte{}
		for _, n := range names {
			data, ok := files[n]
			if !ok {
				return nil, fmt.Errorf("the file %s is missing", n)
			}
			kept["setup/inputs/files/"+n] = data
		}
		return kept, nil
	}
}

// stopTable is the stop table of the runner, the block setup-stop, as the
// session's input/stop.tsv.
func stopTable(rows []setup.StopRow) ([]byte, error) {
	var b bytes.Buffer
	if err := (setup.Result{Stops: rows}).Write(&b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// proseText gives the prompt of the session: the paths of its input files, each
// on a line of its own (the vision brief only when the target has one), and of
// the directory of its output.
func proseText(input map[string]string, output string) []byte {
	var b strings.Builder
	b.WriteString("You write the prose files of a project's repository, for LAYUP. Change no file in this repository and make no commit.\n\n")
	b.WriteString("The problem statement: " + input["problem-statement.md"] + "\n")
	if v := input["vision.md"]; v != "" {
		b.WriteString("The vision brief: " + v + "\n")
	}
	b.WriteString("The stop table: " + input["stop.tsv"] + "\n")
	b.WriteString("The output directory: " + output + "\n\n")
	b.WriteString(proseRule)
	return []byte(b.String())
}

// proseRule is the fixed text of the prompt after its paths. The ask of an
// adapted file begins Adapt (internal/setup), and S14 refuses an input that loses
// a place of a marker of the file (docs/spec/setup.md, S14).
const proseRule = `The stop table has one row for each file to write: its where column is the
path of the file in the project, and its ask column says what the file must be.
Read the problem statement, the vision brief when there is one, and the project
in this directory, which is as it is now. Write each file of the table in the
output directory, at its path below it (the file docs/glossary.md goes to
<output directory>/docs/glossary.md); make the directories that you need. Write
the files of the table and no other: a file that is missing, a file that the table
does not name, a link, a file over 1 MiB and a file that is not valid UTF-8 refuse
the whole output. Write plain text or Markdown in the language of the problem
statement, with no placeholder and no angle quote of your own. A row whose ask
says Adapt names a file that this directory already has: give that file with
your changes, and keep each marker that it holds (a text between the angle
quotes U+2039 and U+203A), in every place, byte for byte, with its line ends.
Then end.
`
