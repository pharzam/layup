# The test runs of T-zck8

## Red, at base `9e980b9` (2026-10-07T09:04Z)

`sh runs/T-zck8/sections.sh`: each of the 18 headings is missing.

```
FAIL docs/spec/run.md: ## The command (the command `layup run --new OWNER/NAME` and the restart `layup run TARGET`: argu
FAIL docs/spec/run.md: ### The steps of `layup run --new` (Start 1 to 3: the baseline, the root push, the read-back, the
FAIL docs/spec/run.md: ### The restart (the restart: rebuild the clones, take the lease, stop on another LAYUP version)
FAIL docs/spec/run.md: ### The lease and fencing (the lease with heartbeat and fencing)
FAIL docs/spec/run.md: ### A human decision (a human decision: an author ID in `approvers.tsv`, an empty App field, the 
FAIL docs/spec/run.md: ### Copy before read (copy before read: each comment copied before a step acts on it; an edit is 
FAIL docs/spec/run.md: ### The Intake and control issues (the Intake issue and the control issue, and the first notice o
FAIL docs/spec/run.md: ### Input states (the input states that the happy path never makes)
FAIL docs/spec/run.md: ## NFR-006 — The target's pin at Start (the target's pin at Start (`NFR-006`))
FAIL docs/spec/run.md: ## The acceptance tests of M2a (the acceptance test of each part: unit, integration, e2e, uat)
FAIL docs/spec/run.md: ## The findings of the first pilot (the findings of the first pilot that `M2a` does not change)
FAIL docs/spec/forge.md: ## The six capabilities (the six capabilities of the forge interface)
FAIL docs/spec/forge.md: ## The App identity (the App identity: the JWT, the installation token, the host register of th
FAIL docs/spec/forge.md: ## The calls of M2a (the GitHub calls of `M2a`)
FAIL docs/spec/forge.md: ## Forge errors (a forge error stops a step, never a pass)
FAIL docs/spec/forge.md: ## The test of the adapter (the test of the adapter against a loopback server)
FAIL docs/spec/records.md: ## NFR-001 — The records of Start (the one-writer section of `NFR-001` for `layup run`, and t
FAIL docs/spec/packages.md: ### The table of M2a (the packages of `M2a` and rule 5 for the forge adapter)
exit 1
```
