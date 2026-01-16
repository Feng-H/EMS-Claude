-- EMS Seed Data
-- This file contains sample data for development and testing

-- =====================================================
-- Organization
-- =====================================================

-- Bases
INSERT INTO bases (code, name) VALUES
('BASE01', '华东基地'),
('BASE02', '华南基地'),
('BASE03', '华北基地');

-- Factories
INSERT INTO factories (base_id, code, name) VALUES
(1, 'FAC01', '苏州工厂'),
(1, 'FAC02', '杭州工厂'),
(2, 'FAC03', '深圳工厂'),
(3, 'FAC04', '天津工厂');

-- Workshops
INSERT INTO workshops (factory_id, code, name) VALUES
(1, 'WS01', '机加车间'),
(1, 'WS02', '装配车间'),
(1, 'WS03', '焊接车间'),
(2, 'WS04', '涂装车间'),
(2, 'WS05', '注塑车间');

-- =====================================================
-- Users (password is "password123" for all users)
-- Hash generated with bcrypt cost 10
-- =====================================================

INSERT INTO users (username, password_hash, name, role, factory_id, phone) VALUES
-- Administrators
('admin', '$2a$10$YourHashedPasswordHere', '系统管理员', 'admin', NULL, '13800000001'),

-- Supervisors
('supervisor01', '$2a$10$YourHashedPasswordHere', '设备主管-苏州', 'supervisor', 1, '13800000011'),
('supervisor02', '$2a$10$YourHashedPasswordHere', '设备主管-杭州', 'supervisor', 2, '13800000012'),

-- Engineers
('engineer01', '$2a$10$YourHashedPasswordHere', '设备工程师-苏州', 'engineer', 1, '13800000021'),
('engineer02', '$2a$10$YourHashedPasswordHere', '设备工程师-杭州', 'engineer', 2, '13800000022'),

-- Maintenance Workers
('maint01', '$2a$10$YourHashedPasswordHere', '维修工-张三', 'maintenance', 1, '13800000101'),
('maint02', '$2a$10$YourHashedPasswordHere', '维修工-李四', 'maintenance', 1, '13800000102'),
('maint03', '$2a$10$YourHashedPasswordHere', '维修工-王五', 'maintenance', 2, '13800000103'),

-- Operators
('oper01', '$2a$10$YourHashedPasswordHere', '操作工-赵六', 'operator', 1, '13800000201'),
('oper02', '$2a$10$YourHashedPasswordHere', '操作工-孙七', 'operator', 1, '13800000202'),
('oper03', '$2a$10$YourHashedPasswordHere', '操作工-周八', 'operator', 2, '13800000203');

-- =====================================================
-- Equipment Types
-- =====================================================

INSERT INTO equipment_types (name, category) VALUES
('CNC加工中心', '机加设备'),
('数控车床', '机加设备'),
('焊接机器人', '焊接设备'),
('注塑机', '注塑设备'),
('涂装线', '涂装设备'),
('装配线', '装配设备'),
('空压机', '公用设备'),
('行车', '起重设备');

-- =====================================================
-- Inspection Templates
-- =====================================================

-- CNC Inspection Template
INSERT INTO inspection_templates (name, equipment_type_id) VALUES
('CNC加工中心日常点检', 1);

INSERT INTO inspection_items (template_id, name, method, criteria, sequence_order) VALUES
(1, '检查液压油位', '观察油标', '在上下限之间', 1),
(1, '检查导轨润滑油', '观察油标', '油位正常', 2),
(1, '检查冷却液浓度', '用浓度计测量', '8-12%', 3),
(1, '检查主轴运转', '听觉检查', '无异响', 4),
(1, '检查安全门锁', '操作检查', '正常锁定', 5),
(1, '清理机床铁屑', '目视检查', '铁屑清理干净', 6);

-- Welding Robot Inspection Template
INSERT INTO inspection_templates (name, equipment_type_id) VALUES
('焊接机器人日常点检', 3);

INSERT INTO inspection_items (template_id, name, method, criteria, sequence_order) VALUES
(2, '检查焊枪电缆', '目视检查', '无破损', 1),
(2, '检查导电嘴', '目视检查', '磨损正常', 2),
(3, '检查保护气压力', '读压力表', '0.4-0.5MPa', 3),
(2, '检查机器人零点', '示教器检查', '零点正常', 4),
(2, '清理焊渣飞溅', '目视/清理', '清洁无飞溅', 5);

-- =====================================================
-- Equipment
-- =====================================================

INSERT INTO equipment (code, name, type_id, workshop_id, qr_code, spec, status) VALUES
-- CNC equipment
('EQ001', 'CNC加工中心-01', 1, 1, 'QR_EQ001', 'VMC850,8000rpm', 'running'),
('EQ002', 'CNC加工中心-02', 1, 1, 'QR_EQ002', 'VMC850,8000rpm', 'running'),
('EQ003', 'CNC加工中心-03', 1, 1, 'QR_EQ003', 'VMC1060,6000rpm', 'maintenance'),

-- Lathe equipment
('EQ004', '数控车床-01', 2, 1, 'QR_EQ004', 'CK6150,3000rpm', 'running'),
('EQ005', '数控车床-02', 2, 1, 'QR_EQ005', 'CK6150,3000rpm', 'stopped'),

-- Welding robots
('EQ006', '焊接机器人-01', 3, 3, 'QR_EQ006', 'OTC AI-20', 'running'),
('EQ007', '焊接机器人-02', 3, 3, 'QR_EQ007', 'OTC AI-20', 'running'),

-- Injection molding machines
('EQ008', '注塑机-01', 4, 5, 'QR_EQ008', '海天HTF300', 'running'),
('EQ009', '注塑机-02', 4, 5, 'QR_EQ009', '海天HTF200', 'running'),

-- Assembly line
('EQ010', '装配线-01', 6, 2, 'QR_EQ010', '20米流水线', 'running');

-- Update equipment type with template references
UPDATE equipment_types SET inspection_template_id = 1 WHERE id = 1;
UPDATE equipment_types SET inspection_template_id = 2 WHERE id = 3;

-- Assign dedicated maintenance workers
UPDATE equipment SET dedicated_maintenance_id = 6 WHERE id IN (1, 2, 4); -- maint01
UPDATE equipment SET dedicated_maintenance_id = 7 WHERE id IN (3, 5); -- maint02
UPDATE equipment SET dedicated_maintenance_id = 8 WHERE id IN (6, 7, 8, 9); -- maint03

-- =====================================================
-- Maintenance Plans
-- =====================================================

INSERT INTO maintenance_plans (name, equipment_type_id, level, cycle_days, flexible_days, work_hours) VALUES
('CNC一级保养', 1, 1, 7, 2, 1.0),
('CNC二级保养', 1, 2, 30, 5, 4.0),
('CNC精度维护', 1, 3, 90, 10, 8.0),
('焊接机器人一级保养', 3, 1, 7, 2, 0.5),
('焊接机器人二级保养', 3, 2, 30, 5, 2.0),
('注塑机一级保养', 4, 1, 7, 2, 1.0),
('注塑机二级保养', 4, 2, 30, 5, 3.0);

-- =====================================================
-- Sample Inspection Tasks
-- =====================================================

INSERT INTO inspection_tasks (equipment_id, template_id, assigned_to, scheduled_date, status) VALUES
(1, 1, 9, CURRENT_DATE, 'pending'),
(2, 1, 9, CURRENT_DATE, 'pending'),
(4, 1, 10, CURRENT_DATE, 'pending'),
(6, 2, 9, CURRENT_DATE, 'completed');

-- =====================================================
-- Sample Repair Orders
-- =====================================================

INSERT INTO repair_orders (equipment_id, fault_description, reporter_id, assigned_to, status, priority, fault_code) VALUES
(3, '主轴异响，需检查轴承', 9, 6, 'in_progress', 1, 'MCH-001'),
(5, '伺服报警', 10, 7, 'pending', 2, 'SRV-002'),
(8, '合模异常', 10, NULL, 'pending', 2, 'MOLD-001');

-- Repair logs
INSERT INTO repair_logs (order_id, user_id, action, content) VALUES
(1, 9, 'created', '报修单已创建'),
(1, 6, 'accepted', '维修工已接单'),
(1, 6, 'progress', '正在检查主轴'),
(2, 10, 'created', '报修单已创建'),
(3, 10, 'created', '报修单已创建');

-- =====================================================
-- Spare Parts
-- =====================================================

INSERT INTO spare_parts (code, name, specification, unit, factory_id, safety_stock) VALUES
('SP001', '主轴轴承', '7014C', '个', 1, 5),
('SP002', '导轨滑块', 'HGH30', '个', 1, 20),
('SP003', '焊枪导电嘴', 'φ1.2', '个', 1, 50),
('SP004', '密封圈', 'O型圈-50', '个', 1, 100),
('SP005', '液压油', 'HM-46', '升', 1, 200),
('SP006', '冷却液', '水溶性', '升', 1, 500);

-- Spare part inventory
INSERT INTO spare_part_inventory (spare_part_id, factory_id, quantity) VALUES
(1, 1, 8),
(2, 1, 35),
(3, 1, 60),
(4, 1, 150),
(5, 1, 300),
(6, 1, 600);

-- =====================================================
-- Knowledge Articles
-- =====================================================

INSERT INTO knowledge_articles (title, equipment_type_id, fault_phenomenon, cause_analysis, solution, source_type, created_by, tags) VALUES
('CNC主轴异响处理方法', 1, '主轴运转时有异常响声', '1. 轴承润滑不足\n2. 轴承磨损\n3. 主轴同心度偏差', '1. 检查润滑油路\n2. 如润滑正常则需更换轴承\n3. 检查主轴安装精度', 'manual', 2, ARRAY['主轴', '异响', '轴承']),
('焊接机器人焊偏问题', 3, '焊缝偏离预定位置', '1. 焊枪校准偏差\n2. 工装定位不准\n3. 编程坐标错误', '1. 重新校准焊枪TCP\n2. 检查工装夹具\n3. 检查并修正程序坐标', 'manual', 2, ARRAY['焊接', '焊偏', '校准']),
('注塑机合模不紧', 4, '合模后仍有缝隙', '1. 合模力不足\n2. 模板平行度差\n3. 锁模机构磨损', '1. 调整合模力\n2. 检查模板平行度\n3. 检查锁模油缸和铰链', 'manual', 2, ARRAY['注塑', '合模', '锁模']);
