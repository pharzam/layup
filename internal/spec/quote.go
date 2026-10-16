package spec

// The quote coverage of the revised draft (docs/spec/loop.md, REQ-012 — The
// revision of the draft; task T-3q0b, #350): a requirement that covers an
// answer's fact also names the fact of the answer's quote — the one numbered
// fact whose span holds the whole quoted byte interval, found through the
// questions' chain. An error is never read as "no quote"; only genuinely
// absent paired offsets exempt (architecture.md §7 step 5).

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// The row columns the check reads, as the blocks questions, answers and
// draft of records.md give them: a question's id, source, ref and quote
// offsets; an answers row's question; a draft row's id and covers.
const (
	qID, qSource, qRef = 0, 1, 2
	qStart, qEnd       = 6, 7
	aQuestion          = 0
	dID, dCovers       = 0, 3
)

// qIDForm is the form of an answer fact's ID: Q-NNN, three digits or more,
// the value 1 or more.
var qIDForm = regexp.MustCompile(`^Q-(00[1-9]|0[1-9][0-9]|[1-9][0-9]{2,})$`)

// posInt is the form of a fact number: 1 or more, no leading zero.
var posInt = regexp.MustCompile(`^[1-9][0-9]*$`)

// QuoteFactOf gives the number of the one fact whose span holds the whole
// quoted byte interval [start,end) of the problem statement psb. The empty,
// the out-of-bounds, the unpaired, the contained-in-none and the
// contained-in-several are errors, so no error is silently read as no quote.
func QuoteFactOf(start, end int, spans [][]string, psb []byte) (string, error) {
	if start < 0 || start >= end || end > len(psb) {
		return "", fmt.Errorf("the offsets %d and %d are no quote of a problem statement of %d bytes", start, end, len(psb))
	}
	var hits []string
	for _, s := range spans {
		if s[colFact] == "" {
			continue
		}
		ss, _ := strconv.Atoi(s[colStart])
		se, _ := strconv.Atoi(s[colEnd])
		if ss <= start && end <= se {
			hits = append(hits, s[colFact])
		}
	}
	switch len(hits) {
	case 0:
		return "", fmt.Errorf("the quote at %d-%d is in no fact's span", start, end)
	case 1:
		return hits[0], nil
	default:
		return "", fmt.Errorf("the quote at %d-%d is in %d facts' spans", start, end, len(hits))
	}
}

// CheckQuoteFacts checks the persisted covers of the revised draft (rows as
// the block draft gives them): each token is a positive span-fact number or a
// Q-NNN of the answers; a numbered token names a fact of the checked spans; a
// Q-NNN names a row of the answers and of questions.tsv; and a covered answer
// with a quote asks the row's covers to name its quote's fact too. The rows
// of questions, answers and spans are of one records snapshot.
func CheckQuoteFacts(revised, questions, answers, spans [][]string, psb []byte) []Finding {
	byID := map[string][]string{}
	for _, q := range questions {
		byID[q[qID]] = q
	}
	answered := map[string]bool{}
	for _, a := range answers {
		answered[a[aQuestion]] = true
	}
	facts := map[string]bool{}
	for _, s := range spans {
		if s[colFact] != "" {
			facts[s[colFact]] = true
		}
	}
	var findings []Finding
	at := func(r []string, what string) Finding {
		return Finding{"covers", fmt.Sprintf("the row %s: %s", r[dID], what)}
	}
	for _, r := range revised {
		covers := map[string]bool{}
		for _, t := range strings.Fields(r[dCovers]) {
			covers[t] = true
		}
		for _, t := range strings.Fields(r[dCovers]) {
			switch {
			case qIDForm.MatchString(t):
				if !answered[t] {
					findings = append(findings, at(r, fmt.Sprintf("covers names %s, which is no fact of the answers", t)))
					continue
				}
				origin, f := quoteQuestion(t, byID)
				if f != nil {
					findings = append(findings, at(r, f.Detail))
					continue
				}
				s, e := origin[qStart], origin[qEnd]
				if s == "" && e == "" {
					continue // genuinely no quote: it traces to the answer alone
				}
				start, err1 := strconv.Atoi(s)
				end, err2 := strconv.Atoi(e)
				if s == "" || e == "" {
					findings = append(findings, at(r, fmt.Sprintf("the question %s has one quote offset", t)))
				} else if err1 != nil || err2 != nil {
					findings = append(findings, at(r, fmt.Sprintf("the question %s has quote offsets that are no numbers", t)))
				} else if fact, err := QuoteFactOf(start, end, spans, psb); err != nil {
					findings = append(findings, at(r, err.Error()))
				} else if !covers[fact] {
					findings = append(findings, at(r, fmt.Sprintf("covers an answer's fact %s and not the fact %s of that answer's quote", t, fact)))
				}
			case posInt.MatchString(t):
				if !facts[t] {
					findings = append(findings, at(r, fmt.Sprintf("covers names %s, which is no fact of spec/spans.tsv", t)))
				}
			default:
				findings = append(findings, at(r, fmt.Sprintf("covers names %s, which is no fact", t)))
			}
		}
	}
	return findings
}

// quoteQuestion resolves a question's quote: a follow-up asks again the
// question its ref names, whose quote is the one that applies; the chain ends
// at a question that is no follow-up. An ID that names no row, a ref that
// names none, or a cycle is a finding.
func quoteQuestion(id string, byID map[string][]string) ([]string, *Finding) {
	seen := map[string]bool{id: true}
	q, ok := byID[id]
	if !ok {
		return nil, &Finding{"covers", fmt.Sprintf("covers names %s, which is no row of questions.tsv", id)}
	}
	for q[qSource] == "follow-up" {
		next, ok := byID[q[qRef]]
		if !ok {
			return nil, &Finding{"covers", fmt.Sprintf("the follow-up %s asks again %s, which is no row of questions.tsv", q[qID], q[qRef])}
		}
		if seen[q[qRef]] {
			return nil, &Finding{"covers", fmt.Sprintf("the follow-up %s asks again a chain that cycles", q[qID])}
		}
		seen[q[qRef]] = true
		q = next
	}
	return q, nil
}
