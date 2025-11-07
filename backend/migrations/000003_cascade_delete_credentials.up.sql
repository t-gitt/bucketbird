-- Change credential foreign key constraint to CASCADE DELETE
ALTER TABLE buckets DROP CONSTRAINT IF EXISTS buckets_credential_id_fkey;
ALTER TABLE buckets ADD CONSTRAINT buckets_credential_id_fkey
    FOREIGN KEY (credential_id) REFERENCES credentials(id) ON DELETE CASCADE;
