CREATE TABLE request_logs (
    id SERIAL PRIMARY KEY,
    method TEXT,
    path TEXT,
    ip_address TEXT,
    user_agent TEXT,
    platform TEXT,
    duration_ms INTEGER,
    created_at TIMESTAMP DEFAULT (timezone('America/New_York', now()))
);
