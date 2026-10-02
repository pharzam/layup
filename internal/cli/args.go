package cli

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
)

// A command is one row of the command table.
type command struct {
	words []string // the words that name it, for example psb check
	args  []string // the names of its positional arguments, in their order
	flags []flag   // its flags: each one is required and given once
	help  string   // its line in the usage
	run   func(call) int
}

// A flag is a flag of a command: --name VALUE or --name=VALUE.
type flag struct {
	name  string // without the two dashes
	value string // the name of its value in the usage, for example REV
}

// A call is one run of a command: its positional arguments, its flags by
// name, and the two outputs.
type call struct {
	args           []string
	flags          map[string]string
	stdout, stderr io.Writer
}

// parse finds the command of args in table, by the longest match of its
// words, and checks the rest of args by the argument rules of
// docs/spec/README.md (Commands). An error is a usage error: its text is the
// reason.
func parse(table []command, args []string) (command, call, error) {
	c, ok := match(table, args)
	if !ok {
		return command{}, call{}, unknown(table, args)
	}
	in := call{flags: map[string]string{}}
	rest := args[len(c.words):]
	for i := 0; i < len(rest); i++ {
		a := rest[i]
		if !strings.HasPrefix(a, "-") {
			in.args = append(in.args, a)
			continue
		}
		name, value, inline := strings.Cut(strings.TrimPrefix(a, "--"), "=")
		if !strings.HasPrefix(a, "--") || !slices.ContainsFunc(c.flags, func(f flag) bool { return f.name == name }) {
			return command{}, call{}, fmt.Errorf("unknown flag %q", a)
		}
		if _, twice := in.flags[name]; twice {
			return command{}, call{}, fmt.Errorf("flag --%s given twice", name)
		}
		if !inline && i+1 < len(rest) && !strings.HasPrefix(rest[i+1], "-") {
			i++
			value = rest[i]
		}
		if value == "" {
			return command{}, call{}, fmt.Errorf("flag --%s has no value", name)
		}
		in.flags[name] = value
	}
	for _, f := range c.flags {
		if _, ok := in.flags[f.name]; !ok {
			return command{}, call{}, fmt.Errorf("missing flag --%s", f.name)
		}
	}
	if len(in.args) < len(c.args) {
		return command{}, call{}, fmt.Errorf("missing argument %s", c.args[len(in.args)])
	}
	if len(in.args) > len(c.args) {
		return command{}, call{}, fmt.Errorf("extra argument %q", in.args[len(c.args)])
	}
	return c, in, nil
}

// match gives the command whose words are the longest start of args.
func match(table []command, args []string) (command, bool) {
	var best command
	found := false
	for _, c := range table {
		if len(c.words) > len(best.words) && len(c.words) <= len(args) && slices.Equal(c.words, args[:len(c.words)]) {
			best, found = c, true
		}
	}
	return best, found
}

// unknown gives the reason when no command matches: the first word, with the
// second when the first starts a command.
func unknown(table []command, args []string) error {
	if len(args) == 0 {
		return errors.New("no command")
	}
	name := args[0]
	if len(args) > 1 && slices.ContainsFunc(table, func(c command) bool { return c.words[0] == args[0] }) {
		name += " " + args[1]
	}
	return fmt.Errorf("unknown command %q", name)
}

// usage gives the usage text of table: each command with its arguments and
// its help line, then the exit codes.
func usage(table []command) string {
	var b strings.Builder
	b.WriteString("usage: layup <command> [<argument>...]\n\ncommands:\n")
	lines := make([]string, len(table))
	width := 0
	for i, c := range table {
		parts := append(slices.Clone(c.words), c.args...)
		for _, f := range c.flags {
			parts = append(parts, "--"+f.name+" "+f.value)
		}
		lines[i] = strings.Join(parts, " ")
		width = max(width, len(lines[i]))
	}
	for i, c := range table {
		fmt.Fprintf(&b, "  %-*s  %s\n", width, lines[i], c.help)
	}
	b.WriteString("\nexit codes:\n" +
		"  0  every row passed (psb check: no gap)\n" +
		"  1  a row failed or did not run (psb check: a gap)\n" +
		"  2  a usage or input error\n")
	return b.String()
}
