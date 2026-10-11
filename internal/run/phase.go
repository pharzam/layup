package run

// This file holds a session of a phase (docs/spec/intake.md, The sessions of
// Intake; Input states, rows 1 and 3; row 42a of the plan, task T-gmey,
// #215): a role session with no attempt and no event, which reads its input
// from input/ and writes one output file, whose records the step keeps.

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"path/filepath"
	"strings"

	"github.com/pharzam/layup/internal/ledger"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/session"
)

// PhaseSpec is one session of a phase: the step and its role; the harness
// register row and the row of models.tsv of its pair, as PhasePair gives
// them, and models.tsv, for the probe of a version with no probe that passed;
// the base (the head of the default branch) and the records commit of the
// prompt; the input files by name, byte for byte from the records; the name
// of the output file under result/; Prompt, the fixed text of the step built
// from the absolute paths of the input files (by name) and of the output
// file; and Keep, which reads the output file by its block, runs the step's
// check, and gives the records to commit, by path, or refuses the file.
type PhaseSpec struct {
	Step, Role    string
	Harness       []string
	Model         []string
	Models        [][]string
	Base, Records string
	Input         map[string][]byte
	Output        string
	Prompt        func(input map[string]string, output string) []byte
	Keep          func(id string, data []byte) (map[string][]byte, error)
}

// phaseOf holds what a session of a phase adds to a TaskSpec.
type phaseOf struct {
	step   string
	input  map[string][]byte
	output string
	prompt func(input map[string]string, output string) []byte
	keep   func(id string, data []byte) (map[string][]byte, error)
}

// text gives the prompt of the session whose directory is root.
func (p *phaseOf) text(root string) []byte {
	paths := map[string]string{}
	for name := range p.input {
		paths[name] = filepath.Join(root, "input", name)
	}
	return p.prompt(paths, filepath.Join(root, "result", p.output))
}

// PhaseError is the end of a session of a phase that its step reports as
// fail: a refused start, a class that is not done, or an output file that
// is refused; Reason names which, with its detail.
type PhaseError struct{ Session, Reason string }

func (e *PhaseError) Error() string { return "the session " + e.Session + ": " + e.Reason }

// PhasePair gives the pair of a session of a phase for role: the first pair
// of records:routing.tsv for the role and the tier reasoning (route.Pair)
// whose harness is in the register harnesses and is admitted at the version
// of its last row of records:harnesses.tsv (decided here, task T-gmey: no
// version command runs before the session, whose own version check reads the
// version again and probes a new one). It gives the harness register row and
// the row of models.tsv; with none, route.ErrNoPair, and no session starts.
func (s *Sessions) PhasePair(ctx context.Context, role string, harnesses, models [][]string) ([]string, []string, error) {
	routing, err := s.readTable(ctx, "routing.tsv", records.ReadRouting)
	if err != nil {
		return nil, nil, err
	}
	probes, err := s.readTable(ctx, "harnesses.tsv", records.ReadHarnesses)
	if err != nil {
		return nil, nil, err
	}
	last := map[string]string{} // the version of each harness's last row
	for _, r := range probes {
		last[r[1]] = r[2]
	}
	row := func(harness string) []string {
		for _, h := range harnesses {
			if h[regHarness] == harness {
				return h
			}
		}
		return nil
	}
	harness, model, err := route.Pair(routing, role, "reasoning", func(h, m string) bool {
		v, probed := last[h]
		return probed && row(h) != nil && route.Admitted(probes, models, h, v, m)
	})
	if err != nil {
		return nil, nil, err
	}
	for _, m := range models {
		if m[0] == harness && m[1] == model {
			return row(harness), m, nil
		}
	}
	return nil, nil, route.ErrNoPair // not reached: Admitted holds only for a row of models
}

// PhaseSession runs one session of a phase through the one call of
// TaskSession and gives its ID: its task ID drawn as a probe's, attempt 1,
// its start row alone, no event and no push; the probe of its version first
// when the last probe at that version did not pass. A refused start, a class
// that is not done and a refused output file are a *PhaseError.
func (s *Sessions) PhaseSession(ctx context.Context, p PhaseSpec) (string, error) {
	s.fns()
	if err := s.milestoneCap(ctx, pairOf(p.Harness, p.Model).Cap); err != nil {
		return "", err
	}
	task, err := s.probeTaskID(ctx)
	if err != nil {
		return "", err
	}
	spec := TaskSpec{Task: task, Attempt: 1, Role: p.Role, Base: p.Base, Records: p.Records, Pair: pairOf(p.Harness, p.Model),
		phase: &phaseOf{step: p.Step, input: p.Input, output: p.Output, prompt: p.Prompt, keep: p.Keep}}
	spec.Admit = s.Admit(ctx, ProbeSpec{Harness: p.Harness, Models: p.Models, Base: p.Base, Records: p.Records})
	spec.End = func(id string, d session.Dir, r session.Run, err error) error {
		return s.endPhase(ctx, spec, id, d, r, err)
	}
	return s.TaskSession(ctx, spec)
}

// milestoneCap asks the budget of milestone 0 before a session of Intake
// whose spend cap is sessionCap ("" for none) starts (docs/spec/intake.md,
// The budget of Intake; task T-ybd5): intake.cap of start.tsv, the rows of
// telemetry.tsv at the run's records base, the time of Start and now. A
// refusal wraps ledger.ErrMilestoneCap, and no session starts.
func (s *Sessions) milestoneCap(ctx context.Context, sessionCap string) error {
	start, err := s.readTable(ctx, "start/start.tsv", func(d []byte) ([][]string, error) {
		rows, _, err := records.ReadStartAsWritten(d)
		return rows, err
	})
	if err != nil {
		return err
	}
	capValue := ""
	for _, r := range start {
		if r[0] == "intake.cap" {
			capValue = r[1]
		}
	}
	telemetry, err := s.readTable(ctx, "telemetry.tsv", records.ReadTelemetry)
	if err != nil {
		return err
	}
	started, err := s.Store.Started(ctx)
	if err != nil {
		return err
	}
	return ledger.MilestoneCap(capValue, telemetry, sessionCap, started, s.Now())
}

// endPhase is the end of a session of a phase: the class, with the output
// file read by the rules of a result file; for done, the step's Keep; one
// records commit of the kept records and the telemetry row, or of the
// telemetry row alone when the class is not done or Keep refuses the file
// (decided here, task T-gmey); then its comment on the control issue.
func (s *Sessions) endPhase(ctx context.Context, spec TaskSpec, id string, d session.Dir, r session.Run, err error) error {
	data, readErr := session.ReadOutput(d, spec.phase.output)
	class := session.Class(r, err, readErr)
	if class == "" {
		return err
	}
	start, err := s.startRowOf(ctx, id)
	if err != nil {
		return err
	}
	files, reason, refused := map[string][]byte{}, "", ""
	switch {
	case class == "no-result":
		reason = class + ": " + readErr.Error()
	case class != "done":
		reason = class
	default:
		kept, err := spec.phase.keep(id, data)
		if err != nil {
			reason, refused = "refused file: "+err.Error(), ", refused file"
		} else {
			maps.Copy(files, kept)
		}
	}
	tel, err := s.telemetry(ctx, spec, start, d, r)
	if err != nil {
		return err
	}
	files["telemetry.tsv"] = tel
	if err := s.commit(ctx, files, fmt.Sprintf("layup run: the end of %s, the session of %s, %s", id, spec.phase.step, class)); err != nil {
		return err
	}
	// A comment that fails keeps the session's reason (round 1 of #215).
	var end error
	if reason != "" {
		end = &PhaseError{Session: id, Reason: reason}
	}
	line := fmt.Sprintf("%s: %s session of %s: %s%s\n", id, spec.Role, spec.phase.step, class, refused)
	if _, err := s.Forge.Comment(ctx, s.Control, line); err != nil {
		return errors.Join(end, err)
	}
	return end
}

// refusePhase gives the error of a session of a phase whose start was
// refused, with the reason and its value: no record (decided here, task
// T-gmey: no event holds it, and the row of a refused probe is the probe's
// own).
func refusePhase(id string, refusal session.Refusal) error {
	return &PhaseError{Session: id, Reason: strings.TrimSpace("refused " + refusal.Reason + " " + refusal.Value)}
}
