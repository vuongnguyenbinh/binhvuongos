-- Consolidate roles: owner | core_staff | ctv | staff → owner | manager | staff
UPDATE users SET role = 'manager' WHERE role = 'core_staff';
UPDATE users SET role = 'staff'   WHERE role = 'ctv';
-- Any unexpected legacy values default to staff
UPDATE users SET role = 'staff'
  WHERE role NOT IN ('owner', 'manager', 'staff');

ALTER TABLE users
  ADD CONSTRAINT chk_user_role
  CHECK (role IN ('owner', 'manager', 'staff'));

-- Update default for future rows
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'staff';
