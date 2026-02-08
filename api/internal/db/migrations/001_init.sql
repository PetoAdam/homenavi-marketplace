CREATE TABLE IF NOT EXISTS integrations (
  id TEXT NOT NULL,
  version TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  manifest_url TEXT NOT NULL,
  manifest JSONB,
  image TEXT NOT NULL,
  images JSONB NOT NULL DEFAULT '[]',
  assets JSONB NOT NULL DEFAULT '{}',
  listen_path TEXT NOT NULL,
  repo_url TEXT,
  release_tag TEXT,
  publisher TEXT,
  verified BOOLEAN NOT NULL DEFAULT FALSE,
  latest BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id, version)
);

CREATE UNIQUE INDEX IF NOT EXISTS integrations_listen_path_unique
  ON integrations (listen_path)
  WHERE latest = TRUE;

CREATE INDEX IF NOT EXISTS integrations_latest_idx
  ON integrations (latest);
