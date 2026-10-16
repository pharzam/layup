#!/usr/bin/env python3
"""Mutation audit of T-4cep (#384): each variant changes one source file, builds, vets, runs the tests, and is
restored. A variant counts as killed only if `go build` and `go vet` pass on it and a named test fails in its own
output; the exit status of every run is kept, and the hash of the file before and after proves the restore.
Run from the root of the worktree: python3 -I runs/T-4cep/mutate.py OUTDIR [TAG...]; with tags, only those variants run."""
import hashlib, os, shutil, subprocess, sys, time

OUT = sys.argv[1]
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
 ("m10-check-asked-of-the-last-commit", RUN, "\tif head == last {\n\t\treturn \"\", true\n\t}", "\tif head == last && j.landing == nil {\n\t\treturn \"\", true\n\t}"),
 ("m11-check-error-is-a-pass", RUN, "if err := j.landing(head, last); err != nil {", "if err := j.landing(head, last); err != nil && false {"),
 ("m12-check-asked-of-no-head", RUN, "if head == \"\" || j.landing == nil {", "if j.landing == nil {"),
 ("m13-probe-on-the-last-commit", RUN, "sha, err := j.probe(head)", "sha, err := j.probe(last)"),
 ("m14-no-note-on-admit", RUN, "if j.note != nil {\n\t\tj.note(\"tree: \"", "if false && j.note != nil {\n\t\tj.note(\"tree: \""),
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
INTEGRATION = [
 ["go", "test", "-count=1", "-tags=integration", "./internal/setup", "-run", "TestTheLandingScript$|TestTheLandingCommitIsAdmitted|TestCheckCodeOwnerLanding$|TestAHeadAbove|TestATargetSetUp"],
 ["go", "test", "-count=1", "-tags=integration", "./internal/run", "-run", "TestTheTreeRowOnRealObjects"],
 ["go", "test", "-count=1", "-tags=integration", "./internal/git", "-run", "TestParentsAndTheWholeTree"],
]

def sha(p): return hashlib.sha256(open(p, "rb").read()).hexdigest()
def run(tag, name, cmd):
    t0 = time.time()
    r = subprocess.run(cmd, capture_output=True, text=True)
    open(f"{OUT}/{tag}.{name}.out", "w").write(f"$ {' '.join(cmd)}\n--- stdout\n{r.stdout}\n--- stderr\n{r.stderr}\n--- rc={r.returncode}\n")
    return r.returncode, time.time() - t0, r.stdout + r.stderr
def fails(text): return sorted({l.split()[2].rstrip(":") for l in text.splitlines() if l.startswith("--- FAIL")})
def log(s): print(s, flush=True)

def steps_for(tag):
    for rc_name, cmd in [("unit", UNIT)] + [(f"int{i}", c) for i, c in enumerate(INTEGRATION)]:
        yield rc_name, cmd

ONLY = sys.argv[2:]
if ONLY:
    V = [v for v in V if v[0] in ONLY]
files = {v[1] for v in V}
before = {f: sha(f) for f in files}
log("hash before: " + " ".join(f"{f}={h[:12]}" for f, h in sorted(before.items())))
for tag, f, old, new in V:
    src = open(f).read()
    assert src.count(old) == 1, (tag, src.count(old))
    h0 = sha(f)
    backup = f"{OUT}/backup-{tag}"
    shutil.copy2(f, backup)
    try:
        open(f, "w").write(src.replace(old, new))
        res = {}
        killed_by = []
        for name, cmd in [("build", ["go", "build", "./..."]), ("vet", ["go", "vet", "-tags=integration", "./internal/..."])]:
            rc, dt, out = run(tag, name, cmd)
            res[name] = rc
            if rc != 0:
                break
        valid = res.get("build") == 0 and res.get("vet") == 0
        if valid:
            for name, cmd in steps_for(tag):
                rc, dt, out = run(tag, name, cmd)
                res[name] = rc
                if rc != 0:
                    killed_by = fails(out)
                    break
        log(f"{tag}: rc={res} valid={valid} killed_by={killed_by if valid else 'NOT CREDITED (does not build or vet)'}")
    finally:
        shutil.copy2(backup, f)
        os.remove(backup)
        log(f"{tag}: restored, hash equal to before: {sha(f) == h0}")
after = {f: sha(f) for f in files}
log(f"hash after equals before for every file: {after == before}")
for name, cmd in [("build", ["go", "build", "./..."])] + list(steps_for("control")):
    rc, dt, _ = run("control", name, cmd)
    log(f"control {name}: rc={rc} ({dt:.0f}s)")
log("done")
