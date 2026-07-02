---
type: validate
story: S2-02
---
# S2-02 Validation — Spec draft.6 (T17–T20)

## Pre-flight

- [ ] Confirmed §8.5/§8.6 already occupy those numbers (`spec/index.md:400,415`); the three
      new subsections go at §8.7/§8.8/§8.9 (the renumbering decision), NOT at 8.5/8.6.
- [ ] Grep-located EVERY `draft.5` occurrence before the T20 bump (§0.9, M-003): spec, cli,
      examples, `_config.yml`/front-matter, CHANGELOG, `internal/**/lock.go`.
- [ ] Confirmed T17–T19 add NO new keyword and exactly three namespaces; `substrate.*` /
      "substrate-neutral" NOT renamed (§0.2, §0.3).
- [ ] Version still `draft.5` during T17–T19 (bump is T20 only, §0.1).

## Final sign-off

| # | Task | Command | Expected |
|---|------|---------|----------|
| 1 | T17 §8.7 present | `grep -cE '^### 8\.7 .*substrate\.<platform>' spec/index.md` | `>= 1` |
| 2 | T17 AWS keys | `grep -cE 'roleArn\|networkMode\|code\.s3\.uri' spec/index.md` | `>= 1` each |
| 3 | T18 §8.8 present | `grep -cE '^### 8\.8 .*agent\.auth' spec/index.md` | `>= 1` |
| 4 | T18 fail-closed jwt | `grep -c 'MUST NOT expose the invocation endpoint' spec/index.md` | `>= 1` |
| 5 | T19 §8.9 present | `grep -cE '^### 8\.9 .*substrate\.runtime\.language' spec/index.md` | `>= 1` |
| 6 | no collision | `grep -cE '^### 8\.5 \|^### 8\.6 ' spec/index.md` | `2` (existing untouched) |
| 7 | namespace intact | `rg -c 'POLICY substrate\.' spec/index.md` | `>= 1` |
| 8 | keyword count | `grep -oE '\b(IDENTITY\|CAPABILITY\|SOP\|POLICY)\b' spec/index.md \| sort -u \| wc -l` | `4` |
| 9 | T20 version bump | `grep -rhoE 'draft\.[0-9]+' . \| grep -v .git \| sort -u` (AFTER T20) | single `draft.6` |
| 10 | T20 lock.go bumped | `grep -c '0.1.0-draft.6' internal/*/lock.go internal/**/lock.go 2>/dev/null` | `>= 1` |
| 11 | T8 subsection | `grep -c 'Reproducible builds' spec/index.md` and `grep -c 'format TODO' spec/index.md` | `>= 1` each |
| 12 | T8 slogan kept | `grep -rc 'reproducible' index.md 2>/dev/null` (homepage slogan present) | `>= 1` |
| 13 | code-reviewer commented block | `grep -cE '#.*substrate\.aws\.\|#.*agent\.auth\.' examples/Agentfile.code-reviewer` | `>= 1` |
| 14 | code-reviewer still lints | `go run ./cmd/agentrc lint examples/Agentfile.code-reviewer` | exit 0 |
| 15 | site build (CI) | `bundle exec jekyll build --trace` | exit 0 |
</content>
