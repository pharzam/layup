package setup

// The steps of the setup in phase 2 (docs/spec/scaffold.md, REQ-002 — The setup
// under layup run, The steps of the setup in phase 2; task T-ktgx, #273):
// layup run (row 54 of M2d, task T-eahr) calls the step runner of this package
// with the inputs of phase 2 given in Go, not as a file. StepsPhase2 takes the steps of phase 1 and changes
// the six that differ: S01 to S04 and S13 have the code of phase 2 (S13: the two
// rulesets, task T-9v92), and S15 gives no records commit and no command. S05 to
// S12 and S14 are phase 1's, with the same Calls; the commands of S13 are row 52c.

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/work"
)

// Phase2 is what layup run gives the steps in place of a person, the forge and
// the baseline's repository: the values of the records of Start, the questions
// of the batch and the App's ID, which the work area and its files do not hold.
type Phase2 struct {
	Name          string // OWNER/NAME of the target
	Visibility    string // forge.visibility of start.tsv: public or private
	PinSource     string // the rows pin.source, pin.commit, pin.tree and pin.time of start.tsv
	PinCommit     string
	PinTree       string
	PinTime       string
	BriefSum      string     // psb.sha256 of start.tsv
	AppID         int64      // app_id of the host's forge register (not of the records of Start): the pin of the layup/ checks of S13 and its bypass actor
	Host          string     // the --host DIR of layup run, the root of DIR/work/OWNER/NAME/ and DIR/applies/OWNER/NAME/: the paths of the commands of S13
	DefaultBranch string     // the default branch of the forge's repository: the branch that the commands of S13 read and push to
	CodeOwner     string     // the login of the Operator's row of approvers.tsv: the owner of the gate paths in .github/CODEOWNERS (S13)
	Migrate       *Migration // non-nil for a target set up before the code-owner rule: S13 writes the files of M2d into the tree, and the update by the recorded ruleset ID
	// Questions are the Q- questions of the batch, each with the question of
	// its row of questions.tsv: the raw fact of the answers names them (S04).
	Questions []work.Question
}

// runActor is the actor of S02, S03 and S15 in the step table in phase 2
// (scaffold.md, after the table of the steps): the step is layup run's.
const runActor = "layup-run"

// startRef is the ref of a value that layup run takes from the records of Start.
const startRef = "records:start/start.tsv"

// StepsPhase2 gives the steps of the setup in phase 2: the steps of phase 1
// with S01 to S04 and S13 run by the code of phase 2, and S15 with no records commit
// and no command (docs/spec/scaffold.md, The steps of the setup in phase 2). The
// evidence of S04 stays that of phase 1: the checks pin and facts.
func StepsPhase2(p Phase2, c Calls) map[string]Step {
	m := Steps(Brief{}, c)
	set := func(id string, change func(*Step)) {
		s := m[id]
		change(&s)
		m[id] = s
	}
	set("S01", func(s *Step) {
		s.Run = func(in Input) Outcome { return runS01Phase2(p, in) }
		s.Unchanged = func(r work.Record) error { return briefUnchanged(Brief{Sum: p.BriefSum}, r) }
	})
	set("S02", func(s *Step) { s.Actor, s.Run = runActor, func(Input) Outcome { return runS02Phase2(p) } })
	set("S03", func(s *Step) { s.Actor, s.Run, s.Commands = runActor, runS03Phase2, nil })
	set("S04", func(s *Step) { s.Run = func(in Input) Outcome { return runS04With(p.asked, in) } })
	set("S13", func(s *Step) {
		s.Run = func(in Input) Outcome { return runS13Phase2(p, in) }
		s.Commands = commandsS13Phase2(p)
	})
	set("S15", func(s *Step) { s.Actor, s.Commands, s.Records = runActor, nil, nil })
	return m
}

// runS01Phase2 is S01 in phase 2: it asks the stack and nothing else, runs no
// check of the problem statement, and takes the name, the visibility and the
// baseline from the records of Start. The Q- rows of the answers are the
// questions of the batch, which it takes as given. The checks of the values are
// those of phase 1, which S12 and the pin file build on; a value that fails is
// a fail, not an input error, as layup run made the work area.
func runS01Phase2(p Phase2, in Input) Outcome {
	var stack []string
	for _, r := range in.Answers {
		if r[0] == work.S01Questions[0].ID && r[1] != "" {
			stack = r
		}
	}
	if stack == nil {
		q := work.S01Questions[0]
		return Outcome{Kind: Stop, Stops: []StopRow{{"S01", q.ID, q.Text, ""}}}
	}
	switch {
	case !hasEntry(stack[1]):
		return Outcome{Kind: Fail, Evidence: "the stack " + stack[1] + " has no entry in the catalog"}
	case !ownerName(p.Name):
		return Outcome{Kind: Fail, Evidence: "name: " + p.Name + " is not OWNER/NAME"}
	case p.Visibility != "public" && p.Visibility != "private":
		return Outcome{Kind: Fail, Evidence: "visibility: " + p.Visibility + " is not public or private"}
	case !baselineURL(p.PinSource):
		return Outcome{Kind: Fail, Evidence: "baseline: " + p.PinSource + " is not a URL with the scheme https, http, git or file"}
	}
	return Outcome{Kind: Done, Evidence: "S01-stack answered; the stack has a catalog entry", Values: [][]string{
		append([]string{"stack", stack[1]}, sourceOf(stack)...),
		{"name", p.Name, "computed", startRef},
		{"visibility", p.Visibility, "computed", startRef},
		{"baseline", p.PinSource, "computed", startRef},
		{"brief.sha256", p.BriefSum, "computed", "sha256 " + work.BriefPath},
	}}
}

// runS02Phase2 is S02 in phase 2: it resolves nothing. The pin is that of
// start.tsv (NFR-006: no run resolves the commit again), so it calls no git,
// clones nothing, reads no clock and writes nothing. The forms of the values
// are checked by internal/records at each restart (CheckStart).
func runS02Phase2(p Phase2) Outcome {
	return Outcome{Kind: Done, Evidence: "the commit and the tree of start.tsv", Values: [][]string{
		{"pin.source", p.PinSource, "computed", startRef},
		{"pin.commit", p.PinCommit, "computed", startRef},
		{"pin.tree", p.PinTree, "computed", startRef},
		{"pin.time", p.PinTime, "computed", startRef},
	}}
}

// runS03Phase2 is S03 in phase 2: it makes no commit and writes no command.
// WORK/target is the clone of the default branch with main at the root commit
// that Start read back, so it checks that the tree of main is pin.tree.
func runS03Phase2(in Input) Outcome {
	tree, ok := in.Record.Value("S02", "pin.tree")
	if !ok {
		return Outcome{Kind: Fail, Evidence: "the record has no row S02 pin.tree"}
	}
	got, err := sys.revParse(filepath.Join(in.Dir, work.TargetPath), "refs/heads/main^{tree}")
	switch {
	case err != nil:
		return Outcome{Kind: Fail, Evidence: "the tree of main: " + firstLine(err)}
	case got != tree:
		return Outcome{Kind: Fail, Evidence: fmt.Sprintf("the root tree %s is not pin.tree %s", got, tree)}
	}
	return Outcome{Kind: Done, Evidence: "root tree = pin.tree"}
}

// asked gives the questions that the raw fact of S04 records in phase 2:
// S01-stack, then each question of the batch, each with its text (a question
// with no row in the answers has no fact). A Q- row whose question the batch
// does not hold is a fail: layup run gave the work area, so it is its fault, not
// an input error.
func (p Phase2) asked(in Input) ([]work.Question, Outcome, bool) {
	for _, r := range in.Answers {
		if strings.HasPrefix(r[0], "Q-") && !slices.ContainsFunc(p.Questions, func(q work.Question) bool { return q.ID == r[0] }) {
			return nil, Outcome{Kind: Fail, Evidence: "the question " + r[0] + " is not one of the questions of the batch, so its text is not known"}, false
		}
	}
	return append([]work.Question{work.S01Questions[0]}, p.Questions...), Outcome{}, true
}
