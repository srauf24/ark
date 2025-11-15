-- +migrate Up
-- Create all base tables for Garden Journal MVP

-- ========== PLANTS TABLE ==========
CREATE TABLE plants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    species TEXT NOT NULL,
    location TEXT,
    planted_date TIMESTAMPTZ,
    notes TEXT,
    metadata JSONB,
    sort_order SERIAL
);

-- ========== OBSERVATIONS TABLE ==========
CREATE TABLE observations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id TEXT NOT NULL,
    plant_id UUID NOT NULL REFERENCES plants(id) ON DELETE CASCADE,
    date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    height_cm DECIMAL CHECK (height_cm >= 0),
    notes TEXT,
    sort_order SERIAL
);

-- ========== WEATHER SNAPSHOTS TABLE ==========
CREATE TABLE weather_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    date TIMESTAMPTZ NOT NULL,
    city TEXT NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    temp_max REAL,
    precip_mm REAL,
    sunshine_hrs REAL
);

-- ========== INDEXES ==========
CREATE UNIQUE INDEX plants_unique_name_per_user ON plants(user_id, name);
CREATE INDEX idx_plants_user_id ON plants(user_id);
CREATE INDEX idx_plants_user_sort ON plants(user_id, sort_order);
CREATE INDEX idx_observations_plant_sort ON observations(plant_id, sort_order);
CREATE INDEX idx_observations_user_plant_date_desc
    ON observations(user_id, plant_id, date DESC);

-- ========== TRIGGERS ==========
CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_plants
    BEFORE UPDATE ON plants
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();

CREATE TRIGGER set_updated_at_observations
    BEFORE UPDATE ON observations
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();
