CREATE TABLE IF NOT EXISTS app_subscriptions (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),

	app_id uuid NOT NULL REFERENCES apps (id),
    

	metadata jsonb NOT NULL DEFAULT '{}'::jsonb,

	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);