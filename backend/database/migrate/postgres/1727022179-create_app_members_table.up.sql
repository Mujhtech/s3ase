CREATE TABLE IF NOT EXISTS app_members (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),

	user_id CHAR(26) NOT NULL REFERENCES users (id),
    app_id CHAR(26) NOT NULL REFERENCES apps (id),
    role app_role NOT NULL DEFAULT 'owner',

	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);