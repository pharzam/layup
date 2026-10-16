package spec

// The confirmed inventory (docs/spec/loop.md, REQ-012 — The rendered commit
// and the confirmed inventory; the blocks inventory and inventory-index of
// records.md, which this package owns; task T-3q0b, #350): a new version at
// each bet, versions from 1, consecutive, the earlier bytes never rewritten;
// the index holds the SHA-256 of the version's exact bytes.

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"slices"
	"strconv"

	"github.com/pharzam/layup/internal/tsv"
)

// InventorySchema is the form of spec/inventory-<n>.tsv: the block inventory.
var InventorySchema = tsv.Schema{Name: "inventory", Location: "records:spec/inventory-<n>.tsv", Columns: []tsv.Column{
	{Name: "fact", Type: "text", Key: true},
	{Name: "class", Type: "enum(need|constraint|measure|context)"},
	{Name: "scope", Type: "enum(in|out)"},
	{Name: "copy", Type: "int"},
}}

// InventoryIndexSchema is the form of spec/inventory.tsv: the block
// inventory-index.
var InventoryIndexSchema = tsv.Schema{Name: "inventory-index", Location: "records:spec/inventory.tsv", Columns: []tsv.Column{
	{Name: "n", Type: "int", Key: true},
	{Name: "sha256", Type: "sha256"},
	{Name: "bet", Type: "text"},
}}

// The columns of an inventory row, and of an index row.
const (
	invFact = iota
	invClass
	invScope
	invCopy
)
const (
	idxN = iota
	idxSHA
	idxBet
)

// ReadInventory reads an inventory by its schema and its owner rules (the
// decided-here of records.md): no column holds —, a fact has one row (the
// key), and the copy is a positive comment ID.
func ReadInventory(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, InventorySchema)
	if err != nil {
		return nil, err
	}
	for i, r := range rows {
		for c, name := range []string{"fact", "class", "scope", "copy"} {
			if r[c] == "" {
				return nil, fmt.Errorf("line %d, column %q: — is not a value of the inventory", i+2, name)
			}
		}
		if r[invCopy] == "0" {
			return nil, fmt.Errorf("line %d, column \"copy\": a copy is the comment ID of the bet's or the decision's copy, 1 or more", i+2)
		}
	}
	return rows, nil
}

// ReadInventoryIndex reads the index by its schema and its owner rules:
// sha256 and bet never —.
func ReadInventoryIndex(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, InventoryIndexSchema)
	if err != nil {
		return nil, err
	}
	for i, r := range rows {
		for _, c := range []int{idxSHA, idxBet} {
			if r[c] == "" {
				return nil, fmt.Errorf("line %d, column %q: — is not a value of the index", i+2, InventoryIndexSchema.Columns[c].Name)
			}
		}
	}
	return rows, nil
}

// InventoryFact is a fact the inventory records: its ID (a number of the span
// facts, or the Q-NNN of a fact of the answers), its class from the
// classified inputs — never inferred, the schemas hold none — and its text,
// which the caller supplies for rendering and which is no column.
type InventoryFact struct{ ID, Class, Text string }

// NextInventory gives the next version of the confirmed inventory and the
// index with its new row: version n = len(index)+1, n >= 2 only with the rows
// of version n-1, the index's n in order 1..n-1. The first version is the
// facts with scope in; a later one carries the previous rows forward and adds
// each fact the records know since. The answer's explicit changes apply; an
// unknown fact or change is a finding. Every row names this version's copy,
// and the hash is of the version's exact complete bytes.
func NextInventory(previous, index [][]string, facts []InventoryFact, answer Answer, copyID int) (version, nextIndex []byte, findings []Finding) {
	if copyID <= 0 {
		findings = append(findings, Finding{"copy", fmt.Sprintf("the copy %d is no positive comment ID", copyID)})
	}
	if answer.Milestone == "" {
		findings = append(findings, Finding{"index", "the bet's milestone, for the index's bet column, is not set"})
	}
	for i, r := range index {
		if want := strconv.Itoa(i + 1); len(r) <= idxN || r[idxN] != want {
			findings = append(findings, Finding{"index", fmt.Sprintf("the index's row %d names version %s, want %s", i+1, r[idxN], want)})
		}
	}
	n := len(index) + 1
	if n > 1 && len(previous) == 0 {
		findings = append(findings, Finding{"version", fmt.Sprintf("version %d has no version %d", n, n-1)})
	}
	if findings != nil {
		return nil, nil, findings
	}
	known := map[string]bool{}
	var rows [][]string
	add := func(id, class, scope string) {
		known[id] = true
		rows = append(rows, []string{id, class, scope, strconv.Itoa(copyID)})
	}
	if n > 1 {
		for _, r := range previous {
			add(r[invFact], r[invClass], r[invScope])
		}
	}
	for _, f := range facts {
		if !slices.Contains(classes, f.Class) {
			findings = append(findings, Finding{"fact", fmt.Sprintf("the fact %s has the class %q, none of the four", f.ID, f.Class)})
			continue
		}
		if !known[f.ID] {
			add(f.ID, f.Class, "in")
		}
	}
	for id, ch := range answer.Facts {
		if !known[id] {
			findings = append(findings, Finding{"fact", fmt.Sprintf("the answer names the fact %s, which the records do not know", id)})
			continue
		}
		if ch.Class != "" && !slices.Contains(classes, ch.Class) {
			findings = append(findings, Finding{"fact", fmt.Sprintf("the answer gives the fact %s the class %q, none of the four", id, ch.Class)})
			continue
		}
		for _, r := range rows {
			if r[invFact] == id {
				if ch.Out {
					r[invScope] = "out"
				}
				if ch.Class != "" {
					r[invClass] = ch.Class
				}
			}
		}
	}
	if findings != nil {
		return nil, nil, findings
	}
	var vb, ib bytes.Buffer
	if err := tsv.Write(&vb, InventorySchema, rows); err != nil {
		return nil, nil, []Finding{{"version", err.Error()}}
	}
	version = vb.Bytes()
	idx := make([][]string, len(index)+1)
	copy(idx, index)
	idx[len(index)] = []string{strconv.Itoa(n), fmt.Sprintf("%x", sha256.Sum256(version)), answer.Milestone}
	if err := tsv.Write(&ib, InventoryIndexSchema, idx); err != nil {
		return nil, nil, []Finding{{"index", err.Error()}}
	}
	return version, ib.Bytes(), nil
}
