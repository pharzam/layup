//go:build integration

package run

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestMain gives the package one temporary HOME: t.Setenv panics under
// t.Parallel, and one HOME keeps the host's ~/.gitconfig out, as the per-test
// HOME did. It reads go env GOCACHE before changing HOME so the sub-runs of go
// keep the host's build cache whether or not GOCACHE was exported.
//
// The serial rule (task T-05f4, #439): a top-level integration test of
// internal/run (files internal/run/*_integration_test.go, build tag
// integration) is serial if, directly or through any helper it calls, it calls
// t.Setenv or t.Chdir, or assigns a package-level variable of package run
// (today verifyRun, internal/run/verify.go:29, and betStopAt, internal/run/bet.go). Every other top-level
// integration test gets t.Parallel() as its first statement. Subtests are not
// changed.
//
// The 18 serial top-level integration tests are:
//   - TestAFrozenFileThatIsALinkIsRefusedAndWritesNothingOutside (frozen_integration_test.go): t.Setenv of HOME in checkout
//   - TestAFrozenFileThroughALinkInsideTheCheckoutIsRefused (frozen_integration_test.go): t.Setenv of HOME in checkout
//   - TestAFrozenFileThroughAParentLinkIsRefusedAndWritesNothingOutside (frozen_integration_test.go): t.Setenv of HOME in checkout
//   - TestARecoveryKeepsTheAnswerAlreadyApplied (bet_integration_test.go): writes betStopAt
//   - TestARenderedPushToAnotherBranchTipIsRefused (bet_integration_test.go): writes betStopAt
//   - TestTheBetStepIsAtomicToItsLastAct (bet_integration_test.go): writes betStopAt
//   - TestAPlainCheckoutGetsTheFrozenFilesInOneCommit (frozen_integration_test.go): t.Setenv of HOME in checkout
//   - TestARefusedFrozenListWritesNoFile (frozen_integration_test.go): t.Setenv of HOME in checkout
//   - TestAPushIntoTheTargetOfAWorldStartsNoMaintenance (run_integration_test.go): t.Setenv of GIT_TRACE2_EVENT
//   - TestARecordsCommitRemovesFiles (store_integration_test.go): t.Setenv of TMPDIR
//   - TestThePinnedTreeOfTheMarkers (markers_integration_test.go): t.Setenv of TMPDIR
//   - TestAProseSessionWhoseCommentFailsLeavesNoFileInTheRecords (setup_integration_test.go): writes verifyRun in scaffoldWorld
//   - TestAProseTextWhoseCheckFailsIsRemovedAndANewSessionStarts (setup_integration_test.go): writes verifyRun in scaffoldWorld
//   - TestTheFullInputOfAHeadAfterIntakeAndTheScaffold (specinput_integration_test.go): writes verifyRun (and in scaffoldWorld)
//   - TestTheStepForgeCheckChecksTheForgeAndCommitsItsTableAtEachRestart (forgecheck_integration_test.go): writes verifyRun (and in scaffoldWorld)
//   - TestTheStepScaffoldPushWaitsAndRecordsTheAppliesOfTheOperator (scaffoldpush_integration_test.go): writes verifyRun (and in scaffoldWorld)
//   - TestTheStepSetupRunsS01ToS14 (setup_integration_test.go): writes verifyRun in scaffoldWorld
//   - TestTheStepVerifyCommitsTheRecordsOfTheSetup (verify_integration_test.go): t.Setenv of GOCACHE, GOENV, GOFLAGS, GOTOOLCHAIN and GOPROXY in goEnv, and writes verifyRun (and in scaffoldWorld)
func TestMain(m *testing.M) {
	cache, err := exec.Command("go", "env", "GOCACHE").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "go env GOCACHE: %v\n", err)
		os.Exit(1)
	}
	if err := os.Setenv("GOCACHE", strings.TrimSpace(string(cache))); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	home, err := os.MkdirTemp("", "layup-run-home-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.Setenv("HOME", home); err != nil {
		os.RemoveAll(home)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(home)
	os.Exit(code)
}
