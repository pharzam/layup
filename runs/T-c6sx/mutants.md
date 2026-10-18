# T-c6sx: the compiled mutants

Each row is from a real run of `mutants.py` (kept outside the tree, with its output) on a copy of the committed
green head `b98ccc9` made by `git archive`: one change to one file, then `go build ./...`, `go vet ./...`,
and `go test -count=1 ./internal/run`; when the unit level passes, the named integration tests of this task
(`TMPDIR=/home/layup/.cache/layup-session/tmp-green go test -count=1 -tags=integration -timeout 20m ./internal/run -run 'TestEachStepsPromptFileFromARealRecordsTree|TestPromptFileRefusalsHaveTheirControls|TestAHarnessThatFollowsTheInstructionGivesValidHandoffs|TestTheInstructionBlocksOfLoopMd|TestTheSpecimenOfTheSpecificationInstructionIsAcceptedByThePinnedPrdLint'`) run too, and the column test says which level
the exit status is of. A mutant that does not build or vet counts as not run. The hashes are of the whole file, the
original and the mutated one. The restore is checked by hash against the head's file and followed by a green control
(below the table). Every mutant that built was killed by its first failing assertion.

M13 of the plan (the base guard of `TaskSession`) is not in this task: it moved to #415 (row 66r). Two rows
were run again after a fault of the mutant itself, not of the code: the first M7 removed the only use of a variable
(build failure, not run), and the first M22 indexed an empty ledger (a panic in every other test); the rows below are
the second runs. The plan's M15 (a label hard-coded) is replaced by a mutant that hard-codes today's list and is
killed by the test that gives `labelList` another map (condition 7).

| # | build | vet | test | failing assertion (first line) | sha256 before | sha256 after |
| - | ----- | --- | ---- | ------------------------------ | ------------- | ------------ |
| M1 | 0 | 0 | 1 (unit) | instructions_test.go:87: Text differs from testdata/prompts/plan-of-a-task.txt: got 5826 bytes, want 4649; first difference at byte 59: got "\nThe work of this step:\nPlan the task {t", want "Read each file at its path: its SHA-25 | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | edd0ed9cb87188f9eef4dd6f178cffde8cf028f1db21b321c31c289881b40c8a |
| M2 | 0 | 0 | 1 (unit) | instructions_test.go:87: Text differs from testdata/prompts/plan-of-a-task.txt: got 4635 bytes, want 4649; first difference at byte 3142: got "result/result.tsv: UTF-8 text with no by", want "/s/S-aaaaaaaa/result/result.tsv: UTF-8 | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | 4271663c7f35f87abb4a37d647997c49f4b39c0135515fb424375c0b7f63d674 |
| M3 | 0 | 0 | 1 (unit) | instructions_test.go:87: Text differs from testdata/prompts/verification.txt: got 4622 bytes, want 4622; first difference at byte 1633: got "fedcba9876543210fedcba9876543210fedcba98", want "0123456789abcdef0123456789abcdef01234567 | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | 547b9327ae0cc44777fc5076b31b150eb35c5956d4837b56f181130977713977 |
| M4 | 0 | 0 | 1 (unit) | instructions_test.go:87: Text differs from testdata/prompts/verification.txt: got 4604 bytes, want 4622; first difference at byte 2049: got "plan.md, in this form, each line startin", want "review-record.md, in this form, each lin | 894e4d4ddf8cf9d979e8b57dffce20ceb7cf26b2ecaeec4503ad8d041dc199cb | ca0a23a4420a5f2b37c097f3fea968e95ba8588160a256849306b78629f31d63 |
| M5 | 0 | 0 | 1 (unit) | instructions_test.go:87: Text differs from testdata/prompts/plan-of-a-task.txt: got 4248 bytes, want 4649; first difference at byte 750: got "\ndecision-records:\n- docs/adr/0019-phase", want "- /s/S-aaaaaaaa/input/draft.tsv 7eb2c | 894e4d4ddf8cf9d979e8b57dffce20ceb7cf26b2ecaeec4503ad8d041dc199cb | 8ef5fd163fe0b053d90379613da408e326e4653e2f49b835a7e136f5c18f4cb2 |
| M6 | 0 | 0 | 1 (unit) | instructions_test.go:87: Text differs from testdata/prompts/plan-review.txt: got 4604 bytes, want 4649; first difference at byte 1539: got "\nhandoffs-and-answers:\n- /s/S-aaaaaaaa/i", want "the survivors of the filter:\n- db: pg  | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | bbeeaa27bfae3a93ecdf156121385cba90b4a566ef4ac6227d4abff2d970643f |
| M7 | 0 | 0 | 1 (unit) | instructions_test.go:380: a survivor outside the table: <nil>; want an error holding "decision db: the survivor mongo is not an option of the table" | 894e4d4ddf8cf9d979e8b57dffce20ceb7cf26b2ecaeec4503ad8d041dc199cb | b2b7cb5ec9ac37dcb5f4d11a4e5420ca04993c61ef22639c435d0852796c4c5b |
| M8 | 0 | 0 | 1 (unit) | instructions_test.go:380: no survivor: <nil>; want an error holding "decision db: the result has no survivor" | 894e4d4ddf8cf9d979e8b57dffce20ceb7cf26b2ecaeec4503ad8d041dc199cb | 4d76a15c6e4eea62afc87895e12956d81f4f1583ebbc415cd873a56ef4dd2f03 |
| M9 | 0 | 0 | 1 (integration) | promptfile_integration_test.go:367: a draft with a bad kind: <nil>; want an error holding "spec/draft.tsv: " | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | b909d540b262a4591935a2774315e1ca612828b065f560d460c2190cd5acc4a1 |
| M10 | 0 | 0 | 1 (integration) | promptfile_integration_test.go:283: plan-of-a-task: the records have no tasks/T-sp01/tests.tsv | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | 9e9b2c861e96b34d75d181671240b2820eb983304487437ede7f553a97cbcf9c |
| M11 | 0 | 0 | 1 (unit) | instructions_test.go:253: checkStepKind(specification, architecture) = <nil>; want ok false | 894e4d4ddf8cf9d979e8b57dffce20ceb7cf26b2ecaeec4503ad8d041dc199cb | eaa57892dac3ad3d90fe6adc2d5cbe0f344a487e91b0c045fb7b713815723862 |
| M12 | 0 | 0 | 1 (integration) | promptfile_integration_test.go:346: the step shape-activation-batch: git ls-tree -r -z --full-tree --end-of-options 0000000000000000000000000000000000000000 -- .: exit status 128: fatal: not a tree object | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | 693053e96cb6629a995faa6a914b02755489d2e85336c74e55e575eec20e1ddd |
| M14 | 0 | 0 | 1 (unit) | instructions_test.go:87: Text differs from testdata/prompts/close-out.txt: got 3454 bytes, want 3655; first difference at byte 1263: got "findings of the last attempt:\n- gate:sta", want "- /s/S-aaaaaaaa/input/sessions.tsv e63bc" | 894e4d4ddf8cf9d979e8b57dffce20ceb7cf26b2ecaeec4503ad8d041dc199cb | bd81e332c0d227e08bfdc8a364f2016b0926a13c9d5f198448492d925d95c7db |
| M15 | 0 | 0 | 1 (unit) | instructions_test.go:117: labelList = "architecture-boundary, domain-or-business, environment-or-infrastructure, interface-contract"; want the sorted keys, joined by a comma and a space | 894e4d4ddf8cf9d979e8b57dffce20ceb7cf26b2ecaeec4503ad8d041dc199cb | 7899bae2c832f154c01609ba4a5a25c7785e3986b3fb26030a535a6d1406dea0 |
| M16 | 0 | 0 | 1 (integration) | promptfile_integration_test.go:362: the missing spec/draft.tsv: the records have no spec/spans.tsv; want an error holding "the records have no spec/draft.tsv" | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | bec642f3f24d1b4169045be85397bf3c0ff36f31efaed4ca5e83825d920be905 |
| M17 | 0 | 0 | 1 (unit) | instructions_test.go:87: Text differs from testdata/prompts/plan-of-a-task.txt: got 4649 bytes, want 4649; first difference at byte 2435: got " - <id> <path>\nThe word after size: is e", want "\n- <id> <path>\nThe word after size: | 894e4d4ddf8cf9d979e8b57dffce20ceb7cf26b2ecaeec4503ad8d041dc199cb | cea300e3cde30fbd46a8d63584f9ca4e4d81358959fb6024715ec10c22878a77 |
| M18 | 0 | 0 | 1 (unit) | instructions_test.go:462: attemptBase = "aaaa", true; want the base of attempt 2 | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | c7620d06a5a88efc2b8734316a286f9194a5caaca678f76855e11fbcce2daf04 |
| M19 | 0 | 0 | 1 (unit) | instructions_test.go:87: Text differs from testdata/prompts/plan-of-a-task.txt: got 4649 bytes, want 4649; first difference at byte 295: got "ests.tsv 11d5bbdbcdf1fa5aff71bba2bd6f088", want "asks.tsv 299e5653428469c74c66c8063cb135 | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | fdc224c6caf24599e259a3e1d9133496877087fbec84dfb47baf21b3c16f6465 |
| M20 | 0 | 0 | 1 (unit) | instructions_test.go:380: a decision that is no path word: the records have no spec/options/../db.tsv; want an error holding "decision ../db" | 894e4d4ddf8cf9d979e8b57dffce20ceb7cf26b2ecaeec4503ad8d041dc199cb | b30bdda74c4561e9b3ce33e0bc8120202934e371aed7b631ac6aac4fdb691fa3 |
| M21 | 0 | 0 | 1 (unit) | instructions_test.go:87: Text differs from testdata/prompts/specification.txt: got 5757 bytes, want 5757; first difference at byte 2819: got "3#<n> F-0003 \|\n\n\| Requirement \| Acceptan", want "4#<n> F-0003 \|\n\n\| Requirement  | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | c53dcfd3906838edaa1f2b97e4d9179d1689d073941faa4e0b6c7b7134f17c09 |
| M22 | 0 | 0 | 1 (unit) | instructions_test.go:87: Text differs from testdata/prompts/close-out.txt: got 3655 bytes, want 3655; first difference at byte 1285: got "telemetry.tsv 4692305d77aea295bf6112ac97", want "sessions.tsv e63bc74763f411275824ff25e0a" | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | 6a162d221bd24d4aa86b18e722d31da540a2d41b55779643367419806fd805c4 |


## Round 1: the mutants of the required inputs and of the citation contract

Same method, same script. First the reader-bypass and missing-record mutants of finding 3 on the head `cbb40b0`
(the tests before the round), then all of them and the new rules (M30 to M35) on the head `2b01ec7`.
M25 and M29 were already killed by the old tests; M23, M24, M26, M27 and M28 passed every old test.

On `cbb40b0`:

| # | build | vet | test | failing assertion (first line) | sha256 before | sha256 after |
| - | ----- | --- | ---- | ------------------------------ | ------------- | ------------ |
| M23 | 0 | 0 | 0 (integration) | SURVIVED: all tests pass | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | ccf8b15a58e885005554269e969a126af9b9c9a0e8bd6b771f1aa4636f138bfd |
| M24 | 0 | 0 | 0 (integration) | SURVIVED: all tests pass | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | 984fc9a7c6cd140e4f529a09cdbccb4033bb56e1e5c8291bd86d36750bfe56bb |
| M25 | 0 | 0 | 1 (integration) | promptfile_integration_test.go:392: close-out with no ledger: the records have no telemetry.tsv; want an error holding "the records have no sessions.tsv" | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | dfaa6e38899f6cbbe7e43ca75bd0cc020c30f38eabe7dd9c47c8103d88e3c179 |
| M26 | 0 | 0 | 0 (integration) | SURVIVED: all tests pass | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | 8c4e5ef8f0e99c6b839dc8d0790a0555dbb2ffe88f70f668fe23883fbf12f9ba |
| M27 | 0 | 0 | 0 (integration) | SURVIVED: all tests pass | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | 2d25799a1b02ec0e2d370e6c775ef59cf51ac9eba6af2b57368368806aa0d978 |
| M28 | 0 | 0 | 0 (integration) | SURVIVED: all tests pass | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | c66f967028ed69a3454badd681cf15c3f388f6be71af5837249d3ca90f8a0f5e |
| M29 | 0 | 0 | 1 (integration) | promptfile_integration_test.go:401: a refused telemetry table: <nil>; want an error holding "telemetry.tsv: " | ad675b51b351c8962d50f5ab7fe2237d5775adf5c8f7afc990cc59bfc3704e69 | 64f02121d5472e550b3c0d999124bfdb8b1535f042a4b117fbb19dad39272176 |

On `2b01ec7`:

| # | build | vet | test | failing assertion (first line) | sha256 before | sha256 after |
| - | ----- | --- | ---- | ------------------------------ | ------------- | ------------ |
| M23 | 0 | 0 | 1 (integration) | promptfile_integration_test.go:374: the removed start/problem-statement.md: <nil>; want an error holding "the records have no start/problem-statement.md" | 0a524319e09f21f0293947a02b4bfcee9a87aeaf0bc8e6c49abb42caad9fce7b | 618941576430ee5d6fca240d7eb3b8905f8b756f865bf0e8a8181b77985e371c |
| M24 | 0 | 0 | 1 (integration) | promptfile_integration_test.go:422: close-out without telemetry.tsv: <nil>; want an error holding "the records have no telemetry.tsv" | 0a524319e09f21f0293947a02b4bfcee9a87aeaf0bc8e6c49abb42caad9fce7b | bb96f77deb6d399028c705c09c1cf8a30a27e1184284781687200a7a62a76b2b |
| M25 | 0 | 0 | 1 (integration) | promptfile_integration_test.go:409: close-out with no ledger: the records have no telemetry.tsv; want an error holding "the records have no sessions.tsv" | 0a524319e09f21f0293947a02b4bfcee9a87aeaf0bc8e6c49abb42caad9fce7b | f30a610ad2649ee6808d2d83001d0e2d52812f900927375651416bf8c31775ad |
| M26 | 0 | 0 | 1 (integration) | promptfile_integration_test.go:379: a refused spec/spans.tsv: <nil>; want an error holding "spec/spans.tsv: " | 0a524319e09f21f0293947a02b4bfcee9a87aeaf0bc8e6c49abb42caad9fce7b | b276762aa7848fc5ebd1f760314a69d297384ab74e4045f72ca252f255190495 |
| M27 | 0 | 0 | 1 (integration) | promptfile_integration_test.go:379: a refused intake/answers.tsv: <nil>; want an error holding "intake/answers.tsv: " | 0a524319e09f21f0293947a02b4bfcee9a87aeaf0bc8e6c49abb42caad9fce7b | e4fae12a6b51ba6391bef0764339fbaed35146f494b51430e3c8c23ea2f0f3d7 |
| M28 | 0 | 0 | 0 (integration) | SURVIVED: all tests pass | 0a524319e09f21f0293947a02b4bfcee9a87aeaf0bc8e6c49abb42caad9fce7b | 904680b2da75a47269b8f4dc2cf8a111bd2e2ae228b764744df8d48a8f9ed099 |
| M29 | 0 | 0 | 1 (integration) | promptfile_integration_test.go:427: a refused telemetry table: <nil>; want an error holding "telemetry.tsv: " | 0a524319e09f21f0293947a02b4bfcee9a87aeaf0bc8e6c49abb42caad9fce7b | b4185fb48b0f2ffb36c99cd175faa6b3b317d3f26cc7e5c104df63475ba2c922 |
| M30 | 0 | 0 | 1 (unit) | instructions_test.go:481: the rule covers fails on the specimen: the Facts cell of REQ-001 cites F-0003, which is no fact of the inventory or the answers | faa3cd726f05ecda1671992d6f11a33cc4a99298726662da05dfcb95a67b7a0d | 5d73754eaec4f33c0c5d59f42402a27fd9256480464bf9bf3a0ad81c74e7529b |
| M31 | 0 | 0 | 1 (unit) | instructions_test.go:87: Text differs from testdata/prompts/specification.txt: got 5757 bytes, want 5757; first difference at byte 2819: got "3#<n> F-0003 \|\n\n\| Requirement \| Acceptan", want "4#<n> F-0003 \|\n\n\| Requirement  | 0a524319e09f21f0293947a02b4bfcee9a87aeaf0bc8e6c49abb42caad9fce7b | db435d191fbf7442f89299ec2d118ba56a42e298f3f2a98ab01fc89ea6aed554 |
| M32 | 0 | 0 | 1 (unit) | instructions_test.go:524: after the first bet with no numbered facts file: PSBRecords [F-0005]; want none (the prediction is for the head before the bet) | 0a524319e09f21f0293947a02b4bfcee9a87aeaf0bc8e6c49abb42caad9fce7b | 5b3da234464ad61b695a652c458a2b6c8da5956437612e5002ce1f6655c303ad |
| M33 | 0 | 0 | 1 (unit) | instructions_test.go:503: factRecordIDs = [F-0001 F-0002 F-0003 F-0004 F-0009]; want [F-0001 F-0002 F-0003 F-0004] | 0a524319e09f21f0293947a02b4bfcee9a87aeaf0bc8e6c49abb42caad9fce7b | 43602e57874f46ac89761fce3c99e2cbc7da106220cfa6350d9f8a3cd1d7612c |
| M34 | 0 | 0 | 1 (unit) | instructions_test.go:519: with a numbered facts file: PSBRecords [F-0008]; want it unchanged | 0a524319e09f21f0293947a02b4bfcee9a87aeaf0bc8e6c49abb42caad9fce7b | a2d4bc02d86a66e072371e531baba27d3608de53d84d625eecd8fbcd6fcbac61 |
| M35 | 0 | 0 | 1 (unit) | instructions_test.go:514: before the bet: Records [], PSBRecords [F-0001]; want the records and the predicted ID | 0a524319e09f21f0293947a02b4bfcee9a87aeaf0bc8e6c49abb42caad9fce7b | f40cbf258e0aabd26f94408b97a54f727130bb398d08849db9bf67ee19242c6a |

Restores: each checked by hash against the head's file, followed by a green control of `go test -count=1 ./internal/run`
(exit 0 for every mutant). M23 to M35 killed except M28, which survives and is equivalent: `PromptFile` reads
`sessions.tsv` by `ReadSessions` at stage 2 for every step, before `taskInputs` reads it again, so a table that the
reader refuses never reaches the second read (the case `a refused sessions table` of
`TestPromptFileRefusalsHaveTheirControls` shows the refusal comes from that earlier read).

Restore check of every mutant above (M1, M2, M3, M4, M5, M6, M7, M8, M9, M10, M11, M12, M14, M15, M16, M17, M18, M19, M20, M21, M22; M7, M22, M23 to M29 counted once for each run): the file restored byte for byte by hash (True) and the control `go test -count=1 ./internal/run` exit 0. The line of each run is in the output of the script, kept outside the tree.
