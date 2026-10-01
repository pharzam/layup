# Slice B review record

Reviewer: Claude Fable 5.1 (`claude-fable-5-1`), a fresh `claude -p` session in a disposable clone, read-only tools, no forge token, `--setting-sources project,local`. Not a gate round.

## Round 1

2026-09-29, head `66b7333`; 5 min 24 s; 23,674 output tokens; USD 3.18 at list price. Word for word:

~~~~text
# Slice B review — round 1
Reviewer: claude-fable-5-1. Head: 66b7333. Verdict: `material findings: 9`.

## Material findings
### M1. One step per gate kind is not one required check per kind
- Where: architecture.md:349-350 "one CI job with one step per gate kind, so each kind is its own required check"; ADR-0016 decision 2 says the same.
- Source: GitHub Actions makes one check per job ("checks created by workflow jobs", docs.github.com, required status checks). A step is never a check. The baseline's own S13 and ci.yml:14 tie each required check to a job name. selection-v2.md row 0016 says "one CI job per gate kind".
- Why it matters: the ruleset cannot require "boundary" on its own. Also, a failing `gofmt` step hides whether the boundary step ran (FT1).
- Fix: use one job per gate kind, and name each job in the ruleset.

### M2. The records branch becomes the default branch
- Where: architecture.md:258 "makes the records branch with its first commit" on an empty repository. Then :311 "pushes the verified setup as the first commit of the default branch".
- Source: GitHub makes the first branch pushed to an empty repository its default branch. Only an admin can change it (docs.github.com, "Changing the default branch"). The App has no administration permission (O-92).
- Why it matters: `layup-records` becomes the default branch. The Operator's push creates a second branch that is not the default. A ruleset aimed at the default branch then protects the records branch.
- Fix: step 6 makes the Operator set the default branch. Step 8 reads the default branch back, or the Operator pushes the default branch before `--new` writes records.

### M3. Pushing the adapted setup as the root commit breaks the baseline's pin check
- Where: architecture.md:293 "`setup/steps.tsv` with the changes below" lists changes to S02 and S12 only. Then :311 "The Operator pushes the verified setup as the first commit".
- Source: S03 says "commit the unmodified copy as the root commit", with the evidence "root tree equals the Armature commit tree". `docs/setup/setup-check.sh:59-89` fails `pin` unless the one root commit's tree equals the pinned tree. S13 PUTs `branch-protection.json`, which is a second protection that the rulesets do not replace. R12, I8.
- Why it matters: in the target, `sh docs/setup/setup-check.sh` exits non-zero with "tree: pin names X, root commit has Y". The target fails its own baseline check.
- Fix: push the unmodified copy as the root commit and the setup as later commits. Or list S03, S04 and S13 as changed, and change the pin check. Name both in the list of changes.

### M4. How `layup/gates` treats a pending gate kind is not defined, and the activation batch cannot pass
- Where: ADR-0016 decision 3 "`not-active` is never a pass". Also architecture.md:365-368, which reads the base branch's manifest and gate files.
- Source: FT1, NFR-004/I5, and the required check `layup/gates` with no bypass (:393-395).
- Why it matters: before the first bet, layout, boundary and contract are pending, so every pull request gets `not-active`. If `not-active` fails `layup/gates`, no PR merges. That includes the Intake, spec and activation PRs, and the native step passes them (:354-355). The activation batch is always `not-active`: its tests do not exist on the base, and FT4 forbids running the head's tests. If `not-active` does not fail the check, FT1 is broken.
- Fix: define the map from each kind's result to the one status. For example: pending and no path in scope gives success with that reason; pending and a path in scope gives failure. Give the activation batch its own rule, such as `layup/rules` with the approved hash plus a run of the new tests on the head's known-bad commit.

### M5. A rule batch must be approved before it is pushed, but it is reviewed after
- Where: architecture.md:387-391 "Before `layup run` pushes … A change to a rule path outside an approved batch refuses the result". Also :397-398, where the hash is "recorded with the approval comment". W-04 step 2 puts "the batch lands at that planned point" inside a `model` step.
- Source: O-69, D03, ADR-0017 decisions 2 and 4. Lens 1: the tag must be honest.
- Why it matters: the architect's gate tests exist only on the host. The approver cannot see them in a PR before approving, and the push is refused until approval. One of the two orders must give way, and the text does not say which. W-04 has no `human` approval step and no merge step.
- Fix: add a state "batch proposed". In it, the push goes to a batch PR that `layup/rules` fails until the approval comment records the head's rule-file hash. Add `human` and `code` steps to W-04 for the approval and the merge.

### M6. The Intake answers have no trusted author
- Where: architecture.md:120-121 "each named at Intake in `approvers.tsv`". Also :284-289: the approvers come from the answer comments, and "`layup run` copies each comment … checks that every ID has an answer".
- Source: §3:126 says a human decision needs an author ID in `approvers.tsv`. Decision Points 1 and 2 (`F-0001#10`, `#11`); FT3. W-02 picks a public repository, so anyone can comment.
- Why it matters: when the answers arrive, `approvers.tsv` is empty. Any GitHub user can post "the answers" and name themselves the approver.
- Fix: `layup run --new` takes the Operator's ID (the installing account) and the idea owner's forge ID. It writes them to `approvers.tsv` in the first records commit. A comment from anyone else is recorded as input, never as an answer.

### M7. ADR-0017 says the merge actor is the Operator, which contradicts O-95 and §3
- Where: ADR-0017:57-59 "The activity shows the Operator's account for LAYUP's merges too, so the audit ties each merge … not to an actor name."
- Source: O-95 (installation token, so the App's bot). architecture.md:156-162: each default-branch update must be a `pr_merge` "by the App's bot", or by an approver for a workflows batch.
- Why it matters: the two texts give opposite audit rules. Under the ADR's rule, a merge done by the Operator's own login passes the audit.
- Fix: rewrite decision 5 as §3 states it (the actor ID, the App's bot, or an approver for an O-93 batch).

### M8. `code` steps that read meaning
- Where: architecture.md:279 "each marker that no file answers … so code adds them". Also :300-302 and W-02 step 2 (`code`): "takes each value from … a fact line".
- Source: Lens 1. S10 itself says to "separate those a file or command answers from those that need a decision". NFR-003/I4.
- Why it matters: deciding whether a PSB line answers a marker, or which fact line fills it, is reading meaning. Code either guesses or cannot do it.
- Fix: a session proposes a value's source with a byte-exact quote, and code checks the quote. Otherwise, allow only Intake answers by question ID and catalog entries as sources for `code`.

### M9. The 100 % detection is shown only for gates that fail everything
- Where: architecture.md:308-310 and :409-410, "`layup setup verify` runs each known-bad fixture … (100 % detection, `F-0003#64`)". ADR-0016 decision 4.
- Source: `F-0003#64`, "A set of known-bad commits against the gates". FT1.
- Why it matters: at setup, layout, boundary and contract are pending. They fail any change in scope, so their fixtures "fail" and prove nothing. The architect's tests from the first bet are never run against a known-bad commit. W-04's boundary gate has no detection evidence.
- Fix: record a pending kind's fixture result as `not-active`. Make the activation batch carry one known-bad commit per activated kind, and run them at activation.

## Notes
- N1. W-03 step 1 and W-04 step 3 say "commits in `work/`". §4 names the clone `repo/`.
- N2. architecture.md:321-322, "a status set with a user access token", is stale after O-95.
- N3. The native gate checks are not pinned to GitHub Actions as their source, although S13 of the baseline pins them. Pin them.
- N4. The ruleset requires `layup/verify` from setup, but no slice B step posts it. Mark it `later: slice D`, or say what posts it before then.
- N5. W-01 step 9 checks only that each answer is present and parses. An unclear answer opens no new gap. Say whether a session reads the answers.
- N6. The native CI runs the workflow and the `docs/gates.tsv` of the PR head. So FT4 holds only through `layup/gates` and `layup/rules`, and not once LAYUP is absent. State this.
- N7. The docs for the repository activity API name no App permission. Add a read of it with the App's token to the probe in step 8.
- N8. The GitHub docs say: "If a check and a commit status have the same name, both must pass". The App could post a status named like a native job whose run did not happen. N3 closes this.

## Checklist rows
- S1 answered (M6 and M8 affect it). S2 not answered (M2, M3, M8). S3 not answered (M5, M7, M9). S4 not answered (M1, M4).
- R01 answered. R02 not answered (M3). R03 not answered (M4). R09 answered (the tag is M8). R10 not answered (M4). R12 not answered (M3).
- I3 not answered (M5). I4 (in part) answered. I5 (in part) not answered (M4). I7 answered. I8 not answered (M3).
- K08 not answered (M6). K09 answered. K10 answered. K11 answered. K12 answered (the text changes at plan step 10). K13 answered. K14 answered for the native CI, not for `layup gate` (M4). K15 answered. K16 not answered in full (M3: S03 and S13 are not in the list). K17 known limit L-B1 (acceptable). K18 answered (a reading the Operator must still confirm).
- P02 answered. P12 answered. D03 (the batch) not answered (M5). D07 answered. D18 answered.
- FT1 not answered (M4, M9). FT4 answered (see N6). FT5 answered.
- L-B1 known limit (acceptable). L-B2 known limit (acceptable).

## Existing solutions
- Go boundary tools already exist: golangci-lint `depguard`, `go-arch-lint` and `arch-go` check import rules from a config file. The catalog could use one of them before an architect session writes a custom test.
- Organisation rulesets can require a workflow that runs from a protected source repository, not the PR head. That is FT4 done by GitHub, but it needs an organisation owner.
- OpenSSF Allstar watches branch protection settings all the time. It could narrow L-A4 and the audit.
- The design takes the batch hash from AgentPlane and sdlc-gh, and Shape Up's rabbit holes as gap kinds. It uses them correctly.~~~~

### The author's answer to round 1

All nine material findings and notes N1 to N8 are applied in one commit.

| Finding | Fix |
| ------- | --- |
| M1, N3, N8 | One CI job per gate kind, each pinned to GitHub Actions in the ruleset (§6; ADR-0016 d2). |
| M2, M3 | The Operator pushes the unmodified pinned copy as the root commit first (S03, check `pin`), so the default branch exists before the records branch; the setup commits follow on top (§5 Start and Scaffold; ADR-0016). |
| M4 | A table maps each kind's result to `layup/gates`; a pending kind passes only with no product path; the activation batch runs its approved gate files (§6). |
| M5 | A rule batch is pushed to `batch/<point>` with its own pull request before approval; the approval records its hash; W-04 has the approval and the merge steps (§6; ADR-0017 d4). |
| M6 | `layup run --new` takes the Operator's and the idea owner's logins, and the first records commit holds their IDs; only their comments are answers (§5; W-01 steps 1, 2, 9). |
| M7 | ADR-0017 d5 follows O-95: the merge actor is the App's bot, or an approver for a workflows batch. |
| M8 | Code takes a value only from an Intake answer, a catalog entry or a fact source that the Operator accepted; a session proposes the source with a byte-exact quote (§5). |
| M9 | A pending kind's fixture counts as `not-active`; the activation batch carries one known-bad commit per kind, and each must fail (§6; ADR-0016 d4). |
| N1, N2, N4 to N7 | `repo/`; the stale user-token text removed; `layup/verify` comes from §8; a review session reads the answers; FT4 with LAYUP absent stated; the probe reads the activity API. |

## Round 2

2026-09-29, head `cedff56`; same reviewer and setup. Word for word:

~~~~text
# Slice B review — round 2
Reviewer: claude-fable-5-1. Head: cedff56. Verdict: `material findings: 3`.

**Round-1 closure.** Six findings are closed: M1, M2, M4, M5, M6 and M7. M8 is closed too, because code now only checks a quote, and the Operator accepts the source. M3 is closed for the order of the root commit, but the fix names a check that the baseline does not have (new M1 below). M9 is closed for the pending kinds, but the same hole remains for the kinds that are active at setup (new M2 below). Notes N1 to N8 are applied.

## Material findings
### M1. `docs/setup/` and `setup-check.sh` are LAYUP's files, not the pinned baseline's, so the target's setup evidence has no producer
- **Where:**
  - architecture.md:268-269: "its root tree equals the pinned tree, as step S03 and the baseline's check `pin` need".
  - architecture.md:313: "`setup/steps.tsv` with the changes named here".
  - architecture.md:324-325: "Step S12 no longer adds LAYUP's own `setup-check` job".
  - W-02 step 5: "the pinned baseline's discipline tests, its check `pin` included".
- **Source:** the pinned root commit `d2516fd` ("initialize from Armature kit at a959655") has no `docs/setup/` at all. `setup-check.sh` came from LAYUP's own task T-r7zg (`5d89181`), and `steps.tsv` came from T-9mmm. The script is LAYUP-specific:
  - Check `identity` (setup-check.sh:559-561) says "the repository is LAYUP, not the kit".
  - Check `facts` hard-codes LAYUP's F-0001, F-0003 and F-0004.
  - Check `ci` (:448-449) fails unless `ci.yml` runs `sh docs/setup/setup-check.sh`.

  Rows S04 to S15 of `steps.tsv` take "setup-check X OK" as their evidence. The rows that break are R02, R12, I8, K16 and FT5, and Setup Correctness (`F-0003#63`).
- **Why it matters:** there are two cases, and both fail.
  - If `setup-check.sh` is copied into the target, it is a LAYUP file (FT5). Once S12 has changed, `sh docs/setup/setup-check.sh` exits non-zero on check `ci` and on check `identity`.
  - If it is not copied, nothing produces the evidence that S04 to S15 name, and "the baseline's check `pin`" does not exist.

  W-01 step 6 also takes "the questions of setup steps S01 and S10 from the pinned commit", but those steps are not in the pinned commit.
- **Fix:**
  - Say that `layup setup verify` does each of the checks from outside, for this target: pin, markers, facts and the others.
  - Say that the target gets neither `setup-check.sh` nor its job.
  - Change the evidence column of S04 to S15 to match, and list that change with the S02, S12 and S13 changes.
  - Correct "the baseline's check `pin`", and say that S01 and S10 come from LAYUP's `steps.tsv`.

### M2. The Go gates that are active at setup do not give the result the design states, and setup has no clean run
- **Where:**
  - architecture.md:378: "`gofmt -l` and `go vet ./...` from the setup … test quality is `go test -count=1 ./...`".
  - architecture.md:326-328 ("the known-bad fixtures of §6").
  - W-02 step 5: "the known-bad fixture of each active kind (`gofmt`, `go vet`) makes it fail".
  - ADR-0016 decision 4.
- **Source:** FT1; `F-0003#64`; NFR-004 and I5; ADR-0016: "We reject … a gate that passes when it finds nothing to check without saying so".
- **Why it matters:** there are four problems, and all come from standard Go behaviour.
  - **`gofmt -l` exits 0** even when it lists badly formatted files. So its known-bad fixture does not fail, which is not the result W-02 step 5 states.
  - **No `go.mod` at setup.** The setup writes no `go.mod`, and the pinned copy has none. So `go vet ./...` and `go test ./...` exit 1 on every commit. The `go vet` fixture then "fails" and proves nothing. The test-quality job is a required check with no bypass, so it fails every Intake and spec pull request, and nothing merges.
  - **With a `go.mod` but no package,** both commands only warn that "./..." matched no packages, and exit 0. That is a silent pass on nothing.
  - **No clean run.** Setup verify never runs an active gate on the clean setup tree to see it pass. Only the activation batch has that rule ("the head must pass").
- **Fix:**
  - Make the catalog command explicit, for example `test -z "$(gofmt -l .)"`.
  - Give the active kinds the same rule as a pending kind: pass with the reason "no product path", and state how the job finds "no product path".
  - Say whether the setup writes `go.mod`.
  - Add to `layup setup verify`: each active gate must pass on the clean tree and fail on its fixture.

### M3. The claimed row K11 has no step, and the start value has no source
- **Where:**
  - architecture.md:345-349: "the idea owner posts the baseline numbers … `layup report` compares a measure with its start value only when both exist".
  - W-02 lists K11 under "Checklist rows", but none of its steps does it. §14 has no row for it.
- **Source:** `F-0004#11`: "The first pilot measures the baseline with the current process; then the idea owner sets each start value in one batch, records it in the pilot's PRD". Checklist K11 covers both the baseline and "the idea owner's batch of start values". Lens 2 (a row answered by name only).
- **Why it matters:**
  - Nothing says who sets the start value, when, or where it is recorded.
  - The report can never have "both", so every §7.2 measure prints "not comparable" for the whole pilot.
  - The step "records it in the pilot's PRD" is missing.
- **Fix:** add a `human` step, where the idea owner posts the start values in one batch after the baseline, and a `code` step, where `layup run` records them. The recording can go in the records and in the target's PRD through a pull request. If this belongs to slice G, mark it `later: slice G` in W-02 and in §14.

## Notes
- **N1. Binding the approval to a commit.** ADR-0017 d4 says "The approval comment records the tree hash". §6:427 says that `layup run` records it when the approval arrives. Pick one. Better: the approval request names the head SHA, and code refuses the approval when the head has moved since then.
- **N2. Two meanings of `not-active`.** It means "an active kind that did not run" (§6 table) and also "a pending kind's fixture" (§6:458). ADR-0016 d3 lists only `pass`, `fail` and `not-active`, and has no value for "pending, no product path". Name that result.
- **N3. Where the batch sets `active`.** Is `active` set in the approved batch head, or after the merge ("the manifest then says `active`")? It has to be in the head:
  - A later edit breaks the hash.
  - A pending native job fails a batch whose test files are in the product's scope.
- **N4. Where the known-bad commits live.** The design does not say where the activation's known-bad commits are kept. They must not merge with the batch, and they must stay reproducible. FT6 applies.
- **N5. Pinning a required check to an App.** docs.github.com, "Available rules for rulesets": to select an app as the source, the app "must have recently submitted a check run". Baseline S13 waits for each job to report once, and §5 step 6 drops that wait. Say how the `layup/` checks and the native jobs are pinned before they have reported. For example, the API with `integration_id`, or a first status on the setup head.
- **N6. Approvers outside `approvers.tsv`.** The intent form names "the approver of each planned approval point". §3 lets only the Operator and the idea owner decide. Say whether a third login named there is added to `approvers.tsv`, or refused.
- **N7. The pre-push hook.** If S04 installs the baseline hooks in LAYUP's clone, `.githooks/pre-push` refuses the Operator's push of the setup commits to `main` (§5 Scaffold step 6). Say which clone the printed command runs in.
- **N8. Plan and visibility come late.**
  - The repository's visibility and plan are fixed at Start step 1, but K15 is checked only at Scaffold step 7, after all of the Intake. Check it at Start.
  - S01 still asks for "the name and visibility" of a repository that already exists.
- **N9. FT4 with LAYUP absent.** When LAYUP is absent, the `layup/` checks are removed, so a pull request can edit `docs/gates.tsv` and the native CI runs it. §6:400-402 says this. Record it as a known limit in §15.
- **N10. Check against docs.github.com: both are right.** `GET /repos/{owner}/{repo}/activity` needs Contents: read, and works with an installation token (`pr_merge` and `merge_queue_merge` are valid types). Commit statuses need Commit statuses: write.

## Checklist rows
- S1 answered. S2 not answered (M1). S3 answered (N1). S4 not answered (M2).
- R01 answered. R02 not answered (M1). R03 answered (the Go command is M2). R09 answered. R10 not answered (M2). R12 not answered (M1).
- I3 answered. I4 answered. I5 not answered (M2). I7 answered. I8 not answered (M1).
- K08 answered. K09 answered. K10 answered. K11 not answered (M3). K12 answered. K13 answered. K14 answered. K15 answered (N8). K16 not answered (M1). K17 known limit L-B1 (acceptable). K18 answered (the Operator must still confirm the reading).
- P02 answered. P12 answered. D03 (the batch) answered (N1). D07 answered. D18 answered.
- FT1 not answered (M2). FT4 answered (N9). FT5 not answered (M1: whether a LAYUP check script goes into the target is not decided).
- L-B1 known limit (acceptable). L-B2 known limit (acceptable).

## Existing solutions
- **Copier:** keeps the template commit and an answers file (`.copier-answers.yml`) in the generated repository. That is the pin-plus-values record of S2. Its update flow is prior art for a later baseline update.
- **The `gofmt` idiom:** the common CI form is `test -z "$(gofmt -l .)"`. golangci-lint can run gofmt, vet and depguard in one versioned tool, with non-zero exits.
- **Required workflows:** organisation rulesets can run a protected workflow from outside the pull request's head. That is FT4 with LAYUP absent, but it needs an organisation (N9).
- **OpenSSF Allstar and Scorecard:** they watch branch-protection drift, which could narrow L-A4.~~~~

### The author's answer to round 2

All three material findings and notes N1 to N9 are applied in the next commit. This was the second round; the fixes go to the whole-design review if the Operator applies O-96's rule to this slice too (Q-16 on #72).

| Finding | Fix |
| ------- | --- |
| M1 | `setup/steps.tsv` and `setup-check.sh` are named as LAYUP's files; neither goes into a target; `layup setup verify` does their checks from outside, and the evidence of S04 to S15 changes to it (§5; ADR-0016 d5). |
| M2 | `test -z "$(gofmt -l .)"`; the setup writes `go.mod`; an active kind with no path in scope is `clear` with that reason; setup verify runs each active gate on the clean tree (pass) and on its fixture (fail) (§6; ADR-0016 d3, d4). |
| M3 | The start values are later: slice G, in §5, W-02 and the checklist rows. |
| N1 to N9 | The approval request names the head SHA and code refuses a moved head; `clear` names the result; `active` is set in the batch head; known-bad patches live on the records branch; checks pinned by `integration_id`; a third approver named by the idea owner joins `approvers.tsv`; the push runs from a clone with no hooks; visibility and plan are checked at Start; L-B3 records FT4 with LAYUP absent. |
