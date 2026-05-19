# Release Candidate Checklist

- [ ] `go test ./...` passes on release candidate commit.
- [ ] CI jobs pass: lint, unit-test, benchmark-regression, robot, security.
- [ ] OpenAPI file is valid and matches runtime behavior.
- [ ] Release artifacts generated: binary, SBOM, provenance, release notes.
- [ ] OWASP A01-A10 suite passes.
- [ ] DAST-lite checks pass against candidate deployment.
- [ ] Runtime security policies and schema are present and reviewed.
- [ ] Rollback instructions verified.
