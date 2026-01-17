-- EMS Database Schema
-- PostgreSQL 15+

-- Drop existing tables (for clean setup)
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

-- Drop custom types
DROP TYPE IF EXISTS user_role CASCADE;
DROP TYPE IF EXISTS repair_status CASCADE;

-- =====================================================
-- Custom Types
-- =====================================================

CREATE TYPE user_role AS ENUM ('admin', 'supervisor', 'engineer', 'maintenance', 'operator');
CREATE TYPE user_approval_status AS ENUM ('pending', 'approved', 'rejected');
CREATE TYPE repair_status AS ENUM ('pending', 'assigned', 'in_progress', 'testing', 'confirmed', 'audited', 'closed');

-- =====================================================
-- Organization Structure
-- =====================================================

CREATE TABLE bases (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE factories (
    id SERIAL PRIMARY KEY,
    base_id INTEGER NOT NULL REFERENCES bases(id) ON DELETE RESTRICT,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(base_id, code)
);

CREATE TABLE workshops (
    id SERIAL PRIMARY KEY,
    factory_id INTEGER NOT NULL REFERENCES factories(id) ON DELETE RESTRICT,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(factory_id, code)
);

-- =====================================================
-- Users
-- =====================================================

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    role user_role NOT NULL DEFAULT 'operator',
    factory_id INTEGER REFERENCES factories(id) ON DELETE SET NULL,
    phone VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    approval_status user_approval_status DEFAULT 'approved',
    must_change_password BOOLEAN DEFAULT false,
    first_login BOOLEAN DEFAULT true,
    rejection_reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Equipment Management
-- =====================================================

CREATE TABLE equipment_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    inspection_template_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE equipment (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    type_id INTEGER NOT NULL REFERENCES equipment_types(id) ON DELETE RESTRICT,
    workshop_id INTEGER NOT NULL REFERENCES workshops(id) ON DELETE RESTRICT,
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    spec TEXT,
    purchase_date DATE,
    status VARCHAR(20) DEFAULT 'running', -- running/stopped/maintenance/scrapped
    dedicated_maintenance_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Inspection Management
-- =====================================================

CREATE TABLE inspection_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    equipment_type_id INTEGER NOT NULL REFERENCES equipment_types(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE inspection_items (
    id SERIAL PRIMARY KEY,
    template_id INTEGER NOT NULL REFERENCES inspection_templates(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    method VARCHAR(500),
    criteria VARCHAR(500),
    sequence_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE inspection_tasks (
    id SERIAL PRIMARY KEY,
    equipment_id INTEGER NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    template_id INTEGER NOT NULL REFERENCES inspection_templates(id) ON DELETE RESTRICT,
    assigned_to INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    scheduled_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending/in_progress/completed/overdue
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    latitude DECIMAL(10, 7),
    longitude DECIMAL(10, 7),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE inspection_records (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES inspection_tasks(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES inspection_items(id) ON DELETE RESTRICT,
    result VARCHAR(10) NOT NULL, -- OK/NG
    remark TEXT,
    photo_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Repair Management
-- =====================================================

CREATE TABLE repair_orders (
    id SERIAL PRIMARY KEY,
    equipment_id INTEGER NOT NULL REFERENCES equipment(id) ON DELETE RESTRICT,
    fault_description TEXT NOT NULL,
    reporter_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    assigned_to INTEGER REFERENCES users(id) ON DELETE SET NULL,
    status repair_status DEFAULT 'pending',
    priority INTEGER DEFAULT 3, -- 1=high, 2=medium, 3=low
    photos TEXT[], -- Array of photo URLs
    fault_code VARCHAR(50),
    solution TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    confirmed_at TIMESTAMP,
    audited_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE repair_logs (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES repair_orders(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    action VARCHAR(50),
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Maintenance Management
-- =====================================================

CREATE TABLE maintenance_plans (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    equipment_type_id INTEGER NOT NULL REFERENCES equipment_types(id) ON DELETE CASCADE,
    level INTEGER NOT NULL, -- 1=primary, 2=secondary, 3=precision
    cycle_days INTEGER NOT NULL,
    flexible_days INTEGER DEFAULT 3,
    work_hours DECIMAL(5,1),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE maintenance_tasks (
    id SERIAL PRIMARY KEY,
    plan_id INTEGER NOT NULL REFERENCES maintenance_plans(id) ON DELETE CASCADE,
    equipment_id INTEGER NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    assigned_to INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    scheduled_date DATE NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending/in_progress/completed
    completed_at TIMESTAMP,
    actual_hours DECIMAL(5,1),
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Spare Parts Management
-- =====================================================

CREATE TABLE spare_parts (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    specification TEXT,
    unit VARCHAR(20),
    factory_id INTEGER REFERENCES factories(id) ON DELETE SET NULL,
    safety_stock INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE spare_part_inventory (
    id SERIAL PRIMARY KEY,
    spare_part_id INTEGER NOT NULL REFERENCES spare_parts(id) ON DELETE CASCADE,
    factory_id INTEGER NOT NULL REFERENCES factories(id) ON DELETE CASCADE,
    quantity INTEGER DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(spare_part_id, factory_id)
);

CREATE TABLE spare_part_consumptions (
    id SERIAL PRIMARY KEY,
    spare_part_id INTEGER NOT NULL REFERENCES spare_parts(id) ON DELETE RESTRICT,
    order_id INTEGER REFERENCES repair_orders(id) ON DELETE SET NULL,
    task_id INTEGER REFERENCES maintenance_tasks(id) ON DELETE SET NULL,
    quantity INTEGER NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Knowledge Base
-- =====================================================

CREATE TABLE knowledge_articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    equipment_type_id INTEGER REFERENCES equipment_types(id) ON DELETE SET NULL,
    fault_phenomenon TEXT,
    cause_analysis TEXT,
    solution TEXT,
    source_type VARCHAR(20), -- repair/manual/other
    source_id INTEGER,
    tags TEXT[],
    created_by INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- Indexes for Performance
-- =====================================================

-- Equipment indexes
CREATE INDEX idx_equipment_qr ON equipment(qr_code);
CREATE INDEX idx_equipment_workshop ON equipment(workshop_id);
CREATE INDEX idx_equipment_type ON equipment(type_id);
CREATE INDEX idx_equipment_status ON equipment(status);
CREATE INDEX idx_equipment_maintenance ON equipment(dedicated_maintenance_id);

-- Inspection indexes
CREATE INDEX idx_inspection_task_assigned ON inspection_tasks(assigned_to, status, scheduled_date);
CREATE INDEX idx_inspection_task_equipment ON inspection_tasks(equipment_id, scheduled_date);
CREATE INDEX idx_inspection_task_status ON inspection_tasks(status);
CREATE INDEX idx_inspection_record_task ON inspection_records(task_id);

-- Repair indexes
CREATE INDEX idx_repair_equipment_status ON repair_orders(equipment_id, status);
CREATE INDEX idx_repair_assigned ON repair_orders(assigned_to, status);
CREATE INDEX idx_repair_reporter ON repair_orders(reporter_id);
CREATE INDEX idx_repair_status ON repair_orders(status);
CREATE INDEX idx_repair_created ON repair_orders(created_at DESC);
CREATE INDEX idx_repair_logs_order ON repair_logs(order_id);

-- Maintenance indexes
CREATE INDEX idx_maintenance_task_assigned ON maintenance_tasks(assigned_to, status, scheduled_date);
CREATE INDEX idx_maintenance_task_equipment ON maintenance_tasks(equipment_id, scheduled_date);
CREATE INDEX idx_maintenance_task_status ON maintenance_tasks(status);

-- Spare parts indexes
CREATE INDEX idx_spare_part_inventory_sp ON spare_part_inventory(spare_part_id);
CREATE INDEX idx_spare_part_inventory_factory ON spare_part_inventory(factory_id);
CREATE INDEX idx_spare_part_consumption_sp ON spare_part_consumptions(spare_part_id);
CREATE INDEX idx_spare_part_consumption_order ON spare_part_consumptions(order_id);

-- Knowledge base indexes (for full-text search)
CREATE INDEX idx_knowledge_title ON knowledge_articles USING gin(to_tsvector('english', title));
CREATE INDEX idx_knowledge_tags ON knowledge_articles USING gin(tags);

-- User indexes
CREATE INDEX idx_users_factory ON users(factory_id);
CREATE INDEX idx_users_role ON users(role);

-- Organization indexes
CREATE INDEX idx_factories_base ON factories(base_id);
CREATE INDEX idx_workshops_factory ON workshops(factory_id);

-- =====================================================
-- Functions for Automatic Timestamp Updates
-- =====================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply to tables with updated_at
CREATE TRIGGER update_bases_updated_at BEFORE UPDATE ON bases
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_factories_updated_at BEFORE UPDATE ON factories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_workshops_updated_at BEFORE UPDATE ON workshops
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_equipment_types_updated_at BEFORE UPDATE ON equipment_types
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_equipment_updated_at BEFORE UPDATE ON equipment
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_inspection_templates_updated_at BEFORE UPDATE ON inspection_templates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_inspection_items_updated_at BEFORE UPDATE ON inspection_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_inspection_tasks_updated_at BEFORE UPDATE ON inspection_tasks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_repair_orders_updated_at BEFORE UPDATE ON repair_orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_maintenance_plans_updated_at BEFORE UPDATE ON maintenance_plans
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_maintenance_tasks_updated_at BEFORE UPDATE ON maintenance_tasks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_spare_parts_updated_at BEFORE UPDATE ON spare_parts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_knowledge_articles_updated_at BEFORE UPDATE ON knowledge_articles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
