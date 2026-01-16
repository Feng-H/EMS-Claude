-- =====================================================
-- EMS Database Index Optimization
-- =====================================================
-- This script adds indexes to optimize query performance
-- for high-traffic tables and common search patterns
-- =====================================================

-- =====================================================
-- Equipment Tables Indexes
-- =====================================================

-- Equipment table indexes
-- QR code lookup (high frequency)
CREATE INDEX IF NOT EXISTS idx_equipment_qr_code ON equipment(qr_code);

-- Workshop filtering for equipment lists
CREATE INDEX IF NOT EXISTS idx_equipment_workshop_id ON equipment(workshop_id);

-- Equipment type filtering
CREATE INDEX IF NOT EXISTS idx_equipment_type_id ON equipment(type_id);

-- Status filtering (running/stopped/maintenance/scrapped)
CREATE INDEX IF NOT EXISTS idx_equipment_status ON equipment(status);

-- Dedicated maintenance engineer lookup
CREATE INDEX IF NOT EXISTS idx_equipment_dedicated_maintenance ON equipment(dedicated_maintenance_id);

-- Composite index for workshop + status (common filter combination)
CREATE INDEX IF NOT EXISTS idx_equipment_workshop_status ON equipment(workshop_id, status);

-- =====================================================
-- Organization Tables Indexes
-- =====================================================

-- Factories by base
CREATE INDEX IF NOT EXISTS idx_factories_base_id ON factories(base_id);

-- Workshops by factory
CREATE INDEX IF NOT EXISTS idx_workshops_factory_id ON workshops(factory_id);

-- =====================================================
-- User Tables Indexes
-- =====================================================

-- Username lookup for authentication
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Factory filtering for users
CREATE INDEX IF NOT EXISTS idx_users_factory_id ON users(factory_id);

-- Role filtering for permissions
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- =====================================================
-- Inspection Tables Indexes
-- =====================================================

-- Inspection tasks by assigned user and status (for "My Tasks" queries)
CREATE INDEX IF NOT EXISTS idx_inspection_tasks_assigned_status ON inspection_tasks(assigned_to, status);

-- Inspection tasks by scheduled date
CREATE INDEX IF NOT EXISTS idx_inspection_tasks_scheduled_date ON inspection_tasks(scheduled_date);

-- Inspection tasks by equipment
CREATE INDEX IF NOT EXISTS idx_inspection_tasks_equipment_id ON inspection_tasks(equipment_id);

-- Composite index for assigned user + date + status (mobile dashboard)
CREATE INDEX IF NOT EXISTS idx_inspection_tasks_user_date_status ON inspection_tasks(assigned_to, scheduled_date, status);

-- Inspection records by task
CREATE INDEX IF NOT EXISTS idx_inspection_records_task_id ON inspection_records(task_id);

-- Inspection template by equipment type
CREATE INDEX IF NOT EXISTS idx_inspection_templates_equipment_type ON inspection_templates(equipment_type_id);

-- =====================================================
-- Repair Order Tables Indexes
-- =====================================================

-- Repair orders by equipment (for equipment history)
CREATE INDEX IF NOT EXISTS idx_repair_orders_equipment_id ON repair_orders(equipment_id);

-- Repair orders by assigned engineer and status (for "My Tasks")
CREATE INDEX IF NOT EXISTS idx_repair_orders_assigned_status ON repair_orders(assigned_to, status);

-- Repair orders by status
CREATE INDEX IF NOT EXISTS idx_repair_orders_status ON repair_orders(status);

-- Repair orders by priority
CREATE INDEX IF NOT EXISTS idx_repair_orders_priority ON repair_orders(priority);

-- Repair orders by reporter
CREATE INDEX IF NOT EXISTS idx_repair_orders_reporter_id ON repair_orders(reporter_id);

-- Composite index for equipment + status + created_at (equipment repair history)
CREATE INDEX IF NOT EXISTS idx_repair_orders_equipment_status_created ON repair_orders(equipment_id, status, created_at DESC);

-- Composite index for assigned + status + priority (work queue sorting)
CREATE INDEX IF NOT EXISTS idx_repair_orders_assigned_status_priority ON repair_orders(assigned_to, status, priority, created_at);

-- Repair logs by order
CREATE INDEX IF NOT EXISTS idx_repair_logs_order_id ON repair_logs(order_id);

-- Repair logs by user
CREATE INDEX IF NOT EXISTS idx_repair_logs_user_id ON repair_logs(user_id);

-- =====================================================
-- Maintenance Tables Indexes
-- =====================================================

-- Maintenance plans by equipment type
CREATE INDEX IF NOT EXISTS idx_maintenance_plans_equipment_type ON maintenance_plans(equipment_type_id);

-- Maintenance tasks by equipment
CREATE INDEX IF NOT EXISTS idx_maintenance_tasks_equipment_id ON maintenance_tasks(equipment_id);

-- Maintenance tasks by assigned user and status
CREATE INDEX IF NOT EXISTS idx_maintenance_tasks_assigned_status ON maintenance_tasks(assigned_to, status);

-- Maintenance tasks by scheduled date range
CREATE INDEX IF NOT EXISTS idx_maintenance_tasks_scheduled_date ON maintenance_tasks(scheduled_date);

-- Maintenance tasks by due date
CREATE INDEX IF NOT EXISTS idx_maintenance_tasks_due_date ON maintenance_tasks(due_date);

-- Composite index for assigned + status + due_date (work planning)
CREATE INDEX IF NOT EXISTS idx_maintenance_tasks_user_status_due ON maintenance_tasks(assigned_to, status, due_date);

-- =====================================================
-- Spare Parts Tables Indexes
-- =====================================================

-- Spare parts by code (search)
CREATE INDEX IF NOT EXISTS idx_spare_parts_code ON spare_parts(code);

-- Spare parts by factory
CREATE INDEX IF NOT EXISTS idx_spare_parts_factory_id ON spare_parts(factory_id);

-- Spare parts by name (text search)
CREATE INDEX IF NOT EXISTS idx_spare_parts_name_trgm ON spare_parts USING gin(name gin_trgm_ops);

-- Inventory by spare part
CREATE INDEX IF NOT EXISTS idx_spare_part_inventory_part_id ON spare_part_inventory(spare_part_id);

-- Inventory by factory
CREATE INDEX IF NOT EXISTS idx_spare_part_inventory_factory_id ON spare_part_inventory(factory_id);

-- Low stock alert: quantity + safety_stock
CREATE INDEX IF NOT EXISTS idx_spare_part_inventory_stock_check ON spare_part_inventory(spare_part_id, factory_id, quantity);

-- Consumption by spare part
CREATE INDEX IF NOT EXISTS idx_spare_part_consumptions_part_id ON spare_part_consumptions(spare_part_id);

-- Consumption by repair order
CREATE INDEX IF NOT EXISTS idx_spare_part_consumptions_order_id ON spare_part_consumptions(order_id);

-- Consumption by maintenance task
CREATE INDEX IF NOT EXISTS idx_spare_part_consumptions_task_id ON spare_part_consumptions(task_id);

-- =====================================================
-- Knowledge Base Tables Indexes
-- =====================================================

-- Knowledge articles by equipment type
CREATE INDEX IF NOT EXISTS idx_knowledge_articles_equipment_type ON knowledge_articles(equipment_type_id);

-- Knowledge articles by source type
CREATE INDEX IF NOT EXISTS idx_knowledge_articles_source_type ON knowledge_articles(source_type);

-- Knowledge articles by creator
CREATE INDEX IF NOT EXISTS idx_knowledge_articles_created_by ON knowledge_articles(created_by);

-- Full-text search on title, fault_phenomenon, solution
CREATE INDEX IF NOT EXISTS idx_knowledge_articles_search ON knowledge_articles USING gin(
    to_tsvector('simple', coalesce(title, '') || ' ' || coalesce(fault_phenomenon, '') || ' ' || coalesce(solution, ''))
);

-- Tag search using GIN
CREATE INDEX IF NOT EXISTS idx_knowledge_articles_tags ON knowledge_articles USING gin(tags);

-- =====================================================
-- Statistics and Analytics Indexes
-- =====================================================

-- Created_at indexes for trend analysis
CREATE INDEX IF NOT EXISTS idx_repair_orders_created_at ON repair_orders(created_at);
CREATE INDEX IF NOT EXISTS idx_inspection_tasks_created_at ON inspection_tasks(created_at);
CREATE INDEX IF NOT EXISTS idx_maintenance_tasks_created_at ON maintenance_tasks(created_at);

-- =====================================================
-- Partial Indexes for Common Active Queries
-- =====================================================

-- Only index pending repair orders for queue views
CREATE INDEX IF NOT EXISTS idx_repair_orders_pending
    ON repair_orders(assigned_to, priority, created_at)
    WHERE status IN ('pending', 'assigned', 'in_progress');

-- Only index pending inspection tasks
CREATE INDEX IF NOT EXISTS idx_inspection_tasks_pending
    ON inspection_tasks(assigned_to, scheduled_date)
    WHERE status IN ('pending', 'overdue');

-- Only index pending maintenance tasks
CREATE INDEX IF NOT EXISTS idx_maintenance_tasks_pending
    ON maintenance_tasks(assigned_to, due_date)
    WHERE status IN ('pending', 'overdue');

-- Only index active equipment
CREATE INDEX IF NOT EXISTS idx_equipment_active
    ON equipment(workshop_id, type_id, status)
    WHERE status != 'scrapped';

-- =====================================================
-- Covering Indexes for Specific Queries
-- =====================================================

-- Dashboard equipment count by status (no table lookup needed)
CREATE INDEX IF NOT EXISTS idx_equipment_status_covering
    ON equipment(status) INCLUDE (id);

-- Recent repair orders for dashboard (no table lookup needed)
CREATE INDEX IF NOT EXISTS idx_repair_orders_recent_covering
    ON repair_orders(created_at DESC) INCLUDE (id, equipment_id, status, fault_description);

-- =====================================================
-- Notes
-- =====================================================
-- 1. This script uses CREATE INDEX IF NOT EXISTS to be idempotent
-- 2. Some indexes use gin_trgm_ops which requires pg_trgm extension
-- 3. Run this after enabling the pg_trgm extension:
--    CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- 4. Monitor index usage with:
--    SELECT schemaname, tablename, indexname, idx_scan
--    FROM pg_stat_user_indexes
--    ORDER BY idx_scan ASC;
-- 5. Reindex periodically with:
--    REINDEX DATABASE CONCURRENTLY ems_db;
-- =====================================================
