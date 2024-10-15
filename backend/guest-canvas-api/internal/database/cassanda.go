package database

import (
	"os"

	"github.com/gocql/gocql"
)

func InitCassandra() (*gocql.Session, error) {
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_ADDR"))
	cluster.Keyspace = os.Getenv("CASSANDRA_KEYSPACE")
	cluster.Consistency = gocql.Quorum
	return cluster.CreateSession()
}
