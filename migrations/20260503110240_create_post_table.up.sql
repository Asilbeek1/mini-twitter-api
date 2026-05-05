CREATE TABLE posts (
    id SERIAL PRIMARY KEY,

    user_id INTEGER NOT NULL,

    title TEXT NOT NULL,
    description TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id)
);