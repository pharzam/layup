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
