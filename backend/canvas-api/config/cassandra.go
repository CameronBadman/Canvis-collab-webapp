package config

import (
	"github.com/gocql/gocql"
	"log"
)

// SetupCassandraSession initializes a connection to Cassandra
func SetupCassandraSession() (*gocql.Session, error) {
	cluster := gocql.NewCluster("cassandra.db.svc.cluster.local") // Use the service name for internal K8s access
	cluster.Keyspace = "canvas_collab"                            // Update with your keyspace name
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		log.Printf("Failed to connect to Cassandra: %v", err)
		return nil, err
	}
	return session, nil
}
