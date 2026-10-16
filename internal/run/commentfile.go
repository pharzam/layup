package run

// This file holds the end check of a comment row and the kept copy of its file
// (docs/spec/loop.md, The artifact that is a comment; task T-q5k4, #360, row
// 66d of the plan of M2e).

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/pharzam/layup/internal/session"
)

// MaxCommentChars is the size limit of a comment's body, in characters, that
// loop.md names (GitHub's own limit).
const MaxCommentChars = 65536

// A HandoffRefusal is a file that cannot be posted as a comment: the reason
// of the event refused, handoff, and what is wrong with the file.
type HandoffRefusal struct{ Reason, Value string }

func (e HandoffRefusal) Error() string { return e.Reason + ": " + e.Value }

// KeptPath is the path in the records of the copy of a comment row's file:
// tasks/<task>/results/<session>/<file> (docs/spec/records.md, the records
// branch).
func KeptPath(task, session, file string) string {
	return "tasks/" + task + "/results/" + session + "/" + file
}

// CheckCommentBody refuses a file that cannot be the body of a comment: an
// empty one (GitHub refuses it), one that is not valid UTF-8 (JSON would
// replace the bad bytes, and the three hashes of the handoff would then
// differ), and one over MaxCommentChars characters (not bytes). The refusal is
// a HandoffRefusal, the reason of the event refused (loop.md).
func CheckCommentBody(b []byte) error {
	switch {
	case len(b) == 0:
		return HandoffRefusal{"handoff", "the file of a comment is empty"}
	case !utf8.Valid(b):
		return HandoffRefusal{"handoff", "the file of a comment is not valid UTF-8"}
	case utf8.RuneCount(b) > MaxCommentChars:
		return HandoffRefusal{"handoff", fmt.Sprintf("the file of a comment has %d characters, over %d", utf8.RuneCount(b), MaxCommentChars)}
	}
	return nil
}

// keepComments checks each comment row of a result against its file in the
// session's result/ directory and gives the files to keep, by their path in the
// records, byte for byte; bad is true when the file of a row is missing, is not
// a regular file or has another SHA-256 than the row's, and that row's file is
// not kept. The artifact rows are checked against the head, not here.
func keepComments(d session.Dir, task, id string, rows [][]string) (map[string][]byte, bool) {
	kept, bad := map[string][]byte{}, false
	for _, r := range rows {
		if r[0] != "comment" {
			continue
		}
		file := strings.TrimPrefix(r[2], "result/")
		data, err := session.CommentFile(d, file)
		if sum := sha256.Sum256(data); err != nil || hex.EncodeToString(sum[:]) != r[3] {
			bad = true
			continue
		}
		kept[KeptPath(task, id, file)] = data
	}
	return kept, bad
}
