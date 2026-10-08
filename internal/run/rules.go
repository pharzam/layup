package run

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/pharzam/layup/internal/forge"
)

// The roles that the rules of M2c name (docs/spec/run.md, A human decision).
var (
	IntakeRoles     = []string{"operator", "idea-owner"}
	AcceptanceRoles = []string{"idea-owner"}
)

// Decision says whether a row of copies.tsv is a human decision for a rule
// that names roles: the first copy of the comment (seen 1; an edit is never a
// decision, and a comment edited before its first copy is that copy), made by
// no App (app —), by an author ID that approvers.tsv holds in one of the roles.
// A review, a review comment, a commit and a reaction are never rows of
// copies.tsv.
func Decision(copy []string, approvers [][]string, roles ...string) bool {
	if copy[1] != "1" || copy[5] != "—" {
		return false
	}
	for _, a := range approvers {
		for _, role := range roles {
			if a[0] == copy[3] && a[1] == role {
				return true
			}
		}
	}
	return false
}

// Copy gives the row of copies.tsv and the body file of a comment of issue,
// copied at now, or false when the SHA-256 of its body equals that of its last
// copy (docs/spec/run.md, Copy before read). The first copy is seen 1, an edit
// the next seen; created is the forge's time of the last edit, else of the
// comment.
func Copy(copies [][]string, c forge.Comment, issue int, now time.Time) ([]string, []byte, bool) {
	id := strconv.FormatInt(c.ID, 10)
	sum := sha256.Sum256([]byte(c.Body))
	digest := hex.EncodeToString(sum[:])
	seen := 0
	var last []string
	for _, r := range copies {
		if r[0] == id {
			if n, _ := strconv.Atoi(r[1]); n > seen {
				seen, last = n, r
			}
		}
	}
	if last != nil && last[8] == digest {
		return nil, nil, false
	}
	app, created := "—", c.Updated
	if c.App != "" {
		app = c.App
	}
	if created.IsZero() {
		created = c.Created
	}
	stamp := func(t time.Time) string { return t.UTC().Format("2006-01-02T15:04:05Z") }
	row := []string{id, strconv.Itoa(seen + 1), strconv.Itoa(issue), strconv.FormatInt(c.AuthorID, 10), c.AuthorLogin,
		app, stamp(created), stamp(now), digest, fmt.Sprintf("copies/%s-%d.md", id, seen+1)}
	return row, []byte(c.Body), true
}
