-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(100),
    total_error INT,
    total_words INT,
    typed_words INT,
    total_time INT,
    total_time_taken_by_user INT,
    wpm INT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Insert initial data
INSERT INTO users (
    user_id,
    total_error,
    total_words,
    typed_words,
    total_time,
    total_time_taken_by_user,
    wpm,
    created_at
)
VALUES
('user_001', 5, 100, 95, 60, 55, 85, NOW()),
('user_002', 3, 120, 118, 70, 65, 90, NOW()),
('user_003', 8, 80, 72, 50, 48, 75, NOW());
