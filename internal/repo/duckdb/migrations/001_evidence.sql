CREATE TABLE IF NOT EXISTS evidence_records (
  id TEXT PRIMARY KEY,
  external_id TEXT NOT NULL,
  source TEXT NOT NULL,
  type TEXT NOT NULL,
  hash TEXT NOT NULL,
  payload_json TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL
);
