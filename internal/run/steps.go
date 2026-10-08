package run

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/tsv"
)

// RunStepsSchema is the form of the table that layup run prints: the block
// run-steps of docs/spec/run.md (The command).
var RunStepsSchema = tsv.Schema{Name: "run-steps", Location: "stdout", Columns: []tsv.Column{
	{Name: "step", Type: "enum(forge|plan|baseline|root-push|read-back|records|issues|watch|clone|version|lease|phase)", Key: true},
	{Name: "result", Type: "enum(done|fail)"},
	{Name: "detail", Type: "text"},
}}

// Step is a row of run-steps.
type Step struct{ Name, Result, Detail string }

// WriteSteps writes the table run-steps.
func WriteSteps(w io.Writer, rows []Step) error {
	t := make([][]string, len(rows))
	for i, r := range rows {
		t[i] = []string{r.Name, r.Result, r.Detail}
	}
	return tsv.Write(w, RunStepsSchema, t)
}

// Config is what internal/cli gives layup run (docs/spec/run.md, The command):
// the target OWNER/NAME, the host directory DIR, the values of the flags, the
// two host registers as read, the forge adapter, the run's clock and progress
// function, its run ID, its host name, LAYUP's version and LAYUP's own pin.
type Config struct {
	Owner, Name         string
	Dir                 string
	PSB, Vision         []byte // Vision is nil with no --vision
	Operator, IdeaOwner string
	Plan, IntakeCap     string
	LeaseH, WatchT      int // minutes
	Register            forge.Register
	Harnesses           [][]string // the rows of harnesses.tsv
	Forge               forge.Forge
	Clock               Clock
	Progress            func(string)
	RunID, HostName     string
	Version             string
	Baseline            string // the source of LAYUP's own pin
	LayupPin            string // the commit of LAYUP's own pin
}

// state is one run.
type state struct {
	cfg                            Config
	permissions                    string
	operatorID, ideaOwnerID, botID int64
	visibility, defaultBranch      string
	pinCommit, pinTree             string
	pinTime                        time.Time
	root, target                   string
	store                          *recordsStore
	start, approvers, copies       [][]string
	lease                          LeaseRow
	stopBeat                       func() error       // ends the heartbeat and gives its error
	cancel                         context.CancelFunc // ends the context of the steps
	mu                             sync.Mutex
	lostBeat                       error // a beat that was lost
}

// Start runs the steps of layup run --new (docs/spec/run.md, The steps of
// layup run --new) and gives the rows of run-steps; it stops at the first fail.
func Start(ctx context.Context, cfg Config) []Step {
	r := &state{cfg: cfg}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return r.run(ctx, cancel, []func(context.Context) Step{
		func(ctx context.Context) Step { return r.forgeStep(ctx, true) },
		func(context.Context) Step { return r.planStep() },
		r.baselineStep, r.rootPushStep, r.readBackStep, r.recordsStep,
		func(ctx context.Context) Step { return r.issuesStep(ctx, "issues") },
		func(ctx context.Context) Step { return r.watchStep(ctx, "watch") },
		r.leaseStep,
	})
}

// Restart runs the steps of layup run TARGET (docs/spec/run.md, The restart).
func Restart(ctx context.Context, cfg Config) []Step {
	r := &state{cfg: cfg}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return r.run(ctx, cancel, []func(context.Context) Step{
		func(ctx context.Context) Step { return r.forgeStep(ctx, false) },
		r.cloneStep, r.versionStep, r.takeStep, r.phaseStep,
	})
}

// run runs the steps in order, stops at the first fail, and ends the heartbeat
// before it returns. A lost beat ends the context of the step that runs, and
// that step fails with the loss.
func (r *state) run(ctx context.Context, cancel context.CancelFunc, steps []func(context.Context) Step) []Step {
	var rows []Step
	r.cancel = cancel
	defer r.endBeat()
	for _, step := range steps {
		s := step(ctx)
		if lost := r.lost(); s.Result == "fail" && lost != nil {
			s.Detail = lost.Error()
		}
		s.Detail = r.clean(s.Detail)
		rows = append(rows, s)
		if s.Result == "fail" {
			break
		}
	}
	return rows
}

func (r *state) lost() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lostBeat
}

func done(name, detail string) Step { return Step{name, "done", detail} }

func fail(name string, err error) Step {
	return Step{name, "fail", strings.SplitN(err.Error(), "\n", 2)[0]}
}

func failf(name, format string, a ...any) Step { return Step{name, "fail", fmt.Sprintf(format, a...)} }

// clean takes the paths of the host out of a detail (run-steps: "no path of
// the host").
func (r *state) clean(s string) string {
	if r.cfg.Dir != "" {
		s = strings.ReplaceAll(s, r.cfg.Dir, "DIR")
	}
	return strings.ReplaceAll(s, strings.TrimSuffix(os.TempDir(), "/"), "TMPDIR")
}

func (r *state) url() string {
	return r.cfg.Register.Web + "/" + r.cfg.Owner + "/" + r.cfg.Name + ".git"
}

// auth gives git the token of the forge: one that has five minutes or more.
func (r *state) auth(ctx context.Context) (git.Auth, error) {
	t, err := r.cfg.Forge.Token(ctx)
	if err != nil {
		return git.Auth{}, err
	}
	return git.Auth{Web: r.cfg.Register.Web, Token: t.Token}, nil
}

// who is the App's bot at a time (run.md, step 6).
func (r *state) who(at time.Time) git.Identity {
	bot := r.cfg.Register.AppSlug + "[bot]"
	return git.Identity{Name: bot, Email: fmt.Sprintf("%d+%s@users.noreply.github.com", r.botID, bot), Time: at}
}

func (r *state) newStore(read string) {
	r.store = &recordsStore{dir: r.target, url: r.url(), auth: r.auth, read: read,
		who: func() git.Identity { return r.who(r.cfg.Clock.Now()) }}
}

// forgeStep is step 1 of Start (isNew) and of the restart.
func (r *state) forgeStep(ctx context.Context, isNew bool) Step {
	const name = "forge"
	f := r.cfg.Forge
	if _, err := f.Token(ctx); err != nil {
		return fail(name, err)
	}
	if m := forge.Missing(f.Capabilities()); len(m) > 0 {
		return failf(name, "the adapter does not declare the capability %s", m[0])
	}
	inst, err := f.Installation(ctx)
	if err != nil {
		return fail(name, err)
	}
	if m := forge.CheckPermissions(inst.Permissions); len(m) > 0 {
		return failf(name, "the installation lacks the permissions %s", strings.Join(m, ", "))
	}
	var pairs []string
	for n, level := range inst.Permissions {
		pairs = append(pairs, n+":"+level)
	}
	sort.Strings(pairs)
	r.permissions = strings.Join(pairs, " ")
	repo, err := f.Repository(ctx)
	if err != nil {
		return fail(name, err)
	}
	r.visibility, r.defaultBranch = repo.Visibility, repo.DefaultBranch
	records := slices.Contains(repo.Branches, "layup-records")
	switch {
	case isNew && records:
		return failf(name, "the repository has the branch layup-records; --new takes an empty repository (O-163)")
	case isNew && repo.HasCommit:
		return failf(name, "the repository has a commit; --new takes an empty repository (O-163)")
	case !isNew && !records:
		return failf(name, "the repository has no branch layup-records")
	}
	if isNew {
		for _, p := range []struct {
			login string
			id    *int64
		}{{r.cfg.Operator, &r.operatorID}, {r.cfg.IdeaOwner, &r.ideaOwnerID}} {
			if *p.id, err = f.UserID(ctx, p.login); err != nil {
				return failf(name, "the login %s: %s", p.login, strings.SplitN(err.Error(), "\n", 2)[0])
			}
		}
	}
	if r.botID, err = f.UserID(ctx, r.cfg.Register.AppSlug+"[bot]"); err != nil {
		return failf(name, "the App's bot %s[bot]: %s", r.cfg.Register.AppSlug, strings.SplitN(err.Error(), "\n", 2)[0])
	}
	state := "empty"
	if !isNew {
		state = "with layup-records"
	}
	return done(name, fmt.Sprintf("the installation token, the six capabilities and the permissions of M2a; a %s repository, %s", r.visibility, state))
}

// planStep is step 2: the plan check of §5 (K15).
func (r *state) planStep() Step {
	if r.visibility == "public" || r.cfg.Plan == "team" || r.cfg.Plan == "enterprise" {
		return done("plan", fmt.Sprintf("the plan %s on a %s repository", r.cfg.Plan, r.visibility))
	}
	return failf("plan", "the plan %s has no rulesets or draft pull requests on a %s repository", r.cfg.Plan, r.visibility)
}

// baselineStep is step 3 (as S02): the latest commit of the baseline, its
// tree, and a copy with no .git. A Start before step 6 holds no records, so it
// removes the directories of an earlier Start first.
func (r *state) baselineStep(context.Context) Step {
	const name = "baseline"
	r.root = filepath.Join(r.cfg.Dir, "roots", r.cfg.Owner, r.cfg.Name)
	r.target = filepath.Join(r.cfg.Dir, "targets", r.cfg.Owner, r.cfg.Name)
	for _, d := range []string{r.root, r.target} {
		if err := os.RemoveAll(d); err != nil {
			return fail(name, err)
		}
	}
	if err := os.MkdirAll(filepath.Dir(r.root), 0o755); err != nil {
		return fail(name, err)
	}
	commit, err := git.LsRemote(r.cfg.Baseline, "HEAD")
	if err != nil {
		return fail(name, err)
	}
	r.pinTime = r.cfg.Clock.Now().UTC().Truncate(time.Second)
	if err := git.Clone(r.cfg.Baseline, r.root, git.Auth{}); err != nil {
		return fail(name, err)
	}
	if err := git.CheckoutDetach(r.root, commit); err != nil {
		return fail(name, err)
	}
	if r.pinTree, err = git.RevParse(r.root, commit+"^{tree}"); err != nil {
		return fail(name, err)
	}
	if err := os.RemoveAll(filepath.Join(r.root, ".git")); err != nil {
		return fail(name, err)
	}
	r.pinCommit = commit
	return done(name, "the baseline at "+commit)
}

// readEvery is the period of each wait of the steps.
const waitEvery = 10 * time.Second

// rootPushStep is step 4 (as S03): the root commit as the App's bot at
// pin.time, the command that the Operator runs, and the wait for its push.
func (r *state) rootPushStep(ctx context.Context) Step {
	const name = "root-push"
	if err := git.Init(r.root); err != nil {
		return fail(name, err)
	}
	if err := git.Add(r.root); err != nil {
		return fail(name, err)
	}
	if err := git.Commit(r.root, "chore: the unmodified baseline at "+r.pinCommit, r.who(r.pinTime)); err != nil {
		return fail(name, err)
	}
	if tree, err := git.RevParse(r.root, "HEAD^{tree}"); err != nil || tree != r.pinTree {
		return failf(name, "the root tree %s is not pin.tree %s", tree, r.pinTree)
	}
	pin := "not the commit of LAYUP's own pin " + short(r.cfg.LayupPin)
	if r.pinCommit == r.cfg.LayupPin {
		pin = "the commit of LAYUP's own pin " + short(r.cfg.LayupPin)
	}
	r.cfg.Progress(fmt.Sprintf("root-push: the root commit is the unmodified baseline at %s, %s; push it with your own login: git -C %s push -- %s %s",
		r.pinCommit, pin, r.root, r.url(), "main"))
	for {
		repo, err := r.cfg.Forge.Repository(ctx)
		if err != nil {
			return fail(name, err)
		}
		if slices.Contains(repo.Branches, r.defaultBranch) {
			return done(name, "the default branch "+r.defaultBranch+" exists")
		}
		r.cfg.Progress("root-push: waiting for the push of the root commit")
		if err := r.cfg.Clock.Sleep(ctx, waitEvery); err != nil {
			return fail(name, err)
		}
	}
}

func short(id string) string {
	if len(id) > 7 {
		return id[:7]
	}
	return id
}

// readBackStep is step 5: the clone made with Init and Fetch, and the checks
// of the default branch.
func (r *state) readBackStep(ctx context.Context) Step {
	const name = "read-back"
	a, err := r.auth(ctx)
	if err != nil {
		return fail(name, err)
	}
	ref := "refs/heads/" + r.defaultBranch
	if err := git.Init(r.target); err != nil {
		return fail(name, err)
	}
	if err := git.Fetch(r.target, r.url(), ref, a); err != nil {
		return fail(name, err)
	}
	repo, err := r.cfg.Forge.Repository(ctx)
	if err != nil {
		return fail(name, err)
	}
	if !slices.Equal(repo.Branches, []string{r.defaultBranch}) {
		return failf(name, "the repository has the branches %s; want %s only", strings.Join(repo.Branches, ", "), r.defaultBranch)
	}
	if repo.Visibility != r.visibility {
		return failf(name, "the visibility is %s; step 1 read %s", repo.Visibility, r.visibility)
	}
	head, err := git.RevParse(r.target, ref)
	if err != nil {
		return fail(name, err)
	}
	roots, err := git.RootCommits(r.target, ref)
	if err != nil {
		return fail(name, err)
	}
	tree, err := git.RevParse(r.target, head+"^{tree}")
	if err != nil {
		return fail(name, err)
	}
	if !slices.Equal(roots, []string{head}) || tree != r.pinTree {
		return failf(name, "the default branch is not one root commit with the tree pin.tree %s", r.pinTree)
	}
	return done(name, "one branch, one root commit, root tree = pin.tree")
}

func sha(data []byte) string { s := sha256.Sum256(data); return hex.EncodeToString(s[:]) }

// startRows gives the rows of start.tsv of step 6, in the order of the block.
func (r *state) startRows() [][]string {
	vision := "" // the empty value, which start.tsv writes as —
	if r.cfg.Vision != nil {
		vision = sha(r.cfg.Vision)
	}
	rows := [][]string{
		{"layup.version", r.cfg.Version, "run"},
		{"psb.sha256", sha(r.cfg.PSB), "command"},
		{"vision.sha256", vision, "command"},
		{"forge.plan", r.cfg.Plan, "command"},
		{"forge.visibility", r.visibility, "forge"},
		{"app.permissions", r.permissions, "forge"},
		{"operator.id", strconv.FormatInt(r.operatorID, 10), "forge"},
		{"idea-owner.id", strconv.FormatInt(r.ideaOwnerID, 10), "forge"},
		{"intake.cap", r.cfg.IntakeCap, "command"},
		{"lease.H", strconv.Itoa(r.cfg.LeaseH), "command"},
		{"watch.T", strconv.Itoa(r.cfg.WatchT), "command"},
	}
	for _, h := range r.cfg.Harnesses {
		rows = append(rows, []string{"harness." + h[0] + ".cap", h[1], "register"}, []string{"harness." + h[0] + ".wall", h[2], "register"})
	}
	return append(rows,
		[]string{"pin.source", r.cfg.Baseline, "run"},
		[]string{"pin.commit", r.pinCommit, "run"},
		[]string{"pin.tree", r.pinTree, "run"},
		[]string{"pin.time", stamp(r.pinTime), "run"},
		[]string{"issue.intake", "", "forge"},
		[]string{"issue.control", "", "forge"},
		[]string{"watch", "", "run"})
}

func (r *state) value(name string) string {
	for _, row := range r.start {
		if row[0] == name {
			return row[1]
		}
	}
	return ""
}

func (r *state) set(name, value string) {
	for _, row := range r.start {
		if row[0] == name {
			row[1] = value
		}
	}
}

func (r *state) harnessIDs() []string {
	var ids []string
	for _, h := range r.cfg.Harnesses {
		ids = append(ids, h[0])
	}
	return ids
}

// writeStart commits start.tsv as it is now.
func (r *state) writeStart(ctx context.Context, message string) error {
	if err := records.CheckStart(r.start, r.harnessIDs()); err != nil {
		return err
	}
	data, err := table(records.StartSchema, r.start)
	if err != nil {
		return err
	}
	return r.write(ctx, map[string][]byte{"start/start.tsv": data}, message)
}

// write makes a records commit through fencing.
func (r *state) write(ctx context.Context, files map[string][]byte, message string) error {
	return Fenced(ctx, r.store, r.cfg.RunID, func(ctx context.Context) error { return r.store.commit(ctx, files, message, "") })
}

// recordsStep is step 6: the first records commit, an orphan commit as S15
// makes it, and its push; the heartbeat starts.
func (r *state) recordsStep(ctx context.Context) Step {
	const name = "records"
	now := r.cfg.Clock.Now()
	r.start = r.startRows()
	if err := records.CheckStart(r.start, r.harnessIDs()); err != nil {
		return fail(name, err)
	}
	r.approvers = [][]string{
		{strconv.FormatInt(r.operatorID, 10), "operator", r.cfg.Operator, stamp(now), "start"},
		{strconv.FormatInt(r.ideaOwnerID, 10), "idea-owner", r.cfg.IdeaOwner, stamp(now), "start"},
	}
	r.lease = LeaseRow{Run: r.cfg.RunID, Host: r.cfg.HostName, Version: r.cfg.Version, Started: now.UTC().Truncate(time.Second), State: "held"}
	files := map[string][]byte{"README.md": []byte(startReadme), "start/problem-statement.md": r.cfg.PSB}
	if r.cfg.Vision != nil {
		files["start/vision.md"] = r.cfg.Vision
	}
	for p, t := range map[string]struct {
		s    tsv.Schema
		rows [][]string
	}{"start/start.tsv": {records.StartSchema, r.start}, "approvers.tsv": {records.ApproversSchema, r.approvers}} {
		data, err := table(t.s, t.rows)
		if err != nil {
			return fail(name, err)
		}
		files[p] = data
	}
	data, err := leaseTable(r.lease)
	if err != nil {
		return fail(name, err)
	}
	files["lease.tsv"] = data
	r.newStore("")
	if err := r.store.commit(ctx, files, "records: the Start of "+r.cfg.Owner+"/"+r.cfg.Name, r.defaultBranch); err != nil {
		return fail(name, err)
	}
	r.beat()
	return done(name, "the first records commit: the pin, the briefs, approvers.tsv, start.tsv, the lease")
}

// beat starts the heartbeat of the lease: a goroutine with its own context; a
// lost beat ends the run's context.
func (r *state) beat() {
	hb, stop := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		err := Heartbeat(hb, r.store, r.cfg.Clock, &r.lease, time.Duration(r.cfg.LeaseH)*time.Minute)
		if err != nil {
			r.mu.Lock()
			r.lostBeat = err
			r.mu.Unlock()
			r.cancel()
		}
		result <- err
	}()
	r.stopBeat = func() error {
		stop()
		r.stopBeat = nil
		return <-result
	}
}

func (r *state) endBeat() {
	if r.stopBeat != nil {
		r.stopBeat()
	}
}

// issuesStep is step 7: for each issue whose row is not a number, a records
// commit with opening, the issue, and a records commit with its number.
func (r *state) issuesStep(ctx context.Context, name string) Step {
	for _, i := range []struct{ row, title, body string }{
		{"issue.intake", intakeTitle, intakeBody}, {"issue.control", controlTitle, controlBody},
	} {
		if _, err := strconv.Atoi(r.value(i.row)); err == nil {
			continue
		}
		if r.value(i.row) != "opening" { // a restart finds the announcement made
			r.set(i.row, "opening")
			if err := r.writeStart(ctx, "records: "+i.row+" opening"); err != nil {
				return fail(name, err)
			}
		}
		n, err := r.cfg.Forge.OpenIssue(ctx, i.title, i.body)
		if err != nil {
			return fail(name, err)
		}
		r.set(i.row, strconv.Itoa(n))
		if err := r.writeStart(ctx, fmt.Sprintf("records: %s %d", i.row, n)); err != nil {
			return fail(name, err)
		}
	}
	return done(name, "the Intake issue "+r.value("issue.intake")+", the control issue "+r.value("issue.control"))
}

// watchStep is step 8: the comments of the control issue every ten seconds for
// up to watch.T, each copied and pushed before it is read, until a notice of
// the App watch_slug comes.
func (r *state) watchStep(ctx context.Context, name string) Step {
	confirmed := false
	if slug := r.cfg.Register.WatchSlug; slug != "" {
		issue, err := strconv.Atoi(r.value("issue.control"))
		if err != nil {
			return failf(name, "no control issue")
		}
		end := r.cfg.Clock.Now().Add(time.Duration(r.cfg.WatchT) * time.Minute)
		for !confirmed {
			comments, err := r.cfg.Forge.Comments(ctx, issue)
			if err != nil {
				return fail(name, err)
			}
			for _, c := range comments {
				if err := r.copy(ctx, c, issue); err != nil {
					return fail(name, err)
				}
				confirmed = confirmed || c.App == slug
			}
			if confirmed || !r.cfg.Clock.Now().Before(end) {
				break
			}
			r.cfg.Progress("watch: waiting for the notice of " + slug)
			if err := r.cfg.Clock.Sleep(ctx, waitEvery); err != nil {
				return fail(name, err)
			}
		}
	}
	result := "not-confirmed"
	if confirmed {
		result = "confirmed"
	}
	r.set("watch", result)
	if err := r.writeStart(ctx, "records: watch "+result); err != nil {
		return fail(name, err)
	}
	return done(name, result)
}

// copy pushes the copy of a comment before any step reads it (copy before
// read).
func (r *state) copy(ctx context.Context, c forge.Comment, issue int) error {
	row, body, ok := Copy(r.copies, c, issue, r.cfg.Clock.Now())
	if !ok {
		return nil
	}
	data, err := table(records.CopiesSchema, append(slices.Clone(r.copies), row))
	if err != nil {
		return err
	}
	if err := r.write(ctx, map[string][]byte{"copies.tsv": data, row[9]: body}, "records: the copy of the comment "+row[0]); err != nil {
		return err
	}
	r.copies = append(r.copies, row)
	return nil
}

// leaseStep is step 9: the heartbeat ends, and the lease is released.
func (r *state) leaseStep(ctx context.Context) Step {
	if r.stopBeat != nil {
		if err := r.stopBeat(); err != nil {
			return fail("lease", err)
		}
	}
	if err := Release(ctx, r.store, &r.lease); err != nil {
		return fail("lease", err)
	}
	return done("lease", "released")
}

// cloneStep is step 2 of the restart: a fresh clone from the forge, and the
// records of the Start read.
func (r *state) cloneStep(ctx context.Context) Step {
	const name = "clone"
	r.target = filepath.Join(r.cfg.Dir, "targets", r.cfg.Owner, r.cfg.Name)
	if err := os.RemoveAll(r.target); err != nil {
		return fail(name, err)
	}
	if err := os.MkdirAll(filepath.Dir(r.target), 0o755); err != nil {
		return fail(name, err)
	}
	a, err := r.auth(ctx)
	if err != nil {
		return fail(name, err)
	}
	if err := git.Clone(r.url(), r.target, a); err != nil {
		return fail(name, err)
	}
	const rev = "refs/remotes/origin/layup-records"
	head, err := git.RevParse(r.target, rev)
	if err != nil {
		return fail(name, err)
	}
	read := func(path string, reader func([]byte) ([][]string, error)) ([][]string, error) {
		data, err := git.Show(r.target, head, path)
		if err != nil {
			return nil, err
		}
		rows, err := reader(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		return rows, nil
	}
	if r.start, err = read("start/start.tsv", func(d []byte) ([][]string, error) { return records.ReadStart(d, r.harnessIDs()) }); err != nil {
		return fail(name, err)
	}
	if r.approvers, err = read("approvers.tsv", records.ReadApprovers); err != nil {
		return fail(name, err)
	}
	if _, err := read("lease.tsv", records.ReadLease); err != nil {
		return fail(name, err)
	}
	if _, err := git.Show(r.target, head, "copies.tsv"); err == nil {
		if r.copies, err = read("copies.tsv", records.ReadCopies); err != nil {
			return fail(name, err)
		}
	}
	// A restart has no --lease-h and no --watch-t: the Start values are in
	// start.tsv.
	if r.cfg.LeaseH, err = strconv.Atoi(r.value("lease.H")); err != nil {
		return fail(name, err)
	}
	if r.cfg.WatchT, err = strconv.Atoi(r.value("watch.T")); err != nil {
		return fail(name, err)
	}
	r.newStore(head)
	return done(name, "start.tsv, approvers.tsv and lease.tsv read")
}

// versionStep is step 3 of the restart.
func (r *state) versionStep(context.Context) Step {
	if v := r.value("layup.version"); v != r.cfg.Version {
		return failf("version", "the target was started by LAYUP %s; this is LAYUP %s", v, r.cfg.Version)
	}
	return done("version", "LAYUP "+r.cfg.Version)
}

// takeStep is step 4 of the restart: the lease, at once when it is released,
// else by the takeover rule; the heartbeat starts.
func (r *state) takeStep(ctx context.Context) Step {
	me := LeaseRow{Run: r.cfg.RunID, Host: r.cfg.HostName, Version: r.cfg.Version}
	var err error
	if r.lease, err = Take(ctx, r.store, r.cfg.Clock, me, time.Duration(r.cfg.LeaseH)*time.Minute, r.cfg.Progress); err != nil {
		var held *HeldError
		var lost *LostError
		if errors.As(err, &held) || errors.As(err, &lost) {
			return failf("lease", "%v", err)
		}
		return fail("lease", err)
	}
	r.beat()
	return done("lease", "held by this run")
}

// phaseStep is step 5 of the restart: the steps of Start that are not done,
// run again; then, in M2a, the lease is released.
func (r *state) phaseStep(ctx context.Context) Step {
	const name = "phase"
	var again []string
	if _, err := strconv.Atoi(r.value("issue.intake")); err != nil {
		again = append(again, "issues")
	} else if _, err := strconv.Atoi(r.value("issue.control")); err != nil {
		again = append(again, "issues")
	}
	if r.value("watch") == "" {
		again = append(again, "watch")
	}
	if slices.Contains(again, "issues") {
		if s := r.issuesStep(ctx, name); s.Result == "fail" {
			return s
		}
	}
	if slices.Contains(again, "watch") {
		if s := r.watchStep(ctx, name); s.Result == "fail" {
			return s
		}
	}
	if s := r.leaseStep(ctx); s.Result == "fail" {
		return Step{name, "fail", s.Detail}
	}
	if len(again) == 0 {
		return done(name, "each step of Start is done; the lease is released")
	}
	return done(name, "ran again: "+strings.Join(again, ", ")+"; the lease is released")
}
