package run

// This file holds the step markers of Intake (docs/spec/scaffold.md, The step
// markers; The steps of the Scaffold; Input states; row 53 of the plan, task
// T-4c5c, #276): the markers of the pinned tree, one session of a phase of
// systems-architect that proposes a value for each that the problem statement
// states, with a quote that code finds once, and the record of both.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/verify"
)

// markersFile is the form of result/markers.tsv: three columns, the marker's
// key, the proposed value and the quote of the problem statement that states
// it. It is no block of its own: the quote is not kept, as the offsets of the
// record give it (docs/spec/scaffold.md, The step markers, item 2).
var markersFile = tsv.Schema{Name: "markers", Location: "host:sessions/<session>/result/markers.tsv", Columns: []tsv.Column{
	{Name: "id", Type: "id(M-xxxxxxxx)", Key: true},
	{Name: "value", Type: "text"},
	{Name: "quote", Type: "text"},
}}

// markersList gives the file that the session reads, input/markers.tsv: the
// block markers with value, start, end and session as the empty value, one row
// for each place, in the order of the stop table of S10. Its form is the
// writer's: the reader of the block refuses a row with no session.
func markersList(places []setup.MarkerPlace) ([]byte, error) {
	rows := make([][]string, len(places))
	for i, p := range places {
		rows[i] = []string{p.ID, p.File, strconv.Itoa(p.Line), p.Marker, "", "", "", ""}
	}
	return table(records.MarkersSchema, rows)
}

// markersKeep gives the Keep of the step markers on the problem statement psb
// and the places of the scan: the session's file read by its form; a key that
// the list lacks, a key twice (the form's key), an empty value and an empty
// quote refuse the file, and so does a quote that quoteAt refuses (it occurs
// nowhere or more than once, or holds a tab). The rows of the file may be in
// any order (decided here, task T-4c5c): the record has one row for each place,
// in the scan's order, with the value and the offsets of its proposal, or the
// empty value in those three columns, and the session ID; the table passes the
// rules of the block, row by row.
func markersKeep(psb []byte, places []setup.MarkerPlace) func(id string, data []byte) (map[string][]byte, error) {
	return func(id string, data []byte) (map[string][]byte, error) {
		rows, err := tsv.Read(data, markersFile)
		if err != nil {
			return nil, err
		}
		type proposal struct {
			value      string
			start, end int
		}
		proposed := map[string]proposal{}
		for i, r := range rows {
			key, value, quote := r[0], r[1], r[2]
			if !slices.ContainsFunc(places, func(p setup.MarkerPlace) bool { return p.ID == key }) {
				return nil, fmt.Errorf("row %d: the key %s is not in the list", i+1, key)
			}
			if value == "" {
				return nil, fmt.Errorf("row %d: the value is empty", i+1)
			}
			start, end, err := quoteAt(psb, quote)
			if err != nil {
				return nil, fmt.Errorf("row %d: %w", i+1, err)
			}
			proposed[key] = proposal{value, start, end}
		}
		out := make([][]string, len(places))
		for i, p := range places {
			row := []string{p.ID, p.File, strconv.Itoa(p.Line), p.Marker, "", "", "", id}
			if q, ok := proposed[p.ID]; ok {
				row[4], row[5], row[6] = q.value, strconv.Itoa(q.start), strconv.Itoa(q.end)
			}
			out[i] = row
		}
		table, err := table(records.MarkersSchema, out)
		if err != nil {
			return nil, err
		}
		if _, err := records.ReadMarkers(table); err != nil {
			return nil, err
		}
		return map[string][]byte{"intake/markers.tsv": table}, nil
	}
}

// markersPrompt is the fixed text of the prompt of the step markers, with
// {input}, {list} and {output}; each path is on a line of its own.
const markersPrompt = `You propose values for the markers of a project's files, for LAYUP. Change no
file in this repository and make no commit.

The problem statement: {input}
The list of markers: {list}
The output file: {output}

The list holds one marker in each row: its id, its file, its line and its text
with the two angle quotes; its other columns are empty. For each marker whose
value the problem statement states, write one row of the output file, as
tab-separated values, with the header line "id<TAB>value<TAB>quote": the id of
the marker; the value that the problem statement gives it, on one line; and the
quote, which is text of one line of the problem statement, byte for byte, with
no tab, that states the value and occurs once in the file: quote enough to be
unique. Write no row for a marker that the problem statement states no value
for, and none for an id that is not in the list. With no proposal, write the
header line alone. Then end.
`

func markersText(input map[string]string, output string) []byte {
	return []byte(strings.NewReplacer("{input}", input["problem-statement.md"], "{list}", input["markers.tsv"], "{output}", output).Replace(markersPrompt))
}

// pinnedPlaces gives the places of the markers of the pinned tree: the tree of
// the one root commit of the default branch, which must be pin.tree, in a
// scratch work tree (the clone has none), by the scanner of S10 with MK_EXEMPT
// at the engine's version, one place for each file and key. It is the root
// commit's tree and not the head's, which after Scaffold 6 is the setup head
// (docs/spec/scaffold.md, The step markers, item 1).
func (r *state) pinnedPlaces() ([]setup.MarkerPlace, error) {
	roots, err := git.RootCommits(r.target, "refs/remotes/origin/"+r.defaultBranch)
	if err != nil {
		return nil, err
	}
	if len(roots) != 1 {
		return nil, fmt.Errorf("the default branch has %d root commits; want one", len(roots))
	}
	tree, err := git.RevParse(r.target, roots[0]+"^{tree}")
	if err != nil {
		return nil, err
	}
	if want := r.value("pin.tree"); tree != want {
		return nil, fmt.Errorf("the root tree %s is not pin.tree %s", tree, want)
	}
	dir, err := os.MkdirTemp("", "layup-markers-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "tree")
	if err := git.WorktreeAdd(r.target, path, roots[0]); err != nil {
		return nil, err
	}
	defer git.WorktreeRemove(r.target, path)
	found, err := verify.Markers(path)
	if err != nil {
		return nil, err
	}
	marks := make([]setup.Marker, len(found))
	for i, m := range found {
		marks[i] = setup.Marker{File: m.File, Line: m.Line, Col: m.Col, Text: m.Text}
	}
	return setup.MarkerPlaces(marks), nil
}

// markersStep is the step markers of Intake, between completeness and batch:
// "done before" when intake/markers.tsv exists; else the scan of the pinned
// tree; with no place, one records commit of the table with no row and no
// session; else one session of systems-architect on the problem statement and
// the list, and the records commit of its file (docs/spec/scaffold.md, The step
// markers). No admitted pair is fail (pair); a session that ends other than
// done, a refused file and the cap of milestone 0 are fail with their text.
func (r *state) markersStep(ctx context.Context) Step {
	const name = "markers"
	if data, err := r.store.ReadFile(ctx, "intake/markers.tsv"); err != nil {
		return fail(name, err)
	} else if data != nil {
		return done(name, "done before")
	}
	psb, err := r.store.ReadFile(ctx, "start/problem-statement.md")
	if err != nil {
		return fail(name, err)
	}
	if psb == nil {
		return failf(name, "the records have no start/problem-statement.md")
	}
	places, err := r.pinnedPlaces()
	if err != nil {
		return fail(name, err)
	}
	if len(places) == 0 {
		data, err := table(records.MarkersSchema, nil)
		if err != nil {
			return fail(name, err)
		}
		if err := r.write(ctx, map[string][]byte{"intake/markers.tsv": data}, "layup run: the markers of the pinned tree"); err != nil {
			return fail(name, err)
		}
		return done(name, "no marker in the pinned tree; intake/markers.tsv has no row")
	}
	list, err := markersList(places)
	if err != nil {
		return fail(name, err)
	}
	s := r.sessions()
	h, m, err := s.PhasePair(ctx, "systems-architect", r.cfg.Harnesses, r.cfg.Models, "")
	if errors.Is(err, route.ErrNoPair) {
		return failf(name, "pair: %v", err)
	} else if err != nil {
		return fail(name, err)
	}
	base, err := git.RevParse(r.target, "refs/remotes/origin/"+r.defaultBranch)
	if err != nil {
		return fail(name, err)
	}
	if _, err := s.PhaseSession(ctx, PhaseSpec{Step: name, Role: "systems-architect", Harness: h, Model: m, Models: r.cfg.Models,
		Base: base, Records: r.store.Base(), Input: map[string][]byte{"problem-statement.md": psb, "markers.tsv": list},
		Output: "markers.tsv", Prompt: markersText, Keep: markersKeep(psb, places)}); err != nil {
		return fail(name, err)
	}
	data, err := r.store.ReadFile(ctx, "intake/markers.tsv")
	if err != nil {
		return fail(name, err)
	}
	rows, err := records.ReadMarkers(data)
	if err != nil {
		return fail(name, err)
	}
	with := 0
	for _, row := range rows {
		if row[4] != "" {
			with++
		}
	}
	return done(name, fmt.Sprintf("%d markers, %d with a proposal, in intake/markers.tsv", len(rows), with))
}
