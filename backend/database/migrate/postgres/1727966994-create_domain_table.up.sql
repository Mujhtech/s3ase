CREATE TABLE IF NOT EXISTS domains (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),

	app_id uuid NOT NULL REFERENCES apps (id),
    created_by uuid NOT NULL REFERENCES users (id),
    domain TEXT NOT NULL UNIQUE,
	cname_record TEXT NOT NULL,
	txt_record TEXT NOT NULL,
	status domain_status NOT NULL DEFAULT 'pending',

	metadata jsonb NOT NULL DEFAULT '{}'::jsonb,

	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);