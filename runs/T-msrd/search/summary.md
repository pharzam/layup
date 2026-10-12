# T-msrd — the public-solution search of M2d: the result

Two blind searchers screened the same sources ([`sources.md`](sources.md)) and ran
their own web searches in the words of the problem:
[`searcher-A.md`](searcher-A.md), GPT-6 Sol on Devin, and
[`searcher-B.md`](searcher-B.md), Claude Sonnet 5.5 on Claude Code. Their brief
([`brief.md`](brief.md)) held the facts, the milestone row and the sources: no
name of a candidate, nothing of the architecture or the specification. This file
is the author's reading of the two files, written after both ended, and
committed before the decisions of S3.

**Both conclude: no public solution can be the base of `M2d`.** No candidate can
be a dependency of a controller in Go with the standard library only
(`F-0004#1`). A: "No checked candidate combines all of the following in one
public solution." B: "Tools usable as a dependency: none. … Only API usage,
patterns and formats transfer." What transfers is GitHub's own API, and the
known limits that its documentation states.

## What the search gives each part of M2d

| Part | Candidates (searcher) | What they bear on in `M2d` (the decisions come in S3) |
| ---- | --------------------- | ----------------------------------------------------- |
| The setup from a pinned baseline, each value with its source | git-scaffold (A); Copier (A, B); ChimeraTK/project-template (A) | A lock of the template's full commit and a file of the answers, both in the project's Git (git-scaffold, Copier); git-scaffold refuses a missing required argument rather than invent one, but takes a template's default with no source. No candidate records a source for each value or keeps a value with no source as an open gap (A, B) |
| A second setup run on an earlier setup | git-scaffold `init --existing` and its update (A); Copier's update by a three-way merge (A, B); ChimeraTK (A) | Each needs its own earlier record and may end in conflicts for a person (Copier: "do not change the answers file by hand"); none reads a record of another tool or version, or checks the gates after the update (A, B) |
| The stop at a human decision, the resume from a record; an agent's text checked by code | Agentplane, Squelette (A, B); gh-aw's validated "safe outputs" (B) | A transition that code checks, and a status rebuilt from the records in the repository; none keeps that record on a branch that one identity alone may update (A, B) |
| The ruleset file that a program writes and a human applies | the GitHub rules REST API (A, B); gh-migrate-rulesets (A, B); rulesets-as-code-action, repo-ruleset-bot (A); github-settings-as-code, the Terraform GitHub provider's `github_repository_ruleset` (B) | The body of a ruleset: `enforcement`, `conditions`, `bypass_actors` (`actor_type` `Integration`, …, `bypass_mode` `always`, `pull_request`, `exempt`), and `required_status_checks[]` with an optional `integration_id`. Each tool applies with an admin token, which the controller does not hold (A, B); creating a ruleset needs Administration (write) (B, from a summary, not confirmed) |
| The read-back of the effective rules | the rules API, `GET /repos/{owner}/{repo}/rules/branches/{branch}` (A, B) | It returns the active rules of a branch, also of one that does not exist yet, inherited rules included; a ruleset in `evaluate` or `disabled` is not returned; `integration_id` is optional, so the read-back must compare it; the bypass list is read only from the ruleset itself (A, B). No tool compares the read-back with the file that the program wrote (B) |
| The two rulesets: no bypass on the gates, the App alone on the records branch | the rules API (A) | The rule `update` lets only the bypass actors update a matching ref, so the records branch takes the App as its bypass actor in its own ruleset, apart from the ruleset of the gates, which has none; no public example sets both (A) |
| The probe push refused for a rule violation | the rules API, the rule-suites API (A) | A push to a ref that does not exist is a **creation**, which a rule of updates or of pull requests may not cover: the ruleset that covers `layup-probe` needs the rule `creation` (A). The rule suites (`GET /repos/{owner}/{repo}/rulesets/rule-suites`) record a `fail` evaluation with its ref, actor and SHAs, a check beside the push's own result (A). Neither searcher found an official page of the refusal text `GH013` (B) |
| The audit of who updated each branch | the repository activity API, `GET /repos/{owner}/{repo}/activity` (A, B) | Its items give the actor, the ref, `before`, `after`, the time and the type (`push`, `force_push`, `branch_creation`, `branch_deletion`, `pr_merge`, `merge_queue_merge`); it lists accepted updates, not refused ones; its permission and retention are not stated on the page (A, B) |
| Continuous checks of the settings | Allstar (A, B); OpenSSF Scorecard, legitify (B) | Read-only reports or repairs after a change; a repair is not a prevention, and Allstar's hosted App is retired (A, B) |
| The repository passes its gates with the controller absent | — | Nothing found (A, B) |

## The limits of the search

- Searcher A ran on GPT-6 Sol, a model of the reasoning tier, as Devin's paid
  models answered again; searcher B on Claude Sonnet 5.5, as OpenCode's login
  failed, a model of neither tier list.
- Both runs were short: A 3 min 20 s, B 2 min 49 s. B says it read the five lists
  in full ("Five lists, read in full"); A gives a section per list. A file's
  claim of a full read is the searcher's own.
- B opened no README of a candidate (its Limits); A read the candidates' own
  pages. The permission that creates a ruleset (Administration, write) and the
  text of a refusal by a rule are not confirmed by an official page; the build
  task reads GitHub's documentation at its date, as `forge.md` asks.
