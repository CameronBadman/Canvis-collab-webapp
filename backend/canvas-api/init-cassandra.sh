#!/bin/bash

echo "Initializing Cassandra..."

cqlsh $CASSANDRA_HOST $CASSANDRA_PORT -e "
CREATE KEYSPACE IF NOT EXISTS $CASSANDRA_KEYSPACE
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

USE $CASSANDRA_KEYSPACE;

CREATE TABLE IF NOT EXISTS canvases (
    id UUID PRIMARY KEY,
    user_id UUID,
    name TEXT,
    svg_data TEXT,
    created_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS ON canvases (user_id);"

echo "Cassandra initialization completed"
