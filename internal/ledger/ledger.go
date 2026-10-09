// Package ledger is the writer of the telemetry record of milestone M2b
// (docs/spec/session.md, REQ-011 — The writer of the telemetry record; row 35
// of the plan, task T-4c3q, #163): one row of telemetry.tsv per session, from
// the start row, layup run's own times, the usage report and host:prices.tsv.
// It reads no file and no clock; internal/run commits the row (rows 36b, 38).
package ledger

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/records"
)

// A Session is what the row takes from a session's start: the start row's
// columns, the register row's billing, the task's requirements (none before
// M2e), and layup run's own times; FirstOutput is the zero time when the
// process wrote nothing.
type Session struct {
	ID, Task, Role, Harness, Model, Billing string
	Requirements                            []string
	Start, FirstOutput, End                 time.Time
}

// A Usage is what the usage report gave (row 34b reads it): each token count,
// or nil when the report has none; the status and its reason; a reported cost
// as the text of the JSON number of the report, with its currency, or "" for
// none; and the models that the report names.
type Usage struct {
	In, Out, Cache *int64
	Status, Reason string
	Cost, Currency string
	Models         []string
}

const timeForm = "2006-01-02T15:04:05Z"

// The columns of a row of prices.tsv, as records.ReadPrices gives them.
const (
	priceID = iota
	priceHarness
	priceModel
	priceClass
	pricePrice
	priceCurrency
	priceSource
	priceDate
)

// Row gives the row of telemetry.tsv of one session, as records.ReadTelemetry
// gives a row ("" for the empty value), after records.CheckTelemetry passes
// it; a row that it refuses is an error, never a row. The money is reported
// when the report gives a cost and the billing is api; else computed from the
// tokens and the price rows of the session's harness and model, one per class,
// of the newest date among those rows; else unknown (docs/spec/session.md).
func Row(s Session, u Usage, prices [][]string) ([]string, error) {
	start, end := s.Start.UTC().Truncate(time.Second), s.End.UTC().Truncate(time.Second)
	first, latency := "", ""
	if !s.FirstOutput.IsZero() {
		f := s.FirstOutput.UTC().Truncate(time.Second)
		first, latency = f.Format(timeForm), strconv.FormatInt(int64(f.Sub(start)/time.Second), 10)
	}
	tokens := func(n *int64) string {
		if n == nil {
			return ""
		}
		return strconv.FormatInt(*n, 10)
	}
	money, currency, status, price, err := moneyOf(s, u, prices)
	if err != nil {
		return nil, err
	}
	r := []string{s.ID, s.Task, strings.Join(s.Requirements, " "), s.Role, s.Harness, s.Model, s.Billing,
		start.Format(timeForm), first, end.Format(timeForm), latency, strconv.FormatInt(int64(end.Sub(start)/time.Second), 10),
		tokens(u.In), tokens(u.Out), tokens(u.Cache), u.Status, u.Reason, money, currency, status, price}
	if err := records.CheckTelemetry(r); err != nil {
		return nil, fmt.Errorf("the telemetry row of %s: %w", s.ID, err)
	}
	return r, nil
}

// moneyOf gives money, currency, money_status and price.
func moneyOf(s Session, u Usage, prices [][]string) (string, string, string, string, error) {
	if u.Cost != "" && s.Billing == "api" {
		cost, ok := new(big.Rat).SetString(u.Cost)
		if !ok || cost.Sign() < 0 {
			return "", "", "", "", fmt.Errorf("the reported cost %q is not a number of 0 or more", u.Cost)
		}
		text, err := decimal(cost)
		if err != nil {
			return "", "", "", "", err
		}
		return text, u.Currency, "reported", "", nil
	}
	if m, cur, id, ok := computed(s, u, prices); ok {
		text, err := decimal(m)
		if err != nil {
			return "", "", "", "", err
		}
		return text, cur, "computed", id, nil
	}
	return "", "", "unknown", "", nil
}

// computed gives the money of the tokens at the newest price rows of the
// session's harness and model, and false for each reason of an unknown money:
// tokens that are not observed, a report that names more than one model or
// one model other than the session's, a class with no row or two rows of that
// date, or rows of two currencies.
func computed(s Session, u Usage, prices [][]string) (*big.Rat, string, string, bool) {
	if u.Status != "observed" || u.In == nil || u.Out == nil || u.Cache == nil ||
		len(u.Models) > 1 || (len(u.Models) == 1 && u.Models[0] != s.Model) {
		return nil, "", "", false
	}
	newest := ""
	for _, r := range prices {
		if r[priceHarness] == s.Harness && r[priceModel] == s.Model && r[priceDate] > newest {
			newest = r[priceDate]
		}
	}
	byClass := map[string][]string{}
	currencies := map[string]bool{}
	for _, r := range prices {
		if r[priceHarness] != s.Harness || r[priceModel] != s.Model || r[priceDate] != newest {
			continue
		}
		if byClass[r[priceClass]] != nil {
			return nil, "", "", false
		}
		byClass[r[priceClass]] = r
		currencies[r[priceCurrency]] = true
	}
	if len(byClass) != 3 || len(currencies) != 1 {
		return nil, "", "", false
	}
	sum := new(big.Rat)
	for class, n := range map[string]int64{"in": *u.In, "out": *u.Out, "cache": *u.Cache} {
		p, ok := new(big.Rat).SetString(byClass[class][pricePrice])
		if !ok {
			return nil, "", "", false
		}
		sum.Add(sum, new(big.Rat).Mul(p, new(big.Rat).SetInt64(n)))
	}
	sum.Quo(sum, big.NewRat(1000000, 1))
	return sum, byClass["in"][priceCurrency], byClass["in"][priceID], true
}

// decimal writes a number of a finite decimal form exactly, in the form of
// the type decimal: no trailing zero, and one decimal digit at least.
func decimal(r *big.Rat) (string, error) {
	scaled, places := new(big.Rat).Set(r), 0
	ten := big.NewRat(10, 1)
	for !scaled.IsInt() {
		if places > 60 {
			return "", fmt.Errorf("%s has no finite decimal form", r.RatString())
		}
		scaled.Mul(scaled, ten)
		places++
	}
	digits := scaled.Num().String()
	if places == 0 {
		return digits + ".0", nil
	}
	for len(digits) <= places {
		digits = "0" + digits
	}
	return digits[:len(digits)-places] + "." + digits[len(digits)-places:], nil
}
