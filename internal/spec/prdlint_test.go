package spec

// The unit tests of the prd-lint seam (task T-4wyg, #351): the anchored
// classifier and the verdict's exact accept condition (review conditions C1
// and C2). The real sh run is the integration test.

import (
	"strings"
	"testing"
)

const (
	waitMoSCoW = `FAIL  PRD-0001-x.md: REQ-001 MoSCoW "" not in Must|Should|Could|Won't`
	waitPhase  = `FAIL  PRD-0001-x.md: NFR-001 has no phase value`
)

func TestEmptyCellFailureIsAnchored(t *testing.T) {
	for _, l := range []string{waitMoSCoW, waitPhase,
		`FAIL  PRD-1234-a-2.md: REQ-999 has no phase value`} {
		if !emptyCellFailure.MatchString(l) {
			t.Errorf("the classifier refuses the producer's line %q", l)
		}
	}
	for _, l := range []string{
		`FAIL  PRD-0001-x.md: REQ-001 MoSCoW "—" not in Must|Should|Could|Won't`,
		`FAIL  PRD-0001-x.md: REQ-001 is Won't but Phase is not —`,
		`FAIL  PRD-0001-x.md: REQ-001 cites no resolvable F-NNNN fact`,
		`FAIL  PRD-bad.md: filename must be PRD-NNNN-kebab-case.md`,
		`FAIL  PRD directory not found: /tmp/x`,
		`FAIL  PRD-0001-x.md: REQ-001 has no phase value, really`,
		` FAIL  PRD-0001-x.md: REQ-001 has no phase value`,
		`FAIL  PRD-0001-x.md: REQ-001 has no phase value `,
		`prd-lint: OK`,
		`awk: syntax error`,
	} {
		if emptyCellFailure.MatchString(l) {
			t.Errorf("the classifier accepts the non-empty-cell line %q", l)
		}
	}
}

// waitingVerdict is the pure decision over the producer's exit code and the
// non-empty lines of both streams: exit 0 accepts; exit 1 accepts only when
// every line is a classified empty-cell failure and at least one exists — a
// vacuous or a foreign line refuses (C1).
func TestWaitingVerdict(t *testing.T) {
	for _, c := range []struct {
		name        string
		code        int
		stdout      []string
		stderr      []string
		want        bool
		wantOffense string // "" when want is true
	}{
		{"clean", 0, []string{"prd-lint: OK"}, nil, true, ""},
		{"all empty cells", 1, []string{waitMoSCoW, waitPhase}, nil, true, ""},
		{"a stray stderr line refuses", 1, []string{waitMoSCoW}, []string{`FAIL  PRD-bad.md: filename must be PRD-NNNN-kebab-case.md`}, false, "filename"},
		{"a non-FAIL stderr refuses", 1, []string{waitMoSCoW}, []string{"awk: syntax error"}, false, "awk"},
		{"no lines refuses", 1, nil, nil, false, ""},
		{"ok text under exit 1 refuses", 1, []string{"prd-lint: OK"}, nil, false, "OK"},
		{"a real violation refuses", 1, []string{waitMoSCoW, `FAIL  PRD-0001-x.md: REQ-002 cites no resolvable F-NNNN fact`}, nil, false, "resolvable"},
		{"the em-dash MoSCoW refuses", 1, []string{`FAIL  PRD-0001-x.md: REQ-001 MoSCoW "—" not in Must|Should|Could|Won't`}, nil, false, "MoSCoW"},
		{"a Won't phase is no empty cell", 1, []string{`FAIL  PRD-0001-x.md: REQ-001 is Won't but Phase is not —`}, nil, false, "Won't"},
	} {
		t.Run(c.name, func(t *testing.T) {
			ok, offend := waitingVerdict(c.code, c.stdout, c.stderr)
			if ok != c.want {
				t.Errorf("waitingVerdict gives %v, want %v", ok, c.want)
			}
			if !c.want && c.wantOffense != "" {
				joined := strings.Join(offend, "\n")
				if !strings.Contains(joined, c.wantOffense) {
					t.Errorf("the offending lines %q name none of %q", joined, c.wantOffense)
				}
			}
		})
	}
}
