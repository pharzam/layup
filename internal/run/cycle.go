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

// The command `cycle <task> because <reason>` of the Operator (docs/spec/loop.md, One more cycle at the cap; task T-0xr7, #493): one more
// fix-and-review cycle for a Shape task that is at the review-cycle cap. The copy of the comment in copies.tsv is the record, and nothing else is
// written: the grant is read, on every view, from the copies and from the task's own events.

// cycleTask is the form id(T-xxxx) of a task ID (internal/tsv, idPattern).
var cycleTask = regexp.MustCompile(`^T-[a-z0-9]{4}$`)

// A cycleLine is the two parts of a command: its task and its reason.
type cycleLine struct{ Task, Reason string }

// parseCycle reads the body of a comment as `cycle <task> because <reason>`: one line, once its trailing white space is trimmed, with one space between
// the parts, a task of the form id(T-xxxx), and a reason that is not empty and not —. A body whose first word is not cycle is no command (false, no
// error); one that starts with the word and breaks the form is a command that is refused (true, the error).
func parseCycle(body string) (cycleLine, bool, error) {
	text := strings.TrimRight(body, " \t\r\n")
	first := text
	if i := strings.IndexAny(text, " \t\r\n"); i >= 0 {
		first = text[:i]
	}
	if first != "cycle" {
		return cycleLine{}, false, nil
	}
	const form = "the form is `cycle <task> because <reason>`, on one line"
	if strings.ContainsAny(text, "\r\n") {
		return cycleLine{}, true, fmt.Errorf("%s: this comment has more than one line", form)
	}
	f := strings.Split(text, " ")
	if len(f) < 4 || !cycleTask.MatchString(f[1]) || f[2] != "because" {
		return cycleLine{}, true, fmt.Errorf("%s, with a task of the form T-xxxx", form)
	}
	reason := strings.Join(f[3:], " ")
	if reason == "" || reason == "—" {
		return cycleLine{}, true, fmt.Errorf("%s: the reason is not empty and not —", form)
	}
	return cycleLine{Task: f[1], Reason: reason}, true, nil
}

// judgeCommand judges a copy against the first three items of the order of docs/spec/loop.md: the form; the author (a first copy, by no App, by an
// account of approvers.tsv in the role operator); the issue (the control issue). It gives the command, whether the body is one, and the refusal.
// The task, the cap and the bound are judged against the records of the task (judgeGrants and the boundary).
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

// capRound gives the review record of the cap round: the round number shapeCycleCap+1 among the rounds that reviewRounds counts, in the order of the
// sessions. It gives the comment of its valid handoff and its verdict; found is false when the task has fewer rounds.
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

// judgeGrants judges the commands that name the task of v against its recorded history, in the order of the comment IDs, and gives the comment of the
// grant (0 for none) and the refusals of the items cap and bound. The order of the items for a command is: form, author, issue (judgeCommand), task
// (the boundary), cap, bound. A command is eligible only after the cap round, that is the second counted review round of the task, and only if that
// round did not end `not mergeable, findings recorded`, whose wait grants no cycle and so consumes no grant; its comment ID is strictly greater than
// that of the cap round's record (the IDs of one forge grow with time). The earliest eligible command is the grant, and each later one is refused
// with the item bound. Nothing in the judgement reads the position, so a grant stays accepted after its attempt starts.
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
		if err != nil {
			return nil, fmt.Errorf("the copy %s: %w", row[copyBody], err)
		}
		if body == nil {
			return nil, fmt.Errorf("the copy %s is not in the records", row[copyBody])
		}
		if _, isCycle, _ := parseCycle(string(body)); isCycle {
			out = append(out, setInput{Row: row, Body: string(body)})
		}
	}
	return out, nil
}

// cycleBoundary is what the boundary calls after the copies are made: it judges every command of the copies and reports each refusal as a progress line
// `cycle: refused comment N: <item>: <reason>`, the items being form, author, issue, task, cap and bound, in that order of judgement. A refusal changes and
// records nothing, and the same command is judged again at the next boundary with the same result; a grant needs no write, because its copy is the record.
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
	data, err := r.store.ReadFile(ctx, tasksPath)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if data != nil {
		if tasks, err = ReadRegister(data); err != nil {
			return nil, fmt.Errorf("%s: %w", tasksPath, err)
		}
	}
	var refused []Refusal
	judged := map[string]bool{}
	for _, c := range cands {
		line, _, ref := judgeCommand(c, control, r.approvers)
		if ref == nil && !slices.ContainsFunc(tasks, func(t Task) bool { return t.ID == line.Task }) {
			id, _ := strconv.ParseInt(c.Row[copyComment], 10, 64)
			ref = &Refusal{Comment: id, Item: "task", Reason: fmt.Sprintf("%s is not a task of tasks.tsv", line.Task)}
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
