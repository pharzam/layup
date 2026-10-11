package cli

import (
	"strings"
	"testing"
)

// A table of the span check that the command cannot write gives exit 2 and the
// error on standard error, with a pass and with a fail (decided here, task
// T-hpid, as for layup psb check).
func TestSpecCheckGivesTwoWhenItCannotWriteTheTable(t *testing.T) {
	saved := readFile
	t.Cleanup(func() { readFile = saved })
	const head = "start\tend\tfact\tclass\treason\ttext\n"
	for _, spans := range []string{head + "0\t2\t1\tneed\t—\tAb\n", head} {
		readFile = func(name string) ([]byte, error) {
			if name == "spans.tsv" {
				return []byte(spans), nil
			}
			return []byte("Ab\n"), nil
		}
		var errOut strings.Builder
		if code := Run([]string{"spec", "check", "--facts", "spans.tsv", "--psb", "psb.md"}, refusingWriter{}, &errOut); code != 2 || errOut.String() != "layup: the output refuses the write\n" {
			t.Errorf("%q: exit %d, stderr %q; want 2 and the error of the output", spans, code, errOut.String())
		}
	}
}
