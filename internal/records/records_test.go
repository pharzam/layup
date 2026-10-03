package records

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// dash is the empty value of a field, from its code point.
var dash = string(rune(0x2014))

// fileOf gives the text of a record of schema s with the rows, each a map of
// its columns; a column that a row does not name is the empty value.
func fileOf(s tsv.Schema, rows ...map[string]string) string {
	var names []string
	for _, c := range s.Columns {
		names = append(names, c.Name)
	}
	text := strings.Join(names, "\t") + "\n"
	for _, r := range rows {
		var fields []string
		for _, n := range names {
			v, ok := r[n]
			if !ok || v == "" {
				v = dash
			}
			fields = append(fields, v)
		}
		text += strings.Join(fields, "\t") + "\n"
	}
	return text
}

// session gives a complete observed row of telemetry.tsv, a reported cost,
// with the changes; "" is the empty value.
func session(change map[string]string) map[string]string {
	r := map[string]string{"session": "S-0123abcd", "task": "T-tmhw", "requirements": "REQ-011", "role": "builder", "harness": "claude-code",
		"model": "claude-opus-5-5", "billing": "api", "start": "2026-10-03T00:00:00Z", "first_output": "2026-10-03T00:00:05Z",
		"end": "2026-10-03T00:10:00Z", "latency_s": "5", "duration_s": "600", "tokens_in": "1000", "tokens_out": "200",
		"tokens_cache": "3000", "tokens_status": "observed", "tokens_reason": "", "money": "0.42", "currency": "USD",
		"money_status": "reported", "price": ""}
	for k, v := range change {
		r[k] = v
	}
	return r
}

// A row of each status passes: observed, partial and unavailable tokens; a
// reported, a computed and an unknown cost; a session with no first output;
// a subscription, computed or unknown (D2 of #94).
func TestTelemetryRowsThatPass(t *testing.T) {
	for name, change := range map[string]map[string]string{
		"observed and reported":   nil,
		"partial":                 {"tokens_out": "", "tokens_cache": "", "tokens_status": "partial", "tokens_reason": "the harness reports input tokens only"},
		"unavailable and unknown": {"tokens_in": "", "tokens_out": "", "tokens_cache": "", "tokens_status": "unavailable", "tokens_reason": "the harness reports no tokens", "money": "", "currency": "", "money_status": "unknown"},
		"computed":                {"money_status": "computed", "price": "P-001"},
		"no first output":         {"first_output": "", "latency_s": ""},
		"a subscription":          {"billing": "subscription", "money_status": "computed", "price": "P-002"},
		"a subscription, unknown": {"billing": "subscription", "money": "", "currency": "", "money_status": "unknown"},
	} {
		rows, err := ReadTelemetry([]byte(fileOf(TelemetrySchema, session(change))))
		if err != nil || len(rows) != 1 {
			t.Errorf("%s: %d rows, %v; want the row", name, len(rows), err)
		}
	}
}

// Each broken row rule is an error that names the line and the column (D2 and
// D3 of #94).
func TestTelemetryRowsThatBreakARule(t *testing.T) {
	for _, c := range []struct {
		name   string
		change map[string]string
		column string
	}{
		{"a session ID that is not hexadecimal", map[string]string{"session": "S-0123wxyz"}, "session"},
		{"a latency with no first output", map[string]string{"first_output": ""}, "latency_s"},
		{"a first output with no latency", map[string]string{"latency_s": ""}, "latency_s"},
		{"a latency that is not first_output - start", map[string]string{"latency_s": "6"}, "latency_s"},
		{"a duration that is not end - start", map[string]string{"duration_s": "601"}, "duration_s"},
		{"an end before the start", map[string]string{"end": "2026-10-02T23:59:00Z", "first_output": "", "latency_s": ""}, "end"},
		{"a first output after the end", map[string]string{"first_output": "2026-10-03T00:11:00Z", "latency_s": "660"}, "first_output"},
		{"observed with a token column missing", map[string]string{"tokens_cache": ""}, "tokens_status"},
		{"unavailable with a token column", map[string]string{"tokens_out": "", "tokens_cache": "", "tokens_status": "unavailable", "tokens_reason": "none"}, "tokens_status"},
		{"partial with the three token columns", map[string]string{"tokens_status": "partial", "tokens_reason": "some"}, "tokens_status"},
		{"partial with no reason", map[string]string{"tokens_cache": "", "tokens_status": "partial"}, "tokens_reason"},
		{"observed with a reason", map[string]string{"tokens_reason": "all of them"}, "tokens_reason"},
		{"a cost of 0 that is unknown", map[string]string{"money": "0.00", "money_status": "unknown"}, "money"},
		{"a reported cost with no money", map[string]string{"money": ""}, "money"},
		{"an unknown cost with a currency", map[string]string{"money": "", "money_status": "unknown"}, "currency"},
		{"a currency that is not ISO 4217", map[string]string{"currency": "usd"}, "currency"},
		{"a computed cost with no price row", map[string]string{"money_status": "computed"}, "price"},
		{"a price row of a reported cost", map[string]string{"price": "P-001"}, "price"},
		{"a price that is not a row ID", map[string]string{"money_status": "computed", "price": "the list price"}, "price"},
		{"a subscription that is reported", map[string]string{"billing": "subscription"}, "money_status"},
		// Round 1 of #94: a row with no start and no end, and each column whose
		// block rule has no clause for the empty value.
		{"no start and no end", map[string]string{"start": "", "end": "", "first_output": "", "latency_s": "", "duration_s": "0"}, "start"},
		{"the empty value of task", map[string]string{"task": ""}, "task"},
		{"the empty value of role", map[string]string{"role": ""}, "role"},
		{"the empty value of harness", map[string]string{"harness": ""}, "harness"},
		{"the empty value of model", map[string]string{"model": ""}, "model"},
		{"the empty value of billing", map[string]string{"billing": ""}, "billing"},
		{"the empty value of end", map[string]string{"end": ""}, "end"},
		{"the empty value of duration_s", map[string]string{"duration_s": ""}, "duration_s"},
		{"the empty value of tokens_status", map[string]string{"tokens_status": "", "tokens_reason": "a reason"}, "tokens_status"},
		{"the empty value of money_status", map[string]string{"money_status": ""}, "money_status"},
	} {
		broken := session(map[string]string{"session": "S-89abcdef"})
		for k, v := range c.change {
			broken[k] = v
		}
		_, err := ReadTelemetry([]byte(fileOf(TelemetrySchema, session(nil), broken)))
		var e *tsv.Error
		if !errors.As(err, &e) || e.Line != 3 || e.Column != c.column {
			t.Errorf("%s: %v; want an error of line 3, column %q", c.name, err, c.column)
		}
	}
	// The row function that a writer calls before it writes a row (note 4 of
	// the plan review of #94).
	var rows [][]string
	for _, change := range []map[string]string{nil, {"money_status": "computed"}} {
		r := session(change)
		var row []string
		for _, col := range TelemetrySchema.Columns {
			row = append(row, r[col.Name])
		}
		rows = append(rows, row)
	}
	var re *RowError
	if err := CheckTelemetry(rows[0]); err != nil {
		t.Errorf("CheckTelemetry of a good row: %v", err)
	}
	// A time with an offset, which the type time refuses, is not a time of a
	// row either (note 3 of round 1).
	offset := slices.Clone(rows[0])
	offset[slices.IndexFunc(TelemetrySchema.Columns, func(c tsv.Column) bool { return c.Name == "start" })] = "2026-10-03T01:00:00+01:00"
	if err := CheckTelemetry(offset); !errors.As(err, &re) || re.Column != "start" {
		t.Errorf("CheckTelemetry of a start with an offset: %v; want a *RowError of the column start", err)
	}
	if err := CheckTelemetry(rows[1]); !errors.As(err, &re) || re.Column != "price" {
		t.Errorf("CheckTelemetry of a computed row with no price: %v; want a *RowError of the column price", err)
	}
}

// price gives a row of prices.tsv, with the changes.
func price(change map[string]string) map[string]string {
	r := map[string]string{"id": "P-001", "harness": "claude-code", "model": "claude-opus-5-5", "class": "in", "price": "15.00", "currency": "USD",
		"source": "https://example.invalid/pricing", "date": "2026-10-01T00:00:00Z"}
	for k, v := range change {
		r[k] = v
	}
	return r
}

// A row of prices.tsv passes with an http or an https source; a bad ID, a
// class out of in, out and cache, a price with an exponent, a currency that is
// not ISO 4217 and a source that is not an http or https URL are refused, each
// with its line and its column (D2 of #94, note 3 of its plan review).
func TestPriceRows(t *testing.T) {
	rows, err := ReadPrices([]byte(fileOf(PricesSchema, price(nil), price(map[string]string{"id": "P-002", "source": "http://example.invalid/pricing"}))))
	if err != nil || len(rows) != 2 {
		t.Errorf("two good rows: %d rows, %v", len(rows), err)
	}
	for _, c := range []struct {
		name   string
		change map[string]string
		column string
	}{
		{"a bad ID", map[string]string{"id": "P-1"}, "id"},
		{"a class out of in, out and cache", map[string]string{"class": "total"}, "class"},
		{"a price with an exponent", map[string]string{"price": "1e3"}, "price"},
		{"a currency in lower case", map[string]string{"currency": "usd"}, "currency"},
		{"a source of another scheme", map[string]string{"source": "ftp://example.invalid/pricing"}, "source"},
		{"a source that is not a URL", map[string]string{"source": "the pricing page"}, "source"},
		// Round 1 of #94: no column of prices.tsv holds the empty value.
		{"the empty value of harness", map[string]string{"harness": ""}, "harness"},
		{"the empty value of model", map[string]string{"model": ""}, "model"},
		{"the empty value of class", map[string]string{"class": ""}, "class"},
		{"the empty value of price", map[string]string{"price": ""}, "price"},
		{"the empty value of currency", map[string]string{"currency": ""}, "currency"},
		{"the empty value of source", map[string]string{"source": ""}, "source"},
		{"the empty value of date", map[string]string{"date": ""}, "date"},
	} {
		_, err := ReadPrices([]byte(fileOf(PricesSchema, price(map[string]string{"id": "P-009"}), price(c.change))))
		var e *tsv.Error
		if !errors.As(err, &e) || e.Line != 3 || e.Column != c.column {
			t.Errorf("%s: %v; want an error of line 3, column %q", c.name, err, c.column)
		}
	}
}
