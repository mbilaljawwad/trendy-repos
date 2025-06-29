-- +migrate Up
-- Database initialization script for trendy-repos application

-- Set database encoding and connection parameters
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create repositories table
CREATE TABLE IF NOT EXISTS repositories (
    id SERIAL PRIMARY KEY,
    github_id BIGINT UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    url VARCHAR(500) NOT NULL,
    stars_count INTEGER DEFAULT 0,
    language VARCHAR(100),
    topics TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_repositories_github_id ON repositories(github_id);
CREATE INDEX IF NOT EXISTS idx_repositories_stars_count ON repositories(stars_count DESC);
CREATE INDEX IF NOT EXISTS idx_repositories_language ON repositories(language);
CREATE INDEX IF NOT EXISTS idx_repositories_created_at ON repositories(created_at DESC);

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update the updated_at column
CREATE TRIGGER update_repositories_updated_at
    BEFORE UPDATE ON repositories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +migrate Down
-- Drop trigger
DROP TRIGGER IF EXISTS update_repositories_updated_at ON repositories;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_repositories_created_at;
DROP INDEX IF EXISTS idx_repositories_language;
DROP INDEX IF EXISTS idx_repositories_stars_count;
DROP INDEX IF EXISTS idx_repositories_github_id;

-- Drop table
DROP TABLE IF EXISTS repositories;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp"; 