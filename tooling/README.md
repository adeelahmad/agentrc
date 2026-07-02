# agentrc tooling

A Go implementation of the `agentrc` CLI and a BuildKit frontend for the
Dockerfile-shaped Agentfile (0.1.0-draft.6): four new keywords —
`IDENTITY`, `CAPABILITY`, `SOP`, `POLICY` — plus an `ADD --remote`
extension, layered on real Dockerfile instructions.

## Layout

- `internal/agentfile` — extracts the agentrc-specific surface from an
  Agentfile and builds `org.agentrc.*` labels. It does **not** implement
  its own Dockerfile parser: it uses BuildKit's own
  `frontend/dockerfile/parser` to tokenize the file, classifies each
  top-level instruction as one of the four agentrc keywords / `ADD
  --remote` or a standard Dockerfile instruction, and blanks out the
  agentrc-specific lines (preserving line numbers) into `CleanedSource` — a
  plain Dockerfile ready for BuildKit's real compiler.
- `internal/llb` — compiles `CleanedSource` with BuildKit's own
  `dockerfile2llb.Dockerfile2LLB` (full, correct standard Dockerfile
  semantics: multi-stage builds, `ARG`/`ENV`, `COPY`, `HEALTHCHECK`, …for
  free), then layers on the agentrc effects: `org.agentrc.*` labels, the
  embedded `/mnt/SOP` file, and `ADD --remote --cached` fetch + embed.
- `cmd/agentrc-frontend` — the BuildKit gateway frontend, built on
  `dockerui` (the same library the real `dockerfile` frontend uses).
- `cmd/agentrc` — the CLI: `init`, `lint`, `lock`, `build`, `inspect`,
  `push`, `pull` (`sign`/`verify`/`run` are stubs — see below).

## Why reuse BuildKit's own Dockerfile compiler?

The spec is explicit that the Agentfile is "Dockerfile-shaped": everything
except the four new keywords keeps exact standard Dockerfile semantics, and
"two build paths ... MUST produce identical artifacts." Reimplementing
multi-stage builds, `ARG` scoping, `COPY`/`ADD` flag handling, and
`HEALTHCHECK` from scratch would be a large, error-prone undertaking with
no upside — BuildKit's own parser/compiler packages
(`frontend/dockerfile/parser`, `frontend/dockerfile/instructions`,
`frontend/dockerfile/dockerfile2llb`) are public Go packages built for
exactly this. `agentrc build` (the native CLI path) goes further and
literally shells out to `docker build` through the same frontend image, so
identical output is guaranteed by construction rather than by parallel
maintenance.

## Building

```bash
go build -o bin/agentrc ./cmd/agentrc
go build -o bin/agentrc-frontend ./cmd/agentrc-frontend
```

## Using the CLI

```bash
./bin/agentrc init                                   # scaffold ./Agentfile
./bin/agentrc lint examples/Agentfile.code-reviewer   # validate
./bin/agentrc lock examples/Agentfile.code-reviewer   # write agentrc.lock (Resolved Manifest)
./bin/agentrc inspect examples/Agentfile.code-reviewer  # local: human-readable summary
./bin/agentrc build examples/Agentfile.code-reviewer -t ghcr.io/org/code-reviewer:1.0
./bin/agentrc push ghcr.io/org/code-reviewer:1.0
./bin/agentrc inspect ghcr.io/org/code-reviewer:1.0   # remote: prints org.agentrc.* labels
```

`agentrc inspect` accepts either a local Agentfile path (prints a static
summary from the extracted file) or a registry reference (pulls the real,
already-built image's config and prints its `org.agentrc.*` labels — this
is the "review before run" path from profiles/oci-package.md §8).

`agentrc build`, `push`, and `pull` shell out to `docker build`/`push`/`pull`
— an agentrc artifact is a completely ordinary OCI image, so there is no
custom registry client to maintain; the local Docker/OCI credential store
(`docker login`) is reused automatically.

## Using the BuildKit frontend directly

Build and load the frontend image locally first:

```bash
docker build -t local/agentrc-frontend:dev -f Dockerfile.frontend .
```

Then either let `agentrc build` invoke it (the default `--frontend-image`
is `local/agentrc-frontend:dev`), or drive it directly:

```bash
docker build -f examples/Agentfile.minimal \
  --build-arg BUILDKIT_SYNTAX=local/agentrc-frontend:dev \
  -t hello:dev examples
```

Once the frontend image is published to a registry, an Agentfile's first
line can be changed to `# syntax=<registry>/<repo>/agentrc-frontend:<tag>`
so plain `docker build -f Agentfile .` dispatches to it automatically,
without `--build-arg`.

## Directive → build-time semantics

- **Standard Dockerfile instructions** (`FROM`, `CMD`, `COPY`, `ADD` without
  `--remote`, `HEALTHCHECK`, `LABEL`, `ENV`, `ARG`, `WORKDIR`, `USER`,
  `EXPOSE`, `RUN`) keep their real Dockerfile semantics exactly, via
  `dockerfile2llb.Dockerfile2LLB` — no reinterpretation.
- **`IDENTITY`** → `org.agentrc.identity.<key>` labels.
- **`CAPABILITY`** → `org.agentrc.capability.<value>=true` labels.
- **`SOP`** (inline, heredoc, or file-backed via `COPY`/`ADD --remote
  .../mnt/SOP`) → always embedded at `/mnt/SOP`; label is a pointer +
  digest (`org.agentrc.sop`, `org.agentrc.sop.sha256`), never the full text.
- **`POLICY`** → `org.agentrc.<key>=<value>`; the `network` namespace's
  `dns:<host>:<port>` value becomes `org.agentrc.network.dns.<host>=<port>`.
  A `POLICY` value that's a URL under `agent.hooks.*` /
  `agent.interrupt_endpoint` auto-derives an attributed network-egress
  label per spec/index.md §8.5.
- **`ADD --remote`** → `--cached` (default) fetches at build time (http/https
  only in this implementation) and embeds the content, emitting a resolved
  digest + `.origin` label; `--runtime` records a `runtime:<url>` reference
  only. `--fail-if-unavailable` (default) aborts the build if resolution
  fails; `--warn-if-unavailable` degrades to a runtime-style reference with
  a warning instead.
- **`--policy-mode digest`** pulls the `agent.*`/`substrate.*`/`model.*`/
  `network.*` labels out into a JSON manifest layer and replaces them with
  a single `org.agentrc.policy.manifest.sha256` pointer label
  (spec/index.md §9.4 leaves the default mode as an open decision; this
  implementation defaults to `inline`).

## Testing

```bash
go build ./...
go vet ./...
go test -race ./...
```

Parser, label, and LLB-translation tests run directly against the real
fixtures in `examples/Agentfile.*`; LLB tests use a fake image-meta-resolver
so they never touch the network (the examples reference fictional base
images like `ghcr.io/acme/pii-redacted-base`).

There is no automated end-to-end `docker build` test in this pass —
verify the frontend manually with the commands above once a working
Docker + BuildKit setup (e.g. `docker buildx`) is available. This
environment's Docker CLI plugin was architecture-mismatched
(darwin/arm64 client, linux/amd64 daemon) so that step could not be run
here.

## Known limitations / follow-ups

- `agentrc sign` / `agentrc verify` / `agentrc run` are stubs (clear error,
  not a partial implementation) — signing needs its own design pass, and
  `run` would mean shipping an execution runtime, which contradicts
  agentrc's own non-goals unless scoped very deliberately.
- `ADD --remote --cached` only fetches `http://`/`https://` sources at
  build time; other schemes (e.g. `mcp://`) can only be used with
  `--runtime`, or degrade to a runtime-style reference under
  `--warn-if-unavailable`.
- `agentrc lock` resolves what it can (base image digest, http(s) `--cached`
  resource digests) and warns rather than failing on anything it can't
  reach — useful for the fictional URLs in the example Agentfiles, but
  means a lock file isn't a guarantee everything was actually pinned.
- `ADD --remote` flag re-parsing (`internal/agentfile/extract.go`'s
  `parseRemoteAdd`) assumes a flag's raw source text matches BuildKit's
  normalized `child.Flags` entry byte-for-byte. Quoted flag values
  (`--chmod="755"`) or a quoted/spaced `<source>`/`<destination>` — both
  legal, if unusual, Dockerfile syntax — aren't supported; use the
  unquoted forms shown in the examples.
- `ADD --remote --cached` fetches (both the frontend and `agentrc lock`) are
  scheme-gated to `http(s)://`, reject loopback/link-local/private/
  metadata-shaped hosts, re-validate every redirect hop, and cap response
  size at 200MiB — best-effort SSRF/DoS mitigation for a URL taken directly
  from Agentfile content, not a substitute for network-level egress
  controls in a shared build service.
