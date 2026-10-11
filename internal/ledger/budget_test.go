package ledger

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/records"
)

// telemetryRow gives a row of telemetry.tsv as records.ReadTelemetry gives it,
// with the money columns of the test and the empty value elsewhere.
func telemetryRow(session, money, currency, status string) []string {
	row := make([]string, len(records.TelemetrySchema.Columns))
	for i, c := range records.TelemetrySchema.Columns {
		switch c.Name {
		case "session":
			row[i] = session
		case "money":
			row[i] = money
		case "currency":
			row[i] = currency
		case "money_status":
			row[i] = status
		}
	}
	return row
}

// The budget of milestone 0 (docs/spec/intake.md, The budget of Intake): the
// known money since Start plus the new session's cap against MONEY, and the
// hours since Start against HOURS; "passes" is "is above".
func TestMilestoneCap(t *testing.T) {
	start := time.Date(2026, 10, 11, 0, 0, 0, 0, time.UTC)
	known := [][]string{telemetryRow("S-00000001", "1.25", "USD", "reported"), telemetryRow("S-00000002", "0.75", "USD", "computed")}
	unknown := append(append([][]string{}, known...), telemetryRow("S-00000003", "", "", "unknown"))
	for _, c := range []struct {
		name      string
		cap       string
		telemetry [][]string
		session   string
		now       time.Time
		want      string // "" for none, else a part of the error
	}{
		{"below", "5.00,8.0", known, "2.99", start, ""},
		{"equal is inside", "5.00,8.0", known, "3.00", start, ""},
		{"above", "5.00,8.0", known, "3.01", start, "the known money 2.0 plus the session's cap 3.01 passes 5.00 USD"},
		{"an unknown row adds nothing", "5.00,8.0", unknown, "3.00", start, ""},
		{"a session with no cap adds nothing", "2.00,8.0", known, "", start, ""},
		{"the known money alone above", "1.99,8.0", known, "", start, "the known money 2.0 plus the session's cap 0.0 passes 1.99 USD"},
		{"no row", "0.10,8.0", nil, "0.10", start, ""},
		{"the hours exactly", "5.00,8.0", known, "", start.Add(8 * time.Hour), ""},
		{"the hours and one second", "5.00,8.0", known, "", start.Add(8*time.Hour + time.Second), "the cap of milestone 0: the time since Start passes 8.0 hours"},
		{"the hours of a decimal", "5.00,0.5", known, "", start.Add(31 * time.Minute), "the cap of milestone 0: the time since Start passes 0.5 hours"},
		{"a known row not in USD", "5.00,8.0", append(append([][]string{}, known...), telemetryRow("S-00000004", "0.10", "EUR", "computed")), "", start,
			"the session S-00000004 has money in EUR, not USD"},
	} {
		err := MilestoneCap(c.cap, c.telemetry, c.session, start, c.now)
		switch {
		case c.want == "" && err != nil:
			t.Errorf("%s: %v; want none", c.name, err)
		case c.want != "" && (!errors.Is(err, ErrMilestoneCap) || !strings.HasPrefix(err.Error(), "the cap of milestone 0: ") || !strings.Contains(err.Error(), c.want)),
			strings.HasPrefix(c.want, "the cap of milestone 0: ") && err.Error() != c.want: // a detail with no clock reading (round 1 of #216)
			t.Errorf("%s: %v; want ErrMilestoneCap with %q", c.name, err, c.want)
		}
	}
}
