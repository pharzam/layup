package tsv

import (
	"fmt"
	"io/fs"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// fieldType is a type of the closed list of docs/spec/README.md (The types).
type fieldType struct {
	expr  string            // as the schema writes it
	valid func(string) bool // reports whether a value is of the type
}

var timeForm = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$`)

// simple holds the check of each type that takes no argument.
var simple = map[string]func(string) bool{
	"text":    oneLine,
	"int":     regexp.MustCompile(`^(0|[1-9][0-9]*)$`).MatchString,
	"decimal": regexp.MustCompile(`^(0|[1-9][0-9]*)\.[0-9]+$`).MatchString,
	"bool":    func(v string) bool { return v == "yes" || v == "no" },
	"sha1":    regexp.MustCompile(`^[0-9a-f]{40}$`).MatchString,
	"sha256":  regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString,
	// The form refuses a fraction and an offset, which time.Parse takes;
	// time.Parse refuses a month 13 or a 30 February.
	"time": func(v string) bool {
		_, err := time.Parse("2006-01-02T15:04:05Z", v)
		return err == nil && timeForm.MatchString(v)
	},
	// One form for each path: no empty, . or .. part, and no / at the start or
	// the end (io/fs.ValidPath, less its root ".").
	"path": func(v string) bool { return v != "." && fs.ValidPath(v) && oneLine(v) },
}

// oneLine reports whether v is text of one line: not empty, valid UTF-8, and
// with no tab, line feed or carriage return.
func oneLine(v string) bool {
	return v != "" && utf8.ValidString(v) && !strings.ContainsAny(v, "\t\n\r")
}

// parseType parses a type of the closed list: text, int, decimal, bool, time,
// sha1, sha256, path, enum(a|b|c), id(<pattern>) or list(<type>).
func parseType(expr string) (fieldType, error) {
	if valid, ok := simple[expr]; ok {
		return fieldType{expr, valid}, nil
	}
	kind, arg, ok := strings.Cut(expr, "(")
	arg, closed := strings.CutSuffix(arg, ")")
	switch {
	case !ok || !closed:
	case kind == "enum":
		words := map[string]bool{}
		for _, w := range strings.Split(arg, "|") {
			if w == "" || strings.ContainsAny(w, " \t()") || words[w] {
				return fieldType{}, fmt.Errorf("type %q: the word %q is empty, repeats, or holds a space or a bracket", expr, w)
			}
			words[w] = true
		}
		return fieldType{expr, func(v string) bool { return words[v] }}, nil
	case kind == "id":
		re, err := idPattern(arg)
		if err != nil {
			return fieldType{}, fmt.Errorf("type %q: %v", expr, err)
		}
		return fieldType{expr, re.MatchString}, nil
	case kind == "list":
		elem, err := parseType(arg)
		if err != nil || strings.HasPrefix(arg, "list(") {
			return fieldType{}, fmt.Errorf("type %q: the type of its values is not a type of the closed list other than a list", expr)
		}
		return fieldType{expr, func(v string) bool {
			for _, e := range strings.Split(v, " ") {
				if !elem.valid(e) {
					return false
				}
			}
			return true
		}}, nil
	}
	return fieldType{}, fmt.Errorf("type %q is not on the closed list", expr)
}

// idPattern makes the pattern of id(<pattern>) a regular expression: a run of
// N is that many digits or more, each x is one lowercase letter or digit,
// <word> is one or more lowercase letters or '-', and each other character
// stands for itself.
func idPattern(p string) (*regexp.Regexp, error) {
	if p == "" {
		return nil, fmt.Errorf("the pattern is empty")
	}
	re := "^"
	for i := 0; i < len(p); i++ {
		switch c := p[i]; {
		case strings.HasPrefix(p[i:], "<word>"):
			re, i = re+"[a-z-]+", i+len("<word>")-1
		case c == 'N':
			k := len(p[i:]) - len(strings.TrimLeft(p[i:], "N"))
			re, i = re+fmt.Sprintf("[0-9]{%d,}", k), i+k-1
		case c == 'x':
			re += "[a-z0-9]"
		case c <= ' ' || c >= utf8.RuneSelf || strings.IndexByte("<>()|", c) >= 0:
			return nil, fmt.Errorf("the pattern %q holds %q, which is not N, x, <word> or a character that stands for itself", p, c)
		default:
			re += regexp.QuoteMeta(string(c))
		}
	}
	return regexp.Compile(re + "$")
}

// check returns an error when value is not a value of the type. The empty
// string is never one: a record holds an empty value as —, which Read and Write
// take before they call check.
func (t fieldType) check(value string) error {
	if t.valid == nil || !t.valid(value) {
		return fmt.Errorf("%q is not a value of %s", value, t.expr)
	}
	return nil
}
