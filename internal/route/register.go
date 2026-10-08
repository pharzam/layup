// Package route reads the harness register of the host (docs/spec/records.md,
// the block harness-register). In milestone M2a it is only that reader (row
// 22b of the plan, task T-1g1q, #137); the probe, admission and routing come
// in M2b (docs/spec/packages.md).
package route

import "github.com/pharzam/layup/internal/tsv"

// HarnessRegisterSchema is the form of host:registers/harnesses.tsv: the
// columns that Start reads; M2b adds its own.
var HarnessRegisterSchema = tsv.Schema{Name: "harness-register", Location: "host:registers/harnesses.tsv", Columns: []tsv.Column{
	{Name: "harness", Type: "id(<word>)", Key: true},
	{Name: "cap", Type: "decimal"},
	{Name: "wall", Type: "int"},
}}

// ReadHarnesses reads the harness register by its schema and the rules of its
// block, and gives its rows and the harness IDs in their order: wall is 1 or
// more and never the empty value; cap may be the empty value; a register with
// no row is allowed (docs/spec/run.md, Input states). An error names its line.
func ReadHarnesses(data []byte) ([][]string, []string, error) {
	rows, err := tsv.Read(data, HarnessRegisterSchema)
	if err != nil {
		return nil, nil, err
	}
	ids := make([]string, 0, len(rows))
	for i, r := range rows {
		if r[2] == "" || r[2] == "0" {
			return nil, nil, &tsv.Error{Line: i + 2, Column: "wall", Reason: "the wall-clock limit is 1 minute or more"}
		}
		ids = append(ids, r[0])
	}
	return rows, ids, nil
}
