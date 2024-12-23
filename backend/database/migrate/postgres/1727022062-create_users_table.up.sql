CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  email VARCHAR(255) NOT NULL,
  email_verified BOOLEAN NOT NULL DEFAULT false,
  given_name VARCHAR(255) NOT NULL,
  display_name VARCHAR(255) NULL,
  avatar_url TEXT NULL,
  authentication_method authentication_method NOT NULL DEFAULT 'password',
  password TEXT NULL,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT NULL
);

