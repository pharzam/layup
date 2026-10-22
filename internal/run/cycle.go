package run

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/records"
)

// The command `cycle <task> because <reason>` of the Operator (docs/spec/loop.md, One more cycle at the cap; task T-0xr7, #493): one more fix-and-review
// cycle for a Shape task at the review-cycle cap. The copy of the comment is the record; the grant is read from the copies and the task's events.

// cycleTask is the form id(T-xxxx) of a task ID (internal/tsv, idPattern).
var cycleTask = regexp.MustCompile(`^T-[a-z0-9]{4}$`)

// A cycleLine is the two parts of a command: its task and its reason.
type cycleLine struct{ Task, Reason string }

// parseCycle reads a body as `cycle <task> because <reason>`: one line, trailing white space trimmed, one space between the parts, a task of the form
// id(T-xxxx), a reason that is not empty and not —. A first word other than cycle is no command (false, no error); a broken form is refused (true, the error).
func parseCycle(body string) (cycleLine, bool, error) {
	text := strings.TrimRight(body, " \t\r\n")
	first := text
	if i := strings.IndexAny(text, " \t\r\n"); i >= 0 {
		first = text[:i]
	}
	if first != "cycle" {
		return cycleLine{}, false, nil
	}
	const form = "the form is `cycle <task> because <reason>`, on one line, with a task of the form T-xxxx"
	f := strings.Split(text, " ")
	switch reason := strings.Join(f[min(3, len(f)):], " "); {
	case strings.ContainsAny(text, "\r\n"):
		return cycleLine{}, true, fmt.Errorf("%s: this comment has more than one line", form)
	case len(f) < 4 || !cycleTask.MatchString(f[1]) || f[2] != "because" || f[3] == "": // f[3] is empty after a second space
		return cycleLine{}, true, fmt.Errorf("%s", form)
	case reason == "" || reason == "—":
		return cycleLine{}, true, fmt.Errorf("%s: the reason is not empty and not —", form)
	default:
		return cycleLine{Task: f[1], Reason: reason}, true, nil
	}
}

// judgeCommand judges a copy by the first three items of the order of loop.md: the form; the author (a first copy, by no App, by an account of the role
// operator); the issue (the control issue). The task, the cap and the bound are judged against the records (the boundary and judgeGrants).
func judgeCommand(c setInput, control int, approvers [][]string) (cycleLine, bool, *Refusal) {
	line, isCycle, err := parseCycle(c.Body)
	if !isCycle {
		return cycleLine{}, false, nil
	}
	id, _ := strconv.ParseInt(c.Row[copyComment], 10, 64)
	refuse := func(item, format string, a ...any) (cycleLine, bool, *Refusal) {
		return line, true, &Refusal{Comment: id, Item: item, Reason: fmt.Sprintf(format, a...)}
	}
	switch {
	case err != nil:
		return refuse("form", "%v", err)
	case !Decision(c.Row, approvers, "operator"):
		return refuse("author", "a grant is the first copy of a comment by no App, by an account of the role operator in approvers.tsv")
	case c.Row[copyIssue] != strconv.Itoa(control):
		return refuse("issue", "a grant is a comment of the control issue %d, and this one is of the issue %s", control, c.Row[copyIssue])
	}
	return line, true, nil
}

// capRound gives the comment of the valid handoff and the verdict of the cap round: round number shapeCycleCap+1 among those reviewRounds counts, in the
// order of the sessions; found is false when the task has fewer.
func capRound(v loopView) (comment int64, verdict string, found bool) {
	n := 0
	for _, s := range v.Sessions {
		rec, ok := v.Records[s[sessionID]]
		if _, verifier := v.qaKind(s); !ok || !verifier || rec.Verdict == nothingMaterial {
			continue
		}
		if n++; n == shapeCycleCap+1 {
			detail, _ := v.lastHandoff(s[sessionID])
			var cid, seen int64
			if _, err := fmt.Sscanf(detail, "valid comment %d-%d", &cid, &seen); err != nil {
				return 0, "", false
			}
			return cid, rec.Verdict, true
		}
	}
	return 0, "", false
}

// judgeGrants judges the commands that name the task of v against its recorded history, in the order of the comment IDs, and gives the grant (a comment
// ID, 0 for none) and the refusals of the items cap and bound (the order: form, author, issue, task, cap, bound). A command is eligible if the task has a
// cap round whose verdict is not `not mergeable, findings recorded` (that wait grants no cycle and consumes no grant) and its ID is strictly greater than
// that of the cap round's record (the IDs of one forge grow with time). The earliest eligible command is the grant; each later one is refused with the
// item bound. Nothing reads the position, so a grant stays accepted after its attempt starts.
func judgeGrants(v loopView, cands []setInput) (grant int64, refused []Refusal) {
	capID, verdict, found := capRound(v)
	sorted := slices.Clone(cands)
	idOf := func(c setInput) int64 { n, _ := strconv.ParseInt(c.Row[copyComment], 10, 64); return n }
	sort.SliceStable(sorted, func(i, j int) bool { return idOf(sorted[i]) < idOf(sorted[j]) })
	for _, c := range sorted {
		line, isCycle, ref := judgeCommand(c, v.Control, v.Approvers)
		if !isCycle || ref != nil || line.Task != v.Task {
			continue
		}
		id := idOf(c)
		refuse := func(item, format string, a ...any) {
			refused = append(refused, Refusal{Comment: id, Item: item, Reason: fmt.Sprintf(format, a...)})
		}
		switch {
		case !found:
			refuse("cap", "%s has not come to the review-cycle cap: it has fewer than %d rounds with findings", v.Task, shapeCycleCap+1)
		case verdict == "not mergeable, findings recorded":
			refuse("cap", "the cap round of %s ended not mergeable, findings recorded, a wait that grants no cycle", v.Task)
		case id <= capID:
			refuse("cap", "the comment %d is not later than the comment %d of the record of the cap round", id, capID)
		case grant != 0:
			refuse("bound", "%s has the grant of the comment %d, and one grant is the bound", v.Task, grant)
		default:
			grant = id
		}
	}
	return grant, refused
}

// grantOf is 1 when the task of v has a grant, else 0.
func grantOf(v loopView) int {
	if g, _ := judgeGrants(v, v.Cycles); g != 0 {
		return 1
	}
	return 0
}

// cycleCommands gives the copies of comments by no App whose body begins with the word cycle, with the bodies, in the order of the file.
func (s *Sessions) cycleCommands(ctx context.Context, copies [][]string) ([]setInput, error) {
	var out []setInput
	for _, row := range copies {
		if row[copyApp] != "" {
			continue
		}
		body, err := s.Store.ReadFile(ctx, row[copyBody])
		if err != nil || body == nil {
			return nil, fmt.Errorf("the copy %s: %v", row[copyBody], err)
		}
		if _, isCycle, _ := parseCycle(string(body)); isCycle {
			out = append(out, setInput{Row: row, Body: string(body)})
		}
	}
	return out, nil
}

// cycleBoundary is what the boundary calls once the copies are made: it judges every command of the copies and reports each refusal as a progress line
// `cycle: refused comment N: <item>: <reason>`. A refusal changes and records nothing and is judged again at the next boundary; a grant needs no write.
func (r *state) cycleBoundary(ctx context.Context) ([]Refusal, error) {
	control, err := strconv.Atoi(r.value("issue.control"))
	if err != nil {
		return nil, fmt.Errorf("no control issue")
	}
	s := r.sessions()
	copies, err := s.readTable(ctx, "copies.tsv", records.ReadCopies)
	if err != nil {
		return nil, err
	}
	cands, err := s.cycleCommands(ctx, copies)
	if err != nil || len(cands) == 0 {
		return nil, err
	}
	var refused []Refusal
	judged := map[string]bool{}
	for _, c := range cands {
		line, _, ref := judgeCommand(c, control, r.approvers)
		if _, err := s.taskRow(ctx, line.Task); ref == nil && err != nil {
			id, _ := strconv.ParseInt(c.Row[copyComment], 10, 64)
			ref = &Refusal{Comment: id, Item: "task", Reason: err.Error()}
		}
		if ref != nil {
			refused = append(refused, *ref)
			continue
		}
		if judged[line.Task] {
			continue
		}
		judged[line.Task] = true
		v, err := s.shapeView(ctx, line.Task)
		if err != nil {
			return refused, err
		}
		v.Control, v.Approvers = control, r.approvers
		_, more := judgeGrants(v, cands)
		refused = append(refused, more...)
	}
	sort.SliceStable(refused, func(i, j int) bool { return refused[i].Comment < refused[j].Comment })
	for _, ref := range refused {
		r.cfg.Progress(fmt.Sprintf("cycle: refused comment %d: %s: %s", ref.Comment, ref.Item, ref.Reason))
	}
	return refused, nil
}
