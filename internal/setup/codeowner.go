package setup

// The code-owner rule of S13 in phase 2 (docs/spec/scaffold.md, The code-owner rule; docs/spec/loop.md, The
// gate jobs read the base, part 2; task T-jew3, #356): the CODEOWNERS file of the four path classes, the
// commands that update the live ruleset of a target that M2d set up, and the script that lands the file on
// its base.

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// Migration marks a target that was set up before the code-owner rule: the ID of the ruleset that the
// Operator applied, the `id` of default.json in setup/rulesets.txt of its records.
type Migration struct{ RulesetID int64 }

// The files of the rule: CODEOWNERS in the tree of a target, its copy for the landing in the work area, and
// the script of the landing.
const (
	codeOwnersPath    = ".github/CODEOWNERS"
	codeOwnersOutPath = "out/CODEOWNERS"
	landScriptPath    = "out/land-code-owner.sh"
)

// landApplicable says whether the landing script may run. It is false until the tree check admits the commit
// that the script makes (docs/spec/scaffold.md, The tree check; the question Q1 of #356 is open), so the text
// that S13 writes refuses at its first command.
const landApplicable = false

// ownerLogin reports whether s is the login of a GitHub account: letters, digits and single hyphens, not at
// an end, at most 39 characters (the bound of ownerName).
func ownerLogin(s string) bool { return len(s) <= 39 && ownerForm.MatchString(s) }

// configPathForm is a path that CODEOWNERS holds as one root-anchored pattern with no escape.
var configPathForm = regexp.MustCompile(`^[A-Za-z0-9._-]+(?:/[A-Za-z0-9._-]+)*$`)

// CodeOwners gives the text of .github/CODEOWNERS (loop.md, The gate jobs read the base, part 2): one
// root-anchored line, owned by login, for .github/, docs/gates.tsv, docs/gates/ and each config path of the
// manifest, in the manifest's order, that no earlier line covers by the directory prefix at a slash. A login
// or a path that the file cannot hold is an error, never quoted.
func CodeOwners(login string, configs []string) ([]byte, error) {
	if !ownerLogin(login) {
		return nil, fmt.Errorf("CODEOWNERS: the login %q is not a GitHub login", login)
	}
	covered := []string{".github", "docs/gates.tsv", "docs/gates"}
	lines := []string{"/.github/", "/docs/gates.tsv", "/docs/gates/"}
	for _, c := range configs {
		if !configPathForm.MatchString(c) || strings.Contains("/"+c+"/", "/./") || strings.Contains("/"+c+"/", "/../") {
			return nil, fmt.Errorf("CODEOWNERS: the config path %q cannot be written as a root-anchored pattern", c)
		}
		if !coveredBy(covered, c) {
			covered, lines = append(covered, c), append(lines, "/"+c)
		}
	}
	var b strings.Builder
	b.WriteString("# Written by layup setup (S13): a change to a gate path waits for the Operator's review (docs/spec/loop.md).\n")
	for _, l := range lines {
		b.WriteString(l + " @" + login + "\n")
	}
	return []byte(b.String()), nil
}

// coveredBy reports whether path is one of the paths or lies below one, at a slash.
func coveredBy(paths []string, path string) bool {
	for _, p := range paths {
		if path == p || strings.HasPrefix(path, p+"/") {
			return true
		}
	}
	return false
}

// reduceFilter keeps the fields of a write of a ruleset; `gh api --jq` applies it to a read.
const reduceFilter = `{name,target,enforcement,conditions,rules,bypass_actors}`

// migrationPaths gives the work area, its clone of the target, and the directory of the apply outputs.
func migrationPaths(p Phase2) (area, target, applies string) {
	area = filepath.Join(p.Host, "work", p.Name)
	return area, filepath.Join(area, "target"), filepath.Join(p.Host, "applies", p.Name)
}

// headGuard is the read of the head of the default branch on the forge and the stop unless it is the head of
// layup-setup: the order-2 command of a target that is already pushed, with no push (plan review N6).
func headGuard(p Phase2) string {
	_, target, _ := migrationPaths(p)
	t := quote(target)
	inMessage := p.DefaultBranch
	if !shellSafe(inMessage) {
		inMessage = `"` + quote(inMessage) + `"`
	}
	return strings.Join([]string{
		"remote=$(git -C " + t + " ls-remote origin " + quote("refs/heads/"+p.DefaultBranch) + ")",
		`head=$(printf '%s\n' "$remote" | cut -f1)`,
		"setup=$(git -C " + t + " rev-parse refs/heads/layup-setup)",
		`if [ "$head" != "$setup" ]; then`,
		"\techo \"the head of " + inMessage + " on the forge is '$head', not the head of layup-setup '$setup'\" >&2",
		"\texit 1",
		"fi",
	}, "\n")
}

// updateCommands gives the commands of S13 for a target set up before the rule (scaffold.md, The code-owner
// rule of a target that M2d set up): the guard of the head, then the update of the ruleset of the recorded ID.
// The update does not send the file: it reads the live ruleset, refuses unless it is the one that the records
// name, and writes it back whole with require_code_owner_review true and nothing else changed, so a rule, a
// parameter, a condition or a bypass actor that the file never held is kept whatever GitHub does with an
// omitted field. Each output goes to a temporary file that mv renames when the command succeeds, and a step is
// done only when its output file is missing.
func updateCommands(p Phase2) []Command {
	_, _, applies := migrationPaths(p)
	rs := fmt.Sprintf("/repos/%s/rulesets/%d", p.Name, p.Migrate.RulesetID)
	file := func(kind string) string { return quote(filepath.Join(applies, "default."+kind+".json")) }
	tmp := func(kind string) string { return quote(filepath.Join(applies, "default."+kind+".json.tmp")) }
	pulls := `[.rules[] | select(.type == "pull_request")] | length`
	body := fmt.Sprintf(`if .id != %d or .name != "%s" or .source != "%s" or .enforcement != "active" or (%s) != 1 then `+
		`error("the live ruleset is not the one that setup/rulesets.txt records: id \(.id), name \(.name), source \(.source), enforcement \(.enforcement), \(%s) pull_request rules") `+
		`else %s | .rules |= map(if .type == "pull_request" then .parameters.require_code_owner_review = true else . end) end`,
		p.Migrate.RulesetID, rulesetName, p.Name, pulls, pulls, reduceFilter)
	return []Command{
		{Order: 2, Comment: "the head of the default branch on the forge: a stop unless it is the head of layup-setup, which is already pushed; nothing is pushed", Text: headGuard(p)},
		{Order: 4, Comment: "the directory of the apply outputs", Text: "mkdir -p " + quote(applies)},
		{Order: 4, Comment: "the live ruleset of the default branch, as it is now, kept as it was before the update", Text: "if [ ! -e " + file("before") + " ]; then\n\tgh api " + rs + " > " + tmp("before") + "\n\tmv " + tmp("before") + " " + file("before") + "\nfi"},
		{Order: 4, Comment: "the update of the ruleset of the recorded ID, with your own login, when its output file is missing: the body is the live ruleset with require_code_owner_review true and nothing else changed; the refusal is an error of the filter",
			Text: "if [ ! -e " + file("update") + " ]; then\n\tgh api " + rs + " --jq " + quote(body) + " > " + tmp("body") + "\n\tmv " + tmp("body") + " " + file("body") +
				"\n\tgh api --method PUT " + rs + " --input " + file("body") + " > " + tmp("update") + "\n\tmv " + tmp("update") + " " + file("update") + "\nfi"},
		{Order: 4, Comment: "the read back: the live ruleset, reduced to the fields of a write, must equal the body that was sent",
			Text: "if [ ! -e " + file("after") + " ]; then\n\tgh api " + rs + " --jq " + quote(reduceFilter) + " > " + tmp("after") +
				"\n\tif ! cmp -s " + tmp("after") + " " + file("body") + "; then\n\t\techo 'the ruleset " + fmt.Sprint(p.Migrate.RulesetID) + " read back differs from the body that was sent: compare' " + tmp("after") + " with " + file("body") + " >&2\n\t\texit 1\n\tfi\n\tmv " +
				tmp("after") + " " + file("after") + "\nfi"},
	}
}

// landScript gives out/land-code-owner.sh (scaffold.md, The code-owner rule of a target that M2d set up): it
// adds .github/CODEOWNERS to the head of layup-setup in one commit and pushes it as a fast-forward of the
// default branch, with the ruleset disabled for the push. It starts from the AFTER object, the updated
// protection that the update kept, and puts exactly that back, never the object before the update. When
// applicable is false its first command refuses, naming the open question, and the rest is only to be read.
func landScript(p Phase2, applicable bool) string {
	area, target, applies := migrationPaths(p)
	rs := fmt.Sprintf("/repos/%s/rulesets/%d", p.Name, p.Migrate.RulesetID)
	t := quote(target)
	var b strings.Builder
	b.WriteString("set -eu\n")
	if !applicable {
		b.WriteString("# not applicable until the tree check admits the commit that this script makes (Q1 of #356, open): the next forge-check would stop on the row tree\n")
		b.WriteString("echo 'not applicable: the tree check admits no head but the last commit of setup/tree.tsv (Q1 of #356, open)' >&2\nexit 1\n")
	}
	b.WriteString("after=" + quote(filepath.Join(applies, "default.after.json")) + "\n")
	b.WriteString("off=" + quote(filepath.Join(applies, "default.landing-off.json")) + "\n")
	b.WriteString(`# the AFTER object: the update was done, and it is the protection to put back
if [ ! -e "$after" ]; then echo "$after is missing: run commands.sh first" >&2; exit 1; fi
if ! grep -q '"require_code_owner_review":true' "$after"; then echo "$after does not require the code owner: run commands.sh first" >&2; exit 1; fi
if [ "$(grep -o '"enforcement":"active"' "$after" | wc -l)" != 1 ]; then echo "$after does not hold the enforcement active exactly once" >&2; exit 1; fi
`)
	b.WriteString("# the head of the default branch is the head of layup-setup, and the live ruleset is the AFTER object\n" + headGuard(p) + "\n")
	b.WriteString("live=$(gh api " + rs + " --jq " + quote(reduceFilter) + ")\n")
	b.WriteString(`if [ "$live" != "$(cat "$after")" ]; then echo "the live ruleset is not the object that the update left ($after): nothing was changed" >&2; exit 1; fi` + "\n")
	b.WriteString("# one commit on the head that adds .github/CODEOWNERS, made before the ruleset is touched\n")
	b.WriteString("blob=$(git -C " + t + " hash-object -w " + quote(filepath.Join(area, filepath.FromSlash(codeOwnersOutPath))) + ")\n")
	b.WriteString("dir=$(mktemp -d)\n")
	b.WriteString("GIT_INDEX_FILE=\"$dir/index\" git -C " + t + " read-tree refs/heads/layup-setup\n")
	b.WriteString("GIT_INDEX_FILE=\"$dir/index\" git -C " + t + " update-index --add --cacheinfo \"100644,$blob," + codeOwnersPath + "\"\n")
	b.WriteString("tree=$(GIT_INDEX_FILE=\"$dir/index\" git -C " + t + " write-tree)\n")
	b.WriteString("commit=$(git -C " + t + " commit-tree \"$tree\" -p refs/heads/layup-setup -m 'chore: the code owners of the gate paths')\n")
	b.WriteString("rm -rf \"$dir\"\n")
	b.WriteString("# the ruleset disabled for the push; the trap is armed only after the disable succeeded, and puts the AFTER object back\n")
	b.WriteString("sed 's/\"enforcement\":\"active\"/\"enforcement\":\"disabled\"/' \"$after\" > \"$off.tmp\"\nmv \"$off.tmp\" \"$off\"\n")
	b.WriteString("gh api --method PUT " + rs + " --input \"$off\" > /dev/null\n")
	b.WriteString("restore() {\n\tstatus=$?\n\ttrap - EXIT\n")
	b.WriteString("\tif gh api --method PUT " + rs + " --input \"$after\" > /dev/null && [ \"$(gh api " + rs + " --jq " + quote(reduceFilter) + ")\" = \"$(cat \"$after\")\" ]; then\n\t\texit \"$status\"\n\tfi\n")
	b.WriteString("\techo \"the ruleset " + fmt.Sprint(p.Migrate.RulesetID) + " of " + p.Name + " could not be put back as the object of $after: enforcement may be off; put it back now with: gh api --method PUT " + rs + " --input $after\" >&2\n\texit 1\n}\n")
	b.WriteString("trap restore EXIT\n")
	b.WriteString("git -C " + t + " push origin \"$commit:refs/heads/\"" + shellWord(p.DefaultBranch) + "\n")
	return b.String()
}
