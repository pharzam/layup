# T-92lm — the test runs

Task `T-92lm` (#326), base `b7f3f7d`. Times are UTC on 2026-10-14, read from the clock after each step. The code of round 1 is `b3f78e7`, that of round 2 `14b7333`, the final code `52d11f5`.

## Step 3: the tests, red first

The plan review (comment 6175863181) and the answer (6175893691, 18:53:25Z) came first. **Declarations that hold no fix**, so that a test fails by an assertion and not by a build error:
in `internal/verify` the type `AdaptedPhrase`, `Words` and the three functions returning nothing (18:54:18Z), and in `internal/run` the first parameter `name` of `proseText`, which its
old body ignores, with the closure that passes it. For round 1 (comment 6176783005, `material`, four findings) the field `Tail` left empty, `Stem` and `AdaptedExempt` returning nothing.

| Group | Run | Result, read for its cause |
| ----- | --- | -------------------------- |
| `internal/verify`, 4 new tests | 18:54:24Z | 4 FAIL by an assertion on missing data: "0 phrases for 16 patterns", "0 paths for 9 files", "a list is empty", "no line". |
| `internal/run`, 7 new prompt tests | 18:56:58Z | 6 FAIL, "the prompt does not say" 39 sentences of the rules; the 5 older tests of the prose session PASS (retained). |
| `TestTheStepSetupRunsS01ToS14`, integration | 18:57:58Z | FAIL (50.5 s): the prompt that the scripted session read has no name of the target, no branch, no phrase, no request to search. |
| round 1: `internal/verify` | 19:42:04Z | 3 FAIL: "the tail \"\" is none of word, open and stem" for the 16 patterns, "no name", "a list is empty". The probes of what stands before a phrase PASS at red: they hold already (guards). |
| round 1: `internal/run`, unit | 19:43:17Z | 4 FAIL (words, identity, onboarding and guardrails, the search). |
| round 1: the search command, integration | 19:43:22Z | FAIL: the old command does not find "fillable". |

**Guards that passed at red:** `TestTheProsePromptHoldsNoAngleQuote`. **Changed after a red run:** in round 1's code two sentences of the prompt and the test of the allowed files (mutants 22, 23 and 7
of that set show them red or alive), and in round 2 the negative half of the test of the search command (mutant 38 survived the first text: `grep -o` prints only the match, so the test now reads
a file of whole words inside longer ones).

## Mutants (step 3, each alone on a copy of the tree, a build failure counted as not run)

Round 1's set of 29 (26 caught, 2 survived and 1 not run on the first run; the weak test and two mutants repaired) was caught in full on `b3f78e7`. The set of round 2 has 40 mutants
on `14b7333`; a first run stopped on a mutant whose old text wrapped, the run after it left mutant 38 alive, and with the test strengthened all 40 are caught.

| Caught by | Mutants |
| --------- | ------- |
| TestEachPatternOfAdaptedHasAPhraseThatNamesIt | pattern-added-with-an-empty-phrase, pattern-added-with-the-phrase-of-another, a-form-that-is-no-hit, a-phrase-with-an-angle-quote, a-tail-that-the-check-does-not-have, a-stem-that-is-not-the-stem, the-last-phrase-left-out, the-case-of-a-phrase-lost, the-tail-of-a-phrase-lost, the-forms-of-a-phrase-not-read |
| TestTheFilesThatMayNameTheBaselineAreTheOnesCheckAdaptedHonors | an-allowed-file-that-the-check-does-not-honor, an-allowed-file-left-out |
| TestThePassedNamesOfCheckAdaptedAreNoHit; TestTheHitsOfAdapted, TestTheTextOfAdapted | a-name-passed-that-the-check-does-not-pass; a-passed-name-left-out-of-the-check |
| TestTheListsOfThePromptAreCopies | the-passed-names-are-not-a-copy, the-lines-of-the-baseline-are-not-a-copy |
| TestTheProsePromptHasTheNameOfItsTarget; TestTheStepSetupRunsS01ToS14 | the-name-of-the-target-ignored; the-name-not-passed-to-the-text |
| TestTheProsePromptGivesTheWordsOfCheckAdapted | a-phrase-left-out-of-the-prompt, the-open-group-merged-into-the-whole-group, the-capital-word-read-as-a-phrase-of-the-list, the-allowed-files-left-out, the-passed-names-left-out, a-letter-that-is-not-ASCII-as-in-the-first-draft, the-first-sentence-of-rule-1-as-it-was-drafted |
| TestTheProsePromptGivesTheRulesOfTheOtherChecks | the-lines-of-the-baseline-left-out, the-branch-left-out, the-keeping-of-the-baseline-rules-left-out, the-form-of-ci-left-out, a-whole-line-reading-of-the-identity-lines, the-end-of-the-line-of-a-check-value-left-out, the-ban-of-the-onboarding-file-left-out, rule-6-left-out |
| TestTheProsePromptAsksForASearchBeforeTheEnd (the first two also HoldsNoAngleQuote and the words test) | the-request-to-search-left-out, an-angle-quote-in-the-text, the-end-left-out |
| TestTheProsePromptLetsAnAdaptedFileKeepItsMarkers | the-marker-sentence-left-out |
| TestTheSearchCommandOfTheProsePromptFindsWhatCheckAdaptedFinds (integration) | the-aid-without-the-end-of-a-word, the-aid-without-the-letters-of-a-stem, the-aid-without-the-join-of-lines |

## The gate on `52d11f5`, from 20:08:22Z (each exit code)

go build ./... 0; go vet ./... 0; gofmt -l internal cmd (empty) 0; go test ./... 0; go test -tags=integration ./... 0; go test -tags=e2e ./... 0; adr-lint 0; prd-lint 0; link-lint 0; run-discipline-tests 0; nested-checkout-check 0; setup-check 0; git diff --check origin/main HEAD 0.

## The rules of the prompt against their sources (note 4 of the plan review; round 1 compared them again from the source)

| Clause | Source | Comparison |
| ------ | ------ | ---------- |
| The phrases, by tail | `adHits` and `tail` (`adapted.go:48`); `AdaptedPhrases` | Each form is a hit of its pattern and of no other; before a phrase an ASCII letter, digit or underscore makes another word and a letter that is not ASCII does not; after it, by tail: whole word, open end (no right boundary), stem (the beginning shared by the forms, any ASCII letters after it, a marker before the next full stop). All probed by a test. The check also lets `#` or `-` before `fill`, `replace` and `fill in`, and `-` after `the template`, pass: the prompt is stricter and does not say so. |
| Case, line ends, blanks, paragraph | `hitsOf`, `paragraphHits` | ASCII A to Z read as small; lines of a paragraph joined by a space, blanks one space, a blank line ends the paragraph: stated. |
| A marker is one unit | `unitText` | A marker runs to its first close quote, also over lines; stated. |
| The names that pass | `kitExempt`, `passes`, `AdaptedExempt` | A hit of `kit` that starts one of the two names is dropped, whatever follows: stated, and a test. |
| The capital word | `adAllowed`, rule 2, `AdaptedAllowed` | Read as it is, outside the nine files: stated. |
| A citation | `citation`, `resolveFact`, three consumers | `F-`, four digits, `#`, digits; one record file and a line of the number: stated, with its files. Leading zeros on the fact line are accepted by the check: the prompt is stricter. |
| `README.md`, `AGENTS.md` | `identityTarget`, `identityFindings` | The name and `layup-records` are substrings of `README.md`; the link; the three texts are substrings anywhere in either file: stated "anywhere, also inside a longer line". |
| The onboarding file | `onboardingFindings` | Any U+2039 is a finding; the prompt bans both quotes, which is stricter; the link is exact. |
| The guardrails file | `guardrailsEntry`, `guardrailsFindings` | An entry runs to the next entry or heading; the value is the text after the last `Check: ` of the first such line, to the end of that line; three forms, an existing file, a job: stated. |
| Markers | `LostMarkers`, `CheckAsked` | An adapted file loses no place; a new marker makes S10 stop on an `M-` row that the batch did not ask: stricter than `LostMarkers`, which allows more places. |
| The search command | `grep -o -i -E`, `\\<`, `\\>` | Joins every line, also over a blank line; takes a letter that is not ASCII as part of a word; shows a word inside a marker and the passed names; not the baseline's name: all stated; a test runs it as printed. |
