DROP INDEX IF EXISTS webhooks_app_active_idx;
DROP INDEX IF EXISTS api_keys_app_active_idx;
DROP INDEX IF EXISTS folders_app_active_idx;
DROP INDEX IF EXISTS files_app_folder_active_idx;
DROP INDEX IF EXISTS domains_app_unique_active;
DROP INDEX IF EXISTS app_members_user_app_unique_active;
DROP INDEX IF EXISTS apps_slug_unique_active;
DROP INDEX IF EXISTS users_email_unique_active;

ALTER TABLE files DROP CONSTRAINT IF EXISTS files_folder_id_fkey;
ALTER TABLE files
  ADD CONSTRAINT files_folder_id_fkey
  FOREIGN KEY (folder_id) REFERENCES folders (id);

ALTER TABLE folders DROP CONSTRAINT IF EXISTS folders_parent_id_fkey;
ALTER TABLE folders
  ADD CONSTRAINT folders_parent_id_fkey
  FOREIGN KEY (parent_id) REFERENCES users (id);
