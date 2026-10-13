package run

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

// The step scaffold-push (docs/spec/scaffold.md, The commands for the Operator, with its
// paragraph "Decided here" of task T-hfn1; the four input states of the step; the acceptance
// row "The apply outputs"; row 56 of the plan, #279): the record of the apply outputs, and
// the step with its seams as fakes. A failure and a wait never reach the records commit.

const (
	defaultOut = "{\"source\":\"acme/target\",\"bypass_actors\":[]}\n"
	recordsOut = "{\"source\":\"acme/target\",\"bypass_actors\":[{\"actor_id\":42}]}"
)

// The record is, for default.json and then records.json, the line == <file>, the output byte
// for byte, and one line feed that is the separator, whether or not the output ends with one.
func TestTheRecordOfTheAppliesFramesEachOutputByteForByte(t *testing.T) {
	for _, c := range []struct{ name, d, r, want string }{
		{"each ends with a line feed", defaultOut, recordsOut + "\n", "== default.json\n" + defaultOut + "\n== records.json\n" + recordsOut + "\n\n"},
		{"neither ends with a line feed", strings.TrimSuffix(defaultOut, "\n"), recordsOut, "== default.json\n" + strings.TrimSuffix(defaultOut, "\n") + "\n== records.json\n" + recordsOut + "\n"},
		{"the first ends with one, the second does not", defaultOut, recordsOut, "== default.json\n" + defaultOut + "\n== records.json\n" + recordsOut + "\n"},
		{"the second ends with one, the first does not", strings.TrimSuffix(defaultOut, "\n"), recordsOut + "\n", "== default.json\n" + strings.TrimSuffix(defaultOut, "\n") + "\n== records.json\n" + recordsOut + "\n\n"},
	} {
		got, err := rulesetsRecord("acme/target", map[string][]byte{"default.json": []byte(c.d), "records.json": []byte(c.r)})
		if err != nil || string(got) != c.want {
			t.Errorf("%s: %q, %v; want %q", c.name, got, err, c.want)
		}
	}
}

// An output that is not a JSON object whose source is OWNER/NAME is refused, naming its file.
func TestAnApplyOutputThatIsNotARulesetOfTheTargetIsRefusedNamingItsFile(t *testing.T) {
	const notObject, otherSource = "is not a JSON object", "its source is"
	for _, c := range []struct{ name, file, text, want string }{
		{"an array", "default.json", `[{"source":"acme/target"}]`, notObject},
		{"a string", "records.json", `"acme/target"`, notObject},
		{"invalid JSON", "default.json", `{"source":`, notObject},
		{"an empty file", "records.json", ``, notObject},
		{"a number", "default.json", `7`, notObject},
		{"null", "records.json", `null`, notObject},
		{"another repository", "records.json", `{"source":"acme/other"}`, otherSource},
		{"the owner in another case", "default.json", `{"source":"Acme/target"}`, otherSource},
		{"no source", "default.json", `{"id":1}`, otherSource},
		{"a source that is not a string", "records.json", `{"source":7}`, otherSource},
	} {
		files := both()
		files[c.file] = []byte(c.text)
		got, err := rulesetsRecord("acme/target", files)
		if err == nil || !strings.Contains(err.Error(), c.file) || !strings.Contains(err.Error(), c.want) || got != nil {
			t.Errorf("%s in %s: %q, %v; want an error that names %s and says %q, and no record", c.name, c.file, got, err, c.file, c.want)
		}
	}
}

// pushFake is the world of a job: what each seam gives, and what the job asked of them.
type pushFake struct {
	last, head                                                    string
	lastErr, areaErr, headErr, appliesErr, commandsErr, commitErr error
	applies                                                       map[string][]byte
	called, notes                                                 []string
	areaFor                                                       string
	commits                                                       []map[string][]byte
}

const pushCommands = "set -eu\n# the push\ngit push origin layup-setup:main\n"

func (f *pushFake) job() scaffoldJob {
	return scaffoldJob{
		source: "acme/target",
		last:   func() (string, error) { f.called = append(f.called, "last"); return f.last, f.lastErr },
		area:   func(last string) error { f.called = append(f.called, "area"); f.areaFor = last; return f.areaErr },
		head:   func() (string, error) { f.called = append(f.called, "head"); return f.head, f.headErr },
		applies: func() (map[string][]byte, error) {
			f.called = append(f.called, "applies")
			return f.applies, f.appliesErr
		},
		commands: func() (string, string, error) {
			f.called = append(f.called, "commands")
			return "DIR/work/acme/target/out/commands.sh", pushCommands, f.commandsErr
		},
		print: func(note string) { f.notes = append(f.notes, note) },
		commit: func(files map[string][]byte) error {
			f.called = append(f.called, "commit")
			f.commits = append(f.commits, files)
			return f.commitErr
		},
	}
}

func both() map[string][]byte {
	return map[string][]byte{"default.json": []byte(defaultOut), "records.json": []byte(recordsOut)}
}

const sha1 = "0123456789012345678901234567890123456789"
const sha2 = "abcdefabcdefabcdefabcdefabcdefabcdefabcd"

func TestTheStepScaffoldPushWaitsAndPrintsTheWholeRuleOfTheOperator(t *testing.T) {
	for _, c := range []struct {
		name string
		f    pushFake
	}{
		{"the head of the forge is not the last commit", pushFake{last: sha1, head: sha2, applies: both()}},
		{"the forge has no head for the default branch", pushFake{last: sha1, head: "", applies: both()}},
		{"the head is right and default.json is missing", pushFake{last: sha1, head: sha1, applies: map[string][]byte{"records.json": []byte(recordsOut)}}},
		{"the head is right and records.json is missing", pushFake{last: sha1, head: sha1, applies: map[string][]byte{"default.json": []byte(defaultOut)}}},
		{"nothing yet", pushFake{last: sha1, head: sha2}},
	} {
		f := c.f
		got := f.job().do()
		if got != (Step{"scaffold-push", "done", "waiting for the Operator's pushes and ruleset apply"}) {
			t.Errorf("%s: %+v; want done, waiting", c.name, got)
		}
		if want := []string{"last", "area", "head", "applies", "commands"}; !slices.Equal(f.called, want) || len(f.commits) != 0 || f.areaFor != sha1 {
			t.Errorf("%s: the calls %q, the commits %q, the area for %q; want %q, no commit, and the area for the last commit", c.name, f.called, f.commits, f.areaFor, want)
		}
		if len(f.notes) != 4 || !strings.Contains(f.notes[0], "DIR/work/acme/target/out/commands.sh") || !strings.Contains(f.notes[0], "your own login") ||
			!strings.Contains(f.notes[0], "no hooks installed") || f.notes[1] != "set -eu" || f.notes[2] != "# the push" || f.notes[3] != "git push origin layup-setup:main" {
			t.Errorf("%s: the notes %q; want the instruction (the file, your own login, no hooks installed) and each line of the file", c.name, f.notes)
		}
	}
}

func TestTheStepScaffoldPushGoesOnWithOneCommitOfTheRecordOfTheApplies(t *testing.T) {
	f := &pushFake{last: sha1, head: sha1, applies: both()}
	got := f.job().do()
	want, _ := rulesetsRecord("acme/target", both())
	if got != (Step{"scaffold-push", "done", "the apply of the rulesets recorded"}) || len(f.notes) != 0 {
		t.Errorf("the step %+v and the notes %q; want done and nothing printed", got, f.notes)
	}
	if calls := []string{"last", "area", "head", "applies", "commit"}; !slices.Equal(f.called, calls) {
		t.Errorf("the calls %q; want %q (the work area is made sure of before the outcome)", f.called, calls)
	}
	if len(f.commits) != 1 || len(f.commits[0]) != 1 || string(f.commits[0]["setup/rulesets.txt"]) != string(want) {
		t.Errorf("the commits %q; want one of setup/rulesets.txt alone, %q", f.commits, want)
	}
}

func TestEachFailureOfTheStepScaffoldPushPublishesNoCommit(t *testing.T) {
	boom := errors.New("boom\nthe second line")
	bad := func(file, text string) map[string][]byte {
		m := both()
		m[file] = []byte(text)
		return m
	}
	for _, c := range []struct {
		name   string
		f      pushFake
		detail string
		calls  []string
	}{
		{"the tree of the setup is not read", pushFake{lastErr: boom}, "setup/tree.tsv: boom", []string{"last"}},
		{"the work area cannot be made sure of", pushFake{last: sha1, areaErr: errors.New("layup-setup is at x, setup/tree.tsv ends at y\nmore")},
			"layup-setup is at x, setup/tree.tsv ends at y", []string{"last", "area"}},
		{"the forge cannot be read", pushFake{last: sha1, headErr: boom}, "the forge: boom", []string{"last", "area", "head"}},
		{"the apply outputs cannot be read", pushFake{last: sha1, head: sha1, appliesErr: boom}, "the apply outputs: boom", []string{"last", "area", "head", "applies"}},
		{"commands.sh cannot be read while the step waits", pushFake{last: sha1, head: sha2, applies: both(), commandsErr: boom}, "commands.sh: boom", []string{"last", "area", "head", "applies", "commands"}},
		{"an output is not a JSON object", pushFake{last: sha1, head: sha1, applies: bad("default.json", "[]")}, "default.json", []string{"last", "area", "head", "applies"}},
		{"an output is of another repository", pushFake{last: sha1, head: sha1, applies: bad("records.json", `{"source":"acme/other"}`)}, "records.json", []string{"last", "area", "head", "applies"}},
		{"the records commit fails", pushFake{last: sha1, head: sha1, applies: both(), commitErr: boom}, "the records commit: boom", []string{"last", "area", "head", "applies", "commit"}},
	} {
		f := c.f
		got := f.job().do()
		if got.Name != "scaffold-push" || got.Result != "fail" || !strings.Contains(got.Detail, c.detail) || !slices.Equal(f.called, c.calls) {
			t.Errorf("%s: %+v after the calls %q; want fail with %q after %q", c.name, got, f.called, c.detail, c.calls)
		}
		if c.name != "the records commit fails" && len(f.commits) != 0 {
			t.Errorf("%s: a commit was made", c.name)
		}
	}
}
