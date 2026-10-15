package handoff

import (
	"bytes"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
)

// h64 is a SHA-256 of one letter repeated, and h40 a commit.
func h64(c string) string { return strings.Repeat(c, 64) }
func h40(c string) string { return strings.Repeat(c, 40) }

// The thirteen default rows, as docs/spec/loop.md gives them (the integration
// test reads them from the specification).
var wantDefault = []Row{
	{"task plan", "plan review", []string{"plan"}, []string{"schema", "artifacts", "plan-fields"}},
	{"plan review", "developer", []string{"plan-review"}, []string{"schema", "artifacts", "verdict-field"}},
	{"developer (build task)", "verifier", []string{"commits", "task-file"}, []string{"schema", "artifacts", "tests-red-green"}},
	{"developer (specification)", "verifier", []string{"commits", "task-file"}, []string{"schema", "artifacts", "spec-sections"}},
	{"developer (architecture)", "verifier", []string{"commits", "task-file"}, []string{"schema", "artifacts"}},
	{"developer (rule batch)", "verifier", []string{"commits", "task-file", "patch"}, []string{"schema", "artifacts", "patch-per-kind"}},
	{"developer (milestone plan)", "verifier", []string{"rows"}, []string{"schema", "artifacts", "register-rows"}},
	{"code (rendered task)", "verifier", []string{"rendering"}, []string{"artifacts", "rendering-equal"}},
	{"verifier", "close-out", []string{"review-record"}, []string{"schema", "artifacts", "nothing-material"}},
	{"verifier", "developer", []string{"review-record"}, []string{"schema", "artifacts", "findings-located"}},
	{"specification", "milestone plan", []string{"task-file"}, []string{"artifacts", "heading-ids"}},
	{"asking role", "owner role", []string{"question"}, []string{"schema", "question-row"}},
	{"owner role", "asking role", []string{"answer"}, []string{"schema", "answer-row"}},
}

func TestTheDefaultRowsAreTheThirteenOfTheSpecification(t *testing.T) {
	got := Default()
	if len(got) != 13 {
		t.Fatalf("%d default rows, want 13", len(got))
	}
	for i, w := range wantDefault {
		if got[i].From != w.From || got[i].To != w.To || !slices.Equal(got[i].Artifacts, w.Artifacts) || !slices.Equal(got[i].Items, w.Items) {
			t.Errorf("row %d: %+v, want %+v", i+1, got[i], w)
		}
	}
	if err := CheckTable(got); err != nil {
		t.Errorf("the default table fails its own check: %v", err)
	}
}

// A caller may change what Default returns without changing the next call.
func TestDefaultGivesACopy(t *testing.T) {
	a := Default()
	a[0].Items[0] = "changed"
	a[0].From = "changed"
	if b := Default(); b[0].From != "task plan" || b[0].Items[0] != "schema" {
		t.Fatalf("Default shares its rows: %+v", b[0])
	}
}

// The lists are copies: a caller that changes one does not change the next.
func TestTheListsAreCopies(t *testing.T) {
	k, it := Kinds(), Items()
	k[0], it[0] = "changed", "changed"
	if Kinds()[0] != "plan" || Items()[0] != "schema" {
		t.Fatal("Kinds or Items shares its list")
	}
}

func TestTheListsAreTheWordsOfTheSpecification(t *testing.T) {
	if got, want := Kinds(), []string{"plan", "plan-review", "review-record", "task-file", "commits", "rows", "patch", "rendering", "question", "answer"}; !slices.Equal(got, want) {
		t.Errorf("kinds %v, want %v", got, want)
	}
	if got, want := Items(), []string{"schema", "artifacts", "plan-fields", "verdict-field", "tests-red-green", "spec-sections", "patch-per-kind", "register-rows", "rendering-equal", "nothing-material", "findings-located", "heading-ids", "question-row", "answer-row"}; !slices.Equal(got, want) {
		t.Errorf("items %v, want %v", got, want)
	}
}

// The list of items cannot drift from the words that the block events takes in
// the detail of a handoff (records.CheckEvent).
func TestEachItemIsAWordOfTheDetailOfAnEvent(t *testing.T) {
	event := func(detail string) error {
		return records.CheckEvent([]string{"1", "handoff", "1", "S-12345678", "", "", detail, "2026-10-15T12:00:00Z"})
	}
	for _, item := range append(Items(), "pair") {
		if err := event("invalid " + item); err != nil {
			t.Errorf("invalid %s: %v", item, err)
		}
	}
	if err := event("invalid not-an-item"); err == nil {
		t.Error("an item outside the table is accepted by the event")
	}
}

func TestWriteThenReadGivesTheDefaultTable(t *testing.T) {
	var b bytes.Buffer
	if err := WriteTable(&b, Default()); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(b.String(), "from\tto\tartifacts\titems\n") {
		t.Fatalf("header %q", b.String())
	}
	got, err := ReadTable(b.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 13 || got[3].From != "developer (specification)" || !slices.Equal(got[3].Items, []string{"schema", "artifacts", "spec-sections"}) {
		t.Fatalf("read back %+v", got)
	}
}

func tableText(rows ...string) []byte {
	return []byte("from\tto\tartifacts\titems\n" + strings.Join(rows, "\n") + "\n")
}

func TestReadTableRefuses(t *testing.T) {
	for name, tc := range map[string]struct{ data, why string }{
		"an unknown kind":             {string(tableText("a\tb\tplan blueprint\tschema")), "line 2: column artifacts"},
		"an unknown item":             {string(tableText("a\tb\tplan\tschema vibes")), "line 2: column items"},
		"the word pair as an item":    {string(tableText("a\tb\tplan\tschema pair")), "line 2: column items"},
		"no artifact":                 {string(tableText("a\tb\t—\tschema")), "line 2: column artifacts"},
		"no item":                     {string(tableText("a\tb\tplan\t—")), "line 2: column items"},
		"a bad row after a good one":  {string(tableText("a\tb\tplan\tschema", "c\td\tplan\tvibes")), "line 3: column items"},
		"a pair twice":                {string(tableText("a\tb\tplan\tschema", "a\tb\tplan\tschema")), "line 3"},
		"a wrong header":              {"from\tto\titems\tartifacts\na\tb\tplan\tschema\n", "header"},
		"a row with too few fields":   {string(tableText("a\tb\tplan")), "line 2"},
		"a kind written with a comma": {string(tableText("a\tb\tplan,task-file\tschema")), "line 2"},
	} {
		_, err := ReadTable([]byte(tc.data))
		if err == nil || !strings.Contains(strings.ToLower(err.Error()), tc.why) {
			t.Errorf("%s: err %v, want one that names %q", name, err, tc.why)
		}
	}
}

func TestCheckTable(t *testing.T) {
	def := Default()
	without := func(i int) []Row { return slices.Concat(def[:i:i], def[i+1:]) }
	extra := append(slices.Clone(def), Row{"a", "b", []string{"plan"}, []string{"schema"}})
	for name, tc := range map[string]struct {
		rows []Row
		why  string // "" for a table that passes
	}{
		"the default":                           {def, ""},
		"an extra pair with words of the lists": {extra, ""},
		"a default pair missing":                {without(3), "developer (specification)"},
		"the last default pair missing":         {without(12), "owner role"},
		"an empty table":                        {nil, "task plan"},
		"a pair listed twice":                   {append(slices.Clone(def), def[0]), "twice"},
		"an unknown kind in an extra pair":      {append(slices.Clone(def), Row{"a", "b", []string{"blueprint"}, []string{"schema"}}), "blueprint"},
		"an unknown item in a default pair":     {func() []Row { r := Default(); r[0].Items = append(r[0].Items, "vibes"); return r }(), "vibes"},
		"a default pair with no items":          {func() []Row { r := Default(); r[0].Items = nil; return r }(), "items"},
		"a default pair with no artifacts":      {func() []Row { r := Default(); r[0].Artifacts = nil; return r }(), "artifacts"},
	} {
		err := CheckTable(tc.rows)
		switch {
		case tc.why == "" && err != nil:
			t.Errorf("%s: %v", name, err)
		case tc.why != "" && (err == nil || !strings.Contains(err.Error(), tc.why)):
			t.Errorf("%s: err %v, want one that names %q", name, err, tc.why)
		}
	}
}

// A handoff of the row "verifier, close-out" that is valid, and each way to
// break it.
func good() Handoff {
	return Handoff{
		From: "verifier", To: "close-out", SchemaOK: true,
		Artifacts: []Artifact{{Kind: "review-record", Comment: &Comment{ID: 12, Seen: 1, File: h64("a"), Copy: h64("a"), Forge: h64("a")}}},
		Items:     map[string]bool{"nothing-material": true},
	}
}

func TestACheckedHandoffIsValidOnlyByItsTable(t *testing.T) {
	table := Default()
	if v := Check(table, good()); !v.Valid || v.Detail != "valid comment 12-1" {
		t.Fatalf("a good handoff: %+v", v)
	}
	for name, tc := range map[string]struct {
		change func(*Handoff)
		detail string
	}{
		"a wrong schema":                                          {func(h *Handoff) { h.SchemaOK = false }, "invalid schema comment 12-1"},
		"an item that is false":                                   {func(h *Handoff) { h.Items["nothing-material"] = false }, "invalid nothing-material comment 12-1"},
		"an item with no result":                                  {func(h *Handoff) { delete(h.Items, "nothing-material") }, "invalid nothing-material comment 12-1"},
		"no artifact":                                             {func(h *Handoff) { h.Artifacts = nil }, "invalid artifacts"},
		"an artifact of another kind":                             {func(h *Handoff) { h.Artifacts[0].Kind = "plan" }, "invalid artifacts comment 12-1"},
		"a copy that is not the session's file":                   {func(h *Handoff) { h.Artifacts[0].Comment.Copy = h64("b") }, "invalid artifacts comment 12-1"},
		"a forge body that differs from the copy":                 {func(h *Handoff) { h.Artifacts[0].Comment.Forge = h64("b") }, "invalid artifacts comment 12-1"},
		"a session file that differs from the copy and the forge": {func(h *Handoff) { h.Artifacts[0].Comment.File = h64("b") }, "invalid artifacts comment 12-1"},
		"a comment hash that is not a SHA-256":                    {func(h *Handoff) { c := h.Artifacts[0].Comment; c.File, c.Copy, c.Forge = "x", "x", "x" }, "invalid artifacts comment 12-1"},
		"a comment artifact with no comment":                      {func(h *Handoff) { h.Artifacts[0].Comment = nil }, "invalid artifacts"},
		"a pair with no row":                                      {func(h *Handoff) { h.To = "merge" }, "invalid pair comment 12-1"},
		"a pair in the wrong direction":                           {func(h *Handoff) { h.From, h.To = h.To, h.From }, "invalid pair comment 12-1"},
	} {
		h := good()
		tc.change(&h)
		v := Check(table, h)
		want := tc.detail
		if name == "no artifact" || name == "a comment artifact with no comment" {
			want = "invalid artifacts"
		}
		if v.Valid || v.Detail != want {
			t.Errorf("%s: %+v, want invalid %q", name, v, want)
		}
	}
}

// The first failing item of the row's order names the detail, so a handoff
// whose schema and item both fail says schema.
func TestTheFirstFailingItemNamesTheDetail(t *testing.T) {
	h := good()
	h.SchemaOK = false
	h.Items["nothing-material"] = false
	h.Artifacts[0].Comment.Copy = h64("b")
	if v := Check(Default(), h); v.Detail != "invalid schema comment 12-1" {
		t.Fatalf("%+v", v)
	}
}

// A field that an agent writes never makes a handoff valid: the items
// "schema" and "artifacts" are computed, whatever the map says, and a hash that
// the session names is compared with the one that code found.
func TestWhatAnAgentWritesNeverMakesAHandoffValid(t *testing.T) {
	build := Handoff{
		From: "developer (build task)", To: "verifier", SchemaOK: true,
		Artifacts: []Artifact{{Kind: "commits", Path: "task/T-x/1", Named: h40("a"), Found: h40("a")}, {Kind: "task-file", Path: "docs/tasks/T-x.md", Named: h64("c"), Found: h64("c")}},
		Items:     map[string]bool{"tests-red-green": true},
	}
	if v := Check(Default(), build); !v.Valid || v.Detail != "valid" {
		t.Fatalf("a good build handoff: %+v", v)
	}
	lie := build
	lie.Artifacts = []Artifact{build.Artifacts[0], {Kind: "task-file", Path: "docs/tasks/T-x.md", Named: h64("c"), Found: h64("d")}}
	lie.Items = map[string]bool{"tests-red-green": true, "artifacts": true, "schema": true}
	if v := Check(Default(), lie); v.Valid || v.Detail != "invalid artifacts" {
		t.Errorf("a named hash that differs from the found one, with the item claimed true: %+v", v)
	}
	noFile := build
	noFile.Artifacts = []Artifact{build.Artifacts[0], {Kind: "task-file", Path: "docs/tasks/T-x.md", Named: h64("c")}}
	if v := Check(Default(), noFile); v.Valid || v.Detail != "invalid artifacts" {
		t.Errorf("a path that code did not find: %+v", v)
	}
	noSchema := build
	noSchema.SchemaOK = false
	noSchema.Items = map[string]bool{"schema": true, "tests-red-green": true}
	if v := Check(Default(), noSchema); v.Valid || v.Detail != "invalid schema" {
		t.Errorf("schema claimed true in the map: %+v", v)
	}
	short := build
	short.Artifacts = []Artifact{{Kind: "commits", Path: "task/T-x/1", Named: "abc", Found: "abc"}, build.Artifacts[1]}
	if v := Check(Default(), short); v.Valid || v.Detail != "invalid artifacts" {
		t.Errorf("a hash that is not a full object name: %+v", v)
	}
}

// Each pair of the default table is checked with its own kinds and items.
func TestEachDefaultPairNeedsItsOwnItemsAndKinds(t *testing.T) {
	if len(Default()) != 13 {
		t.Fatalf("%d default rows, want 13", len(Default()))
	}
	for _, r := range Default() {
		h := Handoff{From: r.From, To: r.To, SchemaOK: true, Items: map[string]bool{}}
		for _, k := range r.Artifacts {
			a := Artifact{Kind: k, Path: "p/" + k, Named: h64("e"), Found: h64("e")}
			if k == "plan" || k == "plan-review" || k == "review-record" {
				a = Artifact{Kind: k, Comment: &Comment{ID: 7, Seen: 2, File: h64("e"), Copy: h64("e"), Forge: h64("e")}}
			}
			if k == "commits" {
				a.Named, a.Found = h40("e"), h40("e")
			}
			h.Artifacts = append(h.Artifacts, a)
		}
		for _, it := range r.Items {
			if it != "schema" && it != "artifacts" {
				h.Items[it] = true
			}
		}
		if v := Check(Default(), h); !v.Valid {
			t.Errorf("%s -> %s: %+v", r.From, r.To, v)
		}
		for _, it := range r.Items {
			if it == "schema" || it == "artifacts" {
				continue
			}
			g := h
			g.Items = map[string]bool{}
			for k, v := range h.Items {
				g.Items[k] = v && k != it
			}
			if v := Check(Default(), g); v.Valid || !strings.HasPrefix(v.Detail, "invalid "+it) {
				t.Errorf("%s -> %s without %s: %+v", r.From, r.To, it, v)
			}
		}
		if len(r.Artifacts) > 0 && slices.Contains(r.Items, "artifacts") {
			g := h
			g.Artifacts = h.Artifacts[1:]
			if v := Check(Default(), g); v.Valid || v.Detail != "invalid artifacts" {
				t.Errorf("%s -> %s without its first artifact: %+v", r.From, r.To, v)
			}
		}
	}
}

// A replacement table is used by Check, not the default.
func TestCheckUsesTheTableItIsGiven(t *testing.T) {
	table := []Row{{"a", "b", []string{"question"}, []string{"question-row"}}}
	h := Handoff{From: "a", To: "b", SchemaOK: true, Artifacts: []Artifact{{Kind: "question", Named: h64("f"), Found: h64("f")}}, Items: map[string]bool{"question-row": true}}
	if v := Check(table, h); !v.Valid || v.Detail != "valid" {
		t.Fatalf("%+v", v)
	}
	if v := Check(Default(), h); v.Valid || v.Detail != "invalid pair" {
		t.Fatalf("the default table has no pair a, b: %+v", v)
	}
}

// A replacement row that omits "schema" or "artifacts" does not let a bad typed
// result, or a bad artifact, through: both guard every handoff of a session,
// whatever its row lists; only code's own handoffs (a rendered task) have no
// session result.
func TestAReplacementRowCannotLetABadResultThrough(t *testing.T) {
	table := []Row{
		{"a", "b", []string{"task-file"}, []string{"question-row"}},
		{"code (rendered task)", "verifier", []string{"rendering"}, []string{"rendering-equal"}},
	}
	file := func(named, found string) Artifact {
		return Artifact{Kind: "task-file", Path: "f", Named: named, Found: found}
	}
	base := func() Handoff {
		return Handoff{From: "a", To: "b", SchemaOK: true, Artifacts: []Artifact{file(h64("1"), h64("1"))}, Items: map[string]bool{"question-row": true}}
	}
	if v := Check(table, base()); !v.Valid {
		t.Fatalf("the base handoff: %+v", v)
	}
	for name, tc := range map[string]struct {
		change func(*Handoff)
		detail string
	}{
		"a bad schema": {func(h *Handoff) { h.SchemaOK = false }, "invalid schema"},
		"a named artifact with a hash that differs": {func(h *Handoff) { h.Artifacts[0] = file(h64("1"), h64("2")) }, "invalid artifacts"},
		"an extra artifact that is bad": {func(h *Handoff) {
			h.Artifacts = append(h.Artifacts, Artifact{Kind: "patch", Named: h64("1"), Found: h64("9")})
		}, "invalid artifacts"},
		"an extra artifact of an unknown kind": {func(h *Handoff) {
			h.Artifacts = append(h.Artifacts, Artifact{Kind: "blueprint", Named: h64("1"), Found: h64("1")})
		}, "invalid artifacts"},
		"an extra comment with no comment": {func(h *Handoff) { h.Artifacts = append(h.Artifacts, Artifact{Kind: "plan"}) }, "invalid artifacts"},
	} {
		h := base()
		tc.change(&h)
		if v := Check(table, h); v.Valid || v.Detail != tc.detail {
			t.Errorf("%s: %+v, want %q", name, v, tc.detail)
		}
	}
	// A row that does not list "artifacts" does not need its kind present.
	none := base()
	none.Artifacts = nil
	if v := Check(table, none); !v.Valid {
		t.Errorf("a row with no artifacts item and no artifact: %+v", v)
	}
	// Code's own handoff has no session result: SchemaOK is not read.
	rendered := Handoff{From: "code (rendered task)", To: "verifier", CodeWritten: true, Artifacts: []Artifact{{Kind: "rendering", Named: h64("3"), Found: h64("3")}}, Items: map[string]bool{"rendering-equal": true}}
	if v := Check(table, rendered); !v.Valid {
		t.Errorf("a code-written handoff: %+v", v)
	}
	rendered.CodeWritten = false
	if v := Check(table, rendered); v.Valid || v.Detail != "invalid schema" {
		t.Errorf("the same handoff written by a session: %+v", v)
	}
}

// Each of the three comment kinds needs its three hashes equal: each is broken
// alone, and a hash that is empty or not a lowercase SHA-256 fails.
func TestEachCommentKindNeedsItsThreeHashes(t *testing.T) {
	for _, kind := range []string{"plan", "plan-review", "review-record"} {
		table := []Row{{"a", "b", []string{kind}, []string{"schema", "artifacts"}}}
		mk := func(file, copy, forge string) Handoff {
			return Handoff{From: "a", To: "b", SchemaOK: true, Artifacts: []Artifact{{Kind: kind, Comment: &Comment{ID: 3, Seen: 2, File: file, Copy: copy, Forge: forge}}}}
		}
		if v := Check(table, mk(h64("a"), h64("a"), h64("a"))); !v.Valid || v.Detail != "valid comment 3-2" {
			t.Errorf("%s with three equal hashes: %+v", kind, v)
		}
		for name, h := range map[string]Handoff{
			"the file differs":      mk(h64("b"), h64("a"), h64("a")),
			"the copy differs":      mk(h64("a"), h64("b"), h64("a")),
			"the forge differs":     mk(h64("a"), h64("a"), h64("b")),
			"all three are empty":   mk("", "", ""),
			"the forge is empty":    mk(h64("a"), h64("a"), ""),
			"three equal but upper": mk(strings.ToUpper(h64("a")), strings.ToUpper(h64("a")), strings.ToUpper(h64("a"))),
			"three equal but 63":    mk(h64("a")[:63], h64("a")[:63], h64("a")[:63]),
			"three equal but 40":    mk(h40("a"), h40("a"), h40("a")),
		} {
			if v := Check(table, h); v.Valid || v.Detail != "invalid artifacts comment 3-2" {
				t.Errorf("%s, %s: %+v", kind, name, v)
			}
		}
	}
	// A file artifact with an empty or malformed hash fails too.
	table := []Row{{"a", "b", []string{"task-file"}, []string{"schema", "artifacts"}}}
	for name, a := range map[string]Artifact{
		"both empty":        {Kind: "task-file", Path: "f"},
		"found empty":       {Kind: "task-file", Path: "f", Named: h64("1")},
		"named empty":       {Kind: "task-file", Path: "f", Found: h64("1")},
		"equal but upper":   {Kind: "task-file", Path: "f", Named: strings.ToUpper(h64("a")), Found: strings.ToUpper(h64("a"))},
		"equal but not hex": {Kind: "task-file", Path: "f", Named: strings.Repeat("g", 64), Found: strings.Repeat("g", 64)},
	} {
		if v := Check(table, Handoff{From: "a", To: "b", SchemaOK: true, Artifacts: []Artifact{a}}); v.Valid || v.Detail != "invalid artifacts" {
			t.Errorf("%s: %+v", name, v)
		}
	}
}

// The word pair is a detail of the event for a pair with no row, not an item
// of a table: a table that lists it is refused.
func TestThePairIsNotAnItemOfTheTable(t *testing.T) {
	if slices.Contains(Items(), "pair") {
		t.Fatal("pair is in the list of items")
	}
	if err := CheckTable(append(Default(), Row{"a", "b", []string{"plan"}, []string{"pair"}})); err == nil || !strings.Contains(err.Error(), "pair") {
		t.Errorf("CheckTable accepts the item pair: %v", err)
	}
}
