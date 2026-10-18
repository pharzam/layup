package run

// This file holds the composition of the question and the answer comment of
// a session (docs/spec/loop.md, The artifact that is a comment; task T-gyvo,
// #405, row 66o of the plan of M2e). It composes only: the post is row 66k's.

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/tsv"
)

// QuestionComment is the body of the comment of a question row of a result:
// the first line, a blank line, and the row's text, with no final newline. The
// ID is the session's and the row's n, never its text or reason. start is the
// session's row of sessions.tsv and row a row of its result, as the records'
// readers give them; either off its schema, its types or its rules, a row of
// another kind, or a body CheckCommentBody refuses (a HandoffRefusal) gives an
// error and no bytes.
func QuestionComment(start, row []string) ([]byte, error) {
	id, err := composedID(start, row, "question")
	if err != nil {
		return nil, err
	}
	return body(start, "question "+id, row)
}

// AnswerComment is the body of the comment of an answer row of a result, as
// QuestionComment is of a question row. questionID, the ID of the question that
// the answer answers, is the caller's: it is checked by its form and neither
// read from the row's reason nor resolved against the records.
func AnswerComment(start, row []string, questionID string) ([]byte, error) {
	id, err := composedID(start, row, "answer")
	if err != nil {
		return nil, err
	}
	if !answerForm.MatchString(questionID) {
		return nil, fmt.Errorf("the question ID %q is not S- and 8 lowercase hexadecimal characters, a dot and the n of the row", questionID)
	}
	return body(start, "answer "+id+" of "+questionID, row)
}

// composedID checks the start row and the row of kind and gives the ID
// <session>.<n> of the row.
func composedID(start, row []string, kind string) (string, error) {
	if err := readerForm(records.SessionsSchema, start, records.CheckSession); err != nil {
		return "", fmt.Errorf("the start row: %w", err)
	}
	// Text that is not valid UTF-8 is a body that CheckCommentBody refuses (a
	// HandoffRefusal), before the reader's schema check would call it a bad field.
	if len(row) > 2 && !utf8.ValidString(row[2]) {
		return "", CheckCommentBody([]byte(row[2]))
	}
	if err := readerForm(records.ResultSchema, row, records.CheckResultRow); err != nil {
		return "", fmt.Errorf("the result row: %w", err)
	}
	if row[0] != kind {
		return "", fmt.Errorf("the result row is a %s row, not a %s row", row[0], kind)
	}
	id := start[sessionID] + "." + row[1]
	if !answerForm.MatchString(id) {
		return "", fmt.Errorf("the ID %q is not a session and an n of 1 or more", id)
	}
	return id, nil
}

// readerForm checks that row is what the records' reader gives for one row of
// the schema s: it reads the row as a table of its own, with the reader's
// schema and type checks, then applies the rules of the row. It repairs
// nothing: a field the reader would read as another value (— is the empty
// value) is refused.
func readerForm(s tsv.Schema, row []string, rules func([]string) error) error {
	if len(row) != len(s.Columns) {
		return fmt.Errorf("%d columns; the schema %s has %d", len(row), s.Name, len(s.Columns))
	}
	names := make([]string, len(row))
	fields := slices.Clone(row)
	for i, c := range s.Columns {
		names[i] = c.Name
		if fields[i] == "" {
			fields[i] = "—"
		}
	}
	data := strings.Join(names, "\t") + "\n" + strings.Join(fields, "\t") + "\n"
	rows, err := tsv.Read([]byte(data), s)
	switch {
	case err != nil:
		return err
	case len(rows) != 1 || !slices.Equal(rows[0], row):
		return errors.New("the row is not what the reader gives for it")
	}
	return rules(row)
}

// body composes the comment: the first line, a blank line and the row's text,
// checked by CheckCommentBody.
func body(start []string, what string, row []string) ([]byte, error) {
	b := []byte(firstLinePrefix(start) + what + "\n\n" + row[2])
	if err := CheckCommentBody(b); err != nil {
		return nil, err
	}
	return b, nil
}
