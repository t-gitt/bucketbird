-- Update credentials table: remove old plaintext columns
ALTER TABLE credentials DROP COLUMN IF EXISTS access_key;
ALTER TABLE credentials DROP COLUMN IF EXISTS secret_key;

-- Update profiles table: add user_id column and foreign key
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'profiles' AND column_name = 'user_id'
    ) THEN
        ALTER TABLE profiles ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;

        -- Populate user_id for existing profiles (assuming id = user_id in old schema)
        UPDATE profiles SET user_id = id WHERE user_id IS NULL;

        -- Make user_id NOT NULL after populating
        ALTER TABLE profiles ALTER COLUMN user_id SET NOT NULL;
    END IF;
END $$;

-- Add unique constraint on user_id
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'profiles_user_id_unique'
    ) THEN
        ALTER TABLE profiles ADD CONSTRAINT profiles_user_id_unique UNIQUE (user_id);
    END IF;
END $$;

-- Add index on user_id
CREATE INDEX IF NOT EXISTS profiles_user_id_idx ON profiles(user_id);
