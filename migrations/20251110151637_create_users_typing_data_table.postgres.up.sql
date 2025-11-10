CREATE TABLE user_typing_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(100) NOT NULL,
    total_error INT NOT NULL,
    total_words INT NOT NULL,
    typed_words INT NOT NULL,
    total_time INT NOT NULL,
    total_time_taken_by_user INT NOT NULL,
    wpm INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_email
        FOREIGN KEY (email)
        REFERENCES users(email)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
