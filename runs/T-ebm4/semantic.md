# T-ebm4 author semantic comparison

This is the author's preparation for independent semantic review, not that review
or a verdict. Sources are the canonical "Reviewing until findings decay" section
of docs/engineering-discipline.md at baseline
abdbe996282b33080406223ae8ba19e5a3634360, the supplied plan review and amendment 1
on #451, and the Operator's replacement decision on the time to overrule
(comment 6301262948). Baseline locations below stay stable as the prose rewraps.

## Canonical text and authority

| Changed clause | Source and clause-by-clause comparison |
| --- | --- |
| Gate step 5: freeze and one blind round; protocol defines cycles, verdicts and ends. | Baseline canonical section, lines 268-354: first freeze and independent round remain. The pointer includes all endings, including end 3's unread fix, instead of promising a next round after every fix. The revealed-defect and independence clauses after the replaced sentences remain. |
| Records finding, exclusions, disputed classification, permission and exception sentences: straight quotes. | Baseline lines 331, 336, 339, 340, 344, 346 and 347: only apostrophe and quotation characters change. "And changes nothing else", exclusions including authority-changing text, the Operator's retained authority and end 3's exception all retain their words and effects. |
| Tolerance: straight apostrophe. | Baseline line 399: 20 percent, floor 60, cap 500 lines and 4 files stay unchanged; the exact rule still uses "20 %". |
| End 3: shared waiting-rule pointer. | Baseline lines 342-344 plus the Operator's decision: the end remains posted with findings and fixes and may be overruled before merge. The pointer locates the single shared wait below; it defines no second duration or window. |
| Delegated figure or end: author posts a new issue comment naming it. | Amendment 1, condition 1: both delegations require a fresh post; a local draft, an unnamed reference or an edit of an earlier comment supplies no post. |
| Ten minutes from the new comment's UTC posting time. | Operator decision and amendment 1: use that new comment's created_at; the posting time selects the window. A daytime post waits the full ten minutes even if the deadline crosses 19:00 UTC. |
| Day from 05:00 inclusive to 19:00 exclusive UTC; corresponding UTC+3 times. | Operator decision and confirmed boundary reading: 08:00-22:00 UTC+3 equals 05:00-19:00 UTC. Inclusion at 05:00 and exclusion at 19:00 leave each instant in exactly one window. |
| Night from 19:00 inclusive to 05:00 exclusive UTC; no wait. | Same decision: 22:00-08:00 UTC+3 equals 19:00-05:00 UTC. The merge may follow the required new post with no wait in that window. |
| An edit of an earlier comment posts nothing. | Amendment 1 removes the old-created_at loophole identified by the plan review. Editing a 10:00 comment at 14:00 cannot satisfy this delegation's posting step. |

For both a figure and an end, the amended sentence gives these same outcomes:

| New-comment UTC posting time | Earliest merge under this sentence |
| --- | --- |
| 04:59:59 | Immediately after the post |
| 05:00:00 | 05:10:00 |
| 18:59:59 | 19:09:59 |
| 19:00:00 | Immediately after the post |
| 14:00:00 | 14:10:00 |

An edit at 14:00 of a comment created at 10:00 supplies no post. A new comment at
14:00 starts the wait until 14:10. These are source comparisons, not a claim that
a delegation has been posted. This new rule takes effect when #451 lands and
supplies no authority for its own landing.

## Pointers and vocabulary

| Changed passage | Source and clause-by-clause comparison |
| --- | --- |
| R12 growth sentence | Baseline canonical budget, lines 398-419: size still requires the plan's bound; growth follows advance approval within tolerance or the canonical overrun disposition. At maximum 650/20, 700/20 is approved in advance; R12 no longer independently orders a new slice for those 50 lines. The warning that plan approval is not unlimited implementation approval remains. |
| Named-base pitfall | Baseline canonical lines 376-396: whole-branch unit, evidence and bookkeeping exclusions, execution inputs and mixed files are all reached by the budget pointer. The pitfall's historical measurements and procedure for a changed base remain; none of those figures is recalculated here. |
| Size class | Baseline canonical lines 398-400: maximum plus tolerance determines the limit; 651 lines is not automatically forbidden under a 650 maximum. Plan size classes and the target's model-selection size field retain their distinct meanings. |
| Issue split | Baseline end 4, lines 351-354: successor carries the stopped work, children take the open findings it does not take; the linked source sets the disposition. The independent "always" for an unapproved overrun is removed. The example chooses successor F1 and child F2, without implying every stopped branch must split. |
| Material | Baseline lines 294-327 and 355-362: the materiality test and note distinction remain. The outcome pointer covers clean verdicts and the other permitted endings instead of describing a clean verdict as the sole stop. |
| Frozen head | Baseline lines 268-307 and end 3: named review commit remains the definition. The pointer governs fixes, close-out, merges and exceptions; no unconditional next round is added. The specific example still illustrates an ordinary fix followed by review. |
| Cycle cap example | Baseline cycle mechanics: first round is Cycle 0; one fix-and-review cycle is Cycle 1. With cap 1, this is a concrete ordinary case, not a requirement to run another round after an end-3 fix. |
| Ceiling example | Baseline budget: 20 percent of 650 is 130 lines, plus 4 files. Figure 700/20 is within limit 780/24. |
| Records finding row | Baseline lines 331-340 and amendment 1: whose fix restores source or form AND changes nothing else. The link retains exclusions; the collision rejects classification by location alone. The traceability example changes a name to an existing test without changing that test or other behaviour. |
| Tolerance row | Baseline lines 398-407: size allowance above maximum, distinct from materiality. Example 650/20 gives 130/4, inside the unchanged floor and cap. |
| Input row | Baseline lines 383-389: execution data is distinct from evidence audited only for form, citations or evidence. Fixture is input, audit-only log is not, dual-use file remains input. |
| Margin row | Baseline line 401 and K5 output at docs/tests/freeze-check.sh:569: limit minus figure independently in lines and files. For 650/20 + 130/4 - 700/20, margin is 80/4. No new margin definition is introduced. |
| Coordinated Universal Time row | Amendment 1, condition 3: full term, UTC abbreviation, role in posting/waiting, collision with fixed offset UTC+3 and two concrete conversions. It defines vocabulary, not a second waiting rule. |

## Occurrence audit and preserved text

mentions.tsv identifies 1,699 distinct locations against the named baseline: all
913 primary matches, all 789 supplementary matches, deduplicated, and named
neighbouring/context lines. Its search column distinguishes those sources. The
exact searches and measured statuses are in test-runs.md. Classification follows
the surrounding passage's purpose, not the token alone. No existing accepted ADR,
task record or earlier run contract is rewritten; only accepted ADR bodies are
claimed immutable. Product budgets and target behaviour are distinct from the
budget for builders landing this branch.

The explicitly marked summaries at README.md:33-38,
docs/onboarding-for-engineers.md:178-180 and docs/glossary.md:166 retain their
canonical pointers and supply no independent authority requiring a next round.
The unmarked gate overview and operative/glossary restatements named in the plan
are replaced. AGENTS.md:59, the plan's estimate/cap pointers, the DoD checklist and
the K5 description already defer to the canonical source and remain compatible.

Condition 2 drops decision 8. The complete guardrails validation section is
byte-for-byte unchanged from the baseline, as checked by C2-validation-preserved.
Its clean-verdict test determines whether the result passes validation. A
permitted landing with a known limit or unread fix need not pass that test. No
conflict is asserted and no separate issue is needed for that cell.

The unchanged task check also protects count exclusions, input and dual-use
definitions, margin, end-3 eligibility and exclusions, and ends 1, 2 and 4.
N4's suggested edits to runs/T-3ff9/checks.sh are declined: that historical
executable records #447's contract and is preserved. The supplied plan review
reports it already failing on main because freeze-check changed after its base;
Stage B has not rerun it or represented it as green. Its note and the promise in
docs/tasks/T-3ff9.md:41 must be carried into this task's close-out record.

The reference-sweep lessons already in guardrails cover following paraphrases
and reading their neighbourhoods. This task adds no duplicate pitfall.

## Evidence scope and remaining stage work

The only declared evidence records are mentions.tsv, semantic.md, conditions.md
and test-runs.md under runs/T-ebm4/. Inspection of checks.sh shows it reads source
documents and baseline blobs, not these records as configuration or expected
output. Inspection of freeze-check.sh K6 shows an audit of the conditions table,
paths and line ranges at HEAD. A repository search found no other path-specific
consumer outside these records. Their use is audit evidence under ADR-0028.
checks.sh stays counted; no existing script, hook or CI file changes.

Stage B leaves the operative changes and records for the coordinator's commits.
The task file, completed-log line, committed-head pre-freeze run, independent
review and verdict/resource close-out remain for their assigned later stages.
The confirmed maximum is 400 lines over 10 files against current main at landing,
with ADR-0028 exclusions, and cap 1. No freeze or final landing figure is claimed
by this author comparison.
