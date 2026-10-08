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
