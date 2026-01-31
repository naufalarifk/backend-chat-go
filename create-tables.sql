DROP TABLE IF EXISTS message;

CREATE TABLE message (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    sender VARCHAR(128) NOT NULL,
    content TEXT NOT NULL,
    room_id VARCHAR(128) NOT NULL,
    timestamp TIMESTAMP(6) NOT NULL,
    INDEX idx_room_id (room_id),
    INDEX idx_timestamp (timestamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Sample data
INSERT INTO message 
(id, sender, content, room_id, timestamp)
VALUES 
('550e8400-e29b-41d4-a716-446655440000', 'JohnDoe', 'Hello there!', 'default', NOW()),
('550e8400-e29b-41d4-a716-446655440001', 'JaneDoe', 'Hi John!', 'default', NOW());

DESCRIBE message;