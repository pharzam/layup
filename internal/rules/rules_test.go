package rules

import (
	"strings"
	"testing"
)

const header = "pattern\texception\tsource"

// register is a register of the entries of phase 1 (docs/spec/setup.md, The
// rule-path register), with the exception of docs/guardrails.md.
var register = header + "\n" + strings.Join([]string{
	".github/\t—\tbaseline", ".githooks/\t—\tbaseline", ".gitattributes\t—\tbaseline", "AGENTS.md\t—\tbaseline",
	"CLAUDE.md\t—\tbaseline", "docs/engineering-discipline.md\t—\tbaseline", "docs/issue-workflow.md\t—\tbaseline",
	"docs/ci/\t—\tbaseline", "docs/tests/\t—\tbaseline", "docs/setup/setup-check.sh\t—\tbaseline",
	"docs/gates.tsv\t—\tcatalog", "go.mod\t—\tcatalog", "docs/setup/\t—\tarchitecture", "docs/facts/\t—\tarchitecture",
	"docs/guardrails.md\tadded lines in section 2\tbaseline",
}, "\n") + "\n"

func read(t *testing.T, data string) []Entry {
	t.Helper()
	reg, err := ReadRegister([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	return reg
}

func TestReadRegister(t *testing.T) {
	reg := read(t, register)
	if len(reg) != 15 || reg[14] != (Entry{Pattern: "docs/guardrails.md", Exception: ExceptionGuardrails}) || reg[0] != (Entry{Pattern: ".github/"}) {
		t.Fatalf("the register: %+v", reg)
	}
	for name, data := range map[string]string{
		"a wrong header":                  strings.Replace(register, "exception", "except", 1),
		"an exception of another path":    strings.Replace(register, "AGENTS.md\t—", "AGENTS.md\tadded lines in section 2", 1),
		"another exception of guardrails": strings.Replace(register, "added lines in section 2", "added lines anywhere", 1),
		"a source of no list":             strings.Replace(register, "go.mod\t—\tcatalog", "go.mod\t—\tthe Operator", 1),
		"a pattern twice":                 register + "AGENTS.md\t—\tcatalog\n",
	} {
		if _, err := ReadRegister([]byte(data)); err == nil {
			t.Errorf("%s: read, want an error", name)
		}
	}
}

func TestMatch(t *testing.T) {
	reg := read(t, register)
	for path, want := range map[string]string{
		".github/workflows/ci.yml":      ".github/",
		"docs/ci/review-record-lint.sh": "docs/ci/",
		"AGENTS.md":                     "AGENTS.md",
		"go.mod":                        "go.mod",
		"docs/guardrails.md":            "docs/guardrails.md",
		"internal/tool/run.sh":          ".sh",
		"build.sh":                      ".sh",
		"docs/setup/setup-check.sh":     "docs/setup/setup-check.sh",
		"docs/cinema.md":                "",
		"docs/ci":                       "",
		"AGENTS.md.bak":                 "",
		"internal/records/session.go":   "",
		"docs/tests.md":                 "",
		"run.shell":                     "",
	} {
		e, ok := Match(reg, path)
		switch {
		case want == "" && ok:
			t.Errorf("%s matches %+v, want no entry", path, e)
		case want == ".sh" && !ok:
			t.Errorf("%s: no match, want the kind .sh", path)
		case want != "" && want != ".sh" && (!ok || e.Pattern != want):
			t.Errorf("%s matches %+v (%v), want %s", path, e, ok, want)
		}
	}
}

// head is a docs/guardrails.md of 9 lines: §1 on lines 1 and 2, §2 on lines 3
// to 7 with a sub-heading on line 5, §3 on lines 8 and 9.
const head = "## 1. Decisions\none\n## 2. Pitfalls\n- a\n### A part\n- b\n- c\n## 3. Validation\nthree\n"

// diff gives a -U0 diff of docs/guardrails.md with the hunks.
func diff(hunks ...string) []byte {
	return []byte("diff --git a/docs/guardrails.md b/docs/guardrails.md\nindex 1111111..2222222 100644\n" +
		"--- a/docs/guardrails.md\n+++ b/docs/guardrails.md\n" + strings.Join(hunks, ""))
}

func TestGuardrailsAdditions(t *testing.T) {
	for name, c := range map[string]struct {
		diff []byte
		want bool
	}{
		"one line added in §2":                     {diff("@@ -5,0 +6 @@\n+- b\n"), true},
		"two hunks in §2, under a sub-heading":     {diff("@@ -3,0 +4 @@\n+- a\n", "@@ -5,0 +6,2 @@\n+- b\n+- c\n"), true},
		"an added line with no newline at its end": {diff("@@ -6,0 +7 @@\n+- c\n\\ No newline at end of file\n"), true},
		"a line added in §1":                       {diff("@@ -1,0 +2 @@\n+one\n"), false},
		"a line added in §3":                       {diff("@@ -8,0 +9 @@\n+three\n"), false},
		"an added line at the line of ## 3.":       {diff("@@ -7,0 +8 @@\n+## 3. Validation\n"), false},
		"an added line that is ## 2.":              {diff("@@ -2,0 +3 @@\n+## 2. Pitfalls\n"), false},
		"a removed line":                           {diff("@@ -4 +3,0 @@\n-- a\n"), false},
		"a changed line":                           {diff("@@ -4 +4 @@\n-- old\n+- a\n"), false},
		"no hunk: a mode change":                   {[]byte("diff --git a/docs/guardrails.md b/docs/guardrails.md\nold mode 100644\nnew mode 100755\n"), false},
		"a binary diff":                            {[]byte("diff --git a/docs/guardrails.md b/docs/guardrails.md\nindex 1..2 100644\nBinary files a/docs/guardrails.md and b/docs/guardrails.md differ\n"), false},
		"text that is no diff":                     {[]byte("hello\n"), false},
		"an empty diff":                            {nil, false},
	} {
		if got := GuardrailsAdditions(c.diff, []byte(head)); got != c.want {
			t.Errorf("%s: %v, want %v", name, got, c.want)
		}
	}
	if GuardrailsAdditions(diff("@@ -5,0 +6 @@\n+- b\n"), []byte("## 1. Decisions\none\n## 3. Validation\n- b\n- b\n- b\n")) {
		t.Error("a head with no ## 2.: true, want false")
	}
}

func TestCheck(t *testing.T) {
	reg := read(t, register)
	in2 := Change{Diff: diff("@@ -5,0 +6 @@\n+- b\n"), Head: []byte(head)}
	in3 := Change{Diff: diff("@@ -8,0 +9 @@\n+three\n"), Head: []byte(head)}
	for name, c := range map[string]struct {
		names        []string
		guardrails   Change
		reason, path string
	}{
		"no rule path":                           {[]string{"internal/x/x.go", "README.md"}, Change{}, "", ""},
		"a workflow before a rule path":          {[]string{"AGENTS.md", ".github/workflows/ci.yml"}, Change{}, "workflow", ".github/workflows/ci.yml"},
		"a rule path":                            {[]string{"internal/x/x.go", "docs/tests/test-levels.md"}, Change{}, "rule-path", "docs/tests/test-levels.md"},
		"a script anywhere":                      {[]string{"tools/gen.sh"}, Change{}, "rule-path", "tools/gen.sh"},
		"additions in §2 of the guardrails":      {[]string{"docs/guardrails.md", "internal/x/x.go"}, in2, "", ""},
		"an addition in §3 of the guardrails":    {[]string{"docs/guardrails.md"}, in3, "rule-path", "docs/guardrails.md"},
		"the guardrails with no diff given":      {[]string{"docs/guardrails.md"}, Change{}, "rule-path", "docs/guardrails.md"},
		"the first rule path in the order given": {[]string{"go.mod", "AGENTS.md"}, Change{}, "rule-path", "go.mod"},
	} {
		if reason, path := Check(c.names, reg, c.guardrails); reason != c.reason || path != c.path {
			t.Errorf("%s: %q %q, want %q %q", name, reason, path, c.reason, c.path)
		}
	}
	// The register's exception is what lets the guardrails pass: with no
	// exception on its row, the same diff is a rule-path change.
	noException := read(t, strings.Replace(register, "docs/guardrails.md\tadded lines in section 2", "docs/guardrails.md\t—", 1))
	if reason, _ := Check([]string{"docs/guardrails.md"}, noException, in2); reason != "rule-path" {
		t.Errorf("the guardrails with no exception in the register: %q, want rule-path", reason)
	}
}
