-- 用户审核和密码管理功能迁移脚本
-- 执行方式: psql -h localhost -U ems -d ems_db -f db/migrations/001_add_user_approval.sql

-- 添加新的审核状态枚举类型
DO $$ BEGIN
    CREATE TYPE user_approval_status AS ENUM ('pending', 'approved', 'rejected');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- 为 users 表添加新字段
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS approval_status user_approval_status DEFAULT 'approved',
    ADD COLUMN IF NOT EXISTS must_change_password BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS first_login BOOLEAN DEFAULT true,
    ADD COLUMN IF NOT EXISTS rejection_reason TEXT;

-- 为现有用户设置默认值
UPDATE users
SET approval_status = 'approved',
    must_change_password = false,
    first_login = false
WHERE approval_status IS NULL;

-- 为新字段添加注释
COMMENT ON COLUMN users.approval_status IS '用户审核状态: pending-待审核, approved-已通过, rejected-已拒绝';
COMMENT ON COLUMN users.must_change_password IS '是否必须修改密码';
COMMENT ON COLUMN users.first_login IS '是否首次登录';
COMMENT ON COLUMN users.rejection_reason IS '账号申请被拒绝的原因';

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_users_approval_status ON users(approval_status);
CREATE INDEX IF NOT EXISTS idx_users_must_change_password ON users(must_change_password) WHERE must_change_password = true;

-- 迁移完成提示
DO $$
BEGIN
    RAISE NOTICE '========================================';
    RAISE NOTICE '数据库迁移完成！';
    RAISE NOTICE '========================================';
    RAISE NOTICE '已添加功能:';
    RAISE NOTICE '  - 用户审核状态';
    RAISE NOTICE '  - 首次登录强制修改密码';
    RAISE NOTICE '  - 账号申请功能';
    RAISE NOTICE '';
    RAISE NOTICE '现有用户已设置为已通过审核状态';
    RAISE NOTICE '========================================';
END $$;
