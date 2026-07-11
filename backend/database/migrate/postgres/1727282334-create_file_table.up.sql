CREATE TABLE IF NOT EXISTS files (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),

	uploaded_by uuid NOT NULL REFERENCES users (id),
    app_id uuid NOT NULL REFERENCES apps (id),
    folder_id uuid NOT NULL REFERENCES folders (id),

	name text NOT NULL,
	mime_type text NOT NULL,
	extension text NOT NULL,
	size bigint NOT NULL,
	is_public boolean NOT NULL DEFAULT false,
	public_id text NOT NULL UNIQUE,
	status file_status NOT NULL DEFAULT 'pending',
	
	metadata jsonb NOT NULL DEFAULT '{}'::jsonb,

	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);