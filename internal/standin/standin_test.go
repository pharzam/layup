package standin

import (
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/work"
)

// The pin of a work area has the five keys of NFR-006, in the order of
// LAYUP's own pin file, with the method of a target and the date of pin.time.
func TestThePinText(t *testing.T) {
	got := pinText("file:///b", "c1", "t1")
	want := "source=file:///b\ncommit=c1\ntree=t1\nmethod=git clone file:///b, checkout c1, .git removed\ndate=2026-10-02\n"
	if got != want {
		t.Errorf("pinText:\n%s\nwant\n%s", got, want)
	}
}

// The record of a work area has the rows of S01 and S02, by its schema; an
// option changes a value, or removes a row.
func TestTheRecordRows(t *testing.T) {
	rows := recordRows("file:///b", "c1", "t1", map[string]string{"S02 pin.commit": "c2", "S01 name": ""})
	r := work.Record(rows)
	if v, ok := r.Value("S02", "pin.commit"); !ok || v != "c2" {
		t.Errorf("pin.commit: %q, %v; want the changed value c2", v, ok)
	}
	if _, ok := r.Value("S01", "name"); ok {
		t.Error("the row S01 name is there; the option removes it")
	}
	for _, name := range []string{"pin.source", "pin.tree", "pin.time"} {
		if v, ok := r.Value("S02", name); !ok || v == "" {
			t.Errorf("no value at S02 %s", name)
		}
	}
	for _, row := range rows {
		if len(row) != len(work.RecordSchema.Columns) || !strings.HasPrefix(row[0], "S0") {
			t.Errorf("the row %q is not a row of the schema setup-record", row)
		}
	}
}
