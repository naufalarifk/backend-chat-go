DROP TABLE IF EXISTS message;
CREATE TABLE message (
    ID INT AUTO_INCREMENT NOT NULL,
    Sender VARCHAR(128) NOT NULL,
    Content VARCHAR(255) NOT NULL,
    RoomId INT NOT NULL,
    Timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (ID)
);


INSERT INTO message 
(Sender, Content, RoomId)
VALUES ('JohnDoe', 'Hello there!', 1);