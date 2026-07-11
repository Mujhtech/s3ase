DROP INDEX IF EXISTS app_subscriptions_app_active_idx;
ALTER TABLE app_subscriptions
  DROP COLUMN IF EXISTS status,
  DROP COLUMN IF EXISTS plan;
