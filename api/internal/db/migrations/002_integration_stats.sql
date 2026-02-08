ALTER TABLE integrations
  ADD COLUMN IF NOT EXISTS downloads BIGINT NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS trending_score DOUBLE PRECISION NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS featured BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS integrations_featured_idx
  ON integrations (featured)
  WHERE featured = TRUE;

CREATE INDEX IF NOT EXISTS integrations_downloads_idx
  ON integrations (downloads DESC);

CREATE INDEX IF NOT EXISTS integrations_trending_idx
  ON integrations (trending_score DESC);
