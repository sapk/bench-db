package cassandra

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/gocql/gocql"
	"github.com/moby/moby/client"
	"github.com/stretchr/testify/assert"

	"github.com/sapk/bench-db/database"
)

//Cassandra represent a cassandra type
type Cassandra struct {
	dockerClient      *client.Client
	dockerContainerID string
	session           *gocql.Session
}

//Name return cassandra
func (c *Cassandra) Name() string {
	return "Cassandra"
}

//Setup setup the cassandra database
func (c *Cassandra) Setup(tb testing.TB) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if !assert.NoError(tb, err, "Failed to connect to docker") {
		assert.FailNow(tb, "Failed setup docker")
	}
	c.dockerClient = cli

	//_, err = cli.ImagePull(ctx, "docker.io/library/cassandra", types.ImagePullOptions{})
	_, err = cli.ImagePull(ctx, "cassandra", types.ImagePullOptions{})
	if !assert.NoError(tb, err, "Failed to pull docker image cassandra") {
		assert.FailNow(tb, "Failed setup docker")
	}

	cont, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "cassandra",
	}, nil, nil, "")
	assert.NoError(tb, err, "Failed to create docker container cassandra")

	//tb.Log(cont.Name, cont.State)
	err = cli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{})
	assert.NoError(tb, err, "Failed to start docker container cassandra")

	c.dockerContainerID = cont.ID
	/* TODO
	var buf bytes.Buffer
	err = client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: &buf,
		Stderr:       true,
		Stdout:       true,
		Stream:       true,
	})
	assert.Nil(tb, err, "Failed to attach to docker container cassandra")
	*/
	time.Sleep(60 * time.Second)
}

//Init the cassandra database structure
func (c *Cassandra) Init(tb testing.TB) {
	session := c.getSession(tb)

	tb.Log("Setting up keyspace benchmark")
	err := session.Query(`CREATE KEYSPACE benchmark WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`).Exec()
	if err != nil && !strings.Contains(err.Error(), "no response received from cassandra within timeout period") {
		assert.NoError(tb, err, "Failed to create Cassandra keyspace benchmark")
	}
	time.Sleep(5 * time.Second)
	err = session.Query(`CREATE TABLE benchmark.tweet(timeline text, id UUID, text text, PRIMARY KEY(id))`).Exec()
	if err != nil && !strings.Contains(err.Error(), "no response received from cassandra within timeout period") {
		assert.NoError(tb, err, "Failed to create Cassandra table benchmark.tweet")
	}
	time.Sleep(15 * time.Second)
}

//Clean the cassandra database structure
func (c *Cassandra) Clean(tb testing.TB) {
	session := c.getSession(tb)
	tb.Log("Clearing Cassandra eyspace benchmark")
	err := session.Query(`DROP KEYSPACE IF EXISTS benchmark`).Exec()
	if err != nil && !strings.Contains(err.Error(), "no response received from cassandra within timeout period") {
		assert.NoError(tb, err, "Failed to drop keyspace")
	}
	time.Sleep(15 * time.Second)
}

//Destroy the cassandra container
func (c *Cassandra) Destroy(tb testing.TB) {
	ctx := context.Background()
	timeout := 30 * time.Second
	err := c.dockerClient.ContainerStop(ctx, c.dockerContainerID, &timeout)
	assert.Nil(tb, err, "Failed to stop container")

	err = c.dockerClient.ContainerRemove(ctx, c.dockerContainerID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		//	RemoveLinks:   true,
		Force: true,
	})
	assert.NoError(tb, err, "Failed to remove container")
	//time.Sleep(5 * time.Second)
}

//Benchs return bench list
func (c *Cassandra) Benchs(tb testing.TB) []database.Bench {
	s := c.getSession(tb)
	return []database.Bench{
		BenchAddTweet{
			session: s,
		},
	}
}

func (c *Cassandra) getSession(tb testing.TB) *gocql.Session {
	//TODO manage multiple IP
	//TODO check if session is closed
	if c.session != nil {
		return c.session
	}

	ctx := context.Background()
	cont, err := c.dockerClient.ContainerInspect(ctx, c.dockerContainerID)
	assert.Nil(tb, err, "Failed to inspect container")
	nets := ""
	for _, n := range cont.NetworkSettings.Networks {
		nets += n.IPAddress //TODO join
	}
	tb.Log("DEBUG", nets)
	cluster := gocql.NewCluster(nets)
	session, err := cluster.CreateSession()
	assert.NoError(tb, err, "Failed to connect to Cassandra cluster")
	assert.NotNil(tb, session, "Failed to connect to Cassandra cluster")
	return session
}

//BenchAddTweet represent a cassandra bench
type BenchAddTweet struct {
	session *gocql.Session
}

//Run the benchmark
func (b BenchAddTweet) Run(tb testing.TB) {
	//TODO and use bench N
	tb.Skip()
}
