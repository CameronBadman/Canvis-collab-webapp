#!/bin/bash

set -e

host="$CASSANDRA_HOST"
port="$CASSANDRA_PORT"

until nc -z $host $port; do
  echo "Cassandra is unavailable - sleeping"
  sleep 5
done

echo "Cassandra is up - executing command"
exec "$@"
