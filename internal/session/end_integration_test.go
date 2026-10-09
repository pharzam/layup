//go:build integration

package session

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/tsv"
)

const validResult = "kind\tn\tvalue\tsha256\treason\nstatus\t1\tcompleted\t—\tdone\n"

// The demo: with a fake harness program, each way a session ends gets its
// class.
func TestEachEndGetsItsClass(t *testing.T) {
	write := `printf 'kind\tn\tvalue\tsha256\treason\nstatus\t1\tcompleted\t—\tdone\n' > "$M/result.tsv"`
	for _, c := range []struct{ name, body, want string }{
		{"done", write + "; exit 0", "done"},
		{"no-result", "exit 0", "no-result"},
		{"crash by its exit", write + "; exit 3", "crash"},
		{"crash by a signal of its own", write + "; kill -SEGV $$", "crash"},
		{"wall", write + "; trap 'exit 0' INT; while :; do sleep 0.05; done", "wall"},
		{"output", write + "; while :; do echo aaaaaaaaa; done", "output"},
		{"start", "", "start"},
	} {
		d := dirOf(t.TempDir(), "S-1a2b3c4d")
		os.MkdirAll(d.Result, 0o700)
		s := spec(t, fake(t, c.body), d.Result, 300*time.Millisecond)
		if c.want == "start" {
			s.Words = []string{filepath.Join(t.TempDir(), "none")}
		}
		r, err := Process(s)
		_, resultErr := ResultOf(d, false)
		if got := Class(r, err, resultErr); got != c.want {
			t.Errorf("%s: %q (%+v, %v, %v), want %q", c.name, got, r, err, resultErr, c.want)
		}
	}
}

func TestResultOf(t *testing.T) {
	d := dirOf(t.TempDir(), "S-1a2b3c4d")
	os.MkdirAll(d.Result, 0o700)
	task, probe := filepath.Join(d.Result, "result.tsv"), filepath.Join(d.Result, "probe.tsv")
	os.WriteFile(task, []byte(validResult), 0o600)
	os.WriteFile(probe, []byte("kind\tvalue\ntoken\t0123456789abcdef\nfile\tAGENTS.md\n"), 0o600)
	if rows, err := ResultOf(d, false); err != nil || len(rows) != 1 {
		t.Errorf("a valid result.tsv: %q, %v", rows, err)
	}
	if rows, err := ResultOf(d, true); err != nil || len(rows) != 2 {
		t.Errorf("a valid probe.tsv: %q, %v", rows, err)
	}
	host := filepath.Join(t.TempDir(), "host.tsv")
	os.WriteFile(host, []byte(validResult), 0o600)
	// A valid result of n bytes: its status row's reason is padded.
	sized := func(n int) string {
		return strings.TrimSuffix(validResult, "done\n") + strings.Repeat("d", n-len(validResult)+4) + "\n"
	}
	// Each case asserts the words of its own rule, so another rule cannot
	// refuse it in its place.
	want := map[string]string{"missing": "no such file", "a link": "too many levels of symbolic links",
		"a directory": "is not a regular file", "over 1 MiB": "is over 1 MiB", "malformed": "line 1",
		"two status rows": "one status row"}
	for name, set := range map[string]func(){
		"missing":     func() { os.Remove(task) },
		"a link":      func() { os.Remove(task); os.Symlink(host, task) },
		"a directory": func() { os.Remove(task); os.Mkdir(task, 0o700) },
		"over 1 MiB":  func() { os.WriteFile(task, []byte(sized(1<<20+1)), 0o600) },
		"malformed":   func() { os.WriteFile(task, []byte("kind\tn\nstatus\t1\n"), 0o600) },
		"two status rows": func() {
			os.WriteFile(task, []byte(validResult+"status\t2\tfailed\t—\tagain\n"), 0o600)
		},
	} {
		os.RemoveAll(task)
		set()
		if _, err := ResultOf(d, false); err == nil || !strings.Contains(err.Error(), want[name]) {
			t.Errorf("a result.tsv %s: %v, want an error with %q", name, err, want[name])
		}
	}
	// A file of exactly 1 MiB is read: the limit is "over 1 MiB".
	os.RemoveAll(task)
	pad := sized(1 << 20)
	os.WriteFile(task, []byte(pad), 0o600)
	if _, err := ResultOf(d, false); err != nil || len(pad) != 1<<20 {
		t.Errorf("a result.tsv of 1 MiB (%d bytes): %v", len(pad), err)
	}
}

func TestUsageOfAFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "stdout")
	os.WriteFile(p, []byte("not json\n"+subagent), 0o600)
	u, err := UsageOf("claude-result", p)
	if err != nil || u.Status != "observed" || *u.Cache != 52469 {
		t.Errorf("%+v, %v", show(u), err)
	}
	if _, err := UsageOf("claude-result", filepath.Join(t.TempDir(), "none")); err == nil {
		t.Error("a missing stdout: no error")
	}
	// none reads nothing: a stdout that cannot be opened is no error.
	if u, err := UsageOf("none", filepath.Join(t.TempDir(), "none")); err != nil || u.Reason != "the harness reports none" {
		t.Errorf("none with no stdout: %+v, %v", show(u), err)
	}
}

func TestTheSchemaOfProbeResultEqualsItsBlock(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	if err := tsv.Compare(blocks["probe-result"], ProbeResultSchema); err != nil {
		t.Error(err)
	}
}
