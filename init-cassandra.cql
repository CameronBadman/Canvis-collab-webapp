CREATE KEYSPACE IF NOT EXISTS myapp_dev
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE myapp_dev;

CREATE TABLE IF NOT EXISTS canvases (
    id UUID PRIMARY KEY,
    user_id UUID,
    name TEXT,
    svg_data TEXT,
    created_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS ON canvases (user_id);