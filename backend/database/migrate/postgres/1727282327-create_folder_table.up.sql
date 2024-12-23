CREATE TABLE IF NOT EXISTS folders (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),

    app_id uuid NOT NULL REFERENCES apps (id),
	created_by uuid NOT NULL REFERENCES users (id),
    parent_id uuid NULL REFERENCES users (id) DEFAULT NULL,

	name TEXT NOT NULL,
	description TEXT NULL DEFAULT NULL,

	metadata jsonb NOT NULL DEFAULT '{}'::jsonb,

	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);