package cli

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"crypto/rsa"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/forge/github"
	"github.com/pharzam/layup/internal/route"
	lrun "github.com/pharzam/layup/internal/run"
)

// LAYUP's own pin (docs/setup/armature.pin): the baseline of a Start
// (docs/spec/run.md, step 3). They are variables, not constants, so a test
// with no network can build a binary whose baseline is a local repository
// (go build -ldflags -X); TestThePinVariablesEqualArmaturePin holds them to
// the file.
var (
	pinSource = "https://github.com/pharzam/armature"
	pinCommit = "a95965534b14b0bf14ad74da0c9a45b5f4aedf88"
)

// The run and its checks, and the context of Ctrl-C; the unit tests replace
// them.
var (
	runStart      = lrun.Start
	runRestart    = lrun.Restart
	runCheckGit   = lrun.CheckGit
	notifyContext = signal.NotifyContext
)

// runNew is layup run --new: each input checked before the first step (exit 2
// with no table), then the steps of Start, the table on standard output and
// the progress on standard error.
func runNew(in call) int {
	cfg, key, err := startInputs(in)
	if err != nil {
		fmt.Fprintf(in.stderr, "layup: %v\n", err)
		return exitUsage
	}
	return runSteps(in, cfg, key, runStart)
}

// runAgain is layup run TARGET.
func runAgain(in call) int {
	cfg, key, err := hostInputs(in, in.args[0])
	if err != nil {
		fmt.Fprintf(in.stderr, "layup: %v\n", err)
		return exitUsage
	}
	return runSteps(in, cfg, key, runRestart)
}

// startInputs checks the flags of Start and builds the Config.
func startInputs(in call) (lrun.Config, *rsa.PrivateKey, error) {
	f := in.flags
	if err := lrun.CheckValues(f["plan"], f["intake-cap"], f["lease-h"], f["watch-t"]); err != nil {
		return lrun.Config{}, nil, err
	}
	cfg, key, err := hostInputs(in, f["new"])
	if err != nil {
		return lrun.Config{}, nil, err
	}
	if cfg.PSB, err = readBriefFile(f["psb"]); err != nil {
		return lrun.Config{}, nil, err
	}
	if v, ok := f["vision"]; ok {
		if cfg.Vision, err = readBriefFile(v); err != nil {
			return lrun.Config{}, nil, err
		}
	}
	cfg.Operator, cfg.IdeaOwner, cfg.Plan, cfg.IntakeCap = f["operator"], f["idea-owner"], f["plan"], f["intake-cap"]
	cfg.LeaseH, _ = strconv.Atoi(f["lease-h"])
	cfg.WatchT, _ = strconv.Atoi(f["watch-t"])
	return cfg, key, nil
}

// hostInputs checks OWNER/NAME, the host directory, its two registers, the
// key file and git, and builds the Config that Start and the restart share.
func hostInputs(in call, target string) (lrun.Config, *rsa.PrivateKey, error) {
	owner, name, ok := strings.Cut(target, "/")
	if !ok || owner == "" || name == "" || strings.ContainsAny(name, "/ ") || strings.Contains(owner, " ") {
		return lrun.Config{}, nil, fmt.Errorf("%q is not OWNER/NAME", target)
	}
	dir, err := filepath.Abs(in.flags["host"])
	if err != nil {
		return lrun.Config{}, nil, err
	}
	data, err := os.ReadFile(filepath.Join(dir, "registers", "forge.tsv"))
	if err != nil {
		return lrun.Config{}, nil, err
	}
	reg, err := forge.ReadForgeRegister(data)
	if err != nil {
		return lrun.Config{}, nil, fmt.Errorf("%s: %w", filepath.Join(dir, "registers", "forge.tsv"), err)
	}
	if data, err = os.ReadFile(filepath.Join(dir, "registers", "harnesses.tsv")); err != nil {
		return lrun.Config{}, nil, err
	}
	harnesses, _, err := route.ReadHarnesses(data)
	if err != nil {
		return lrun.Config{}, nil, fmt.Errorf("%s: %w", filepath.Join(dir, "registers", "harnesses.tsv"), err)
	}
	key, err := forge.CheckKeyFile(reg.KeyFile)
	if err != nil {
		return lrun.Config{}, nil, err
	}
	if err := runCheckGit(); err != nil {
		return lrun.Config{}, nil, err
	}
	id := make([]byte, 8)
	if _, err := rand.Read(id); err != nil {
		return lrun.Config{}, nil, err
	}
	host, err := os.Hostname()
	if err != nil {
		return lrun.Config{}, nil, err
	}
	return lrun.Config{Owner: owner, Name: name, Dir: dir, Register: reg, Harnesses: harnesses, Clock: realClock{},
		RunID: hex.EncodeToString(id), HostName: host, Version: Version, Baseline: pinSource, LayupPin: pinCommit}, key, nil
}

// readBriefFile reads a brief: a readable file of valid UTF-8 (docs/spec/run.md,
// the exit codes); Start copies it byte for byte.
func readBriefFile(path string) ([]byte, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, err
	}
	if !utf8.Valid(data) {
		return nil, fmt.Errorf("%s: not valid UTF-8 (line %d)", path, invalidLine(data))
	}
	return data, nil
}

// runSteps runs the steps with the progress of layup run, ends the beats,
// prints the table and gives its exit code. Ctrl-C ends the context: the run
// gives the rows it reached.
func runSteps(in call, cfg lrun.Config, key *rsa.PrivateKey, steps func(context.Context, lrun.Config) []lrun.Step) int {
	p := newProgress(in.stderr, "run")
	cfg.Step, cfg.Progress = p.step, p.note
	// The adapter makes its own client (forge.md, The App identity); its
	// rate-limit wait prints through the same progress.
	cfg.Forge = github.New(github.Config{API: cfg.Register.API, AppID: cfg.Register.AppID, Key: key,
		Owner: cfg.Owner, Name: cfg.Name, Progress: p.note})
	ctx, stop := notifyContext(context.Background(), os.Interrupt)
	defer stop()
	rows := steps(ctx, cfg)
	p.end()
	if err := lrun.WriteSteps(in.stdout, rows); err != nil {
		fmt.Fprintf(in.stderr, "layup: %v\n", err)
		return exitUsage
	}
	results := make([]string, len(rows))
	for i, r := range rows {
		results[i] = r.Result
	}
	return exitCode(results)
}

// realClock is the run's clock; Sleep ends early with its context.
type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }
func (realClock) Sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
