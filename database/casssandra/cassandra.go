package cassandra

import (
	"strings"
	"testing"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
)

//Cassandra represent a cassandra type
type Cassandra struct {
	dockerClient    *docker.Client
	dockerContainer *docker.Container
	session         *gocql.Session
}

//Name return cassandra
func (c *Cassandra) Name() string {
	return "Cassandra"
}

//Setup setup the cassandra database
func (c *Cassandra) Setup(tb testing.TB) {
	endpoint := "unix:///var/run/docker.sock" //TODO sue env Docker endpoint //	client, _ := docker.NewClientFromEnv()
	client, err := docker.NewClient(endpoint)
	if !assert.NoError(tb, err, "Failed to connect to docker") {
		assert.FailNow(tb, "Failed setup docker")
	}
	c.dockerClient = client

	err = client.PullImage(docker.PullImageOptions{
		Repository: "cassandra",
	}, docker.AuthConfiguration{})
	if !assert.NoError(tb, err, "Failed to pull docker image cassandra") {
		assert.FailNow(tb, "Failed setup docker")
	}

	container, err := client.CreateContainer(docker.CreateContainerOptions{
		//Name: "cassandra_database_bench",
		Config: &docker.Config{
			Image: "cassandra",
		},
	})
	assert.NoError(tb, err, "Failed to create docker container cassandra")

	tb.Log(container.Name, container.State)
	//container.
	err = client.StartContainer(container.ID, &docker.HostConfig{})
	assert.NoError(tb, err, "Failed to start docker container cassandra")

	c.dockerContainer = container
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
	time.Sleep(15 * time.Second)
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
	time.Sleep(5 * time.Second)
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
	err := c.dockerClient.StopContainer(c.dockerContainer.ID, 30)
	assert.Nil(tb, err, "Failed to stop container")
	err = c.dockerClient.RemoveContainer(docker.RemoveContainerOptions{
		ID:            c.dockerContainer.ID,
		Force:         true,
		RemoveVolumes: true,
	})
	assert.NoError(tb, err, "Failed to remove container")
	//time.Sleep(5 * time.Second)
}

func (c *Cassandra) getSession(tb testing.TB) *gocql.Session {
	//TODO manage multiple IP
	//TODO check if session is closed
	if c.session != nil {
		return c.session
	}
	cluster := gocql.NewCluster(c.dockerContainer.NetworkSettings.IPAddress)
	session, err := cluster.CreateSession()
	assert.NoError(tb, err, "Failed to connect to Cassandra cluster")
	assert.NotNil(tb, session, "Failed to connect to Cassandra cluster")
	return session
}
