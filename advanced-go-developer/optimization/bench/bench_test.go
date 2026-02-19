package bench

import (
	"testing"
)

func BenchmarkLoadBalancers(b *testing.B) {
	n := 50
	reqsN := 1000
	var resultVacuum *Connection
	conns := makeConnections(n)
	b.Run("chan", func(b *testing.B) {
		balancerChan := NewLoadBalancerChan(conns)
		balancerChan.Init()
		defer balancerChan.Close()
		for j := 0; j < b.N; j++ {
			for i := 0; i < reqsN; i++ {
				c := balancerChan.NextConn()
				resultVacuum = c
			}
		}
	})

	b.Run("atomic", func(b *testing.B) {

		for j := 0; j < b.N; j++ {
			balancerA := NewLoadBalancerAtomic(conns)
			for i := 0; i < reqsN; i++ {
				c := balancerA.NextConn()
				_ = c
			}
		}
	})

	b.Run("mutex", func(b *testing.B) {
		balancerM := NewLoadBalancerMutex(conns)
		for j := 0; j < b.N; j++ {

			for i := 0; i < reqsN; i++ {
				c := balancerM.NextConn()
				_ = c
			}
		}
	})
	b.StopTimer()

	_ = resultVacuum
}

func makeConnections(n int) []*Connection {
	var conns []*Connection
	for i := 0; i < n; i++ {
		conns = append(conns, &Connection{})
	}

	return conns
}
