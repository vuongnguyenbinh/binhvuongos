CREATE TABLE settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO settings (key, value) VALUES
    ('smtp_host', ''),
    ('smtp_port', '587'),
    ('smtp_user', ''),
    ('smtp_pass', ''),
    ('smtp_from', ''),
    ('notion_api_key', ''),
    ('notion_database_ids', '{}'),
    ('n8n_webhook_url', ''),
    ('google_oauth_client_id', ''),
    ('google_oauth_client_secret', ''),
    ('unsplash_keywords', 'vietnam,hanoi,nature');
