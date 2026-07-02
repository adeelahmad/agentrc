---
type: plan
story: S2-05
scope: "tests only"
---
# S2-05 Test Plan — CLI docs table (RED)

Tests only. Each check is a function in `scripts/verify-sprint2.sh`. Every bullet FAILS now
because `cli.md` still lists `run` as `planned` and carries `--substrate`.

## T25 — CLI docs table

- [ ] `scripts/verify-sprint2.sh::t25_run_implemented` — assert the `run` row status reads
      "implemented" with the "reference translators" qualifier. RED now: FAILS — `cli.md:103`
      says `planned`.
- [ ] `scripts/verify-sprint2.sh::t25_positioning_line_verbatim` — assert the §0.8 line
      "Reference translators — a proof of concept until platforms read `org.agentrc.*` labels
      natively. Not production runners." appears verbatim above the table. RED now: FAILS — absent.
- [ ] `scripts/verify-sprint2.sh::t25_sign_verify_stay_planned` — assert `sign` and `verify`
      rows still read `planned`. RED now: green today; regression guard against flipping them.
- [ ] `scripts/verify-sprint2.sh::t25_no_substrate_in_cli` — assert `rg -- '--substrate' cli.md`
      count == 0. RED now: FAILS — `cli.md:86/103/120` carry `--substrate`.

## Suite guard

- [ ] `scripts/verify-sprint2.sh::v8_terminology_split` — assert `--substrate` gone from
      CLI/docs (→0) and `POLICY substrate.` intact in spec (§V.8). RED now: FAILS pre-rename;
      shared guard across S2-03/S2-05.
</content>
