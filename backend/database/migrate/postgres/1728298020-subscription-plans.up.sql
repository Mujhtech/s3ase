ALTER TABLE app_subscriptions
  ADD COLUMN IF NOT EXISTS plan TEXT NOT NULL DEFAULT 'basic',
  ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'active';

CREATE UNIQUE INDEX IF NOT EXISTS app_subscriptions_app_active_idx
  ON app_subscriptions (app_id) WHERE deleted_at IS NULL;
