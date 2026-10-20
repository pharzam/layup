# T-c826 merged-task replays

These are clean historical heads in disposable local shared clones under
`/home/layup/.cache/layup-session/tmp-green/T-c826-stage-c`, using the Stage C
checker as an external script. Target HEADs are unchanged; the tested checker
is the uncommitted implementation identified by hash in test-runs.md. Git clone,
checkout, consumer inspection, committed record reads, gross numstat and checker
commands each have measured UTC endpoints and direct status in that record.

## Classification before declarations

For T-3ff9, inspected test-runs.md (run times and results), conditions.md
(condition/evidence mapping), semantic.md (comparison with authoritative clauses)
and mentions.tsv (reference audit). Searched committed path consumers with
`git grep`; the hits outside those records are task-document citations and task-check
invocations. Inspected the full T-3ff9 checks.sh: it reads operative documents and
Git objects, not these records as execution inputs. It remains counted.

For T-b6fn, inspected test-runs.md (runs, diagnostics and mutation evidence),
conditions.md (citations) and real-heads.md (historical replay outputs). The
committed path search finds only task-document evidence citations outside the
records. The checker reads Markdown in both task directories for K3 selectors,
K4 restoration evidence and K6 condition citations; those are audits of evidence,
not execution-data uses under the canonical input definition. None of the declared
records is fixture data, expected output, executable content or configuration.
The classification covers inspected repository consumers, not unknown external
uses. Declarations are the author's classification, as the checker comments state.

## Arithmetic and other checks

The figures below remove only whole declared evidence records. They do not remove
mixed-file close-out lines, and they do not retroactively decide a historical
landing under this task's delegated tolerance.

| Task | Head | Named base | Gross lines/files | Net figure | Maximum | Tolerance | Limit | Margin | K5 | Other checks / process exit |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| T-3ff9 | `0dfac539061cda59b233ab222dca76d111c72db8` | `df9bfacf33054741564dbf0254a741f0d4be1e0f` | 547/16 | 318/12 | 650/20 | 130/4 | 780/24 | 462/12 | ok | K1–K4 and K6 ok; exit 0 |
| T-b6fn | `78386339e297ed1e1a4bd8ff35d2af5189bbfd6f` | `99ff86235b726f8842120dafb2356b1aff035238` | 2363/238 | 2023/235 | 1300/170 | 260/4 | 1560/174 | -463/-61 | FAIL | K1–K4 and K6 ok; exit 1 |

### T-3ff9

Declared evidence: `runs/T-3ff9/test-runs.md`, `runs/T-3ff9/conditions.md`, `runs/T-3ff9/semantic.md`, `runs/T-3ff9/mentions.tsv`.

Command in the clean historical clone:

```sh
sh /home/layup/projects/layup/.worktree/T-c826/docs/tests/freeze-check.sh --task T-3ff9 --base df9bfacf33054741564dbf0254a741f0d4be1e0f --budget 650/20 --evidence runs/T-3ff9/test-runs.md --evidence runs/T-3ff9/conditions.md --evidence runs/T-3ff9/semantic.md --evidence runs/T-3ff9/mentions.tsv
```

2026-10-20T21:21:49Z to 2026-10-20T21:21:50Z, direct exit 0. Complete checker output:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
ok K5: projection: figure 318/12, maximum 650/20, tolerance 130/4, limit 780/24, margin 462/12 (the figure at landing is still measured at landing)
ok K6: conditions table
```

### T-b6fn

Declared evidence: `runs/T-b6fn/test-runs.md`, `runs/T-b6fn/conditions.md`, `runs/T-b6fn/real-heads.md`.

Command in the clean historical clone:

```sh
sh /home/layup/projects/layup/.worktree/T-c826/docs/tests/freeze-check.sh --task T-b6fn --base 99ff86235b726f8842120dafb2356b1aff035238 --budget 1300/170 --evidence runs/T-b6fn/test-runs.md --evidence runs/T-b6fn/conditions.md --evidence runs/T-b6fn/real-heads.md
```

2026-10-20T21:21:50Z to 2026-10-20T21:21:51Z, direct exit 1. Complete checker output:

```text
ok K1: white space
ok K2: registration of new tests
ok K3: test selectors
ok K4: restore evidence of mutants
FAIL K5: -:0: projection: figure 2023/235, maximum 1300/170, tolerance 260/4, limit 1560/174, margin -463/-61 (the figure at landing is still measured at landing)
ok K6: conditions table
```

T-3ff9's historical task record uses its earlier evidence-inclusive 650/20 bound;
this replay demonstrates the new projection on its data and grants no new
historical approval. T-b6fn's recorded landing gross figure is 2363/238 and its
task record names an Operator approval of that exact figure. Here its inspected
records remove 340 lines and 3 files, leaving 2023/235; both dimensions exceed
1560/174. The process failure is K5 alone, which is reported as failure rather
than treated as a failed replay environment or an unrelated K1–K6 defect.
