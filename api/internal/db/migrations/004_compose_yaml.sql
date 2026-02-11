ALTER TABLE integrations
  ADD COLUMN IF NOT EXISTS compose_file TEXT;
