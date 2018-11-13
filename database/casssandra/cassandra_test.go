package cassandra

import "testing"

func TestCassandra(t *testing.T) {
	c := Cassandra{}
	c.Setup(t)
	//defer c.Destroy(t)
	c.Init(t)

	for _, b := range c.Benchs(t) {
		b.Run(t)
	}

	c.Clean(t)
	c.Destroy(t)
}
