package run

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
)

// approvers is approvers.tsv of a Start: one account in two roles, and an
// approver.
var approvers = [][]string{
	{"101", "operator", "pharzam", "2026-10-08T12:00:00Z", "start"},
	{"101", "idea-owner", "pharzam", "2026-10-08T12:00:00Z", "start"},
	{"202", "approver", "bob", "2026-10-08T13:00:00Z", "6040085862"},
	{"303", "idea-owner", "carol", "2026-10-08T13:00:00Z", "6040085863"},
}

func copyRow(comment, seen, author, app string) []string {
	return []string{comment, seen, "2", author, "x", app, "2026-10-08T12:00:00Z", "2026-10-08T12:00:05Z", sha256Hex("b"), "copies/" + comment + "-" + seen + ".md"}
}

func sha256Hex(s string) string { sum := sha256.Sum256([]byte(s)); return hex.EncodeToString(sum[:]) }

func TestADecisionIsAFirstCopyByAnApproverInTheRole(t *testing.T) {
	for _, c := range []struct {
		name  string
		row   []string
		roles []string
		want  bool
	}{
		{"the operator in Intake", copyRow("901", "1", "101", "—"), IntakeRoles, true},
		{"the idea owner in an acceptance", copyRow("901", "1", "303", "—"), AcceptanceRoles, true},
		{"the operator's account in an acceptance, as idea owner too", copyRow("901", "1", "101", "—"), AcceptanceRoles, true},
		{"an approver in Intake", copyRow("901", "1", "202", "—"), IntakeRoles, false},
		{"an ID not in approvers.tsv", copyRow("901", "1", "404", "—"), IntakeRoles, false},
		{"a comment that an App made", copyRow("901", "1", "101", "layup-watch"), IntakeRoles, false},
		{"an edit (seen 2)", copyRow("901", "2", "101", "—"), IntakeRoles, false},
		{"a comment edited before its first copy is its first copy", append(copyRow("901", "1", "101", "—")[:6:6], "2026-10-08T12:30:00Z", "2026-10-08T12:30:05Z", sha256Hex("b"), "copies/901-1.md"), IntakeRoles, true},
	} {
		if got := Decision(c.row, approvers, c.roles...); got != c.want {
			t.Errorf("%s: Decision = %v, want %v", c.name, got, c.want)
		}
	}
	if !reflect.DeepEqual(IntakeRoles, []string{"operator", "idea-owner"}) || !reflect.DeepEqual(AcceptanceRoles, []string{"idea-owner"}) {
		t.Errorf("the roles: Intake %v, acceptance %v", IntakeRoles, AcceptanceRoles)
	}
}

func comment(body string, edited time.Time) forge.Comment {
	created := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	return forge.Comment{ID: 901, AuthorID: 101, AuthorLogin: "pharzam", Created: created, Updated: edited, Body: body}
}

func TestCopyBeforeRead(t *testing.T) {
	now := time.Date(2026, 10, 8, 12, 0, 5, 0, time.UTC)
	first := comment("the answer\n", time.Time{}) // a stand-in forge that leaves Updated zero
	row, body, ok := Copy(nil, first, 2, now)
	want := []string{"901", "1", "2", "101", "pharzam", "—", "2026-10-08T12:00:00Z", "2026-10-08T12:00:05Z", sha256Hex("the answer\n"), "copies/901-1.md"}
	if !ok || !reflect.DeepEqual(row, want) || string(body) != "the answer\n" {
		t.Fatalf("the first copy = %q, %q, %v\nwant %q", row, body, ok, want)
	}
	if err := records.CheckCopy(row); err != nil {
		t.Errorf("the row breaks the block copies: %v", err)
	}
	table := [][]string{row}
	if _, _, ok := Copy(table, first, 2, now.Add(time.Minute)); ok {
		t.Error("an unchanged comment got a new row")
	}
	edited := comment("the answer, changed\n", time.Date(2026, 10, 8, 12, 10, 0, 0, time.UTC))
	edited.App = "layup-watch"
	row2, body2, ok := Copy(table, edited, 2, now.Add(time.Hour))
	want2 := []string{"901", "2", "2", "101", "pharzam", "layup-watch", "2026-10-08T12:10:00Z", "2026-10-08T13:00:05Z", sha256Hex("the answer, changed\n"), "copies/901-2.md"}
	if !ok || !reflect.DeepEqual(row2, want2) || string(body2) != "the answer, changed\n" {
		t.Fatalf("the copy of an edit = %q, %v\nwant %q", row2, ok, want2)
	}
	if !reflect.DeepEqual(table[0], want) {
		t.Errorf("the first copy changed: %q", table[0])
	}
	if err := records.CheckCopies(append(table, row2)); err != nil {
		t.Errorf("the table breaks the block copies: %v", err)
	}
	other := forge.Comment{ID: 902, AuthorID: 202, AuthorLogin: "bob", Created: now, Body: "x"}
	if row3, _, ok := Copy(append(table, row2), other, 2, now); !ok || row3[1] != "1" {
		t.Errorf("a new comment = %q, %v; want seen 1", row3, ok)
	}
}
