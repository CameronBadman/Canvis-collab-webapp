#!/bin/bash
set -e

# Wait for Cassandra to be ready
until cqlsh -e "describe cluster"; do
    echo "Waiting for Cassandra..."
    sleep 5
done

# Create the keyspace if it doesn't exist
cqlsh -e "CREATE KEYSPACE IF NOT EXISTS cassandra-Db WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};"
echo "Cassandra initialization completed"
