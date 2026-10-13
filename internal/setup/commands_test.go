package setup

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/work"
)

// The commands for the Operator (task T-56dj, #275; docs/spec/scaffold.md, The
// commands for the Operator): commands.sh starts with set -eu in both phases, and in
// phase 2 its commands read the head of the default branch, push the setup commits
// once, and apply each ruleset only when its output file is missing. The file is run
// with sh and with bash, with stand-ins of git and gh that log each call and answer
// from the environment. Each case breaks one rule only.

// stubGit is the stand-in of git: it logs its arguments, skips -C and its directory,
// and answers ls-remote, rev-list and rev-parse from the environment.
const stubGit = `#!/bin/sh
echo "git $*" >> "$STUB_LOG"
[ "$1" = "-C" ] && shift 2
case "$1" in
ls-remote)
	[ -n "${STUB_LSREMOTE_FAIL:-}" ] && exit "$STUB_LSREMOTE_FAIL"
	[ -n "${STUB_HEAD:-}" ] && printf '%s\trefs/heads/main\n' "$STUB_HEAD"
	exit 0 ;;
rev-list) printf '%s\n' "$STUB_ROOT" ;;
rev-parse) printf '%s\n' "$STUB_SETUP" ;;
push) [ -n "${STUB_PUSH_FAIL:-}" ] && exit 1; echo pushed >> "$STUB_PUSHED" ;;
remote | add) ;;
*) echo "the stand-in of git does not know: $*" >&2; exit 99 ;;
esac
`

// stubGh is the stand-in of gh: it logs its arguments, answers a ruleset apply with a
// small JSON object by the file of --input, and, for the file named by STUB_GH_FAIL_ON,
// writes a part of an answer and exits 1.
const stubGh = `#!/bin/sh
echo "gh $*" >> "$STUB_LOG"
case "$*" in
"api --method POST /repos/acme/widget/rulesets --input "* | "api --method POST repos/acme/widget/rulesets --input "*) ;;
*) echo "the stand-in of gh does not know: $*" >&2; exit 99 ;;
esac
for last; do :; done
name=$(basename "$last")
if [ "${STUB_GH_FAIL_ON:-}" = "$name" ]; then printf '{"id":'; exit 1; fi
case "$name" in
ruleset-default.json) echo '{"id":1,"source":"acme/widget","name":"default"}' ;;
ruleset-records.json) echo '{"id":2,"source":"acme/widget","name":"records"}' ;;
esac
`

const (
	rootSHA  = "1111111111111111111111111111111111111111"
	setupSHA = "2222222222222222222222222222222222222222"
	otherSHA = "3333333333333333333333333333333333333333"
)

// phase2Commands gives the phase 2 inputs of these tests, with the host dir given.
func phase2Commands(host string) Phase2 {
	p := phase2App(42)
	p.Host = host
	return p
}

// s13Record is a record in which S13 is done, so its commands are written.
func s13Record() work.Record {
	return work.Record{{"S01", "name", "acme/widget", "computed", "records:start/start.tsv"}, {"S13", "done", "the ruleset files", "step", "—"}}
}

// shell is a directory with the stand-ins on a PATH of its own, and the files that
// they write.
type shell struct {
	t    *testing.T
	dir  string
	host string
	log  string
}

func newShell(t *testing.T) *shell {
	t.Helper()
	d := t.TempDir()
	bin := filepath.Join(d, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, text := range map[string]string{"git": stubGit, "gh": stubGh} {
		if err := os.WriteFile(filepath.Join(bin, name), []byte(text), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return &shell{t: t, dir: d, host: filepath.Join(d, "host"), log: filepath.Join(d, "calls.log")}
}

func (s *shell) applies() string { return filepath.Join(s.host, "applies", "acme", "widget") }

// run runs script with the shell sh, from the directory of the test, with the
// answers env (NAME=value), and gives its exit status and its standard error.
func (s *shell) run(sh, script string, env ...string) (int, string) {
	s.t.Helper()
	file := filepath.Join(s.dir, "commands.sh")
	if err := os.WriteFile(file, []byte(script), 0o644); err != nil {
		s.t.Fatal(err)
	}
	cmd := exec.Command(sh, file)
	cmd.Dir = s.dir
	cmd.Env = append(os.Environ(), "PATH="+filepath.Join(s.dir, "bin")+":"+os.Getenv("PATH"), "STUB_LOG="+s.log, "STUB_PUSHED="+filepath.Join(s.dir, "pushed.log"),
		"STUB_ROOT="+rootSHA, "STUB_SETUP="+setupSHA)
	cmd.Env = append(cmd.Env, env...)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err == nil {
		return 0, stderr.String()
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode(), stderr.String()
	}
	s.t.Fatalf("%s did not run: %v", sh, err)
	return -1, ""
}

// calls gives the lines of the log of the stand-ins, with the directory of the test cut.
func (s *shell) calls() []string {
	data, _ := os.ReadFile(s.log)
	var out []string
	for _, l := range strings.Split(strings.TrimSpace(strings.ReplaceAll(string(data), s.dir, "~")), "\n") {
		if l != "" {
			out = append(out, l)
		}
	}
	return out
}

func (s *shell) pushes() int {
	data, _ := os.ReadFile(filepath.Join(s.dir, "pushed.log"))
	return strings.Count(string(data), "pushed")
}

func (s *shell) exists(name string) bool {
	_, err := os.Stat(filepath.Join(s.applies(), name))
	return err == nil
}

// shells gives the shells of the cases: sh always, and bash where the host has it.
func shells(t *testing.T) []string {
	t.Helper()
	out := []string{"sh"}
	if _, err := exec.LookPath("bash"); err == nil {
		out = append(out, "bash")
	}
	return out
}

// commandsOf gives the text of commands.sh that the writer makes from the steps and the record.
func commandsOf(steps map[string]Step, record work.Record) string { return commandsText(steps, record) }

// The text of phase 2, byte for byte.
func TestTheCommandsOfS13InPhase2(t *testing.T) {
	got := commandsOf(StepsPhase2(phase2Commands("/h"), Calls{}), s13Record())
	want := `set -eu
# the head of the default branch on the forge, and the push of the setup commits onto it: a fast-forward from the root commit when the head is that commit, nothing when it is the head of layup-setup, and a stop when it is neither
remote=$(git -C '/h/work/acme/widget/target' ls-remote origin 'refs/heads/main')
head=$(printf '%s\n' "$remote" | cut -f1)
root=$(git -C '/h/work/acme/widget/target' rev-list --max-parents=0 refs/heads/layup-setup)
setup=$(git -C '/h/work/acme/widget/target' rev-parse refs/heads/layup-setup)
if [ "$head" = "$root" ]; then
	git -C '/h/work/acme/widget/target' push origin layup-setup:main
elif [ "$head" != "$setup" ]; then
	echo "the head of main on the forge is '$head', neither the root commit '$root' nor the head of layup-setup '$setup'" >&2
	exit 1
fi
# the directory of the apply outputs
mkdir -p '/h/applies/acme/widget'
# the apply of the ruleset of the default branch, with the Operator's login, when its output file is missing: the output goes to a temporary file, renamed when gh exits 0
if [ ! -e '/h/applies/acme/widget/default.json' ]; then
	gh api --method POST /repos/acme/widget/rulesets --input '/h/work/acme/widget/out/ruleset-default.json' > '/h/applies/acme/widget/default.json.tmp'
	mv '/h/applies/acme/widget/default.json.tmp' '/h/applies/acme/widget/default.json'
fi
# the apply of the ruleset of the records branch, with the Operator's login, when its output file is missing: the output goes to a temporary file, renamed when gh exits 0
if [ ! -e '/h/applies/acme/widget/records.json' ]; then
	gh api --method POST /repos/acme/widget/rulesets --input '/h/work/acme/widget/out/ruleset-records.json' > '/h/applies/acme/widget/records.json.tmp'
	mv '/h/applies/acme/widget/records.json.tmp' '/h/applies/acme/widget/records.json'
fi
`
	if got != want {
		t.Errorf("commands.sh of phase 2:\n%s\nwant\n%s", got, want)
	}
	// four entries: the head and the push in order 2, the directory and the two applies in order 4
	var orders []int
	for _, c := range StepsPhase2(phase2Commands("/h"), Calls{})["S13"].Commands(s13Record()) {
		orders = append(orders, c.Order)
		if c.Comment == "" {
			t.Errorf("a command of S13 with no comment: %q", c.Text)
		}
	}
	if !reflect.DeepEqual(orders, []int{2, 4, 4, 4}) {
		t.Errorf("the orders of the commands of S13 in phase 2 are %v; want [2 4 4 4]", orders)
	}
}

// commands.sh starts with set -eu in both phases, and is empty when no done step has a command.
func TestTheFileStartsWithSetEUInBothPhases(t *testing.T) {
	one := commandsOf(Steps(Brief{}, Calls{}), work.Record{{"S01", "name", "acme/widget", "computed", "x"}, {"S02", "pin.commit", rootSHA, "computed", "x"}, {"S03", "done", "root", "step", "—"}})
	two := commandsOf(StepsPhase2(phase2Commands("/h"), Calls{}), s13Record())
	for name, text := range map[string]string{"phase 1": one, "phase 2": two} {
		if first, _, _ := strings.Cut(text, "\n"); first != "set -eu" {
			t.Errorf("%s: the first line of commands.sh is %q; want set -eu", name, first)
		}
	}
	if got := commandsOf(StepsPhase2(phase2Commands("/h"), Calls{}), work.Record{{"S01", "name", "acme/widget", "computed", "x"}}); got != "" {
		t.Errorf("no done step has a command: commands.sh is %q; want it empty", got)
	}
}

// A first run, with the head at the root commit: one push, then both applies, whose
// outputs are the files of DIR/applies and leave no temporary file.
func TestAFirstRunPushesOnceAndAppliesBothRulesets(t *testing.T) {
	script := commandsOf(StepsPhase2(phase2Commands("/h"), Calls{}), s13Record())
	for _, sh := range shells(t) {
		s := newShell(t)
		script := strings.ReplaceAll(script, "/h/", s.host+"/")
		code, stderr := s.run(sh, script, "STUB_HEAD="+rootSHA)
		got := s.calls()
		want := []string{"git -C ~/host/work/acme/widget/target ls-remote origin refs/heads/main", "git -C ~/host/work/acme/widget/target rev-list --max-parents=0 refs/heads/layup-setup",
			"git -C ~/host/work/acme/widget/target rev-parse refs/heads/layup-setup", "git -C ~/host/work/acme/widget/target push origin layup-setup:main",
			"gh api --method POST /repos/acme/widget/rulesets --input ~/host/work/acme/widget/out/ruleset-default.json",
			"gh api --method POST /repos/acme/widget/rulesets --input ~/host/work/acme/widget/out/ruleset-records.json"}
		if code != 0 || !reflect.DeepEqual(got, want) || s.pushes() != 1 {
			t.Errorf("%s, a first run: exit %d (%s), calls\n%s\nwant exit 0 and\n%s", sh, code, stderr, strings.Join(got, "\n"), strings.Join(want, "\n"))
		}
		for name, out := range map[string]string{"default.json": `{"id":1,"source":"acme/widget","name":"default"}`, "records.json": `{"id":2,"source":"acme/widget","name":"records"}`} {
			data, _ := os.ReadFile(filepath.Join(s.applies(), name))
			if strings.TrimSpace(string(data)) != out || s.exists(name+".tmp") {
				t.Errorf("%s: %s is %q (a temporary file: %v); want %q and no temporary file", sh, name, data, s.exists(name+".tmp"), out)
			}
		}
	}
}

// A run again after the push (the head is the head of layup-setup) pushes nothing and
// applies both; after one apply (default.json exists) it applies records.json only; after
// both it applies nothing and pushes nothing.
func TestARunAgainPushesNothingAndAppliesOnlyWhatIsMissing(t *testing.T) {
	script := commandsOf(StepsPhase2(phase2Commands("/h"), Calls{}), s13Record())
	for _, sh := range shells(t) {
		for _, c := range []struct {
			name    string
			have    []string
			applied []string
		}{
			{"after the push", nil, []string{"ruleset-default.json", "ruleset-records.json"}},
			{"after one apply", []string{"default.json"}, []string{"ruleset-records.json"}},
			{"after both applies", []string{"default.json", "records.json"}, nil},
		} {
			s := newShell(t)
			os.MkdirAll(s.applies(), 0o755)
			for _, f := range c.have {
				os.WriteFile(filepath.Join(s.applies(), f), []byte(`{"id":9}`), 0o644)
			}
			code, stderr := s.run(sh, strings.ReplaceAll(script, "/h/", s.host+"/"), "STUB_HEAD="+setupSHA)
			var gh []string
			for _, l := range s.calls() {
				if strings.HasPrefix(l, "gh ") {
					gh = append(gh, filepath.Base(l))
				}
			}
			if code != 0 || s.pushes() != 0 || !reflect.DeepEqual(gh, c.applied) {
				t.Errorf("%s, a run again %s: exit %d (%s), %d pushes, applies %v; want exit 0, no push, applies %v", sh, c.name, code, stderr, s.pushes(), gh, c.applied)
			}
			if c.name == "after one apply" {
				if data, _ := os.ReadFile(filepath.Join(s.applies(), "default.json")); string(data) != `{"id":9}` {
					t.Errorf("%s: default.json was written again: %q", sh, data)
				}
			}
		}
	}
}

// A gh that writes part of its output and exits 1 leaves no output file, stops the rest,
// and the run again applies that ruleset (and the next).
func TestAFailedApplyLeavesNoOutputFileAndIsDoneAgain(t *testing.T) {
	script := commandsOf(StepsPhase2(phase2Commands("/h"), Calls{}), s13Record())
	for _, sh := range shells(t) {
		s := newShell(t)
		script := strings.ReplaceAll(script, "/h/", s.host+"/")
		code, _ := s.run(sh, script, "STUB_HEAD="+setupSHA, "STUB_GH_FAIL_ON=ruleset-default.json")
		var gh int
		for _, l := range s.calls() {
			if strings.HasPrefix(l, "gh ") {
				gh++
			}
		}
		if code == 0 || s.exists("default.json") || s.exists("records.json") || gh != 1 {
			t.Errorf("%s, a failed apply: exit %d, default.json %v, records.json %v, %d gh calls; want a failure, no output file and one call", sh, code, s.exists("default.json"), s.exists("records.json"), gh)
		}
		if code, stderr := s.run(sh, script, "STUB_HEAD="+setupSHA); code != 0 || !s.exists("default.json") || !s.exists("records.json") {
			t.Errorf("%s, the run again: exit %d (%s), default.json %v, records.json %v; want both applied", sh, code, stderr, s.exists("default.json"), s.exists("records.json"))
		}
	}
}

// A head that is neither the root commit nor the head of layup-setup, also an empty one, stops the
// file with exit 1, and the message names the head, the root commit and the head of layup-setup;
// nothing is pushed and nothing is applied.
func TestAHeadThatIsNeitherStopsTheFile(t *testing.T) {
	script := commandsOf(StepsPhase2(phase2Commands("/h"), Calls{}), s13Record())
	for _, sh := range shells(t) {
		for name, head := range map[string]string{"another commit": otherSHA, "no head": ""} {
			s := newShell(t)
			code, stderr := s.run(sh, strings.ReplaceAll(script, "/h/", s.host+"/"), "STUB_HEAD="+head)
			want := "the head of main on the forge is '" + head + "', neither the root commit '" + rootSHA + "' nor the head of layup-setup '" + setupSHA + "'"
			var gh int
			for _, l := range s.calls() {
				if strings.HasPrefix(l, "gh ") {
					gh++
				}
			}
			if code != 1 || strings.TrimSpace(stderr) != want || s.pushes() != 0 || gh != 0 || s.exists("default.json") {
				t.Errorf("%s, %s: exit %d, stderr %q, %d pushes, %d gh calls; want exit 1, %q, and nothing done", sh, name, code, strings.TrimSpace(stderr), s.pushes(), gh, want)
			}
		}
	}
}

// A command that fails stops the commands after it: a read of the forge that fails stops the file
// with its status, and a push that fails stops it before the applies.
func TestACommandThatFailsStopsTheRest(t *testing.T) {
	script := commandsOf(StepsPhase2(phase2Commands("/h"), Calls{}), s13Record())
	for _, sh := range shells(t) {
		s := newShell(t)
		code, _ := s.run(sh, strings.ReplaceAll(script, "/h/", s.host+"/"), "STUB_LSREMOTE_FAIL=128")
		if code != 128 || len(s.calls()) != 1 || s.pushes() != 0 {
			t.Errorf("%s, a failed ls-remote: exit %d, calls %v; want exit 128 and the one call", sh, code, s.calls())
		}
		s = newShell(t)
		code, _ = s.run(sh, strings.ReplaceAll(script, "/h/", s.host+"/"), "STUB_HEAD="+rootSHA, "STUB_PUSH_FAIL=1")
		var gh int
		for _, l := range s.calls() {
			if strings.HasPrefix(l, "gh ") {
				gh++
			}
		}
		if code == 0 || gh != 0 || s.exists("default.json") {
			t.Errorf("%s, a failed push: exit %d, %d gh calls; want a failure and no apply", sh, code, gh)
		}
	}
}

// The three steps of phase 1 give five commands, which hold under set -eu: the file runs to the end
// with the stand-ins, in the order of setup.md; and a command that fails stops the rest.
func TestThePhase1CommandsHoldUnderSetEU(t *testing.T) {
	record := work.Record{{"S01", "name", "acme/widget", "computed", "x"}, {"S02", "pin.commit", rootSHA, "computed", "x"},
		{"S03", "done", "root", "step", "—"}, {"S13", "done", "ruleset", "step", "—"}, {"S15", "done", "records", "step", "—"}}
	script := commandsOf(Steps(Brief{}, Calls{}), record)
	for _, sh := range shells(t) {
		s := newShell(t)
		code, stderr := s.run(sh, script)
		want := []string{"git -C target remote add origin https://github.com/acme/widget.git", "git -C target push origin main", "git -C target push origin layup-setup:main",
			"git -C target push origin layup-records", "gh api --method POST repos/acme/widget/rulesets --input out/ruleset-default.json"}
		if code != 0 || !reflect.DeepEqual(s.calls(), want) {
			t.Errorf("%s, phase 1: exit %d (%s), calls\n%s\nwant exit 0 and\n%s", sh, code, stderr, strings.Join(s.calls(), "\n"), strings.Join(want, "\n"))
		}
		// the first command fails: the rest does not run
		s = newShell(t)
		code, _ = s.run(sh, strings.Replace(script, "git -C target remote add", "git -C target unknown-subcommand add", 1))
		if code != 99 || len(s.calls()) != 1 {
			t.Errorf("%s, phase 1 with a first command that fails: exit %d, calls %v; want exit 99 and the one call", sh, code, s.calls())
		}
	}
}

// The default branch is the one that the forge gives, not main: the read, the push and the
// message name it, a branch with a character that sh reads is quoted, and the file runs with it.
func TestTheCommandsNameTheDefaultBranchOfTheForge(t *testing.T) {
	for name, branch := range map[string]string{"a plain name": "trunk", "a name with a dollar": "a$b", "a name with a dollar and a double quote": "a$b\"c"} {
		p := phase2Commands("/h")
		p.DefaultBranch = branch
		script := commandsOf(StepsPhase2(p, Calls{}), s13Record())
		for _, sh := range shells(t) {
			s := newShell(t)
			code, _ := s.run(sh, strings.ReplaceAll(script, "/h/", s.host+"/"), "STUB_HEAD="+rootSHA)
			var read, push string
			for _, l := range s.calls() {
				if strings.Contains(l, " ls-remote ") {
					read = l
				}
				if strings.Contains(l, " push ") {
					push = l
				}
			}
			if code != 0 || !strings.HasSuffix(read, "ls-remote origin refs/heads/"+branch) || !strings.HasSuffix(push, "push origin layup-setup:"+branch) {
				t.Errorf("%s, %s: exit %d, the read %q, the push %q; want the branch %q in both", sh, name, code, read, push, branch)
			}
			s = newShell(t)
			code, stderr := s.run(sh, strings.ReplaceAll(script, "/h/", s.host+"/"), "STUB_HEAD="+otherSHA)
			if want := "the head of " + branch + " on the forge is '" + otherSHA + "'"; code != 1 || !strings.HasPrefix(strings.TrimSpace(stderr), want) {
				t.Errorf("%s, %s: exit %d, stderr %q; want exit 1 and a message that begins %q", sh, name, code, strings.TrimSpace(stderr), want)
			}
		}
	}
}

// S13 in phase 2 is a fail when an input of its commands is not usable: a host directory that is
// not an absolute path (the file would depend on the directory that the Operator runs it from) and an
// empty default branch (the commands would name no ref); no file is written.
func TestPhase2S13RefusesInputsThatMakeNoCommands(t *testing.T) {
	for _, c := range []struct {
		name     string
		change   func(*Phase2)
		evidence string
	}{
		{"a host directory that is relative", func(p *Phase2) { p.Host = "host/dir" }, "S13: the host directory of the commands, host/dir, is not an absolute path"},
		{"no host directory", func(p *Phase2) { p.Host = "" }, "S13: the host directory of the commands, , is not an absolute path"},
		{"no default branch", func(p *Phase2) { p.DefaultBranch = "" }, "S13: the default branch of the forge is empty"},
	} {
		p := phase2Commands("/h")
		c.change(&p)
		f := &fakeRepo{shows: map[string]string{"docs/gates.tsv": twoKinds}}
		f.install(t)
		o := StepsPhase2(p, Calls{})["S13"].Run(Input{Dir: "w", Record: lastRecord()})
		if o.Kind != Fail || o.Evidence != c.evidence || len(f.files) != 0 {
			t.Errorf("%s: %s %q, %d files; want fail %q and no file", c.name, o.Kind, o.Evidence, len(f.files), c.evidence)
		}
	}
}
