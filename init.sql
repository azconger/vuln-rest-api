-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    role VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert seed data
INSERT INTO users (username, email, role) VALUES
    ('admin', 'admin@example.com', 'admin'),
    ('user1', 'user1@example.com', 'user'),
    ('user2', 'user2@example.com', 'user'),
    ('moderator', 'mod@example.com', 'moderator')
ON CONFLICT (username) DO NOTHING; 