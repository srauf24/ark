---- tern migration up

-- Enable UUID generation extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable trigram fuzzy text search extension
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Create or replace trigger function for auto-updating updated_at
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

---- tern migration down

-- TODO: Rollback content will be added in subsequent steps
