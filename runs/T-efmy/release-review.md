## Release review — phase 1

| Field | Value |
| ----- | ----- |
| Commit reviewed | `621af09bce4e3ef897a9fa655cabb59229f700ba` |
| Reviewer | Claude Fable 5.1 on Claude Code |
| Lens | REQ-015 and REQ-017 over each code path of the release |
| Briefed on | `.review-in/brief.md`; `.review-in/release-check.txt` (exit 0) and `.review-in/release-check.sh`; `.review-in/contract.md` (the plan D1 to D6, the plan review, the author's answer); `.review-in/inventory.md`; `.review-in/issue-body.md`; `internal/git/git.go` in full; `internal/gate/scratch.go`; `internal/verify/scripts.go`; the eight files of `internal/catalog/go/`; `internal/setup/` (four files); each other non-test Go file of `cmd/` and `internal/` (36 files, 7,625 lines, each read in full); `go.mod`; `PRD-0001` §6 rows REQ-015 to REQ-018, §7 criteria, §9; `F-0003` facts 53 and 55; `docs/spec/packages.md` (NFR-007 rules 1 to 5, the table of phase 1, the calls of `internal/git`); `docs/spec/setup.md` (the boundary of phase 1, S02, S03, S13, S15, `commands.sh`); `docs/architecture.md` §1. |
| Independence claimed | Context: a fresh session, read-only, with no round before this one. Execution: a separate run with its own record; the check was run again from a scratch path. Model: Claude Fable 5.1 differs from the authors' Claude Opus 5.5 (the same provider). |
| REQ-015 | `holds` |
| REQ-017 | `holds` |

### The scenario (uat)

```text
Scenario: The release of phase 1 holds REQ-015 and REQ-017

Given the code of main at 621af09 (the release of phase 1)
When a reviewer whose model differs from the authors' reads each code path, with the output of the release check
Then no code path modifies base LLM weights or trains a model (REQ-015)
And no code path calls a cloud provider or hosting platform to modify it (REQ-017)

Accepted by Claude Fable 5.1 on Claude Code, on 2026-10-03
Covers REQ-015, REQ-017
```

A fresh session of a model other than the authors' stands for the person of `docs/tests/template-uat.md`, by row 19 of the plan, the external input of `gov-release-review`, and Bootstrap mode rule 3.

### The reading

I applied D2. For REQ-015, a code path breaks the row when it holds a machine-learning library, a call to a model or a training service, a program that trains, or a write of model files. For REQ-017, a code path breaks the row when LAYUP itself opens a connection to a cloud provider or a hosting platform, or starts a program that changes its infrastructure, its services or its settings; a command that `layup setup` writes into `commands.sh` for the Operator is not a call by a code path, and `git ls-remote` and `git clone` of S02 read the baseline and change nothing. I agree with this reading: it follows the §7 text of `PRD-0001` ("calls ... to modify it"), it keeps the forge adapter of phase 2 (architecture §1, O-129) out of a `Won't` row, and it matches the boundary of phase 1 in `setup.md` ("makes no forge call"). Each path that a stricter reading would catch is a note below.

### Findings

1. **REQ-015, the dependencies.** `go.mod:1-3` has no `require` line. `go list -deps ./...` names only the standard library and the twelve packages of this module (release check (2); the run below). No machine-learning library, no client of a model or of a training service is in the release. Basis: a code path cannot modify weights or train a model without such a library, a program or a connection. `holds`.

2. **REQ-015, the programs.** The release starts three programs: `git` with fixed `-c` values and a fixed environment (`internal/git/git.go:58-87`); `sh -c` of a gate command (`internal/gate/scratch.go:190`); `sh` of a baseline script (`internal/verify/scripts.go:61`). None is a training program. The commands that LAYUP itself writes into a target are `out=$(gofmt -l .) && test -z "$out" && go vet ./...` and `go test -count=1 ./...` (`internal/catalog/go/kinds.tsv:2` and `:6`), and the job script of a target starts only `sh`, `awk`, `tail`, `iconv`, `tr`, `mktemp`, `rm`, `git` and the kind's tool (`internal/catalog/go/files/.github/gates.sh.tmpl:7-8`, `:170`). `holds`.

3. **REQ-015, the files that the release writes.** Each write is a text, TSV or JSON file: the files of the catalog entry and the manifest (`internal/setup/final.go:103-107`), the facts and the pin (`internal/setup/steps.go:446-450`), the briefs and the prose files (`internal/setup/scaffold.go:88-93`, `:251`), the protection and the ruleset bodies (`internal/setup/final.go:270-275`), the register and the records commit (`internal/setup/final.go:432`, `:512-517`), the scratch tree of a gate (`internal/gate/scratch.go:126-141`), and a fixture patch (`internal/verify/gates.go:142`). No model file. `holds`.

4. **REQ-015, the word model.** The word `model` is a column name of the schemas `telemetry` and `prices` (`internal/records/records.go:28`, `:50`, `:211-212`): a record of the name of a session's model, read and checked by `CheckTelemetry` and `CheckPrice`; no call. `holds`.

5. **REQ-017, the connections.** No package depends on `net`, `net/http` or `crypto/tls` (release check (2); the run below). `internal/records` imports `net/url` (`internal/records/records.go:11`) for `url.Parse` of the source of a price row (`records.go:294`); `net/url` pulls in `net/netip` (the run below). Both packages parse text and open no connection. NFR-007 rule 5 names the three packages above and not these two. `note`, no change to the row.

6. **REQ-017, the git calls that reach a remote.** `LsRemote` (`internal/git/git.go:128-140`) and `Clone` (`git.go:143`) are the two calls that reach a remote, both from S02 (`internal/setup/steps.go:268`, `:273`), to read the baseline. The environment allows the protocols `file`, `git`, `http` and `https` (`git.go:69`), reads no configuration or credential of the host (`git.go:64-69`: `HOME=/dev/null`, `GIT_CONFIG_GLOBAL=/dev/null`, `GIT_TERMINAL_PROMPT=0`), and starts no `ssh`. No `push`, `fetch`, `pull`, `send-pack` or `remote` call exists (release check (4)). A reading that counts each connection to a hosting platform as a call would catch S02 when the baseline is on GitHub; the call reads and changes nothing, so even then it does not modify the platform. `note`.

7. **REQ-017, the commands of `commands.sh`.** `layup setup` writes five commands and never runs one: `git -C target remote add origin 'https://github.com/OWNER/NAME.git'` and `git -C target push origin main` (`internal/setup/steps.go:351-352`, S03); `git -C target push origin layup-setup:main` and `gh api --method POST 'repos/OWNER/NAME/rulesets' --input out/ruleset-default.json` (`internal/setup/final.go:288-289`, S13); `git -C target push origin layup-records` (`internal/setup/final.go:542`, S15). The writer is `writeCommands` (`internal/setup/setup.go:458-471`), which writes the file `out/commands.sh` (`internal/work/work.go:25`) and starts nothing. LAYUP also authors the body of the ruleset (`internal/setup/final.go:198-246`) and the protection file (`final.go:160-192`). Under D2 these are the Operator's commands, run with the Operator's login. A stricter reading, which counts the text of a command that changes the settings of a repository on GitHub, would catch the `gh api` command of S13 and, as a change of refs, the three pushes. `note`.

8. **REQ-017, the workflow that S12 writes.** `internal/catalog/go/files/.github/workflows/gates.yml.tmpl` is a GitHub Actions workflow with `permissions: contents: read` (`:10-11`), which uses `actions/checkout@v4` and `actions/setup-go@v5` (`:18`, `:22`). The platform runs it on the target's pull requests; LAYUP starts nothing. A stricter reading would count a file that the platform runs; it changes no setting of the platform. `note`.

9. **REQ-017, the gate command and the baseline scripts.** The `command` of a kind comes from the target's manifest at the base (`internal/gate/gate.go:105-111`) and runs with the environment of `layup` (`internal/gate/scratch.go:190-192`); the two scripts are files of the target (`internal/verify/scripts.go:20-21`). These are inputs of the target, not code paths of LAYUP. The two commands that LAYUP writes (finding 2) run `gofmt`, `go vet` and `go test` on a module with no requirement (`internal/catalog/go/files/go.mod.tmpl`), so they download nothing. `note`.

10. **REQ-017, the other data that names GitHub.** The commit identity `layup-agent[bot]` with a GitHub noreply address (`internal/setup/setup.go:123`), the ID of the GitHub Actions app in the JSON bodies (`internal/setup/final.go:144`), the module path `github.com/OWNER/NAME` (`final.go:53`), and the test-only package `internal/standin` with `file://` URLs and `github.invalid` (`internal/standin/standin.go:107`, `:252`) are strings in files and commits; none is a call. `note`.

### Runs

| Command | Result |
| ------- | ------ |
| `sh release-check.sh 621af09bce4e3ef897a9fa655cabb59229f700ba` (a copy under `rel96-scratch/`, run from the clone, `HOME` set to `rel96-scratch/home`, `GOCACHE`, `GOTOOLCHAIN=local`, `GOPROXY=off`) | exit 0, `== PASS`; the same lines as `.review-in/release-check.txt` |
| `go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./...` | each import is the standard library or this module; `net/url` only in `internal/records` |
| `go list -deps ./... \| grep -x -E 'net\|net/http\|crypto/tls'` (inside the check) | none |
| `go list -deps -f '{{.ImportPath}}: {{join .Imports " "}}' ./... \| grep ' net/netip'` | the one importer is `net/url` |
| `git grep -n -E 'import "C"\|//go:cgo\|//go:linkname\|unsafe\.' -- 'cmd/*.go' 'internal/*.go' ':!*_test.go'` | none |
| `git grep -n -i -E 'http\|url\|dial\|socket\|curl\|wget\|ssh\|gh \|aws\|gcloud\|azure\|terraform\|kubectl\|docker\|openai\|anthropic\|claude\|model\|weights\|train\|token\|secret\|password\|credential\|api[ _-]?key' -- 'cmd/*.go' 'internal/*.go' ':!*_test.go'` | the hits of findings 4 to 7 and 10; no other |
| `git grep -n -E 'Order: [0-9]\|git -C\|gh api\|push\|remote add' -- 'internal/setup/*.go' ':!*_test.go'` | the five command texts of finding 7, and comments |
| `git status --porcelain \| wc -l` | 0: no file of the repository changed |

*Posted for a reviewer session (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI, a fresh read-only session in a clone at `621af09`, from 04:32:54 UTC, 4 min 42 s, its record at 4 min 22 s, 1,504,470 tokens, USD 6.39) through the `layup-agent` App. The text is the session's file, unchanged. Skipped before it, with no record (Bootstrap mode rule 4): Devin (its usage quota), 04:26:30 to 04:27:44, and OpenCode (Grok 4.7: no output in five minutes), 04:27:49 to 04:32:49. This is the release review of D4 of the plan, the uat of the task, not a review round of its pull request; a copy is `runs/T-efmy/release-review.md`.*
