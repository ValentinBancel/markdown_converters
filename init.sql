-- Initialize the markdown_converters database
-- This script runs when the PostgreSQL container starts for the first time

-- Create the database if it doesn't exist
SELECT 'CREATE DATABASE markdown_converters'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'markdown_converters')\gexec

-- Connect to the database
\c markdown_converters;

-- Create extension for UUID generation (optional, for future use)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- The tables will be created automatically by GORM Auto-Migrate
-- but we can add any additional setup here

-- Insert some sample data for testing (optional)
-- This will be handled by the application, but can be useful for initial testing