#!/usr/bin/env python3
"""Mutation audit of T-4cep (#384): one replacement per variant, then go build, go vet and the tests; a variant is valid only if the
build and the vet pass on it, and killed only by a `--- FAIL` line of its own output. The hash of the file before and after proves the
restore. From the root of the worktree: python3 -I runs/T-4cep/mutate.py OUTDIR [TAG...] (with tags, only those variants run)."""
import hashlib, os, shutil, subprocess, sys

OUT, ONLY = sys.argv[1], sys.argv[2:]
os.makedirs(OUT, exist_ok=True)

SETUP, RUN, GIT = "internal/setup/codeowner.go", "internal/run/forgecheck.go", "internal/git/git.go"
RUNCO = "internal/run/codeowner.go"
V = [
 ("m01-parent-not-compared", SETUP, "case parents[0] != last:", "case false:"),
 ("m02-more-than-one-parent", SETUP, "case len(parents) != 1:", "case len(parents) < 1:"),
 ("m03-tree-not-compared", SETUP, "if len(differ) > 0 {", "if false && len(differ) > 0 {"),
 ("m04-mode-ignored", SETUP, "case e.Mode != landingMode || e.Type != \"blob\":", "case e.Mode != landingMode && e.Type != \"blob\":"),
 ("m05-bytes-not-compared", SETUP, "case e.Object != git.BlobID(want):", "case false:"),
 ("m06-last-may-have-the-file", SETUP, "if _, has := l[codeOwnersPath]; has {", "if _, has := l[codeOwnersPath]; has && false {"),
 ("m07-config-paths-ignored", SETUP, "want, err := CodeOwners(login, configs)", "want, err := CodeOwners(login, configs[:0])"),
 ("m08-login-is-the-first-row", RUNCO, "return operatorLogin(approvers)", "return approvers[0][2], nil"),
 ("m09-record-gate-ignored", RUNCO, "ok {\n\t\treturn \"\", errors.New(\"the setup record has the row", "ok && false {\n\t\treturn \"\", errors.New(\"the setup record has the row"),
 ("m10-check-asked-of-the-last-commit", RUN, "\tif head == last {\n\t\treturn \"\", true\n\t}", "\tif head == last && false {\n\t\treturn \"\", true\n\t}"),
 ("m11-check-error-is-a-pass", RUN, "if err := j.landing(head, last); err != nil {", "if err := j.landing(head, last); err != nil && false {"),
 ("m12-check-asked-of-no-head", RUN, "if head == \"\" {\n\t\treturn plain, false", "if false {\n\t\treturn plain, false"),
 ("m13-probe-on-the-last-commit", RUN, "sha, err := j.probe(head)", "sha, err := j.probe(last)"),
 ("m14-no-note-on-admit", RUN, "j.note(\"tree: \" + head + \" is the code-owner commit above \" + last)", "_ = head"),
 ("m15-detail-on-admit", RUN, "\treturn \"\", true\n}\n\n// lastOfTree", "\treturn \"admitted\", true\n}\n\n// lastOfTree"),
 ("m16-script-refuses-again", SETUP, "b.WriteString(\"set -eu\\n\")", "b.WriteString(\"set -eu\\nexit 1\\n\")"),
 ("m17-script-adds-another-path", SETUP, "+ codeOwnersPath + \"\\\"\\n\")\n\tb.WriteString(\"tree=", "+ \"docs/extra.md\\\"\\n\")\n\tb.WriteString(\"tree="),
 ("m18-script-commit-on-the-wrong-parent", SETUP, "-p refs/heads/layup-setup -m", "-p refs/heads/layup-setup~1 -m"),
 ("m19-trap-missing", SETUP, "b.WriteString(\"trap restore EXIT\\n\")", "b.WriteString(\"\")"),
 ("m20-restore-the-before-object", SETUP, "--input \\\"$after\\\" > /dev/null && [", "--input \\\"${after%default.after.json}default.before.json\\\" > /dev/null && ["),
 ("m21-no-disable-before-the-push", SETUP, "b.WriteString(\"gh api --method PUT \" + rs + \" --input \\\"$off\\\" > /dev/null\\n\")", "b.WriteString(\"\")"),
 ("m22-git-parents-include-the-commit", GIT, "return f[1:], nil", "return f, nil"),
 ("m23-blob-header-without-the-length", GIT, "fmt.Fprintf(h, \"blob %d\\x00\", len(data))", "fmt.Fprintf(h, \"blob %d\\x00\", 0)"),
]
UNIT = ["go", "test", "-count=1", "./internal/setup", "./internal/run", "./internal/git"]
INT = ["go", "test", "-count=1", "-tags=integration"]
STEPS = [("unit", UNIT),
 ("int0", INT + ["./internal/setup", "-run", "TestTheLandingScript$|TestTheLandingCommitIsAdmitted|TestCheckCodeOwnerLanding$|TestAHeadAbove|TestATargetSetUp"]),
 ("int1", INT + ["./internal/run", "-run", "TestTheTreeRowOnRealObjects"]), ("int2", INT + ["./internal/git", "-run", "TestParentsAndTheWholeTree"])]

def sha(p): return hashlib.sha256(open(p, "rb").read()).hexdigest()
def run(tag, name, cmd):
    r = subprocess.run(cmd, capture_output=True, text=True)
    open(f"{OUT}/{tag}.{name}.out", "w").write(f"$ {' '.join(cmd)}\n{r.stdout}\n{r.stderr}\n--- rc={r.returncode}\n")
    return r.returncode, r.stdout + r.stderr
def log(s): print(s, flush=True)

V = [v for v in V if not ONLY or v[0] in ONLY]
before = {f: sha(f) for f in {v[1] for v in V}}
for tag, f, old, new in V:
    src, h0 = open(f).read(), sha(f)
    assert src.count(old) == 1, (tag, src.count(old))
    try:
        open(f, "w").write(src.replace(old, new))
        res, killed = {}, []
        for name, cmd in [("build", ["go", "build", "./..."]), ("vet", ["go", "vet", "-tags=integration", "./internal/..."])] + STEPS:
            res[name], out = run(tag, name, cmd)
            if res[name] != 0:
                killed = sorted({l.split()[2].rstrip(":") for l in out.splitlines() if l.startswith("--- FAIL")})
                break
        valid = res.get("build") == 0 and res.get("vet") == 0
        log(f"{tag}: rc={res} valid={valid} killed_by={killed if valid else 'NOT CREDITED (does not build or vet)'}")
    finally:
        open(f, "w").write(src)
        log(f"{tag}: restored, hash equal to before: {sha(f) == h0}")
log(f"hash after equals before for every file: {before == {f: sha(f) for f in before}}")
for name, cmd in [("build", ["go", "build", "./..."])] + STEPS:
    log(f"control {name}: rc={run('control', name, cmd)[0]}")
