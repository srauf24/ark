-- ============================================
-- ARK Schema Migration Test Data Script
-- Run after: task migrations:up
-- Cleanup: DELETE FROM assets WHERE user_id LIKE 'test-%';
-- ============================================

\echo '=== Test 1: Insert Assets with Defaults ==='
INSERT INTO assets (user_id, name)
VALUES
  ('test-user-1', 'Homelab Server'),
  ('test-user-1', 'NAS Device')
RETURNING id, name, created_at, updated_at;

\echo '\n=== Test 2: Insert Asset with All Fields ==='
INSERT INTO assets (user_id, name, type, hostname, metadata)
VALUES (
  'test-user-1',
  'Production VM',
  'vm',
  'prod-vm-01.local',
  '{"ip": "192.168.1.50", "cpu": "4 cores", "ram": "16GB", "os": "Ubuntu 22.04"}'::jsonb
)
RETURNING *;

\echo '\n=== Test 3: Insert Logs with Generated content_vector ==='
WITH asset AS (
  SELECT id FROM assets WHERE name = 'Homelab Server' AND user_id = 'test-user-1'
)
INSERT INTO asset_logs (asset_id, user_id, content, tags)
SELECT
  id,
  'test-user-1',
  content,
  tags
FROM asset, (VALUES
  ('Fixed nginx by updating /etc/nginx/nginx.conf and restarting service', ARRAY['nginx', 'fix', 'web-server']),
  ('Installed Docker and configured daemon.json for logging', ARRAY['docker', 'install', 'config']),
  ('Updated system packages with apt update && apt upgrade', ARRAY['system', 'update', 'maintenance'])
) AS logs(content, tags)
RETURNING id, content, tags, content_vector;

\echo '\n=== Test 4: Verify Trigger - Update Asset ==='
-- Note original updated_at
SELECT id, name, updated_at FROM assets WHERE name = 'Homelab Server' AND user_id = 'test-user-1';

-- Wait 2 seconds
SELECT pg_sleep(2);

-- Update asset
UPDATE assets
SET name = 'Homelab Server - Updated'
WHERE name = 'Homelab Server' AND user_id = 'test-user-1'
RETURNING id, name, created_at, updated_at;

\echo '\n=== Test 5: Verify Trigger - Update Log ==='
WITH log AS (
  SELECT id, content, updated_at
  FROM asset_logs
  WHERE content LIKE '%nginx%' AND user_id = 'test-user-1'
  LIMIT 1
)
SELECT * FROM log;

-- Wait 2 seconds
SELECT pg_sleep(2);

-- Update log
WITH log AS (
  SELECT id FROM asset_logs WHERE content LIKE '%nginx%' AND user_id = 'test-user-1' LIMIT 1
)
UPDATE asset_logs
SET content = 'Fixed nginx - UPDATED CONTENT'
FROM log
WHERE asset_logs.id = log.id
RETURNING asset_logs.id, asset_logs.content, asset_logs.updated_at, asset_logs.content_vector;

\echo '\n=== Test 6: Full-Text Search ==='
SELECT
  content,
  ts_rank(content_vector, query) AS rank
FROM asset_logs, to_tsquery('english', 'nginx | docker') AS query
WHERE content_vector @@ query
AND user_id = 'test-user-1'
ORDER BY rank DESC;

\echo '\n=== Test 7: Tag Search ==='
SELECT content, tags
FROM asset_logs
WHERE tags @> ARRAY['nginx']
AND user_id = 'test-user-1';

\echo '\n=== Test 8: Fuzzy Asset Name Search ==='
SELECT
  name,
  similarity(name, 'Homelb') AS sim
FROM assets
WHERE name % 'Homelb'
AND user_id = 'test-user-1'
ORDER BY sim DESC;

\echo '\n=== Test 9: Cascade Delete - Delete Asset ==='
WITH asset AS (
  SELECT id FROM assets WHERE name LIKE '%Homelab%' AND user_id = 'test-user-1' LIMIT 1
)
SELECT
  (SELECT COUNT(*) FROM asset_logs WHERE asset_id = asset.id) AS logs_before_delete
FROM asset;

-- Delete asset
WITH asset AS (
  SELECT id FROM assets WHERE name LIKE '%Homelab%' AND user_id = 'test-user-1' LIMIT 1
)
DELETE FROM assets WHERE id = (SELECT id FROM asset);

-- Check logs are cascade deleted
SELECT COUNT(*) AS logs_after_delete
FROM asset_logs
WHERE user_id = 'test-user-1';

\echo '\n=== Test 10: Foreign Key Violation ==='
\echo 'Attempting to insert log with invalid asset_id (should fail)...'
INSERT INTO asset_logs (asset_id, user_id, content)
VALUES ('00000000-0000-0000-0000-000000000000', 'test-user-1', 'Invalid log');
-- Expected: ERROR - foreign key constraint violation

\echo '\n=== Cleanup Test Data ==='
DELETE FROM assets WHERE user_id LIKE 'test-%';

\echo '\n=== Verify Cleanup ==='
SELECT
  (SELECT COUNT(*) FROM assets WHERE user_id LIKE 'test-%') AS assets_count,
  (SELECT COUNT(*) FROM asset_logs WHERE user_id LIKE 'test-%') AS logs_count;

\echo '\n=== Test Data Script Complete ==='
