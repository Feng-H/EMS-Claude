-- EMS Database Schema (Final Unified Version)
-- PostgreSQL 15+

-- Drop existing tables (for clean setup)
DROP TABLE IF EXISTS agent_push_subscriptions CASCADE;
DROP TABLE IF EXISTS agent_messages CASCADE;
DROP TABLE IF EXISTS agent_conversations CASCADE;
DROP TABLE IF EXISTS agent_experiences CASCADE;
DROP TABLE IF EXISTS agent_knowledges CASCADE;
DROP TABLE IF EXISTS agent_skills CASCADE;
DROP TABLE IF EXISTS agent_usage CASCADE;
DROP TABLE IF EXISTS agent_evidence_links CASCADE;
DROP TABLE IF EXISTS agent_artifacts CASCADE;
DROP TABLE IF EXISTS agent_sessions CASCADE;
DROP TABLE IF EXISTS equipment_runtime_snapshots CASCADE;
DROP TABLE IF EXISTS repair_cost_details CASCADE;
DROP TABLE IF EXISTS manual_chunks CASCADE;
DROP TABLE IF EXISTS manual_documents CASCADE;
DROP TABLE IF EXISTS spare_part_consumptions CASCADE;
DROP TABLE IF EXISTS spare_part_inventory CASCADE;
DROP TABLE IF EXISTS spare_parts CASCADE;
DROP TABLE IF EXISTS knowledge_articles CASCADE;
DROP TABLE IF EXISTS maintenance_tasks CASCADE;
DROP TABLE IF EXISTS maintenance_plans CASCADE;
DROP TABLE IF EXISTS repair_logs CASCADE;
DROP TABLE IF EXISTS repair_orders CASCADE;
DROP TABLE IF EXISTS inspection_records CASCADE;
DROP TABLE IF EXISTS inspection_tasks CASCADE;
DROP TABLE IF EXISTS inspection_items CASCADE;
DROP TABLE IF EXISTS inspection_templates CASCADE;
DROP TABLE IF EXISTS equipment CASCADE;
DROP TABLE IF EXISTS equipment_types CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS workshops CASCADE;
DROP TABLE IF EXISTS factories CASCADE;
DROP TABLE IF EXISTS bases CASCADE;

-- =====================================================
-- Custom Types & Organization
-- =====================================================

CREATE TYPE user_role AS ENUM ('admin', 'supervisor', 'engineer', 'maintenance', 'operator');
CREATE TYPE user_approval_status AS ENUM ('pending', 'approved', 'rejected');
CREATE TYPE repair_status AS ENUM ('pending', 'assigned', 'in_progress', 'testing', 'confirmed', 'audited', 'closed');

CREATE TABLE bases (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE factories (
    id SERIAL PRIMARY KEY,
    base_id INTEGER NOT NULL REFERENCES bases(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    UNIQUE(base_id, code)
);

CREATE TABLE workshops (
    id SERIAL PRIMARY KEY,
    factory_id INTEGER NOT NULL REFERENCES factories(id),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    UNIQUE(factory_id, code)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    role user_role NOT NULL DEFAULT 'operator',
    factory_id INTEGER REFERENCES factories(id),
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    approval_status user_approval_status DEFAULT 'approved',
    must_change_password BOOLEAN DEFAULT false,
    first_login BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Equipment & Financials
-- =====================================================

CREATE TABLE equipment_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    inspection_template_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE equipment (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    type_id INTEGER NOT NULL REFERENCES equipment_types(id),
    workshop_id INTEGER NOT NULL REFERENCES workshops(id),
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    spec TEXT,
    purchase_price DECIMAL(12,2) DEFAULT 0,
    purchase_date DATE,
    service_life_years INTEGER DEFAULT 10,
    scrap_value DECIMAL(12,2) DEFAULT 0,
    hourly_loss DECIMAL(10,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'running',
    dedicated_maintenance_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Tasks: Inspection, Repair, Maintenance
-- =====================================================

CREATE TABLE inspection_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    equipment_type_id INTEGER NOT NULL REFERENCES equipment_types(id)
);

CREATE TABLE inspection_items (
    id SERIAL PRIMARY KEY,
    template_id INTEGER NOT NULL REFERENCES inspection_templates(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    method TEXT,
    criteria TEXT,
    sequence_order INTEGER DEFAULT 0
);

CREATE TABLE inspection_tasks (
    id SERIAL PRIMARY KEY,
    equipment_id INTEGER NOT NULL REFERENCES equipment(id),
    template_id INTEGER NOT NULL REFERENCES inspection_templates(id),
    assigned_to INTEGER NOT NULL REFERENCES users(id),
    scheduled_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE inspection_records (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES inspection_tasks(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES inspection_items(id),
    result VARCHAR(10) NOT NULL, -- OK/NG
    remark TEXT,
    photo_url VARCHAR(500)
);

CREATE TABLE repair_orders (
    id SERIAL PRIMARY KEY,
    equipment_id INTEGER NOT NULL REFERENCES equipment(id),
    fault_description TEXT NOT NULL,
    reporter_id INTEGER NOT NULL REFERENCES users(id),
    assigned_to INTEGER REFERENCES users(id),
    status repair_status DEFAULT 'pending',
    priority INTEGER DEFAULT 3,
    fault_code VARCHAR(50),
    solution TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    confirmed_at TIMESTAMP,
    audited_at TIMESTAMP,
    closed_at TIMESTAMP,
    photos TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE repair_logs (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES repair_orders(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    action VARCHAR(50),
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE maintenance_plans (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    equipment_type_id INTEGER NOT NULL REFERENCES equipment_types(id) ON DELETE CASCADE,
    level INTEGER NOT NULL,
    cycle_days INTEGER NOT NULL,
    flexible_days INTEGER DEFAULT 3,
    work_hours DECIMAL(5,2) DEFAULT 0
);

CREATE TABLE maintenance_plan_items (
    id SERIAL PRIMARY KEY,
    plan_id INTEGER NOT NULL REFERENCES maintenance_plans(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    content TEXT,
    standard TEXT,
    sequence_order INTEGER DEFAULT 0
);

CREATE TABLE maintenance_tasks (
    id SERIAL PRIMARY KEY,
    plan_id INTEGER NOT NULL REFERENCES maintenance_plans(id) ON DELETE CASCADE,
    equipment_id INTEGER NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    assigned_to INTEGER NOT NULL REFERENCES users(id),
    scheduled_date DATE NOT NULL,
    due_date DATE,
    status VARCHAR(20) DEFAULT 'pending',
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    actual_hours DECIMAL(5,2),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE maintenance_records (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES maintenance_tasks(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL, -- Logical reference to plan_item
    result VARCHAR(20) NOT NULL, -- 合格/不合格
    remark TEXT,
    photo_url VARCHAR(500)
);

-- =====================================================
-- Inventory & Knowledge
-- =====================================================

CREATE TABLE spare_parts (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    specification TEXT,
    unit VARCHAR(20),
    category VARCHAR(50),
    safety_stock INTEGER DEFAULT 0,
    factory_id INTEGER REFERENCES factories(id)
);

CREATE TABLE spare_part_inventory (
    id SERIAL PRIMARY KEY,
    spare_part_id INTEGER NOT NULL REFERENCES spare_parts(id) ON DELETE CASCADE,
    factory_id INTEGER NOT NULL REFERENCES factories(id) ON DELETE CASCADE,
    quantity INTEGER DEFAULT 0,
    UNIQUE(spare_part_id, factory_id)
);

CREATE TABLE spare_part_consumptions (
    id SERIAL PRIMARY KEY,
    spare_part_id INTEGER NOT NULL REFERENCES spare_parts(id),
    order_id INTEGER REFERENCES repair_orders(id),
    task_id INTEGER REFERENCES maintenance_tasks(id),
    quantity INTEGER NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE knowledge_articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    equipment_type_id INTEGER REFERENCES equipment_types(id),
    fault_phenomenon TEXT,
    cause_analysis TEXT,
    solution TEXT,
    source_type VARCHAR(50),
    source_id INTEGER,
    tags TEXT[],
    created_by INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Agent Intelligent Engine (Phase 2 & 3)
-- =====================================================

CREATE TABLE agent_skills (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    applicable_to JSONB,
    applicable_scenarios JSONB,
    steps JSONB NOT NULL,
    version INTEGER DEFAULT 1,
    status VARCHAR(20) DEFAULT 'draft',
    usage_count INTEGER DEFAULT 0,
    success_rate DECIMAL(5,4) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE agent_knowledges (
    id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    type VARCHAR(50),
    summary TEXT,
    details TEXT,
    related_skill_id VARCHAR(100),
    confidence DECIMAL(5,4) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'draft',
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    verified_by INTEGER REFERENCES users(id),
    verified_at TIMESTAMP
);

CREATE TABLE agent_experiences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    category VARCHAR(50),
    content TEXT,
    weight DECIMAL(5,4) DEFAULT 1.0,
    decay_rate DECIMAL(5,4) DEFAULT 0.01,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE agent_conversations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    title VARCHAR(200),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE agent_messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER NOT NULL REFERENCES agent_conversations(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    image_url VARCHAR(500),
    tool_calls JSONB,
    skill_id VARCHAR(100),
    knowledge_ids JSONB DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE repair_cost_details (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES repair_orders(id) ON DELETE CASCADE,
    spare_part_cost DECIMAL(12,2) DEFAULT 0,
    labor_cost DECIMAL(12,2) DEFAULT 0,
    other_cost DECIMAL(12,2) DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'CNY'
);

-- =====================================================
-- Indexes
-- =====================================================
CREATE INDEX idx_agent_messages_conv ON agent_messages(conversation_id);
CREATE INDEX idx_agent_knowledges_status ON agent_knowledges(status);
CREATE INDEX idx_equipment_finance ON equipment(purchase_price, purchase_date);
