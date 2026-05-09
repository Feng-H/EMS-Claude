-- db/seeds/agent_skills_v2.sql

INSERT INTO agent_skills (name, description, applicable_to, applicable_scenarios, steps, status) VALUES
(
    '设备生命周期经济性分析', 
    '从财务视角分析设备的全生命周期成本（TCO），判断老旧设备是该“继续维修”还是“报废替换”。',
    '{"equipment_types": ["all"]}',
    '["设备报废评估", "大修决策", "年度预算规划"]',
    '[
        {"step": 1, "action": "获取设备财务数据", "tool": "get_equipment_financials"},
        {"step": 2, "action": "获取累计维修成本", "tool": "get_repair_costs"},
        {"step": 3, "action": "经济性对比分析", "logic": "对比(累计维修成本 + 预计停机损失) 与 (设备原值 - 残值)"},
        {"step": 4, "action": "提供决策建议", "logic": "如果维修成本占比超过 40%，发出预警并建议评估投资回报率（ROI）"}
    ]',
    'active'
),
(
    '故障根因与可靠性分析',
    '分析设备的历史报修记录，计算 MTTR/MTBF，识别高频故障模式及其根本原因。',
    '{"equipment_types": ["all"]}',
    '["重复性故障诊断", "可靠性提升", "备件库存优化"]',
    '[
        {"step": 1, "action": "获取历史报修详情", "tool": "query_repair_orders"},
        {"step": 2, "action": "计算可靠性指标", "logic": "计算平均修复时间 (MTTR) 和平均故障间隔时间 (MTBF)"},
        {"step": 3, "action": "故障模式识别", "logic": "统计高频故障代码 (fault_code) 或关键词"},
        {"step": 4, "action": "知识库方案匹配", "tool": "search_manual_knowledge"}
    ]',
    'active'
),
(
    '维保合规与健康度综合评估',
    '综合考量巡检、保养和实时亚健康征兆，对设备的运行健康度进行评分。',
    '{"equipment_types": ["all"]}',
    '["日常健康巡检", "预防性维护规划", "生产计划调整参考"]',
    '[
        {"step": 1, "action": "检查保养合规性", "tool": "get_maintenance_compliance"},
        {"step": 2, "action": "获取实时健康数据", "tool": "get_equipment_health"},
        {"step": 3, "action": "亚健康征兆识别", "tool": "detect_symptoms"},
        {"step": 4, "action": "综合评分与建议", "logic": "综合各项数据给出一个 0-100 的健康评分，并提供针对性的维护建议"}
    ]',
    'active'
);
