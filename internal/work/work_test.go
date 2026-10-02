package work

import (
	"errors"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

const record = "step\tname\tvalue\tsource\tref\n" +
	"S01\tname\tacme\tanswer\tS01-name\n" +
	"S02\tpin.commit\t0123456789abcdef0123456789abcdef01234567\tcomputed\tgit ls-remote u HEAD\n" +
	"S02\tdone\tthe commit and the tree\tstep\t—\n"

const answers = "question\tanswer\tby\tsource\tquestion_text\n" +
	"S01-name\tacme\toperator\thttps://example.invalid/c/1\t—\n"

// A value of the record is read by its step and its name, the key of the
// record.
func TestARecordValueByStepAndName(t *testing.T) {
	r, err := ParseRecord([]byte(record))
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		step, name, want string
		ok               bool
	}{
		{"S01", "name", "acme", true},
		{"S02", "pin.commit", "0123456789abcdef0123456789abcdef01234567", true},
		{"S01", "pin.commit", "", false},
		{"S02", "pin.tree", "", false},
	} {
		if got, ok := r.Value(c.step, c.name); got != c.want || ok != c.ok {
			t.Errorf("Value(%s, %s) = %q, %v; want %q, %v", c.step, c.name, got, ok, c.want, c.ok)
		}
	}
	a, err := ParseAnswers([]byte(answers))
	if err != nil || len(a) != 1 || a[0][0] != "S01-name" || a[0][4] != "" {
		t.Errorf("ParseAnswers: %q, %v; want the one row, with the empty question text", a, err)
	}
}

// A record that does not match its schema is an error of internal/tsv, which
// names the line and the column.
func TestARecordOfAnotherFormIsAnError(t *testing.T) {
	for _, c := range []struct {
		name  string
		parse func() error
		line  int
	}{
		{"a source of no kind", func() error {
			_, err := ParseRecord([]byte("step\tname\tvalue\tsource\tref\nS01\tname\tacme\tguess\tS01-name\n"))
			return err
		}, 2},
		{"a step of another form", func() error {
			_, err := ParseRecord([]byte("step\tname\tvalue\tsource\tref\nS1\tname\tacme\tanswer\tS01-name\n"))
			return err
		}, 2},
		{"a step and a name twice", func() error {
			_, err := ParseRecord([]byte(record + "S01\tname\tother\tanswer\tS01-name\n"))
			return err
		}, 5},
		{"an answer by no one of the two", func() error {
			_, err := ParseAnswers([]byte("question\tanswer\tby\tsource\tquestion_text\nS01-name\tacme\tsomeone\tu\t—\n"))
			return err
		}, 2},
		{"the header of the record in the answers", func() error {
			_, err := ParseAnswers([]byte(record))
			return err
		}, 1},
	} {
		var e *tsv.Error
		if err := c.parse(); !errors.As(err, &e) || e.Line != c.line {
			t.Errorf("%s: %v; want a *tsv.Error at line %d", c.name, err, c.line)
		}
	}
}
