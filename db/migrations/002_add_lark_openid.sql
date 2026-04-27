-- Add lark_openid to users table for Lark integration
ALTER TABLE users ADD COLUMN IF NOT EXISTS lark_openid VARCHAR(100) UNIQUE;

-- Add index for faster lookup
CREATE INDEX IF NOT EXISTS idx_users_lark_openid ON users(lark_openid);
