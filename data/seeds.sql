CREATE TABLE IF NOT EXISTS short_urls (
    id VARCHAR(10) PRIMARY KEY,
    original_url VARCHAR(200) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Insert test data
INSERT INTO short_urls (id, original_url, created_at) VALUES
    ('abc123', 'https://www.google.com', NOW()),
    ('def456', 'https://www.github.com', NOW()),
    ('ghi789', 'https://www.example.com', NOW())
ON CONFLICT (id) DO NOTHING;