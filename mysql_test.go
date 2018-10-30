package main

import (
	"database/sql"
	"os"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var (
	mysqlSetupLock sync.Mutex
	mysqlState     = "unknown"
)

func setupMySQL(tb testing.TB) {
	mysqlSetupLock.Lock()
	defer mysqlSetupLock.Unlock()

	if mysqlState == "ready" {
		tb.Log("Skipping setup of mysql")
		return
	}

	//Setup mysql database
	db, err := sql.Open("mysql", os.Getenv("MYSQL_URL"))
	assert.Nil(tb, err, "Failed to connect to mysql")

	clearMySQL(tb, db)
	tb.Log("Setting up keyspace benchmark")
	_, err = db.Exec(`CREATE DATABASE benchmark;`)
	assert.Nil(tb, err, "Failed to create database")
	_, err = db.Exec(`USE benchmark;`)
	assert.Nil(tb, err, "Failed to move to table")

	_, err = db.Exec(`CREATE TABLE tweet (id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY, timeline VARCHAR(30), text VARCHAR(30))`)
	//TODO index on  timeline
	assert.Nil(tb, err, "Failed to create table")
	mysqlState = "ready"
	db.Close()
}

func clearMySQL(tb testing.TB, db *sql.DB) {
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

func BenchmarkMySQL(b *testing.B) {
	b.StopTimer()
	if os.Getenv("MYSQL_URL") == "" {
		b.Skip("Env. variable MYSQL_URL not set -> Skipping MySQL tests")
	}

	setupMySQL(b)
	db, err := sql.Open("mysql", os.Getenv("MYSQL_URL"))
	assert.Nil(b, err, "Failed to connect to mysql")
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
