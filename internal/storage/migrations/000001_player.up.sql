CREATE TABLE players (
                       id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                       player_id VARCHAR(50) UNIQUE NOT NULL,
                       images VARCHAR(60) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);