---
type: report
sprint: 1
task: T8
---

# T8 — Lockfile Investigation Report

Investigate-and-report-only. This document records what `arc lock` actually
does, whether `arc build` consumes its output, the mismatch with the homepage
slogan, and an A/B recommendation for the owner. No decision is taken here; no
site/spec/CLI content is changed. Every format detail below is quoted from the
real code; nothing is invented.

Sources investigated:
- `cmd/agentrc/lock.go` — the `lock` command implementation.
- `cmd/agentrc/build.go` — the `build` command implementation.
- `schemas/agentrc-lock.schema.json` — the "agentrc Resolved Manifest" schema.
- `internal/agentfile/` — the extract/label/resource types the command calls.
- `index.md`, `cli.md`, `docs/package.md` — the site references to a lockfile.

## What `arc lock` emits

**Output filename.** Default `agentrc.lock`, overridable with `--out`.
`cmd/agentrc/lock.go:159-161`:

```go
if out == "" {
    out = "agentrc.lock"
}
```

Flag: `cmd/agentrc/lock.go:170` — `--out` "output path (default: agentrc.lock)".
There is also `--policy-mode` (`inline`|`digest`, default `inline`) at
`cmd/agentrc/lock.go:171`.

**On-disk format.** Pretty-printed JSON (2-space indent). The struct is
serialized with `json.MarshalIndent(manifest, "", "  ")` at
`cmd/agentrc/lock.go:155`. The command's `Short` help calls its purpose "Pin
`ADD --remote` resources to digests for reproducible builds"
(`cmd/agentrc/lock.go:61`).

**Top-level structure — the "Resolved Manifest".** The Go type
`resolvedManifest` (`cmd/agentrc/lock.go:26-35`) mirrors
`schemas/agentrc-lock.schema.json` (title: "agentrc Resolved Manifest"). Top-level
fields:

| Field | JSON key | Source in lock.go | Meaning |
|-------|----------|-------------------|---------|
| Version | `version` | hardcoded `"0.1.0-draft.5"` (L77) | manifest schema version |
| AgentfileSHA256 | `agentfile_sha256` | `hashHex([]byte(f.Source))` (L78) | SHA-256 of the Agentfile source |
| Timestamp | `timestamp` | `time.Now().UTC()` RFC 3339 (L79) | when resolved |
| PolicyMode | `policy_mode` | from `--policy-mode` flag (L80) | how POLICY requests are encoded |
| LabelsDigest | `labels_digest` | `sha256:` + hash of marshaled labels (L149-153) | integrity digest over `ai.agentrc.*` labels |
| Base | `base` | `{ref, digest}` (L83-91) | resolved `FROM` ref and (best-effort) registry digest |
| Resources | `resources` | array (see below) | resolved COPY/ADD resources |
| SOP | `sop` | `{sha256}` (L140-147) | digest of the embedded SOP, never the full text |

Note: the schema declares an optional `build_id` field, but `lock.go` never
populates it — the tool does not emit it today.

**The record set (`resources[]`).** Each entry is a `resolvedResource`
(`cmd/agentrc/lock.go:42-50`) with `name`, `kind`, `dest`, `delivery`,
`digest`, `origin`, `fail_mode`. Records are built from two passes:

- **Local resources** (`f.LocalResources`, from `COPY`) — `delivery: "local"`
  (L96-109). The `/mnt/SOP` COPY is diverted into the top-level `sop` object
  (hashed from the local file) rather than listed as a generic resource.
- **Remote adds** (`f.RemoteAdds`, from `ADD --remote`) — L111-138.
  `runtime` adds get `delivery: "runtime"` with no digest; otherwise the tool
  HTTP-fetches the resource (`fetchHTTP`, bounded to 200 MiB), records
  `delivery: "cached"` and `digest: "sha256:<hex>"`. `fail_mode` is `"fail"` or
  `"warn"` depending on `ra.FailIfUnavailable`.

So the lock pins: the Agentfile source hash, the base image ref+digest, each
resource's content digest (for embedded/cached resources) or a runtime pointer,
the SOP digest, and a digest over the emitted label set.

Digest resolution is best-effort: unreachable bases/resources emit a `warning:`
to stderr and the corresponding `digest` is simply omitted rather than failing
the command.

## How `arc build` consumes it

**It does not.** `arc build` (`cmd/agentrc/build.go`) never reads, parses, or
references `agentrc.lock` or the `resolvedManifest` type. What `build` actually
does (`build.go:26-63`):

1. Reads and extracts the Agentfile (`readAndExtract`, L36).
2. Runs local validation (`agentfile.Validate`, L40-42).
3. Shells out to `docker build -f <file> --build-arg BUILDKIT_SYNTAX=<frontend>
   --build-arg AGENTRC_POLICY_MODE=<mode> -- <context>` (L44-61).

A repo-wide search confirms nothing reads the lock back: no `os.ReadFile` /
`json.Unmarshal` of the manifest exists anywhere in `cmd/` or `internal/`, and
`resolvedManifest` is referenced only inside `lock.go`.

**Implication.** The lockfile is currently **produced but not consumed**.
`arc lock` writes `agentrc.lock`, but `arc build` ignores it — builds are not
pinned to, gated by, or reproducible-from the lock today. The "reproducible
builds" framing in the command help and schema describes an intended future
relationship, not a wired-up one. Any documentation must state that `build` does
not yet read the lock; the manifest is presently an informational/audit artifact
only.

## The gap

The homepage core slogan (`index.md:70`) states, verbatim:

> The Agentfile declares one agent. **The lockfile pins dependencies.** The
> package makes it portable. The policy makes boundaries reviewable. The
> registry makes it shareable. Compatible runners execute it.

The CLI page lists `agentrc lock` / `arc lock` as `implemented` (`cli.md:81`),
and `docs/package.md:10` mentions "no lockfile-as-a-package". But the
**specification contains zero lockfile content** — there is no §9 (or other)
spec section defining `agentrc.lock`, the Resolved Manifest, its fields, or how
a build relates to it. So a first-class concept in the marketing slogan
("the lockfile pins dependencies") has:

- a working CLI producer (`arc lock` → `agentrc.lock`),
- a JSON schema (`schemas/agentrc-lock.schema.json`),
- **no consumer** in `arc build`, and
- **no normative spec text**.

The mismatch: the slogan implies the lockfile is a load-bearing, build-time
dependency-pinning mechanism, while in reality it is an unconsumed audit
artifact undocumented by the spec.

## Options (for the owner — NO decision taken)

Both options are presented verbatim in intent. The owner decides at Sprint 2
planning. Neither is selected here.

### Option A — Add an informative spec subsection under §9

Add a spec subsection under §9 ("Reproducible builds / `agentrc.lock`")
documenting **only what the tooling actually does** (filename `agentrc.lock`,
JSON Resolved Manifest, the fields/records above), marked
**"Status: informative in this draft; format TODO"**. Keep the homepage slogan.

- **Pros:** Closes the "slogan with no spec" gap; gives readers a real
  reference for the produced artifact; keeps the marketing narrative intact;
  honest about maturity via the Status marker.
- **Cons:** Documents a format that `build` does not yet consume, which can read
  as over-promising unless the "not consumed by build today" caveat is explicit;
  adds spec surface that may churn once the format is finalized ("format TODO").

### Option B — Cut the slogan sentence; lockfile lives on `/cli/` only

Remove "The lockfile pins dependencies." from the homepage slogan. The lockfile
remains documented only as a CLI command on `/cli/` (where `arc lock` is already
listed).

- **Pros:** Removes the unmet promise immediately; no new spec surface; the
  slogan only claims what the spec actually defines; lowest risk of
  over-promising.
- **Cons:** Drops a genuinely differentiating concept from the top-level
  narrative; the working `arc lock` tool becomes less discoverable; may need
  re-adding later once the format/consumer land, undoing the slogan edit.

### recommendation (owner decides at Sprint 2 planning)

Lean **Option A** — the tool and schema already exist, so an *informative*,
clearly-marked ("format TODO", "not yet consumed by `build`") §9 subsection is
more honest and useful than deleting a real, differentiating capability from the
slogan; but this is the owner's call at Sprint 2 planning, not a decision taken
here.
