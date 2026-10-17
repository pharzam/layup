package spec

// The answer of the bet (docs/spec/loop.md, REQ-012 — The bet; task T-3q0b,
// #350): one comment — a line for each requirement `<REQ-or-NFR-ID>:
// <Must|Should|Could|Won't>, phase <tag|—>`, a line for each fact to change
// `fact <number-or-Q-NNN>: out of scope` or `fact <id>: <class>`, a line
// `option <decision>: <option>` for each decision with more than one
// survivor, and the verdict `bet` or `no bet`, last.

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var (
	reqLine  = regexp.MustCompile(`^(Must|Should|Could|Won't), phase (\S+)$`)
	factForm = regexp.MustCompile(`^fact (\S+): (.+)$`)
	optForm  = regexp.MustCompile(`^option (\S+): (.+)$`)

	// reqHead heads a requirement line of an answer or of a copy: a
	// requirement ID, a colon and a space. Other lines of a copy — `fact`,
	// `option`, the verdict, prose — are not requirement lines.
	reqHead = regexp.MustCompile(`^(REQ|NFR)-[0-9]{3}: `)
)

// parseReqLine parses one requirement line of a copy: the ID of reqHead, then
// the grammar of reqLine — the one form, shared by the bet's answer and the
// full check's classes rule (task T-4wyg). The phase is `""` for a `—` tag, as
// ReqAnswer writes it.
func parseReqLine(l string) (id string, a ReqAnswer, ok bool) {
	h := reqHead.FindString(l)
	if h == "" {
		return "", a, false
	}
	id = strings.TrimSuffix(h, ": ")
	m := reqLine.FindStringSubmatch(strings.TrimPrefix(l, h))
	if m == nil {
		return id, a, false
	}
	phase := m[2]
	if phase == "—" {
		phase = ""
	}
	return id, ReqAnswer{Class: m[1], Phase: phase}, true
}

// ReqAnswer is the priority the answer gives one requirement: its MoSCoW
// class and its phase ("" for the — that Won't asks).
type ReqAnswer struct{ Class, Phase string }

// FactAnswer is the change the answer gives one fact: out of scope, or a new
// class ("" for none).
type FactAnswer struct {
	Out   bool
	Class string
}

// Answer is the parsed answer of a bet: the requirements' priorities, the
// facts' changes, the chosen options and the verdict. Milestone is not of the
// comment: the caller sets it for the index's bet column.
type Answer struct {
	Reqs      map[string]ReqAnswer
	Facts     map[string]FactAnswer
	Options   map[string]string
	Verdict   string // "bet" or "no-bet"
	Milestone string
}

// ParseAnswer parses the comment that answers a bet's brief: each of
// requiredIDs (from the brief, by the caller) has one line with a phase tag
// of phaseTags — a project-defined tag, not necessarily a number — or — for
// Won't; fact lines name facts of factIDs; option lines choose among
// survivors. A `bet` needs an option line for each decision that has more
// than one survivor; a `no bet` needs none.
func ParseAnswer(body []byte, requiredIDs, phaseTags, factIDs []string, survivors map[string][]string) (Answer, []Finding) {
	a := Answer{Reqs: map[string]ReqAnswer{}, Facts: map[string]FactAnswer{}, Options: map[string]string{}}
	var findings []Finding
	verdict := false
	for _, raw := range strings.Split(strings.TrimRight(string(body), "\n"), "\n") {
		l := strings.TrimRight(raw, " \r")
		if strings.TrimSpace(l) == "" {
			continue
		}
		if verdict {
			findings = append(findings, Finding{"verdict", fmt.Sprintf("a line after the verdict: %q", l)})
			continue
		}
		switch {
		case l == "bet" || l == "no bet":
			verdict = true
			a.Verdict = strings.Replace(l, " ", "-", 1)
		case strings.HasPrefix(l, "fact "):
			m := factForm.FindStringSubmatch(l)
			switch {
			case m == nil:
				findings = append(findings, Finding{"line", fmt.Sprintf("the line %q is no fact line", l)})
			case !slices.Contains(factIDs, m[1]):
				findings = append(findings, Finding{"line", fmt.Sprintf("the fact %s of %q is no fact the records know", m[1], l)})
			case m[2] != "out of scope" && !slices.Contains(classes, m[2]):
				findings = append(findings, Finding{"line", fmt.Sprintf("the line %q is `out of scope` or a class of the four", l)})
			default:
				if _, dup := a.Facts[m[1]]; dup {
					findings = append(findings, Finding{"line", fmt.Sprintf("the fact %s has a line already", m[1])})
				} else if m[2] == "out of scope" {
					a.Facts[m[1]] = FactAnswer{Out: true}
				} else {
					a.Facts[m[1]] = FactAnswer{Class: m[2]}
				}
			}
		case strings.HasPrefix(l, "option "):
			m := optForm.FindStringSubmatch(l)
			switch {
			case m == nil:
				findings = append(findings, Finding{"line", fmt.Sprintf("the line %q is no option line", l)})
			case !slices.Contains(survivors[m[1]], m[2]):
				findings = append(findings, Finding{"line", fmt.Sprintf("the option %s of %q is no survivor of that decision", m[2], m[1])})
			default:
				if _, dup := a.Options[m[1]]; dup {
					findings = append(findings, Finding{"line", fmt.Sprintf("the decision %s has a line already", m[1])})
				} else {
					a.Options[m[1]] = m[2]
				}
			}
		default:
			id, rest, ok := strings.Cut(l, ": ")
			m := reqLine.FindStringSubmatch(rest)
			switch {
			case !ok || m == nil:
				findings = append(findings, Finding{"line", fmt.Sprintf("the line %q is of no form of the answer", l)})
			case !slices.Contains(requiredIDs, id):
				findings = append(findings, Finding{"line", fmt.Sprintf("the line %q is for %s, which the brief did not name", l, id)})
			default:
				if _, dup := a.Reqs[id]; dup {
					findings = append(findings, Finding{"line", fmt.Sprintf("the requirement %s has a line already", id)})
					break
				}
				class, tag := m[1], m[2]
				switch {
				case class == "Won't" && tag != "—":
					findings = append(findings, Finding{"line", fmt.Sprintf("the line %q: a Won't has `phase —`", l)})
				case class != "Won't" && tag == "—":
					findings = append(findings, Finding{"line", fmt.Sprintf("the line %q: `phase —` is for Won't alone", l)})
				case class != "Won't" && !slices.Contains(phaseTags, tag):
					findings = append(findings, Finding{"line", fmt.Sprintf("the line %q: the phase %s is no tag of the target", l, tag)})
				default:
					phase := tag
					if tag == "—" {
						phase = ""
					}
					a.Reqs[id] = ReqAnswer{class, phase}
				}
			}
		}
	}
	for _, id := range requiredIDs {
		if _, ok := a.Reqs[id]; !ok {
			findings = append(findings, Finding{"requirement", fmt.Sprintf("the requirement %s has no line", id)})
		}
	}
	if !verdict {
		findings = append(findings, Finding{"verdict", "the last line is not `bet` or `no bet`"})
	} else if a.Verdict == "bet" {
		for d := range survivors {
			if _, ok := a.Options[d]; !ok && len(survivors[d]) > 1 {
				findings = append(findings, Finding{"line", fmt.Sprintf("the decision %s has more than one survivor and no option line", d)})
			}
		}
	}
	return a, findings
}
