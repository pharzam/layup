# T-qjjd: the conditions of the second plan review of #475, and their evidence

One row per condition of the plan review (comment 6336309085). This table is check K6's input; K6 proves that each cited file exists, not that it meets the condition.

| Condition | Evidence |
| --- | --- |
| 1: one input contract for the replacement rule (events, session rows with roles and bases, kept-result status, task kind), given to `findingsOf`, `PromptFile` (now with the kind) and `AnswerAccepted` through `Position` and `OpenQuestions`; the answer-session and earlier-asking exclusions are kept | `internal/run/shapeloop.go:326`, `internal/run/shapeloop.go:364`, `internal/run/shapeloop.go:253`, `internal/run/promptfile.go:569`, `internal/run/answer.go:125`, `internal/run/tier.go:147`, `internal/run/shapeloop.go:418` |
| 2: r2 counts the eligible end that caused each earlier replacement; the reverse-stage case, the single-fact negatives, the crashed answer session of an eligible role, two consecutive replacements through the findings reader, the rendered answer IDs and the answer acceptance, and the mixed oracle "ended no-result" | `internal/run/shapeloop.go:326`, `internal/run/shapeloop_test.go:367`, `internal/run/shapeloop_test.go:307`, `internal/run/shapeloop_test.go:406`, `internal/run/answer_test.go:225`, `internal/run/questionloop_integration_test.go:459`, `internal/run/shapeloop_integration_test.go:649` |
| 3: the tests first, with only the scaffolding to compile and the harness controls in the same slice; assertion failures recorded before the rule; guards, behavioural reds and compile failures kept apart | `runs/T-qjjd/test-runs.md:1`, `internal/run/shapetask_integration_test.go:64` |
| 4: the session table says an eligible task-session end gets the new attempt, with a pointer to the rule's home | `docs/spec/session.md:322`, `docs/spec/loop.md:2835` |
