-- Migration 004: Sync Agent Schema with model.go
-- Added missing tables and fields for Phase 2/3 Agent features

-- 1. User API Keys
CREATE TABLE IF NOT EXISTS user_api_keys (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(100),
    description TEXT,
    scopes TEXT, -- Comma separated scopes
    rate_limit INTEGER DEFAULT 0,
    last_used_at TIMESTAMP,
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 2. Agent Push Subscriptions
CREATE TABLE IF NOT EXISTS agent_push_subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    push_type VARCHAR(50) NOT NULL,
    enabled BOOLEAN DEFAULT true,
    scope TEXT,
    webhook_url VARCHAR(500),
    secret VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 3. Agent Push Logs
CREATE TABLE IF NOT EXISTS agent_push_logs (
    id SERIAL PRIMARY KEY,
    subscription_id INTEGER NOT NULL REFERENCES agent_push_subscriptions(id) ON DELETE CASCADE,
    artifact_id INTEGER,
    payload TEXT,
    status VARCHAR(20), -- success, failed
    retry_count INTEGER DEFAULT 0,
    error_message TEXT,
    delivered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 4. Agent Sessions & Artifacts
CREATE TABLE IF NOT EXISTS agent_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scenario VARCHAR(100),
    factory_id INTEGER REFERENCES factories(id),
    language VARCHAR(10),
    input_snapshot TEXT,
    status VARCHAR(20),
    trace_id VARCHAR(100) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS agent_artifacts (
    id SERIAL PRIMARY KEY,
    session_id INTEGER NOT NULL REFERENCES agent_sessions(id) ON DELETE CASCADE,
    artifact_type VARCHAR(50),
    title VARCHAR(200),
    summary TEXT,
    result_json TEXT,
    risk_level VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS agent_evidence_links (
    id SERIAL PRIMARY KEY,
    artifact_id INTEGER NOT NULL REFERENCES agent_artifacts(id) ON DELETE CASCADE,
    evidence_type VARCHAR(50),
    source_table VARCHAR(100),
    source_id INTEGER,
    excerpt TEXT,
    score DECIMAL(5,4),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 5. Agent Usage Tracking
CREATE TABLE IF NOT EXISTS agent_usage (
    id SERIAL PRIMARY KEY,
    session_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    scenario VARCHAR(50),
    model VARCHAR(100),
    prompt_tokens INTEGER DEFAULT 0,
    completion_tokens INTEGER DEFAULT 0,
    total_tokens INTEGER DEFAULT 0,
    response_time_ms BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 6. Manuals & Runtime Snapshots
CREATE TABLE IF NOT EXISTS manual_documents (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    equipment_type_id INTEGER REFERENCES equipment_types(id),
    file_path VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS manual_chunks (
    id SERIAL PRIMARY KEY,
    document_id INTEGER NOT NULL REFERENCES manual_documents(id) ON DELETE CASCADE,
    section_title VARCHAR(200),
    content TEXT,
    page_number INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS equipment_runtime_snapshots (
    id SERIAL PRIMARY KEY,
    equipment_id INTEGER NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    status VARCHAR(20),
    load_factor DECIMAL(5,2),
    snapshot_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Indexes for new tables
CREATE INDEX IF NOT EXISTS idx_api_keys_user ON user_api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_key ON user_api_keys(key);
CREATE INDEX IF NOT EXISTS idx_push_subs_user ON agent_push_subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_push_logs_sub ON agent_push_logs(subscription_id);
CREATE INDEX IF NOT EXISTS idx_agent_sessions_user ON agent_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_agent_artifacts_session ON agent_artifacts(session_id);
CREATE INDEX IF NOT EXISTS idx_runtime_snapshots_equip ON equipment_runtime_snapshots(equipment_id, snapshot_date);
