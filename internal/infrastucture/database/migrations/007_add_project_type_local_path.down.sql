-- Rollback: Remove project type and local path columns
-- Note: This will lose data for local projects!

-- First, delete local projects since they can't work without these columns
DELETE FROM projects WHERE type = 'local';

-- Remove columns
ALTER TABLE projects DROP COLUMN IF EXISTS local_path;
ALTER TABLE projects DROP COLUMN IF EXISTS type;

-- Make repo_url and branch NOT NULL again
ALTER TABLE projects ALTER COLUMN repo_url SET NOT NULL;
ALTER TABLE projects ALTER COLUMN branch SET NOT NULL;
