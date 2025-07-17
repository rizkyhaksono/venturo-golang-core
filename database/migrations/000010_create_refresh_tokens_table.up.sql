CREATE TABLE refresh_tokens (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    hashed_token VARCHAR(255) NOT NULL UNIQUE
);