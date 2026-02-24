-- Add project type and local path columns
ALTER TABLE projects ADD COLUMN IF NOT EXISTS type VARCHAR(10) DEFAULT 'github';
ALTER TABLE projects ADD COLUMN IF NOT EXISTS local_path TEXT;

-- Update existing projects to have type 'github'
UPDATE projects SET type = 'github' WHERE type IS NULL;

-- Make repo_url and branch nullable for local projects
ALTER TABLE projects ALTER COLUMN repo_url DROP NOT NULL;
ALTER TABLE projects ALTER COLUMN branch DROP NOT NULL;
