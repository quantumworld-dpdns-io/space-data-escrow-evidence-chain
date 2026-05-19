# Release Evidence Index

Use this index to collect all mandatory v1 release evidence.

## Required Workflow Runs
- CI (`.github/workflows/ci.yml`)
- Security (`.github/workflows/security.yml`)
- Release Readiness (`.github/workflows/release-readiness.yml`)
- Release (`.github/workflows/release.yml`)
- Post-Deploy Verify (`.github/workflows/post-deploy-verify.yml`)

## Required Artifacts
- API binary
- SBOM (`artifacts/sbom.spdx.json`)
- Provenance (`artifacts/provenance.json`)
- Release notes
- Sign-off bundle (from template)
