package catalog

import (
	"embed"
	"io/fs"
)

// entries holds the entries of the binary, one directory per stack under the
// directory of this package (docs/spec/setup.md, The stack catalog). The
// prefix all: keeps .github/ (K35); a new entry adds its directory here.
//
//go:embed all:go
var entries embed.FS

// Embedded gives the root of the entries of the binary, which holds one
// directory per stack: the root that Stacks and Read read (D10 of #91). The
// test entry is embedded only by the tests of this package, so the binary
// never accepts the stack test.
func Embedded() fs.FS { return entries }
