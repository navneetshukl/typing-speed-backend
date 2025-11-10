CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    avg_speed INTEGER DEFAULT 0 CHECK (avg_speed >= 0),
    avg_accuracy INTEGER DEFAULT 0 CHECK (avg_accuracy BETWEEN 0 AND 100),
    total_test INTEGER DEFAULT 0 CHECK (total_test >= 0),
    level INTEGER DEFAULT 0 CHECK (level >= 0),
    last_test_time TIMESTAMPTZ,
    streak INTEGER DEFAULT 0 CHECK (streak >= 0)
);

INSERT INTO users (name, email, password, created_at)
VALUES
('Navneet Shukla', 'navneet@example.com', '$2a$12$somehashedvalue', NOW()),
('John Doe', 'john.doe@example.com', '$2a$12$anotherhash', NOW()),
('Jane Smith', 'jane.smith@example.com', '$2a$12$hash3', NOW());
