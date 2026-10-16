import subprocess, hashlib, shutil, os, time, sys, json
ROOT='/home/layup/projects/layup/.worktree/T-jew3'
os.chdir(ROOT)
F='internal/setup/codeowner.go'
V=[
 ("m7-name-guard", ' or .name != "%s"', ' or (.name != "%s" and false)'),
 ("m8-source-guard", ' or .source != "%s"', ' or (.source != "%s" and false)'),
 ("m9-one-pull-request-rule", ' or (%s) != 1 then', ' or ((%s) != 1 and false) then'),
 ("m10-restore-before-file", '--input \\"$after\\" > /dev/null && [', '--input \\"${after%default.after.json}default.before.json\\" > /dev/null && ['),
]
STEPS=[("build",["go","build","./..."]),("vet",["go","vet","./internal/setup"]),
 ("unit",["go","test","-count=1","./internal/setup"]),
 ("integration",["go","test","-count=1","-tags=integration","./internal/setup","-run","TestTheUpdateRefusals|TestTheUpdateKeepsWhatIsUnknown|TestTheLandingScript"])]
def sha(p): return hashlib.sha256(open(p,'rb').read()).hexdigest()
def run(tag,name,cmd):
    t0=time.time(); r=subprocess.run(cmd,capture_output=True,text=True)
    open(f'/tmp/mut2/{tag}.{name}.out','w').write(f'$ {" ".join(cmd)}\n--- stdout\n{r.stdout}\n--- stderr\n{r.stderr}\n--- rc={r.returncode}\n')
    return r.returncode, time.time()-t0, r.stdout+r.stderr
summary=[]
def log(s):
    print(s,flush=True); summary.append(s)
before=sha(F); log(f"hash before all: {before}")
for tag,old,new in V:
    src=open(F).read(); h0=sha(F)
    assert src.count(old)==1,(tag,src.count(old))
    bdir=f'/tmp/mut2/backup-{tag}'; os.makedirs(bdir,exist_ok=True); shutil.copy2(F,f'{bdir}/codeowner.go')
    try:
        open(F,'w').write(src.replace(old,new))
        res={}
        for name,cmd in STEPS:
            rc,dt,out=run(tag,name,cmd); res[name]=(rc,dt,out)
            if name in("build","vet") and rc!=0: break
            if name=="unit" and rc!=0: break   # killed at the unit level: integration not needed
        row={k:v[0] for k,v in res.items()}
        failing=[l.strip() for k,v in res.items() if k in("unit","integration") and v[0]!=0 for l in v[2].splitlines() if l.startswith("--- FAIL") or "codeowner_integration_test.go" in l][:6]
        log(f"{tag}: rc={row} valid_build={res['build'][0]==0 and res['vet'][0]==0} assertions={failing}")
    finally:
        shutil.copy2(f'{bdir}/codeowner.go',F)
        h1=sha(F); log(f"{tag}: restored, hash equal to before: {h1==h0}")
        shutil.rmtree(bdir)
log(f"hash after all: {sha(F)} equal={sha(F)==before}")
for name,cmd in STEPS:
    rc,dt,_=run('control',name,cmd); log(f"control {name}: rc={rc} ({dt:.0f}s)")
log("done")
