package ledger

// This file holds the budget of milestone 0 (docs/spec/intake.md, The budget
// of Intake; Input states, row 7; row 42b of the plan, task T-ybd5, #216): the
// rule that internal/run asks before each session of Intake starts. It reads
// no file and no clock: internal/run gives it the records and the time.

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/records"
)

// ErrMilestoneCap is the refusal of a session of Intake by the cap of
// milestone 0: the step is fail, and no session starts.
var ErrMilestoneCap = errors.New("the cap of milestone 0")

// column gives the index of the column name of telemetry.tsv.
func column(name string) int {
	for i, c := range records.TelemetrySchema.Columns {
		if c.Name == name {
			return i
		}
	}
	return -1
}

// MilestoneCap gives nil when a new session of Intake may start under the cap
// capValue of milestone 0 (intake.cap of start.tsv, MONEY,HOURS: US dollars
// and hours), else an error that wraps ErrMilestoneCap. The known money is the
// sum of money of each row of telemetry (rows as records.ReadTelemetry gives
// them, the probes' included) whose money_status is reported or computed, in
// exact decimals; sessionCap, the new session's spend cap ("" for none, which
// adds nothing), is added to it. The sum passes MONEY when it is above it
// (equal is inside, decided here, task T-ybd5), and the time from start to
// now passes HOURS the same way. A known row whose currency is not USD cannot
// be compared with MONEY, and refuses on the safe side (decided here).
func MilestoneCap(capValue string, telemetry [][]string, sessionCap string, start, now time.Time) error {
	moneyText, hoursText, _ := strings.Cut(capValue, ",")
	money, okMoney := new(big.Rat).SetString(moneyText)
	hours, okHours := new(big.Rat).SetString(hoursText)
	if !okMoney || !okHours {
		return fmt.Errorf("intake.cap %q is not MONEY,HOURS", capValue)
	}
	session, status, amount, currency := column("session"), column("money_status"), column("money"), column("currency")
	sum := new(big.Rat)
	for _, r := range telemetry {
		if len(r) != len(records.TelemetrySchema.Columns) {
			return fmt.Errorf("a row of telemetry.tsv of %d columns", len(r))
		}
		if r[status] != "reported" && r[status] != "computed" {
			continue
		}
		if r[currency] != "USD" {
			return fmt.Errorf("%w: the session %s has money in %s, not USD", ErrMilestoneCap, r[session], r[currency])
		}
		m, ok := new(big.Rat).SetString(r[amount])
		if !ok {
			return fmt.Errorf("the money %q of the session %s is not a decimal", r[amount], r[session])
		}
		sum.Add(sum, m)
	}
	known, err := decimal(sum)
	if err != nil {
		return err
	}
	total, capText := new(big.Rat).Set(sum), "0.0"
	if sessionCap != "" {
		c, ok := new(big.Rat).SetString(sessionCap)
		if !ok {
			return fmt.Errorf("the session's cap %q is not a decimal", sessionCap)
		}
		total.Add(total, c)
		capText = sessionCap
	}
	if total.Cmp(money) > 0 {
		return fmt.Errorf("%w: the known money %s plus the session's cap %s passes %s USD", ErrMilestoneCap, known, capText, moneyText)
	}
	elapsed := now.Sub(start)
	if new(big.Rat).SetFrac64(int64(elapsed), int64(time.Hour)).Cmp(hours) > 0 {
		return fmt.Errorf("%w: %s since Start passes %s hours", ErrMilestoneCap, elapsed.Round(time.Second), hoursText)
	}
	return nil
}
