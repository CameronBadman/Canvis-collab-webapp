-- schema.cql
CREATE KEYSPACE IF NOT EXISTS canvas_collab
    WITH replication = {'class': 'NetworkTopologyStrategy', 'datacenter1': 3};

USE canvas_collab;

CREATE TABLE IF NOT EXISTS users (
                                     user_id UUID PRIMARY KEY,
                                     username TEXT,
                                     email TEXT
);

CREATE TABLE IF NOT EXISTS canvases (
                                        user_id UUID,
                                        canvas_id UUID,
                                        canvas_name TEXT,
                                        created_at TIMESTAMP,
                                        PRIMARY KEY (user_id, canvas_id)
);

CREATE TABLE IF NOT EXISTS svg_data (
                                        canvas_id UUID,
                                        svg_id UUID,
                                        svg_content TEXT,
                                        created_at TIMESTAMP,
                                        action TEXT,
                                        PRIMARY KEY (canvas_id, svg_id)
);
