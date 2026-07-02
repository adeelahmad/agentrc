---
type: validate
story: S2-05
---
# S2-05 Validation — CLI docs table (T25)

## Pre-flight

- [ ] S2-03 (rename) + S2-04 (backends) landed — the table reflects real status (M-002).
- [ ] Grep-located the `run`/`sign`/`verify` rows in `cli.md` (lines 99/100/103) before editing.
- [ ] Confirmed the §0.8 line will be pasted VERBATIM (no paraphrase; §0.6-standards).

## Final sign-off

| # | Task | Command | Expected |
|---|------|---------|----------|
| 1 | run implemented | `grep -c 'reference translators' cli.md` | `>= 1` |
| 2 | §0.8 verbatim | `grep -c 'Reference translators — a proof of concept until platforms read' cli.md` | `>= 1` |
| 3 | sign planned | `grep -E 'arc sign.*planned\|sign.*\` planned' cli.md \| wc -l` | `>= 1` |
| 4 | verify planned | `grep -E 'arc verify.*planned\|verify.*\` planned' cli.md \| wc -l` | `>= 1` |
| 5 | no --substrate | `rg -c -- '--substrate' cli.md 2>/dev/null \|\| echo 0` | `0` |
| 6 | namespace prose intact | `grep -c 'Agentfile never names a substrate' cli.md` | `>= 1` |
| 7 | site build (CI) | `bundle exec jekyll build --trace` | exit 0 |
</content>
