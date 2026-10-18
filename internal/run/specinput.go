package run

// The input of the full spec check at a head (docs/spec/loop.md, REQ-012, The
// full spec check; task T-cav3, #412, row 66p of the plan of M2e): each field
// of spec.FullInput is built from one named record or from the tree at the
// head, and a caller's literal of a field is never read.

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/spec"
	"github.com/pharzam/layup/internal/tsv"
)

// SourceError is a refusal of the input of the full check. Field is the field
// of spec.FullInput (or base, head) and Source is "base" or "head" for a SHA
// that is not a full SHA (a caller's error), "records:<path>" for a record, or
// "target:<what>@<sha>" for the tree of the run's clone (a posted error of the
// status; the status is row 66f's).
type SourceError struct {
	Field, Source string
	Err           error
}

func (e *SourceError) Error() string {
	return fmt.Sprintf("the input of the full spec check: %s from %s: %v", e.Field, e.Source, e.Err)
}
func (e *SourceError) Unwrap() error { return e.Err }

// numberedFacts is the name of the numbered facts record of the problem
// statement, the one file of docs/facts/ that row 66j's rendered commit writes.
var numberedFacts = regexp.MustCompile(`^F-[0-9]{4}-numbered-facts\.md$`)

func isNumberedFacts(base string) bool { return numberedFacts.MatchString(base) }

// psbRecordsOf gives the record IDs of the numbered facts files: none for no
// file, F-NNNN for one, and an error for two or more.
func psbRecordsOf(files map[string][]byte) ([]string, error) {
	switch len(files) {
	case 0:
		return nil, nil
	case 1:
		for name := range files {
			return []string{name[:len("F-NNNN")]}, nil
		}
	}
	return nil, fmt.Errorf("%d numbered facts records at the head; want at most one", len(files))
}

// specInput gives the FullInput of head against base: the records of the run's
// records base at the call, the tree of the run's clone, and the caller's copy.
// It writes nothing, and on any error gives the zero input.
func (s *Sessions) specInput(ctx context.Context, base, head string, copy []byte) (spec.FullInput, error) {
	for _, c := range []struct{ field, sha string }{{"base", base}, {"head", head}} {
		if !forge.FullSHA(c.sha) {
			return spec.FullInput{}, &SourceError{c.field, c.field, fmt.Errorf("%q is not a full SHA", c.sha)}
		}
	}
	var in spec.FullInput
	if err := s.specRecords(ctx, &in); err != nil {
		return spec.FullInput{}, err
	}
	if err := specTree(s.Clone, base, head, &in); err != nil {
		return spec.FullInput{}, err
	}
	in.Copy = copy
	return in, nil
}

// specRecords sets the fields of in that the records hold: all required but
// the index of the inventory and the events of a task, whose absence is "before
// the first bet" and "not merged".
func (s *Sessions) specRecords(ctx context.Context, in *spec.FullInput) error {
	// raw reads the required record path and keeps its exact bytes.
	raw := func(field, path string, read func([]byte) ([][]string, error)) ([][]string, []byte, error) {
		var data []byte
		rows, err := readRequired(ctx, s.Store, path, func(b []byte) ([][]string, error) {
			data = b
			if read == nil {
				return nil, nil
			}
			return read(b)
		})
		if err != nil {
			return nil, nil, &SourceError{field, "records:" + path, err}
		}
		return rows, data, nil
	}
	var err error
	if _, in.PSB, err = raw("PSB", "start/problem-statement.md", nil); err != nil {
		return err
	}
	for _, f := range []struct {
		field, path string
		read        func([]byte) ([][]string, error)
		to          *[][]string
	}{
		{"Spans", "spec/spans.tsv", records.ReadSpans, &in.Spans},
		{"Questions", "questions.tsv", records.ReadQuestions, &in.Questions},
		{"Answers", "intake/answers.tsv", records.ReadAnswers, &in.Answers},
	} {
		if *f.to, _, err = raw(f.field, f.path, f.read); err != nil {
			return err
		}
	}
	if err := s.specIndex(ctx, in, raw); err != nil {
		return err
	}
	register, data, err := raw("Register", tasksPath, func(b []byte) ([][]string, error) { return tsv.Read(b, TasksSchema) })
	if err != nil {
		return err
	}
	in.Register = data
	for _, r := range register {
		events, err := s.readTable(ctx, eventsPath(r[0]), records.ReadEvents)
		if err != nil {
			return &SourceError{"Merged", "records:" + eventsPath(r[0]), err}
		}
		for _, e := range events {
			if e[1] == "merged" {
				in.Merged = append(in.Merged, r[0])
				break
			}
		}
	}
	return nil
}

// specIndex sets Index and Versions: no spec/inventory.tsv is before the first
// bet (both nil); the file is the rows of the index, with no row an error, and
// each version it names is required, as its exact bytes.
func (s *Sessions) specIndex(ctx context.Context, in *spec.FullInput, raw func(string, string, func([]byte) ([][]string, error)) ([][]string, []byte, error)) error {
	const path = "spec/inventory.tsv"
	data, err := s.Store.ReadFile(ctx, path)
	if err != nil {
		return &SourceError{"Index", "records:" + path, err}
	}
	if data == nil {
		return nil
	}
	index, err := spec.ReadInventoryIndex(data)
	if err == nil && len(index) == 0 {
		err = fmt.Errorf("the index has a header and no row")
	}
	if err != nil {
		return &SourceError{"Index", "records:" + path, fmt.Errorf("%s: %w", path, err)}
	}
	in.Index, in.Versions = index, map[int][]byte{}
	for _, r := range index {
		n, err := strconv.Atoi(r[0])
		if err != nil {
			return &SourceError{"Index", "records:" + path, err}
		}
		if _, in.Versions[n], err = raw("Versions", fmt.Sprintf("spec/inventory-%d.tsv", n), nil); err != nil {
			return err
		}
	}
	return nil
}

// specTree sets the fields of in that the tree at head holds, read in the
// run's clone: the changed paths, the specification and PRD files, and the
// record of the problem statement (one numbered facts file in docs/facts/).
func specTree(clone, base, head string, in *spec.FullInput) error {
	fail := func(field, what string, err error) error {
		return &SourceError{field, "target:" + what + "@" + head, err}
	}
	var err error
	if in.Changed, err = git.DiffNames(clone, base, head); err != nil {
		return fail("Changed", "diff", err)
	}
	if in.SpecFiles, err = specFilesAt(clone, head); err != nil {
		return fail("SpecFiles", "docs/spec", err)
	}
	if in.PRDs, err = filesAt(clone, head, "docs/prd", func(b string) bool { return prdName.MatchString(b) }); err != nil {
		return fail("PRDs", "docs/prd", err)
	}
	facts, err := filesAt(clone, head, "docs/facts", isNumberedFacts)
	if err == nil {
		in.PSBRecords, err = psbRecordsOf(facts)
	}
	if err != nil {
		return fail("PSBRecords", "docs/facts", err)
	}
	return nil
}

// prdName is the name of a PRD file, the selector of the full check (PRD-*.md).
var prdName = regexp.MustCompile(`^PRD-.*\.md$`)
