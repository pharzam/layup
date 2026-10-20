## The problem
### The requirements
- REQ-001: Alpha Criterion: a check
- NFR-002: p99 under 200 ms Criterion: measured at the edge
### The success criteria
- faster
The problem in prose.
## The appetite
appetite: USD 8, 12 hours
Why: small.
## The solution
### The verified heads
- T-sp01 (specification): [pull request #7](https://github.com/acme/target/pull/7), head 4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f4f
- T-ar01 (architecture): [pull request #8](https://github.com/acme/target/pull/8), head a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1
### Decision db
- pg (a survivor)
- sqlite (a survivor)
- mysql (dropped: NFR-002: fail)
### Decision web
- go (the one survivor, decided at Shape)
- rust (dropped: Inv-1: fail; ADR-0003: fail)
## The rabbit holes and the no-gos
### The questions still open
- Q-002: The question Q-002
- Q-003: The question Q-003
- Q-007: The question Q-007
architecture change: no
No rabbit holes but one.
## The facts
- 1 (need): The tool runs offline.
- 2 (constraint): The data stays in Postgres.
- Q-001 (context): 5.0
- Q-004 (constraint): Postgres only
