# T-3xvc: the compiled mutants

Code head `21a3ae79` (the docs commit; the code is that of `357a757c` and the test that of `a58bc0af`). The tree was copied (`tar` without `.git` and `runs/`) to `/home/layup/.cache/layup-session/t3/scratch`, outside the worktree, by `/home/layup/.cache/layup-session/t3/mutants.sh`; each mutant is one replacement of one text, built (`go build ./...`), vetted (`go vet -tags=e2e ./cmd/layup/ ./internal/run/`), run with `go test -count=1 -tags=e2e -run TestTheSpecificationTaskEndToEnd ./cmd/layup/` (the test builds the binary from the mutated tree, so the binary carries the mutant), restored by copying the original file back, and its hash read again. The complete output is `/home/layup/.cache/layup-session/t3/mutants-out.txt` (run at 10:11 to 10:12 UTC).

**Result: 4 mutants built and vetted; 4 killed; every restore matched its hash: True; the control before (`ok 11.437s`) and after (`ok 21.422s`) was green: True.**

| # | Mutant (the one change) | File | Killed by | First failing assertion | sha256 before | during | after | Restore ok |
| - | ----------------------- | ---- | --------- | ----------------------- | ------------- | ------ | ----- | ---------- |
| E1 | `afterScaffold` drops `r.loopStep` | `internal/run/steps.go` | the row sequence of stage 6 | `loop_e2e_test.go:425: stage 6: exit 0, the rows [... "shape" "bet"]; want ... "shape" "bet" "loop"` | b3ee6748df8dab62 | f78919d00ca4ff16 | b3ee6748df8dab62 | True |
| E2 | `case adv.Merged:` returns `done` after the first merge | `internal/run/loop.go` | the row `loop` of stage 6 (the detail names only the specification task) | `loop_e2e_test.go:427: stage 6: the rows map[... "bet":["done" "the bet 1: the rendered commit 5d23e2cbec63 ...` (the row `loop` is not `the loop merged <S>, <A>`) | a9d9aaadd137843d | b1267380ec56abad | a9d9aaadd137843d | True |
| E3 | the SHA that `MergeTask` passes to the forge is `0000...` (40 zeros) in place of `head`, **after** the head comparison (condition 5: a mutation of `ready.go:213`, not of the head the statuses were read at) | `internal/run/ready.go` | the fake's 409 and the row `loop` | `loop_e2e_test.go:425: stage 6: exit 1 ... loop fail the merge of the pull request 1 at 983dbaa0... was refused and its head did not move: Merge: status 409` | 17e36140730b0505 | 32f274eaf008ea34 | 17e36140730b0505 | True |
| C1 | the fake answers the required job `static` with `conclusion: failure` (fixture control, NFR-004) | `cmd/layup/loop_e2e_test.go` | the row `loop` of stage 6, which is not `the loop merged <S>, <A>` (the binary merges nothing at a head whose required job is a failure) | `loop_e2e_test.go:427: stage 6: the rows map[... "bet":["done" "the bet 1: the rendered commit 0721b7f2d2d7 ...` | a06b51919fa1409d | 4ca55efd57753d8e | a06b51919fa1409d | True |

The patches (unified diff, from `mutants-out.txt`):

```diff
E1  steps.go:  -	return []func(context.Context) Step{r.atBoundary("shape", r.shapeStep), r.atBoundary("bet", r.betStep), r.loopStep}
               +	return []func(context.Context) Step{r.atBoundary("shape", r.shapeStep), r.atBoundary("bet", r.betStep)}
E2  loop.go:   -			merged, stuck = append(merged, t.ID), 0
               -			continue
               +			merged, stuck = append(merged, t.ID), 0
               +			return done(name, "the loop merged "+strings.Join(merged, ", "))
E3  ready.go:  -	merge, err := s.Forge.Merge(ctx, pr.Number, head)
               +	merge, err := s.Forge.Merge(ctx, pr.Number, strings.Repeat("0", 40))
C1  loop_e2e_test.go: "conclusion": "success"  ->  "conclusion": map[bool]string{true: "failure", false: "success"}[c.Context == "static"]
```

The original E3 of the plan (the verified head in place of the head of the statuses) was rewritten as condition 5 asks: passing the verified head to `MergeTask` stops at the head comparison of `markReady` before any merge request, so it cannot reach the fake's 409.
