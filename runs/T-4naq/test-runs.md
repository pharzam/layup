# T-4naq: the test runs

Task `T-4naq` (#428), row 66ab of M2e (working name K7). Times are `date -u` readings taken after each step; an exit status is the direct one of the command, never of a pipe.

## S0. The reconciliation

Main at the start: `ebbb581` (the merge of #455, K4), which holds K3 (`4259e05`) and K4. The branch `T-4naq` stood at the claim `d9852d2`; `git merge origin/main` was a fast-forward, so no merge commit was made.
Two things in the brief differ from the issue: it names the comments of #425 (K4's issue), whose IDs (6262367609, 6262491026, 6262589337) are the plan, the plan review and the Operator's decisions of #428 (this task), which were read.

| Cited by the plan (6262367609) | At `ebbb581` | Divergence |
| --- | --- | --- |
| `ReadBets` / `BetsSchema`, verdict `bet`, kind `first` | `internal/run/bets.go:14`, `:45` | none; `bet.go` reads the bets with it |
| the Intake class of a fact: a span fact from `spec/spans.tsv`, a `Q-NNN` by `factOfQuote`, else `context` | `answerFactClass` (`betbrief.go:196`) with `factRows` (`:219`), `factOfQuote` (`revision.go:34`); `(*state).factsOf` (`bet.go`) gives `[]spec.InventoryFact` with the Intake class | already extracted by 66u/66j: no extraction, K7 calls `factsOf` |
| `NextInventory` (`internal/spec/inventory.go:103`); version n > 1 refused | `inventory.go:103`; the bet wraps it as `nextVersion` (`rendered.go:46`) and refuses a later bet (`bet.go`, "a later bet's rendered task is row 66t's") | none |
| the precedent of a binding in the register commit | `archStage` (`archtasks.go:121`) with `openSpec.extra` | none |
| 66r's opening core `openTask(ctx, openSpec)` with the reuse-by-kind in `openShapeTask` | `shapetask.go:77`, `:87`, `:100`; `reuse` is a field of `openSpec`, false for the architecture tasks of the bet | none for the reuse. **One gap:** `openTask` draws its own ID (`s.taskID`), and K7's ID is reserved in `spec/reopenings.tsv` before the issue, so `openSpec` gains `id` |
| `taskID` / `drawTaskID` | `probe.go:182`, `:163`; draws avoid `sessions.tsv` and `tasks.tsv` | extended: draws also avoid `spec/reopenings.tsv`, and a batch of draws avoids itself (review condition 4) |
| `finalID` / `finalIDs` | `shapetask.go:40`, `:54` | none |
| `AdvanceShapeTask` | `shapeloop.go`; K7 calls none of it | none |
| `stepInstructions`, `inputsWanted`, `buildPrompt`, `PromptFile`, `taskInputs`, `factsRecords`, `nextFactsID` | `instructions.go:29`, `:158`; `promptfile.go:172`, `:437`, `:600`, `:312`, `:669` | lines moved only |
| `FinishTask`, `MergeTask`, `RecordMerged` | `finish.go:45`, `ready.go:195`, `:257` | **FinishTask treats every task of kind `specification` as the original one**: `finish.go:87` asks for the rendered commit of the first bet and `shapeloop.go:743` for the pinned `prd-lint` of a specification that waits for its bet. A reopened task is of kind `specification` and has no rendered commit; 66k tells the two apart with `reopenedTasks` (K7's) when it drives a reopened task. Not changed here (reported) |
| `registerTask` `register.go:267`, `registerFiles` `:213`, `AddTask` `:194`, `loopPermissionsStep` `:254`, `parameterFiles` `parameters.go:374`, `Fenced` `lease.go:118`, `specInput` `specinput.go:60`, `FullInput` `full.go:63`, `needs` `:494`, `measures` `:520`, the waits heading `full.go:352` | all at the lines in the column | none |
| K3 (`UpToDate` `baseup.go:43`, `TipOf` `:20`, the opening pin `pinPath` `redo.go:312`) | landed | K7 uses none. `redoIssue` (`redo.go:152`) composes a specification task's issue only from an opening pin; a reopened task has none (it is not `reuse`), so a redo of a reopened task is an error before any write: a known limit of M2e, stated in `loop.md` and tested |
| the `merged` event (`records.ReadEvents`) | `ready.go:257` writes it; column 5 is the merge commit, column 6 the pull request | none |

No step of the plan is impossible, and the plan, the review and the decisions settle every decision this task makes, with one exception that I record and settle by the plan review's own rule (condition 4, "an existing task row must agree with the reservation"): a restart that finds a register row of a reserved task ID compares its kind, requirements and predecessors with the reservation and stops on a difference.
