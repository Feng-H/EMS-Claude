-- Add lark bot configuration columns to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS lark_app_id VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS lark_app_secret VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS lark_verification_token VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS lark_encrypt_key VARCHAR(100);

-- Add indexes for app_id if we want to look up by it
CREATE INDEX IF NOT EXISTS idx_users_lark_app_id ON users(lark_app_id);
