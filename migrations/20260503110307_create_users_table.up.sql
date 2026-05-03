CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,

    first_name TEXT NOT NULL,
    second_name TEXT,

    password_hash TEXT NOT NULL,
    
    role TEXT NOT NULL DEFAULT 'user',

    is_email_verified BOOLEAN DEFAULT 0,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);