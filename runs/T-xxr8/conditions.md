# Conditions and acceptance evidence for T-xxr8

Stages A and B address the sole numbered condition of the plan review supplied for
[#457](https://github.com/pharzam/layup/issues/457). Its four explicit edits
and preservation requirements are itemized below. The supplied host results
complete acceptance criterion 1 and the host scan part of criterion 2. The
required PR security job still needs CI evidence. These records do not claim
a freeze or completed review.

| Condition | Evidence |
| --- | --- |
| Condition 1: name each edit in plan step 3. | `runs/T-xxr8/test-runs.md`, Plan amendment 1, names each change and every historical exclusion. |
| Condition 1: setup states go 1.26.10 since T-xxr8; keep the catalog minor version 1.26. | `docs/spec/setup.md:1212-1214`; `go.mod:3`; `internal/catalog/go/kinds.tsv:2-6`. |
| Condition 1: the guardrail check's full patch example becomes go 1.26.10. | `docs/guardrails.md:396-398`. |
| Condition 1: retain the T-w89c advisory history in the guardrails. | `docs/guardrails.md:391-395`; `runs/T-xxr8/test-runs.md`, Plan amendment 1. |
| Condition 1: retain all existing matches under runs, the T-w89c task, and its completed-log entry. | `docs/tasks/T-w89c.md`; `docs/tasks/completed.md:83`; `runs/T-xxr8/test-runs.md`, Plan amendment 1 and Stage A final scope verification. |
| Acceptance criterion 1: build, vet and all three test levels pass on Go 1.26.10. | `runs/T-xxr8/test-runs.md`, Stage B host green: Go 1.26.10, build, vet, unit, integration and E2E all exit 0. The earlier sandbox socket failures remain recorded separately. |
| Acceptance criterion 2: govulncheck clean and the PR's required security job passes. | `runs/T-xxr8/test-runs.md`, Baseline security red and Stage B host green: govulncheck exits 0 with "No vulnerabilities found." on Go 1.26.10. Freeze and CI pending records that the required PR security job result and run URL remain outstanding before the freeze. |
