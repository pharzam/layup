# T-92lm — the test runs

Task `T-92lm` (#326), base `b7f3f7d`. Times are UTC, read from the clock after each step: 2026-10-14 up to the record of round 2, 2026-10-15 after it (the rehearsal fixes from 03:39Z). The code that round 1 read is `b3f78e7`, that round 2 read `4fbfc06`, the final code `35c03f8`.

## Step 3: the tests, red first

The plan review (comment 6175863181) and the answer (6175893691, 18:53:25Z) came first. **Declarations that hold no fix**, so that a test fails by an assertion and not by a build error: in `internal/verify` the
type `AdaptedPhrase`, `Words` and the three functions returning nothing (18:54:18Z), and in `internal/run` the first parameter `name` of `proseText`, which its old body ignores, with the closure that passes it;
for round 1 (comment 6176783005) the field `Tail` left empty, `Stem` and `AdaptedExempt` returning nothing.

| Group | Run | Result, read for its cause |
| ----- | --- | -------------------------- |
| `internal/verify`, 4 new tests | 18:54:24Z | 4 FAIL by an assertion on missing data: "0 phrases for 16 patterns", "0 paths for 9 files", "a list is empty", "no line". |
| `internal/run`, 7 new prompt tests | 18:56:58Z | 6 FAIL, "the prompt does not say" 39 sentences of the rules; the 5 older tests of the prose session PASS (retained). |
| `TestTheStepSetupRunsS01ToS14`, integration | 18:57:58Z | FAIL (50.5 s): the prompt that the scripted session read has no name of the target, no branch, no phrase, no request to search. |
| round 1: `internal/verify` | 19:42:04Z | 3 FAIL: "the tail \"\" is none of word, open and stem" for the 16 patterns, "no name", "a list is empty". The probes of what stands before a phrase PASS at red: they hold already (guards). |
| round 1: `internal/run`, unit | 19:43:17Z | 4 FAIL (words, identity, onboarding and guardrails, the search). |
| round 1: the search command, integration | 19:43:22Z | FAIL: the old command does not find "fillable". |
| round 2: `internal/run`, unit | 20:39Z | 2 FAIL (the words, the search) by the new sentences; the probes of the six contexts of a stem in `internal/verify` PASS at red (guards: they describe the check). |
| rehearsal fixes: `internal/run`, unit | 03:39Z | 3 FAIL (the words, the other checks, the search) by the three new sentences (the stem as a start of a word, the stop by the period, the path of a `Check:` value); in `internal/verify` the two new contexts of a stem and the path case PASS at red (guards: they describe the check). |
| rehearsal 2: `internal/run`, unit | 04:11Z | FAIL (the other checks) by the sentence on the job name; the test of a `Check:` value in `internal/verify` PASSES at red (a guard: it describes the check). |
| rehearsal 3: `internal/run`, unit | 04:40Z | FAIL (the other checks) by the two sentences of rule 5 (the next line that starts with `#`, a comment in a code block too; a regular file, not a directory); the test of the entry that a `#` line ends, in `internal/verify`, PASSES at red (a guard: it describes the check). |
| fourth look: `internal/run`, unit | 05:05Z | FAIL (the other checks) by the sentence on the form of the path (from the root, no `.` or `..` element, no `//`, no `/` at its start or end); the test of the six forms in `internal/verify` PASSES at red (a guard: it describes the check, as `fs.ValidPath` reads it). |

**Guards that passed at red:** `TestTheProsePromptHoldsNoAngleQuote`. **Changed after a red run:** in round 1's code two sentences of the prompt and the test of the allowed files (mutants 22, 23 and 7
of that set show them red or alive), and in round 2 the negative half of the test of the search command (mutant 38 survived the first text: `grep -o` prints only the match, so the test now reads
a file of whole words inside longer ones).

**Three private rehearsals.** While the Operator's answer at the cap was awaited, a read of `c4fc063` by GPT-6.1 Sol (03:28:14Z to 03:37:54Z, 242,830 tokens; no record posted, as no round was authorised) found five findings and a note: the stem instruction and the stopping rule (a digit after the stem; "sentence" for "period"), a `Check:` path with a space and a parenthesis in it, the output limits of the process in the protocol, the audit record of attempt 1, the budget measured after a merge of `main`, and two stale descriptions. They are settled in `4657919` and `6835776`; the merge of `main` was dropped from the branch, so the budget is measured against its base `b7f3f7d`. A second read of `6835776` (04:02:26Z to 04:10:38Z, 238,187 tokens) found one finding, the job name of a `Check:` value, and the dates of this record: settled in `7824855`. A third read of `2e86501` (04:31:04Z to 04:39:29Z, 250,509 tokens) found one finding, rule 5, which said that an entry runs to the next entry or heading where the check ends it at any raw line that starts with `#`, a comment inside a code block too: settled in `ab3420d`, with the sentence on a regular file (a directory is no file for the check). The author ended the reads there, as the findings had decayed (five, one, one); the next read is a round, which the Operator authorises. A fourth look, the author's own and not a read by a reviewer, compared the path of a `Check:` value with the check once more and found that the check reads it as `fs.ValidPath` does, so `./docs/x.sh`, `docs//x.sh`, `docs/../x.sh`, `/docs/x.sh` and `docs/x.sh/` are no file for it where the prompt said only "the path of a regular file": settled in `35c03f8`.

## Mutants (step 3, each alone on a copy of the tree, a build failure counted as not run)

Round 1's set of 29 (26 caught, 2 survived and 1 not run on the first run; the weak test and two mutants repaired) was caught in full on `b3f78e7`. The set of round 2 had 40 mutants
on `4fbfc06` (one survived, mutant 38 below, and the test was strengthened); the final set has 51 on `35c03f8`, each caught.

| Caught by | Mutants |
| --------- | ------- |
| TestEachPatternOfAdaptedHasAPhraseThatNamesIt | pattern-added-with-an-empty-phrase, pattern-added-with-the-phrase-of-another, a-form-that-is-no-hit, a-phrase-with-an-angle-quote, a-tail-that-the-check-does-not-have, a-stem-that-is-not-the-stem, the-last-phrase-left-out, the-case-of-a-phrase-lost, the-tail-of-a-phrase-lost, the-forms-of-a-phrase-not-read |
| TestTheFilesThatMayNameTheBaselineAreTheOnesCheckAdaptedHonors | an-allowed-file-that-the-check-does-not-honor, an-allowed-file-left-out |
| TestThePassedNamesOfCheckAdaptedAreNoHit; TestTheHitsOfAdapted, TestTheTextOfAdapted | a-name-passed-that-the-check-does-not-pass; a-passed-name-left-out-of-the-check |
| TestTheListsOfThePromptAreCopies | the-passed-names-are-not-a-copy, the-lines-of-the-baseline-are-not-a-copy |
| TestTheProsePromptHasTheNameOfItsTarget; TestTheStepSetupRunsS01ToS14 | the-name-of-the-target-ignored; the-name-not-passed-to-the-text |
| TestTheProsePromptGivesTheWordsOfCheckAdapted | a-phrase-left-out-of-the-prompt, the-open-group-merged-into-the-whole-group, the-capital-word-read-as-a-phrase-of-the-list, the-allowed-files-left-out, the-passed-names-left-out, a-letter-that-is-not-ASCII-as-in-the-first-draft, the-first-sentence-of-rule-1-as-it-was-drafted, the-stem-context-as-a-marker-only, blanks-as-in-the-first-drafts, the-passed-names-without-the-words-that-start-with-them, the-stem-as-a-word-of-ASCII-letters |
| TestTheProsePromptGivesTheRulesOfTheOtherChecks | the-lines-of-the-baseline-left-out, the-branch-left-out, the-keeping-of-the-baseline-rules-left-out, the-form-of-ci-left-out, a-whole-line-reading-of-the-identity-lines, the-end-of-the-line-of-a-check-value-left-out, the-ban-of-the-onboarding-file-left-out, rule-6-left-out, the-path-restriction-left-out, the-form-of-the-path-left-out, the-raw-heading-boundary-as-a-heading, a-file-that-exists-for-a-regular-file, the-job-name-left-out-of-the-restriction |
| TestTheProsePromptAsksForASearchBeforeTheEnd (the second also HoldsNoAngleQuote and the words test) | the-request-to-search-left-out, an-angle-quote-in-the-text, the-end-left-out, the-stopping-rule-of-the-second-draft, the-stopping-rule-in-its-sentence |
| TestTheProsePromptLetsAnAdaptedFileKeepItsMarkers | the-marker-sentence-left-out |
| TestTheSearchCommandOfTheProsePromptFindsWhatCheckAdaptedFinds (integration) | the-aid-without-the-end-of-a-word, the-aid-without-the-letters-of-a-stem, the-aid-without-the-join-of-lines |

## The gate on `35c03f8`, from 05:06:33Z (each exit code)

go build ./... 0; go vet ./... 0; gofmt -l internal cmd (empty) 0; go test ./... 0; go test -tags=integration ./... 0; go test -tags=e2e ./... 0; adr-lint 0; prd-lint 0; link-lint 0; run-discipline-tests 0; nested-checkout-check 0; setup-check 0; git diff --check origin/main HEAD 0.

## The rules of the prompt against their sources (note 4 of the plan review; rounds 1 to 3 compared them again from the source)

| Clause | Source | Comparison |
| ------ | ------ | ---------- |
| The phrases, by tail | `adHits` and `tail` (`adapted.go:48`) | Each form is a hit of its pattern and of no other; before a phrase an ASCII letter, digit or underscore makes another word, a letter that is not ASCII does not; after it, by tail (whole word, open end, stem): all probed by a test. `#` or `-` before `fill`, `replace` and `fill in`, and `-` after `the template`, pass in the check: the prompt is stricter. |
| Case, lines, blanks, paragraph, marker | `hitsOf`, `paragraphHits`, `unitText` | ASCII A to Z read as small; a paragraph's lines joined, spaces and tabs one space, a blank line ends the paragraph; a marker is one unit; a stem is the letters at the start of a word, then any text with no period, then an angle quote U+2039 (eight contexts probed): stated. |
| The passed names, the capital word | `kitExempt`, `passes`; `adAllowed`, rule 2 | A hit of `kit` that starts one of the two names is dropped, whatever follows; the capital word is read as it is outside the nine files: stated. |
| A citation | `citation`, `resolveFact`, three consumers | `F-`, four digits, `#`, digits; one record file and a line of the number: stated, with its files. Leading zeros on the fact line pass in the check: stricter. |
| `README.md`, `AGENTS.md`, the onboarding file | `identityTarget`, `identityFindings`, `onboardingFindings` | The name and `layup-records` are substrings of `README.md`; the link; the three texts anywhere in either file: stated. Any U+2039 in the onboarding file is a finding; the prompt bans both quotes: stricter. |
| The guardrails file | `guardrailsEntry`, `guardrailsFindings` | An entry runs to the next entry or to the next raw line that starts with `#` (a heading, a comment in a code block too); the value is the text after the last `Check: ` of the first such line, to its end; three forms, a regular file written from the root as `fs.ValidPath` reads it, a job, and no space and parenthesis in the path or in the job name (the check splits at the first and at the last): stated. |
| Markers | `LostMarkers`, `CheckAsked` | An adapted file loses no place; a new marker makes S10 stop on an `M-` row that the batch did not ask: stricter than `LostMarkers`. |
| The search command | `grep -o -i -E`, `\<`, `\>` | Joins every line, also over a blank line; takes a letter that is not ASCII as part of a word; shows a word inside a marker and the passed names; not the capital word: all stated; a test runs it as printed. |
