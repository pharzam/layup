package records

import (
	"strings"
	"testing"
)

// table writes a record of the schema's header and the rows, a field per tab;
// "—" is the empty value.
func table(header string, rows ...string) []byte {
	return []byte(strings.Join(append([]string{header}, rows...), "\n") + "\n")
}

const (
	sha1A   = "0123456789abcdef0123456789abcdef01234567"
	sha256A = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
)

// startRows is a valid start.tsv for the harness register claude, devin.
func startRows() []string {
	return []string{
		"layup.version\tv0.1.0\trun",
		"psb.sha256\t" + sha256A + "\tcommand",
		"vision.sha256\t—\tcommand",
		"forge.plan\tfree\tcommand",
		"forge.visibility\tpublic\tforge",
		"app.permissions\tcontents:write issues:write metadata:read\tforge",
		"operator.id\t123\tforge",
		"idea-owner.id\t123\tforge",
		"intake.cap\t50.0,8.0\tcommand",
		"lease.H\t5\tcommand",
		"watch.T\t10\tcommand",
		"harness.claude.cap\t10.0\tregister",
		"harness.claude.wall\t60\tregister",
		"harness.devin.cap\t—\tregister",
		"harness.devin.wall\t30\tregister",
		"pin.source\thttps://github.com/pharzam/armature\trun",
		"pin.commit\t" + sha1A + "\trun",
		"pin.tree\t" + sha1A + "\trun",
		"pin.time\t2026-10-07T12:00:00Z\trun",
		"issue.intake\t1\tforge",
		"issue.control\topening\tforge",
		"watch\t—\trun",
	}
}

const startHeader = "name\tvalue\tsource"

// set gives startRows with the row of name replaced by line.
func set(name, line string) []string {
	rows := startRows()
	for i, r := range rows {
		if strings.HasPrefix(r, name+"\t") {
			rows[i] = line
		}
	}
	return rows
}

// without gives startRows without the row of name.
func without(name string) []string {
	var rows []string
	for _, r := range startRows() {
		if !strings.HasPrefix(r, name+"\t") {
			rows = append(rows, r)
		}
	}
	return rows
}

func TestAValidStartIsRead(t *testing.T) {
	rows, ids, err := ReadStartAsWritten(table(startHeader, startRows()...))
	if err != nil || len(rows) != len(startRows()) || strings.Join(ids, " ") != "claude devin" {
		t.Fatalf("%d rows, the harnesses %q, %v; want each row and claude devin", len(rows), ids, err)
	}
	noHarness := append(append([]string{}, startRows()[:11]...), startRows()[15:]...)
	if _, ids, err := ReadStartAsWritten(table(startHeader, noHarness...)); err != nil || len(ids) != 0 {
		t.Errorf("no harness row: the harnesses %q, %v; want none", ids, err)
	}
	// The harnesses are those of the rows, in their order, whatever the
	// register is now (#199).
	swapped := startRows()
	swapped[11], swapped[12], swapped[13], swapped[14] = swapped[13], swapped[14], swapped[11], swapped[12]
	if _, ids, err := ReadStartAsWritten(table(startHeader, swapped...)); err != nil || strings.Join(ids, " ") != "devin claude" {
		t.Errorf("devin before claude: the harnesses %q, %v; want devin claude", ids, err)
	}
}

func TestStartRefusesEachBrokenRule(t *testing.T) {
	order := startRows()
	order[0], order[1] = order[1], order[0]
	wallFirst := startRows()
	wallFirst[11], wallFirst[12] = wallFirst[12], wallFirst[11]
	// renamed gives startRows with the two rows of claude under the ID id, so
	// that only the form of the ID is wrong.
	renamed := func(id string) []string {
		rows := startRows()
		rows[11] = "harness." + id + ".cap\t10.0\tregister"
		rows[12] = "harness." + id + ".wall\t60\tregister"
		return rows
	}
	for _, c := range []struct {
		name string
		rows []string
	}{
		{"an unknown name", append(startRows(), "colour\tblue\trun")},
		{"a missing name", without("pin.tree")},
		{"the rows out of order", order},
		{"a harness with a wall row and no cap row", without("harness.devin.cap")},
		{"a harness with a cap row and no wall row", without("harness.devin.wall")},
		{"wall before cap", wallFirst},
		{"a harness ID that is not of the form <word>", renamed("Claude")},
		{"an empty harness ID", renamed("")},
		{"the empty value where the block does not allow it", set("layup.version", "layup.version\t—\trun")},
		{"a SHA-256 that is not one", set("psb.sha256", "psb.sha256\tabc\tcommand")},
		{"a SHA-1 that is not one", set("pin.commit", "pin.commit\t"+sha256A+"\trun")},
		{"a time that is not one", set("pin.time", "pin.time\t2026-10-07 12:00\trun")},
		{"an ID that is not an int", set("operator.id", "operator.id\tpharzam\tforge")},
		{"lease.H that is not an int", set("lease.H", "lease.H\t5.0\tcommand")},
		{"a wall that is not an int", set("harness.claude.wall", "harness.claude.wall\tan hour\tregister")},
		{"a cap that is not a decimal", set("harness.claude.cap", "harness.claude.cap\t10\tregister")},
		{"an intake cap that is not two decimals", set("intake.cap", "intake.cap\t50.0\tcommand")},
		{"a permission that is not name:level", set("app.permissions", "app.permissions\tcontents issues:write\tforge")},
		{"an issue that is neither an int nor opening", set("issue.intake", "issue.intake\tsoon\tforge")},
		{"a watch that is neither confirmed nor not-confirmed", set("watch", "watch\tyes\trun")},
		{"the wrong source", set("pin.source", "pin.source\thttps://x\tcommand")},
	} {
		if _, _, err := ReadStartAsWritten(table(startHeader, c.rows...)); err == nil {
			t.Errorf("%s: read, want an error", c.name)
		}
	}
}

const approversHeader = "id\trole\tlogin\tsince\tsource"

func TestApproversRefusesEachBrokenRule(t *testing.T) {
	ok := []string{"123\toperator\tpharzam\t2026-10-07T12:00:00Z\tstart", "123\tidea-owner\tpharzam\t2026-10-07T12:00:00Z\tstart", "456\tapprover\tbob\t2026-10-08T09:00:00Z\t6040085862"}
	if _, err := ReadApprovers(table(approversHeader, ok...)); err != nil {
		t.Fatal(err)
	}
	for name, row := range map[string]string{
		"an empty login":                      "123\toperator\t—\t2026-10-07T12:00:00Z\tstart",
		"an empty since":                      "123\toperator\tpharzam\t—\tstart",
		"a source that is not start or an ID": "123\toperator\tpharzam\t2026-10-07T12:00:00Z\tthe Operator",
		"an empty source":                     "123\toperator\tpharzam\t2026-10-07T12:00:00Z\t—",
	} {
		if _, err := ReadApprovers(table(approversHeader, row)); err == nil {
			t.Errorf("%s: read, want an error", name)
		}
	}
}

const leaseHeader = "run\thost\tversion\tstarted\theartbeat\tstate"

func TestLeaseRefusesEachBrokenRule(t *testing.T) {
	ok := "0123456789abcdef\thost-1\tv0.1.0\t2026-10-07T12:00:00Z\t0\theld"
	if _, err := ReadLease(table(leaseHeader, ok)); err != nil {
		t.Fatal(err)
	}
	for name, rows := range map[string][]string{
		"no row":                    nil,
		"two rows":                  {ok, "fedcba9876543210\thost-2\tv0.1.0\t2026-10-07T13:00:00Z\t0\theld"},
		"a run ID of 15 characters": {"0123456789abcde\thost-1\tv0.1.0\t2026-10-07T12:00:00Z\t0\theld"},
		"a run ID in capitals":      {"0123456789ABCDEF\thost-1\tv0.1.0\t2026-10-07T12:00:00Z\t0\theld"},
		"an empty host":             {"0123456789abcdef\t—\tv0.1.0\t2026-10-07T12:00:00Z\t0\theld"},
		"an empty version":          {"0123456789abcdef\thost-1\t—\t2026-10-07T12:00:00Z\t0\theld"},
		"an empty started":          {"0123456789abcdef\thost-1\tv0.1.0\t—\t0\theld"},
		"an empty heartbeat":        {"0123456789abcdef\thost-1\tv0.1.0\t2026-10-07T12:00:00Z\t—\theld"},
		"an empty state":            {"0123456789abcdef\thost-1\tv0.1.0\t2026-10-07T12:00:00Z\t0\t—"},
	} {
		if _, err := ReadLease(table(leaseHeader, rows...)); err == nil {
			t.Errorf("%s: read, want an error", name)
		}
	}
}

const copiesHeader = "comment\tseen\tissue\tauthor_id\tauthor_login\tapp\tcreated\tcopied\tsha256\tbody"

// copyRow is a valid row of copies.tsv for one comment and one seen.
func copyRow(comment, seen, app string) string {
	return comment + "\t" + seen + "\t2\t123\tpharzam\t" + app + "\t2026-10-07T12:00:00Z\t2026-10-07T12:00:05Z\t" + sha256A + "\tcopies/" + comment + "-" + seen + ".md"
}

func TestCopiesRefusesEachBrokenRule(t *testing.T) {
	ok := []string{copyRow("901", "1", "—"), copyRow("901", "2", "—"), copyRow("902", "1", "layup-watch")}
	if _, err := ReadCopies(table(copiesHeader, ok...)); err != nil {
		t.Fatal(err)
	}
	for name, rows := range map[string][]string{
		"seen does not start at 1": {copyRow("901", "2", "—")},
		"a gap in seen":            {copyRow("901", "1", "—"), copyRow("901", "3", "—")},
		"seen 0":                   {copyRow("901", "0", "—")},
		"a body at another path":   {strings.Replace(copyRow("901", "1", "—"), "copies/901-1.md", "copies/901.md", 1)},
		"an empty author login":    {strings.Replace(copyRow("901", "1", "—"), "\tpharzam\t", "\t—\t", 1)},
		"an empty sha256":          {strings.Replace(copyRow("901", "1", "—"), sha256A, "—", 1)},
	} {
		if _, err := ReadCopies(table(copiesHeader, rows...)); err == nil {
			t.Errorf("%s: read, want an error", name)
		}
	}
}

func TestCheckValueGivesTheFormOfAName(t *testing.T) {
	for _, c := range []struct {
		name, value string
		ok          bool
	}{
		{"intake.cap", "50.0,8.0", true}, {"intake.cap", "50,8", false}, {"intake.cap", "50.0", false},
		{"lease.H", "5", true}, {"lease.H", "5.0", false}, {"watch.T", "10", true}, {"watch.T", "-1", false},
		{"pin.commit", sha1A, true}, {"no.such.name", "x", false},
	} {
		if err := CheckValue(c.name, c.value); (err == nil) != c.ok {
			t.Errorf("CheckValue(%s, %q) = %v; want ok %v", c.name, c.value, err, c.ok)
		}
	}
}
