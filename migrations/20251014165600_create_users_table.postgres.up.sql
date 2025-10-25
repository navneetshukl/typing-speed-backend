CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO users (name, email, password, created_at)
VALUES
('Navneet Shukla', 'navneet@example.com', 'password123', NOW()),
('John Doe', 'john.doe@example.com', 'securepass', NOW()),
('Jane Smith', 'jane.smith@example.com', 'mypassword', NOW());
