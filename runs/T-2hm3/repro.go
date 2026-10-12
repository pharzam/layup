//go:build ignore

// repro is the made setup of T-2hm3 (#297): a bare target with two loose
// objects in objects/17, a push into it, and the removal of its directory at
// once, N times; it counts the removals that fail. The git on PATH is the one
// that runs.
//
//	go run runs/T-2hm3/repro.go [-n 100] [-fix]
//
// -fix sets receive.autogc and maintenance.auto to false in the target first.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func git(dir, stdin string, args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=/dev/null", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t"}
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "git %v: %v\n%s", args, err, out)
		os.Exit(2)
	}
	return strings.TrimSpace(string(out))
}

func main() {
	n := flag.Int("n", 100, "the number of runs")
	fix := flag.Bool("fix", false, "set receive.autogc and maintenance.auto to false in the target")
	flag.Parse()
	fmt.Println(git("/", "", "--version"))
	fails := 0
	for i := 0; i < *n; i++ {
		d, err := os.MkdirTemp("", "repro")
		if err != nil {
			panic(err)
		}
		bare := filepath.Join(d, "acme", "target.git")
		os.MkdirAll(filepath.Dir(bare), 0o755)
		git("/", "", "init", "-q", "--bare", "-b", "main", bare)
		if *fix {
			git(bare, "", "config", "receive.autogc", "false")
			git(bare, "", "config", "maintenance.auto", "false")
		}
		// The estimate of the geometric repack of git 2.55 counts objects/17
		// only: two loose objects there are more than it allows.
		for k := 0; ; k++ {
			if ents, _ := os.ReadDir(filepath.Join(bare, "objects", "17")); len(ents) >= 2 {
				break
			}
			git(bare, fmt.Sprintf("blob %d of run %d\n", k, i), "hash-object", "-w", "--stdin")
		}
		src := filepath.Join(d, "src")
		git("/", "", "init", "-q", "-b", "main", src)
		os.WriteFile(filepath.Join(src, "a"), []byte("a\n"), 0o644)
		git(src, "", "add", "-A")
		git(src, "", "-c", "maintenance.auto=false", "commit", "-qm", "one")
		git(src, "", "-c", "maintenance.auto=false", "push", "-q", "file://"+bare, "HEAD:refs/heads/main")
		if err := os.RemoveAll(d); err != nil {
			fails++
			fmt.Printf("run %d: RemoveAll: %v\n", i+1, err)
		}
	}
	fmt.Printf("fix=%v: %d of %d removals failed\n", *fix, fails, *n)
	if fails > 0 {
		os.Exit(1)
	}
}
