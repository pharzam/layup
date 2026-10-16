package session

import (
	_ "embed"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/pharzam/layup/internal/ledger"
)

func TestClass(t *testing.T) {
	sent := []syscall.Signal{syscall.SIGINT}
	for _, c := range []struct {
		name                string
		r                   Run
		startErr, resultErr error
		want                string
	}{
		{"exit 0 and a valid result", Run{}, nil, nil, "done"},
		{"a program that cannot start", Run{}, ErrStart, errors.New("no file"), "start"},
		{"a non-zero exit", Run{Exit: 3}, nil, nil, "crash"},
		{"a signal that layup run did not send", Run{Exit: -1, Signal: syscall.SIGSEGV}, nil, nil, "crash"},
		{"exit 0 with no valid result", Run{}, nil, errors.New("no file"), "no-result"},
		// A stop that layup run sent is its class, whatever the exit and the result.
		{"stopped at wall, exit 0, a valid result", Run{StoppedBy: "wall", Sent: sent}, nil, nil, "wall"},
		{"stopped at an output cap, a non-zero exit", Run{StoppedBy: "output", Sent: sent, Exit: 2}, nil, errors.New("no file"), "output"},
		{"an error of layup run's own", Run{}, errors.New("disk full"), nil, ""},
	} {
		if got := Class(c.r, c.startErr, c.resultErr); got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
}

//go:embed testdata/claude-result-one-model.jsonl
var oneModel string

//go:embed testdata/claude-result-subagent.jsonl
var subagent string

func n(v int64) *int64 { return &v }

func TestUsageOfClaudeResult(t *testing.T) {
	for _, c := range []struct {
		name, stdout string
		want         Usage
	}{
		// The sums worked out by hand from each recorded result event of Claude
		// Code 2.1.295 (runs/T-bpxg/test-runs.md).
		{"one model", oneModel, Usage{In: n(290), Out: n(39003), Cache: n(775649 + 145871), Status: "observed",
			Cost: "5.0643822499999995", Currency: "USD", Models: []string{"claude-fable-5-1"}}},
		{"a subagent on a second model", subagent, Usage{In: n(34 + 2), Out: n(157 + 4), Cache: n(31161 + 9863 + 0 + 11445), Status: "observed",
			Cost: "0.27055324999999997", Currency: "USD", Models: []string{"claude-fable-5-1", "claude-opus-5-5"}}},
		{"the last result object, after lines that are not JSON", "not json\n{\"type\":\"result\",\"total_cost_usd\":9,\"modelUsage\":{\"m\":{\"inputTokens\":1,\"outputTokens\":1,\"cacheReadInputTokens\":1,\"cacheCreationInputTokens\":1}}}\n[1]\n{\"type\":\"assistant\"}\n" + subagent + "{oops\n",
			Usage{In: n(36), Out: n(161), Cache: n(52469), Status: "observed", Cost: "0.27055324999999997", Currency: "USD", Models: []string{"claude-fable-5-1", "claude-opus-5-5"}}},
		{"a line of 100 KiB before the result", strings.Repeat("x", 100<<10) + "\n" + subagent,
			Usage{In: n(36), Out: n(161), Cache: n(52469), Status: "observed", Cost: "0.27055324999999997", Currency: "USD", Models: []string{"claude-fable-5-1", "claude-opus-5-5"}}},
		{"no result object", "{\"type\":\"assistant\"}\nnot json\n", Usage{Status: "unavailable", Reason: "stdout holds no result object"}},
		// The cost is the text of the JSON number: 1.50 stays 1.50.
		{"a result with no modelUsage", "{\"type\":\"result\",\"total_cost_usd\":1.50}\n",
			Usage{Status: "unavailable", Reason: "the report names no model", Cost: "1.50", Currency: "USD"}},
		{"a result with an empty modelUsage", "{\"type\":\"result\",\"modelUsage\":{}}\n", Usage{Status: "unavailable", Reason: "the report names no model"}},
		{"a model that lacks a field", "{\"type\":\"result\",\"modelUsage\":{\"b\":{\"inputTokens\":1,\"outputTokens\":2,\"cacheReadInputTokens\":3},\"a\":{\"inputTokens\":4,\"outputTokens\":5,\"cacheReadInputTokens\":6,\"cacheCreationInputTokens\":7}}}\n",
			Usage{In: n(5), Out: n(7), Status: "partial", Reason: "b has no cacheCreationInputTokens", Models: []string{"a", "b"}}},
		// A null is no integer: the class is not summed, never 0.
		{"a field that is null", "{\"type\":\"result\",\"modelUsage\":{\"a\":{\"inputTokens\":null,\"outputTokens\":2,\"cacheReadInputTokens\":3,\"cacheCreationInputTokens\":4}}}\n",
			Usage{Out: n(2), Cache: n(7), Status: "partial", Reason: "a has no inputTokens", Models: []string{"a"}}},
		{"each class lacked", "{\"type\":\"result\",\"modelUsage\":{\"a\":{\"inputTokens\":\"x\",\"webSearchRequests\":0}}}\n",
			Usage{Status: "unavailable", Reason: "a has no inputTokens, outputTokens, cacheReadInputTokens, cacheCreationInputTokens", Models: []string{"a"}}},
	} {
		got, err := usageOf("claude-result", strings.NewReader(c.stdout))
		if err != nil || !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: %+v, %v; want %+v", c.name, show(got), err, show(c.want))
		}
	}
}

// show gives u with its counts as values, for a message.
func show(u Usage) []any {
	v := func(p *int64) any {
		if p == nil {
			return nil
		}
		return *p
	}
	return []any{v(u.In), v(u.Out), v(u.Cache), u.Status, u.Reason, u.Cost, u.Currency, u.Models}
}

func TestUsageOfNone(t *testing.T) {
	got, err := usageOf("none", strings.NewReader(subagent))
	want := Usage{Status: "unavailable", Reason: "the harness reports none"}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("none: %+v, %v; want %+v", show(got), err, show(want))
	}
	if _, err := usageOf("other", strings.NewReader("")); err == nil {
		t.Error("a format of neither: no error")
	}
}

// Usage has the fields of ledger.Usage in name, type and order, so
// internal/run converts one to the other.
var _ = ledger.Usage(Usage{})

func TestCheckProbeResult(t *testing.T) {
	for name, data := range map[string]string{
		"a token and two files": "kind\tvalue\ntoken\t0123456789abcdef\nfile\tAGENTS.md\nfile\tCLAUDE.md\n",
		"no token":              "kind\tvalue\nfile\tAGENTS.md\n",
	} {
		if _, err := ReadProbeResult([]byte(data)); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
	for name, data := range map[string]string{
		"two tokens":     "kind\tvalue\ntoken\t0123456789abcdef\ntoken\tfedcba9876543210\n",
		"a wrong header": "kind\tpath\ntoken\t0123456789abcdef\n",
		"an empty value": "kind\tvalue\nfile\t—\n",
		"another kind":   "kind\tvalue\nrule\tAGENTS.md\n",
	} {
		if _, err := ReadProbeResult([]byte(data)); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
}

// resultDir gives a session directory whose result/ holds the files by name.
func resultDir(t *testing.T, files map[string]string) Dir {
	t.Helper()
	d := Dir{Result: filepath.Join(t.TempDir(), "result")}
	if err := os.MkdirAll(d.Result, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, text := range files {
		if err := os.WriteFile(filepath.Join(d.Result, name), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return d
}

func TestACommentFileIsReadByteForByte(t *testing.T) {
	for _, text := range []string{"a plan\n", "no final newline", "crlf\r\nlines\r\n", "trailing  \n  spaces \t\n", "ü 𝄞 non-ASCII\n", "a\x00nul\n"} {
		d := resultDir(t, map[string]string{"plan.md": text})
		got, err := CommentFile(d, "plan.md")
		if err != nil || string(got) != text {
			t.Errorf("%q: got %q, %v; want the bytes of the file", text, got, err)
		}
	}
}

func TestACommentFileThatIsNotARegularFileOfOneNameIsRefused(t *testing.T) {
	d := resultDir(t, map[string]string{"plan.md": "text\n"})
	outside := filepath.Join(t.TempDir(), "host-file")
	if err := os.WriteFile(outside, []byte("a file of the host\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(d.Result, "link.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(d.Result, "dir.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	big := strings.Repeat("x", maxResult+1)
	if err := os.WriteFile(filepath.Join(d.Result, "big.md"), []byte(big), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"link.md", "dir.md", "big.md", "missing.md", "", ".", "..", "a/plan.md", "../plan.md", "plan.md/"} {
		if got, err := CommentFile(d, name); err == nil {
			t.Errorf("%q: read %d bytes, want an error", name, len(got))
		}
	}
	// A link above result/ is followed, as for the result file: result/ itself
	// is what the session could have made.
	real := resultDir(t, map[string]string{"plan.md": "through a link\n"})
	linked := Dir{Result: filepath.Join(t.TempDir(), "result")}
	if err := os.Symlink(real.Result, linked.Result); err != nil {
		t.Fatal(err)
	}
	if got, err := CommentFile(linked, "plan.md"); err != nil || string(got) != "through a link\n" {
		t.Errorf("a link above result/: %q, %v; want the file", got, err)
	}
}
