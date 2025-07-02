CREATE TABLE photos
(
    id        INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    player_id VARCHAR(50) NOT NULL,
    image     TEXT               NOT NULL,
    type      VARCHAR(60),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);