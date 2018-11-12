package cassandra

import "testing"

func TestCassandra(t *testing.T) {
	c := Cassandra{}
	c.Setup(t)
	c.Init(t)

	c.Clean(t)
	c.Destroy(t)
}
