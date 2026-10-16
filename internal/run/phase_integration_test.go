//go:build integration

package run

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/ledger"
	"github.com/pharzam/layup/internal/route"
)

// headsForge is the stand-in forge that keeps, for each comment, the head of
// layup-records at the time it was posted, so a test sees the order of the
// records commit and the comment.
type headsForge struct {
	*standInForge
	t     *testing.T
	bare  string
	heads []string
}

func (f *headsForge) Comment(ctx context.Context, issue int, body string) (int64, error) {
	f.heads = append(f.heads, strings.TrimSpace(gitOut(f.t, f.bare, "rev-parse", "layup-records")))
	return f.standInForge.Comment(ctx, issue, body)
}

// phaseHarness gives a harness register row whose program answers a probe
// as probeHarness does, and runs body for a session of a phase, with $in the
// input file and $out the output file that the test's prompt names.
func phaseHarness(t *testing.T, body string) []string {
	bin := t.TempDir()
	os.WriteFile(filepath.Join(bin, "version"), []byte("#!/bin/sh\necho 1.0.0\n"), 0o755)
	os.WriteFile(filepath.Join(bin, "harness"), []byte(`#!/bin/sh
p=$(cat)
case "$p" in
"This is a probe"*)
	result=$(printf '%s\n' "$p" | sed -n 's/^Write one file, \(.*\), of tab-separated values, with the header line$/\1/p')
	token=$(printf '%s\n' "$p" | sed -n 's/.*one row "token<TAB>\([0-9a-f]*\)".*/\1/p')
	printf 'kind\tvalue\ntoken\t%s\nfile\tAGENTS.md\n' "$token" > "$result" ;;
*)
	in=$(printf '%s\n' "$p" | sed -n 's/^input //p')
	out=$(printf '%s\n' "$p" | sed -n 's/^output //p')
	`+body+` ;;
esac
`), 0o755)
	return []string{"fake", "", "1", filepath.Join(bin, "harness"), "stdin", filepath.Join(bin, "version"), "", "", "LAYUP-TEST-RULES.md", "", "none", "subscription", ""}
}

// phaseWorld gives the world of the sessions with the routing register of
// one pair for systems-architect, probed and passed by the step probe, a
// forge that keeps the head of the records at each comment, and a clock a
// minute after Start.
func phaseWorld(t *testing.T, body string) (*world, *Sessions, *headsForge, worldAt, []string) {
	w, r, f, a := sessionWorld(t)
	hf := &headsForge{standInForge: f, t: t, bare: w.bare}
	r.Forge = hf
	h := phaseHarness(t, body)
	routing := [][]string{{"systems-architect", "reasoning", "1", "fake", "m1"}}
	if p, _, _, err := r.ProbeStep(context.Background(), [][]string{h}, probeModels, routing, a.base, r.Store.(*recordsStore).Base); err != nil || p != 1 {
		t.Fatalf("the step probe: %d passed, %v", p, err)
	}
	hf.heads, f.posted = nil, nil
	// The world's Start is on its own clock: the sessions run a minute after
	// it, inside the hours of milestone 0 (row 42b).
	started := startOf(t, r)
	r.Now = func() time.Time { return started.Add(time.Minute) }
	return w, r, hf, worldAt{a.base, r.Store.(*recordsStore).Base()}, h
}

// phaseSpec gives the session of the step meaning on the pair that PhasePair
// gives from the harness row, with the problem statement of the records as its input; Keep refuses
// a file that holds "bad", and gives intake/phase.tsv, the session ID and the
// file.
func phaseSpec(t *testing.T, w *world, r *Sessions, a worldAt, harness []string) PhaseSpec {
	t.Helper()
	h, m, err := r.PhasePair(context.Background(), "systems-architect", [][]string{harness}, probeModels, "")
	if err != nil {
		t.Fatalf("PhasePair: %v", err)
	}
	return PhaseSpec{Step: "meaning", Role: "systems-architect", Harness: h, Model: m, Models: probeModels, Base: a.base, Records: a.records,
		Input:  map[string][]byte{"problem-statement.md": []byte(show(t, w.bare, "layup-records", "start/problem-statement.md"))},
		Output: "meaning.tsv",
		Prompt: func(in map[string]string, out string) []byte {
			return []byte("input " + in["problem-statement.md"] + "\noutput " + out + "\n")
		},
		Keep: func(id string, data []byte) (map[string][]byte, error) {
			if bytes.Contains(data, []byte("bad")) {
				return nil, errors.New("a bad row")
			}
			return map[string][]byte{"intake/phase.tsv": append([]byte(id+"\t"), data...)}, nil
		}}
}

// tree gives the paths of layup-records.
func tree(t *testing.T, w *world) []string {
	return strings.Fields(gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "layup-records"))
}

// The demo of row 42a: a session of a phase on a scripted harness reads its
// input from input/ and commits its records and its telemetry row in one
// records commit; it has its own task ID, a start row and no event; then its
// comment; its directory is removed.
func TestASessionOfAPhaseCommitsItsRecordsAndItsTelemetryRow(t *testing.T) {
	seen := filepath.Join(t.TempDir(), "seen")
	w, r, f, a, h := phaseWorld(t, `cp "$in" `+seen+` && printf 'kind\tvalue\nnote\tread\n' > "$out"`)
	before := show(t, w.bare, "layup-records", "sessions.tsv")
	spec := phaseSpec(t, w, r, a, h)
	id, err := r.PhaseSession(context.Background(), spec)
	if err != nil {
		t.Fatalf("PhaseSession: %q, %v", id, err)
	}
	if got, err := os.ReadFile(seen); err != nil || !bytes.Equal(got, spec.Input["problem-statement.md"]) {
		t.Errorf("the input the harness read: %q, %v; want the record byte for byte", got, err)
	}
	if got := show(t, w.bare, "layup-records", "intake/phase.tsv"); got != id+"\tkind\tvalue\nnote\tread\n" {
		t.Errorf("the record kept: %q", got)
	}
	commit := strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%H", "layup-records", "--", "intake/phase.tsv"))
	if changed := strings.Fields(gitOut(t, w.bare, "show", "--format=", "--name-only", commit)); !slices.Equal(changed, []string{"intake/phase.tsv", "telemetry.tsv"}) {
		t.Errorf("the commit of the records changed %q; want the records and the telemetry row in one", changed)
	}
	// Its own task ID: T- and four characters, held by no row before it.
	sessions := show(t, w.bare, "layup-records", "sessions.tsv")
	var task string
	for _, l := range strings.Split(sessions, "\n") {
		if c := strings.Split(l, "\t"); c[0] == id {
			task = c[1]
			if c[2] != "1" || c[3] != "systems-architect" || c[4] != "fake" || c[6] != "m1" || c[7] != a.base {
				t.Errorf("the start row %q", c)
			}
		}
	}
	if !regexp.MustCompile(`^T-[0-9a-z]{4}$`).MatchString(task) || strings.Contains(before, "\t"+task+"\t") {
		t.Errorf("the task ID %q: want T- and four characters, held by no earlier row of sessions.tsv:\n%s", task, before)
	}
	if tel := show(t, w.bare, "layup-records", "telemetry.tsv"); !strings.Contains(tel, "\n"+id+"\t"+task+"\t—\tsystems-architect\tfake\tm1\tsubscription\t") {
		t.Errorf("telemetry.tsv has no row of %s:\n%s", id, tel)
	}
	// No event: no events.tsv but the one of T-ab12.
	for _, p := range tree(t, w) {
		if strings.HasSuffix(p, "events.tsv") && p != "tasks/T-ab12/events.tsv" {
			t.Errorf("a session of a phase wrote %s", p)
		}
	}
	if len(f.posted[2]) != 1 || f.posted[2][0] != id+": systems-architect session of meaning: done\n" {
		t.Errorf("the comments %q", f.posted[2])
	}
	if len(f.heads) != 1 || f.heads[0] != commit {
		t.Errorf("the head of the records at the comment %q; want the records commit %s", f.heads, commit)
	}
	if _, err := os.Stat(filepath.Join(r.Host, "sessions", id)); !os.IsNotExist(err) {
		t.Errorf("the session directory after the call: %v; want it removed", err)
	}
}

// Row 3 of Input states: an output file that is missing, a link, over 1 MiB,
// or that its check refuses: the step's error names the session; the end
// commits the telemetry row alone, and nothing of the file.
func TestAnOutputFileThatIsRefusedIsNoRecord(t *testing.T) {
	for _, c := range []struct{ name, body, comment, reason string }{
		{"missing", `true`, "no-result", "no-result"},
		{"a link", `ln -s /etc/hostname "$out"`, "no-result", "no-result"},
		{"over 1 MiB", `head -c 1048577 /dev/zero | tr '\0' x > "$out"`, "no-result", "no-result"},
		{"refused by its check", `printf 'kind\tvalue\nnote\tbad\n' > "$out"`, "done, refused file", "refused file"},
	} {
		t.Run(c.name, func(t *testing.T) {
			w, r, f, a, h := phaseWorld(t, c.body)
			id, err := r.PhaseSession(context.Background(), phaseSpec(t, w, r, a, h))
			var pe *PhaseError
			if !errors.As(err, &pe) || pe.Session != id || id == "" || !strings.HasPrefix(pe.Reason, c.reason) {
				t.Fatalf("PhaseSession: %q, %v; want a PhaseError of the session, %s", id, err, c.reason)
			}
			if slices.Contains(tree(t, w), "intake/phase.tsv") {
				t.Errorf("the refused file is a record")
			}
			last := strings.Fields(gitOut(t, w.bare, "show", "--format=", "--name-only", "layup-records"))
			if !slices.Equal(last, []string{"telemetry.tsv"}) || !strings.Contains(show(t, w.bare, "layup-records", "telemetry.tsv"), "\n"+id+"\t") {
				t.Errorf("the last records commit changed %q; want the telemetry row of %s alone", last, id)
			}
			if len(f.posted[2]) != 1 || f.posted[2][0] != id+": systems-architect session of meaning: "+c.comment+"\n" {
				t.Errorf("the comments %q", f.posted[2])
			}
		})
	}
}

// Row 1 of Input states: with no admitted pair for the role, no session
// starts: ErrNoPair, and no records commit.
func TestNoAdmittedPairStartsNoSession(t *testing.T) {
	w, r, _, _, h := phaseWorld(t, `true`)
	commits := gitOut(t, w.bare, "rev-list", "--count", "layup-records")
	for _, c := range []struct {
		name   string
		role   string
		models [][]string
	}{
		{"a role with no row", "product-owner", probeModels},
		{"a model of use no", "systems-architect", [][]string{{"fake", "m1", "1000", "https://example.invalid", "2026-10-09T12:00:00Z", "no", ""}}},
	} {
		if _, _, err := r.PhasePair(context.Background(), c.role, [][]string{h}, c.models, ""); !errors.Is(err, route.ErrNoPair) {
			t.Errorf("%s: %v; want ErrNoPair", c.name, err)
		}
	}
	if got := gitOut(t, w.bare, "rev-list", "--count", "layup-records"); got != commits {
		t.Errorf("the records commits: %s, then %s; want none", commits, got)
	}
	// A harness with no row of harnesses.tsv is not admitted.
	_, r2, _, _ := sessionWorld(t)
	if err := r2.commit(context.Background(), map[string][]byte{"routing.tsv": []byte("role\ttier\tposition\tharness\tmodel\nsystems-architect\treasoning\t1\tfake\tm1\n")}, "the routing"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := r2.PhasePair(context.Background(), "systems-architect", [][]string{h}, probeModels, ""); !errors.Is(err, route.ErrNoPair) {
		t.Errorf("a harness never probed: %v; want ErrNoPair", err)
	}
}

// A refused start of a session of a phase (here its version check) writes
// no record and removes its directory; the step's error names the session
// and the refusal.
func TestARefusedStartOfASessionOfAPhaseIsNoRecord(t *testing.T) {
	w, r, f, a, h := phaseWorld(t, `true`)
	spec := phaseSpec(t, w, r, a, h)
	spec.Harness = append([]string{}, h...)
	spec.Harness[regVersion] = "false"
	commits := gitOut(t, w.bare, "rev-list", "--count", "layup-records")
	id, err := r.PhaseSession(context.Background(), spec)
	var pe *PhaseError
	if !errors.As(err, &pe) || pe.Session != id || id == "" || !strings.HasPrefix(pe.Reason, "refused version") {
		t.Fatalf("PhaseSession: %q, %v; want a PhaseError of the session, refused version", id, err)
	}
	if got := gitOut(t, w.bare, "rev-list", "--count", "layup-records"); got != commits {
		t.Errorf("the records commits: %s, then %s; want none", commits, got)
	}
	if len(f.posted[2]) != 0 {
		t.Errorf("the comments %q; want none", f.posted[2])
	}
	if _, err := os.Stat(filepath.Join(r.Host, "sessions", id)); !os.IsNotExist(err) {
		t.Errorf("the session directory after the call: %v; want it removed", err)
	}
}

// A comment that fails after a refused file keeps the session's reason: the
// step's error names the session and the forge (round 1 of #215, note 5).
func TestACommentThatFailsKeepsTheReasonOfTheSession(t *testing.T) {
	w, r, f, a, h := phaseWorld(t, `printf 'kind\tvalue\nnote\tbad\n' > "$out"`)
	f.fail = "Comment"
	id, err := r.PhaseSession(context.Background(), phaseSpec(t, w, r, a, h))
	var pe *PhaseError
	if !errors.As(err, &pe) || pe.Session != id || !strings.HasPrefix(pe.Reason, "refused file") || !strings.Contains(err.Error(), "Comment") {
		t.Errorf("PhaseSession: %q, %v; want the PhaseError of the session, refused file, with the error of the comment", id, err)
	}
}

// The demo of row 42b: no session of Intake starts when the known money since
// Start plus the new session's cap passes the money cap of milestone 0
// (intake.cap 50.0,8.0 of the world), or when the hours since Start pass its
// hours; no records commit and no session directory.
func TestTheCapOfMilestoneZeroStartsNoSession(t *testing.T) {
	for _, c := range []struct {
		name, cap string
		later     time.Duration
		want      string
	}{
		{"the money", "50.01", 0, "the known money 0.0 plus the session's cap 50.01 passes 50.0 USD"},
		{"the hours", "", 8*time.Hour + time.Minute, "the time since Start passes 8.0 hours"},
	} {
		t.Run(c.name, func(t *testing.T) {
			w, r, f, a, h := phaseWorld(t, `printf 'kind\tvalue\nnote\tread\n' > "$out"`)
			h[regCap] = c.cap
			spec := phaseSpec(t, w, r, a, h)
			spec.Harness = h
			started := startOf(t, r)
			r.Now = func() time.Time { return started.Add(c.later) }
			commits := gitOut(t, w.bare, "rev-list", "--count", "layup-records")
			id, err := r.PhaseSession(context.Background(), spec)
			if !errors.Is(err, ledger.ErrMilestoneCap) || id != "" || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("PhaseSession: %q, %v; want ErrMilestoneCap with %q and no session", id, err, c.want)
			}
			if got := gitOut(t, w.bare, "rev-list", "--count", "layup-records"); got != commits {
				t.Errorf("the records commits: %s, then %s; want none", commits, got)
			}
			if left, _ := os.ReadDir(filepath.Join(r.Host, "sessions")); len(left) != 0 || len(f.posted[2]) != 0 {
				t.Errorf("the session directories %v and the comments %q; want none", left, f.posted[2])
			}
		})
	}
	// Inside the cap, the session starts: the money and the hours at the cap
	// exactly.
	w, r, _, a, h := phaseWorld(t, `printf 'kind\tvalue\nnote\tread\n' > "$out"`)
	// The time of Start is the committer time of the first commit of the
	// records branch, not of its head.
	root := strings.TrimSpace(gitOut(t, w.bare, "rev-list", "--max-parents=0", "layup-records"))
	head := strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%ct", "layup-records"))
	if s := strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%ct", root)); s == head || strconv.FormatInt(startOf(t, r).Unix(), 10) != s {
		t.Errorf("the time of Start %v; want the first commit's, %s (the head's is %s)", startOf(t, r), s, head)
	}
	h[regCap] = "50.0"
	started := startOf(t, r)
	r.Now = func() time.Time { return started.Add(96 * time.Hour) }
	spec := phaseSpec(t, w, r, a, h)
	spec.Harness = h
	if id, err := r.PhaseSession(context.Background(), spec); err != nil || id == "" {
		t.Errorf("a session at the cap exactly: %q, %v; want it run", id, err)
	}
}

// startOf gives the time of Start of the world, the first commit of its
// records branch, whose clock is the world's own.
func startOf(t *testing.T, r *Sessions) time.Time {
	t.Helper()
	started, err := r.Store.Started(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return started
}

// filesSpec is phaseSpec for a session whose output is a directory of files,
// as the prose session of the step setup has (docs/spec/scaffold.md, The prose
// rows; task T-eahr): the paths of its two files, and Keep, which refuses a file
// that holds "bad" and gives each at setup/inputs/files/<path>.
func filesSpec(t *testing.T, w *world, r *Sessions, a worldAt, harness []string) PhaseSpec {
	t.Helper()
	s := phaseSpec(t, w, r, a, harness)
	s.Step, s.Output, s.Files, s.Keep = "setup", "files", []string{"README.md", "docs/glossary.md"}, nil
	s.KeepFiles = func(id string, files map[string][]byte) (map[string][]byte, error) {
		kept := map[string][]byte{}
		for p, data := range files {
			if bytes.Contains(data, []byte("bad")) {
				return nil, errors.New("a bad file " + p)
			}
			kept["setup/inputs/files/"+p] = data
		}
		return kept, nil
	}
	return s
}

// A session of a phase whose output is files commits all of them, with the
// telemetry row, in one records commit, byte for byte.
func TestASessionOfAPhaseWithFilesCommitsThemInOneRecordsCommit(t *testing.T) {
	w, r, f, a, h := phaseWorld(t, `mkdir -p "$out/docs" && printf '# Readme\n' > "$out/README.md" && printf '# Glossary\r\n' > "$out/docs/glossary.md"`)
	id, err := r.PhaseSession(context.Background(), filesSpec(t, w, r, a, h))
	if err != nil {
		t.Fatalf("PhaseSession: %q, %v", id, err)
	}
	for path, want := range map[string]string{"setup/inputs/files/README.md": "# Readme\n", "setup/inputs/files/docs/glossary.md": "# Glossary\r\n"} {
		if got := show(t, w.bare, "layup-records", path); got != want {
			t.Errorf("%s: %q; want %q", path, got, want)
		}
	}
	if changed := strings.Fields(gitOut(t, w.bare, "show", "--format=", "--name-only", "layup-records")); !slices.Equal(changed, []string{"setup/inputs/files/README.md", "setup/inputs/files/docs/glossary.md", "telemetry.tsv"}) {
		t.Errorf("the last records commit changed %q; want the two files and the telemetry row in one", changed)
	}
	if len(f.posted[2]) != 1 || f.posted[2][0] != id+": systems-architect session of setup: done\n" {
		t.Errorf("the comments %q", f.posted[2])
	}
	if _, err := os.Stat(filepath.Join(r.Host, "sessions", id)); !os.IsNotExist(err) {
		t.Errorf("the session directory after the call: %v; want it removed", err)
	}
}

// An output of files that is refused is no record: the whole set is read and
// checked first (the cases of the reader are in internal/session), and the end
// commits the telemetry row alone, whichever file is the cause (docs/spec/scaffold.md, Input states: "nothing of it is committed").
func TestAnOutputOfFilesThatIsRefusedIsNoRecord(t *testing.T) {
	good := `mkdir -p "$out/docs" && printf 'g\n' > "$out/docs/glossary.md" && `
	for _, c := range []struct{ name, body, comment string }{
		{"a file that no name covers", good + `printf 'r\n' > "$out/README.md" && printf 'x\n' > "$out/extra.md"`, "no-result"},
		{"not valid UTF-8", good + `printf 'a\377b' > "$out/README.md"`, "no-result"},
		{"refused by its check", good + `printf 'a bad one\n' > "$out/README.md"`, "done, refused file"},
	} {
		t.Run(c.name, func(t *testing.T) {
			w, r, f, a, h := phaseWorld(t, c.body)
			id, err := r.PhaseSession(context.Background(), filesSpec(t, w, r, a, h))
			var pe *PhaseError
			if !errors.As(err, &pe) || pe.Session != id || id == "" {
				t.Fatalf("PhaseSession: %q, %v; want a PhaseError of the session", id, err)
			}
			for _, p := range tree(t, w) {
				if strings.HasPrefix(p, "setup/inputs/files/") {
					t.Errorf("the refused output left %s in the records", p)
				}
			}
			last := strings.Fields(gitOut(t, w.bare, "show", "--format=", "--name-only", "layup-records"))
			if !slices.Equal(last, []string{"telemetry.tsv"}) {
				t.Errorf("the last records commit changed %q; want the telemetry row alone", last)
			}
			if len(f.posted[2]) != 1 || f.posted[2][0] != id+": systems-architect session of setup: "+c.comment+"\n" {
				t.Errorf("the comments %q", f.posted[2])
			}
		})
	}
}

// The session of a phase runs on the clone and the branch of its Sessions, with
// repo/ on its task branch at Base: the prose session of the step setup gets
// the work area's target at the head of layup-setup (task T-eahr).
func TestASessionOfAPhaseRunsOnTheCloneAndTheBranchOfItsSessions(t *testing.T) {
	seen := filepath.Join(t.TempDir(), "seen")
	w, r, _, a, h := phaseWorld(t, `{ git rev-parse --abbrev-ref HEAD; git rev-parse HEAD; cat docs/setup/armature.pin; git branch --list; } > `+seen+` && mkdir -p "$out" && printf 'x\n' > "$out/README.md"`)
	other := t.TempDir()
	gitOut(t, other, "init", "-q", "-b", "main")
	os.WriteFile(filepath.Join(other, "README.md"), []byte("# Root\n"), 0o644)
	gitOut(t, other, "add", "-A")
	gitOut(t, other, "commit", "-q", "-m", "root")
	gitOut(t, other, "switch", "-q", "-c", "layup-setup")
	os.MkdirAll(filepath.Join(other, "docs", "setup"), 0o755)
	os.WriteFile(filepath.Join(other, "docs", "setup", "armature.pin"), []byte("the pin\n"), 0o644)
	gitOut(t, other, "add", "-A")
	gitOut(t, other, "commit", "-q", "-m", "chore: setup S04")
	head := strings.TrimSpace(gitOut(t, other, "rev-parse", "HEAD"))
	r.Clone, r.Branch = other, "layup-setup"
	spec := filesSpec(t, w, r, a, h)
	spec.Files, spec.Base = []string{"README.md"}, head
	if id, err := r.PhaseSession(context.Background(), spec); err != nil {
		t.Fatalf("PhaseSession: %q, %v", id, err)
	}
	got, err := os.ReadFile(seen)
	lines := strings.Split(strings.TrimRight(string(got), "\n"), "\n")
	if err != nil || len(lines) < 4 || !regexp.MustCompile(`^task/T-[0-9a-z]{4}/1$`).MatchString(lines[0]) || lines[1] != head || lines[2] != "the pin" {
		t.Errorf("what the session saw in repo/: %q, %v; want its task branch at %s, with docs/setup/armature.pin", got, err, head)
	}
	if branches := strings.Join(lines[3:], " "); strings.Contains(branches, "main") {
		t.Errorf("the branches of repo/: %q; want the one branch of the clone and the task branch", branches)
	}
}
