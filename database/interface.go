package database

import "testing"

//Database represent database
type Database interface {
	Name() string
	Setup(testing.TB)
	Init(testing.TB)
	Clean(testing.TB)
	Destroy(testing.TB)
	Benchs() []Bench
}

//Benchs a bench to run on database
type Bench interface {
	Run(testing.TB)
}
