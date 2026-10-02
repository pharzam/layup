package verify

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/pharzam/layup/internal/work"
)

const (
	commitA = "a95965534b14b0bf14ad74da0c9a45b5f4aedf88"
	treeA   = "cb89dd90fb277f5d5767688e3e56bb2872814180"
	treeB   = "8ffb250afd584da8b418bc220fb6d72e802924ce"
	source  = "https://github.com/pharzam/armature"
)

// stubHistory answers the git questions of check pin.
type stubHistory struct {
	shallow bool
	roots   []string
	trees   map[string]string
	err     error
}

func (h stubHistory) IsShallow() (bool, error)       { return h.shallow, h.err }
func (h stubHistory) RootCommits() ([]string, error) { return h.roots, h.err }
func (h stubHistory) Tree(c string) (string, error) {
	if h.err != nil {
		return "", h.err
	}
	return h.trees[c], nil
}

// oneRoot is a history with one root commit whose tree is treeA.
var oneRoot = stubHistory{roots: []string{"r1"}, trees: map[string]string{"r1": treeA}}

func pinFile(lines ...string) fstest.MapFS {
	return fstest.MapFS{pinPath: {Data: []byte(strings.Join(lines, "\n") + "\n")}}
}

var goodPin = []string{"# a comment line with commit=x", "source=" + source, "commit=" + commitA, "tree=" + treeA,
	"method=npx degit pharzam/armature#" + commitA, "date=2026-09-23"}

func same(t *testing.T, name string, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s:\n got %q\nwant %q", name, got, want)
	}
}

// Each finding of check_pin, in the order of the sh function; the first is
// the reason of the row.
func TestPinFindings(t *testing.T) {
	short := append([]string{"commit=a959655"}, goodPin...)
	for _, c := range []struct {
		name string
		fsys fstest.MapFS
		h    history
		want []string
	}{
		{"a good pin", pinFile(goodPin...), oneRoot, nil},
		{"no pin file", fstest.MapFS{}, oneRoot, []string{"missing: docs/setup/armature.pin is absent"}},
		{"a key twice, the first commit short", pinFile(short...), oneRoot,
			[]string{"key: commit appears 2 times, expected 1", "commit: not 40 hexadecimal characters: a959655"}},
		{"a key that is missing", pinFile(goodPin[:5]...), oneRoot, []string{"key: date appears 0 times, expected 1"}},
		{"a shallow clone", pinFile(goodPin...), stubHistory{shallow: true, roots: []string{"r1"}},
			[]string{"shallow: the clone is shallow, so the root commit is not visible (fetch the full history)"}},
		{"two root commits", pinFile(goodPin...), stubHistory{roots: []string{"r1", "r2"}},
			[]string{"tree: expected one root commit, found 2"}},
		{"git fails", pinFile(goodPin...), stubHistory{err: errors.New("fatal")},
			[]string{"tree: expected one root commit, found 0"}},
		{"another tree", pinFile(goodPin...), stubHistory{roots: []string{"r1"}, trees: map[string]string{"r1": treeB}},
			[]string{"tree: pin names " + treeA + ", root commit has " + treeB}},
	} {
		same(t, c.name, pinFindings(c.fsys, c.h), c.want)
	}
}

// Each finding of check_kit_history, in the order of the sh function.
func TestKitHistoryFindings(t *testing.T) {
	const link = "github.com/pharzam/armature"
	good := fstest.MapFS{
		"docs/tasks/T-0001.md":    {Data: []byte("# T-0001\n")},
		"docs/tasks/T-0002.md":    {Data: []byte("# T-0002\n")},
		"docs/tasks/backlog.md":   {Data: []byte("- **T-0001** — a task\n")},
		"docs/tasks/completed.md": {Data: []byte("- **2026-01-01** — **T-0002** — a task\n")},
	}
	with := func(files map[string]string) fstest.MapFS {
		m := fstest.MapFS{}
		for k, v := range good {
			m[k] = v
		}
		for k, v := range files {
			m[k] = &fstest.MapFile{Data: []byte(v)}
		}
		return m
	}
	linked := "- **T-0001** — a task (https://" + link + "/issues/1)\n"
	linkedDone := "- **2026-01-01** — **T-0002** — a task (https://" + link + "/issues/2)\n"
	for _, c := range []struct {
		name string
		fsys fstest.MapFS
		link string
		want []string
	}{
		{"no history of the kit", good, link, nil},
		{"the two directories of the kit", with(map[string]string{"docs/decisions/D-0001.md": "x", "docs/audit/README.md": "x"}), link,
			[]string{"decisions: docs/decisions/ exists (kit step 4 deletes it)", "audit: docs/audit/ exists (kit step 4 deletes it)"}},
		{"two task files with no line, in the order of their names", with(map[string]string{"docs/tasks/T-zzzz.md": "x", "docs/tasks/T-aaaa.md": "x"}), link,
			[]string{"orphan: docs/tasks/T-aaaa.md has no line with T-aaaa in backlog.md or completed.md",
				"orphan: docs/tasks/T-zzzz.md has no line with T-zzzz in backlog.md or completed.md"}},
		{"a directory with the name of a task file", with(map[string]string{"docs/tasks/T-dir.md/x": "x"}), link, nil},
		{"both indexes link the baseline", with(map[string]string{"docs/tasks/backlog.md": linked, "docs/tasks/completed.md": linkedDone}), link,
			[]string{"kit-link: docs/tasks/backlog.md links " + link + "/ (a kit task or note)",
				"kit-link: docs/tasks/completed.md links " + link + "/ (a kit task or note)"}},
		{"a link to the repository itself (round 1, finding 1)", with(map[string]string{"docs/tasks/backlog.md": "- **T-0001** ([the baseline](https://" + link + "))\n"}), link,
			[]string{"kit-link: docs/tasks/backlog.md links " + link + "/ (a kit task or note)"}},
		{"a repository whose name only starts with the same text", with(map[string]string{
			"docs/tasks/backlog.md": "- **T-0001** (https://" + link + "2/issues/1, https://" + link + "-docs/x)\n"}), link, nil},
		{"other repositories (round 2, finding 1)", with(map[string]string{"docs/tasks/backlog.md": "- **T-0001** (https://" + link +
			".testing/issues/1, https://" + link + ".gitlab/x, https://other" + link + "/issues/1)\n"}), link, nil},
		{"another host that ends with the host of the baseline (round 3, finding 1)", with(map[string]string{
			"docs/tasks/backlog.md": "- **T-0001** (https://my." + link + "/issues/1, https://www." + link + ")\n"}), link, nil},
		{"a link with .git", with(map[string]string{"docs/tasks/backlog.md": "- **T-0001** (https://" + link + ".git)\n"}), link,
			[]string{"kit-link: docs/tasks/backlog.md links " + link + "/ (a kit task or note)"}},
		{"a link at the end of a sentence", with(map[string]string{"docs/tasks/backlog.md": "- **T-0001**: see " + link + ".\n"}), link,
			[]string{"kit-link: docs/tasks/backlog.md links " + link + "/ (a kit task or note)"}},
		{"no link text reads no link", with(map[string]string{"docs/tasks/backlog.md": linked}), "", nil},
	} {
		same(t, c.name, kitHistoryFindings(c.fsys, c.link), c.want)
	}
}

// Each finding of check_identity, in the order of the sh function.
func TestIdentityFindings(t *testing.T) {
	const pinLink = "the [pin](docs/setup/armature.pin)\n"
	for _, c := range []struct {
		name string
		fsys fstest.MapFS
		want []string
	}{
		{"the target's entry files", fstest.MapFS{"README.md": {Data: []byte(pinLink)}, "AGENTS.md": {Data: []byte("x\n")}}, nil},
		{"each phrase of the kit", fstest.MapFS{
			"README.md": {Data: []byte("This repository is a generic **template**. Agent context for **Armature**.\n")},
			"AGENTS.md": {Data: []byte("A domain-free **template**, not a product.\n")}},
			[]string{`kit: README.md says the repository is the Armature kit ("Agent context for **Armature**")`,
				`kit: README.md says the repository is the Armature kit ("This repository is a generic **template**")`,
				`kit: AGENTS.md says the repository is the Armature kit ("A domain-free **template**, not a product")`,
				"pin: README.md has no link to docs/setup/armature.pin"}},
		{"no README.md", fstest.MapFS{"AGENTS.md": {Data: []byte("x\n")}}, []string{"pin: README.md has no link to docs/setup/armature.pin"}},
	} {
		same(t, c.name, identityFindings(c.fsys), c.want)
	}
}

// targetRecord is a record with the rows that the three checks read.
func targetRecord(t *testing.T, drop string, change map[string]string) work.Record {
	t.Helper()
	rows := [][]string{
		{"S01", "name", "acme"},
		{"S02", "pin.source", source},
		{"S02", "pin.commit", commitA},
		{"S02", "pin.tree", treeA},
		{"S02", "pin.time", "2026-10-02T09:30:00Z"},
	}
	var r work.Record
	for _, row := range rows {
		if row[1] == drop {
			continue
		}
		if v, ok := change[row[1]]; ok {
			row[2] = v
		}
		r = append(r, []string{row[0], row[1], row[2], "computed", "x"})
	}
	return r
}

var targetPin = []string{"source=" + source, "commit=" + commitA, "tree=" + treeA,
	"method=git clone " + source + ", checkout " + commitA + ", .git removed", "date=2026-10-02"}

// The target's part of check pin: the pin equals the record rows of S02, its
// date is the date of pin.time, and its method is the text of NFR-006 (D3).
func TestThePinOfATarget(t *testing.T) {
	npx := replaced(targetPin, 3, "method=npx degit pharzam/armature#"+commitA)
	for _, c := range []struct {
		name   string
		fsys   fstest.MapFS
		record work.Record
		want   []string
	}{
		{"the pin of the record", pinFile(targetPin...), targetRecord(t, "", nil), nil},
		{"no pin file: the core says it", fstest.MapFS{}, targetRecord(t, "", nil), nil},
		{"another commit in the record", pinFile(targetPin...), targetRecord(t, "", map[string]string{"pin.commit": strings.Repeat("1", 40)}),
			[]string{"commit: the pin names " + commitA + ", the record row pin.commit is " + strings.Repeat("1", 40),
				"method: the pin names git clone " + source + ", checkout " + commitA + ", .git removed, not git clone " + source + ", checkout " + strings.Repeat("1", 40) + ", .git removed"}},
		{"the method of LAYUP's own pin", pinFile(npx...), targetRecord(t, "", nil),
			[]string{"method: the pin names npx degit pharzam/armature#" + commitA + ", not git clone " + source + ", checkout " + commitA + ", .git removed"}},
		{"another date", pinFile(replaced(targetPin, 4, "date=2026-09-23")...), targetRecord(t, "", nil),
			[]string{"date: the pin names 2026-09-23, the date of pin.time is 2026-10-02"}},
		{"a time of another form", pinFile(targetPin...), targetRecord(t, "", map[string]string{"pin.time": "2026-10-02 09:30"}),
			[]string{"record: pin.time is not of the form YYYY-MM-DDTHH:MM:SSZ: 2026-10-02 09:30"}},
		{"no row pin.commit", pinFile(targetPin...), targetRecord(t, "pin.commit", nil), []string{"record: no value at S02 pin.commit"}},
	} {
		same(t, c.name, pinTarget(c.fsys, c.record), c.want)
	}
}

// replaced gives a copy of lines with line i replaced.
func replaced(lines []string, i int, line string) []string {
	out := append([]string{}, lines...)
	out[i] = line
	return out
}

// The link text of the kit-link rule is the target's own baseline (D3).
func TestTheKitLinkOfATarget(t *testing.T) {
	for in, want := range map[string]string{
		"https://github.com/pharzam/armature":     "github.com/pharzam/armature",
		"https://github.com/pharzam/armature.git": "github.com/pharzam/armature",
		"https://github.com/pharzam/armature/":    "github.com/pharzam/armature",
		"file:///tmp/w/baseline":                  "/tmp/w/baseline",
		"":                                        "",
	} {
		if got := kitLink(in); got != want {
			t.Errorf("kitLink(%q) = %q, want %q", in, got, want)
		}
	}
}

// The three checks of a target: the core, then the target's part (D3).
func TestTheChecksOfATarget(t *testing.T) {
	readme := "# acme\n\nThe [pin](docs/setup/armature.pin).\n"
	tree := fstest.MapFS{
		"README.md":             {Data: []byte(readme)},
		"docs/tasks/backlog.md": {Data: []byte("- **T-0001** — a task (https://github.com/pharzam/armature/issues/1)\n")},
		pinPath:                 pinFile(targetPin...)[pinPath],
	}
	in := input{fsys: tree, repo: oneRoot, record: targetRecord(t, "", nil)}
	same(t, "pin", checkPin(in), nil)
	same(t, "kit-history", checkKitHistory(in), []string{"kit-link: docs/tasks/backlog.md links github.com/pharzam/armature/ (a kit task or note)"})
	same(t, "identity", checkIdentity(in), nil)

	in.record = targetRecord(t, "pin.source", map[string]string{"name": "acme-billing"})
	same(t, "kit-history with no pin.source", checkKitHistory(in), []string{"record: no value at S02 pin.source"})
	same(t, "identity with another name", checkIdentity(in), []string{"name: README.md does not hold the name of the target, acme-billing"})
	in.record = targetRecord(t, "name", nil)
	same(t, "identity with no name", checkIdentity(in), []string{"record: no value at S01 name"})
}
