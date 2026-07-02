---
type: tasks
story: S1-01
---
# S1-01 Tasks — Content correctness (T1–T4)

All targets grep-located against the live repo (§0.8, M-003). Do NOT hand-edit inline
`hello` snippets independently — propagate the canonical from `examples/Agentfile.minimal`
(§0.6, M-001).

## T1 — Kill the ghost healthcheck tool (`/mnt/tools/ping`) [P0]

Replace each ghost healthcheck line with `file_read --agentrc-schema`, KEEPING each file's
existing HEALTHCHECK options. `--agentrc-schema` mirrors existing usage and does NOT resolve
§14.2 #3 — add no prose claiming it is decided (§0.5).

Exact edits (grep-located surfaces):
- `index.md:38` — `HEALTHCHECK --interval=60s --timeout=15s CMD /mnt/tools/ping` →
  `HEALTHCHECK --interval=60s --timeout=15s CMD /mnt/tools/file_read --agentrc-schema`
- `spec/index.md:636` — `HEALTHCHECK --interval=60s --timeout=15s --retries=3 CMD /mnt/tools/ping` →
  `... --retries=3 CMD /mnt/tools/file_read --agentrc-schema`
- `spec/index.md:716` — same replacement as line 636 (preserve `--retries=3`)
- `examples/Agentfile.minimal:23` — `... CMD /mnt/tools/ping` →
  `... CMD /mnt/tools/file_read --agentrc-schema`
- `examples/Agentfile.code-reviewer:52` — `... --retries=3 CMD /mnt/tools/ping` →
  `... --retries=3 CMD /mnt/tools/file_read --agentrc-schema`
- `examples/index.md:45` — `... CMD /mnt/tools/ping` →
  `... CMD /mnt/tools/file_read --agentrc-schema`

Each affected file already `COPY`s `file_read`; `Agentfile.secure-workspace` is already
correct (no `tools/ping`). Do not add a new COPY.

Verify: `grep -rn "tools/ping" . | grep -v .git | wc -l` → 0 (the two doc-artifact
mentions in `docs/agents/` memory/work-order/intake describe the ghost as a string and are
outside the site tree, but the §V check counts the whole tree — confirm the only remaining
hits are inside `docs/agents/` planning docs, which are non-published and expected; the six
site surfaces above must be 0).

## T2 — One hello, always with `FROM` [P0]

1. Ensure every rendered `hello` snippet has `FROM python:3.11-slim` as the first
   instruction after the `# syntax=` line so it builds under BuildKit, byte-identical to
   `examples/Agentfile.minimal`. Grep-located hello renders needing a FROM audit (from
   `grep -rln "IDENTITY name=hello" --include="*.md" --include="*.html" .`): `index.md`,
   `docs/quickstart.md`, `docs/agentfile.md`, `cli.md`, `examples/index.md`
   (`examples/Agentfile.minimal` already has `FROM python:3.11-slim` at line 5). For any
   inline hello lacking a `FROM`, insert `FROM python:3.11-slim` within 3 lines above the
   `IDENTITY name=hello` line and align the whole snippet to the canonical file.
2. Add to spec §2 (`spec/index.md`) the MUST sentence: "Exactly as in Dockerfile, every
   Agentfile MUST contain a `FROM` instruction, and `FROM` must be the first instruction
   after the `# syntax=` line, comments, and any `ARG` that `FROM` consumes."

Verify: every located file has `FROM python:3.11-slim` within 3 lines above hello's
`IDENTITY`; the spec sentence is present; each inline snippet diffs clean against
`examples/Agentfile.minimal` (M-001 — diff, don't just grep).

## T3 — Quickstart honesty callouts [P0]

In `docs/quickstart.md`, mirror `/cli/` wording exactly; invent no status claims (M-002).
- At step 2, add: "Status: the agentrc frontend image is not yet published to a public
  registry, so `docker build -f Agentfile .` won't auto-route through it yet. Build the
  frontend locally first (see `tooling/README.md`) or pass
  `--build-arg BUILDKIT_SYNTAX=<your-built-image>`. Details on the CLI page (/cli/)."
- At step 5, add: "Status: `arc run` is planned — agentrc declares agents, it does not ship
  a runtime. See the CLI status table (/cli/)."
- Reword remaining present-tense claims to future tense. Grep-located:
  `docs/quickstart.md:71` "published a stock Docker / BuildKit install will need no extra
  tooling" — ensure it reads as future/conditional, not a present-tense guarantee.

Note (open question #1): follow the work-order wording above; the published-image
discrepancy (`ghcr.io/adeelahmad/agentrc-frontend` was published this session) is flagged to
the owner in `stories.md` › Notes, NOT silently reworded here.

Verify: both callouts present; zero present-tense stock-`docker build` "nothing to
install"/"no extra tooling" claims remain
(`grep -rn "nothing to install\|no extra tooling" docs/` → no present-tense hits).

## T4 — Polish [P3]

- `examples/Agentfile.minimal:8` — description `"Minimal AgentRC agent"` →
  `"Minimal agentrc agent"` (brand casing).
- Align the quickstart hello (`docs/quickstart.md`) to canonical fields per §0.6:
  `IDENTITY name=hello version=0.1 author=acme`, `POLICY model.name claude-sonnet-4`.

Verify: `grep -rn "Minimal AgentRC agent" .` → 0; the quickstart hello diffs clean against
`examples/Agentfile.minimal`.
