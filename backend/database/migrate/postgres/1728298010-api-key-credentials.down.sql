DROP INDEX IF EXISTS api_keys_key_hash_active_idx;
ALTER TABLE api_keys
  DROP COLUMN IF EXISTS last_four,
  DROP COLUMN IF EXISTS key_prefix,
  DROP COLUMN IF EXISTS key_hash;
