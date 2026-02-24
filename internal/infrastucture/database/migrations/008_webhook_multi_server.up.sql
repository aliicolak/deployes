-- Migration: Add multi-server support for webhooks
-- Creates a junction table for webhook-server relationships

-- Create junction table for webhook-server many-to-many relationship
CREATE TABLE IF NOT EXISTS webhook_servers (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    webhook_id VARCHAR(36) NOT NULL,
    server_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (webhook_id) REFERENCES webhooks(id) ON DELETE CASCADE,
    FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE,
    UNIQUE(webhook_id, server_id)
);

-- Migrate existing data from webhooks.server_id to webhook_servers
INSERT INTO webhook_servers (id, webhook_id, server_id, created_at)
SELECT 
    gen_random_uuid()::text as id,
    id as webhook_id,
    server_id,
    created_at
FROM webhooks
WHERE server_id IS NOT NULL AND server_id != ''
ON CONFLICT DO NOTHING;

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_webhook_servers_webhook_id ON webhook_servers(webhook_id);
CREATE INDEX IF NOT EXISTS idx_webhook_servers_server_id ON webhook_servers(server_id);

-- Remove server_id column from webhooks table (PostgreSQL supports DROP COLUMN)
ALTER TABLE webhooks DROP COLUMN IF EXISTS server_id;

