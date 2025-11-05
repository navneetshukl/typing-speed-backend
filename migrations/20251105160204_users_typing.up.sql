CREATE TABLE userTypingData (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(100) NOT NULL,
    total_error INT NOT NULL,
    total_words INT NOT NULL,
    typed_words INT NOT NULL,
    total_time INT NOT NULL,
    total_time_taken_by_user INT NOT NULL,
    wpm INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
