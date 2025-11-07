-- Create demo user
INSERT INTO users (id, email, password_hash, first_name, last_name, is_demo, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001'::uuid,
    'demo@bucketbird.app',
    '$2a$10$demohashdemohashdemohashdemohashdemohashdemohashdemo', -- Dummy hash, demo user doesn't use password
    'Demo',
    'User',
    true,
    NOW(),
    NOW()
) ON CONFLICT (email) DO NOTHING;

-- Create demo S3 credentials (encrypted with sample data)
INSERT INTO credentials (id, user_id, name, provider, region, endpoint, encrypted_access_key, encrypted_secret_key, use_ssl, status, logo, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000002'::uuid,
    '00000000-0000-0000-0000-000000000001'::uuid,
    'AWS Production',
    'Amazon S3',
    'us-east-1',
    's3.amazonaws.com',
    'ENCRYPTED_DEMO_ACCESS_KEY_1',
    'ENCRYPTED_DEMO_SECRET_KEY_1',
    true,
    'active',
    'cloud',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

INSERT INTO credentials (id, user_id, name, provider, region, endpoint, encrypted_access_key, encrypted_secret_key, use_ssl, status, logo, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000003'::uuid,
    '00000000-0000-0000-0000-000000000001'::uuid,
    'MinIO Development',
    'MinIO',
    'local',
    'minio.dev.local:9000',
    'ENCRYPTED_DEMO_ACCESS_KEY_2',
    'ENCRYPTED_DEMO_SECRET_KEY_2',
    false,
    'active',
    'database',
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

-- Create demo buckets
INSERT INTO buckets (id, user_id, credential_id, name, region, description, size_bytes, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000004'::uuid,
    '00000000-0000-0000-0000-000000000001'::uuid,
    '00000000-0000-0000-0000-000000000002'::uuid,
    'product-images',
    'us-east-1',
    'Production product catalog images and thumbnails',
    524288000,
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

INSERT INTO buckets (id, user_id, credential_id, name, region, description, size_bytes, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000005'::uuid,
    '00000000-0000-0000-0000-000000000001'::uuid,
    '00000000-0000-0000-0000-000000000002'::uuid,
    'user-uploads',
    'us-east-1',
    'User-generated content and file uploads',
    1073741824,
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;

INSERT INTO buckets (id, user_id, credential_id, name, region, description, size_bytes, created_at, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000006'::uuid,
    '00000000-0000-0000-0000-000000000001'::uuid,
    '00000000-0000-0000-0000-000000000003'::uuid,
    'dev-testing',
    'local',
    'Local development and testing bucket',
    104857600,
    NOW(),
    NOW()
) ON CONFLICT (id) DO NOTHING;
