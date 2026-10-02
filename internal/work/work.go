// Package work is the work area of one target on the LAYUP host
// (docs/spec/setup.md, The command layup setup): its paths, and the two records
// that layup setup and layup setup verify both read, each by its schema. It is
// the one home of the schemas that internal/setup and internal/verify share
// (K9 of the plan).
package work

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/pharzam/layup/internal/tsv"
)

// The paths of a work area, from its root.
const (
	AnswersPath  = "inputs/answers.tsv"
	RecordPath   = "out/record.tsv"
	CommandsPath = "out/commands.sh"
	TargetPath   = "target"
)

// AnswersSchema is the form of inputs/answers.tsv: the block setup-answers of
// docs/spec/setup.md.
var AnswersSchema = tsv.Schema{Name: "setup-answers", Location: "host:<work>/inputs/answers.tsv", Columns: []tsv.Column{
	{Name: "question", Type: "text", Key: true},
	{Name: "answer", Type: "text"},
	{Name: "by", Type: "enum(operator|idea-owner)"},
	{Name: "source", Type: "text"},
	{Name: "question_text", Type: "text"},
}}

// RecordSchema is the form of out/record.tsv: the block setup-record.
var RecordSchema = tsv.Schema{Name: "setup-record", Location: "records:setup/record.tsv", Columns: []tsv.Column{
	{Name: "step", Type: "id(SNN)", Key: true},
	{Name: "name", Type: "text", Key: true},
	{Name: "value", Type: "text"},
	{Name: "source", Type: "enum(answer|catalog|fact|computed|gap|step)"},
	{Name: "ref", Type: "text"},
}}

// OpenGapsPath is the file of the open gaps of a target, from its root.
const OpenGapsPath = "docs/setup/open-gaps.tsv"

// OpenGapsSchema is the form of docs/setup/open-gaps.tsv of a target: the
// block open-gaps of docs/spec/setup.md, which S11 writes and the checks
// markers and sources read (#87).
var OpenGapsSchema = tsv.Schema{Name: "open-gaps", Location: "target:docs/setup/open-gaps.tsv", NoHeader: true, Columns: []tsv.Column{
	{Name: "file", Type: "path", Key: true},
	{Name: "marker", Type: "text", Key: true},
	{Name: "question", Type: "text"},
}}

// Answers is the rows of answers.tsv, one value per column.
type Answers [][]string

// A Record is the rows of record.tsv, one value per column.
type Record [][]string

// ParseAnswers reads the bytes of answers.tsv by AnswersSchema.
func ParseAnswers(data []byte) (Answers, error) {
	rows, err := tsv.Read(data, AnswersSchema)
	return Answers(rows), err
}

// ParseRecord reads the bytes of record.tsv by RecordSchema.
func ParseRecord(data []byte) (Record, error) {
	rows, err := tsv.Read(data, RecordSchema)
	return Record(rows), err
}

// ReadAnswers reads inputs/answers.tsv of the work area at dir. An error
// names the file from the root of the work area.
func ReadAnswers(dir string) (Answers, error) { return read(dir, AnswersPath, ParseAnswers) }

// ReadRecord reads out/record.tsv of the work area at dir.
func ReadRecord(dir string) (Record, error) { return read(dir, RecordPath, ParseRecord) }

// read reads the file at path of the work area at dir with parse.
func read[T any](dir, path string, parse func([]byte) (T, error)) (T, error) {
	var zero T
	data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(path)))
	if err != nil {
		return zero, fmt.Errorf("%s: %w", path, err)
	}
	v, err := parse(data)
	if err != nil {
		return zero, fmt.Errorf("%s: %w", path, err)
	}
	return v, nil
}

// Value gives the value of the row of step and name, and whether the record
// has that row. The key of the record is the step and the name, so there is
// at most one.
func (r Record) Value(step, name string) (string, bool) {
	for _, row := range r {
		if row[0] == step && row[1] == name {
			return row[2], true
		}
	}
	return "", false
}

// AnswersHash gives the record row "<step> answers.sha256" of a step that read
// the rows of a whose question has one of the prefixes (setup.md, The
// answers): the SHA-256 of those rows, sorted, each its fields joined by tabs,
// and the ref "sha256 <path> <prefix>…" of the record's rule.
func AnswersHash(step string, a Answers, prefixes []string) []string {
	var rows []string
	for _, r := range a {
		if slices.ContainsFunc(prefixes, func(p string) bool { return strings.HasPrefix(r[0], p) }) {
			rows = append(rows, strings.Join(r, "\t"))
		}
	}
	sort.Strings(rows)
	sum := sha256.Sum256([]byte(strings.Join(rows, "\n")))
	return []string{step, "answers.sha256", fmt.Sprintf("%x", sum), "computed",
		"sha256 " + AnswersPath + " " + strings.Join(prefixes, " ")}
}
