package cli

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	lrun "github.com/pharzam/layup/internal/run"
)

// runHost makes a host directory with a forge register, a harness register,
// a key file of mode 0600 made at run time, and a brief; it gives the
// directory and the arguments of a good layup run --new.
func runHost(t *testing.T) (string, []string) {
	t.Helper()
	dir := t.TempDir()
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	keyFile := filepath.Join(dir, "app.pem")
	write := func(path string, data []byte, mode os.FileMode) {
		os.MkdirAll(filepath.Dir(path), 0o755)
		if err := os.WriteFile(path, data, mode); err != nil {
			t.Fatal(err)
		}
	}
	write(keyFile, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}), 0o600)
	write(filepath.Join(dir, "registers", "forge.tsv"), []byte("forge\tapp_id\tapp_slug\tkey_file\twatch_slug\tapi\tweb\n"+
		"github\t42\tlayup-agent\t"+keyFile+"\tlayup-watch\thttps://api.github.com\thttps://github.com\n"), 0o644)
	write(filepath.Join(dir, "registers", "harnesses.tsv"), []byte("harness\tcap\twall\nclaude\t10.0\t60\n"), 0o644)
	write(filepath.Join(dir, "psb.md"), []byte("# The problem\n"), 0o644)
	write(filepath.Join(dir, "vision.md"), []byte("# The vision\n"), 0o644)
	return dir, []string{"run", "--new", "acme/target", "--host", dir, "--psb", filepath.Join(dir, "psb.md"),
		"--vision", filepath.Join(dir, "vision.md"), "--operator", "pharzam", "--idea-owner", "carol", "--plan", "free",
		"--intake-cap", "50.0,8.0", "--lease-h", "5", "--watch-t", "10"}
}

// stubRun replaces lrun.Start, lrun.Restart and run.CheckGit until the test
// ends; it gives the Config that the command built.
func stubRun(t *testing.T, rows []lrun.Step, gitErr error) *lrun.Config {
	t.Helper()
	var got lrun.Config
	savedStart, savedRestart, savedGit, savedNotify := runStart, runRestart, runCheckGit, notifyContext
	stub := func(ctx context.Context, cfg lrun.Config) []lrun.Step {
		got = cfg
		cfg.Step(1, len(rows), "forge")
		if ctx.Err() != nil {
			return []lrun.Step{{Name: "forge", Result: "fail", Detail: ctx.Err().Error()}}
		}
		return rows
	}
	runStart, runRestart = stub, stub
	runCheckGit = func() error { return gitErr }
	t.Cleanup(func() {
		runStart, runRestart, runCheckGit, notifyContext = savedStart, savedRestart, savedGit, savedNotify
	})
	return &got
}

func runLayup(args ...string) (int, string, string) {
	var stdout, stderr bytes.Buffer
	code := Run(args, &stdout, &stderr)
	return code, stdout.String(), stderr.String()
}

func TestRunNewBuildsTheConfigAndPrintsTheTable(t *testing.T) {
	dir, args := runHost(t)
	got := stubRun(t, []lrun.Step{{Name: "forge", Result: "done", Detail: "the token"}, {Name: "plan", Result: "done", Detail: "free"}}, nil)
	code, out, errOut := runLayup(args...)
	if code != 0 || out != "step\tresult\tdetail\nforge\tdone\tthe token\nplan\tdone\tfree\n" {
		t.Fatalf("exit %d, stdout %q, stderr %q", code, out, errOut)
	}
	if !strings.Contains(errOut, "layup run: [1/2] forge\n") {
		t.Errorf("stderr %q; want the progress line of step 1", errOut)
	}
	cfg := *got
	if cfg.Owner != "acme" || cfg.Name != "target" || cfg.Dir != dir || string(cfg.PSB) != "# The problem\n" || string(cfg.Vision) != "# The vision\n" ||
		cfg.Operator != "pharzam" || cfg.IdeaOwner != "carol" || cfg.Plan != "free" || cfg.IntakeCap != "50.0,8.0" || cfg.LeaseH != 5 || cfg.WatchT != 10 {
		t.Errorf("the Config of the flags: %+v", cfg)
	}
	if cfg.Register.AppSlug != "layup-agent" || len(cfg.Harnesses) != 1 || cfg.Harnesses[0][0] != "claude" || cfg.Forge == nil || cfg.Clock == nil {
		t.Errorf("the registers, the adapter or the clock of the Config: %+v", cfg)
	}
	if !regexp.MustCompile(`^[0-9a-f]{16}$`).MatchString(cfg.RunID) || cfg.HostName == "" || cfg.Version != Version {
		t.Errorf("the run ID %q, the host name %q, the version %q", cfg.RunID, cfg.HostName, cfg.Version)
	}
	if cfg.Baseline != pinSource || cfg.LayupPin != pinCommit {
		t.Errorf("the pin of the Config: %q, %q", cfg.Baseline, cfg.LayupPin)
	}
	// With no --vision, the Config has no vision brief.
	var noVision []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--vision" {
			i++
			continue
		}
		noVision = append(noVision, args[i])
	}
	if code, _, errOut := runLayup(noVision...); code != 0 || got.Vision != nil {
		t.Errorf("no --vision: exit %d, vision %q, stderr %q", code, got.Vision, errOut)
	}
	// A brief with a marker is read: run.md refuses only a brief that is not a
	// readable UTF-8 file (condition 5 of the plan review).
	marked := []byte("# The problem\n\n\u2039an open gap\u203a\n")
	os.WriteFile(filepath.Join(dir, "psb.md"), marked, 0o644)
	if code, _, errOut := runLayup(args...); code != 0 || !bytes.Equal(got.PSB, marked) {
		t.Errorf("a brief with a marker: exit %d, stderr %q", code, errOut)
	}
	// A fail row gives exit 1, with the table.
	stubRun(t, []lrun.Step{{Name: "forge", Result: "fail", Detail: "no"}}, nil)
	if code, out, _ := runLayup(args...); code != 1 || !strings.HasSuffix(out, "forge\tfail\tno\n") {
		t.Errorf("a fail row: exit %d, stdout %q", code, out)
	}
}

func TestARelativeHostIsMadeAbsolute(t *testing.T) {
	dir, args := runHost(t)
	got := stubRun(t, []lrun.Step{{Name: "forge", Result: "done"}}, nil)
	wd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(wd) })
	os.Chdir(filepath.Dir(dir))
	args[4] = filepath.Base(dir) // --host
	if code, _, errOut := runLayup(args...); code != 0 || !filepath.IsAbs(got.Dir) {
		t.Errorf("a relative --host: exit %d, Dir %q, stderr %q", code, got.Dir, errOut)
	}
}

func TestRunNewChecksEachInputBeforeTheFirstStep(t *testing.T) {
	for _, c := range []struct {
		name   string
		change func(dir string, args []string) []string
		gitErr error
		want   string
	}{
		{"OWNER/NAME with no slash", func(_ string, a []string) []string { a[2] = "target"; return a }, nil, "OWNER/NAME"},
		{"OWNER/NAME with two slashes", func(_ string, a []string) []string { a[2] = "a/b/c"; return a }, nil, "OWNER/NAME"},
		{"a brief that is missing", func(d string, a []string) []string { a[6] = filepath.Join(d, "none.md"); return a }, nil, "none.md"},
		{"a brief that is not UTF-8", func(d string, a []string) []string {
			os.WriteFile(filepath.Join(d, "psb.md"), []byte("caf\xe9\n"), 0o644)
			return a
		}, nil, "UTF-8"},
		{"a plan of no list", func(_ string, a []string) []string { a[14] = "gold"; return a }, nil, "--plan"},
		{"an intake cap of one value", func(_ string, a []string) []string { a[16] = "50.0"; return a }, nil, "--intake-cap"},
		{"lease.H of 0", func(_ string, a []string) []string { a[18] = "0"; return a }, nil, "--lease-h"},
		{"a forge register with no row", func(d string, a []string) []string {
			os.WriteFile(filepath.Join(d, "registers", "forge.tsv"), []byte("forge\tapp_id\tapp_slug\tkey_file\twatch_slug\tapi\tweb\n"), 0o644)
			return a
		}, nil, "forge.tsv"},
		{"a harness register that is missing", func(d string, a []string) []string {
			os.Remove(filepath.Join(d, "registers", "harnesses.tsv"))
			return a
		}, nil, "harnesses.tsv"},
		{"a harness row with wall empty", func(d string, a []string) []string {
			os.WriteFile(filepath.Join(d, "registers", "harnesses.tsv"), []byte("harness\tcap\twall\nclaude\t10.0\t—\n"), 0o644)
			return a
		}, nil, "harnesses.tsv"},
		{"a key file of mode 0644", func(d string, a []string) []string { os.Chmod(filepath.Join(d, "app.pem"), 0o644); return a }, nil, "mode 0644"},
		{"git older than 2.32", func(_ string, a []string) []string { return a }, errors.New("git 2.32.0 or newer is needed: version 2.31.0"), "2.32.0"},
	} {
		dir, args := runHost(t)
		called := false
		stubRun(t, nil, c.gitErr)
		runStart = func(context.Context, lrun.Config) []lrun.Step { called = true; return nil }
		code, out, errOut := runLayup(c.change(dir, args)...)
		if code != 2 || out != "" || called || !strings.HasPrefix(errOut, "layup: ") || !strings.Contains(errOut, c.want) || strings.Contains(errOut, "usage:") {
			t.Errorf("%s: exit %d, stdout %q, a step ran %v, stderr %q; want exit 2 and an input error that names %q", c.name, code, out, called, errOut, c.want)
		}
	}
}

func TestRestartChecksTheHostAndRunsTheRestart(t *testing.T) {
	dir, _ := runHost(t)
	got := stubRun(t, []lrun.Step{{Name: "forge", Result: "done"}, {Name: "clone", Result: "done"}}, nil)
	code, out, errOut := runLayup("run", "acme/target", "--host", dir)
	if code != 0 || !strings.HasPrefix(out, "step\tresult\tdetail\nforge\tdone\t") || got.Owner != "acme" || got.Name != "target" || got.Dir != dir {
		t.Errorf("the restart: exit %d, stdout %q, stderr %q, Config %+v", code, out, errOut, *got)
	}
	if got.PSB != nil || got.Plan != "" || got.LeaseH != 0 {
		t.Errorf("the restart's Config holds a value of Start: %+v", *got)
	}
	if code, _, errOut := runLayup("run", "acme", "--host", dir); code != 2 || !strings.Contains(errOut, "OWNER/NAME") {
		t.Errorf("a restart of a TARGET with no slash: exit %d, stderr %q", code, errOut)
	}
}

// Ctrl-C ends the context of the run; the run gives the rows it reached, the
// table is printed, and the exit is 1.
func TestCtrlCPrintsTheTableOfTheRowsReached(t *testing.T) {
	_, args := runHost(t)
	stubRun(t, []lrun.Step{{Name: "forge", Result: "done"}}, nil)
	notifyContext = func(ctx context.Context, _ ...os.Signal) (context.Context, context.CancelFunc) {
		ctx, cancel := context.WithCancel(ctx)
		cancel()
		return ctx, cancel
	}
	if code, out, _ := runLayup(args...); code != 1 || out != "step\tresult\tdetail\nforge\tfail\tcontext canceled\n" {
		t.Errorf("Ctrl-C: exit %d, stdout %q", code, out)
	}
}

// The pin of LAYUP (decision 2 of #131) equals docs/setup/armature.pin.
func TestThePinVariablesEqualArmaturePin(t *testing.T) {
	data, err := os.ReadFile("../../docs/setup/armature.pin")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"\nsource=" + pinSource + "\n", "\ncommit=" + pinCommit + "\n"} {
		if !strings.Contains(string(data), want) {
			t.Errorf("armature.pin has no line %q", strings.TrimSpace(want))
		}
	}
}
