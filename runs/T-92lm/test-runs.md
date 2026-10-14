# T-92lm — the test runs

Task `T-92lm` (#326), base `b7f3f7d`; the head of the code and its documents is `b3f78e7`. All times are UTC on 2026-10-14, read from the clock after each step.

## Step 3: the tests, red first

The plan review (4 min 36 s, comment 6175863181) and the author's answer (comment 6175893691, 18:53:25Z) came first. **Declarations that hold no fix**, so that each
test fails by an assertion and not by a build error: in `internal/verify`, the type `AdaptedPhrase` with `Words` returning nothing and the three functions
`AdaptedPhrases`, `AdaptedAllowed` and `KitPhrases` returning nothing (18:54:18Z); in `internal/run`, the first parameter `name` of `proseText`, which its old body
ignores, and the closure of `proseSession` that passes the name (the text is unchanged).

| Group | Run (UTC) | Result, read for its cause |
| ----- | --------- | -------------------------- |
| `internal/verify`, 4 new tests | 18:54:24Z | 4 FAIL, each by an assertion on missing data: "0 phrases for 16 patterns", "0 paths for 9 files of the table", "a list is empty: 0 phrases, 0 files, 0 lines of the baseline", "no line". |
| `internal/run`, 7 new prompt tests, unit | 18:56:58Z | 6 FAIL: "the prompt does not say" each sentence of the rules (39 failed assertions in all), the prompt printed in full; the 5 older tests of the prose session PASS (retained). |
| `TestTheStepSetupRunsS01ToS14` (the world), integration | 18:57:58Z | FAIL (50.5 s): the prompt that the scripted session read does not say "the name of this repository, acme/target, as text", "must name the branch layup-records", the phrase "kit", "kits", nor "Before you end, search each file that you wrote"; its two subtests PASS. |

**A guard that passed at red:** `TestTheProsePromptHoldsNoAngleQuote` (the old prompt has none; the new text must keep it so). **Changed after the red run:** two sentences of the
prompt were tightened while the code was written ("No file may hold any of these phrases" for "A file holds none of them", and the heading of rule 2), shown red by
mutants 22 and 23 below; and `TestTheFilesThatMayNameTheBaselineAreTheOnesCheckAdaptedHonors` read a file that the check did not read (mutant 7 survived the first run), so it now puts a word
of rule 1 in each allowed file too.

## Mutants (step 3, each alone on a copy of the tree, a build failure counted as not run)

A first run, on the code before the last two edits (rule 2 reflowed, the closure named `target`), gave 26 caught, 2 survived and 1 not run: mutant 4 held a Go escape in a raw string (no
angle quote at all), mutant 7 is the weak test above, and mutant 9 dropped the use of `slices`. Two mutants were rewritten and the test strengthened. The run on the final code (19:23:50Z)
gave 28 caught and 1 not run (mutant 12 left `target` unused); rewritten, it is caught (19:24:51Z). **29 caught, 0 survived, 0 not run.**

| # | Mutant | Caught by |
| - | ------ | --------- |
| 1 | pattern-added-with-an-empty-phrase | TestEachPatternOfAdaptedHasAPhraseThatNamesIt |
| 2 | pattern-added-with-the-phrase-of-another | TestEachPatternOfAdaptedHasAPhraseThatNamesIt |
| 3 | a-form-that-is-no-hit | TestEachPatternOfAdaptedHasAPhraseThatNamesIt |
| 4 | a-phrase-with-an-angle-quote | TestEachPatternOfAdaptedHasAPhraseThatNamesIt |
| 5 | the-last-phrase-left-out | TestEachPatternOfAdaptedHasAPhraseThatNamesIt |
| 6 | the-case-of-a-phrase-lost | TestEachPatternOfAdaptedHasAPhraseThatNamesIt |
| 7 | an-allowed-file-that-the-check-does-not-honor | TestTheFilesThatMayNameTheBaselineAreTheOnesCheckAdaptedHonors |
| 8 | an-allowed-file-left-out | TestTheFilesThatMayNameTheBaselineAreTheOnesCheckAdaptedHonors |
| 9 | the-lines-of-the-baseline-are-not-a-copy | TestTheListsOfThePromptAreCopies, TestIdentityFindings |
| 10 | the-forms-of-a-phrase-not-read | TestEachPatternOfAdaptedHasAPhraseThatNamesIt |
| 11 | the-name-of-the-target-ignored | TestTheProsePromptHasTheNameOfItsTarget |
| 12 | the-name-not-passed-to-the-text | TestTheStepSetupRunsS01ToS14 |
| 13 | a-phrase-left-out-of-the-prompt | TestTheProsePromptGivesTheWordsOfCheckAdapted |
| 14 | the-allowed-files-left-out | TestTheProsePromptGivesTheWordsOfCheckAdapted |
| 15 | the-lines-of-the-baseline-left-out | TestTheProsePromptGivesTheRulesOfCheckIdentity |
| 16 | the-branch-left-out | TestTheProsePromptGivesTheRulesOfCheckIdentity |
| 17 | the-request-to-search-left-out | TestTheProsePromptAsksForASearchBeforeTheEnd |
| 18 | the-keeping-of-the-baseline-rules-left-out | TestTheProsePromptGivesTheRuleOfACitation |
| 19 | the-form-of-ci-left-out | TestTheProsePromptGivesTheRulesOfCheckOnboardingAndGuardrails |
| 20 | the-ban-of-the-onboarding-file-left-out | TestTheProsePromptGivesTheRulesOfCheckOnboardingAndGuardrails |
| 21 | an-angle-quote-in-the-text | TestTheProsePromptAsksForASearchBeforeTheEnd, TestTheProsePromptHoldsNoAngleQuote |
| 22 | the-first-sentence-of-rule-1-as-it-was-drafted | TestTheProsePromptGivesTheWordsOfCheckAdapted |
| 23 | the-first-sentence-of-rule-2-as-it-was-drafted | TestTheProsePromptGivesTheRuleOfACitation |
| 24 | the-aid-without-whole-words | TestTheProsePromptAsksForASearchBeforeTheEnd |
| 25 | the-aid-without-the-join-of-lines | TestTheProsePromptAsksForASearchBeforeTheEnd |
| 26 | the-aid-with-the-capital-word | TestTheProsePromptAsksForASearchBeforeTheEnd |
| 27 | the-marker-sentence-left-out | TestTheProsePromptLetsAnAdaptedFileKeepItsMarkers |
| 28 | rule-6-left-out | TestTheProsePromptGivesTheRulesOfCheckOnboardingAndGuardrails |
| 29 | the-end-left-out | TestTheProsePromptAsksForASearchBeforeTheEnd |

## The gate on `b3f78e7`, 19:07:41Z to 19:22:24Z

| Check | Exit | At |
| ----- | ---- | -- |
| go build ./... | 0 | 19:07:41Z |
| go vet ./... | 0 | 19:07:42Z |
| gofmt -l internal cmd (empty) | 0 | 19:07:42Z |
| go test ./... | 0 | 19:07:48Z |
| go test -tags=integration ./... | 0 | 19:19:29Z |
| go test -tags=e2e ./... | 0 | 19:21:32Z |
| adr-lint | 0 | 19:21:33Z |
| prd-lint | 0 | 19:21:33Z |
| link-lint | 0 | 19:21:48Z |
| run-discipline-tests | 0 | 19:21:55Z |
| nested-checkout-check | 0 | 19:21:56Z |
| setup-check | 0 | 19:22:24Z |
| git diff --check origin/main HEAD | 0 | 19:22:24Z |

## The rules of the prompt against their sources (note 4 of the plan review)

Each clause of the prompt (`proseChecksText`, `internal/run/prose.go:119`) with the code that is its home. A difference is named where the prompt is stricter than the check.

| Clause | Source | Variants and exceptions |
| ------ | ------ | ----------------------- |
| 1. The phrases, each with its forms | `adHits` (`adapted.go:48`), `AdaptedPhrases` (`:298`) | Each quoted form is a hit of its pattern and of no other (a test). Patterns with no word end (`your project`, `skip this section`, `delete the one you do not use`) list the plural; `the template` is not a hit before `-` and `fill`, `replace` and `fill in` are none after `#` or `-`: the prompt is stricter and says nothing of it. |
| 1. Case, line ends, blanks | `hitsOf` (`:148`) | The text is read in lower case but for rule 2; the lines of a paragraph join with one space and blanks are one space; **a blank line ends a paragraph**, so a phrase across one is no hit, while the aid's `tr -s` joins it: the aid is stricter, and the prompt calls it an aid. |
| 1. A letter, a digit or an underscore before a phrase | the boundary `[^a-z0-9_]` of each pattern | After a single word the same class must follow (`kits?`, `adopters?`, `optional`, `adapts?`, `you use`): `kits` is a hit, `kitten` is not. |
| 1. A marker is one unit | `unitText` (`:182`) | A marker runs to the first close quote, also over lines; an open quote with no close is no marker. The prompt says "a marker is read as one unit". |
| 1. The capital word and its files | `adAllowed` (`:26`), `AdaptedAllowed` (`:307`), rule 2 | Read as it is; the nine files are the table's; every other file is read. |
| 1. Every file | `adaptedFiles` (`:107`), `adExclude` (`:20`) | The check reads each tracked `.md` file that `adExclude` does not match; the session writes only files of the stop, none under `docs/facts/`. |
| 2. A citation | `citation` (`target.go:210`), `resolveFact` (`:226`) | `F-` and four digits, `#` and digits; one file `docs/facts/F-NNNN-*.md` and a line `^0*N\. `: the number, a full stop and a space. Read by `onboardingFindings` (`:252`), `glossaryFindings` (`:269`) and `guardrailsFindings` (`:283`) only. |
| 2. Keep the rules of the baseline | Inv-7 of `docs/guardrails.md` (`:68`): the domain changes content, never rules | No code reads it; a text assertion, and round 1. |
| 3. `README.md` | `identityTarget` (`checks.go:236`), `identityFindings` (`:131`) | The name is a substring of `README.md`, as is `layup-records`; the link is `](docs/setup/armature.pin)`; `KitPhrases` (`:127`) are searched in `README.md` and `AGENTS.md`. |
| 4. The onboarding file | `onboardingFindings` (`target.go:252`) | Any U+2039 character is a finding, a marker or not; the link `](facts/problem-statement-brief.md)`; the file must exist. |
| 5. The guardrails file | `guardrailsEntry` (`:277`), `guardrailsFindings` (`:283`) | The entry runs to the next entry or a heading; its value is after the last `Check: ` of its first such line; `no check yet`, or `<path> (hook)`, or `<path> (ci:<job>)` with a regular file and a non-empty job; no entry passes. |
| 6. Markers | `LostMarkers` (`adapted.go:252`, S14), `CheckAsked` (`internal/setup/setup.go:526`, S10) | An adapted file loses no place of a marker; a marker that the pinned tree lacks makes S10 stop on an `M-` row that the batch did not ask: `fail`. |
| The aid | `unitText` | `-i` for the case, `-w` for the boundary (the class `[^A-Za-z0-9_]`, as the patterns), `-F` for the text, `tr -s ' \t\n' ' '` for the line ends; it cannot read the two phrases about a marker, and it is not the check. |
