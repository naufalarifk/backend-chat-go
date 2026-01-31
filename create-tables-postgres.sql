-- PostgreSQL Schema for Realtime Chat Backend
-- This file will be automatically executed when the container first starts

-- Drop existing table if it exists
DROP TABLE IF EXISTS message CASCADE;

-- Create message table with PostgreSQL-specific features
CREATE TABLE message (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender VARCHAR(128) NOT NULL,
    content TEXT NOT NULL,
    room_id VARCHAR(128) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB DEFAULT '{}',
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_message_room_id ON message(room_id);
CREATE INDEX idx_message_timestamp ON message(timestamp DESC);
CREATE INDEX idx_message_sender ON message(sender);
CREATE INDEX idx_message_metadata ON message USING GIN(metadata);
CREATE INDEX idx_message_room_timestamp ON message(room_id, timestamp DESC);

-- Create function for automatically updating updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to call the function before any update
CREATE TRIGGER update_message_updated_at
    BEFORE UPDATE ON message
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert sample data for testing
INSERT INTO message (sender, content, room_id) VALUES
('JohnDoe', 'Hello there!', 'default'),
('JaneDoe', 'Hi John! How are you?', 'default'),
('Alice', 'Welcome to the chat!', 'general'),
('Bob', 'Thanks Alice!', 'general');

-- Create a view for recent messages (optional, useful for queries)
CREATE OR REPLACE VIEW recent_messages AS
SELECT 
    id,
    sender,
    content,
    room_id,
    timestamp,
    metadata
FROM message
WHERE is_deleted = FALSE
ORDER BY timestamp DESC
LIMIT 100;

-- Display table info
SELECT 
    'Messages table created successfully!' as status,
    COUNT(*) as sample_message_count 
FROM message;