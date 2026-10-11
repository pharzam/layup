package session

// This file holds the files of a session of a phase (docs/spec/intake.md, The
// sessions of Intake; row 42a of the plan, task T-gmey, #215): its input/,
// and the read of its output file by the rules of a result file.

import (
	"fmt"
	"os"
	"path/filepath"
)

// fileName refuses a name that is not one plain file name: empty, . or .., or
// with a path separator.
func fileName(name string) error {
	if name == "" || name == "." || name == ".." || filepath.Base(name) != name {
		return fmt.Errorf("%q is not one file name", name)
	}
	return nil
}

// MakeInput makes input/ in the directory d of a session of a phase and
// writes each file in it by its name, byte for byte, mode 0600. It runs right
// after Make (whose last write is prompt.md) and before the start row.
func MakeInput(d Dir, files map[string][]byte) error {
	for name := range files {
		if err := fileName(name); err != nil {
			return fmt.Errorf("an input file: %w", err)
		}
	}
	if err := os.Mkdir(d.Input, 0o700); err != nil {
		return err
	}
	for name, data := range files {
		if err := os.WriteFile(filepath.Join(d.Input, name), data, 0o600); err != nil {
			return err
		}
	}
	return nil
}

// ReadOutput reads the output file result/<name> of a session of a phase by
// the rules of a result file, with no block: the caller reads it by its own.
func ReadOutput(d Dir, name string) ([]byte, error) {
	if err := fileName(name); err != nil {
		return nil, fmt.Errorf("the output file: %w", err)
	}
	return readResult(filepath.Join(d.Result, name))
}
