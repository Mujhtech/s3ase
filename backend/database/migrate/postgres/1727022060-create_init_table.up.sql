CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE authentication_method AS ENUM ('google', 'github', 'password');

CREATE TYPE token_type AS ENUM ('bearer', 'basic', 'api_key');

CREATE TYPE api_key_access AS ENUM ('none', 'read', 'write', 'full');

CREATE TYPE app_role AS ENUM ('owner', 'member');

CREATE TYPE domain_status AS ENUM ('pending', 'review', 'verified', 'failed');

CREATE TYPE domain_verification_type AS ENUM ('dns', 'cname');

CREATE TYPE file_status AS ENUM ('started', 'pending', 'uploading', 'completed', 'cancelled', 'failed');