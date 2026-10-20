# T-c826 independent mutation evidence

The original S5 sweep's **52** variants are killed by their named assertions: 35 checker variants
run the entire 81-case suite, and 17 documentation/CI/preservation variants run
the unchanged PR-numbered task check. No runner metadata, executable-mode,
symlink, pin, blob-comparison or restore-replay rule is introduced.

Each variant starts from a fresh copy outside the worktree under
`/home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/Mxx`. The checker
runner has no FREEZE_CHECK hook, so its entire `docs/tests` tree is copied.
Every shell variant passes `sh -n`. A killed checker variant has suite exit 1,
a named full-output/status mismatch (case harness exit 2), and no syntax or
unbound-variable error. Task-check variants exit 1 on the named assertion.
After each run, the original bytes are restored from that variant's fresh backup;
the restored hash equals both the before hash and the unchanged worktree source,
and the named green control exits 0. Every check's command, actual UTC endpoints
and direct exit are recorded by run ID in [test-runs.md](test-runs.md).

A driver uniqueness assertion stopped before changing M37 or running a test;
that attempt is not a mutant or kill. The resumed edits are scoped to the one
operative job. IDs retain the real scratch/run labels, so they are not contiguous.

The table gives SHA-256 of the one changed file. Failure excerpts below are cut
to the named case or assertion; they are not a claim that only that case failed.

| ID | Single rule removed or changed | Named case / assertion | Changed file | Before | During | After | Verdict / restored control |
| --- | --- | --- | --- | --- | --- | --- | --- |
| M01 | Exclude Markdown records | `good-K5-evidence-kinds` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 1ce27efa4a7a1571bb4b4a9e676df955660fc865883647457512350454e8b721 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M01-restored` exit 0 |
| M02 | Exclude TSV records | `good-K5-evidence-kinds` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | add3f8b0d0b0bc1c02606aad763feb4f443671872f68e53304281c8f9762d056 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M02-restored` exit 0 |
| M03 | Exclude text records | `good-K5-evidence-kinds` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 08924f0f5fb16edb1eaedb068380307d84144c2f40cdd032e7714aae75088562 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M03-restored` exit 0 |
| M04 | Accumulate repeated evidence options | `good-K5-evidence-kinds` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 16d9a4a9cda4ea8a5545921d65724828da4d3ce4b62e4c5f2b89457c6521f933 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M04-restored` exit 0 |
| M05 | Exclude evidence lines | `good-K5-evidence-kinds` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 9cc8d62e3252c4cf05d5c2ad1640fd4fd06b7b42e53fefeb517842bad59188f7 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M05-restored` exit 0 |
| M06 | Exclude evidence files | `good-K5-evidence-kinds` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | e24243a9a4603a733c1486c9e824e701f03d4ab2d96569708bd3b8ff9349de5a | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M06-restored` exit 0 |
| M07 | Count unlisted execution inputs | `good-K5-evidence-inputs` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 85d17c8720b2a6b64f3472f8075557740ac98688e156c8a067d860a7495cc48c | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M07-restored` exit 0 |
| M08 | Accept added committed endpoint (base-only mutant) | `good-K5-evidence-added` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | b0690d5e2f8caf0b853736ae1a17bc0318f37bb2880b4a2a1c267371552dbffa | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M08-restored` exit 0 |
| M09 | Accept deleted committed endpoint (HEAD-only mutant) | `good-K5-evidence-deleted` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | aec01737ea73503610b62d9ccc86b89c3dcd48131f36b2f093181bf0d626d4f5 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M09-restored` exit 0 |
| M10 | Subtract modified-record additions | `good-K5-evidence-modified` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | d82678178713bd21fa342ba17e06ba05e3bda808f0c78b8a3e10c9ebbe3937f0 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M10-restored` exit 0 |
| M11 | Subtract modified-record removals | `good-K5-evidence-modified` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 3f18e4bad02a4e0e959e28cb63f5b199a3b34f0af4b2e9c0cc47251036c2aeee | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M11-restored` exit 0 |
| M12 | Exact task directory | `bad-usage-evidence-other-task` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 03d601de9c49da174810f8d73a151706ffe5dc8d8b6f578769b6dac7df446290 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M12-restored` exit 0 |
| M13 | Reject task prefix neighbour | `bad-usage-evidence-task-prefix` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 8a98a344e686d50851d2c0470f81040356bd2c75f29ad5daff306d31e43a805b | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M13-restored` exit 0 |
| M14 | Reject traversal before results | `bad-usage-evidence-traversal` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 25f3a59c0ef37b10ee3547c59a62ef894e188f9628432b7fbd17af63f0eae8be | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M14-restored` exit 0 |
| M15 | Reject unsupported record suffix | `bad-usage-evidence-suffix` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 0d570f65f47c2da77a92491055469fd737c40330b00c764b98825ab4024de27b | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M15-restored` exit 0 |
| M16 | Require committed presence | `bad-usage-evidence-missing` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | fcc3f9e3f2215137c2081889398c3fda6dcf2d3b0c93f62071da4eb3687ecb75 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M16-restored` exit 0 |
| M17 | Match declared paths literally | `good-K5-evidence-literal` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | fdcd294c0e8d0d060dc963403eb6f65aca2268eb5e82b5b74355a52a6596918b | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M17-restored` exit 0 |
| M18 | Deduplicate declarations | `good-K5-evidence-duplicate` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 4f9468291ca576b0070c2c53e8f32453af0fe575105edf562db3cb6eee687e09 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M18-restored` exit 0 |
| M19 | Line tolerance floor | `good-K5-line-limit` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | e8827598c11e1581efa30ea64c26a128e9b742798b5ddbc742c4088f0bf78d6b | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M19-restored` exit 0 |
| M20 | Line tolerance percentage | `good-K5-percentage` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 173c214565b6b58b7d5b0e526146c1493f9720d7996f3f0d14dcd1a316b93694 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M20-restored` exit 0 |
| M21 | Line tolerance ceiling | `good-K5-ceiling` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 58518e26fdcdd5277cc4cacb050c5a321a160b0624aeb919c1ca8eeb1f51705e | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M21-restored` exit 0 |
| M22 | Whole-line truncation | `good-K5-truncation` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 66f370ee6129f7baa8cb33481e54c28584cc8a41490c3f777362aed09df78c2e | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M22-restored` exit 0 |
| M23 | Four-file tolerance | `good-K5-file-limit` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 257154d7887a830b555f1a36e69d3373b86dc42a6f69254628ea52b1d08b7bc7 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M23-restored` exit 0 |
| M24 | Reject over line limit | `bad-K5-over-lines` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 394b90a85a58bb06b1a52dc52761739c73eeb3fce11a3862363fe76e4cc23e8d | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M24-restored` exit 0 |
| M25 | Accept equal line limit | `good-K5-line-limit` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 2e55dc35c152abb75d39e4d412b23dfc763087839df6345bb8bf37736fd0eb30 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M25-restored` exit 0 |
| M26 | Reject over file limit | `bad-K5-over-files` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 80ee4a1e4d7ec194c8c0c36d590e824c9363a36149650aae150d125f12bb9680 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M26-restored` exit 0 |
| M27 | Accept equal file limit | `good-K5-file-limit` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | fc76815bb438bafacd65089794eb4fa4b33a69398f9df4719708ed68509cde4d | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M27-restored` exit 0 |
| M28 | Line margin includes tolerance | `good-K5-line-limit` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | e2d86514a2c38107b0988991e721682b851f48b927e393a2599417737c8dc2a7 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M28-restored` exit 0 |
| M29 | File margin includes tolerance | `good-K5-file-limit` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | a115d4eaf1f48a76ddf7a92108abf865da14e4c3b3a12e88c0fcaad7e6a3e882 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M29-restored` exit 0 |
| M30 | No synthetic line allowance | `good-K5-no-closeout` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 96c3a461a4e1e0e7cb4802a3d323f47cc231b86e78a0836bb169f089b68cbd22 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M30-restored` exit 0 |
| M31 | No synthetic file allowance | `good-K5-no-closeout` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | d1f3e060b2a2a2ca84b326c5e7090a0e207cbdd15cb51fa19205460e533b98f1 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M31-restored` exit 0 |
| M32 | Reject retired closeout option | `bad-usage-closeout` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 1eec6855ef5cb03a96e9df4cf811e06bfe88a11c09dc92a6e5f7d949a228fedc | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M32-restored` exit 0 |
| M33 | Require evidence argument | `bad-usage-evidence-no-value` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 0692d4430f77242606399ee4270df7536a78443e8c96423d80a7f623bf5cfc3a | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M33-restored` exit 0 |
| M34 | Canonical margin assertion | `canonical margin sentence` | `docs/engineering-discipline.md` | d595f80abf238f38a37d4f55143eda5670555c926b44f282433919de498d5d8e | c41f605ff23e385f47a4c4e60915d4958a11bb77d06cf4091a1cc4cd6c41068c | d595f80abf238f38a37d4f55143eda5670555c926b44f282433919de498d5d8e | KILLED; `M34-restored` exit 0 |
| M35 | Preserve run.sh | `unchanged runner: docs/tests/freeze-check/run.sh` | `docs/tests/freeze-check/run.sh` | 9443c510830f705fdc3f81732270db9f406cb0302cb2f27e18b414392899a4b0 | b19c4ba8eb619c9b13290b8a0063d8377075147d5a69a758bfbb57667f083a96 | 9443c510830f705fdc3f81732270db9f406cb0302cb2f27e18b414392899a4b0 | KILLED; `M35-restored` exit 0 |
| M36 | Preserve run-all.sh | `unchanged runner: docs/tests/freeze-check/run-all.sh` | `docs/tests/freeze-check/run-all.sh` | b892114bdc1719b724c46ba2242941e1f8bcf00d3945ce95877dd265e7a26bb8 | 625395e9114db97ddc14e2f5f33c119d43ee0250ff2021ae7c6e78b14b6d8699 | b892114bdc1719b724c46ba2242941e1f8bcf00d3945ce95877dd265e7a26bb8 | KILLED; `M36-restored` exit 0 |
| M51 | Line limit adds tolerance | `good-K5-line-limit` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | e7880138ac21bc2674d9b23399d0dd7b8d31bab967bebe6ca8883ebf761714af | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M51-restored` exit 0 |
| M52 | File limit adds tolerance | `good-K5-file-limit` | `docs/tests/freeze-check.sh` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | f708171d3c9a38db747457def6014261296c443a3427220b4791a7e4470a5f5c | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; `M52-restored` exit 0 |
| M53 | Retain restore-list path docs/adr/adr-lint.sh | `all nine baseline paths` | `.github/workflows/ci.yml` | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | 7f6d5f03a54c2c16526bee27744e8262a1f0ff5fd431905b3f3019587d411b05 | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | KILLED; `M53-restored` exit 0 |
| M54 | Retain restore-list path docs/tests/freeze-check/run.sh | `all nine baseline paths` | `.github/workflows/ci.yml` | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | 2158b90f4ffa60ac021482515d16e2507e31c26b11ebe6683904099c4d3f09fc | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | KILLED; `M54-restored` exit 0 |
| M55 | Retain restore-list path docs/tests/freeze-check/run-all.sh | `all nine baseline paths` | `.github/workflows/ci.yml` | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | f8f8526016591811c81abcab1a4c848e27c8981d523d920d3da08fb95ac46fab | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | KILLED; `M55-restored` exit 0 |
| M56 | Require pull_request guard | `kept checker requires the pull_request event` | `.github/workflows/ci.yml` | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | 8fcfb19dd42bac517b48caa883395b5f80075faed97fe26c9bbb7720976588d7 | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | KILLED; `M56-restored` exit 0 |
| M57 | Require PR-number guard | `one literal PR number and only the checker path` | `.github/workflows/ci.yml` | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | 17cef89327be58f2654e04319e7905b0e8db1c92339baa619481480e2d215e6e | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | KILLED; `M57-restored` exit 0 |
| M58 | Keep exactly one checker path | `one literal PR number and only the checker path` | `.github/workflows/ci.yml` | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | 93707ee64710a5a6d993938a707252eee6e6e320a514ee64493326f398598afa | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | KILLED; `M58-restored` exit 0 |
| M59 | Match actual approved number | `exception matches the supplied actual migration PR` | `.github/workflows/ci.yml` | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | cf4b65d6e6e7a0de645c555373810cec3b5d815167bda349910049a3b0c403fb | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | KILLED; `M59-restored` exit 0 |
| M60 | Read actual event | `event guard receives the actual GitHub event name` | `.github/workflows/ci.yml` | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | 2e6bd0577aa46bbb36219687909067c3b6329b2861fbfbb324192d4ce436ddb2 | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | KILLED; `M60-restored` exit 0 |
| M61 | Read actual PR number | `number guard receives the actual pull request number` | `.github/workflows/ci.yml` | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | 24330f0dbccd1aa0c8af168a1e311ee0b0b89c73eba0ccf9adb83f3405b7f47d | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | KILLED; `M61-restored` exit 0 |
| M62 | Restore run.sh even if listed | `only the guarded checker skips restoration` | `.github/workflows/ci.yml` | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | 3879df893ff8b24aa0363234511b55d0cc857ea2b966ea0ea4bd399b7904135a | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | KILLED; `M62-restored` exit 0 |
| M63 | Restore run-all.sh even if listed | `only the guarded checker skips restoration` | `.github/workflows/ci.yml` | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | 2b5f02c8f31cc66dfc53aeb0822c84f888b7e44b2d259a59ab67c6ae37316f81 | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | KILLED; `M63-restored` exit 0 |
| M64 | Use default-branch source | `restoration still reads the default branch of this repository` | `.github/workflows/ci.yml` | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | 30384e8a05b1580eb07994140526966094d2d7b20cbeb70832f8ea3cd4fee9fb | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | KILLED; `M64-restored` exit 0 |
| M65 | Read repository default branch | `default branch still comes from the repository event` | `.github/workflows/ci.yml` | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | 5f238ab19eb202f66100c3f64d36038111de5072618b47743a94f546c0fcd737 | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | KILLED; `M65-restored` exit 0 |
| M66 | Name removal issue | `exception comment names the task, literal PR and removal issue` | `.github/workflows/ci.yml` | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | 6273bbf23e5b03ed4d055b0dfc39f062305dcaf3bdeed17b55f48947df7b5018 | a450dce6c00be6a671b409d002eaf0695eba743b3bb540270c58b8dc89060e9f | KILLED; `M66-restored` exit 0 |

## M01 — Exclude Markdown records

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	*.md|*.tsv|*.txt) ;;
+	*.tsv|*.txt) ;;
```

`M01-mutant`: 2026-10-20T21:20:04Z to 2026-10-20T21:20:11Z, exit 1. Named observed failure:

```text
FAIL  good-K5-evidence-kinds
      run: good-K5-evidence-kinds: exit 2 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.gJlPxOyuTt/expect	2026-10-21 00:20:10.749379346 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.gJlPxOyuTt/got	2026-10-21 00:20:10.748695507 +0300
      @@ -1,7 +1,2 @@
      -0
      -ok K1: white space
      -ok K2: registration of new tests
      -ok K3: test selectors
      -ok K4: restore evidence of mutants
      -ok K5: projection: figure 2/2, maximum 0/0, tolerance 60/4, limit 60/4, margin 58/2 (the figure at landing is still measured at landing)
      -ok K6: conditions table
      +2
      +freeze-check: --evidence runs/T-0000/report.md: must end in .md, .tsv or .txt
```

`M01-restored`: 2026-10-20T21:20:11Z to 2026-10-20T21:20:11Z, exit 0; restored SHA-256 equals Before and After in the table.

## M02 — Exclude TSV records

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	*.md|*.tsv|*.txt) ;;
+	*.md|*.txt) ;;
```

`M02-mutant`: 2026-10-20T21:20:12Z to 2026-10-20T21:20:20Z, exit 1. Named observed failure:

```text
FAIL  good-K5-evidence-kinds
      run: good-K5-evidence-kinds: exit 2 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.dBpUKK889P/expect	2026-10-21 00:20:19.288375606 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.dBpUKK889P/got	2026-10-21 00:20:19.287375607 +0300
      @@ -1,7 +1,2 @@
      -0
      -ok K1: white space
      -ok K2: registration of new tests
      -ok K3: test selectors
      -ok K4: restore evidence of mutants
      -ok K5: projection: figure 2/2, maximum 0/0, tolerance 60/4, limit 60/4, margin 58/2 (the figure at landing is still measured at landing)
      -ok K6: conditions table
      +2
      +freeze-check: --evidence runs/T-0000/times.tsv: must end in .md, .tsv or .txt
```

`M02-restored`: 2026-10-20T21:20:20Z to 2026-10-20T21:20:20Z, exit 0; restored SHA-256 equals Before and After in the table.

## M03 — Exclude text records

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	*.md|*.tsv|*.txt) ;;
+	*.md|*.tsv) ;;
```

`M03-mutant`: 2026-10-20T21:20:20Z to 2026-10-20T21:20:29Z, exit 1. Named observed failure:

```text
FAIL  good-K5-evidence-kinds
      run: good-K5-evidence-kinds: exit 2 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.HiZksyYNov/expect	2026-10-21 00:20:27.883310946 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.HiZksyYNov/got	2026-10-21 00:20:27.882372082 +0300
      @@ -1,7 +1,2 @@
      -0
      -ok K1: white space
      -ok K2: registration of new tests
      -ok K3: test selectors
      -ok K4: restore evidence of mutants
      -ok K5: projection: figure 2/2, maximum 0/0, tolerance 60/4, limit 60/4, margin 58/2 (the figure at landing is still measured at landing)
      -ok K6: conditions table
      +2
      +freeze-check: --evidence runs/T-0000/output.txt: must end in .md, .tsv or .txt
```

`M03-restored`: 2026-10-20T21:20:29Z to 2026-10-20T21:20:29Z, exit 0; restored SHA-256 equals Before and After in the table.

## M04 — Accumulate repeated evidence options

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-		evidence=$evidence$2'
+		evidence=$2'
```

`M04-mutant`: 2026-10-20T21:20:29Z to 2026-10-20T21:20:37Z, exit 1. Named observed failure:

```text
FAIL  good-K5-evidence-kinds
      run: good-K5-evidence-kinds: exit 0 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.Txq6kEJAiZ/expect	2026-10-21 00:20:36.552330108 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.Txq6kEJAiZ/got	2026-10-21 00:20:36.551315566 +0300
      @@ -3,5 +3,5 @@
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 2/2, maximum 0/0, tolerance 60/4, limit 60/4, margin 58/2 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 5/4, maximum 0/0, tolerance 60/4, limit 60/4, margin 55/0 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M04-restored`: 2026-10-20T21:20:37Z to 2026-10-20T21:20:37Z, exit 0; restored SHA-256 equals Before and After in the table.

## M05 — Exclude evidence lines

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	fig_l=$((fig_l - $1))
+	fig_l=$fig_l
```

`M05-mutant`: 2026-10-20T21:20:37Z to 2026-10-20T21:20:46Z, exit 1. Named observed failure:

```text
FAIL  good-K5-evidence-kinds
      run: good-K5-evidence-kinds: exit 0 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.LhU9Qkqgg3/expect	2026-10-21 00:20:45.288365663 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.LhU9Qkqgg3/got	2026-10-21 00:20:45.287365663 +0300
      @@ -3,5 +3,5 @@
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 2/2, maximum 0/0, tolerance 60/4, limit 60/4, margin 58/2 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 6/2, maximum 0/0, tolerance 60/4, limit 60/4, margin 54/2 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M05-restored`: 2026-10-20T21:20:46Z to 2026-10-20T21:20:46Z, exit 0; restored SHA-256 equals Before and After in the table.

## M06 — Exclude evidence files

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	fig_f=$((fig_f - $2))
+	fig_f=$fig_f
```

`M06-mutant`: 2026-10-20T21:20:46Z to 2026-10-20T21:20:54Z, exit 1. Named observed failure:

```text
FAIL  good-K5-evidence-kinds
      run: good-K5-evidence-kinds: exit 1 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.CoMlarZQ3d/expect	2026-10-21 00:20:53.822333990 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.CoMlarZQ3d/got	2026-10-21 00:20:53.821362861 +0300
      @@ -1,7 +1,7 @@
      -0
      +1
       ok K1: white space
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 2/2, maximum 0/0, tolerance 60/4, limit 60/4, margin 58/2 (the figure at landing is still measured at landing)
      +FAIL K5: -:0: projection: figure 2/5, maximum 0/0, tolerance 60/4, limit 60/4, margin 58/-1 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M06-restored`: 2026-10-20T21:20:54Z to 2026-10-20T21:20:55Z, exit 0; restored SHA-256 equals Before and After in the table.

## M07 — Count unlisted execution inputs

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-		git diff --no-renames --numstat "$base" HEAD -- ":(literal)$ev" >> "$tmp/k5_evidence"
+		git diff --no-renames --numstat "$base" HEAD -- ":(literal)runs/$task/" >> "$tmp/k5_evidence"
```

`M07-mutant`: 2026-10-20T21:20:55Z to 2026-10-20T21:21:03Z, exit 1. Named observed failure:

```text
FAIL  good-K5-evidence-inputs
      run: good-K5-evidence-inputs: exit 0 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.k9CHrMspmF/expect	2026-10-21 00:21:02.052090756 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.k9CHrMspmF/got	2026-10-21 00:21:02.051300047 +0300
      @@ -3,5 +3,5 @@
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 6/3, maximum 0/0, tolerance 60/4, limit 60/4, margin 54/1 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 0/0, maximum 0/0, tolerance 60/4, limit 60/4, margin 60/4 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M07-restored`: 2026-10-20T21:21:03Z to 2026-10-20T21:21:03Z, exit 0; restored SHA-256 equals Before and After in the table.

## M08 — Accept added committed endpoint (base-only mutant)

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	if ! git cat-file -e "$base:$ev" 2>/dev/null && ! git cat-file -e "HEAD:$ev" 2>/dev/null; then
+	if ! git cat-file -e "$base:$ev" 2>/dev/null; then
```

`M08-mutant`: 2026-10-20T21:21:03Z to 2026-10-20T21:21:11Z, exit 1. Named observed failure:

```text
FAIL  good-K5-evidence-added
      run: good-K5-evidence-added: exit 2 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.ZmnJnBjvzt/expect	2026-10-21 00:21:09.791358208 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.ZmnJnBjvzt/got	2026-10-21 00:21:09.790358208 +0300
      @@ -1,7 +1,2 @@
      -0
      -ok K1: white space
      -ok K2: registration of new tests
      -ok K3: test selectors
      -ok K4: restore evidence of mutants
      -ok K5: projection: figure 1/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 159/13 (the figure at landing is still measured at landing)
      -ok K6: conditions table
      +2
      +freeze-check: --evidence runs/T-0000/added.md: is absent from both base and HEAD
```

`M08-restored`: 2026-10-20T21:21:11Z to 2026-10-20T21:21:11Z, exit 0; restored SHA-256 equals Before and After in the table.

## M09 — Accept deleted committed endpoint (HEAD-only mutant)

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	if ! git cat-file -e "$base:$ev" 2>/dev/null && ! git cat-file -e "HEAD:$ev" 2>/dev/null; then
+	if ! git cat-file -e "HEAD:$ev" 2>/dev/null; then
```

`M09-mutant`: 2026-10-20T21:21:11Z to 2026-10-20T21:21:19Z, exit 1. Named observed failure:

```text
FAIL  good-K5-evidence-deleted
      run: good-K5-evidence-deleted: exit 2 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.HaqERCiihn/expect	2026-10-21 00:21:17.872356142 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.HaqERCiihn/got	2026-10-21 00:21:17.871356142 +0300
      @@ -1,7 +1,2 @@
      -0
      -ok K1: white space
      -ok K2: registration of new tests
      -ok K3: test selectors
      -ok K4: restore evidence of mutants
      -ok K5: projection: figure 1/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 159/13 (the figure at landing is still measured at landing)
      -ok K6: conditions table
      +2
      +freeze-check: --evidence runs/T-0000/test-runs.md: is absent from both base and HEAD
```

`M09-restored`: 2026-10-20T21:21:19Z to 2026-10-20T21:21:19Z, exit 0; restored SHA-256 equals Before and After in the table.

## M10 — Subtract modified-record additions

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-		{ a = ($1 == "-" ? 0 : $1); r = ($2 == "-" ? 0 : $2); l += a + r; f += 1 }
+		{ a = ($1 == "-" ? 0 : $1); r = ($2 == "-" ? 0 : $2); l += r; f += 1 }
```

`M10-mutant`: 2026-10-20T21:21:19Z to 2026-10-20T21:21:27Z, exit 1. Named observed failure:

```text
FAIL  good-K5-evidence-modified
      run: good-K5-evidence-modified: exit 0 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.cA42rRYJYE/expect	2026-10-21 00:21:26.749325569 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.cA42rRYJYE/got	2026-10-21 00:21:26.748354090 +0300
      @@ -3,5 +3,5 @@
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 1/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 159/13 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 2/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 158/13 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M10-restored`: 2026-10-20T21:21:27Z to 2026-10-20T21:21:27Z, exit 0; restored SHA-256 equals Before and After in the table.

## M11 — Subtract modified-record removals

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-		{ a = ($1 == "-" ? 0 : $1); r = ($2 == "-" ? 0 : $2); l += a + r; f += 1 }
+		{ a = ($1 == "-" ? 0 : $1); r = ($2 == "-" ? 0 : $2); l += a; f += 1 }
```

`M11-mutant`: 2026-10-20T21:21:27Z to 2026-10-20T21:21:36Z, exit 1. Named observed failure:

```text
FAIL  good-K5-evidence-modified
      run: good-K5-evidence-modified: exit 0 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.jIqxSGIk4D/expect	2026-10-21 00:21:35.112099050 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.jIqxSGIk4D/got	2026-10-21 00:21:35.111305455 +0300
      @@ -3,5 +3,5 @@
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 1/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 159/13 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 4/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 156/13 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M11-restored`: 2026-10-20T21:21:36Z to 2026-10-20T21:21:36Z, exit 0; restored SHA-256 equals Before and After in the table.

## M12 — Exact task directory

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	runs/"$task"/*) ;;
+	runs/*) ;;
```

`M12-mutant`: 2026-10-20T21:21:36Z to 2026-10-20T21:21:44Z, exit 1. Named observed failure:

```text
FAIL  bad-usage-evidence-other-task
      run: bad-usage-evidence-other-task: exit 0 (wanted 2); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.h8sbHoRxBY/expect	2026-10-21 00:21:41.616219674 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.h8sbHoRxBY/got	2026-10-21 00:21:41.615351155 +0300
      @@ -1,2 +1,7 @@
      -2
      -freeze-check: --evidence runs/T-0001/report.md: must be under runs/T-0000/
      +0
      +ok K1: white space
      +ok K2: registration of new tests
      +ok K3: test selectors
      +ok K4: restore evidence of mutants
      +ok K5: projection: figure 0/0, maximum 100/10, tolerance 60/4, limit 160/14, margin 160/14 (the figure at landing is still measured at landing)
      +ok K6: conditions table
```

`M12-restored`: 2026-10-20T21:21:44Z to 2026-10-20T21:21:44Z, exit 0; restored SHA-256 equals Before and After in the table.

## M13 — Reject task prefix neighbour

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	runs/"$task"/*) ;;
+	runs/"$task"*/*) ;;
```

`M13-mutant`: 2026-10-20T21:21:44Z to 2026-10-20T21:21:53Z, exit 1. Named observed failure:

```text
FAIL  bad-usage-evidence-task-prefix
      run: bad-usage-evidence-task-prefix: exit 0 (wanted 2); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.ZAua4syJ7w/expect	2026-10-21 00:21:50.335349720 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.ZAua4syJ7w/got	2026-10-21 00:21:50.334349720 +0300
      @@ -1,2 +1,7 @@
      -2
      -freeze-check: --evidence runs/T-0000-extra/report.md: must be under runs/T-0000/
      +0
      +ok K1: white space
      +ok K2: registration of new tests
      +ok K3: test selectors
      +ok K4: restore evidence of mutants
      +ok K5: projection: figure 0/0, maximum 100/10, tolerance 60/4, limit 160/14, margin 160/14 (the figure at landing is still measured at landing)
      +ok K6: conditions table
```

`M13-restored`: 2026-10-20T21:21:53Z to 2026-10-20T21:21:53Z, exit 0; restored SHA-256 equals Before and After in the table.

## M14 — Reject traversal before results

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	case /$ev/ in */../*|*/./*|*//*) evidence_error 'path traversal is not allowed' ;; esac
+	: # traversal validator removed
```

`M14-mutant`: 2026-10-20T21:21:53Z to 2026-10-20T21:22:01Z, exit 1. Named observed failure:

```text
FAIL  bad-usage-evidence-traversal
      run: bad-usage-evidence-traversal: exit 2 (wanted 2); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.3B8Rb7CAeJ/expect	2026-10-21 00:21:58.754348531 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.3B8Rb7CAeJ/got	2026-10-21 00:21:58.753348531 +0300
      @@ -1,2 +1,2 @@
       2
      -freeze-check: --evidence runs/T-0000/../T-0000/conditions.md: path traversal is not allowed
      +freeze-check: --evidence runs/T-0000/../T-0000/conditions.md: is absent from both base and HEAD
```

`M14-restored`: 2026-10-20T21:22:01Z to 2026-10-20T21:22:01Z, exit 0; restored SHA-256 equals Before and After in the table.

## M15 — Reject unsupported record suffix

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	case $ev in
-	*.md|*.tsv|*.txt) ;;
-	*) evidence_error 'must end in .md, .tsv or .txt' ;;
-	esac
+	: # suffix validator removed
```

`M15-mutant`: 2026-10-20T21:22:01Z to 2026-10-20T21:22:10Z, exit 1. Named observed failure:

```text
FAIL  bad-usage-evidence-suffix
      run: bad-usage-evidence-suffix: exit 0 (wanted 2); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.r1srKA8R5Y/expect	2026-10-21 00:22:07.121087765 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.r1srKA8R5Y/got	2026-10-21 00:22:07.120197396 +0300
      @@ -1,2 +1,7 @@
      -2
      -freeze-check: --evidence runs/T-0000/checks.sh: must end in .md, .tsv or .txt
      +0
      +ok K1: white space
      +ok K2: registration of new tests
      +ok K3: test selectors
      +ok K4: restore evidence of mutants
      +ok K5: projection: figure 0/0, maximum 100/10, tolerance 60/4, limit 160/14, margin 160/14 (the figure at landing is still measured at landing)
      +ok K6: conditions table
```

`M15-restored`: 2026-10-20T21:22:10Z to 2026-10-20T21:22:10Z, exit 0; restored SHA-256 equals Before and After in the table.

## M16 — Require committed presence

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	if ! git cat-file -e "$base:$ev" 2>/dev/null && ! git cat-file -e "HEAD:$ev" 2>/dev/null; then
-		evidence_error 'is absent from both base and HEAD'
-	fi
+		: # presence validator removed
```

`M16-mutant`: 2026-10-20T21:22:10Z to 2026-10-20T21:22:18Z, exit 1. Named observed failure:

```text
FAIL  bad-usage-evidence-missing
      run: bad-usage-evidence-missing: exit 0 (wanted 2); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.wDbFmZudOI/expect	2026-10-21 00:22:15.608346717 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.wDbFmZudOI/got	2026-10-21 00:22:15.607346717 +0300
      @@ -1,2 +1,7 @@
      -2
      -freeze-check: --evidence runs/T-0000/missing.md: is absent from both base and HEAD
      +0
      +ok K1: white space
      +ok K2: registration of new tests
      +ok K3: test selectors
      +ok K4: restore evidence of mutants
      +ok K5: projection: figure 0/0, maximum 100/10, tolerance 60/4, limit 160/14, margin 160/14 (the figure at landing is still measured at landing)
      +ok K6: conditions table
```

`M16-restored`: 2026-10-20T21:22:18Z to 2026-10-20T21:22:18Z, exit 0; restored SHA-256 equals Before and After in the table.

## M17 — Match declared paths literally

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-		git diff --no-renames --numstat "$base" HEAD -- ":(literal)$ev" >> "$tmp/k5_evidence"
+		git diff --no-renames --numstat "$base" HEAD -- ":(glob)$ev*" >> "$tmp/k5_evidence"
```

`M17-mutant`: 2026-10-20T21:22:18Z to 2026-10-20T21:22:27Z, exit 1. Named observed failure:

```text
FAIL  good-K5-evidence-literal
      run: good-K5-evidence-literal: exit 0 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.vZpowpRKCX/expect	2026-10-21 00:22:26.218345954 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.vZpowpRKCX/got	2026-10-21 00:22:26.217345954 +0300
      @@ -3,5 +3,5 @@
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 3/2, maximum 100/10, tolerance 60/4, limit 160/14, margin 157/12 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 1/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 159/13 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M17-restored`: 2026-10-20T21:22:27Z to 2026-10-20T21:22:27Z, exit 0; restored SHA-256 equals Before and After in the table.

## M18 — Deduplicate declarations

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	grep -Fxq -- "$ev" "$tmp/evidence" || printf '%s\n' "$ev" >> "$tmp/evidence"
+	printf '%s\n' "$ev" >> "$tmp/evidence"
```

`M18-mutant`: 2026-10-20T21:22:27Z to 2026-10-20T21:22:35Z, exit 1. Named observed failure:

```text
FAIL  good-K5-evidence-duplicate
      run: good-K5-evidence-duplicate: exit 0 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.xQCwjq3ufn/expect	2026-10-21 00:22:34.436345560 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.xQCwjq3ufn/got	2026-10-21 00:22:34.436176651 +0300
      @@ -3,5 +3,5 @@
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 1/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 159/13 (the figure at landing is still measured at landing)
      +ok K5: projection: figure -2/0, maximum 100/10, tolerance 60/4, limit 160/14, margin 162/14 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M18-restored`: 2026-10-20T21:22:35Z to 2026-10-20T21:22:36Z, exit 0; restored SHA-256 equals Before and After in the table.

## M19 — Line tolerance floor

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	[ "$tol_l" -ge 60 ] || tol_l=60
+	: # floor removed
```

`M19-mutant`: 2026-10-20T21:22:36Z to 2026-10-20T21:22:44Z, exit 1. Named observed failure:

```text
FAIL  good-K5-line-limit
      run: good-K5-line-limit: exit 1 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.qipRMrxM1X/expect	2026-10-21 00:22:43.976875867 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.qipRMrxM1X/got	2026-10-21 00:22:43.975345314 +0300
      @@ -1,7 +1,7 @@
      -0
      +1
       ok K1: white space
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 160/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 0/13 (the figure at landing is still measured at landing)
      +FAIL K5: -:0: projection: figure 160/1, maximum 100/10, tolerance 20/4, limit 120/14, margin -40/13 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M19-restored`: 2026-10-20T21:22:44Z to 2026-10-20T21:22:44Z, exit 0; restored SHA-256 equals Before and After in the table.

## M20 — Line tolerance percentage

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	tol_l=$((max_l / 5))
+	tol_l=60
```

`M20-mutant`: 2026-10-20T21:22:44Z to 2026-10-20T21:22:53Z, exit 1. Named observed failure:

```text
FAIL  good-K5-percentage
      run: good-K5-percentage: exit 0 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.Yz7yrprGau/expect	2026-10-21 00:22:52.542040354 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.Yz7yrprGau/got	2026-10-21 00:22:52.541100978 +0300
      @@ -3,5 +3,5 @@
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 1/1, maximum 1000/10, tolerance 200/4, limit 1200/14, margin 1199/13 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 1/1, maximum 1000/10, tolerance 60/4, limit 1060/14, margin 1059/13 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M20-restored`: 2026-10-20T21:22:53Z to 2026-10-20T21:22:53Z, exit 0; restored SHA-256 equals Before and After in the table.

## M21 — Line tolerance ceiling

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	[ "$tol_l" -le 500 ] || tol_l=500
+	: # ceiling removed
```

`M21-mutant`: 2026-10-20T21:22:53Z to 2026-10-20T21:23:01Z, exit 1. Named observed failure:

```text
FAIL  good-K5-ceiling
      run: good-K5-ceiling: exit 0 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.JvSC8GeiGk/expect	2026-10-21 00:22:59.779345396 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.JvSC8GeiGk/got	2026-10-21 00:22:59.778345396 +0300
      @@ -3,5 +3,5 @@
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 1/1, maximum 3000/10, tolerance 500/4, limit 3500/14, margin 3499/13 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 1/1, maximum 3000/10, tolerance 600/4, limit 3600/14, margin 3599/13 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M21-restored`: 2026-10-20T21:23:01Z to 2026-10-20T21:23:01Z, exit 0; restored SHA-256 equals Before and After in the table.

## M22 — Whole-line truncation

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	tol_l=$((max_l / 5))
+	tol_l=$(((max_l + 4) / 5))
```

`M22-mutant`: 2026-10-20T21:23:01Z to 2026-10-20T21:23:10Z, exit 1. Named observed failure:

```text
FAIL  good-K5-truncation
      run: good-K5-truncation: exit 0 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.yVKDEzGOVP/expect	2026-10-21 00:23:09.887261266 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.yVKDEzGOVP/got	2026-10-21 00:23:09.886275311 +0300
      @@ -3,5 +3,5 @@
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 1/1, maximum 309/10, tolerance 61/4, limit 370/14, margin 369/13 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 1/1, maximum 309/10, tolerance 62/4, limit 371/14, margin 370/13 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M22-restored`: 2026-10-20T21:23:10Z to 2026-10-20T21:23:10Z, exit 0; restored SHA-256 equals Before and After in the table.

## M23 — Four-file tolerance

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	tol_f=4
+	tol_f=0
```

`M23-mutant`: 2026-10-20T21:23:10Z to 2026-10-20T21:23:19Z, exit 1. Named observed failure:

```text
FAIL  good-K5-file-limit
      run: good-K5-file-limit: exit 1 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.YMeJor5rZO/expect	2026-10-21 00:23:18.192164377 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.YMeJor5rZO/got	2026-10-21 00:23:18.191239437 +0300
      @@ -1,7 +1,7 @@
      -0
      +1
       ok K1: white space
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 5/5, maximum 100/1, tolerance 60/4, limit 160/5, margin 155/0 (the figure at landing is still measured at landing)
      +FAIL K5: -:0: projection: figure 5/5, maximum 100/1, tolerance 60/0, limit 160/1, margin 155/-4 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M23-restored`: 2026-10-20T21:23:19Z to 2026-10-20T21:23:19Z, exit 0; restored SHA-256 equals Before and After in the table.

## M24 — Reject over line limit

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	if [ "$fig_l" -gt "$lim_l" ] || [ "$fig_f" -gt "$lim_f" ]; then
+	if false || [ "$fig_f" -gt "$lim_f" ]; then
```

`M24-mutant`: 2026-10-20T21:23:19Z to 2026-10-20T21:23:27Z, exit 1. Named observed failure:

```text
FAIL  bad-K5-over-lines
      run: bad-K5-over-lines: exit 0 (wanted 1); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.oKIzwdl7cr/expect	2026-10-21 00:23:22.584346556 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.oKIzwdl7cr/got	2026-10-21 00:23:22.583346556 +0300
      @@ -1,7 +1,7 @@
      -1
      +0
       ok K1: white space
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -FAIL K5: -:0: projection: figure 161/1, maximum 100/10, tolerance 60/4, limit 160/14, margin -1/13 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 161/1, maximum 100/10, tolerance 60/4, limit 160/14, margin -1/13 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M24-restored`: 2026-10-20T21:23:27Z to 2026-10-20T21:23:27Z, exit 0; restored SHA-256 equals Before and After in the table.

## M25 — Accept equal line limit

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	if [ "$fig_l" -gt "$lim_l" ] || [ "$fig_f" -gt "$lim_f" ]; then
+	if [ "$fig_l" -ge "$lim_l" ] || [ "$fig_f" -gt "$lim_f" ]; then
```

`M25-mutant`: 2026-10-20T21:23:27Z to 2026-10-20T21:23:36Z, exit 1. Named observed failure:

```text
FAIL  good-K5-line-limit
      run: good-K5-line-limit: exit 1 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.YTmQ1ljMFt/expect	2026-10-21 00:23:35.596073262 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.YTmQ1ljMFt/got	2026-10-21 00:23:35.594347750 +0300
      @@ -1,7 +1,7 @@
      -0
      +1
       ok K1: white space
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 160/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 0/13 (the figure at landing is still measured at landing)
      +FAIL K5: -:0: projection: figure 160/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 0/13 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M25-restored`: 2026-10-20T21:23:36Z to 2026-10-20T21:23:36Z, exit 0; restored SHA-256 equals Before and After in the table.

## M26 — Reject over file limit

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	if [ "$fig_l" -gt "$lim_l" ] || [ "$fig_f" -gt "$lim_f" ]; then
+	if [ "$fig_l" -gt "$lim_l" ] || false; then
```

`M26-mutant`: 2026-10-20T21:23:36Z to 2026-10-20T21:23:44Z, exit 1. Named observed failure:

```text
FAIL  bad-K5-over-files
      run: bad-K5-over-files: exit 0 (wanted 1); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.Z61d4BbNjC/expect	2026-10-21 00:23:39.776258747 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.Z61d4BbNjC/got	2026-10-21 00:23:39.775473441 +0300
      @@ -1,7 +1,7 @@
      -1
      +0
       ok K1: white space
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -FAIL K5: -:0: projection: figure 6/6, maximum 100/1, tolerance 60/4, limit 160/5, margin 154/-1 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 6/6, maximum 100/1, tolerance 60/4, limit 160/5, margin 154/-1 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M26-restored`: 2026-10-20T21:23:44Z to 2026-10-20T21:23:45Z, exit 0; restored SHA-256 equals Before and After in the table.

## M27 — Accept equal file limit

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	if [ "$fig_l" -gt "$lim_l" ] || [ "$fig_f" -gt "$lim_f" ]; then
+	if [ "$fig_l" -gt "$lim_l" ] || [ "$fig_f" -ge "$lim_f" ]; then
```

`M27-mutant`: 2026-10-20T21:23:45Z to 2026-10-20T21:23:53Z, exit 1. Named observed failure:

```text
FAIL  good-K5-file-limit
      run: good-K5-file-limit: exit 1 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.jpHcEQU8Mq/expect	2026-10-21 00:23:52.793161882 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.jpHcEQU8Mq/got	2026-10-21 00:23:52.792225427 +0300
      @@ -1,7 +1,7 @@
      -0
      +1
       ok K1: white space
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 5/5, maximum 100/1, tolerance 60/4, limit 160/5, margin 155/0 (the figure at landing is still measured at landing)
      +FAIL K5: -:0: projection: figure 5/5, maximum 100/1, tolerance 60/4, limit 160/5, margin 155/0 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M27-restored`: 2026-10-20T21:23:53Z to 2026-10-20T21:23:53Z, exit 0; restored SHA-256 equals Before and After in the table.

## M28 — Line margin includes tolerance

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	mar_l=$((lim_l - fig_l))
+	mar_l=$((max_l - fig_l))
```

`M28-mutant`: 2026-10-20T21:23:53Z to 2026-10-20T21:24:02Z, exit 1. Named observed failure:

```text
FAIL  good-K5-line-limit
      run: good-K5-line-limit: exit 0 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.b1DImVTwxh/expect	2026-10-21 00:24:01.764351280 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.b1DImVTwxh/got	2026-10-21 00:24:01.763351280 +0300
      @@ -3,5 +3,5 @@
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 160/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 0/13 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 160/1, maximum 100/10, tolerance 60/4, limit 160/14, margin -60/13 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M28-restored`: 2026-10-20T21:24:02Z to 2026-10-20T21:24:02Z, exit 0; restored SHA-256 equals Before and After in the table.

## M29 — File margin includes tolerance

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	mar_f=$((lim_f - fig_f))
+	mar_f=$((max_f - fig_f))
```

`M29-mutant`: 2026-10-20T21:24:02Z to 2026-10-20T21:24:11Z, exit 1. Named observed failure:

```text
FAIL  good-K5-file-limit
      run: good-K5-file-limit: exit 0 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.9Ap5M6RIZr/expect	2026-10-21 00:24:10.365160648 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.9Ap5M6RIZr/got	2026-10-21 00:24:10.364215499 +0300
      @@ -3,5 +3,5 @@
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 5/5, maximum 100/1, tolerance 60/4, limit 160/5, margin 155/0 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 5/5, maximum 100/1, tolerance 60/4, limit 160/5, margin 155/-4 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M29-restored`: 2026-10-20T21:24:11Z to 2026-10-20T21:24:11Z, exit 0; restored SHA-256 equals Before and After in the table.

## M30 — No synthetic line allowance

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
+	fig_l=$((fig_l + 70))
```

`M30-mutant`: 2026-10-20T21:24:11Z to 2026-10-20T21:24:19Z, exit 1. Named observed failure:

```text
FAIL  good-K5-no-closeout
      run: good-K5-no-closeout: exit 0 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.I5UxAA555y/expect	2026-10-21 00:24:19.358099993 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.I5UxAA555y/got	2026-10-21 00:24:19.357342668 +0300
      @@ -3,5 +3,5 @@
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 6/1, maximum 75/10, tolerance 60/4, limit 135/14, margin 129/13 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 76/1, maximum 75/10, tolerance 60/4, limit 135/14, margin 59/13 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M30-restored`: 2026-10-20T21:24:19Z to 2026-10-20T21:24:20Z, exit 0; restored SHA-256 equals Before and After in the table.

## M31 — No synthetic file allowance

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
+	fig_f=$((fig_f + 2))
```

`M31-mutant`: 2026-10-20T21:24:20Z to 2026-10-20T21:24:28Z, exit 1. Named observed failure:

```text
FAIL  good-K5-no-closeout
      run: good-K5-no-closeout: exit 0 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.aekaUut8e8/expect	2026-10-21 00:24:28.158356301 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.aekaUut8e8/got	2026-10-21 00:24:28.157356301 +0300
      @@ -3,5 +3,5 @@
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 6/1, maximum 75/10, tolerance 60/4, limit 135/14, margin 129/13 (the figure at landing is still measured at landing)
      +ok K5: projection: figure 6/3, maximum 75/10, tolerance 60/4, limit 135/14, margin 129/11 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M31-restored`: 2026-10-20T21:24:28Z to 2026-10-20T21:24:28Z, exit 0; restored SHA-256 equals Before and After in the table.

## M32 — Reject retired closeout option

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
+	--closeout)
+		[ $# -ge 2 ] || usage
+		shift 2 ;;
```

`M32-mutant`: 2026-10-20T21:24:28Z to 2026-10-20T21:24:37Z, exit 1. Named observed failure:

```text
FAIL  bad-usage-closeout
      run: bad-usage-closeout: exit 0 (wanted 2); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.bUXtqnGkcq/expect	2026-10-21 00:24:34.145032057 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.bUXtqnGkcq/got	2026-10-21 00:24:34.144216788 +0300
      @@ -1,2 +1,7 @@
      -2
      -freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--evidence PATH] [--no-conditions]
      +0
      +ok K1: white space
      +ok K2: registration of new tests
      +ok K3: test selectors
      +ok K4: restore evidence of mutants
      +ok K5: projection: figure 0/0, maximum 100/10, tolerance 60/4, limit 160/14, margin 160/14 (the figure at landing is still measured at landing)
      +ok K6: conditions table
```

`M32-restored`: 2026-10-20T21:24:37Z to 2026-10-20T21:24:37Z, exit 0; restored SHA-256 equals Before and After in the table.

## M33 — Require evidence argument

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-		[ $# -ge 2 ] || usage
+		[ $# -ge 2 ] || set -- "$1" runs/T-0000/test-runs.md
```

`M33-mutant`: 2026-10-20T21:24:37Z to 2026-10-20T21:24:45Z, exit 1. Named observed failure:

```text
FAIL  bad-usage-evidence-no-value
      run: bad-usage-evidence-no-value: exit 0 (wanted 2); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.xvUFGhVTaY/expect	2026-10-21 00:24:42.842901851 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.xvUFGhVTaY/got	2026-10-21 00:24:42.842107433 +0300
      @@ -1,2 +1,7 @@
      -2
      -freeze-check: usage: freeze-check.sh --task T-xxxx --base REV --budget LINES/FILES [--evidence PATH] [--no-conditions]
      +0
      +ok K1: white space
      +ok K2: registration of new tests
      +ok K3: test selectors
      +ok K4: restore evidence of mutants
      +ok K5: projection: figure 0/0, maximum 100/10, tolerance 60/4, limit 160/14, margin 160/14 (the figure at landing is still measured at landing)
      +ok K6: conditions table
```

`M33-restored`: 2026-10-20T21:24:45Z to 2026-10-20T21:24:45Z, exit 0; restored SHA-256 equals Before and After in the table.

## M34 — Canonical margin assertion

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-  The margin is the maximum plus the tolerance, minus the landing figure, calculated separately in lines and in files.
+  Margin definition removed.
```

`M34-mutant`: 2026-10-20T21:24:45Z to 2026-10-20T21:24:45Z, exit 1. Named observed failure:

```text
FAIL: canonical margin sentence in the tolerance section
```

`M34-restored`: 2026-10-20T21:24:45Z to 2026-10-20T21:24:46Z, exit 0; restored SHA-256 equals Before and After in the table.

## M35 — Preserve run.sh

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
+# altered runner copy
```

`M35-mutant`: 2026-10-20T21:24:46Z to 2026-10-20T21:24:46Z, exit 1. Named observed failure:

```text
FAIL: unchanged runner: docs/tests/freeze-check/run.sh
```

`M35-restored`: 2026-10-20T21:24:46Z to 2026-10-20T21:24:46Z, exit 0; restored SHA-256 equals Before and After in the table.

## M36 — Preserve run-all.sh

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
+# altered runner copy
```

`M36-mutant`: 2026-10-20T21:24:46Z to 2026-10-20T21:24:46Z, exit 1. Named observed failure:

```text
FAIL: unchanged runner: docs/tests/freeze-check/run-all.sh
```

`M36-restored`: 2026-10-20T21:24:46Z to 2026-10-20T21:24:46Z, exit 0; restored SHA-256 equals Before and After in the table.

## M51 — Line limit adds tolerance

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	lim_l=$((max_l + tol_l))
+	lim_l=$max_l
```

`M51-mutant`: 2026-10-20T21:24:41Z to 2026-10-20T21:24:49Z, exit 1. Named observed failure:

```text
FAIL  good-K5-line-limit
      run: good-K5-line-limit: exit 1 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.ABG81qzJGC/expect	2026-10-21 00:24:48.924088445 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.ABG81qzJGC/got	2026-10-21 00:24:48.923205345 +0300
      @@ -1,7 +1,7 @@
      -0
      +1
       ok K1: white space
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 160/1, maximum 100/10, tolerance 60/4, limit 160/14, margin 0/13 (the figure at landing is still measured at landing)
      +FAIL K5: -:0: projection: figure 160/1, maximum 100/10, tolerance 60/4, limit 100/14, margin -60/13 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M51-restored`: 2026-10-20T21:24:49Z to 2026-10-20T21:24:49Z, exit 0; restored SHA-256 equals Before and After in the table.

## M52 — File limit adds tolerance

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-	lim_f=$((max_f + tol_f))
+	lim_f=$max_f
```

`M52-mutant`: 2026-10-20T21:24:49Z to 2026-10-20T21:24:58Z, exit 1. Named observed failure:

```text
FAIL  good-K5-file-limit
      run: good-K5-file-limit: exit 1 (wanted 0); diff (-EXPECT +got):
      --- /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.BWjYfPIT1b/expect	2026-10-21 00:24:57.610114617 +0300
      +++ /home/layup/.cache/layup-session/tmp-green/T-c826-stage-c/tmp.BWjYfPIT1b/got	2026-10-21 00:24:57.609327060 +0300
      @@ -1,7 +1,7 @@
      -0
      +1
       ok K1: white space
       ok K2: registration of new tests
       ok K3: test selectors
       ok K4: restore evidence of mutants
      -ok K5: projection: figure 5/5, maximum 100/1, tolerance 60/4, limit 160/5, margin 155/0 (the figure at landing is still measured at landing)
      +FAIL K5: -:0: projection: figure 5/5, maximum 100/1, tolerance 60/4, limit 160/1, margin 155/-4 (the figure at landing is still measured at landing)
       ok K6: conditions table
```

`M52-restored`: 2026-10-20T21:24:58Z to 2026-10-20T21:24:58Z, exit 0; restored SHA-256 equals Before and After in the table.

## M53 — Retain restore-list path docs/adr/adr-lint.sh

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-          for f in docs/adr/adr-lint.sh docs/prd/prd-lint.sh \
+          for f in docs/prd/prd-lint.sh \
```

`M53-mutant`: 2026-10-20T21:25:15Z to 2026-10-20T21:25:15Z, exit 1. Named observed failure:

```text
FAIL: all nine baseline paths stay in the restore list
```

`M53-restored`: 2026-10-20T21:25:15Z to 2026-10-20T21:25:15Z, exit 0; restored SHA-256 equals Before and After in the table.

## M54 — Retain restore-list path docs/tests/freeze-check/run.sh

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-                   docs/tests/freeze-check/run.sh \
+                   \
```

`M54-mutant`: 2026-10-20T21:25:15Z to 2026-10-20T21:25:15Z, exit 1. Named observed failure:

```text
FAIL: all nine baseline paths stay in the restore list
```

`M54-restored`: 2026-10-20T21:25:15Z to 2026-10-20T21:25:15Z, exit 0; restored SHA-256 equals Before and After in the table.

## M55 — Retain restore-list path docs/tests/freeze-check/run-all.sh

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-                   docs/tests/freeze-check/run-all.sh; do
+                   ; do
```

`M55-mutant`: 2026-10-20T21:25:15Z to 2026-10-20T21:25:15Z, exit 1. Named observed failure:

```text
FAIL: all nine baseline paths stay in the restore list
```

`M55-restored`: 2026-10-20T21:25:15Z to 2026-10-20T21:25:15Z, exit 0; restored SHA-256 equals Before and After in the table.

## M56 — Require pull_request guard

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-            if [ "$EVENT_NAME" = pull_request ] && [ "$PR_NUMBER" = 453 ] && [ "$f" = docs/tests/freeze-check.sh ]; then
+            if [ "$PR_NUMBER" = 453 ] && [ "$f" = docs/tests/freeze-check.sh ]; then
```

`M56-mutant`: 2026-10-20T21:25:15Z to 2026-10-20T21:25:15Z, exit 1. Named observed failure:

```text
FAIL: kept checker requires the pull_request event
FAIL: one literal PR number and only the checker path can be kept
FAIL: only the guarded checker skips restoration; both runners and other paths restore for every event
```

`M56-restored`: 2026-10-20T21:25:15Z to 2026-10-20T21:25:15Z, exit 0; restored SHA-256 equals Before and After in the table.

## M57 — Require PR-number guard

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-            if [ "$EVENT_NAME" = pull_request ] && [ "$PR_NUMBER" = 453 ] && [ "$f" = docs/tests/freeze-check.sh ]; then
+            if [ "$EVENT_NAME" = pull_request ] && [ "$f" = docs/tests/freeze-check.sh ]; then
```

`M57-mutant`: 2026-10-20T21:25:15Z to 2026-10-20T21:25:15Z, exit 1. Named observed failure:

```text
FAIL: one literal PR number and only the checker path can be kept
FAIL: exception matches the supplied actual migration PR
FAIL: only the guarded checker skips restoration; both runners and other paths restore for every event
FAIL: exception comment names the task, literal PR and removal issue
```

`M57-restored`: 2026-10-20T21:25:15Z to 2026-10-20T21:25:15Z, exit 0; restored SHA-256 equals Before and After in the table.

## M58 — Keep exactly one checker path

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-            if [ "$EVENT_NAME" = pull_request ] && [ "$PR_NUMBER" = 453 ] && [ "$f" = docs/tests/freeze-check.sh ]; then
+            if [ "$EVENT_NAME" = pull_request ] && [ "$PR_NUMBER" = 453 ] && true; then
```

`M58-mutant`: 2026-10-20T21:25:15Z to 2026-10-20T21:25:15Z, exit 1. Named observed failure:

```text
FAIL: one literal PR number and only the checker path can be kept
FAIL: only the guarded checker skips restoration; both runners and other paths restore for every event
```

`M58-restored`: 2026-10-20T21:25:15Z to 2026-10-20T21:25:16Z, exit 0; restored SHA-256 equals Before and After in the table.

## M59 — Match actual approved number

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-            if [ "$EVENT_NAME" = pull_request ] && [ "$PR_NUMBER" = 453 ] && [ "$f" = docs/tests/freeze-check.sh ]; then
+            if [ "$EVENT_NAME" = pull_request ] && [ "$PR_NUMBER" = 454 ] && [ "$f" = docs/tests/freeze-check.sh ]; then
```

`M59-mutant`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 1. Named observed failure:

```text
FAIL: exception matches the supplied actual migration PR
FAIL: exception comment names the task, literal PR and removal issue
```

`M59-restored`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 0; restored SHA-256 equals Before and After in the table.

## M60 — Read actual event

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-          EVENT_NAME: ${{ github.event_name }}
+          EVENT_NAME: pull_request
```

`M60-mutant`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 1. Named observed failure:

```text
FAIL: event guard receives the actual GitHub event name
```

`M60-restored`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 0; restored SHA-256 equals Before and After in the table.

## M61 — Read actual PR number

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-          PR_NUMBER: ${{ github.event.pull_request.number }}
+          PR_NUMBER: 453
```

`M61-mutant`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 1. Named observed failure:

```text
FAIL: number guard receives the actual pull request number
```

`M61-restored`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 0; restored SHA-256 equals Before and After in the table.

## M62 — Restore run.sh even if listed

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
+            [ "$f" = docs/tests/freeze-check/run.sh ] && continue
```

`M62-mutant`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 1. Named observed failure:

```text
FAIL: only the guarded checker skips restoration; both runners and other paths restore for every event
```

`M62-restored`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 0; restored SHA-256 equals Before and After in the table.

## M63 — Restore run-all.sh even if listed

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
+            [ "$f" = docs/tests/freeze-check/run-all.sh ] && continue
```

`M63-mutant`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 1. Named observed failure:

```text
FAIL: only the guarded checker skips restoration; both runners and other paths restore for every event
```

`M63-restored`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 0; restored SHA-256 equals Before and After in the table.

## M64 — Use default-branch source

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-          git clone --quiet --depth 1 --branch "$DEFAULT_BRANCH" \
+          git clone --quiet --depth 1 --branch T-c826 \
```

`M64-mutant`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 1. Named observed failure:

```text
FAIL: restoration still reads the default branch of this repository
```

`M64-restored`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 0; restored SHA-256 equals Before and After in the table.

## M65 — Read repository default branch

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-          DEFAULT_BRANCH: ${{ github.event.repository.default_branch }}
+          DEFAULT_BRANCH: T-c826
```

`M65-mutant`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 1. Named observed failure:

```text
FAIL: default branch still comes from the repository event
```

`M65-restored`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 0; restored SHA-256 equals Before and After in the table.

## M66 — Name removal issue

Single edit (removed lines start with `-`, replacement lines with `+`):

```diff
-        # R4: T-c826 (#449), PR 453; remove in T-szvp (#452).
+        # R4: T-c826 (#449), PR 453; removal reference removed.
```

`M66-mutant`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 1. Named observed failure:

```text
FAIL: exception comment names the task, literal PR and removal issue
```

`M66-restored`: 2026-10-20T21:25:16Z to 2026-10-20T21:25:16Z, exit 0; restored SHA-256 equals Before and After in the table.

## Stage G — round 1 fixes (M67–M70)

Four additional variants are killed (56 variants total across both sweeps).
M69 tests the existing rename rule on the unmodified checker at frozen head
`fdb8d626008762b2e13cbab5fe53c4d5d3226700`, after the new cases were written
and before the Stage G code. M67, M68 and M70 test the fixed checker.
Each variant copies the current test tree to its own fresh scratch directory,
changes only the named checker statement, passes `sh -n`, and runs the unchanged
case harness. Harness exit 2 is a full-output/status mismatch. Each fresh backup
is restored, its SHA-256 must equal Before and the source checker, and the same
case then exits 0. The commands and direct UTC timestamps are in
[test-runs.md](test-runs.md#stage-g--round-1-fixes-cycle-1-of-cap-2).

| ID | Rule changed alone | Named case | Before SHA-256 | During SHA-256 | After SHA-256 | Verdict |
| --- | --- | --- | --- | --- | --- | --- |
| M67 | Skip HEAD object-kind validation | `bad-usage-evidence-directory-head` | d316f5e82ae1a6ae951915bfa926e75e6a8dd0d209286592fd090f456bc1abcc | cd55b52259f4738ca2a7b94307cca5b832e6c1810e4bd318acfcf83abd8f0c61 | d316f5e82ae1a6ae951915bfa926e75e6a8dd0d209286592fd090f456bc1abcc | KILLED; mutant exit 2, restored exit 0 |
| M68 | Skip base object-kind validation | `bad-usage-evidence-directory-base` | d316f5e82ae1a6ae951915bfa926e75e6a8dd0d209286592fd090f456bc1abcc | f94e2170c487aa94f87329260c969eb4354e8144bb7690b25ea531dcbdb3ecf1 | d316f5e82ae1a6ae951915bfa926e75e6a8dd0d209286592fd090f456bc1abcc | KILLED; mutant exit 2, restored exit 0 |
| M69 | Remove --no-renames from the gross diff only | `good-K5-evidence-renamed` | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | 8cf1e3e0a599ccbae6863c86468c55ecaddd41c3c8590bab96b48960cbc64eb6 | 04c5918711f4f07e5fdd423424cbf14b99a43aefd77e76603c183b98cc96e2cf | KILLED; mutant exit 2, restored exit 0 |
| M70 | Restore the wrong trailing-slash diagnostic | `bad-usage-evidence-trailing-slash` | d316f5e82ae1a6ae951915bfa926e75e6a8dd0d209286592fd090f456bc1abcc | 6242c5288d6c732426f87adecb84a97ca432baf5143ed796979be6ed5f5853b9 | d316f5e82ae1a6ae951915bfa926e75e6a8dd0d209286592fd090f456bc1abcc | KILLED; mutant exit 2, restored exit 0 |

### M67 — Skip HEAD object-kind validation

Single edit to `docs/tests/freeze-check.sh`:

```diff
-for ev_rev in "$base" HEAD; do
+for ev_rev in "$base"; do
```

Observed mismatch, 2026-10-20T22:23:19Z to 2026-10-20T22:23:20Z, harness exit 2:

```text
run: bad-usage-evidence-directory-head: exit 0 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/tmp.SakuOlLBZj/expect	2026-10-21 01:23:20.021930410 +0300
+++ /home/layup/.cache/layup-session/tmp-green/tmp.SakuOlLBZj/got	2026-10-21 01:23:20.021071857 +0300
@@ -1,2 +1,7 @@
-2
-freeze-check: --evidence runs/T-0000/d.md: must name a file at HEAD
+0
+ok K1: white space
+ok K2: registration of new tests
+ok K3: test selectors
+ok K4: restore evidence of mutants
+ok K5: projection: figure 0/0, maximum 0/0, tolerance 60/4, limit 60/4, margin 60/4 (the figure at landing is still measured at landing)
+ok K6: conditions table
```

Restored control: 2026-10-20T22:23:20Z to 2026-10-20T22:23:20Z, exit 0; Before = After, as measured above.

### M68 — Skip base object-kind validation

Single edit to `docs/tests/freeze-check.sh`:

```diff
-for ev_rev in "$base" HEAD; do
+for ev_rev in HEAD; do
```

Observed mismatch, 2026-10-20T22:23:20Z to 2026-10-20T22:23:20Z, harness exit 2:

```text
run: bad-usage-evidence-directory-base: exit 0 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/tmp.v0dbA5jHsL/expect	2026-10-21 01:23:20.427832036 +0300
+++ /home/layup/.cache/layup-session/tmp-green/tmp.v0dbA5jHsL/got	2026-10-21 01:23:20.427078467 +0300
@@ -1,2 +1,7 @@
-2
-freeze-check: --evidence runs/T-0000/base-tree.md: must name a file at base
+0
+ok K1: white space
+ok K2: registration of new tests
+ok K3: test selectors
+ok K4: restore evidence of mutants
+ok K5: projection: figure 0/0, maximum 0/0, tolerance 60/4, limit 60/4, margin 60/4 (the figure at landing is still measured at landing)
+ok K6: conditions table
```

Restored control: 2026-10-20T22:23:20Z to 2026-10-20T22:23:20Z, exit 0; Before = After, as measured above.

### M69 — Remove --no-renames from the gross diff only

Single edit to `docs/tests/freeze-check.sh`:

```diff
-git diff --no-renames --numstat "$base" HEAD > "$tmp/k5_gross"
+git diff --numstat "$base" HEAD > "$tmp/k5_gross"
```

Observed mismatch, 2026-10-20T22:22:22Z to 2026-10-20T22:22:23Z, harness exit 2:

```text
run: good-K5-evidence-renamed: exit 0 (wanted 0); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/tmp.FXznD9Hpp7/expect	2026-10-21 01:22:23.041194126 +0300
+++ /home/layup/.cache/layup-session/tmp-green/tmp.FXznD9Hpp7/got	2026-10-21 01:22:23.039903950 +0300
@@ -3,5 +3,5 @@
 ok K2: registration of new tests
 ok K3: test selectors
 ok K4: restore evidence of mutants
-ok K5: projection: figure 3/1, maximum 0/0, tolerance 60/4, limit 60/4, margin 57/3 (the figure at landing is still measured at landing)
+ok K5: projection: figure -3/0, maximum 0/0, tolerance 60/4, limit 60/4, margin 63/4 (the figure at landing is still measured at landing)
 ok K6: conditions table
```

Restored control: 2026-10-20T22:22:23Z to 2026-10-20T22:22:23Z, exit 0; Before = After, as measured above.

### M70 — Restore the wrong trailing-slash diagnostic

Single edit to `docs/tests/freeze-check.sh`:

```diff
-evidence_error 'empty path components are not allowed'
+evidence_error 'path traversal is not allowed'
```

Observed mismatch, 2026-10-20T22:23:20Z to 2026-10-20T22:23:20Z, harness exit 2:

```text
run: bad-usage-evidence-trailing-slash: exit 2 (wanted 2); diff (-EXPECT +got):
--- /home/layup/.cache/layup-session/tmp-green/tmp.uuPi8a3kpe/expect	2026-10-21 01:23:20.769054041 +0300
+++ /home/layup/.cache/layup-session/tmp-green/tmp.uuPi8a3kpe/got	2026-10-21 01:23:20.768054038 +0300
@@ -1,2 +1,2 @@
 2
-freeze-check: --evidence runs/T-0000/: empty path components are not allowed
+freeze-check: --evidence runs/T-0000/: path traversal is not allowed
```

Restored control: 2026-10-20T22:23:20Z to 2026-10-20T22:23:20Z, exit 0; Before = After, as measured above.
