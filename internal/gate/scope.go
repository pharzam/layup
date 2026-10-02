package gate

import (
	"fmt"
	"io/fs"
	"strings"
)

// A pattern is one scope pattern of the manifest (docs/spec/gate.md, A scope
// pattern): P/*.E, which matches each file under the directory P at any depth
// whose name ends with .E (./*.E: anywhere), or a path with no *, which
// matches the path itself and each path under it as a directory.
type pattern struct {
	dir, ext string // for P/*.E: P ("." for anywhere) and .E
	prefix   string // for a pattern with no *
}

// parsePattern reads one scope pattern. Any other form with a * is an error.
func parsePattern(s string) (pattern, error) {
	star := strings.Index(s, "*")
	if star < 0 {
		if !validPath(s) {
			return pattern{}, fmt.Errorf("the scope pattern %q is not a path", s)
		}
		return pattern{prefix: s}, nil
	}
	dir, ext, ok := strings.Cut(s, "/*.")
	if !ok || strings.Count(s, "*") != 1 || ext == "" || strings.ContainsAny(ext, "/*") || !(dir == "." || validPath(dir)) {
		return pattern{}, fmt.Errorf("the scope pattern %q is neither P/*.E nor a path with no *", s)
	}
	return pattern{dir: dir, ext: "." + ext}, nil
}

// validPath reports whether s is a path of a tree: no empty, . or .. part, and
// no / at its start or its end.
func validPath(s string) bool { return s != "." && fs.ValidPath(s) }

// match reports whether the path of a file of the tree is in the scope.
func (p pattern) match(path string) bool {
	if p.ext == "" {
		return path == p.prefix || strings.HasPrefix(path, p.prefix+"/")
	}
	if p.dir != "." && !strings.HasPrefix(path, p.dir+"/") {
		return false
	}
	name := path[strings.LastIndex(path, "/")+1:]
	return strings.HasSuffix(name, p.ext) && len(name) > len(p.ext)
}
