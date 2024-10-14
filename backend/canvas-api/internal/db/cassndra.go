package db

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gocql/gocql"
)

func InitCassandra() (*gocql.Session, error) {
	host := os.Getenv("CASSANDRA_HOST")
	if host == "" {
		return nil, fmt.Errorf("CASSANDRA_HOST is not set")
	}

	port, err := strconv.Atoi(os.Getenv("CASSANDRA_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid CASSANDRA_PORT: %v", err)
	}

	keyspace := os.Getenv("CASSANDRA_KEYSPACE")
	if keyspace == "" {
		return nil, fmt.Errorf("CASSANDRA_KEYSPACE is not set")
	}

	cluster := gocql.NewCluster(host)
	cluster.Port = port
	cluster.Keyspace = keyspace

	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.Timeout = time.Second * 10
	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{
		Min:        time.Second,
		Max:        time.Minute,
		NumRetries: 5,
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	// Create the keyspace if it doesn't exist
	err = session.Query(fmt.Sprintf(`
		CREATE KEYSPACE IF NOT EXISTS %s
		WITH REPLICATION = { 'class': 'SimpleStrategy', 'replication_factor': 1 }
	`, keyspace)).Exec()
	if err != nil {
		return nil, fmt.Errorf("failed to create keyspace: %v", err)
	}

	return session, nil
}
