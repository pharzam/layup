# One-off (T-55n2 plan step 4), run by hand from the repository root: fill the Host column of runs/T-55n2/inventory.md from docs/plan/README.md.
import re
def cells(line):
    s=line.strip()
    if not s.startswith('|'): return None
    return [p.strip().replace('\\|','|') for p in re.split(r'(?<!\\)\|', s.strip('|'))]
def table_after(text, heading):
    lines=text.split('\n'); start=next((i for i,l in enumerate(lines) if l.strip().startswith(heading)),None)
    if start is None: return []
    rows=[];header=None
    for l in lines[start+1:]:
        if l.startswith('#'): break
        c=cells(l)
        if c is None:
            if header is not None and rows: break
            continue
        if header is None: header=c; continue
        if all(re.fullmatch(r':?-+:?', x or '-') for x in c): continue
        rows.append(dict(zip(header,c)))
    return rows
plan=open('docs/plan/README.md').read()
host={}
for r in table_after(plan,'## The tasks of phase 1'):
    for k in re.findall(r'`([a-z0-9][a-z0-9-]*)`', r.get('Items','')): host.setdefault(k,[]).append(f"row {r['#']}")
for r in table_after(plan,'## The hosts of the other items'):
    for k in re.findall(r'`([a-z0-9][a-z0-9-]*)`', r.get('Item','')): host.setdefault(k,[]).append(r['Host'].strip())
inv=open('runs/T-55n2/inventory.md').read().split('\n')
keys=set(re.findall(r'^### `([a-z0-9-]+)` — ', '\n'.join(inv), re.M))
n=0
for i,l in enumerate(inv):
    c=cells(l)
    if c and len(c)==6 and re.fullmatch(r'`[a-z0-9-]+`', c[0]) and c[0].strip('`') in keys:
        k=c[0].strip('`'); h=', '.join(host.get(k,[]))
        parts=l.rstrip().split('|')   # ['', key, title, req, size, after, host, '']
        parts[6]=' '+h.replace('|','\\|')+' '
        inv[i]='|'.join(parts); n+=1
open('runs/T-55n2/inventory.md','w').write('\n'.join(inv))
print('filled',n)
