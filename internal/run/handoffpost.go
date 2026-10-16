package run

// This file holds the post of a comment that is an artifact, its copy, and the
// record of a handoff (docs/spec/loop.md, The artifact that is a comment and
// The checks of a handoff; task T-q5k4, #360, row 66d of the plan of M2e). The
// loop step that calls them is row 66k's.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/handoff"
	"github.com/pharzam/layup/internal/records"
)

// Posted is a comment that holds an artifact: its ID, the number of its copy in
// copies.tsv and the SHA-256 of the file.
type Posted struct {
	ID   int64
	Seen int
	SHA  string
}

// digest is the SHA-256 of the exact bytes, in lowercase hexadecimal.
func digest(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

// The columns of a row of copies.tsv, as records.ReadCopies gives it.
const (
	copyComment = iota
	copySeen
	copyIssue
	copyAuthorID
	copyAuthorLogin
	copyApp
	copyCreated
	copyCopied
	copySHA
	copyBody
)

// copyNamed gives the last row of copies.tsv, in the order of the file, of a
// comment of issue whose body has the SHA-256 sum (loop.md, "The link is
// written by code"): the copy that a handoff of a file with that hash names.
func copyNamed(copies [][]string, issue int, sum string) ([]string, bool) {
	for i := len(copies) - 1; i >= 0; i-- {
		if copies[i][copySHA] == sum && copies[i][copyIssue] == strconv.Itoa(issue) {
			return copies[i], true
		}
	}
	return nil, false
}

// postedOf reads a row of copies.tsv as the comment it copies.
func postedOf(row []string) Posted {
	id, _ := strconv.ParseInt(row[copyComment], 10, 64)
	seen, _ := strconv.Atoi(row[copySeen])
	return Posted{ID: id, Seen: seen, SHA: row[copySHA]}
}

// taskIssue gives the issue of a task from the register.
func (s *Sessions) taskIssue(ctx context.Context, task string) (int, error) {
	data, err := s.Store.ReadFile(ctx, tasksPath)
	if err != nil {
		return 0, err
	}
	if data == nil {
		return 0, fmt.Errorf("the records have no %s", tasksPath)
	}
	tasks, err := ReadRegister(data)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", tasksPath, err)
	}
	for _, t := range tasks {
		if t.ID == task {
			return t.Issue, nil
		}
	}
	return 0, fmt.Errorf("%s has no row of %s", tasksPath, task)
}

// refuseFile writes the event refused, handoff, of a session whose file cannot
// be posted, once (a second call finds it).
func (s *Sessions) refuseFile(ctx context.Context, task, session string) error {
	start, err := s.startRowOf(ctx, session)
	if err != nil {
		return err
	}
	before, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return err
	}
	for _, e := range before {
		if e[1] == "refused" && e[3] == session && e[6] == "handoff" {
			return nil
		}
	}
	attempt, _ := strconv.Atoi(start[sessionAttempt])
	data, err := s.appendEvents(ctx, task, attempt, session, []string{"refused handoff"})
	if err != nil {
		return err
	}
	return s.commit(ctx, map[string][]byte{eventsPath(task): data}, fmt.Sprintf("layup run: the file of %s cannot be a comment", session))
}

// copyFiles gives the files of the copy of c: copies.tsv with its row and the
// body, or nil when the last copy of the comment has this body already.
func (s *Sessions) copyFiles(copies [][]string, c forge.Comment, issue int) (map[string][]byte, []string, error) {
	row, body, ok := Copy(copies, c, issue, s.Now())
	if !ok {
		return nil, nil, nil
	}
	data, err := table(records.CopiesSchema, append(slices.Clone(copies), row))
	if err != nil {
		return nil, nil, err
	}
	return map[string][]byte{"copies.tsv": data, row[copyBody]: body}, row, nil
}

// PostComment posts the file of a comment row of a session as a comment of the
// task's issue, from its kept copy in the records (the session directory is
// gone), and copies the comment by the mechanism of copy before read. A file
// that cannot be a comment is refused before any forge call (the event refused,
// handoff, once). A comment that has a copy of these bytes on the issue, or that
// the App has already posted with these bytes, is not posted again; the one post
// is fenced by a read of the lease, and the copy by Fenced.
func (s *Sessions) PostComment(ctx context.Context, task, session, file string) (Posted, error) {
	s.fns()
	path := KeptPath(task, session, file)
	data, err := s.Store.ReadFile(ctx, path)
	if err != nil {
		return Posted{}, err
	}
	if data == nil {
		return Posted{}, fmt.Errorf("the records have no kept file %s", path)
	}
	if err := CheckCommentBody(data); err != nil {
		var refusal HandoffRefusal
		if errors.As(err, &refusal) {
			if rErr := s.refuseFile(ctx, task, session); rErr != nil {
				return Posted{}, errors.Join(err, rErr)
			}
		}
		return Posted{}, err
	}
	issue, err := s.taskIssue(ctx, task)
	if err != nil {
		return Posted{}, err
	}
	sum := digest(data)
	copies, err := s.readTable(ctx, "copies.tsv", records.ReadCopies)
	if err != nil {
		return Posted{}, err
	}
	if row, ok := copyNamed(copies, issue, sum); ok {
		return postedOf(row), nil
	}
	comments, err := s.Forge.Comments(ctx, issue)
	if err != nil {
		return Posted{}, err
	}
	for _, c := range comments {
		if s.App != "" && c.App == s.App && c.Body == string(data) {
			return s.copyComment(ctx, c, issue, copies, sum)
		}
	}
	lease, err := s.Store.ReadLease(ctx)
	if err != nil {
		return Posted{}, err
	}
	if lease.Run != s.RunID || lease.State != "held" {
		return Posted{}, &LostError{Holder: lease.Run}
	}
	id, err := s.Forge.Comment(ctx, issue, string(data))
	if err != nil {
		return Posted{}, err
	}
	if comments, err = s.Forge.Comments(ctx, issue); err != nil {
		return Posted{}, err
	}
	for _, c := range comments {
		if c.ID == id {
			return s.copyComment(ctx, c, issue, copies, sum)
		}
	}
	return Posted{}, fmt.Errorf("the forge gave the ID %d for the comment and does not list it", id)
}

// copyComment copies c in one fenced records commit and gives the copy that has
// the SHA-256 sum.
func (s *Sessions) copyComment(ctx context.Context, c forge.Comment, issue int, copies [][]string, sum string) (Posted, error) {
	files, row, err := s.copyFiles(copies, c, issue)
	if err != nil {
		return Posted{}, err
	}
	if files == nil {
		if row, ok := copyNamed(copies, issue, sum); ok {
			return postedOf(row), nil
		}
		return Posted{}, fmt.Errorf("the last copy of the comment %d is not of these bytes", c.ID)
	}
	if err := s.commit(ctx, files, "records: the copy of the comment "+row[copyComment]); err != nil {
		return Posted{}, err
	}
	return postedOf(row), nil
}

// A HandoffInput is what the caller knows of one handoff: the task, the session
// whose result it is (none for a handoff that code writes, which has an Attempt
// instead), the pair, the comment file and its kind when the artifact is a
// comment, the other artifacts that the caller found, and the results of the
// items that need a program, a gate run or a read of a commit. An item the
// caller does not pass is false; schema, artifacts, question-row and
// answer-row are computed here, and a value of the caller's for one of them is
// never read. Freeze is the list of tests of a handoff to the verifier with their
// results (task T-6zj3, #361): when it is given, tests-red-green is computed
// from it too, and the first valid handoff writes tasks/<task>/tests.tsv in the
// commit of its event.
type HandoffInput struct {
	Task, Session string
	Attempt       int
	From, To      string
	File, Kind    string
	Artifacts     []handoff.Artifact
	Items         map[string]bool
	CodeWritten   bool
	Freeze        *Freeze
}

// RecordHandoff checks a handoff against the table of the records and writes
// its event, with the copy of the comment when a read finds it edited, in one
// fenced records commit. The schema is the records' result file read by its
// block; the comment of a file is the copy whose SHA-256 is the file's, and its
// three hashes are the file's, the copy's and the body that the forge gives at
// this step. A handoff already recorded with the same outcome is not written
// again (the verdict is returned).
func (s *Sessions) RecordHandoff(ctx context.Context, in HandoffInput) (handoff.Verdict, error) {
	s.fns()
	attempt := in.Attempt
	var result [][]string
	schemaOK := false
	if !in.CodeWritten {
		start, err := s.startRowOf(ctx, in.Session)
		if err != nil {
			return handoff.Verdict{}, err
		}
		if start[sessionTask] != in.Task {
			return handoff.Verdict{}, fmt.Errorf("the session %s is of %s, not of %s", in.Session, start[sessionTask], in.Task)
		}
		attempt, _ = strconv.Atoi(start[sessionAttempt])
		data, err := s.Store.ReadFile(ctx, "tasks/"+in.Task+"/results/"+in.Session+".tsv")
		if err != nil {
			return handoff.Verdict{}, err
		}
		if data != nil {
			if rows, err := records.ReadResult(data); err == nil {
				result, schemaOK = rows, true
			}
		}
	}
	if attempt < 1 {
		return handoff.Verdict{}, fmt.Errorf("the handoff of %s has no attempt", in.Task)
	}
	tdata, err := s.Store.ReadFile(ctx, transitionsPath)
	if err != nil {
		return handoff.Verdict{}, err
	}
	if tdata == nil {
		return handoff.Verdict{}, fmt.Errorf("the records have no %s", transitionsPath)
	}
	rows, err := handoff.ReadTable(tdata)
	if err == nil {
		err = handoff.CheckTable(rows)
	}
	if err != nil {
		return handoff.Verdict{}, fmt.Errorf("%s: %w", transitionsPath, err)
	}
	items := maps.Clone(in.Items)
	if items == nil {
		items = map[string]bool{}
	}
	items["question-row"], items["answer-row"] = QuestionRowOK(result), AnswerRowOK(result)
	// A list of tests at a handoff to the verifier: tests-red-green is computed
	// from the typed results (the caller's own value is never read), and the list
	// is written in the commit of the event (loop.md, The frozen test list).
	var frozen freezing
	if in.Freeze != nil && !in.CodeWritten && in.To == "verifier" {
		var err error
		if frozen, err = s.freeze(ctx, in); err != nil {
			return handoff.Verdict{}, err
		}
		items["tests-red-green"] = frozen.redGreen
	}
	h := handoff.Handoff{From: in.From, To: in.To, CodeWritten: in.CodeWritten, SchemaOK: schemaOK, Artifacts: slices.Clone(in.Artifacts), Items: items}
	files := map[string][]byte{}
	if in.File != "" {
		art, copyFiles, err := s.commentArtifact(ctx, in, result)
		if err != nil {
			return handoff.Verdict{}, err
		}
		if art != nil {
			h.Artifacts = append(h.Artifacts, *art)
		}
		maps.Copy(files, copyFiles)
	}
	v := handoff.Check(rows, h)
	before, err := s.readTable(ctx, eventsPath(in.Task), records.ReadEvents)
	if err != nil {
		return handoff.Verdict{}, err
	}
	for _, e := range before {
		if e[1] == "handoff" && e[3] == in.Session && e[6] == v.Detail {
			if len(files) > 0 {
				return v, s.commit(ctx, files, "records: the copy of the comment of "+in.Task)
			}
			return v, nil
		}
	}
	data, err := s.appendEvents(ctx, in.Task, attempt, in.Session, []string{"handoff " + v.Detail})
	if err != nil {
		return handoff.Verdict{}, err
	}
	files[eventsPath(in.Task)] = data
	if shouldFreeze(in, v) && frozen.rows != nil {
		if files[testsPath(in.Task)], err = WriteTests(frozen.rows); err != nil {
			return handoff.Verdict{}, err
		}
	}
	return v, s.commit(ctx, files, fmt.Sprintf("layup run: the handoff %s to %s of %s, %s", in.From, in.To, in.Task, v.Detail))
}

// commentArtifact gives the comment artifact of in.File and the files of the
// copy that reading the forge adds. A file that the result does not name as a
// comment row is no artifact (nil); a comment with no copy has the ID 0.
func (s *Sessions) commentArtifact(ctx context.Context, in HandoffInput, result [][]string) (*handoff.Artifact, map[string][]byte, error) {
	named := ""
	for _, r := range result {
		if r[0] == "comment" && r[2] == "result/"+in.File {
			named = r[3]
		}
	}
	if named == "" {
		return nil, nil, nil
	}
	issue, err := s.taskIssue(ctx, in.Task)
	if err != nil {
		return nil, nil, err
	}
	copies, err := s.readTable(ctx, "copies.tsv", records.ReadCopies)
	if err != nil {
		return nil, nil, err
	}
	art := &handoff.Artifact{Kind: in.Kind, Comment: &handoff.Comment{File: named}}
	row, ok := copyNamed(copies, issue, named)
	if !ok {
		return art, nil, nil
	}
	id := postedOf(row).ID
	comments, err := s.Forge.Comments(ctx, issue)
	if err != nil {
		return nil, nil, err
	}
	var files map[string][]byte
	for _, c := range comments {
		if c.ID != id {
			continue
		}
		var newRow []string
		if files, newRow, err = s.copyFiles(copies, c, issue); err != nil {
			return nil, nil, err
		}
		if newRow != nil {
			copies = append(copies, newRow)
			row, _ = copyNamed(copies, issue, named)
		}
		art.Comment.Forge = digest([]byte(c.Body))
	}
	p := postedOf(row)
	art.Comment.ID, art.Comment.Seen, art.Comment.Copy = int(p.ID), p.Seen, row[copySHA]
	return art, files, nil
}
