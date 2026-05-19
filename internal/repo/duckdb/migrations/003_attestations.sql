CREATE TABLE IF NOT EXISTS attestations (
  evidence_id TEXT NOT NULL,
  signer TEXT NOT NULL,
  signature TEXT NOT NULL,
  algorithm TEXT NOT NULL,
  ts TIMESTAMP NOT NULL
);
