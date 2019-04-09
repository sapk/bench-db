package main

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/go-kivik/couchdb" // The CouchDB driver
	"github.com/go-kivik/kivik"   // Development version of Kivik
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	couchdbSetupLock sync.Mutex
	couchdbState     = "unknown"
)

func setupCouchDB(tb testing.TB) {
	couchdbSetupLock.Lock()
	defer couchdbSetupLock.Unlock()

	if couchdbState == "ready" {
		tb.Log("Skipping setup as couchdb benchmark ready")
		return
	}

	client, err := kivik.New("couch", "http://"+os.Getenv("COUCHDB_IP")+"/")
	assert.Nil(tb, err, "Failed to connect to couchdb")
	assert.NotNil(tb, client, "Failed to connect to couchdb")
	err = client.Authenticate(context.TODO(), couchdb.BasicAuth("admin", "password"))
	assert.Nil(tb, err, "Failed to auth to couchdb")

	db := client.DB(context.TODO(), "tweets", nil)

	clearCouchDB(tb, db)
	tb.Log("Setting up couchdb benchmark")

	//TODO setup database tweets;
	time.Sleep(5 * time.Second)

	couchdbState = "ready"
	client.Close(context.TODO())
}

func clearCouchDB(tb testing.TB, session *kivik.DB) {
	//TODO
}

func BenchmarkCouchDB(b *testing.B) {
	b.StopTimer()
	if os.Getenv("COUCHDB_IP") == "" {
		b.Skip("Env. variable COUCHDB_IP not set -> Skipping couchdb tests")
	}

	setupCouchDB(b)

	client, err := kivik.New("couch", "http://"+os.Getenv("COUCHDB_IP")+"/")
	assert.Nil(b, err, "Failed to connect to couchdb")
	assert.NotNil(b, client, "Failed to connect to couchdb")
	err = client.Authenticate(context.TODO(), couchdb.BasicAuth("admin", "password"))
	assert.Nil(b, err, "Failed to auth to couchdb")

	db := client.DB(context.TODO(), "tweets", nil)
	defer client.Close(context.TODO())

	b.StartTimer()

	b.Run("AddTweet", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			// insert a tweet
			_, err := db.Put(context.TODO(), uuid.New().String(), map[string]interface{}{
				"timeline": "me",
				"text":     "hello world",
			})
			assert.Nil(b, err, "Failed to add data to cluster")
		}
	})
}
