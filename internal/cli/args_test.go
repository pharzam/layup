package cli

import (
	"slices"
	"strings"
	"testing"
)

// testTable has the shapes of the commands of phase 1: one word, two words,
// a command whose word starts a longer command, and flags.
var testTable = []command{
	{words: []string{"version"}},
	{words: []string{"psb", "check"}, args: []string{"FILE"}},
	{words: []string{"setup"}, args: []string{"WORK"}},
	{words: []string{"setup", "verify"}, args: []string{"WORK"}},
	{words: []string{"gate"}, args: []string{"REPO"}, flags: []flag{{"base", "REV"}, {"head", "REV"}}},
}

func TestParseTakesTheLongestMatchOfTheCommandWords(t *testing.T) {
	for _, c := range []struct {
		args  []string
		words string
		pos   []string
	}{
		{[]string{"setup", "verify", "w"}, "setup verify", []string{"w"}},
		{[]string{"setup", "w"}, "setup", []string{"w"}},
		{[]string{"setup", "./verify"}, "setup", []string{"./verify"}},
		{[]string{"psb", "check", "f.md"}, "psb check", []string{"f.md"}},
		{[]string{"version"}, "version", nil},
	} {
		cmd, in, err := parse(testTable, c.args)
		if err != nil || strings.Join(cmd.words, " ") != c.words || !slices.Equal(in.args, c.pos) {
			t.Errorf("%q: command %q, arguments %q, error %v; want %q, %q", c.args, cmd.words, in.args, err, c.words, c.pos)
		}
	}
}

func TestParseTakesAFlagBeforeOrAfterThePositionalArguments(t *testing.T) {
	for _, args := range [][]string{
		{"gate", "r", "--base", "b", "--head", "h"},
		{"gate", "--base", "b", "r", "--head=h"},
		{"gate", "--head=h", "--base=b", "r"},
	} {
		_, in, err := parse(testTable, args)
		if err != nil || !slices.Equal(in.args, []string{"r"}) || in.flags["base"] != "b" || in.flags["head"] != "h" {
			t.Errorf("%q: arguments %q, flags %v, error %v", args, in.args, in.flags, err)
		}
	}
}

func TestParseGivesTheReasonOfEachUsageError(t *testing.T) {
	for _, c := range []struct {
		args   []string
		reason string
	}{
		{nil, "no command"},
		{[]string{"frobnicate"}, `unknown command "frobnicate"`},
		{[]string{"psb"}, `unknown command "psb"`},
		{[]string{"psb", "chek", "f"}, `unknown command "psb chek"`},
		{[]string{"version", "extra"}, `extra argument "extra"`},
		{[]string{"version", "--x"}, `unknown flag "--x"`},
		{[]string{"psb", "check"}, "missing argument FILE"},
		{[]string{"setup", "verify"}, "missing argument WORK"},
		{[]string{"psb", "check", "a", "b"}, `extra argument "b"`},
		{[]string{"psb", "check", "--x", "f"}, `unknown flag "--x"`},
		{[]string{"psb", "check", "-x", "f"}, `unknown flag "-x"`},
		{[]string{"psb", "check", "--", "f"}, `unknown flag "--"`},
		{[]string{"gate", "r", "--head", "h"}, "missing flag --base"},
		{[]string{"gate", "r", "--base", "b", "--base", "c", "--head", "h"}, "flag --base given twice"},
		{[]string{"gate", "r", "--base", "--head", "h"}, "flag --base has no value"},
		{[]string{"gate", "r", "--head", "h", "--base"}, "flag --base has no value"},
		{[]string{"gate", "r", "--base=", "--head", "h"}, "flag --base has no value"},
	} {
		if _, _, err := parse(testTable, c.args); err == nil || err.Error() != c.reason {
			t.Errorf("%q: error %v, want %q", c.args, err, c.reason)
		}
	}
}

func TestTheUsageListsEachCommandWithItsArgumentsAndTheExitCodes(t *testing.T) {
	u := usage(testTable)
	for _, want := range []string{"usage: layup <command>", "\n  psb check FILE ", "\n  setup verify WORK ",
		"\n  gate REPO --base REV --head REV ", "\n  0  every row passed", "\n  2  a usage or input error\n"} {
		if !strings.Contains(u, want) {
			t.Errorf("the usage has no %q:\n%s", want, u)
		}
	}
}

// runTable has the two rows of layup run: Start, whose selecting flag is
// --new, and the restart, with one positional argument (docs/spec/run.md, The
// command); --vision of Start may be left out.
var runTable = append(slices.Clone(testTable),
	command{words: []string{"run"}, selector: "new", optional: []string{"vision"},
		flags: []flag{{"new", "OWNER/NAME"}, {"host", "DIR"}, {"psb", "FILE"}, {"vision", "FILE"}}},
	command{words: []string{"run"}, args: []string{"TARGET"}, flags: []flag{{"host", "DIR"}}})

func TestTheSelectingFlagAndTheOptionalFlag(t *testing.T) {
	for _, c := range []struct {
		args     []string
		selector string
		pos      []string
	}{
		{[]string{"run", "--new", "a/b", "--host", "d", "--psb", "p"}, "new", nil},
		{[]string{"run", "--host=d", "--psb=p", "--vision=v", "--new=a/b"}, "new", nil},
		{[]string{"run", "a/b", "--host", "d"}, "", []string{"a/b"}},
	} {
		cmd, in, err := parse(runTable, c.args)
		if err != nil || cmd.selector != c.selector || !slices.Equal(in.args, c.pos) {
			t.Errorf("%q: the row of selector %q, arguments %q, error %v; want %q, %q", c.args, cmd.selector, in.args, err, c.selector, c.pos)
		}
	}
	for _, c := range []struct {
		args   []string
		reason string
	}{
		{[]string{"run", "a/b", "--host", "d", "--psb", "p"}, `unknown flag "--psb"`},
		{[]string{"run", "--new", "a/b", "--host", "d", "--psb", "p", "x"}, `extra argument "x"`},
		{[]string{"run", "--new", "a/b", "--host", "d"}, "missing flag --psb"},
		{[]string{"run", "--host", "d"}, "missing argument TARGET"},
	} {
		if _, _, err := parse(runTable, c.args); err == nil || err.Error() != c.reason {
			t.Errorf("%q: error %v, want %q", c.args, err, c.reason)
		}
	}
	u := usage(runTable)
	for _, want := range []string{"\n  run --new OWNER/NAME --host DIR --psb FILE [--vision FILE] ", "\n  run TARGET --host DIR "} {
		if !strings.Contains(u, want) {
			t.Errorf("the usage has no %q:\n%s", want, u)
		}
	}
}
