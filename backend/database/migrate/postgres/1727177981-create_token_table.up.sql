CREATE TABLE IF NOT EXISTS tokens (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),

	is_app BOOLEAN NOT NULL DEFAULT FALSE,
    type  token_type NOT NULL DEFAULT 'bearer',
    value VARCHAR(255) NOT NULL,
    expired_at BIGINT NOT NULL,
    issued_at BIGINT NOT NULL,
    metadata jsonb NOT NULL DEFAULT '{}'::jsonb,

	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);