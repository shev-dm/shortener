CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    code VARCHAR(16) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO links (original_url, code) VALUES
    ('https://www.google.com', 'google'),
    ('https://github.com', 'github');
