CREATE TABLE IF NOT EXISTS api_keys (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),

	app_id uuid NOT NULL REFERENCES apps (id),
    created_by uuid NOT NULL REFERENCES users (id),
	name TEXT NOT NULL,
	description TEXT NULL DEFAULT NULL,
    access api_key_access NOT NULL DEFAULT 'none',
    expired_at BIGINT NULL DEFAULT NULL,
	last_used TIMESTAMP NULL DEFAULT NULL,

	metadata jsonb NOT NULL DEFAULT '{}'::jsonb,

	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);