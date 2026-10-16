package spec

// The rendering of the bet (docs/spec/loop.md, REQ-012 — The rendered commit
// and the confirmed inventory; task T-3q0b, #350): only the MoSCoW and Phase
// cells of the §6/§7 requirement tables, the bounded `Out of scope by the
// bet:` portion of the `**Non-goals.**` paragraph, and the numbered facts
// record — `<ID>. **<class>:** <text>`. Source text and classes are never
// invented: the caller supplies them.

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// RenderPRD renders the target's PRD from the answer of a bet: for each
// requirement of the answer, the row of §6 or §7 whose first cell names it
// gets the answer's MoSCoW and Phase (— for Won't); and the generated portion
// `Out of scope by the bet: <ids>.` of the `**Non-goals.**` paragraph names
// the facts the answer marked out of scope, numbers first, then the answer
// facts. A requirement with no row or with two is a finding; so is a
// requirement the brief did not name. A run on the rendered bytes gives the
// same bytes.
func RenderPRD(prd []byte, answer Answer) ([]byte, []Finding) {
	var findings []Finding
	lines := strings.Split(string(prd), "\n")
	rows := map[string][]int{} // requirement ID -> its row lines of §6/§7
	section := 0
	for i, l := range lines {
		if strings.HasPrefix(l, "## ") {
			section = 0
			if strings.HasPrefix(l, "## 6.") {
				section = 6
			} else if strings.HasPrefix(l, "## 7.") {
				section = 7
			}
		}
		if section == 0 || !strings.HasPrefix(l, "|") || strings.Count(l, "|") < 6 {
			continue
		}
		id := strings.TrimSpace(strings.Split(l, "|")[1])
		if _, ok := answer.Reqs[id]; ok {
			rows[id] = append(rows[id], i)
		}
	}
	for id, ra := range answer.Reqs {
		switch len(rows[id]) {
		case 0:
			findings = append(findings, Finding{"row", fmt.Sprintf("the requirement %s has no row of §6 or §7", id)})
		case 1:
			cells := strings.Split(lines[rows[id][0]], "|")
			phase := ra.Phase
			if phase == "" {
				phase = "—"
			}
			cells[3], cells[4] = " "+ra.Class+" ", " "+phase+" "
			lines[rows[id][0]] = strings.Join(cells, "|")
		default:
			findings = append(findings, Finding{"row", fmt.Sprintf("the requirement %s has %d rows", id, len(rows[id]))})
		}
	}
	lines = nonGoals(lines, outIDs(answer))
	return []byte(strings.Join(lines, "\n")), findings
}

// outIDs gives the facts the answer marked out of scope, sorted for one
// order: the numbers first, then the answer facts.
func outIDs(a Answer) []string {
	var nums, qs []string
	for id, f := range a.Facts {
		if !f.Out {
			continue
		}
		if strings.HasPrefix(id, "Q-") {
			qs = append(qs, id)
		} else {
			nums = append(nums, id)
		}
	}
	sort.Slice(nums, func(i, j int) bool { x, _ := strconv.Atoi(nums[i]); y, _ := strconv.Atoi(nums[j]); return x < y })
	sort.Strings(qs)
	return append(nums, qs...)
}

// nonGoals writes the generated non-goals portion: the `**Non-goals.**`
// paragraph, to its end at a blank line, drops a portion it holds and gets
// ` Out of scope by the bet: <ids>.` appended, or none.
func nonGoals(lines []string, ids []string) []string {
	start := -1
	for i, l := range lines {
		if strings.Contains(l, "**Non-goals.**") {
			start = i
			break
		}
	}
	if start == -1 {
		return lines
	}
	end := start
	for end+1 < len(lines) && strings.TrimSpace(lines[end+1]) != "" {
		end++
	}
	for i := start; i <= end; i++ {
		if j := strings.Index(lines[i], " Out of scope by the bet:"); j >= 0 {
			lines[i] = lines[i][:j]
			lines = append(lines[:i+1], lines[end+1:]...)
			end = i
			break
		}
	}
	if len(ids) > 0 {
		lines[end] += " Out of scope by the bet: " + strings.Join(ids, ", ") + "."
	}
	return lines
}

// RenderFacts renders the numbered facts record from a version of the
// confirmed inventory (its rows) and the text the caller supplies by fact ID:
// one line per fact `<ID>. **<class>:** <text>`, or `**<class> (out of
// scope):**` for a mark of out. A missing class or text is a finding, never
// invented.
func RenderFacts(version [][]string, textByID map[string]string) ([]byte, []Finding) {
	var findings []Finding
	var b strings.Builder
	for _, r := range version {
		id, class, scope := r[invFact], r[invClass], r[invScope]
		if class == "" {
			findings = append(findings, Finding{"class", fmt.Sprintf("the fact %s has no class", id)})
			continue
		}
		text, ok := textByID[id]
		if !ok || text == "" {
			findings = append(findings, Finding{"text", fmt.Sprintf("the fact %s has no supplied text", id)})
			continue
		}
		if scope == "out" {
			class += " (out of scope)"
		}
		fmt.Fprintf(&b, "%s. **%s:** %s\n", id, class, text)
	}
	return []byte(b.String()), findings
}
