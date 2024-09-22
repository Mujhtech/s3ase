CREATE TABLE IF NOT EXITS users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  email VARCHAR(255) NOT NULL,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  display_name VARCHAR(255) NULL,
  avatar_url TEXT NULL,
  authentication_method authentication_method NOT NULL DEFAULT 'password',
  password TEXT NULL,
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP DEFAULT NULL
);

