# Operator answer: revised routing policy for F-0006

I choose option **(a): revise the policy**.

F-0005 stays immutable. The original end of §4 is lost. Do not reconstruct it or change its hash. Capture this answer as a new raw fact, **F-0006**, and use F-0006 as the new operator input for ADR-0012.

The combined routing matrix in this answer supersedes my earlier routing answers where they conflict.

## 1. Answers to the 17 findings

| # | Answer |
|---|---|
| 1 | Routing is deterministic. A script applies the routing matrix. An LLM is a worker and does not enforce the route. |
| 2 | The harness commands are `claude`, `devin`, and `opencode`. Remove AGY. Do not use the names “Davinci CLI” or “OpenCoder CLI.” |
| 3 | Each route names an exact harness, model ID, effort, and account alias. A model class alone is not routable. |
| 4 | Do not use “latest Haiku,” Haiku, or Opus 4.x. A simple edit never moves to a reasoning model. |
| 5 | An extreme situation exists when no normal eligible row remains and the task is about to become stale. The Judge route is the last model route before human escalation. |
| 6 | The Astra Judge model is `gpt-6-astra-high` on Devin. Do not use the unclear names “Astra GPT,” “Astra 6,” or “Astra 6.1.” |
| 7 | Do not use a general “free tier” class. Use only an exact free-model row. A prepaid allowance is not the same as a free model. |
| 8 | Prompt formatting and context packaging are deterministic template work. Remove trade-off scoring. Compare options by stated considerations, not by a numeric score. |
| 9 | Do not require daily or weekly token or cost figures that the harnesses do not provide. Record an unavailable figure as `not provided`. Do not route from a missing figure. |
| 10 | Quota and budget are separate. Remove the 70% and 90% bands. Each harness state is `available`, `exhausted`, `unavailable`, or `unknown`. `unknown` is not `available`. |
| 11 | The Anthropic-only restriction applies only to the `claude` command. Failover takes the next eligible matrix row. If no row remains, create the red-flag issue defined below. |
| 12 | A request for human input is not a stall. A stall is a deadlock or a task that makes no progress within the specified limit. |
| 13 | Do not wait for all harnesses. Skip an exhausted or unavailable harness immediately. |
| 14 | A diagnostic package contains the number of distinct options actually returned. It can contain zero, one, two, or more options. Do not invent options A, B, and C. |
| 15 | PSB Invariant 9 means that a harness is replaceable. Do not rename it “reciprocal verification.” Review independence is a separate rule. A review model must differ from the author’s model. |
| 16 | Never commit secrets or sensitive prompt content. Commit only a redacted manifest with references where necessary. |
| 17 | Use the timeouts, effort settings, account aliases, external-data rule, availability test, and routing matrix in this answer. |

## 2. Source and capture decision

The stored F-0005 paste is incomplete:

- its Markdown fence does not close;
- its separator appears to come from the terminal display;
- I no longer have the missing original text.

Do not claim that the missing text was recovered. F-0006 is a new policy, not a reconstruction of the lost text.

## 3. Harness accounts and quota preference

The account aliases are:

- `claude-max`: Claude Code Max 20x, as stated by the Operator;
- `devin-pro`: Devin Pro;
- `opencode-go`: OpenCode Go $10, as stated by the Operator.

These names record the plans, not login names or email addresses.

A stated plan size does not prove the remaining quota. Use only a current harness response to determine exhaustion.

## 4. Row eligibility

A routing row is eligible only when all these conditions are true:

1. The command is installed.
2. The exact model is listed by that harness.
3. Authentication is valid.
4. The harness has not reported quota exhaustion.
5. The server responds.
6. The binding has not already failed this task.
7. The model family is permitted for the step.
8. For review work, the model is different from the author’s model.
9. The prompt does not send secrets or prohibited sensitive data to an external service.

A model-list entry proves only that the model is listed. It does not prove quota, availability, price, or suitability.

## 5. Main selection rule

Take the first eligible row in the order shown below.

If a harness reports quota exhaustion or its server is unavailable, skip that row immediately.

If a dispatch gives no response within five minutes:

1. cancel it or confirm that it stopped;
2. record the outcome;
3. start the next eligible row.

Never leave two writers active on the same task.

If no row remains, create the red-flag issue and stop the step.

## 6. Routing matrix

### 6.1 Suggest a solution

Use these rows in order:

1. `claude --model claude-opus-5-5 --effort high`
2. Devin `gpt-6-sol-high`
3. Devin `claude-opus-5-5-high`
4. `opencode-go/grok-4.7#high`

Use `xhigh` for Opus 5.5 if the selected solution can become an ADR.

### 6.2 Specification plan

Use these rows in order:

1. `claude --model claude-opus-5-5 --effort high`
2. Devin `claude-opus-5-5-xhigh`
3. Devin `gpt-6-sol-xhigh`
4. `opencode-go/grok-4.7#xhigh`

### 6.3 Review the specification plan

Use these rows in order:

1. Devin `gpt-6-sol-xhigh`
2. `opencode-go/grok-4.7#xhigh`
3. Second lens: `opencode-go/glm-5.3#max`

The reviewer model must differ from the model that wrote the plan. Claude Opus must not review a plan written by Claude Opus.

### 6.4 Review the PR or MR

Use one lens in each round.

**Round 1:**

1. Devin `gpt-6-sol-xhigh`
2. `opencode-go/gpt-6-luna#high`
3. `opencode-go/grok-4.7#xhigh`
4. `opencode-go/glm-5.3#max`
5. `opencode-go/qwen3.8-max#xhigh`

**Round 2 after a fix:**

1. `claude --model claude-opus-5-5 --effort high`, only when the author model is not from the Claude family;
2. Devin `claude-opus-5-5-high`, with the same family restriction;
3. then the eligible OpenCode rows above.

Use `xhigh` for the first blind round on a frozen head. Use `high` for later rounds after a fix.

The reviewer model must always differ from the author model. Independence has priority over effort.

### 6.5 Judge: the last model step before the human

If the author model is not from the Claude family, use:

1. `claude --model claude-fable-5-1 --effort xhigh`
2. Devin `claude-fable-5-1-xhigh`
3. Devin `gpt-6-astra-high`

If the author model is from the Claude family, use:

1. Devin `gpt-6-astra-high`

There is no OpenCode Judge row. Grok is not a Judge.

Do not use Fable or Astra for ordinary work. They are last-resort Judge models.

If no Judge row is eligible, create the red-flag issue immediately.

## 7. Execution classification

Each execution step has exactly one class.

### Hard

The specification plan explicitly marks the step `hard`. No other condition makes it hard until a later operator decision defines one.

Route:

1. Devin `fusion-claude-opus-5-5-high-sidekick-swe-2-medium`
2. If it cannot run, use the Ordinary rows.

### Cheap mechanical

A simple, low-risk edit that needs little or no reasoning.

Examples:

- rename text;
- apply formatting;
- update a known value;
- make a precise one-line change.

Route:

1. Devin `gpt-6-luna-medium`
2. If it cannot run, use the Ordinary rows.

### Free

Use only when:

- the data is not confidential or sensitive;
- the task is suitable for the free model; and
- the specification plan explicitly selects free execution.

Route:

1. `opencode-go/space-bunny-free`
2. If it cannot run, use the Ordinary rows.

### Ordinary

The default for execution work that is not Hard, Cheap mechanical, or explicitly selected as Free.

Route:

1. Devin `swe-2-high`
2. `claude --model claude-sonnet-5 --effort high`
3. `opencode-go/glm-5.3#high`
4. `opencode-go/kimi-k2.7-code`

### Execution priority

Apply these rules in order:

1. If the specification plan says `hard`, use Hard.
2. Otherwise, if the edit is cheap mechanical, use Cheap mechanical.
3. Otherwise, if the plan explicitly selects Free and the data is not confidential or sensitive, use Free.
4. Otherwise, use Ordinary.

A simple edit never steps up to a reasoning model. Execution never uses Fable or Astra.

## 8. Effort rules for judgment work

Solution selection, specification planning, plan review, and PR/MR review are judgment work.

1. Do not use `low` or `medium` for these steps.
2. Use `high` by default.
3. Use `xhigh` when the decision can become an ADR or for the first blind review round.
4. Use `max` only for high-risk work involving security, data loss, or a costly or irreversible operation.
5. Independence has priority over effort.
6. Record the actual model and effort in the ADR-0007 resource record.

## 9. Timeout and stall rules

| State | Limit | Action |
|---|---:|---|
| No response | 5 minutes | Cancel or confirm termination, record it, and use the next row. |
| In progress but no progress | 10 minutes | Stop the run and ask the Judge whether to continue or select another action. |
| No useful time estimate | No unbounded wait | Ask the Judge. |
| Panel | 15 minutes total | Each member has at most 5 minutes. |
| Quota exhaustion or server unavailable | Immediate | Skip to the next row. |
| No eligible harness or model | Immediate | Create the red-flag issue and stop. |

A task that asks for human input is not automatically a stall.

## 10. Red-flag notification

When no eligible row remains, create a Git issue.

The issue must:

- have the `red-flag` label;
- have a severity label;
- identify the blocked task and step;
- list every row attempted or skipped;
- record each harness state;
- record the time and last observed progress;
- state the one decision or action required from the Operator;
- contain no secret or sensitive data.

Use:

- `severity:blocker` when work has stopped because no eligible row exists;
- `severity:high` when a harness failed but work continues on a fallback.

The repository does not currently have these labels. Create `red-flag`, `severity:blocker`, and `severity:high` before the label rule becomes active. Until then, put `RED FLAG` and the severity in the issue title and body.

## 11. Diagnostic package

Never escalate a bare question.

The package contains:

- stall context;
- failed attempts;
- harness and model used for each attempt;
- quota or server reports;
- elapsed time;
- each distinct option actually returned;
- the Judge recommendation, if the Judge ran;
- the decision required from the Operator.

The number of options is not fixed. It can be zero.

## 12. Sensitive data and external services

Before a dispatch to an external service:

1. send only the context required by the step;
2. remove secrets and credentials;
3. do not send confidential or sensitive data to a free model;
4. do not commit sensitive prompt contents;
5. commit only a redacted manifest and repository references.

## 13. Matrix storage

Store the ordered matrix as structured, versioned repository data, not as prose and not only as shell environment variables.

Use:

`docs/setup/routing-matrix.tsv`

Each row records:

- order;
- role or execution class;
- harness;
- exact model ID;
- effort;
- account alias;
- fallback status;
- confidentiality restriction;
- enabled state.

Environment variables can select or override a row at runtime. They are not the system of record.

A model, price, account, or harness change updates the matrix. It does not require rewriting this policy unless the routing rule itself changes.

## 14. Do not use

Do not use:

- GPT-4.1, GPT-5.1, GPT-5.2, or GPT-5.4;
- Opus 4.5 through 4.8;
- Opus 5;
- Sonnet 4.5 or 4.6;
- Haiku;
- deprecated models;
- `-fast` or `-priority` variants;
- `adaptive`, because ADR-0007 must identify the actual model;
- paid Zen models while that account has no funds;
- DeepSeek on OpenCode Go while the region restriction blocks it.

Claude Code can run only models that it lists. Do not try to run Sol, Luna, Astra, SWE-2, Grok, GLM, Qwen, or Kimi through the `claude` command.

## 15. Evidence status

Benchmark values are operator rationale, not verified repository facts. Record each benchmark with its source URL, publication date, read date, harness, model, and effort before citing it in an ADR.

A model listing proves only that the model is listed. A smoke run proves only that the binding returned once. Neither proves general capability.

No smoke-run recency window is set by this policy.

## 16. Supersession

This answer supersedes my earlier operator routing answers where they conflict.

In particular:

- GPT-6 Sol is again permitted as the plan reviewer and first PR/MR reviewer.
- Grok, GLM, Luna, and Qwen are permitted as the named OpenCode review fallbacks.
- AGY remains removed.
- The 70% and 90% quota bands remain removed.
- Fable and Astra remain last-resort Judge models.
- F-0005 remains immutable.

Capture this answer as F-0006 and use it as the current operator input for the routing ADR.
