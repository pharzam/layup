package run

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
)

// gitVersion gives the version of the host's git; the unit tests replace it.
var gitVersion = git.Version

// CheckGit checks, before the first step, that git is found and is
// git.MinVersion or newer (docs/spec/run.md, Input states). internal/cli maps
// its error to exit 2, as it does for layup gate.
func CheckGit() error {
	v, err := gitVersion()
	if err != nil {
		return fmt.Errorf("git %s or newer is needed: %v", git.MinVersion, err)
	}
	if !git.Supported(v) {
		return fmt.Errorf("git %s or newer is needed: version %s", git.MinVersion, v)
	}
	return nil
}

// The plans of a forge (docs/spec/run.md, The command).
var plans = []string{"free", "pro", "team", "enterprise"}

// CheckValues checks the values of the flags --plan, --intake-cap, --lease-h
// and --watch-t before the first step, with the forms of the block start
// (docs/spec/records.md), and gives the first error.
func CheckValues(plan, intakeCap, leaseH, watchT string) error {
	if !slices.Contains(plans, plan) {
		return fmt.Errorf("--plan %q: want free, pro, team or enterprise", plan)
	}
	if err := records.CheckValue("intake.cap", intakeCap); err != nil {
		return fmt.Errorf("--intake-cap: %v (MONEY,HOURS, two decimals)", err)
	}
	for _, f := range []struct{ flag, name, value string }{{"--lease-h", "lease.H", leaseH}, {"--watch-t", "watch.T", watchT}} {
		if n, err := strconv.Atoi(f.value); records.CheckValue(f.name, f.value) != nil || err != nil || n < 1 {
			return fmt.Errorf("%s %q: want whole minutes, 1 or more", f.flag, f.value)
		}
	}
	return nil
}
