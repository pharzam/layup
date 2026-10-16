package handoff

import (
	"fmt"
	"slices"
)

// A Comment is an artifact that is a comment of the task's issue: the plan, the
// plan review or the review record. Its three hashes are the SHA-256 of the
// session's file, of the copy that layup run keeps in the records branch, and
// of the body that the forge gives for the comment at the step; the copy is
// the row ID-Seen of copies.tsv.
type Comment struct {
	ID, Seen          int
	File, Copy, Forge string
}

// An Artifact is one artifact of a handoff. Named is the hash that the session
// named for it, Found the hash that code found (the SHA-256 of the file at the
// head, or a commit's object name); the two are equal only when the artifact
// is what the session said. A comment artifact has a Comment and no hashes of
// its own.
type Artifact struct {
	Kind, Path   string
	Named, Found string
	Comment      *Comment
}

// A Handoff is what the caller (internal/run) knows of one handoff, made by
// code and never from a field that an agent writes: whether the typed result
// passed the block result (SchemaOK), the artifacts with the hashes that were
// named and found, and the result of each other item that needs a program, a
// gate run or a read of a commit (Items). An item with no result is false.
// CodeWritten is true for a handoff that code writes itself and that has no
// session result (the rendered task): SchemaOK is not read for it, and a row
// that lists "schema" cannot be satisfied by it.
type Handoff struct {
	From, To    string
	CodeWritten bool
	SchemaOK    bool
	Artifacts   []Artifact
	Items       map[string]bool
}

// A Verdict is the result of a check: Detail is the detail of the event
// handoff of tasks/<task>/events.tsv.
type Verdict struct {
	Valid  bool
	Detail string
}

// commentKind reports whether an artifact of the kind is a comment.
func commentKind(kind string) bool {
	return kind == "plan" || kind == "plan-review" || kind == "review-record"
}

// hexHash reports whether h is lowercase hexadecimal of the given length.
func hexHash(h string, n int) bool {
	if len(h) != n {
		return false
	}
	for _, c := range h {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}

// artifactOK reports whether one artifact is what its kind says. A comment
// artifact (a plan, a plan review or a review record) has its Comment, whose
// ID and Seen are 1 or more and whose three hashes are one SHA-256; any other
// artifact has no Comment, and the hash that the session named is the one that
// code found: a SHA-256 of a file, or, for the commits, an object name of 40 or
// 64 characters.
func artifactOK(a Artifact) bool {
	if commentKind(a.Kind) {
		c := a.Comment
		return c != nil && c.ID >= 1 && c.Seen >= 1 && hexHash(c.File, 64) && c.File == c.Copy && c.Copy == c.Forge
	}
	if a.Comment != nil || a.Named != a.Found {
		return false
	}
	if a.Kind == "commits" {
		return hexHash(a.Found, 40) || hexHash(a.Found, 64)
	}
	return hexHash(a.Found, 64)
}

// artifactsOK checks the item "artifacts": each artifact that the handoff
// names, those beyond the kinds of the row too, is a known kind and is what it
// says; and, when the row lists the item, each kind of the row has one.
func artifactsOK(row Row, h Handoff) bool {
	for _, a := range h.Artifacts {
		if !slices.Contains(kinds, a.Kind) || !artifactOK(a) {
			return false
		}
	}
	if !slices.Contains(row.Items, "artifacts") {
		return true
	}
	for _, kind := range row.Artifacts {
		if !slices.ContainsFunc(h.Artifacts, func(a Artifact) bool { return a.Kind == kind }) {
			return false
		}
	}
	return true
}

// Check checks a handoff against a table. The pair must be a row of the table
// ("invalid pair" otherwise). Then, whatever the row lists, the schema (the
// caller's result of the check of the block result; a CodeWritten handoff has
// none, so it is invalid on a row that lists "schema") and each artifact that
// the handoff names are checked, in that order, so that a replacement row that omits "schema" or "artifacts" cannot
// let a bad result or a bad artifact through; and then each item of the row in
// the row's order. "artifacts" is computed here from the hashes, and the map
// Items is never read for "schema" or "artifacts", so a caller's claim for them
// cannot make a handoff valid; no field of Handoff is a verdict. The first
// failure names the detail.
func Check(table []Row, h Handoff) Verdict {
	suffix := ""
	for _, a := range h.Artifacts {
		if c := a.Comment; c != nil && commentKind(a.Kind) && c.ID >= 1 && c.Seen >= 1 {
			suffix = fmt.Sprintf(" comment %d-%d", c.ID, c.Seen)
			break
		}
	}
	i := slices.IndexFunc(table, func(r Row) bool { return r.From == h.From && r.To == h.To })
	if i < 0 {
		return Verdict{Detail: "invalid pair" + suffix}
	}
	row := table[i]
	if (h.CodeWritten && slices.Contains(row.Items, "schema")) || (!h.CodeWritten && !h.SchemaOK) {
		return Verdict{Detail: "invalid schema" + suffix}
	}
	if !artifactsOK(row, h) {
		return Verdict{Detail: "invalid artifacts" + suffix}
	}
	for _, item := range row.Items {
		if item != "schema" && item != "artifacts" && !h.Items[item] {
			return Verdict{Detail: "invalid " + item + suffix}
		}
	}
	return Verdict{Valid: true, Detail: "valid" + suffix}
}
