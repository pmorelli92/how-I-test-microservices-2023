CREATE TABLE IF NOT EXISTS posts(
    id VARCHAR PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
)
