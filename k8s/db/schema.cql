CREATE KEYSPACE IF NOT EXISTS canvas_collab
    WITH replication = {'class': 'NetworkTopologyStrategy', 'datacenter1': 3};

USE canvas_collab;

-- Updated users table with TEXT for user_id (Cognito sub)
CREATE TABLE IF NOT EXISTS users (
                                     user_id TEXT PRIMARY KEY,  -- Using Cognito sub as the primary key
                                     username TEXT,
                                     email TEXT
);

-- Create a user-defined type (UDT) for SVG data
CREATE TYPE IF NOT EXISTS svg_data_type (
                                            svg_id UUID,
                                            svg_content TEXT,
                                            created_at TIMESTAMP,
                                            action TEXT
                                        );

-- Update the canvases table to include a frozen set of svg_data_type
CREATE TABLE IF NOT EXISTS canvases (
                                        user_id TEXT,                  -- Reference Cognito sub as the foreign key
                                        canvas_id UUID,                -- Unique identifier for each canvas
                                        canvas_name TEXT,              -- Name of the canvas
                                        created_at TIMESTAMP,          -- Timestamp when the canvas was created
                                        svg_data FROZEN<LIST<FROZEN<svg_data_type>>>, -- List of frozen SVG data
                                        PRIMARY KEY (user_id, canvas_id)
);
