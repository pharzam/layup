# T-hbw8 — the TypeSafe (Jev) documentation the panel read

Fetched on 2026-09-27 from `https://docs.typesafe.ai` (the index is
`https://docs.typesafe.ai/llms.txt`; each page is its path plus `.md`). The pages
were given to every panel member. The facts below are the ones the ADRs of this
task may rest on; each is quoted from the page named.

| Page | Fact used |
| ---- | --------- |
| `concepts/system-one` | "A System One model evaluates a state and returns typed answers and probabilities." "Jev currently accepts text input only." System One models "do not write replies, produce code, or generate explanations of their reasoning." |
| `primitives` and `primitives/choice`, `primitives/noul`, `primitives/score` | Three primitives: Choice ("one of a defined set"), Noul ("probability of yes"), Score ("probability-weighted position on ordered levels"). |
| `api` | One endpoint: `POST https://api.typesafe.ai/v1/systemone`; the `model` field selects the model; `429 Too Many Requests` when a rate limit is exceeded. |
| `models` | Jev 1.13 (`jev-1.13.0`, alias `jev-latest`): price $0.042 per million input tokens, output free; rate limits 250,000 tokens per second and 1,200 requests per minute, which "can change without notice"; context 64k tokens per request, 32k for the state plus the longest question. "The response's `model` field reports the versioned ID that answered"; pin a versioned ID when thresholds are tuned against it. "Jev is not fine-tuned or LoRA-adapted with customer data … the same weights serve every account." |
| `confidence` | "A confidence threshold is not one number"; thresholds depend on the consequences; "Start with conservative thresholds, test with your own data, and adjust." Noul answers carry no separate confidence. |
| `model-jaggedness/jev-1.13` (reviewed 2026-09-17) | Known weak points: literal reading, math and numbers, counting, date and time comparison, indirection, large state with irrelevant detail, adversarial content, contradictory instructions, structural invariants, generation. "Keep the arithmetic in code"; "compare in code". |
| `concepts/how-to-build-with-system-one`, `patterns/composite-scoring`, `cookbooks/sde_cascade` | Code owns the workflow; ask independent questions over one state together; send uncertain or failing cases to a person or a reasoning model. |
