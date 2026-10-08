package route

import (
	"reflect"
	"strings"
	"testing"
)

const harnessHeader = "harness\tcap\twall"

func harnesses(rows ...string) []byte {
	return []byte(strings.Join(append([]string{harnessHeader}, rows...), "\n") + "\n")
}

func TestAGoodHarnessRegisterIsRead(t *testing.T) {
	rows, ids, err := ReadHarnesses(harnesses("claude\t10.0\t60", "devin\t—\t30"))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ids, []string{"claude", "devin"}) || len(rows) != 2 || rows[1][1] != "" {
		t.Fatalf("ids %v, rows %v", ids, rows)
	}
	if rows, ids, err := ReadHarnesses(harnesses()); err != nil || len(rows) != 0 || len(ids) != 0 {
		t.Errorf("an empty register: %v, %v, %v; want it read, with no row", rows, ids, err)
	}
}

func TestHarnessRegisterRefusesEachBrokenRule(t *testing.T) {
	for name, data := range map[string][]byte{
		"two rows for one harness": harnesses("claude\t10.0\t60", "claude\t5.0\t30"),
		"an unknown column":        []byte(harnessHeader + "\tmodel\nclaude\t10.0\t60\tfable\n"),
		"an empty wall":            harnesses("claude\t10.0\t—"),
		"a wall of 0":              harnesses("claude\t10.0\t0"),
		"a harness ID in capitals": harnesses("Claude\t10.0\t60"),
	} {
		if _, _, err := ReadHarnesses(data); err == nil {
			t.Errorf("%s: read, want an error", name)
		} else if !strings.Contains(err.Error(), "line ") {
			t.Errorf("%s: the error %q names no line", name, err)
		}
	}
}
