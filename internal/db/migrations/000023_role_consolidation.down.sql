ALTER TABLE users ALTER COLUMN role SET DEFAULT 'staff';
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_user_role;
-- Forward-only: no automatic role-name rollback
