# T-c3zv: the conditions of the plan review of #473, and their evidence

One row per condition of the plan review (the last comment of #473). This table is check K6's input; K6 proves that each cited file exists, not that it meets the condition.

| Condition | Evidence |
| --- | --- |
| 1: the literal prompt text `the single character "—" (U+2014)`; the test renders both prompts, extracts the token between the quotes, puts its bytes into a fail row read by `spec.ReadOptions`; a `-` control the reader refuses; no invented delimiter, no repair | `internal/run/options.go:276`, `internal/run/options.go:296`, `internal/run/options_test.go:183`, `internal/run/options_test.go:207` |
| 2: "never a hyphen" is not in the pass clause; the refusal rule is about fail evidence only | `internal/run/options.go:275`, `internal/run/options.go:295` |
| 3: the options session of the specification carries the complete rule (both prompts require exactly `—` (U+2014) for a fail's evidence; `-` or other text is refused) | `docs/spec/loop.md:1801` |
| Note 5: the blanket claim that only these two prompts are words-only is not repeated | `docs/guardrails.md:1508` |
| Demo (R11): both prompts name `—` (U+2014) and a fail row with that character is read by `spec.ReadOptions` | `runs/T-c3zv/test-runs.md:1`, `internal/run/options_test.go:183` |
