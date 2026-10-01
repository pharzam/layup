# T-2tc2 — the author's verification (not a review round)

A fresh session of Claude Opus 5.5, the author's model, verified the head of the implementer adversarially (the workflow `rows-1-2-implement`, run `wf_79f86cf4-afc`, 2026-10-01); a fixer then applied each finding that it found real. This is the author's own step: the review round of the gate came later, by another model.

12 findings, 1 of them material.

## The findings

1. **material** — internal/git/git.go:62-69 (environ passes the host's HOME; GIT_SSH_COMMAND is only 'ssh -o BatchMode=yes'); docs/architecture.md:1526-1532 (L-A7); docs/spec/setup.md:87 (S02 sentence); docs/spec/packages.md:170-172; internal/git/git_integration_test.go:281-304
   - Problem: Credentials of the host reach the S02 clone, so the K31 claims are false. L-A7 says 'no credential helper and no ssh agent of the host reach the clone of S02 ... A baseline repository that needs a credential fails S02 at once'. setup.md S02 says 'The clone reads no credential of the host'. packages.md:170-171 says 'a clone that needs a password fails at once'. HTTPS: git's libcurl reads $HOME/.netrc, and environ() passes the host's HOME. ssh: OpenSSH takes the home directory from the password database, not from HOME. So the user's ~/.ssh keys and ~/.ssh/config reach an ssh clone, and an unencrypted key authenticates in BatchMode. TestNoCallAsksForAPassword cannot find this: it uses an isolated HOME with no .netrc, and it checks only the text 'terminal prompts disabled'.
   - Fix proposed: Choose one of two paths. (A) Give git a HOME of the package's own: an empty directory. Do not unset HOME, because libcurl then uses the home from the password database. Give ssh no user file, for example 'ssh -F /dev/null -o BatchMode=yes -o IdentitiesOnly=yes -o IdentityFile=none -o IdentityAgent=none' and a UserKnownHostsFile of the package. Add an integration test that seeds $HOME/.netrc and asserts that the server gets no Authorization header. The answer's list names HOME, so record this change on #79. (B) Rewrite L-A7, the S02 sentence and packages.md:170-172 to name the credential paths that do reach the clone: $HOME/.netrc for HTTPS, and the user's ~/.ssh keys and configuration for ssh. In both paths, update the guardrails §2 check (docs/guardrails.md:171), because 'an environment from a fixed list' that passes HOME does not close .netrc.

2. **note** — internal/git/git.go:141 (CheckoutDetach), :165 (SwitchCreate); docs/spec/packages.md:184-190 (D5: 'It also covers switch (2.23) and init -b (2.28)')
   - Problem: Probable failure on supported versions. git checkout and git switch parse their options with PARSE_OPT_KEEP_DASHDASH. Since git 2.24 (51b4594b40), parse-options keeps --end-of-options in argv for such callers. checkout_main then reads it as a branch or a path ('--detach does not take a path argument'). As far as I know, git changed this only in 2.44 ('parse-options: decouple "--end-of-options" and "--"'). If so, S02 (CheckoutDetach) and S04 (SwitchCreate) fail on git 2.32 to 2.43, for example Ubuntu 22.04 2.34.1, Debian 12 2.39.5 and Ubuntu 24.04 2.43.0. Supported() says yes for these versions. The other uses of --end-of-options (rev-parse 2.30; rev-list, show and diff 2.24; branch, which does not keep '--') are inside the minimum.
   - Fix proposed: Run the integration tests on git 2.32 and on git 2.43 before the minimum is claimed. Or remove --end-of-options from checkout and switch, and pass only a full object ID (resolve it with RevParse first; a full ID never starts with '-'). Or raise MinVersion and give its evidence.

3. **note** — cmd/layup/rules_checker_test.go:213 (starters), :239-254 (scan); docs/spec/packages.md:92-95
   - Problem: The scan does not find three ways to start a program: a composite literal of exec.Cmd followed by Run or Start, syscall.ForkExec and syscall.StartProcess. internal/gate and internal/verify may import os/exec, and syscall has no restriction. So rule 3 ('Only internal/git starts git') can break while TestPackageRules passes. Condition 6 lists only exec.Command, exec.CommandContext, os.StartProcess and syscall.Exec, so the condition is met as written. But the decided text 'a program starts only with exec.Command or exec.CommandContext' states a rule that the scan does not fully check. The inventory item found-boundary-test asked for 'a syscall exec function'.
   - Fix proposed: Add ForkExec and StartProcess to starters["syscall"]. Report each composite literal of exec.Cmd as a start whose program the scan cannot read. Add one unit case for each.

4. **note** — cmd/layup/rules_integration_test.go:61-80 (load); cmd/layup/rules_checker_test.go:160-176; docs/spec/packages.md:30-34 (D6), :90-91
   - Problem: Non-test files behind a build constraint escape rules 4 and 5 and the May import check. D6 binds every non-test Go file. But go list -deps gives Imports and Deps only for the host's GOOS and GOARCH and for no build tag. The source scan reads IgnoredGoFiles, but the import rules do not. So the result depends on the host, and a non-test file behind the integration or e2e tag is never checked.
   - Fix proposed: Parse the imports of the IgnoredGoFiles that load already reads (go/parser with ImportsOnly). Use them for rule 4 and May import, and for rule 5 take the deps of each such import from go list. Or run go list for each tag set and GOOS that the project uses. Or state this limit next to D6.

5. **note** — cmd/layup/rules_test.go:80 (case 'no separator row'); cmd/layup/rules_checker_test.go:77 (separator check), :221-223 (parse error), :66 (break at the next heading)
   - Problem: Some unit cases pass without the branch that they name. The case 'no separator row' has a header and one row only, so the length check (len(table) < 3) refuses it, not the separator check. Also, no unit case reaches the 'does not parse' finding or the break at the next heading.
   - Fix proposed: Give the case 'no separator row' a header and two rows. Add a case with a source that does not parse. Add a case where the section of the heading has no table but a later section has one.

6. **note** — internal/git/git_integration_test.go:223-224, :255-258; internal/git/git.go:157-158; docs/spec/packages.md:147-149
   - Problem: The integration test seeds hostile input (e), GIT_AUTHOR_NAME, but (e) cannot change a commit through the package. Commit appends its own GIT_AUTHOR_* and GIT_COMMITTER_* values last, and os/exec keeps the last value of a duplicate key. So the (e) part of TestAHostileHostChangesNothing passes even when environ() passes the host's GIT_AUTHOR_NAME. The proof of (e) is the unit test TestTheEnvironmentIsAFixedList. packages.md:147-149 says that the integration test seeds each measured input.
   - Fix proposed: State in the test comment and at packages.md:147-149 that the unit test of the environment proves (e). Keep that unit case.

7. **note** — internal/git/git.go:55-56 (config); docs/spec/packages.md:144-149
   - Problem: A file list depends on the host's operating system. On macOS, git init writes core.precomposeunicode=true and core.ignorecase=true into the repository. Add then stores a decomposed (NFD) file name in composed (NFC) form, but a Linux host keeps NFD. For a baseline with an NFD path, the S03 check 'root tree = pin.tree' fails on a macOS host only. Today LAYUP's only non-ASCII path (docs/setup/tests/markers/bad-unlisted/overlay/docs/é.md) is NFC, so the risk is small.
   - Fix proposed: Add -c core.precomposeUnicode=false (APFS keeps a name as written; HFS+ does not, so measure first). Or record this as a known limit next to D3.

8. **note** — docs/tests/traceability.md:35 (T-2tc2/integration/package-rules is 'planned'); docs/prd/PRD-0001-layup.md:241 (NFR-005 Test '—'), :243 (NFR-007 'planned: ...'); docs/tasks/T-2tc2.md (does not exist)
   - Problem: Conditions 4 and 5 and note 7 are not met at this head. #76 has merged (b48764f, PR #98), so the wait of condition 4 is over. The traceability row, the Test cells of NFR-005 and NFR-007, the NFR-001 statement (condition 5) and the rejected alternatives (note 7) must land before the freeze, and none of them is in the diff. The internal/git tests, for example TestAHostileHostChangesNothing, have no traceability row.
   - Fix proposed: Before the freeze, add the edits of conditions 4 and 5 and the task record of note 7, and give the internal/git tests a traceability row.

9. **note** — whole branch (base of the budget)
   - Problem: The plan review sets 1,600 lines over 20 files against base 7cdd346. Against the named base 7cdd346, the branch diff is 5,851 lines over 35 files, because it holds #98. Against b48764f, the diff is 1,534 lines over 13 files. So 66 lines remain for the items of the finding above and for the close-out. The landing figure will most likely be more than the maximum.
   - Fix proposed: Record the change of base (7cdd346 to b48764f) on #79. Then get the Operator's approval of the landing figure, or move the growth to a child issue, before the merge.

10. **note** — docs/spec/packages.md:48-51 (the form of a cell); cmd/layup/rules_checker_test.go:99-111
   - Problem: The decided form and the checker do not agree. The form allows Starts a program to be 'a code span at the start of the cell, whose first word is the program'. The checker also requires that the cell holds no other code span, so it refuses a cell that is valid by the form. The form writes the 'none' mark as a code span (`—`). The checker wants a bare —, and reads a code-span `—` as an allowed import named '—'.
   - Fix proposed: Write the form as the checker reads it: exactly one code span, at the start of the cell; and the bare character — for none.

11. **note** — docs/spec/packages.md:105-107, :119 (Branch, 'Used by': 'a branch at a commit, with no switch'); internal/git/git.go:168-171
   - Problem: D1 decides 'the calls that the steps, the checks and layup gate name'. But no step or check in setup.md and no step in gate.md names git branch. The Used by cell of Branch names no user.
   - Fix proposed: Name the step that needs Branch. Or remove the call, and let a later task add it to the list first.

12. **note** — runs/T-2tc2/red.md:13; docs/guardrails.md:171-172
   - Problem: red.md says 'each of the 19 calls'. The package has 18 calls; there are 19 cases because Add has two. The guardrails check says 'close each file with the setting that its manual names'. But the package closes the system attributes file with GIT_ATTR_NOSYSTEM, and packages.md:166-167 says that this variable is not in the 2.54.0 manual.
   - Fix proposed: Write '18 calls (19 cases)'. Change the check to 'close each file with a documented setting, or with a measured one that the record names'.

## What the fixer did not apply, with its reasons

- Finding 8 is not fixed here. The brief forbids edits of docs/tests/traceability.md, docs/prd/ and docs/tasks/. So the traceability rows, the NFR-005 and NFR-007 Test cells, the NFR-001 statement (condition 5) and the task record with the rejected alternatives (note 7) stay for the author, before the freeze.
- Finding 9 is not fixed here. I cannot post on #79 (no gh, no network). The landing figure is now 1,818 lines (1,809 added, 9 removed) over 15 files against origin/main b48764f, so it is 218 lines over the 1,600-line maximum before the author's close-out edits. Against the review's base 7cdd346, the diff also holds #98.
- This round departs from the approved answer to the plan review, so the departure must go on #79. The answer's fixed list named HOME (passed from the host) and GIT_SSH_COMMAND=ssh -o BatchMode=yes. Now HOME is /dev/null, GIT_SSH_COMMAND is removed, and GIT_ALLOW_PROTOCOL plus the -c values http.emptyAuth=false and core.precomposeUnicode=false are added.
- This round also departs from the brief's verb list and D1 as answered: the verb 'branch' is removed (finding 11), because no step, check or gate step names it.
- For finding 1, I did not use the verifier's exact path A (ssh -F /dev/null with IdentityFile=none and so on). GIT_ALLOW_PROTOCOL stops ssh and the remote helpers from starting at all, so no ssh option is needed. This is measured, and it also closes remote helpers, which the verifier did not name.
- Process slip: one check loop wrote its output to /private/tmp/claude-501/one-check.out, outside the scratch folder. I removed the file at once; nothing else outside the worktree and the scratch folder changed.
- Finding 2 is fixed on a source reading only. No git older than 2.54.0 is on the host, so the failure on git 2.32 to 2.43 is not measured. The fix is not version-dependent.
