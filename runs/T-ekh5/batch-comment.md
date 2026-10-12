LAYUP Intake: the batch of questions

## For the idea owner

Q-001: Which technology stack does the product use (languages, frameworks, tools)?

Q-002: The capped domain path after start needs at least KB 10s + two decisions 20s + generation 60s = 90s before finalization, so a turn started with 60-89s left is bound to fail a stage guard; should the start guard be raised to the sum of required stage caps, or should stage guards use something other than the caps?
`A start-budget guard MUST refuse to start a turn when less than 60 seconds of processing budget remains after reserving finalization time.`

Q-003: Does the 60-second queue time count inside the 180-second deadline that starts at handler entry (leaving at most 120s for execution), or is it a separate allowance on top of a full 180-second execution budget?
`independently of the 180-second total attempt deadline.`

Q-004: Has it been confirmed that the reviewed Jev 0.2.0 Choice response returns a probability for every option (needed for a top-two margin), and what rule applies if it returns only the chosen option or its probability?
`The default intent-confidence policy accepts a choice only when its top probability is at least 0.70 and its margin over the second choice is at least 0.20.`

Q-005: Is "ambiguous intent" exactly a classification that fails the 0.70/0.20 confidence policy, or can it also come from some other signal, given that the four intent classes contain no "ambiguous" value?
``**R-DEC-04** Ambiguous intent MUST produce `ask_clarification / intent_ambiguous` without guessing a domain workflow.``

Q-006: Who will supply the existing KB's wire contract, and by when? Can that KB enforce tenant/subject/role access scope at all? If it can't, what integration boundary is approved, since Gate A cannot complete without it?
`| KB wire contract | Actual method/path, schemas, auth/access propagation, corpus/filter mapping, byte/result limits, revision and failure/billing semantics. This is not yet supplied by the source material. |`

Q-007: What is the mapping from allowlisted JWT roles to permitted corpus IDs and filter values, which roles are on the allowlist, and what does "authorization scope" mean beyond roles?
`**R-AUTH-05** Retrieval permissions MUST be derived server-side from trusted roles and authorization scope rather than from the session identifier.`

Q-008: How is a deterministic preparer meant to pick the "relevant referent" for a pronoun-based Turkish or English follow-up without a model? Is a concrete rule set designed, or should optional LLM rewriting be the expected path for R-KB-02?
`The default query preparer is deterministic: combine the current question with relevant product/session facts and recent user references.`

Q-009: Is the smalltalk acknowledgement a fixed per-locale template, or an LLM-generated reply under a bounded prompt (which would count against the generation call and usage budget)?
``**R-POLICY-02** `smalltalk` MUST use a bounded localized response without a KB call or unrestricted domain generation.``

Q-010: Which conversation-meta cases besides "first user question" need deterministic selection, and what decides when the bounded LLM is used to phrase a meta answer instead?
`Common meta cases, such as the first user question, use deterministic transcript selection.`

Q-011: What threshold triggers summarization (for example, what fraction of which role's window the unsummarized tail may reach), and how large is a summary chunk by default?
`**R-CTX-06** Summary work MUST trigger from configured role budgets and MUST NOT run unconditionally on every user message.`

Q-012: Is a "principal" for quota purposes the (tenant_id, subject) pair, the subject alone, or the full normalized principal including roles and authz_version?
``| `max_active_per_principal` / `max_queued_per_principal` | 10 / 10 |``

Q-013: What happens when a principal already holds 1,000 retained sessions: reject creation (with which status and code, since none is in the catalog), or evict or expire the oldest session?
`a default maximum of 1,000 retained sessions per principal.`

Q-014: By what rule should the 30-second cold estimate be updated from observed durations (for example, which percentile, window or decay), so that "conservative" can be tested?
`update conservatively from observed durations`

Q-015: What concrete test shows that "further progress is demonstrably possible" (for example, the next unsummarized chunk fits the summarizer window), as opposed to a context that will never converge?
`when durable coverage advanced and further progress is demonstrably possible.`

Q-016: Does "answer rate" count any `answer / grounded` outcome, or only answers that are correct and contain the required facts?
`| Answer rate on labelled answerable domain cases | At least 80%. |`

Q-017: What makes a claim "consequential", and who annotates the critical subset and by what criteria?
`At least 95%; zero unsupported consequential claims in the annotated critical subset.`

Q-018: How is the response language judged in Gate C (human review or a named language detector), and how are mixed-language answers, such as Turkish text quoting English source titles, scored?
`| Requested response language | At least 95% for generated text; 100% for deterministic templates. |`

Q-019: What minimum cohort size or suppression rule makes an aggregate "non-linkable", so that A32 can test it?
`Small or identifiable cohorts require aggregation/suppression under the retention policy.`

Q-020: Who writes and reviews the 60 labelled cases and their synthetic KB evidence corpus, and in which domain, given that the real KB contract and content are not yet supplied?
`**R-QUAL-01** The project MUST maintain at least 20 calibration cases and 40 separate held-out synthetic evaluation cases with reviewed labels and fixed evidence snapshots.`

Q-021: Who are the human reviewers for Gate C, who writes and approves the rubric, and is an LLM-based assessor acceptable as the "audited assessor"?
`**R-QUAL-06** Groundedness and summary-fidelity scoring MUST have a documented human-review rubric or audited assessor, with human review of every critical failure and disputed machine score.`

Q-022: For the pilot and for the first real-data deployment, are historical answers snapshots of formerly authorized evidence, or must they be re-authorized at the document level? Who owns that decision?
`**R-SEC-08** A deployment MUST document whether historical answers are snapshots of formerly authorized evidence or require current document-level reauthorization.`

Q-023: Only the 100-turn bound is quantitative while "preserves order within each session" and "bounds waiting work" are qualitative obligations — should this span be split into a measure plus needs, or reclassed?
`The service admits up to 100 active turns by default, preserves order within each session and bounds waiting work.`

Q-024: The next span calls this a planning approximation rather than a guarantee — is it a normative measure or supporting rationale that should be context, which would leave draft R-196 unsupported?
`Under the simplifying 100-worker/30-second-turn estimate, a 60-second wait budget gives about 200 starts before applying headroom; the selected 100-waiter cap leaves headroom.`

Q-025: This describes how the start and stage guards apply the budget rather than stating a measurable bound — should it be classed constraint (or context) instead of measure?
`A cold domain path with two summaries fits within the 180-second budget only when actual queue/stage times permit it; the start guard and stage guards decide before dispatch.`

Q-026: A12 is a required conformance-test proof with no quantitative target — should this A-family row be classed need rather than measure?
`| A12 | Anthropic and OpenAI-compatible adapters pass shared conformance; unsupported mandatory capabilities fail without core changes. |`

Q-027: A26 lists qualitative failure-mode proofs with no numeric threshold — should it be classed need rather than measure?
`| A26 | Fabricated/no citations, cited-subset ordering, wrong-session message references, empty output, length truncation and provider refusal mapping. |`

Q-028: A33 states gate-blocking conditions rather than a measurable target — should it be classed need or constraint rather than measure?
`| A33 | Missing or incompatible KB/identity contract blocks applicable Gate A completion; missing credentials blocks live gates without fake-backed success claims. |`

Q-029: A36 describes required evidence and separation of bootstrap from product acceptance, with no numeric criterion — should it be classed need rather than measure?
`| A36 | Complete requirement-to-test mapping, configuration/source pins, delivery evidence and separation of repository bootstrap from product/integration acceptance. |`

Q-030: A37 states a security invariant plus a qualitative Gate C check — should it be classed need or constraint rather than measure?
`| A37 | Untrusted KB/summary instructions cannot alter configured endpoint, owner, policy or execute tools; live content-following behaviour is additionally measured in Gate C. |`

Q-031: Fact 663 (R-GATE-02) is classed need for a gate-pass obligation while this structurally identical Gate B obligation is measure — does "at least one real LLM" justify measure, or should it be need?
`**R-GATE-03** Gate B MUST separately exercise the actual KB, pinned Jev profile and at least one real LLM using sanitized reproducible integration evidence.`

Q-032: B, in USD: below it the delivery goes on without asking you. A decimal, below U.

Q-033: U, in USD: the upper edge of the band, the money appetite of the delivery. A decimal, above B.

Q-034: The wall-clock appetite of the delivery, in hours. A decimal.

## For the Operator

No setup question in M2c.

Answer each block in one comment, one line per question ID: `Q-NNN: <answer>`.
