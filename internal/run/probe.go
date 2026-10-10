package run

// This file holds the probe of a harness and the step probe of the restart
// (docs/spec/session.md, The probe, The routing register, A comment for a
// session; run.md, The restart; row 38 of the plan, task T-nxe4, #168).

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/session"
	"github.com/pharzam/layup/internal/tsv"
)

// probePrompt is the fixed text of a probe's prompt, the block probe-prompt
// of docs/spec/session.md, with {result} and {token}.
const probePrompt = `This is a probe of LAYUP. Change no file in this repository and make no commit.

Write one file, {result}, of tab-separated values, with the header line
"kind<TAB>value" and these rows: one row "token<TAB>{token}"; and one row
"file<TAB><path>" for each instruction file that you loaded for this session,
with its path relative to this repository, or absolute when it is outside it.
Then end.
`

// errProbed is the end of a probe session whose version's last probe
// passed: it takes the version and runs no probe, and writes no record.
var errProbed = errors.New("the last probe at this version passed")

// The columns of a row of the harness register (route.ReadHarnesses).
const (
	regHarness = iota
	regCap
	regWall
	regCommand
	regPrompt
	regVersion
	regCredential
	regCredentialTo
	regRules
	regPolicy
	regUsage
	regBilling
	regVars
)

// pairOf gives the Pair of a harness register row and a model row
// (route.ReadModels: harness, model, context, …).
func pairOf(h, model []string) Pair {
	wall, context := 0, 0
	fmt.Sscan(h[regWall], &wall)
	fmt.Sscan(model[2], &context)
	cred, credTo := h[regCredential], h[regCredentialTo]
	if credTo == "" {
		credTo = "—"
	}
	return Pair{Harness: h[regHarness], Model: model[1], VersionCommand: h[regVersion], Command: h[regCommand], PromptMode: h[regPrompt],
		Cap: h[regCap], Wall: wall, Rules: strings.Fields(h[regRules]), Policy: strings.Fields(h[regPolicy]), Context: context,
		Session: session.Harness{Credential: cred, CredentialTo: credTo, Vars: strings.Fields(h[regVars])},
		Billing: h[regBilling], Usage: h[regUsage]}
}

// firstModel gives the first model row of harness whose use is yes, in the
// order of models.tsv, or nil for none.
func firstModel(models [][]string, harness string) []string {
	for _, m := range models {
		if m[0] == harness && m[5] == "yes" {
			return m
		}
	}
	return nil
}

// probeReason gives why a probe failed, or "" when it passed
// (docs/spec/session.md, The probe): the class of its end other than done;
// token, a token that is not the prompt's; files, no file row AGENTS.md;
// outside, a file row outside the session directory that is no policy path
// of the row (a relative value is read from repo/); not-used, a model of the
// usage report whose row of models.tsv has use no.
func probeReason(class string, resultErr error, rows [][]string, token string, d session.Dir, policy, used []string, models [][]string, harness string) string {
	if class != "done" {
		return class
	}
	tokens, files := 0, []string{}
	for _, r := range rows {
		if r[0] == "token" {
			tokens++
			if r[1] != token {
				return "token"
			}
		} else {
			files = append(files, r[1])
		}
	}
	if tokens != 1 {
		return "token"
	}
	if !slices.Contains(files, "AGENTS.md") {
		return "files"
	}
	for _, f := range files {
		path := f
		if !filepath.IsAbs(path) {
			path = filepath.Join(d.Repo, path)
		}
		path = filepath.Clean(path)
		if path != d.Root && !strings.HasPrefix(path, d.Root+string(filepath.Separator)) && !slices.Contains(policy, f) {
			return "outside"
		}
	}
	for _, m := range used {
		for _, row := range models {
			if row[0] == harness && row[1] == m && row[5] == "no" {
				return "not-used"
			}
		}
	}
	return ""
}

// routingDiffers reports whether two routing tables hold other rows, in any
// order (the block's rows have no order of their own).
func routingDiffers(a, b [][]string) bool {
	key := func(t [][]string) []string {
		out := make([]string, len(t))
		for i, r := range t {
			out[i] = strings.Join(r, "\t")
		}
		slices.Sort(out)
		return out
	}
	return !slices.Equal(key(a), key(b))
}

// ProbeSpec is one probe: the harness register row and the rows of
// models.tsv, the base commit (the head of the default branch), the records
// commit that the prompt was built from, and whether it sweeps first (the
// step does; a probe inside a task session's version check does not, as the
// task session swept and its own directory is live).
type ProbeSpec struct {
	Harness []string
	Models  [][]string
	Base    string
	Records string
	Sweep   bool
}

// probeTaskID draws the probe's own task ID, T- and four random characters,
// again until no row of sessions.tsv holds it.
func (s *Sessions) probeTaskID(ctx context.Context) (string, error) {
	rows, err := s.readTable(ctx, "sessions.tsv", records.ReadSessions)
	if err != nil {
		return "", err
	}
	const chars = "0123456789abcdefghijklmnopqrstuvwxyz"
	for {
		b := make([]byte, 4)
		if _, err := rand.Read(b); err != nil {
			return "", err
		}
		id := "T-"
		for _, c := range b {
			id += string(chars[int(c)%len(chars)])
		}
		if !slices.ContainsFunc(rows, func(r []string) bool { return r[sessionTask] == id }) {
			return id, nil
		}
	}
}

// Probe runs the probe of one harness and gives its result: passed, failed,
// or "" when the last probe at its version passed (no record). A harness with
// no model of use yes is the caller's to skip.
func (s *Sessions) Probe(ctx context.Context, p ProbeSpec) (string, error) {
	s.fns()
	model := firstModel(p.Models, p.Harness[regHarness])
	if model == nil {
		return "", fmt.Errorf("the harness %s has no model of use yes", p.Harness[regHarness])
	}
	task, err := s.probeTaskID(ctx)
	if err != nil {
		return "", err
	}
	tok := make([]byte, 8)
	if _, err := rand.Read(tok); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tok)
	pair := pairOf(p.Harness, model)
	spec := TaskSpec{Task: task, Attempt: 1, Role: "probe", Base: p.Base, Records: p.Records, Pair: pair,
		probe: &probeOf{token: token, models: p.Models, sweep: p.Sweep}}
	spec.Admit = func(version string) error {
		rows, err := s.readTable(ctx, "harnesses.tsv", records.ReadHarnesses)
		if err != nil {
			return err
		}
		if route.Probed(rows, pair.Harness, version) {
			return errProbed
		}
		return nil
	}
	spec.End = func(id string, d session.Dir, r session.Run, err error) error {
		return s.endProbe(ctx, spec, id, d, r, err)
	}
	id, err := s.TaskSession(ctx, spec)
	if errors.Is(err, errProbed) {
		return "", nil
	}
	rows, rErr := s.readTable(ctx, "harnesses.tsv", records.ReadHarnesses)
	if rErr != nil {
		return "", errors.Join(err, rErr)
	}
	for _, r := range rows {
		if r[0] == id {
			if r[4] == "failed" {
				return "failed", nil // a refused start is a refusal of the probe, not the run's error
			}
			return r[4], err
		}
	}
	if err == nil {
		err = fmt.Errorf("the probe %s wrote no row of harnesses.tsv", id)
	}
	return "", err
}

// probeOf holds what a probe session adds to a TaskSpec.
type probeOf struct {
	token  string
	models [][]string
	sweep  bool
}

// probeText gives the prompt of a probe whose result file is result.
func probeText(result, token string) []byte {
	return []byte(strings.NewReplacer("{result}", result, "{token}", token).Replace(probePrompt))
}

// endProbe is the end of a probe session: the class, probe.tsv and the pass
// rules; one records commit of its row of harnesses.tsv and its telemetry row;
// then its comment on the control issue (a probe has no push, so its records
// are complete here).
func (s *Sessions) endProbe(ctx context.Context, spec TaskSpec, id string, d session.Dir, r session.Run, err error) error {
	_, rows, resultErr := session.ResultFile(d, true)
	class := session.Class(r, err, resultErr)
	if class == "" {
		return err
	}
	start, err := s.startRowOf(ctx, id)
	if err != nil {
		return err
	}
	u, uErr := session.UsageOf(spec.Pair.Usage, filepath.Join(d.Root, "stdout"))
	if uErr != nil && !errors.As(uErr, new(*fs.PathError)) {
		return uErr
	}
	reason := probeReason(class, resultErr, rows, spec.probe.token, d, spec.Pair.Policy, u.Models, spec.probe.models, spec.Pair.Harness)
	var files []string
	for _, row := range rows {
		if row[0] == "file" {
			files = append(files, row[1])
		}
	}
	end := r.End
	if end.IsZero() {
		end = r.Start
	}
	row, err := s.harnessRow(id, spec.Pair.Harness, start[sessionVersion], spec.Pair.Model, reason, files, u.Models, end)
	if err != nil {
		return err
	}
	harnesses, err := s.appendHarness(ctx, row)
	if err != nil {
		return err
	}
	tel, err := s.telemetry(ctx, spec, start, d, r)
	if err != nil {
		return err
	}
	if err := s.commit(ctx, map[string][]byte{"harnesses.tsv": harnesses, "telemetry.tsv": tel}, fmt.Sprintf("layup run: the probe %s of %s", id, spec.Pair.Harness)); err != nil {
		return err
	}
	result := "passed"
	if reason != "" {
		result = "failed " + reason
	}
	_, err = s.Forge.Comment(ctx, s.Control, fmt.Sprintf("%s: probe of %s %s: %s\n", id, spec.Pair.Harness, start[sessionVersion], result))
	return err
}

// harnessRow gives a row of harnesses.tsv; reason "" is passed.
func (s *Sessions) harnessRow(id, harness, version, model, reason string, files, models []string, end time.Time) ([]string, error) {
	result := "passed"
	if reason != "" {
		result = "failed"
	}
	f, err := tsv.JoinList(files)
	if err != nil {
		return nil, err
	}
	m, err := tsv.JoinList(models)
	if err != nil {
		return nil, err
	}
	return []string{id, harness, version, model, result, reason, f, m, end.UTC().Format(timeForm)}, nil
}

// appendHarness gives harnesses.tsv with row added, after CheckHarness.
func (s *Sessions) appendHarness(ctx context.Context, row []string) ([]byte, error) {
	before, err := s.readTable(ctx, "harnesses.tsv", records.ReadHarnesses)
	if err != nil {
		return nil, err
	}
	if err := records.CheckHarness(row); err != nil {
		return nil, err
	}
	return table(records.HarnessesSchema, append(append([][]string{}, before...), row))
}

// refuseProbe writes the row of a probe whose start was refused: failed,
// with the reason and its value, and its version, which is — when the
// version check refused it (no version was read); no start row, no telemetry row, and no comment (decided here,
// task T-nxe4: its record is its row, as a task session's refused start has
// no comment).
func (s *Sessions) refuseProbe(ctx context.Context, spec TaskSpec, id, version string, refusal session.Refusal) error {
	row, err := s.harnessRow(id, spec.Pair.Harness, version, spec.Pair.Model, strings.TrimSpace(refusal.Reason+" "+refusal.Value), nil, nil, s.Now())
	if err != nil {
		return err
	}
	data, err := s.appendHarness(ctx, row)
	if err != nil {
		return err
	}
	return s.commit(ctx, map[string][]byte{"harnesses.tsv": data}, fmt.Sprintf("layup run: the probe %s of %s refused (%s)", id, spec.Pair.Harness, refusal.Reason))
}

// Admit gives the admission hook of a task session on the harness row h: nil
// when the last probe at the version read passed; else it runs the probe
// first (with no sweep), and a probe that fails refuses the start (probe).
func (s *Sessions) Admit(ctx context.Context, p ProbeSpec) func(version string) error {
	return func(version string) error {
		rows, err := s.readTable(ctx, "harnesses.tsv", records.ReadHarnesses)
		if err != nil {
			return err
		}
		if route.Probed(rows, p.Harness[regHarness], version) {
			return nil
		}
		p.Sweep = false
		result, err := s.Probe(ctx, p)
		if err != nil {
			return err
		}
		if result == "failed" {
			return session.Refusal{Reason: "probe"}
		}
		return nil
	}
}

// ProbeStep runs the step probe: it copies the routing register into the
// records when the two differ, skips each harness with no model of use yes,
// and probes each other one whose last probe at its version did not pass. It
// gives the three counts of the step's detail. recordsCommit gives the run's last
// pushed records commit, read at the start of each probe.
func (s *Sessions) ProbeStep(ctx context.Context, harnesses, models, routing [][]string, base string, recordsCommit func() string) (passed, failed, skipped int, err error) {
	s.fns()
	held, err := s.readTable(ctx, "routing.tsv", records.ReadRouting)
	if err != nil {
		return 0, 0, 0, err
	}
	if routingDiffers(held, routing) {
		data, err := table(records.RoutingSchema, routing)
		if err != nil {
			return 0, 0, 0, err
		}
		if err := s.commit(ctx, map[string][]byte{"routing.tsv": data}, "layup run: the routing register"); err != nil {
			return 0, 0, 0, err
		}
	}
	for _, h := range harnesses {
		if firstModel(models, h[regHarness]) == nil {
			skipped++
			continue
		}
		result, err := s.Probe(ctx, ProbeSpec{Harness: h, Models: models, Base: base, Records: recordsCommit(), Sweep: true})
		if err != nil {
			return passed, failed, skipped, err
		}
		switch result {
		case "passed":
			passed++
		case "failed":
			failed++
		}
	}
	return passed, failed, skipped, nil
}

// probeStep is the step probe of the restart (run.md, The restart): after
// lease, before phase. Its detail is the three counts; it fails when a probe
// could not write its records.
func (r *state) probeStep(ctx context.Context) Step {
	const name = "probe"
	control, _ := strconv.Atoi(r.value("issue.control"))
	s := &Sessions{Store: r.store, RunID: r.cfg.RunID, Host: r.cfg.Dir, Target: r.cfg.Owner + "/" + r.cfg.Name, Clone: r.target,
		Branch: r.defaultBranch, Now: r.cfg.Clock.Now, Forge: r.cfg.Forge, Control: control, Prices: r.cfg.Prices}
	base, err := git.RevParse(r.target, "refs/remotes/origin/"+r.defaultBranch)
	if err != nil {
		return fail(name, err)
	}
	passed, failed, skipped, err := s.ProbeStep(ctx, r.cfg.Harnesses, r.cfg.Models, r.cfg.Routing, base, r.store.Base)
	if err != nil {
		return fail(name, err)
	}
	return done(name, fmt.Sprintf("probed and passed %d, probed and failed %d, skipped %d", passed, failed, skipped))
}

// ReadPrices reads host:prices.tsv, DIR/prices.tsv: none for a missing file,
// an empty table (docs/spec/session.md, Input states, task T-fsjp); a file
// that its reader refuses is an error that names it.
func ReadPrices(dir string) ([][]string, error) {
	path := filepath.Join(dir, "prices.tsv")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	rows, err := records.ReadPrices(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return rows, nil
}
