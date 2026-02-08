CREATE TABLE IF NOT EXISTS integration_download_events (
  id BIGSERIAL PRIMARY KEY,
  integration_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS integration_download_events_id_created_idx
  ON integration_download_events (integration_id, created_at DESC);
