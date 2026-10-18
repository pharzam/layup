package run

// This file holds the parse of the three comment forms of the task loop and the
// four handoff items read from them (docs/spec/loop.md, The items a caller
// computes: The fields and the blocks of a comment; task T-4pje, #396, row 66l of
// the plan of M2e). The step that wires the parse into the loop is row 66k's.

// A PlanTest is a line of the plan's tests: block, a test's ID and its path.
type PlanTest struct{ ID, Path string }

// A ReviewRecord is a parsed review record: its verdict, the head it names and
// its located findings.
type ReviewRecord struct {
	Verdict, Head string
	Findings      []Finding
}

// ParsePlanComment parses a plan comment in the pinned form.
func ParsePlanComment(data []byte) (Plan, []PlanTest, bool) { return Plan{}, nil, false }

// ParsePlanReview parses a plan review comment and gives its verdict word.
func ParsePlanReview(data []byte) (string, bool) { return "", false }

// ParseReviewRecord parses a review record comment in the pinned form.
func ParseReviewRecord(data []byte) (ReviewRecord, bool) { return ReviewRecord{}, false }

// PlanFieldsOK is the item plan-fields of a plan comment's kept copy.
func PlanFieldsOK(data []byte) bool { return false }

// VerdictFieldOK is the item verdict-field of a plan review's kept copy.
func VerdictFieldOK(data []byte) bool { return false }

// FindingsLocatedOK is the item findings-located of a review record's kept copy.
func FindingsLocatedOK(data []byte) bool { return false }

// NothingMaterialOK is the item nothing-material of a review record's kept copy
// and the head that code holds for the verifier's session.
func NothingMaterialOK(data []byte, head string) bool { return false }

// OpenMaterial is the predicate that fails uphill: true when the record did not
// parse, its verdict is not "nothing material in scope", or its block holds a
// finding, so no invalid or empty record reads as no open material finding.
func OpenMaterial(r ReviewRecord, ok bool) bool { return false }
