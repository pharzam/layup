package route

import (
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"syscall"

	"github.com/pharzam/layup/internal/tsv"
)

// ModelsSchema is the form of host:registers/models.tsv: the block models.
var ModelsSchema = tsv.Schema{Name: "models", Location: "host:registers/models.tsv", Columns: []tsv.Column{
	{Name: "harness", Type: "id(<word>)", Key: true},
	{Name: "model", Type: "text", Key: true},
	{Name: "context", Type: "int"},
	{Name: "source", Type: "text"},
	{Name: "date", Type: "time"},
	{Name: "use", Type: "enum(yes|no)"},
	{Name: "reason", Type: "text"},
}}

// RoutingRegisterSchema is the form of host:registers/routing.tsv: the block
// routing-register.
var RoutingRegisterSchema = tsv.Schema{Name: "routing-register", Location: "host:registers/routing.tsv", Columns: []tsv.Column{
	{Name: "role", Type: "text", Key: true},
	{Name: "tier", Type: "enum(reasoning|execution)", Key: true},
	{Name: "position", Type: "int", Key: true},
	{Name: "harness", Type: "id(<word>)"},
	{Name: "model", Type: "text"},
}}

// ReadModels reads models.tsv by its schema and the rules of its block: no
// column but reason holds the empty value; reason is the empty value exactly
// when use is yes; source is an http or https URL.
func ReadModels(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, ModelsSchema)
	if err != nil {
		return nil, err
	}
	for i, r := range rows {
		f := fields(ModelsSchema, r)
		bad := func(column, reason string) error { return &tsv.Error{Line: i + 2, Column: column, Reason: reason} }
		for _, c := range []string{"context", "source", "date", "use"} {
			if f[c] == "" {
				return nil, bad(c, "this column never holds the empty value")
			}
		}
		if (f["use"] == "yes") != (f["reason"] == "") {
			return nil, bad("reason", "reason is the empty value exactly when use is yes")
		}
		if u, err := url.Parse(f["source"]); err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			return nil, bad("source", "the source is an http or https URL")
		}
	}
	return rows, nil
}

// ReadRoutingRegister reads routing.tsv of the host by its schema and the
// rules of its block: no column holds the empty value, and the positions of
// each role and tier are 1 to k, each once, in any order of the file.
func ReadRoutingRegister(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, RoutingRegisterSchema)
	if err != nil {
		return nil, err
	}
	count, most, line := map[string]int{}, map[string]int{}, map[string]int{}
	var lists []string
	for i, r := range rows {
		f := fields(RoutingRegisterSchema, r)
		for _, c := range []string{"harness", "model"} {
			if f[c] == "" {
				return nil, &tsv.Error{Line: i + 2, Column: c, Reason: "this column never holds the empty value"}
			}
		}
		k := f["role"] + " " + f["tier"]
		if _, ok := count[k]; !ok {
			lists = append(lists, k)
		}
		p, _ := strconv.Atoi(f["position"])
		if p < 1 {
			return nil, &tsv.Error{Line: i + 2, Column: "position", Reason: "a position is 1 or more"}
		}
		count[k]++
		if p > most[k] {
			most[k], line[k] = p, i+2
		}
	}
	for _, k := range lists {
		if most[k] != count[k] {
			return nil, &tsv.Error{Line: line[k], Column: "position", Reason: fmt.Sprintf("the positions of %s run 1 to %d with no gap", k, count[k])}
		}
	}
	return rows, nil
}

// CheckRegisters checks the three registers across their files: each row of
// models.tsv and routing.tsv names a harness of the register, and each model
// of routing.tsv has a row of its harness in models.tsv. An error names its
// file, line and column.
func CheckRegisters(ids []string, models, routing [][]string) error {
	have := map[string]bool{}
	for i, r := range models {
		f := fields(ModelsSchema, r)
		if !slices.Contains(ids, f["harness"]) {
			return fmt.Errorf("models.tsv: %w", &tsv.Error{Line: i + 2, Column: "harness", Reason: f["harness"] + " is no harness of the register"})
		}
		have[f["harness"]+" "+f["model"]] = true
	}
	for i, r := range routing {
		f := fields(RoutingRegisterSchema, r)
		switch {
		case !slices.Contains(ids, f["harness"]):
			return fmt.Errorf("routing.tsv: %w", &tsv.Error{Line: i + 2, Column: "harness", Reason: f["harness"] + " is no harness of the register"})
		case !have[f["harness"]+" "+f["model"]]:
			return fmt.Errorf("routing.tsv: %w", &tsv.Error{Line: i + 2, Column: "model", Reason: f["model"] + " has no row of " + f["harness"] + " in models.tsv"})
		}
	}
	return nil
}

// CheckCredential checks a credential file that a register row names
// (docs/spec/session.md, The environment and the harness credential): an
// absolute path, a regular file that exists, mode 0600, owned by the user of
// the run. Each error names the file. The owner is read from syscall.Stat_t,
// so the check runs on Unix hosts, as the check of the App's key does.
func CheckCredential(path string) error {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("the credential %s: not an absolute path", path)
	}
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("the credential %s: %w", path, err)
	}
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("the credential %s: its owner cannot be read on this host", path)
	}
	return checkCredentialFile(path, info.Mode(), int(st.Uid), os.Getuid())
}

// checkCredentialFile checks the mode and the owner of a credential file.
func checkCredentialFile(path string, mode fs.FileMode, owner, runner int) error {
	switch {
	case !mode.IsRegular() || mode.Perm() != 0o600:
		return fmt.Errorf("the credential %s, mode %04o: want a regular file of mode 0600", path, mode.Perm())
	case owner != runner:
		return fmt.Errorf("the credential %s: owned by user %d, not by the user of the run (%d)", path, owner, runner)
	}
	return nil
}
