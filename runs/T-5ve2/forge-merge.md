# T-5ve2 forge merge experiment

Operator decision Q4 on #461. Run by the coordinator (Claude Opus 5.5) at about 2026-10-21T13:10Z. Question: do GitHub's own PR merge and "Update branch" honour `merge=union`?

- **Repository:** `pharzam/pharzam-layup-union-probe`, created by the Operator, disposable, public (copies of three public layup files and two probe entries), the App installed. The Operator deletes it.
- **Base:** this branch's `.gitattributes`, `docs/tasks/completed.md` and `docs/tests/traceability.md`, as commit 3c051ca.
- **Branches:** from the base, probe-a added an entry for T-aaaa directly under `## Log` and a row directly under the traceability separator; probe-b did the same for T-bbbb.

| Step | Result |
| --- | --- |
| PR #1 (A) merged with `PUT /pulls/1/merge`, merge_method `merge` | d9b611a |
| PR #2 (B) read after that | `mergeable: false`, `mergeable_state: dirty` |
| `PUT /pulls/2/update-branch` on B | HTTP 422, "merge conflict between base and head" |
| Local `git merge main` on B | exit 0; `git diff --check` clean; both files hold both entries, B's above A's (landing order, most recent first) |
| B pushed, PR #2 merged by GitHub | 2ab1afb |

The branches were deleted afterwards.

**Conclusion.** GitHub's own merge and "Update branch" do not honour `merge=union`. Merge main locally (`git merge origin/main`), never with "Update branch" or a web conflict editor. The repository's strict protection already requires a branch to take main before landing.
