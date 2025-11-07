-- Add is_demo column to users table
ALTER TABLE users ADD COLUMN is_demo BOOLEAN NOT NULL DEFAULT false;
