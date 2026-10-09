package route

import (
	"errors"
	"slices"
	"strconv"

	"github.com/pharzam/layup/internal/records"
)

// ErrNoPair is the refusal pair of docs/spec/session.md: no pair of the role's
// list for the task's tier is admitted. It comes before a session ID is drawn,
// so the event refused that the caller writes has no session (row 36a).
var ErrNoPair = errors.New("no admitted pair in the routing register for the role and the tier")

// Probed reports whether the harness's last probe at version passed: the last
// row of records:harnesses.tsv of that harness and version, in the order of
// the file (one row per probe, appended), has the result passed. No such row
// is false, so the version needs a probe (docs/spec/session.md, Admission).
func Probed(harnesses [][]string, harness, version string) bool {
	passed := false
	for _, r := range harnesses {
		f := fields(records.HarnessesSchema, r)
		if f["harness"] == harness && f["version"] == version {
			passed = f["result"] == "passed"
		}
	}
	return passed
}

// Admitted reports whether a pair of a harness and a model is admitted: the
// harness is Probed at version, and the model's row of models.tsv for that
// harness has use yes. A model with no row, or use no, is not admitted. Code
// alone admits; the caller gives the version that its version check read.
func Admitted(harnesses, models [][]string, harness, version, model string) bool {
	if !Probed(harnesses, harness, version) {
		return false
	}
	for _, r := range models {
		f := fields(ModelsSchema, r)
		if f["harness"] == harness && f["model"] == model {
			return f["use"] == "yes"
		}
	}
	return false
}

// Pair gives the session's pair: of the rows of records:routing.tsv for the
// role and the tier, in the order of position (the file's order is free), the
// first for which admitted holds. With none, ErrNoPair. Where the caller gets
// each harness's version for admitted is the caller's (row 36a).
func Pair(routing [][]string, role, tier string, admitted func(harness, model string) bool) (string, string, error) {
	type pair struct {
		position       int
		harness, model string
	}
	var list []pair
	for _, r := range routing {
		f := fields(records.RoutingSchema, r)
		if f["role"] == role && f["tier"] == tier {
			p, _ := strconv.Atoi(f["position"])
			list = append(list, pair{p, f["harness"], f["model"]})
		}
	}
	slices.SortFunc(list, func(a, b pair) int { return a.position - b.position })
	for _, p := range list {
		if admitted(p.harness, p.model) {
			return p.harness, p.model, nil
		}
	}
	return "", "", ErrNoPair
}
