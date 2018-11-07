package main

import (
	"database/sql"
	"os"
	"sync"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/stretchr/testify/assert"
)

var (
	mssqlSetupLock sync.Mutex
	mssqlState     = "unknown"
)

func setupMSSQL(tb testing.TB) {
	mssqlSetupLock.Lock()
	defer mssqlSetupLock.Unlock()

	if mssqlState == "ready" {
		tb.Log("Skipping setup of mssql")
		return
	}

	//Setup mssql database
	db, err := sql.Open("mssql", os.Getenv("MSSQL_URL"))
	assert.Nil(tb, err, "Failed to connect to mssql")
	assert.NotNil(tb, db, "Failed to connect to mssql")

	clearMSSQL(tb, db)
	tb.Log("Setting up database benchmark")
	_, err = db.Exec(`CREATE DATABASE benchmark;`)
	assert.Nil(tb, err, "Failed to create database")
	_, err = db.Exec(`USE benchmark;`)
	assert.Nil(tb, err, "Failed to move to table")

	_, err = db.Exec(`CREATE TABLE tweet (id INT IDENTITY(1,1) PRIMARY KEY, timeline VARCHAR(30), text VARCHAR(30))`)
	//TODO index on  timeline
	assert.Nil(tb, err, "Failed to create table")
	mssqlState = "ready"
	db.Close()
}

func clearMSSQL(tb testing.TB, db *sql.DB) {
	tb.Log("Clearing database benchmark")
	_, err := db.Exec(`DROP DATABASE IF EXISTS benchmark`)
	assert.Nil(tb, err, "Failed to drop database")
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

func BenchmarkMSSQL(b *testing.B) {
	b.StopTimer()
	if os.Getenv("MSSQL_URL") == "" {
		b.Skip("Env. variable MSSQL_URL not set -> Skipping MSSQL tests")
	}

	setupMSSQL(b)
	db, err := sql.Open("mssql", os.Getenv("MSSQL_URL"))
	assert.Nil(b, err, "Failed to connect to mssql")
	defer db.Close()

	_, err = db.Exec(`USE benchmark;`)
	assert.Nil(b, err, "Failed to move to table")

	b.StartTimer()

	b.Run("AddTweet", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			// insert a tweet
			//TODO use ripe atlas
			_, err = db.Exec(`INSERT INTO tweet (timeline, text) VALUES (?, ?)`, "me", "hello world")
			assert.Nil(b, err, "Failed to add data to cluster")
		}
	})
}
