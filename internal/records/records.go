// Package records holds the schemas and the row rules of the record kinds that
// phase 1 defines and later phases write (docs/spec/records.md): the
// telemetry record and the price table of REQ-011 (task T-tmhw, #94). The one
// writer of the records branch is phase 2's.
package records

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/pharzam/layup/internal/tsv"
)

// TelemetrySchema is the form of telemetry.tsv: the block telemetry of
// docs/spec/records.md.
var TelemetrySchema = tsv.Schema{Name: "telemetry", Location: "records:telemetry.tsv", Columns: []tsv.Column{
	{Name: "session", Type: "id(S-xxxxxxxx)", Key: true},
	{Name: "task", Type: "id(T-xxxx)"},
	{Name: "requirements", Type: "list(text)"},
	{Name: "role", Type: "text"},
	{Name: "harness", Type: "text"},
	{Name: "model", Type: "text"},
	{Name: "billing", Type: "enum(api|subscription)"},
	{Name: "start", Type: "time"},
	{Name: "first_output", Type: "time"},
	{Name: "end", Type: "time"},
	{Name: "latency_s", Type: "int"},
	{Name: "duration_s", Type: "int"},
	{Name: "tokens_in", Type: "int"},
	{Name: "tokens_out", Type: "int"},
	{Name: "tokens_cache", Type: "int"},
	{Name: "tokens_status", Type: "enum(observed|partial|unavailable)"},
	{Name: "tokens_reason", Type: "text"},
	{Name: "money", Type: "decimal"},
	{Name: "currency", Type: "text"},
	{Name: "money_status", Type: "enum(reported|computed|unknown)"},
	{Name: "price", Type: "text"},
}}

// PricesSchema is the form of prices.tsv on the host: the block prices.
var PricesSchema = tsv.Schema{Name: "prices", Location: "host:prices.tsv", Columns: []tsv.Column{
	{Name: "id", Type: "id(P-NNN)", Key: true},
	{Name: "harness", Type: "text"},
	{Name: "model", Type: "text"},
	{Name: "class", Type: "enum(in|out|cache)"},
	{Name: "price", Type: "decimal"},
	{Name: "currency", Type: "text"},
	{Name: "source", Type: "text"},
	{Name: "date", Type: "time"},
}}

// A RowError is a row that breaks a rule of its block: the column and why.
type RowError struct{ Column, Reason string }

func (e *RowError) Error() string { return "column " + e.Column + ": " + e.Reason }

// The forms that the rules of the blocks give in words, not as a type.
var (
	sessionForm  = regexp.MustCompile(`^S-[0-9a-f]{8}$`)
	currencyForm = regexp.MustCompile(`^[A-Z]{3}$`)
	priceIDForm  = regexp.MustCompile(`^P-([0-9]{3}|[1-9][0-9]{3,})$`)
)

// fields gives the value of each column of a row of s by its name, as
// tsv.Read gives it ("" for the empty value).
func fields(s tsv.Schema, r []string) map[string]string {
	m := map[string]string{}
	for i, c := range s.Columns {
		if i < len(r) {
			m[c.Name] = r[i]
		}
	}
	return m
}

// timeForm is the form of the type time of internal/tsv: UTC, to the second.
const timeForm = "2006-01-02T15:04:05Z"

// seconds gives to - from in whole seconds, and whether both are times of the
// form of the type time.
func seconds(from, to string) (int, bool) {
	a, err := time.Parse(timeForm, from)
	b, err2 := time.Parse(timeForm, to)
	return int(b.Sub(a) / time.Second), err == nil && err2 == nil && a.Format(timeForm) == from && b.Format(timeForm) == to
}

// The columns whose block rule has no clause for the empty value, so they
// never hold it (round 1 of #94; docs/spec/README.md: the owner package checks
// a column whose rule forbids it). requirements may be empty: an empty list.
var (
	telemetryRequired = []string{"task", "role", "harness", "model", "billing", "start", "end", "duration_s", "tokens_status", "money_status"}
	pricesRequired    = []string{"harness", "model", "class", "price", "currency", "source", "date"}
)

// missing gives the first column of names whose value is the empty value.
func missing(f map[string]string, names []string) string {
	for _, n := range names {
		if f[n] == "" {
			return n
		}
	}
	return ""
}

// CheckTelemetry checks the row rules of the block telemetry on one row, as
// tsv.Read gives it, that the types do not hold (D2 and D3 of #94): a writer
// calls it before it writes a row. A broken rule is a *RowError.
func CheckTelemetry(r []string) error {
	f := fields(TelemetrySchema, r)
	bad := func(column, reason string) error { return &RowError{column, reason} }
	tokens := 0
	for _, c := range []string{"tokens_in", "tokens_out", "tokens_cache"} {
		if f[c] != "" {
			tokens++
		}
	}
	if c := missing(f, telemetryRequired); c != "" {
		return bad(c, "this column never holds the empty value")
	}
	duration, times := seconds(f["start"], f["end"])
	latency, firstOutput := seconds(f["start"], f["first_output"])
	afterEnd, _ := seconds(f["first_output"], f["end"])
	switch {
	case !times:
		return bad("start", "start and end are times of the form "+timeForm)
	case f["first_output"] != "" && !firstOutput:
		return bad("first_output", "first_output is a time of the form "+timeForm)
	case !sessionForm.MatchString(f["session"]):
		return bad("session", "a session ID is S- and 8 lowercase hexadecimal characters")
	case (f["first_output"] == "") != (f["latency_s"] == ""):
		return bad("latency_s", "latency_s is the empty value exactly when first_output is")
	case duration < 0:
		return bad("end", "the end is before the start")
	case f["first_output"] != "" && (latency < 0 || afterEnd < 0):
		return bad("first_output", "the first output is not between the start and the end")
	case f["duration_s"] != strconv.Itoa(duration):
		return bad("duration_s", fmt.Sprintf("duration_s is end - start, %d seconds", duration))
	case f["first_output"] != "" && f["latency_s"] != strconv.Itoa(latency):
		return bad("latency_s", fmt.Sprintf("latency_s is first_output - start, %d seconds", latency))
	case f["tokens_status"] == "observed" && tokens != 3:
		return bad("tokens_status", "observed has tokens_in, tokens_out and tokens_cache")
	case f["tokens_status"] == "unavailable" && tokens != 0:
		return bad("tokens_status", "unavailable has no token column")
	case f["tokens_status"] == "partial" && (tokens == 0 || tokens == 3):
		return bad("tokens_status", "partial has at least one token column, and not all three")
	case (f["tokens_status"] == "observed") != (f["tokens_reason"] == ""):
		return bad("tokens_reason", "tokens_reason is the empty value exactly when tokens_status is observed")
	case (f["money_status"] == "unknown") != (f["money"] == ""):
		return bad("money", "money is the empty value exactly when money_status is unknown, so an unknown cost is never 0")
	case (f["money_status"] == "unknown") != (f["currency"] == ""):
		return bad("currency", "currency is the empty value exactly when money_status is unknown")
	case f["currency"] != "" && !currencyForm.MatchString(f["currency"]):
		return bad("currency", "a currency is three capital letters, the form of ISO 4217")
	case f["money_status"] == "computed" && !priceIDForm.MatchString(f["price"]):
		return bad("price", "a computed cost names its row of prices.tsv, P-NNN")
	case f["money_status"] != "computed" && f["price"] != "":
		return bad("price", "price is the empty value unless money_status is computed")
	case f["billing"] == "subscription" && f["money_status"] == "reported":
		return bad("money_status", "a subscription session is computed from the list price, or unknown; never reported")
	}
	return nil
}

// CheckPrice checks the row rules of the block prices on one row that the
// types do not hold: the currency, and a source that is an http or https URL.
func CheckPrice(r []string) error {
	f := fields(PricesSchema, r)
	if c := missing(f, pricesRequired); c != "" {
		return &RowError{c, "no column of prices.tsv holds the empty value"}
	}
	if !currencyForm.MatchString(f["currency"]) {
		return &RowError{"currency", "a currency is three capital letters, the form of ISO 4217"}
	}
	if u, err := url.Parse(f["source"]); err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return &RowError{"source", "the source is the http or https URL of the price list"}
	}
	return nil
}

// read reads data by the schema s, then checks each row with check; a broken
// rule is a *tsv.Error with its line and its column.
func read(data []byte, s tsv.Schema, check func([]string) error) ([][]string, error) {
	rows, err := tsv.Read(data, s)
	if err != nil {
		return nil, err
	}
	for i, r := range rows {
		if err := check(r); err != nil {
			var re *RowError
			if errors.As(err, &re) {
				return nil, &tsv.Error{Line: i + 2, Column: re.Column, Reason: re.Reason}
			}
			return nil, err
		}
	}
	return rows, nil
}

// ReadTelemetry reads telemetry.tsv by its schema and its row rules.
func ReadTelemetry(data []byte) ([][]string, error) {
	return read(data, TelemetrySchema, CheckTelemetry)
}

// ReadPrices reads prices.tsv by its schema and its row rules.
func ReadPrices(data []byte) ([][]string, error) { return read(data, PricesSchema, CheckPrice) }
