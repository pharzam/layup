# The test runs of T-bpxg

## The recorded usage reports (2026-10-09T18:14Z)

`claude --version` printed `2.1.295 (Claude Code)`. Two `result` events of
Claude Code 2.1.295 are the test data of `internal/session/testdata/`:

- `claude-result-subagent.jsonl`: recorded by this task, in an empty scratch
  directory, with `ANTHROPIC_DEFAULT_HAIKU_MODEL=claude-fable-5-1`:
  `claude -p --model claude-fable-5-1 --output-format stream-json --verbose
  --permission-mode dontAsk --setting-sources project,local
  --no-session-persistence --max-budget-usd 1 --allowedTools "Agent" "Task"`,
  the prompt asking for one Agent call with the model `opus`. Exit 0;
  `total_cost_usd` 0.27055324999999997; `modelUsage` names
  `claude-fable-5-1` and `claude-opus-5-5`, and no other model. The file is
  the `result` line of the stream, byte for byte. The rest of the stream held
  host paths (the session's directories), so it is not kept.
- `claude-result-one-model.jsonl`: the `result` line of round 1 of
  `T-6sbe` (#160), a review session of this milestone, one model,
  `claude-fable-5-1`. Its field `result`, the text of the review record,
  which no reader reads, is replaced by a marker; each other field is as
  recorded.

Both files were read for a host path, an account detail and an e-mail address
(none), and `gitleaks dir` found no leak.

The sums by hand: one model, 290 in, 39,003 out, 775,649 + 145,871 = 921,520
cache; the subagent event, 34 + 2 = 36 in, 157 + 4 = 161 out,
31,161 + 9,863 + 0 + 11,445 = 52,469 cache.

## Red 1: the tests before the functions (2026-10-09T18:16Z)

`go vet` does not build the tests, for want of the functions they name:

```
vet: internal/session/end_test.go:76:13: undefined: Usage
```

## Two cases made to fail on their own rule (2026-10-09T18:17Z to 18:19Z)

The first green came with no red of the cases of `TestResultOf`: the file over
1 MiB was also malformed, so another rule could refuse it. Each case now
asserts the words of its own error, and the file over 1 MiB is a valid result
of 1,048,577 bytes. Two rules had no case: the cost as the text of the JSON
number (a float gives the same text for the recorded costs), now a case of
`1.50`; and the line of up to 8 MiB, now a line of 100 KiB before the result.

## The mutations (2026-10-09T18:19Z to 18:22Z)

A mutation of each rule of `end.go`, one at a time, run by a script that
restores the file from a copy after each. The unit rules run the unit tests
(`usage-sort` ten times over, as a map's order is random); the others
`go test -tags=integration`. Each test line in full:

```
== class-start: Class: ErrStart not start (exit 1)
end_test.go:33: a program that cannot start: "", want "start"
== class-own: Class: an error of layup run's own read as a class (exit 1)
end_test.go:33: an error of layup run's own: "done", want ""
== class-stop: Class: a stop not first of the classes that follow the start (exit 1)
end_test.go:33: stopped at wall, exit 0, a valid result: "done", want "wall"
== class-exit: Class: a non-zero exit, or a signal (exit -1), not a crash (exit 1)
end_test.go:33: a non-zero exit: "done", want "crash"
end_test.go:33: a signal that layup run did not send: "done", want "crash"
== class-noresult: Class: no valid result file still done (exit 1)
end_test.go:33: exit 0 with no valid result: "done", want "no-result"
== result-link: ResultOf: a link followed (exit 1)
end_integration_test.go:80: a result.tsv a link: <nil>, want an error with "too many levels of symbolic links"
== result-regular: ResultOf: a file that is not regular read (exit 1)
end_integration_test.go:80: a result.tsv a directory: read /tmp/TestResultOf4255018099/001/sessions/S-1a2b3c4d/result/result.tsv: is a directory, want an error with "is not a regular file"
== result-size: ResultOf: no limit of size (exit 1)
end_integration_test.go:80: a result.tsv over 1 MiB: <nil>, want an error with "is over 1 MiB"
== result-sizeeq: ResultOf: a file of exactly 1 MiB refused (exit 1)
end_integration_test.go:88: a result.tsv of 1 MiB (1048576 bytes): the result file /tmp/TestResultOf3916852063/001/sessions/S-1a2b3c4d/result/result.tsv is over 1 MiB
== result-probe: ResultOf: a probe reads result.tsv (exit 1)
end_integration_test.go:54: a valid probe.tsv: [], line 1, column "value": the header row has "n" in its place
== probe-token: ReadProbeResult: two token rows pass (exit 1)
end_test.go:120: two tokens: no error
== usage-last: usageOf: the first result object, not the last (exit 1)
end_test.go:73: the last result object, after lines that are not JSON: [1 1 2 observed  9 USD [m]], <nil>; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
== usage-skip: usageOf: a line that is not JSON ends the read (exit 1)
end_test.go:73: the last result object, after lines that are not JSON: [<nil> <nil> <nil>     []], invalid character 'o' in literal null (expecting 'u'); want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: a line of 100 KiB before the result: [<nil> <nil> <nil>     []], invalid character 'x' looking for beginning of value; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: no result object: [<nil> <nil> <nil>     []], invalid character 'o' in literal null (expecting 'u'); want [<nil> <nil> <nil> unavailable stdout holds no result object   []]
== usage-line: usageOf: lines of the scanner's default 64 KiB (exit 1)
end_test.go:73: a line of 100 KiB before the result: [<nil> <nil> <nil>     []], bufio.Scanner: token too long; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
== usage-cache: usageOf: cacheCreationInputTokens not in the cache class (exit 1)
end_test.go:73: one model: [290 39003 775649 observed  5.0643822499999995 USD [claude-fable-5-1]], <nil>; want [290 39003 921520 observed  5.0643822499999995 USD [claude-fable-5-1]]
end_test.go:73: a subagent on a second model: [36 161 31161 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]], <nil>; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: the last result object, after lines that are not JSON: [36 161 31161 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]], <nil>; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: a line of 100 KiB before the result: [36 161 31161 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]], <nil>; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: a model that lacks a field: [5 7 9 observed    [a b]], <nil>; want [5 7 <nil> partial b has no cacheCreationInputTokens   [a b]]
end_test.go:73: each class lacked: [<nil> <nil> <nil> unavailable a has no inputTokens, outputTokens, cacheReadInputTokens   [a]], <nil>; want [<nil> <nil> <nil> unavailable a has no inputTokens, outputTokens, cacheReadInputTokens, cacheCreationInputTokens   [a]]
== usage-every: usageOf: a lacking field summed as 0 (exit 1)
end_test.go:73: a model that lacks a field: [5 7 16 observed b has no cacheCreationInputTokens   [a b]], <nil>; want [5 7 <nil> partial b has no cacheCreationInputTokens   [a b]]
end_test.go:73: each class lacked: [0 0 0 observed a has no inputTokens, outputTokens, cacheReadInputTokens, cacheCreationInputTokens   [a]], <nil>; want [<nil> <nil> <nil> unavailable a has no inputTokens, outputTokens, cacheReadInputTokens, cacheCreationInputTokens   [a]]
== usage-none0: usageOf: no class summed is partial (exit 1)
end_test.go:73: each class lacked: [<nil> <nil> <nil> partial a has no inputTokens, outputTokens, cacheReadInputTokens, cacheCreationInputTokens   [a]], <nil>; want [<nil> <nil> <nil> unavailable a has no inputTokens, outputTokens, cacheReadInputTokens, cacheCreationInputTokens   [a]]
== usage-nomodel: usageOf: a report with no model read as observed (exit 1)
end_test.go:73: a result with no modelUsage: [0 0 0 observed  1.50 USD []], <nil>; want [<nil> <nil> <nil> unavailable the report names no model 1.50 USD []]
end_test.go:73: a result with an empty modelUsage: [0 0 0 observed    []], <nil>; want [<nil> <nil> <nil> unavailable the report names no model   []]
== usage-sort: usageOf: the models in the map's order (exit 1)
end_test.go:73: a subagent on a second model: [36 161 52469 observed  0.27055324999999997 USD [claude-opus-5-5 claude-fable-5-1]], <nil>; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: a subagent on a second model: [36 161 52469 observed  0.27055324999999997 USD [claude-opus-5-5 claude-fable-5-1]], <nil>; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: a model that lacks a field: [5 7 <nil> partial b has no cacheCreationInputTokens   [b a]], <nil>; want [5 7 <nil> partial b has no cacheCreationInputTokens   [a b]]
end_test.go:73: a model that lacks a field: [5 7 <nil> partial b has no cacheCreationInputTokens   [b a]], <nil>; want [5 7 <nil> partial b has no cacheCreationInputTokens   [a b]]
end_test.go:73: a model that lacks a field: [5 7 <nil> partial b has no cacheCreationInputTokens   [b a]], <nil>; want [5 7 <nil> partial b has no cacheCreationInputTokens   [a b]]
end_test.go:73: a line of 100 KiB before the result: [36 161 52469 observed  0.27055324999999997 USD [claude-opus-5-5 claude-fable-5-1]], <nil>; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: a model that lacks a field: [5 7 <nil> partial b has no cacheCreationInputTokens   [b a]], <nil>; want [5 7 <nil> partial b has no cacheCreationInputTokens   [a b]]
end_test.go:73: the last result object, after lines that are not JSON: [36 161 52469 observed  0.27055324999999997 USD [claude-opus-5-5 claude-fable-5-1]], <nil>; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: a model that lacks a field: [5 7 <nil> partial b has no cacheCreationInputTokens   [b a]], <nil>; want [5 7 <nil> partial b has no cacheCreationInputTokens   [a b]]
end_test.go:73: a subagent on a second model: [36 161 52469 observed  0.27055324999999997 USD [claude-opus-5-5 claude-fable-5-1]], <nil>; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: a model that lacks a field: [5 7 <nil> partial b has no cacheCreationInputTokens   [b a]], <nil>; want [5 7 <nil> partial b has no cacheCreationInputTokens   [a b]]
end_test.go:73: a subagent on a second model: [36 161 52469 observed  0.27055324999999997 USD [claude-opus-5-5 claude-fable-5-1]], <nil>; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: a model that lacks a field: [5 7 <nil> partial b has no cacheCreationInputTokens   [b a]], <nil>; want [5 7 <nil> partial b has no cacheCreationInputTokens   [a b]]
end_test.go:73: a subagent on a second model: [36 161 52469 observed  0.27055324999999997 USD [claude-opus-5-5 claude-fable-5-1]], <nil>; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: a model that lacks a field: [5 7 <nil> partial b has no cacheCreationInputTokens   [b a]], <nil>; want [5 7 <nil> partial b has no cacheCreationInputTokens   [a b]]
end_test.go:73: a model that lacks a field: [5 7 <nil> partial b has no cacheCreationInputTokens   [b a]], <nil>; want [5 7 <nil> partial b has no cacheCreationInputTokens   [a b]]
== usage-cost: usageOf: the cost through a float (exit 1)
end_test.go:73: a result with no modelUsage: [<nil> <nil> <nil> unavailable the report names no model 1.5 USD []], <nil>; want [<nil> <nil> <nil> unavailable the report names no model 1.50 USD []]
== usage-nocost: usageOf: no cost read (exit 1)
end_test.go:73: one model: [290 39003 921520 observed    [claude-fable-5-1]], <nil>; want [290 39003 921520 observed  5.0643822499999995 USD [claude-fable-5-1]]
end_test.go:73: a subagent on a second model: [36 161 52469 observed    [claude-fable-5-1 claude-opus-5-5]], <nil>; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: the last result object, after lines that are not JSON: [36 161 52469 observed    [claude-fable-5-1 claude-opus-5-5]], <nil>; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: a line of 100 KiB before the result: [36 161 52469 observed    [claude-fable-5-1 claude-opus-5-5]], <nil>; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: a result with no modelUsage: [<nil> <nil> <nil> unavailable the report names no model   []], <nil>; want [<nil> <nil> <nil> unavailable the report names no model 1.50 USD []]
== usage-none: usageOf: none with no reason (exit 1)
end_test.go:93: none: [<nil> <nil> <nil> unavailable    []], <nil>; want [<nil> <nil> <nil> unavailable the harness reports none   []]
== usage-format: usageOf: another format read as claude-result (exit 1)
end_test.go:73: one model: [<nil> <nil> <nil>     []], the usage format "claude-result" is not one of the list; want [290 39003 921520 observed  5.0643822499999995 USD [claude-fable-5-1]]
end_test.go:73: a subagent on a second model: [<nil> <nil> <nil>     []], the usage format "claude-result" is not one of the list; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: the last result object, after lines that are not JSON: [<nil> <nil> <nil>     []], the usage format "claude-result" is not one of the list; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: a line of 100 KiB before the result: [<nil> <nil> <nil>     []], the usage format "claude-result" is not one of the list; want [36 161 52469 observed  0.27055324999999997 USD [claude-fable-5-1 claude-opus-5-5]]
end_test.go:73: no result object: [<nil> <nil> <nil>     []], the usage format "claude-result" is not one of the list; want [<nil> <nil> <nil> unavailable stdout holds no result object   []]
end_test.go:73: a result with no modelUsage: [<nil> <nil> <nil>     []], the usage format "claude-result" is not one of the list; want [<nil> <nil> <nil> unavailable the report names no model 1.50 USD []]
end_test.go:73: a result with an empty modelUsage: [<nil> <nil> <nil>     []], the usage format "claude-result" is not one of the list; want [<nil> <nil> <nil> unavailable the report names no model   []]
end_test.go:73: a model that lacks a field: [<nil> <nil> <nil>     []], the usage format "claude-result" is not one of the list; want [5 7 <nil> partial b has no cacheCreationInputTokens   [a b]]
end_test.go:73: each class lacked: [<nil> <nil> <nil>     []], the usage format "claude-result" is not one of the list; want [<nil> <nil> <nil> unavailable a has no inputTokens, outputTokens, cacheReadInputTokens, cacheCreationInputTokens   [a]]
== usage-noresult: usageOf: no result object gives zero counts (exit 1)
end_test.go:73: no result object: [<nil> <nil> <nil> unavailable the report names no model   []], <nil>; want [<nil> <nil> <nil> unavailable stdout holds no result object   []]
```

A first run had three that did not do what they should. `Class` tested a
signal beside a non-zero exit, and a signal always gives exit -1, so the test
of the signal is gone. `ReadProbeResult` refused an empty value, which
`tsv.Read` already refuses in a key column, so that check is gone. The
mutation of the sort did not build, and is now a call that keeps the order.

## Green (2026-10-09T18:23Z)

Each with exit 0: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd`
(empty), `go test ./...`, `go test -tags=integration ./...` (with
`TestPackageRules`, `TestInputRule` and the block test of `internal/tsv`);
`adr-lint`, `prd-lint`, `link-lint`, `setup-check`,
`run-discipline-tests`, `git diff --check`.
