package run

// This file holds the composition of the question and the answer comment of
// a session (docs/spec/loop.md, The artifact that is a comment; task T-gyvo,
// #405, row 66o of the plan of M2e). Composition only: the post is row 66k's.

import "errors"

var errNotBuilt = errors.New("not built")

// QuestionComment is the body of the comment of a question row of a result.
func QuestionComment(start, row []string) ([]byte, error) { return nil, errNotBuilt }

// AnswerComment is the body of the comment of an answer row of a result.
func AnswerComment(start, row []string, questionID string) ([]byte, error) {
	return nil, errNotBuilt
}
