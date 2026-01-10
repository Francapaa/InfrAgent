-- Clients table
CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL,
    google_id VARCHAR(150)

    -- Credentials (el cliente usa esto)
    api_key_hash VARCHAR(255) NOT NULL,  -- bcrypt hash para reportar eventos
    webhook_secret VARCHAR(255) NOT NULL, -- Para firmar requests del agente
    
    -- Webhook configuration
    webhook_url TEXT NOT NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Client configurations
CREATE TABLE IF NOT EXISTS client_configs (
    id SERIAL PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    max_restarts_per_hour INT NOT NULL DEFAULT 3,
    allowed_actions JSONB NOT NULL DEFAULT '["restart", "notify", "wait"]'::jsonb,
    notify_on_nth_restart INT NOT NULL DEFAULT 3,
    cooldown_minutes INT NOT NULL DEFAULT 5,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(client_id)
);

-- Agents table (1 agent per client)
CREATE TABLE IF NOT EXISTS agents (
    id UUID PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    state VARCHAR(50) NOT NULL DEFAULT 'idle', -- idle, analyzing, acting, cooldown
    last_tick_at TIMESTAMP,
    cooldown_until TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(client_id) -- 1 agent per client
);

-- Events (incidents reported by client's SDK)
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    agent_id VARCHAR(255) NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL, -- app_down, high_cpu, error_spike, etc
    service VARCHAR(255) NOT NULL, -- api, db, worker, etc
    severity VARCHAR(50) NOT NULL, -- info, warning, critical
    data JSONB NOT NULL DEFAULT '{}'::jsonb,
    processed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Actions (decisions made by the agent)
CREATE TABLE IF NOT EXISTS actions (
    id UUID PRIMARY KEY,
    agent_id VARCHAR(255) NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    client_id VARCHAR(255) NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL, -- restart, notify, wait, scale, etc
    target VARCHAR(255) NOT NULL, -- api, db, etc
    params JSONB NOT NULL DEFAULT '{}'::jsonb,
    reasoning TEXT NOT NULL,
    confidence FLOAT NOT NULL,
    status VARCHAR(50) NOT NULL, -- pending, success, failed
    result JSONB NOT NULL DEFAULT '{}'::jsonb,
    executed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Notifications sent to clients
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    action_id VARCHAR(255) REFERENCES actions(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL, -- email, slack, webhook
    recipient VARCHAR(255) NOT NULL,
    subject VARCHAR(255),
    body TEXT NOT NULL,
    sent_at TIMESTAMP,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, sent, failed
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_events_client_pending ON events(client_id, processed_at) WHERE processed_at IS NULL;
CREATE INDEX idx_events_agent_pending ON events(agent_id, processed_at) WHERE processed_at IS NULL;
CREATE INDEX idx_events_created ON events(created_at DESC);
CREATE INDEX idx_actions_agent_type_created ON actions(agent_id, type, created_at DESC);
CREATE INDEX idx_actions_client_created ON actions(client_id, created_at DESC);
CREATE INDEX idx_actions_created ON actions(created_at DESC);
CREATE INDEX idx_agents_client ON agents(client_id);
CREATE INDEX idx_notifications_client ON notifications(client_id, created_at DESC);

-- Function to auto-create agent when client is created
CREATE OR REPLACE FUNCTION create_agent_for_client()
RETURNS TRIGGER AS $
BEGIN
    INSERT INTO agents (id, client_id, state)
    VALUES ('agent-' || NEW.id, NEW.id, 'idle');
    
    INSERT INTO client_configs (client_id)
    VALUES (NEW.id);
    
    RETURN NEW;
END;
$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_create_agent_for_client
    AFTER INSERT ON clients
    FOR EACH ROW
    EXECUTE FUNCTION create_agent_for_client();

-- MOCK DATA ESTO LUEGO EN PRODUCCION LO SACAMOS
-- INSERT INTO clients (id, email, company_name, api_key_hash, webhook_secret, webhook_url) VALUES
--     ('client-001', 'test@example.com', 'Test Company', 
--      '$2a$10$test_hash_here', 'whsec_test_secret_123',
--      'http://localhost:8080/webhooks/agent')
-- ON CONFLICT (id) DO NOTHING;