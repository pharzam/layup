package session

// This file holds the end of a session: its class, the result file, and the
// usage report in stdout (docs/spec/session.md, The end of a session and The
// usage report of a harness; records.md, the block probe-result; row 34b of
// the plan, task T-bpxg, #162).

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/tsv"
)

// Class gives the class of a session's end, in this order: start, when
// Process gave ErrStart; wall or output, a stop that layup run sent, whatever
// the exit and the result; crash, a non-zero exit or a signal that layup run
// did not send; no-result, exit 0 with no valid result file (resultErr); and
// done. Another error of Process is layup run's own, not a class: it gives "".
func Class(r Run, startErr, resultErr error) string {
	switch {
	case errors.Is(startErr, ErrStart):
		return "start"
	case startErr != nil:
		return ""
	case r.StoppedBy != "":
		return r.StoppedBy
	case r.Exit != 0: // a signal of its own gives -1
		return "crash"
	case resultErr != nil:
		return "no-result"
	}
	return "done"
}

// maxResult is the largest result file that is read: 1 MiB.
const maxResult = 1 << 20

// ResultOf reads the result file of the session in d: result/probe.tsv for a
// probe, else result/result.tsv, by its block. A file that is missing, is not
// a regular file (a link is not followed, so no file of the host reaches a
// record), is over 1 MiB, or that its reader refuses is an error.
func ResultOf(d Dir, probe bool) ([][]string, error) {
	path, read := filepath.Join(d.Result, "result.tsv"), records.ReadResult
	if probe {
		path, read = filepath.Join(d.Result, "probe.tsv"), ReadProbeResult
	}
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, fmt.Errorf("the result file: %w", err)
	}
	defer f.Close()
	if info, err := f.Stat(); err != nil || !info.Mode().IsRegular() {
		return nil, fmt.Errorf("the result file %s is not a regular file", path)
	}
	data, err := io.ReadAll(io.LimitReader(f, maxResult+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxResult {
		return nil, fmt.Errorf("the result file %s is over 1 MiB", path)
	}
	return read(data)
}

// ProbeResultSchema is the form of result/probe.tsv: the block probe-result
// of docs/spec/records.md, which this package owns.
var ProbeResultSchema = tsv.Schema{Name: "probe-result", Location: "host:sessions/<session>/result/probe.tsv", Columns: []tsv.Column{
	{Name: "kind", Type: "enum(token|file)", Key: true},
	{Name: "value", Type: "text", Key: true},
}}

// ReadProbeResult reads a probe.tsv by its schema, which refuses an empty
// key, and its rule in words: a token row comes at most once (none is the
// probe's failure token, row 38's).
func ReadProbeResult(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, ProbeResultSchema)
	if err != nil {
		return nil, err
	}
	tokens := 0
	for i, r := range rows {
		if r[0] == "token" {
			if tokens++; tokens > 1 {
				return nil, &tsv.Error{Line: i + 2, Column: "kind", Reason: "a token row comes at most once"}
			}
		}
	}
	return rows, nil
}

// A Usage is what the usage report of a session gave, with the fields of
// ledger.Usage in name, type and order, so internal/run converts one to the
// other: each token count, or nil when the report has none; the status and
// its reason; the reported cost as the text of the JSON number, with its
// currency, or "" for none; and the models that the report names, sorted.
type Usage struct {
	In, Out, Cache *int64
	Status, Reason string
	Cost, Currency string
	Models         []string
}

// UsageOf reads the usage report in the file stdout, of the format of the
// harness register row's usage.
func UsageOf(format, stdout string) (Usage, error) {
	f, err := os.Open(stdout)
	if err != nil {
		return Usage{}, err
	}
	defer f.Close()
	return usageOf(format, f)
}

// maxLine is the longest line of stdout that is read: the line cap of
// Process, and its line feed.
const maxLine = 8<<20 + 1

// usageOf reads a usage report from stdout: claude-result, the last line that
// is a JSON object whose type is result (a line that is not JSON is skipped),
// its total_cost_usd in USD and the sums over its modelUsage; none, nothing.
func usageOf(format string, stdout io.Reader) (Usage, error) {
	switch format {
	case "none":
		return Usage{Status: "unavailable", Reason: "the harness reports none"}, nil
	case "claude-result":
	default:
		return Usage{}, fmt.Errorf("the usage format %q is not one of the list", format)
	}
	type result struct {
		Type       string                                `json:"type"`
		Cost       *json.Number                          `json:"total_cost_usd"`
		ModelUsage map[string]map[string]json.RawMessage `json:"modelUsage"`
	}
	var last *result
	lines := bufio.NewScanner(stdout)
	lines.Buffer(make([]byte, 64<<10), maxLine)
	for lines.Scan() {
		var r result
		d := json.NewDecoder(bytes.NewReader(lines.Bytes()))
		d.UseNumber()
		if d.Decode(&r) == nil && r.Type == "result" {
			last = &r
		}
	}
	if err := lines.Err(); err != nil {
		return Usage{}, err
	}
	if last == nil {
		return Usage{Status: "unavailable", Reason: "stdout holds no result object"}, nil
	}
	var u Usage
	if last.Cost != nil {
		u.Cost, u.Currency = last.Cost.String(), "USD"
	}
	for m := range last.ModelUsage {
		u.Models = append(u.Models, m)
	}
	slices.Sort(u.Models)
	if len(u.Models) == 0 {
		u.Status, u.Reason = "unavailable", "the report names no model"
		return u, nil
	}
	// A class is summed only when each model gives each of its fields, as an
	// integer: a sum that leaves a model out would undercount.
	lacks := map[string][]string{} // the fields that each model lacks
	sum := func(fields ...string) *int64 {
		var total int64
		ok := true
		for _, m := range u.Models {
			var missing []string
			for _, f := range fields {
				var v int64
				if json.Unmarshal(last.ModelUsage[m][f], &v) != nil {
					missing = append(missing, f)
					continue
				}
				total += v
			}
			if len(missing) > 0 {
				ok = false
				lacks[m] = append(lacks[m], missing...)
			}
		}
		if !ok {
			return nil
		}
		return &total
	}
	u.In, u.Out, u.Cache = sum("inputTokens"), sum("outputTokens"), sum("cacheReadInputTokens", "cacheCreationInputTokens")
	switch summed := countSet(u.In, u.Out, u.Cache); summed {
	case 3:
		u.Status = "observed"
	case 0:
		u.Status = "unavailable"
	default:
		u.Status = "partial"
	}
	var reasons []string
	for _, m := range u.Models {
		if len(lacks[m]) > 0 {
			reasons = append(reasons, m+" has no "+strings.Join(lacks[m], ", "))
		}
	}
	u.Reason = strings.Join(reasons, "; ")
	return u, nil
}

// countSet gives the number of counts that are not nil.
func countSet(counts ...*int64) int {
	n := 0
	for _, c := range counts {
		if c != nil {
			n++
		}
	}
	return n
}
