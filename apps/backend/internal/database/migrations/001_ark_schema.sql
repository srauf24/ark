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

-- Create assets table
CREATE TABLE assets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id TEXT NOT NULL,
  name TEXT NOT NULL,
  type TEXT,
  hostname TEXT,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Create index on user_id for security and multi-tenancy
CREATE INDEX idx_assets_user_id ON assets(user_id);

-- Create trigram index on name for fuzzy search
CREATE INDEX idx_assets_name_trgm ON assets USING GIN (name gin_trgm_ops);

-- Create partial index on type for filtering
CREATE INDEX idx_assets_type ON assets(type) WHERE type IS NOT NULL;

---- tern migration down

-- TODO: Rollback content will be added in subsequent steps
