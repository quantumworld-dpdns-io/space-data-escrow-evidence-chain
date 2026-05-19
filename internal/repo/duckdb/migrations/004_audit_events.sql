CREATE TABLE IF NOT EXISTS audit_events (
  event_id BIGINT,
  entry TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL
);
