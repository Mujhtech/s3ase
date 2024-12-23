CREATE TABLE IF NOT EXISTS webhooks (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),

	app_id uuid NOT NULL REFERENCES apps (id),
    created_by uuid NOT NULL REFERENCES users (id),
	name TEXT NOT NULL,
	description TEXT NULL DEFAULT NULL,
    url TEXT NOT NULL,

	metadata jsonb NOT NULL DEFAULT '{}'::jsonb,

	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);