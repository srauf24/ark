---- tern migration up

-- Enable UUID generation extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable trigram fuzzy text search extension
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

---- tern migration down

-- TODO: Rollback content will be added in subsequent steps
