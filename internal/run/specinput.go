package run

// The input of the full spec check at a head (docs/spec/loop.md, REQ-012, The
// full spec check; task T-cav3, #412, row 66p of the plan of M2e): each field
// of spec.FullInput is built from one named record or from the tree at the head.

import (
	"context"

	"github.com/pharzam/layup/internal/spec"
)

// SourceError is a refusal of the input of the full check. Source is "base" or
// "head" for a SHA that is not a full SHA (a caller's error), "records:<path>"
// for a record, or "target:<what>@<sha>" for the tree of the run's clone.
type SourceError struct {
	Field, Source string
	Err           error
}

func (e *SourceError) Error() string { return "" }
func (e *SourceError) Unwrap() error { return e.Err }

// specInput gives the FullInput of head against base: the records of the run's
// records base at the call, the tree of the run's clone, and the caller's copy.
// It writes nothing, and on any error gives the zero input.
func (s *Sessions) specInput(ctx context.Context, base, head string, copy []byte) (spec.FullInput, error) {
	return spec.FullInput{}, nil
}

// specRecords sets the fields of in that the records hold.
func (s *Sessions) specRecords(ctx context.Context, in *spec.FullInput) error { return nil }

// specTree sets the fields of in that the tree at head holds.
func specTree(clone, base, head string, in *spec.FullInput) error { return nil }
