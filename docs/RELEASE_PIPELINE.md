# Release Pipeline

This repository now includes release and artifact workflows:

- `.github/workflows/artifacts.yml`: builds API binary on `main`, generates SBOM/provenance placeholders, uploads artifacts.
- `.github/workflows/release.yml`: triggers on `v*` tags, publishes GitHub Release with binary + SBOM + provenance + notes.
- `.github/workflows/changelog.yml`: generates a changelog draft artifact from recent commits.

Generated artifacts:

- `artifacts/evidence-api`
- `artifacts/sbom.spdx.json`
- `artifacts/provenance.json`
- `artifacts/release-notes.md`

Helper scripts:

- `scripts/release/generate_sbom.sh`
- `scripts/release/generate_provenance.sh`
