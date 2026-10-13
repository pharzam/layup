package run

// startReadme is the README of layup-records that step 6 writes: the block
// start-readme of docs/spec/run.md (The README of a target that Start makes).
const startReadme = "# The records of this repository\n\n" +
	"This branch, `layup-records`, holds the records that LAYUP keeps for this\n" +
	"repository. It is an orphan branch: it shares no commit with the default\n" +
	"branch.\n\n" +
	"Only `layup run` writes this branch from Start on. Its first commit holds the\n" +
	"records of the Start.\n\n" +
	"Each file is a table of tab-separated values with a header row, or Markdown,\n" +
	"so a person reads it with no tool. A plain `git clone` carries the branch as\n" +
	"`origin/layup-records`.\n\n" +
	"- `start/start.tsv`: each value of the Start, with its source.\n" +
	"- `start/problem-statement.md`, `start/vision.md`: the two briefs, byte for byte.\n" +
	"- `approvers.tsv`: each account whose comment can decide, by its numeric ID and\n" +
	"  role.\n" +
	"- `lease.tsv`: the run that holds this target.\n"

// The titles and the bodies of the two issues of step 7: the blocks
// intake-issue and control-issue of docs/spec/run.md (The Intake and control
// issues).
const (
	intakeTitle  = "LAYUP Intake"
	controlTitle = "LAYUP control"
	intakeBody   = "This is the Intake issue of LAYUP for this repository.\n\n" +
		"LAYUP records each comment here on its records branch, `layup-records`,\n" +
		"before it acts on the comment. Only the accounts of `approvers.tsv` on that\n" +
		"branch decide.\n"
	controlBody = "This is the control issue of LAYUP for this repository.\n\n" +
		"LAYUP records each comment here on its records branch, `layup-records`,\n" +
		"before it acts on the comment. Only the accounts of `approvers.tsv` on that\n" +
		"branch decide.\n"
)

// scaffoldReadme is the README of layup-records that the step verify writes in place of
// that of Start: the block scaffold-readme of docs/spec/scaffold.md (The README of the Scaffold).
const scaffoldReadme = "# The records of this repository\n" +
	"\n" +
	"This branch, `layup-records`, holds the records that LAYUP keeps for this\n" +
	"repository. It is an orphan branch: it shares no commit with the default\n" +
	"branch.\n" +
	"\n" +
	"Only `layup run` writes this branch from Start on. Its first commit holds the\n" +
	"records of the Start.\n" +
	"\n" +
	"Each file is a table of tab-separated values with a header row, or Markdown,\n" +
	"so a person reads it with no tool. A plain `git clone` carries the branch as\n" +
	"`origin/layup-records`.\n" +
	"\n" +
	"- `start/start.tsv`: each value of the Start, with its source.\n" +
	"- `start/problem-statement.md`, `start/vision.md`: the two briefs, byte for byte.\n" +
	"- `approvers.tsv`: each account whose comment can decide, by its numeric ID and\n" +
	"  role.\n" +
	"- `lease.tsv`: the run that holds this target.\n" +
	"- `questions.tsv`, `intake/`: the questions of Intake and their answers.\n" +
	"- `setup/record.tsv`: each value of the setup, with its source, and one done\n" +
	"  row per step.\n" +
	"- `setup/verify.tsv`: the table of `layup setup verify` before the records\n" +
	"  commit of the setup.\n" +
	"- `setup/tree.tsv`: the commits of the setup that `layup setup verify` read,\n" +
	"  which the Operator then pushes.\n" +
	"- `setup/inputs/files/`: the texts of the setup that a role session wrote.\n" +
	"- `setup/rulesets.txt`: the Operator's apply of the two rulesets, with their\n" +
	"  bypass lists.\n" +
	"- `setup/forge-check.tsv`: the rules read back and the probes of the last start.\n" +
	"- `rule-paths.tsv`: the rule-path register: the paths whose change is a change\n" +
	"  of the rules.\n"
