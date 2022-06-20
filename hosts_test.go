package main

import (
	"testing"
	
	"github.com/matryer/is"
)

func TestHostRingGet(t *testing.T) {
	is := is.New(t)

	addrs := []string{ "host1", "host2", "host3" }
	ring, err := newHostRing(addrs)

	is.NoErr(err)

	for i := 0; i < len(addrs); i++ {
		host := ring.get()
		is.Equal(host.addr, addrs[i])
		is.True(host.alive)
	}
}
