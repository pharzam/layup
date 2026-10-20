package run

// This file holds the brief of the first bet (docs/spec/loop.md, The brief of the first bet; task T-h65e, #418, row 66u of the plan of
// M2e): its inputs, its session, its checks, the appetite after the session's telemetry, the promotion and the basis. Row 66j calls
// betBrief and posts the brief.

import (
	"context"
	"errors"

	"github.com/pharzam/layup/internal/spec"
)

// ErrBriefLinesUnset and ErrBriefOverLimit are the two waits before any session; row 66j maps them to done waiting for the Operator.
var (
	ErrBriefLinesUnset = errors.New("brief.lines is not set")
	ErrBriefOverLimit  = errors.New("the parts of code are over brief.lines")
)

// BriefBasis is what a kept brief rests on, read from bets/<milestone>/basis.tsv and the brief.
type BriefBasis struct {
	Requirements []string            // the final IDs
	Survivors    map[string][]string // the survivors of each decision, in the order of the filter
	Dropped      map[string]string   // the reasons of each dropped option, by <decision>=<option>
	Facts        []spec.Fact         // the number or Q-NNN of each fact with its class
	Open         []string            // the Q-NNN of the questions still open
	Heads        map[string]string   // the verified head of each Shape task, by task ID
	Appetite     spec.Cap            // the numbers of the brief's appetite line
}

// headRef is a Shape task with its verified head and its pull request.
type headRef struct {
	Task, Kind, Head string
	PR               int
}

// briefParts are the four texts that code writes into the brief.
type briefParts struct{ requirements, solution, open, facts string }

const (
	briefPrompt  = ""
	briefRefused = ""
)

func briefPath(milestone string) string   { return "bets/" + milestone + "/brief.md" }
func basisPath(milestone string) string   { return "bets/" + milestone + "/basis.tsv" }
func attemptsDir(milestone string) string { return "bets/" + milestone + "/attempts" }

func checkMilestone(milestone string) error { return nil }

func requirementsPart(ids []string, prd spec.PRDContent) (string, error) { return "", nil }
func solutionPart(owner, name string, heads []headRef, results map[string]spec.Result) string {
	return ""
}
func openPart(questions [][]string) string                          { return "" }
func factsPart(spans, questions, answers [][]string) string         { return "" }
func answerFactClass(id string, questions, spans [][]string) string { return "" }
func linkSet(owner, name string, heads []headRef, blobs map[string][]string) map[string]bool {
	return nil
}
func assemble(p briefParts, problem, appetite, holes string) string { return "" }

// briefKeeper judges the output of a session and a stored candidate.
type briefKeeper struct {
	parts briefParts
	multi map[string][]string // every option of each decision with more than one survivor
	lines *int
	links map[string]bool
}

func (k briefKeeper) judge(assembled string) []spec.Finding { return nil }
func (k briefKeeper) keep(milestone, id string, files map[string][]byte) (map[string][]byte, error) {
	return nil, nil
}
func (k briefKeeper) overLimit() error { return nil }

func appetiteFindings(brief string, telemetry [][]string, bandU, hours string) ([]spec.Finding, error) {
	return nil, nil
}

func briefText(milestone string, lines, room int, finding string) func(input map[string]string, output string) []byte {
	return func(map[string]string, string) []byte { return nil }
}

func (r *state) betBrief(ctx context.Context, milestone, kind string) (brief []byte, done bool, err error) {
	return nil, false, nil
}

func (r *state) briefBasis(ctx context.Context, milestone string) (BriefBasis, error) {
	return BriefBasis{}, nil
}
