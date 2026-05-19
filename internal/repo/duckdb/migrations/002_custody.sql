CREATE TABLE IF NOT EXISTS custody_events (
  evidence_id TEXT NOT NULL,
  actor TEXT NOT NULL,
  action TEXT NOT NULL,
  note TEXT,
  ts TIMESTAMP NOT NULL
);
