#!/usr/bin/env python3
"""Compiled mutants of internal/rules/basic.go, run on a COPY of the tree.

Usage: python3 -I runs/T-9et9/mutate.py [repo-root]

Each mutant changes basic.go of the copy by one exact text replacement that
keeps the code building. A mutant counts as killed only when `go vet` is clean
(it built) and a test then fails by an assertion; a build or vet failure is
"invalid", never a kill. After each mutant the file is written back from the
bytes read before the run and its SHA-256 is checked, and the run starts and
ends with a control on the unmutated copy.
"""
import hashlib, os, shutil, subprocess, sys, tempfile, time

ROOT = os.path.abspath(sys.argv[1] if len(sys.argv) > 1 else os.getcwd())
FILE = "internal/rules/basic.go"
INT = ("TestTheRegisterIsReadAtTheBaseInARealRepo|TestARegisterThatCannotBeReadIsAnErrorInARealRepo|"
       "TestARenderedFactsRecordInARealRepo|TestTheBlobIsComparedNotTheCheckout|TestARenderedPRDFromItsPreImage|"
       "TestALinkASubmoduleAndAnExecutableAreNoRecords|TestADeletionARenameAndADirectoryAtTheRecordPath|"
       "TestTheNamesOfARealDiffAreExact|TestTheBaseMustBeAnAncestorOfTheHeadInARealRepo|TestTheGuardrailsExceptionInARealRepo")
STEPS = [("vet", ["go", "vet", "-tags=integration", "./internal/rules/"]),
         ("unit", ["go", "test", "-count=1", "./internal/rules/"]),
         ("integration", ["go", "test", "-count=1", "-tags=integration", "-run", INT, "./internal/rules/"])]

# (name, [(old, new), ...]); each old text occurs exactly once in basic.go.
M = [
 ("M01 the register read at the head", [("readRegisterAt(r, regRev)", "readRegisterAt(r, head)")]),
 ("M02 the empty register refusal removed", [("err != nil || len(reg) == 0 {", "err != nil || (len(reg) == 0 && false) {")]),
 ("M03 equality after TrimSpace", [("!bytes.Equal(t.Data, want)", "!bytes.Equal(bytes.TrimSpace(t.Data), bytes.TrimSpace(want))")]),
 ("M04 equality on a prefix", [("!bytes.Equal(t.Data, want)", "!bytes.HasPrefix(t.Data, want)")]),
 ("M05 the mode and type check removed", [('case t.Type != "blob" || t.Mode != "100644":\n\t\t\treturn failure(ReasonNotARegularFile', 'case (t.Type != "blob" || t.Mode != "100644") && false:\n\t\t\treturn failure(ReasonNotARegularFile')]),
 ("M06 the type alone dropped from that check", [('case t.Type != "blob" || t.Mode != "100644":\n\t\t\treturn failure(ReasonNotARegularFile', 'case t.Mode != "100644":\n\t\t\treturn failure(ReasonNotARegularFile')]),
 ("M07 an absent entry not refused as absent", [("case !t.Present:", "case !t.Present && false:")]),
 ("M08 return at the first passing record", [("\t\trendered++\n", "\t\trendered++\n\t\treturn Result{Success, ReasonRenderedRecord, \"\"}\n")]),
 ("M09 the first record covers any path", [("rec, named := byPath[t.Path]\n", "rec, named := byPath[t.Path]\n\t\tif !named && len(records) > 0 {\n\t\t\trec, named = records[0], true\n\t\t}\n")]),
 ("M10 a rendering finding ignored", [("case len(findings) > 0:", "case len(findings) > 0 && false:")]),
 ("M11 the PRD pre-image read at the head", [("e, ok, err := entryAt(r, rev, p)", "e, ok, err := entryAt(r, head, p)"), ("pre[p], err = r.Show(rev, p)", "pre[p], err = r.Show(head, p)")]),
 ("M12 the base ancestor check removed", [("} else if !ok {\n\t\treturn errorOf(ReasonBase", "} else if !ok && false {\n\t\treturn errorOf(ReasonBase")]),
 ("M13 the chain of the pre-image removed", [("\t\tif !ok {\n\t\t\treturn nil, errorOf(ReasonInput, \"\")", "\t\tif !ok && false {\n\t\t\treturn nil, errorOf(ReasonInput, \"\")")]),
 ("M14 the pre-image may be the head", [("err != nil || rev == head {", "err != nil {")]),
 ("M15 a clean-form check removed", [("p != path.Clean(p) ||", "(p != path.Clean(p) && false) ||")]),
 ("M16 a leading slash allowed", [('strings.HasPrefix(p, "/") ||', "false ||")]),
 ("M17 a dot part allowed", [('part == "." || part == ".." ||', 'part == "." && false || part == ".." ||')]),
 ("M18 a dot-dot part allowed", [('part == "." || part == ".." ||', 'part == "." || part == ".." && false ||')]),
 ("M19 a .git part allowed", [('|| strings.EqualFold(part, ".git")', "|| (strings.EqualFold(part, \".git\") && false)")]),
 ("M20 a backslash or NUL allowed", [('strings.ContainsAny(p, "\\\\\\x00")', "false")]),
 ("M21 a leading colon allowed", [('strings.HasPrefix(p, ":") ||', "false ||")]),
 ("M22 an error mapped to a success", [("func errorOf(reason, p string) Result { return Result{Error, reason, p} }", "func errorOf(reason, p string) Result { return Result{Success, reason, p} }")]),
 ("M23 the base error mapped to a success", [('return errorOf(ReasonBase, "")', 'return Result{Success, ReasonNoRulePath, ""}')]),
 ("M24 the guardrails exception without its diff", [("if t.Guardrails.Diff, err = r.DiffFile(base, head, p); err != nil {", "if _, err = r.DiffFile(base, head, p); err != nil {")]),
 ("M25 the exception never given to Check", [("Check([]string{t.Path}, reg, t.Guardrails)", "Check([]string{t.Path}, reg, Change{})")]),
 ("M26 the diff of the guardrails taken backwards", [("r.DiffFile(base, head, p)", "r.DiffFile(head, base, p)")]),
 ("M27 the diff names taken backwards", [("r.DiffNames(base, head)", "r.DiffNames(head, base)")]),
 ("M28 the register's commit may be the head", [("if err != nil || regRev == head {", "if err != nil {")]),
 ("M29 a bad path read before it is refused", [("\t\tif !validPath(rec.Path) {\n\t\t\treturn errorOf(ReasonInput, rec.Path)\n\t\t}\n\t}\n\tbase,", "\t\tif !validPath(rec.Path) && false {\n\t\t\treturn errorOf(ReasonInput, rec.Path)\n\t\t}\n\t}\n\tbase,")]),
 ("M30 a PRD with no pre-image accepted", [("rec.Kind == PRD && !hasPre:", "rec.Kind == PRD && !hasPre && false:")]),
 ("M31 a record twice accepted", [("; twice {", "; twice && false {")]),
 ("M32 the workflow reason not a rule path", [('reason == ""', 'reason != "rule-path"')]),
 ("M33 a register that is no regular file accepted", [('case !ok || e.Type != "blob" || (e.Mode != "100644" && e.Mode != "100755"):', 'case !ok || ((e.Type != "blob" || (e.Mode != "100644" && e.Mode != "100755")) && false):')]),
 ("M34 the reason of a rendered record", [("if rendered > 0 {", "if rendered > 1 {")]),
 ("M35 the record's bytes never read", [("if p == GuardrailsPath || named[p] {", "if p == GuardrailsPath {")]),
 ("M36 the guardrails head never read", [("if t.Present && t.Type == \"blob\" {\n\t\t\t\tif t.Guardrails.Head", "if t.Present && t.Type == \"blob\" && false {\n\t\t\t\tif t.Guardrails.Head")]),
 ("M37 the kind not checked", [("rec.Kind != PRD && rec.Kind != Facts,", "rec.Kind != PRD && rec.Kind != Facts && false,")]),
 ("M38 the guardrails path may be a record", [("rec.Path == GuardrailsPath,", "rec.Path == GuardrailsPath && false,")]),
]

def sha(b): return hashlib.sha256(b).hexdigest()
def run(cwd, step, cmd, tag):
    t0 = time.time()
    r = subprocess.run(cmd, cwd=cwd, capture_output=True, text=True)
    out = r.stdout + r.stderr
    print(f"    [{tag}] {step}: rc={r.returncode} ({time.time()-t0:.0f}s)", flush=True)
    return r.returncode, out

def failing(out):
    lines = [l.strip() for l in out.splitlines() if l.lstrip().startswith("--- FAIL")]
    first = next((l.strip() for l in out.splitlines() if "_test.go:" in l), "")
    return lines[:3], first[:160]

def main():
    tmp = tempfile.mkdtemp(prefix="t9et9-mut-")
    work = os.path.join(tmp, "tree")
    shutil.copytree(ROOT, work, ignore=shutil.ignore_patterns(".git", ".worktree", "runs"))
    path = os.path.join(work, FILE)
    pristine = open(path, "rb").read()
    h0 = sha(pristine)
    print(f"copy of the tree: {work}\npristine {FILE} sha256 {h0}", flush=True)
    def control(label):
        for step, cmd in STEPS:
            rc, out = run(work, step, cmd, label)
            if rc != 0:
                print(out); sys.exit(f"the control {label} fails at {step}: no mutant counts")
        print(f"control {label}: green", flush=True)
    control("control-before")
    results = []
    for n, (name, edits) in enumerate(M, 1):
        text = pristine.decode()
        for old, new in edits:
            if text.count(old) != 1:
                sys.exit(f"{name}: the text occurs {text.count(old)} times, want 1: {old!r}")
            text = text.replace(old, new)
        print(f"[{n}/{len(M)}] {name}", flush=True)
        verdict = "survived"
        try:
            open(path, "w").write(text)
            for step, cmd in STEPS:
                rc, out = run(work, step, cmd, "mutant")
                if step == "vet" and rc != 0:
                    verdict = "invalid (does not build or vet)"; print(out); break
                if step != "vet" and rc != 0:
                    names, first = failing(out)
                    verdict = f"killed by {step}: {', '.join(n.split()[2] for n in names) or '?'}; {first}"
                    break
        finally:
            open(path, "wb").write(pristine)
            if sha(open(path, "rb").read()) != h0:
                sys.exit("the restore wrote other bytes")
        print(f"    => {verdict}", flush=True)
        results.append((name, verdict))
    control("control-after")
    print("\nSUMMARY")
    for name, verdict in results:
        print(f"- {name}: {verdict}")
    print(f"killed {sum(v.startswith('killed') for _, v in results)}, survived {sum(v == 'survived' for _, v in results)}, "
          f"invalid {sum(v.startswith('invalid') for _, v in results)} of {len(results)}")
    shutil.rmtree(tmp)

main()
