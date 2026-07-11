ALTER TABLE api_keys
  ADD COLUMN IF NOT EXISTS key_hash TEXT,
  ADD COLUMN IF NOT EXISTS key_prefix TEXT,
  ADD COLUMN IF NOT EXISTS last_four TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS api_keys_key_hash_active_idx
  ON api_keys (key_hash) WHERE deleted_at IS NULL AND key_hash IS NOT NULL;
