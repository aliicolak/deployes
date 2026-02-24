-- Rollback: Revert multi-server webhook support

-- Add server_id column back to webhooks table
ALTER TABLE webhooks ADD COLUMN IF NOT EXISTS server_id VARCHAR(36);

-- Copy data back from junction table (using the first server_id)
UPDATE webhooks w
SET server_id = (
    SELECT ws.server_id 
    FROM webhook_servers ws 
    WHERE ws.webhook_id = w.id 
    LIMIT 1
);

-- Drop junction table and indexes
DROP INDEX IF EXISTS idx_webhook_servers_webhook_id;
DROP INDEX IF EXISTS idx_webhook_servers_server_id;
DROP TABLE IF EXISTS webhook_servers;

