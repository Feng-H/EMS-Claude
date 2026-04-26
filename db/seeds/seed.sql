-- EMS Seed Data (Hyper-Realistic Industrial World)
-- Covers 180 days of history for MTTR/MTBF/TCO visualization

-- 1. Organization & Users
INSERT INTO bases (code, name) VALUES ('BASE-HQ', '集团总部基地');
INSERT INTO factories (base_id, code, name) VALUES (1, 'FAC-SZ', '苏州智能工厂');
INSERT INTO workshops (factory_id, code, name) VALUES (1, 'WS-MCH', '精密机加车间'), (1, 'WS-ASM', '全自动装配车间');

-- Passwords: admin123 ($2b$10$WEQ9sLfa60UCu7RTu02IN.JJgae3Ux99TpVTBt4HeMlxHjt1GmS8S)
INSERT INTO users (username, password_hash, name, role, factory_id, phone) VALUES
('admin', '$2a$10$LwwSxuZWIbYwWaN5zCsD2uXAcZZC9cxMn.g2T.hGuatnKm4UP1Zq6', '系统管理员', 'admin', NULL, '13800000001'),
('supervisor', '$2a$10$LwwSxuZWIbYwWaN5zCsD2uXAcZZC9cxMn.g2T.hGuatnKm4UP1Zq6', '王主管', 'supervisor', 1, '13800000002'),
('maint_li', '$2a$10$LwwSxuZWIbYwWaN5zCsD2uXAcZZC9cxMn.g2T.hGuatnKm4UP1Zq6', '预防型-李四', 'maintenance', 1, '13800000003'),
('maint_zhang', '$2a$10$LwwSxuZWIbYwWaN5zCsD2uXAcZZC9cxMn.g2T.hGuatnKm4UP1Zq6', '救火型-张三', 'maintenance', 1, '13800000004'),
('operator_wang', '$2a$10$LwwSxuZWIbYwWaN5zCsD2uXAcZZC9cxMn.g2T.hGuatnKm4UP1Zq6', '操作工小王', 'operator', 1, '13800000005');

-- 2. Equipment Types & Templates
INSERT INTO equipment_types (name, category) VALUES 
('高精度数控机床', '加工设备'), ('全自动冲床', '成型设备'), ('工业焊接机器人', '机器人');

INSERT INTO inspection_templates (name, equipment_type_id) VALUES ('CNC日检标准', 1), ('冲床周检标准', 2);
INSERT INTO inspection_items (template_id, name, criteria, sequence_order) VALUES 
(1, '液压压力', '4.0-5.0 MPa', 1), (1, '润滑油位', '可见液面', 2), (1, '刀库状态', '无异响', 3),
(2, '离合器间隙', '< 0.5mm', 1), (2, '制动器磨损', '磨损量 < 2mm', 2);

-- 3. Equipment (Financials & Lifecycle)
INSERT INTO equipment (code, name, type_id, workshop_id, qr_code, spec, purchase_price, purchase_date, service_life_years, scrap_value, hourly_loss, status, dedicated_maintenance_id) VALUES
('CNC-001', 'A区-精密机床(李四维护)', 1, 1, 'QR_A', 'VMC850', 280000.00, CURRENT_DATE - INTERVAL '3 years', 8, 28000.00, 150.00, 'running', 3),
('CNC-002', 'B区-精密机床(张三维护)', 1, 1, 'QR_B', 'VMC850', 280000.00, CURRENT_DATE - INTERVAL '3 years', 8, 28000.00, 150.00, 'stopped', 4),
('PRESS-05', 'C区-12年老旧冲床', 2, 1, 'QR_OLD', 'P-200T', 150000.00, CURRENT_DATE - INTERVAL '12 years', 10, 5000.00, 80.00, 'maintenance', 4),
('ROBOT-01', 'D区-码垛机器人', 3, 2, 'QR_R1', 'ABB-6700', 450000.00, CURRENT_DATE - INTERVAL '1 year', 5, 45000.00, 300.00, 'running', 3);

-- 4. Spare Parts & Inventory
INSERT INTO spare_parts (code, name, specification, unit, factory_id, safety_stock) VALUES
('PUMP-01', '高压柱塞泵', 'Rexroth A10V', '台', 1, 2),
('FLT-CHEAP', '普通滤芯(降本件)', 'Generic-100', '个', 1, 50),
('FLT-OEM', '原厂精密滤芯', 'Rexroth-E30', '个', 1, 20),
('OIL-HM46', '液压油', 'Shell HM-46', '桶', 1, 100);

INSERT INTO spare_part_inventory (spare_part_id, factory_id, quantity) VALUES 
(1, 1, 3), (2, 1, 120), (3, 1, 15), (4, 1, 250);

-- 5. Maintenance History (Li Si is diligent, Zhang San is lazy)
INSERT INTO maintenance_plans (name, equipment_type_id, level, cycle_days, flexible_days, work_hours) VALUES
('CNC二级保养', 1, 2, 30, 3, 4.0);

-- Li Si's 6 months of diligent maintenance for CNC-001
INSERT INTO maintenance_tasks (plan_id, equipment_id, assigned_to, scheduled_date, status, completed_at, actual_hours) 
SELECT 1, 1, 3, CURRENT_DATE - (n || ' days')::interval, 'completed', CURRENT_TIMESTAMP - (n || ' days')::interval, 4.2
FROM generate_series(30, 180, 30) AS n;

-- Zhang San's 6 months of lazy maintenance for CNC-002
INSERT INTO maintenance_tasks (plan_id, equipment_id, assigned_to, scheduled_date, status, completed_at, actual_hours) 
SELECT 1, 2, 4, CURRENT_DATE - (n || ' days')::interval, 'completed', CURRENT_TIMESTAMP - (n || ' days')::interval, 0.2
FROM generate_series(30, 180, 30) AS n;

-- 6. Repair History & The Cascade Failure (5.5W Lesson)
-- Zhang San's CNC-002 fails due to poor filter
INSERT INTO repair_orders (equipment_id, fault_description, reporter_id, assigned_to, status, priority, fault_code, started_at, completed_at, confirmed_at, audited_at, closed_at, solution) VALUES
(2, '由于使用劣质滤芯导致的液压系统严重污染，泵轴承毁坏并拉伤油缸', 5, 4, 'closed', 1, 'HYD-ERR-001', CURRENT_TIMESTAMP - INTERVAL '10 days 4 hours', CURRENT_TIMESTAMP - INTERVAL '10 days', CURRENT_TIMESTAMP - INTERVAL '9 days', CURRENT_TIMESTAMP - INTERVAL '8 days', CURRENT_TIMESTAMP - INTERVAL '8 days', '更换高压泵，彻底清洗油路，更换所有滤芯');

INSERT INTO repair_cost_details (order_id, spare_part_cost, labor_cost) VALUES (1, 52000.00, 3000.00);
INSERT INTO spare_part_consumptions (spare_part_id, order_id, quantity, user_id) VALUES (1, 1, 1, 4), (2, 1, 2, 4);

-- Old PRESS-05 constant repairs
INSERT INTO repair_orders (equipment_id, fault_description, reporter_id, assigned_to, status, priority, started_at, completed_at, closed_at, solution)
SELECT 3, '离合器打滑/漏油', 5, 4, 'closed', 2, CURRENT_TIMESTAMP - (n || ' days')::interval, CURRENT_TIMESTAMP - (n || ' days')::interval + INTERVAL '4 hours', CURRENT_TIMESTAMP - (n || ' days')::interval + INTERVAL '1 day', '紧急紧固并补油'
FROM generate_series(15, 180, 20) AS n;

-- 7. Agent Pre-loaded Intelligence
INSERT INTO agent_skills (name, description, steps, status) VALUES 
('级联失效审计', '识别由于低价值易损件失效导致的连锁重大故障风险', '[{"step": 1, "action": "查询高价值备件更换记录", "tool": "get_cost_analysis"}, {"step": 2, "action": "核查保养记录有效性", "tool": "get_maintenance_compliance"}]', 'active');

INSERT INTO agent_knowledges (id, title, type, summary, confidence, status, created_by) VALUES
('k_seed_001', '高压泵级联失效预防', 'root_cause_analysis', '劣质滤芯是导致 5.5 万元柱塞泵报废的根本诱因', 0.98, 'confirmed', 'expert_system');
