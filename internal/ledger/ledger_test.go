package ledger

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/tsv"
)

func at(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func count(n int64) *int64 { return &n }

// session is an api session of claude that wrote its first output 2 seconds
// after its start and ended 90 seconds after it.
func session() Session {
	return Session{ID: "S-1a2b3c4d", Task: "T-ab12", Role: "developer", Harness: "claude", Model: "claude-opus-5-5",
		Billing: "api", Start: at("2026-10-09T12:00:00Z"), FirstOutput: at("2026-10-09T12:00:02Z"), End: at("2026-10-09T12:01:30Z")}
}

// observed is a report of three counts and no cost.
func observed() Usage {
	return Usage{In: count(1000000), Out: count(200000), Cache: count(3000000), Status: "observed", Models: []string{"claude-opus-5-5"}}
}

// prices holds the claude-opus-5-5 rows of two dates; the newest is
// 2026-10-01: in 5.0, out 25.0, cache 0.5 USD per million.
func prices() [][]string {
	return [][]string{
		{"P-001", "claude", "claude-opus-5-5", "in", "15.0", "USD", "https://docs.claude.com/pricing", "2026-09-01T00:00:00Z"},
		{"P-002", "claude", "claude-opus-5-5", "out", "75.0", "USD", "https://docs.claude.com/pricing", "2026-09-01T00:00:00Z"},
		{"P-003", "claude", "claude-opus-5-5", "cache", "1.5", "USD", "https://docs.claude.com/pricing", "2026-09-01T00:00:00Z"},
		{"P-004", "claude", "claude-opus-5-5", "in", "5.0", "USD", "https://docs.claude.com/pricing", "2026-10-01T00:00:00Z"},
		{"P-005", "claude", "claude-opus-5-5", "out", "25.0", "USD", "https://docs.claude.com/pricing", "2026-10-01T00:00:00Z"},
		{"P-006", "claude", "claude-opus-5-5", "cache", "0.5", "USD", "https://docs.claude.com/pricing", "2026-10-01T00:00:00Z"},
		{"P-007", "claude", "claude-fable-5-1", "in", "20.0", "USD", "https://docs.claude.com/pricing", "2026-10-01T00:00:00Z"},
	}
}

// readBack writes the row by the block telemetry and reads it with
// records.ReadTelemetry, so each row of a test is one that the record takes.
func readBack(t *testing.T, row []string) map[string]string {
	t.Helper()
	var buf bytes.Buffer
	if err := tsv.Write(&buf, records.TelemetrySchema, [][]string{row}); err != nil {
		t.Fatal(err)
	}
	rows, err := records.ReadTelemetry(buf.Bytes())
	if err != nil {
		t.Fatalf("the row does not read back: %v\n%s", err, buf.String())
	}
	f := map[string]string{}
	for i, c := range records.TelemetrySchema.Columns {
		f[c.Name] = rows[0][i]
	}
	return f
}

func row(t *testing.T, s Session, u Usage, p [][]string) map[string]string {
	t.Helper()
	r, err := Row(s, u, p)
	if err != nil {
		t.Fatal(err)
	}
	return readBack(t, r)
}

func TestTheColumnsOfASession(t *testing.T) {
	f := row(t, session(), observed(), prices())
	for c, want := range map[string]string{"session": "S-1a2b3c4d", "task": "T-ab12", "requirements": "", "role": "developer",
		"harness": "claude", "model": "claude-opus-5-5", "billing": "api", "start": "2026-10-09T12:00:00Z",
		"first_output": "2026-10-09T12:00:02Z", "end": "2026-10-09T12:01:30Z", "latency_s": "2", "duration_s": "90",
		"tokens_in": "1000000", "tokens_out": "200000", "tokens_cache": "3000000", "tokens_status": "observed", "tokens_reason": ""} {
		if f[c] != want {
			t.Errorf("%s = %q, want %q", c, f[c], want)
		}
	}
	s := session()
	s.FirstOutput, s.Requirements = time.Time{}, []string{"REQ-013", "REQ-011"}
	if f := row(t, s, observed(), prices()); f["first_output"] != "" || f["latency_s"] != "" || f["requirements"] != "REQ-013 REQ-011" {
		t.Errorf("no first output and two requirements: %q %q %q", f["first_output"], f["latency_s"], f["requirements"])
	}
	s = session()
	s.Start, s.FirstOutput = at("2026-10-09T12:00:00.700Z"), time.Time{}
	if f := row(t, s, observed(), prices()); f["start"] != "2026-10-09T12:00:00Z" || f["duration_s"] != "90" {
		t.Errorf("a start with a fraction of a second: %q, %q", f["start"], f["duration_s"])
	}
}

func TestTheMoney(t *testing.T) {
	reported := observed()
	reported.Cost, reported.Currency = "6.13", "USD"
	whole := observed()
	whole.Cost, whole.Currency = "6", "USD"
	tiny := observed()
	tiny.Cost, tiny.Currency = "1e-7", "USD"
	subscription := session()
	subscription.Billing = "subscription"
	for name, c := range map[string]struct {
		s                              Session
		u                              Usage
		money, currency, status, price string
	}{
		"reported by an api session":             {session(), reported, "6.13", "USD", "reported", ""},
		"a whole-number cost":                    {session(), whole, "6.0", "USD", "reported", ""},
		"a cost with an exponent":                {session(), tiny, "0.0000001", "USD", "reported", ""},
		"a subscription with a cost is computed": {subscription, reported, "11.5", "USD", "computed", "P-004"},
		// 1000000 × 5.0 + 200000 × 25.0 + 3000000 × 0.5, per million: 5 + 5 + 1.5.
		"computed from the newest rows": {session(), observed(), "11.5", "USD", "computed", "P-004"},
		// A row of the model under another harness ID on that date is not the
		// session's (the plan review of #163, condition 1).
		"a row of another harness on that date": {session(), observed(), "11.5", "USD", "computed", "P-004"},
	} {
		p := prices()
		if name == "a row of another harness on that date" {
			p = append(p, []string{"P-030", "bedrock", "claude-opus-5-5", "out", "30.0", "USD", "https://x.example/p", "2026-10-01T00:00:00Z"})
		}
		f := row(t, c.s, c.u, p)
		if f["money"] != c.money || f["currency"] != c.currency || f["money_status"] != c.status || f["price"] != c.price {
			t.Errorf("%s: %q %q %q %q, want %q %q %q %q", name, f["money"], f["currency"], f["money_status"], f["price"], c.money, c.currency, c.status, c.price)
		}
	}
}

// An exact decimal: a price of six decimal places times a count that a
// float64 would round.
func TestAComputedMoneyIsExact(t *testing.T) {
	u := Usage{In: count(123456789), Out: count(1), Cache: count(0), Status: "observed"}
	p := [][]string{
		{"P-010", "claude", "claude-opus-5-5", "in", "3.141593", "USD", "https://x.example/p", "2026-10-01T00:00:00Z"},
		{"P-011", "claude", "claude-opus-5-5", "out", "0.000001", "USD", "https://x.example/p", "2026-10-01T00:00:00Z"},
		{"P-012", "claude", "claude-opus-5-5", "cache", "0.1", "USD", "https://x.example/p", "2026-10-01T00:00:00Z"},
	}
	// 123456789 × 3.141593 / 10^6 + 1 × 0.000001 / 10^6
	if f := row(t, session(), u, p); f["money"] != "387.850984124878" {
		t.Errorf("money %q, want 387.850984124878", f["money"])
	}
	// 3 × 0.1 per million: 0.0000003 exactly; a float64 gives 3.0000000000000004e-07.
	u = Usage{In: count(3), Out: count(0), Cache: count(0), Status: "observed"}
	p[0][4] = "0.1"
	if f := row(t, session(), u, p); f["money"] != "0.0000003" {
		t.Errorf("money %q, want 0.0000003", f["money"])
	}
}

func TestEachReasonOfAnUnknownMoney(t *testing.T) {
	without := func(class string) [][]string {
		var out [][]string
		for _, r := range prices() {
			if !(r[3] == class && r[7] == "2026-10-01T00:00:00Z") {
				out = append(out, r)
			}
		}
		return out
	}
	twice := append(prices(), []string{"P-020", "claude", "claude-opus-5-5", "out", "30.0", "USD", "https://x.example/p", "2026-10-01T00:00:00Z"})
	euro := prices()
	euro[5] = []string{"P-006", "claude", "claude-opus-5-5", "cache", "0.5", "EUR", "https://x.example/p", "2026-10-01T00:00:00Z"}
	otherHarness := [][]string{}
	for _, r := range prices() {
		r = append([]string{}, r...)
		r[1] = "bedrock"
		otherHarness = append(otherHarness, r)
	}
	twoModels := observed()
	twoModels.Models = []string{"claude-opus-5-5", "claude-fable-5-1"}
	otherModel := observed()
	otherModel.Models = []string{"claude-fable-5-1"}
	partial := observed()
	partial.Cache, partial.Status, partial.Reason = nil, "partial", "a model of the report has no cache count"
	unavailable := Usage{Status: "unavailable", Reason: "the harness reports none"}
	sub := session()
	sub.Billing = "subscription"
	for name, c := range map[string]struct {
		s Session
		u Usage
		p [][]string
	}{
		"a class with no row of the newest date":           {session(), observed(), without("cache")},
		"a class with two rows of the newest date":         {session(), observed(), twice},
		"rows of two currencies":                           {session(), observed(), euro},
		"rows of the model under another harness only":     {session(), observed(), otherHarness},
		"a report that names two models":                   {session(), twoModels, prices()},
		"a report that names one model, not the session's": {session(), otherModel, prices()},
		"partial tokens":                                   {session(), partial, prices()},
		"no tokens, a subscription":                        {sub, unavailable, prices()},
		"no prices":                                        {session(), observed(), nil},
	} {
		f := row(t, c.s, c.u, c.p)
		if f["money_status"] != "unknown" || f["money"] != "" || f["currency"] != "" || f["price"] != "" {
			t.Errorf("%s: %q %q %q %q, want unknown and no money, currency or price", name, f["money_status"], f["money"], f["currency"], f["price"])
		}
	}
}

func TestARowThatTheRecordRefusesIsAnError(t *testing.T) {
	s := session()
	s.End = at("2026-10-09T11:59:00Z")
	if _, err := Row(s, observed(), prices()); err == nil {
		t.Error("an end before the start: a row, want an error")
	}
	u := observed()
	u.Cost, u.Currency = "six", "USD"
	if _, err := Row(session(), u, prices()); err == nil || !strings.Contains(err.Error(), "six") {
		t.Errorf("a cost that is not a number: %v, want an error that names it", err)
	}
}
