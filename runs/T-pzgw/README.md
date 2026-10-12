# T-pzgw — the evidence of the cause

The first 5 lines of `result/spans.tsv`, as written (tabs shown as tabs) by the session `S-445dee66` (`product-owner`, `claude-opus-5-5`) in run 3
of the demo of `M2c` (#228), at 2026-10-11 21:43Z. It was copied by a watcher
before the session's directory was removed. The problem statement is the public
`PSB-CHAT-001` of `pharzam/chat-orchestrator`.

```text
start	end	fact	class	reason
0	19	—	—	—
21	89	—	—	—
91	136	—	—	—
139	168	1	context	—
```

The whole file had 771 rows: 696 facts, each with its class and the reason `—`,
as asked; and 75 spans that are not facts, each with `—` as its fact, its class
**and its reason**. Line 2 is the first of them, so `CheckSpan` refused the file
there: "reason is the empty value exactly when fact is not". Runs 1 and 2
(sessions `S-2742cb40` and `S-37a7776c`) were refused at the same line and
column.

The red of `TestTheNumberingPromptGivesEachColumnOnce` on the old prompt, at
`347d845`:

```
--- FAIL: TestTheNumberingPromptGivesEachColumnOnce
    numbering_test.go:58: the prompt lacks "A span that is a fact has its number in the column fact, …"
    numbering_test.go:58: the prompt lacks "A span that is not a fact (a heading, a rule line) has \"—\" in the column fact and …"
    numbering_test.go:62: the prompt still holds the sentence that reads two ways: "has \"—\" as its fact and its class, and its reason"
FAIL	github.com/pharzam/layup/internal/run
```

## The other prompts of Intake (D2)

Each prompt of a session of Intake was read for the same fault, a list of
columns in one sentence that can be read as one value for all:

- `meaningPrompt` (`internal/run/meaning.go`), `draftPrompt`
  (`internal/run/draft.go`) and `reviewPrompt` (`internal/run/followup.go`) name
  no `—`, and each column of their header has its own clause;
- `completenessPrompt` (`internal/run/completeness.go`) gives each `—` in
  parentheses beside the item that fills it, for example "uncovered, … (its
  fact, the requirement "—")", one pair for each kind.

None reads two ways, so none changes.
