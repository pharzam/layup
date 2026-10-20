package run

// This file holds the architecture tasks that the bet opens (docs/spec/loop.md, The step bet; task T-2vva).

// settledDecisions gives the decisions that this bet settles, in name order: those of the answer's option lines that have more than one
// survivor (a decision with one survivor was written as decided before the bet).
func settledDecisions(options map[string]string, survivors map[string][]string) []string { return nil }

// archRequirements gives the requirement IDs of a chosen option from the constraint values of its rows: REQ-NNN and NFR-NNN kept when the
// final IDs name them, R-NNN mapped by finalID with its kind in the draft, any other value none. Sorted, once each; none is an error.
func archRequirements(decision, option string, rows, draft [][]string) ([]string, error) {
	return nil, nil
}

// bindingPath is the binding of an architecture task to the option of its decision.
func bindingPath(task, decision string) string {
	return "tasks/" + task + "/options/" + decision + ".tsv"
}
