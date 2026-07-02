---
type: tasks
story: S2-05
---
# S2-05 Tasks — CLI docs table (T25)

Re-derive every status claim from the now-real backends (M-002). §0.8 line VERBATIM.
Runs in Wave 4 (after S2-03 rename + S2-04 backends). Grep-locate the table first (M-003).

## T25 — CLI docs table [P1]

In `cli.md`:
- Update the `run` row (grep-located at `cli.md:103`): status →
  `implemented (local, bedrock, kubernetes — reference translators)`. Drop the pre-rename
  `--isolation`/`--substrate` phrasing in that row; reference `--backend`.
- `sign` (`cli.md:99`) and `verify` (`cli.md:100`) stay `planned` (they remain stubs — do
  NOT flip them; M-002).
- Place the §0.8 positioning line VERBATIM ABOVE the table: "Reference translators — a proof
  of concept until platforms read `org.agentrc.*` labels natively. Not production runners."
- Ensure no `--substrate` mention survives anywhere in `cli.md` (§V.8) and the "Agentfile
  never names a substrate" / `POLICY substrate.*` request wording stays intact (§0.2).

Verify: `run` row shows implemented with the reference-translator qualifier; §0.8 line
verbatim above the table; `sign`/`verify` still `planned`; `rg -- '--substrate' cli.md` → 0.
</content>
