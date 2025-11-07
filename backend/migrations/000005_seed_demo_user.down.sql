-- Remove demo buckets
DELETE FROM buckets WHERE user_id = '00000000-0000-0000-0000-000000000001'::uuid;

-- Remove demo credentials
DELETE FROM credentials WHERE user_id = '00000000-0000-0000-0000-000000000001'::uuid;

-- Remove demo user
DELETE FROM users WHERE email = 'demo@bucketbird.app';
