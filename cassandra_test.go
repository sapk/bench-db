package main

import (
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
)

var (
	cassandraSetupLock sync.Mutex
	cassandraState     = "unknown"
)

func setupCassandra(tb testing.TB) {
	cassandraSetupLock.Lock()
	defer cassandraSetupLock.Unlock()

	if cassandraState == "ready" {
		tb.Log("Skipping setup as keyspace benchmark ready")
		return
	}

	//TODO setup a real cluster of at least 3 nodes
	//Setup cassandra database
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_IP"))
	session, err := cluster.CreateSession()
	assert.Nil(tb, err, "Failed to connect to cluster")
	assert.NotNil(tb, session, "Failed to connect to cluster")

	clearCassandra(tb, session)
	tb.Log("Setting up keyspace benchmark")
	err = session.Query(`CREATE KEYSPACE benchmark WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`).Exec()
	if err != nil && !strings.Contains(err.Error(), "no response received from cassandra within timeout period") {
		assert.Nil(tb, err, "Failed to create keyspace")
	}
	time.Sleep(3 * time.Second)
	err = session.Query(`CREATE TABLE benchmark.tweet(timeline text, id UUID, text text, PRIMARY KEY(id))`).Exec()
	if err != nil && !strings.Contains(err.Error(), "no response received from cassandra within timeout period") {
		assert.Nil(tb, err, "Failed to create table")
	}
	time.Sleep(3 * time.Second)
	/*
		err = session.Query(`CREATE INDEX benchmark.tweet(timeline)`).Exec()
		if err != nil && !strings.Contains(err.Error(), "no response received from cassandra within timeout period") {
			assert.Nil(tb, err, "Failed to create index")
		}
		time.Sleep(3 * time.Second)
	*/
	//TODO manage error and time
	cassandraState = "ready"
	session.Close()
}

func clearCassandra(tb testing.TB, session *gocql.Session) {
	tb.Log("Clearing keyspace benchmark")
	err := session.Query(`DROP KEYSPACE IF EXISTS benchmark`).Exec()
	if err != nil && !strings.Contains(err.Error(), "no response received from cassandra within timeout period") {
		assert.Nil(tb, err, "Failed to drop keyspace")
	}
	time.Sleep(15 * time.Second)
	//TODO manage error and time
}

/*
func TestMain(m *testing.M) {

	//cluster.ConnectTimeout = 15 * time.Second
	//cluster.Timeout = 15 * time.Second //Increase timeout for init
	//assert.Nil(b, err, "Failed to init database")
	time.Sleep(5 * time.Second)

	os.Exit(m.Run())
}
*/

func BenchmarkCassandra(b *testing.B) {
	b.StopTimer()
	if os.Getenv("CASSANDRA_IP") == "" {
		b.Skip("Env. variable CASSANDRA_IP not set -> Skipping cassandra tests")
	}

	setupCassandra(b)
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_IP"))
	//cluster.Keyspace = "benchmark"
	//cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	assert.Nil(b, err, "Failed to init connexion to cluster")
	defer session.Close()

	b.StartTimer()

	b.Run("AddTweet", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			// insert a tweet
			//TODO use ripe atlas
			err = session.Query(`INSERT INTO benchmark.tweet (timeline, id, text) VALUES (?, ?, ?)`, "me", gocql.TimeUUID(), "hello world").Exec()
			assert.Nil(b, err, "Failed to add data to cluster")
		}
	})
}
