# v1.0.0 Release Protocol

## Tagging Rule
- Release tags must follow `vMAJOR.MINOR.PATCH`.
- v1 launch tag: `v1.0.0`.

## Pre-Tag Requirements
1. `go test ./...` green on main.
2. CI: lint, unit-test, benchmark-regression, robot, security all green.
3. Release Candidate checklist completed.
4. Governance sign-off document completed.

## Tag and Release Steps
1. Create annotated tag:
   - `git tag -a v1.0.0 -m "Release v1.0.0"`
2. Push tag:
   - `git push origin v1.0.0`
3. Verify `Release` workflow publishes assets.
4. Trigger `Post-Deploy Verify` workflow for target environment.

## Post-Release Requirements
1. Post-deploy smoke verification passes.
2. Incident/on-call channel notified.
3. Release notes and sign-off bundle archived.
