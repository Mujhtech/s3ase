ALTER TABLE folders DROP CONSTRAINT IF EXISTS folders_parent_id_fkey;
ALTER TABLE folders
  ADD CONSTRAINT folders_parent_id_fkey
  FOREIGN KEY (parent_id) REFERENCES folders (id) ON DELETE CASCADE;

ALTER TABLE files DROP CONSTRAINT IF EXISTS files_folder_id_fkey;
ALTER TABLE files ALTER COLUMN folder_id DROP NOT NULL;
ALTER TABLE files
  ADD CONSTRAINT files_folder_id_fkey
  FOREIGN KEY (folder_id) REFERENCES folders (id) ON DELETE SET NULL;

CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique_active
  ON users (LOWER(email)) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS apps_slug_unique_active
  ON apps (slug) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS app_members_user_app_unique_active
  ON app_members (user_id, app_id) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS domains_app_unique_active
  ON domains (app_id) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS files_app_folder_active_idx
  ON files (app_id, folder_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS folders_app_active_idx
  ON folders (app_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS api_keys_app_active_idx
  ON api_keys (app_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS webhooks_app_active_idx
  ON webhooks (app_id) WHERE deleted_at IS NULL;
