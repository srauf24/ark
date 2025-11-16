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

-- Create trigger to auto-update updated_at on assets table
CREATE TRIGGER set_assets_timestamp
  BEFORE UPDATE ON assets
  FOR EACH ROW
  EXECUTE FUNCTION trigger_set_timestamp();

-- Create asset_logs table
CREATE TABLE asset_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
  user_id TEXT NOT NULL,
  content TEXT NOT NULL,
  tags TEXT[],
  content_vector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', content)) STORED,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Create index on user_id for security and multi-tenancy
CREATE INDEX idx_asset_logs_user_id ON asset_logs(user_id);

-- Create index on asset_id for efficient joins
CREATE INDEX idx_asset_logs_asset_id ON asset_logs(asset_id);

-- Create index on created_at for chronological ordering
CREATE INDEX idx_asset_logs_created_at ON asset_logs(created_at DESC);

-- Create GIN index on content_vector for full-text search
CREATE INDEX idx_asset_logs_content_vector ON asset_logs USING GIN (content_vector);

-- Create partial GIN index on tags for tag filtering
CREATE INDEX idx_asset_logs_tags ON asset_logs USING GIN (tags) WHERE tags IS NOT NULL;

-- Create trigger to auto-update updated_at on asset_logs table
CREATE TRIGGER set_asset_logs_timestamp
  BEFORE UPDATE ON asset_logs
  FOR EACH ROW
  EXECUTE FUNCTION trigger_set_timestamp();

---- tern migration down

-- TODO: Rollback content will be added in subsequent steps
