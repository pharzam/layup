package gate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/pharzam/layup/internal/tsv"
)

// ManifestSchema is the form of docs/gates.tsv: the block gate-manifest of
// docs/spec/gate.md. internal/catalog, which writes the manifest, holds the
// same Go value; the test of each package compares its value with the block.
var ManifestSchema = tsv.Schema{Name: "gate-manifest", Location: "target:docs/gates.tsv", Columns: []tsv.Column{
	{Name: "kind", Type: "id(<word>)", Key: true},
	{Name: "state", Type: "enum(active|pending)"},
	{Name: "tool", Type: "text"},
	{Name: "command", Type: "text"},
	{Name: "scope", Type: "list(text)"},
	{Name: "config", Type: "list(path)"},
}}

// A Kind is one row of the manifest.
type Kind struct {
	Name, State, Tool, Command string
	Scope                      []pattern
	Config                     []string
}

// inScope gives the first path that the scope of the kind matches, and
// whether there is one.
func (k Kind) inScope(paths []string) (string, bool) {
	for _, p := range paths {
		for _, s := range k.Scope {
			if s.match(p) {
				return p, true
			}
		}
	}
	return "", false
}

// readManifest reads the bytes of docs/gates.tsv. A manifest that does not
// match its schema, that has no row, or that holds a scope pattern of another
// form is an error.
func readManifest(data []byte) ([]Kind, error) {
	rows, err := tsv.Read(data, ManifestSchema)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("no kind")
	}
	var kinds []Kind
	for i, r := range rows {
		k := Kind{Name: r[0], State: r[1], Tool: r[2], Command: r[3], Config: strings.Fields(r[5])}
		for _, s := range strings.Fields(r[4]) {
			p, err := parsePattern(s)
			if err != nil {
				return nil, fmt.Errorf("line %d, column \"scope\": %v", i+2, err)
			}
			k.Scope = append(k.Scope, p)
		}
		kinds = append(kinds, k)
	}
	return kinds, nil
}
