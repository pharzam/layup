# T-3ff9 clause-by-clause semantic comparison

Author comparison by GPT-6.1 Sol; the independent frozen-head semantic review is pending. This table is not a review verdict. Sources below abbreviate `docs/adr/0028-delegate-two-bounded-landing-decisions.md` as ADR-0028 and `docs/engineering-discipline.md` as ED. Baseline citations name `f45fffc`; other line numbers refer to the Stage B text. The Operator's option-1 decision on [#447](https://github.com/pharzam/layup/issues/447), supplied with this stage, resolves the ambiguity identified by plan-review condition 2; ADR-0028 alone does not define that boundary.

| ADR-0028 source | Changed clause | Same meaning? | Source still says it? | Unsupported assertion? |
| --- | --- | --- | --- | --- |
| ADR-0028:32 | ED:373: plan review sets the maximum; whole-branch added plus removed lines, files and named base remain the unit. | Yes; authority retained from baseline ED:353. | Yes. | No. |
| ADR-0028:33 | ED:378: Markdown, TSV and text run/time/output evidence under the task's runs directory leaves both totals. | Yes. | Yes. | No. |
| ADR-0028:33 | ED:380: a file under runs executed or consumed as input stays counted; both script and script-read-file examples remain. | Yes, with the Operator's clarification. | Yes; the clarification fixes the meaning of input. | No. |
| ADR-0028:36 | ED:382: executable tests and their inputs stay counted. | Yes. | Yes. | No. |
| ADR-0028:33, :36; Operator option 1 | ED:385: execution data defines input once for both uses. | Yes with the decision; no claim the ADR defined it alone. | Yes in the supplied Operator decision. | No. |
| ADR-0028:33, :40, :57; Operator option 1 | ED:387: inspection solely of form, citations or evidence does not make a record an input. | Yes with the decision; preserves both ADR examples and consequence. | Yes in the supplied Operator decision. | No. |
| Operator option 1; amendment 1 option 1 | ED:388: dual-purpose files remain inputs. | Yes; the narrower reading does not exclude execution data. | Yes; the selected option expressly retains dual use. | No. |
| ADR-0028:34 | ED:390: completed-log line, removed backlog line, own task verdict and resource record leave the count. | Yes; all four are named. | Yes. | No. |
| ADR-0028:34 | ED:392: mixed files lose only bookkeeping lines and still count as touched. | Yes in both dimensions. | Yes. | No. |
| ADR-0028:37 | ED:395: landing above the maximum within tolerance has advance approval. | Yes. | Yes. | No. |
| ADR-0028:37, :39, :52 | ED:396: 20 percent, 60-line floor, 500-line ceiling and four files. | Yes; accepted figures, not the rejected proposal. | Yes. | No. |
| ADR-0028:37 | ED:397: post figure, maximum and margin on the issue before merge. | Yes. | Yes. | No. |
| ADR-0028:30, :56 | ED:398: Operator may overrule the budget delegation there before merge. | Yes; no waiting interval added. | Yes. | No. |
| ADR-0028:39 | ED:401: one fixed bound for every task, no multiple, findings-driven growth or second approval. | Yes. | Yes. | No. |
| ADR-0028:37; baseline ED:356 | ED:404: beyond tolerance is a finding, not a revision; child issue unless one exact landing figure is approved once on the issue. | Yes; preserves the authority beyond delegation. | Yes, including the retained baseline route. | No. |
| ADR-0028:37; baseline ED:364 | ED:406: exact approval reaches only its named figure. | Yes; future growth is not authorized. | Yes, by retention. | No. |
| ADR-0028:37; baseline ED:358 | ED:407: unapproved overrun blocks with the existing verdict independently of the cap. | Yes; overrun here is beyond tolerance or exact approval. | Yes, by retention. | No. |
| ADR-0028:37; baseline ED:361 | ED:411: once per issue; successor needs its own approval and measured carried size in its plan review; discipline, not mechanism. | Yes. | Yes, by retention. | No. |
| ADR-0028:37; baseline ED:370 | ED:413: growth past an approved figure is absorbed, routed with non-material classification, or ends at the cap. | Yes; retained beyond the named figure. | Yes, by retention. | No. |
| ADR-0028:37, :39; review condition 3 | ED:294: either dimension beyond maximum plus tolerance without exact issue approval gives the non-merge verdict. | Yes; inside-tolerance figures no longer trigger this route. | Yes; applies Decision 2 consistently. | No. |
| ADR-0028:37, :39; review condition 3 | ED:360: in-scope fixes use the budget rule's permitted resulting landing figure; otherwise a child issue. | Yes; includes advance and exact approval routes. | Yes; applies Decision 2 consistently. | No. |
| ADR-0028:30, :40 | ED:310: material findings at cap go on the issue; Operator chooses among four ends except default end 3. | Yes. | Yes. | No. |
| ADR-0028:42; baseline ED:321 | ED:322: last-round fix lands and task record discloses that no round read it. | Yes. | Yes. | No. |
| ADR-0028:42 | ED:323: last verdict stays not mergeable, findings recorded. | Yes; no parser vocabulary change. | Yes. | No. |
| ADR-0028:40 | ED:326: at cap with material findings open, every one must be a records finding for the comment-free default. | Yes; nonempty and universal conditions retained. | Yes. | No. |
| ADR-0028:40 | ED:328: records fix restores a named source or record form and changes nothing else. | Yes. | Yes. | No. |
| ADR-0028:40 | ED:329: existing-test traceability/Test cells, index line, schema columns and trailing whitespace examples. | Yes; every example retained. | Yes. | No. |
| ADR-0028:41; Operator option 1 | ED:333: exit code, assertion, verdict, defined check input, test, product code, specification and authority-sentence exclusions. | Yes; every exclusion retained, with a pointer to the one input definition. | Yes; decision clarifies the input term only. | No. |
| ADR-0028:41, :43 | ED:336: excluded effects and disputed classification retain the Operator's end decision. | Yes. | Yes. | No. |
| ADR-0028:42 | ED:339: author records end 3 by ADR-0028 with findings and fixes on the issue. | Yes; literal record tag retained. | Yes. | No. |
| ADR-0028:30, :56 | ED:340: Operator may overrule before merge; default or individual comment permits merge. | Yes. | Yes. | No. |
| ADR-0028:42; baseline ED:324 | ED:343: end 3 remains the exception to re-freeze, last-round and Verdict-field rules. | Yes; mechanism retained while permission changes. | Yes, by retention. | No. |
| ADR-0028:42; baseline ED:327 | ED:345: after that fix no further fix or round; bookkeeping and merge of main remain allowed. | Yes. | Yes, by retention. | No. |
| ADR-0028:59; ED one home per rule | docs/glossary.md:90: cycle-cap definition and example point to canonical values, verdicts and ends. | Yes; no independent end rule. | Yes. | No. |
| ADR-0028:37, :39, :59 | docs/glossary.md:92: ceiling denotes size; budget, tolerance and authority point to the canonical section. | Yes; retires the no-advance-growth claim. | Yes. | No. |
| ADR-0028:37, :39; amendment 1 sweep | docs/glossary.md:99: in-change disposition points to the canonical section. | Yes; removes the inside-budget shorthand. | Yes. | No. |
| ADR-0028:59 | docs/adr/0007-record-task-resource-use.md:7 and docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md:7: Accepted, amended by ADR-0028; bodies untouched. | Yes. | Yes. | No. |
| ADR-0028:59; ADR index convention | docs/adr/README.md:146 and :165: only the two Status cells become Accepted; amended. | Yes. | Yes. | No. |
| ADR-0028:40, :59; ED one home per rule | AGENTS.md:59: protocol pointer replaces unconditional Operator end selection and duplicated cycle/verdict summary. | Yes; remaining review duties retained. | Yes. | No. |
| ADR-0028:37, :59; ED one home per rule | docs/tests/dod-checklist.md:93: landing/overrun pointer replaces held-or-individually-approved test; missing-budget caveat retained. | Yes. | Yes. | No. |
| ADR-0028:37, :39; amendment 1 sweep | docs/plan/README.md:595: expected lines remain estimates; size authority points to canonical rule. | Yes. | Yes. | No. |
| ADR-0028:33, :40; Operator option 1 | docs/guardrails.md:82: ambiguity lesson links both uses to the clarified canonical definition. | Yes; no second definition or new figure. | Yes in the decision and review condition 2. | No. |
| ADR-0028:37, :39; ED one home per rule | docs/guardrails.md:666: named-base procedure retained; overrun disposition becomes a budget pointer. | Yes; removes an unconditional overrun finding claim. | Yes. | No. |
| ADR-0028:42, :59 | docs/ci/review-record-lint.sh:222, :257, :265: parser unchanged because field presence, cycle accounting and verdict vocabulary/position stay valid. | Yes; no claim parser proves semantics. | Yes; checked against parser. | No. |

Author reading: the definition governs counting and end-3 eligibility together; auditing a record alone cannot cancel the ADR's evidence exclusion or examples. Every other eligibility exclusion remains in force. No second reading remains identified by the author; independent confirmation is still required. T-3ff9 itself uses the baseline budget and Operator stopping route until this change lands (ADR-0028:59), not its proposed delegation.
