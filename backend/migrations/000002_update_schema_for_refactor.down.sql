-- Rollback: Add back old columns
ALTER TABLE credentials ADD COLUMN IF NOT EXISTS access_key TEXT NOT NULL DEFAULT '';
ALTER TABLE credentials ADD COLUMN IF NOT EXISTS secret_key TEXT NOT NULL DEFAULT '';

-- Rollback: Remove user_id from profiles
DROP INDEX IF EXISTS profiles_user_id_idx;
ALTER TABLE profiles DROP CONSTRAINT IF EXISTS profiles_user_id_unique;
ALTER TABLE profiles DROP COLUMN IF EXISTS user_id;
