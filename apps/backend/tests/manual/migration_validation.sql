-- ============================================
-- ARK Schema Migration Validation Script
-- Run after: task migrations:up
-- ============================================

\echo '=== Checking Extensions ==='
\dx

\echo '\n=== Checking Functions ==='
\df trigger_set_timestamp

\echo '\n=== Checking Tables ==='
\dt

\echo '\n=== Assets Table Structure ==='
\d assets

\echo '\n=== Asset Logs Table Structure ==='
\d asset_logs

\echo '\n=== All Indexes ==='
\di

\echo '\n=== Assets Indexes ==='
\di assets*

\echo '\n=== Asset Logs Indexes ==='
\di asset_logs*

\echo '\n=== Triggers ==='
SELECT
  tgname AS trigger_name,
  tgrelid::regclass AS table_name,
  tgenabled AS enabled
FROM pg_trigger
WHERE tgname LIKE '%timestamp%'
ORDER BY tgrelid::regclass, tgname;

\echo '\n=== Foreign Key Constraints ==='
SELECT
  conname AS constraint_name,
  conrelid::regclass AS table_name,
  confrelid::regclass AS referenced_table,
  confdeltype AS on_delete_action
FROM pg_constraint
WHERE contype = 'f'
AND conrelid IN ('assets'::regclass, 'asset_logs'::regclass);

\echo '\n=== Generated Columns ==='
SELECT
  table_name,
  column_name,
  is_generated,
  generation_expression
FROM information_schema.columns
WHERE table_name IN ('assets', 'asset_logs')
AND is_generated = 'ALWAYS';

\echo '\n=== Validation Complete ==='
