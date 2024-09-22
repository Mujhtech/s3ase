CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE authentication_method AS ENUM ('google', 'github', 'password');

CREATE TYPE app_role AS ENUM ('owner', 'member');