package bench

import (
	"testing"
)

func BenchmarkLoadBalancers(b *testing.B) {
	n := 50
	reqsN := 1000
	conns := makeConnections(n)

	b.Run("chan", func(b *testing.B) {
		balancer := NewLoadBalancerChan(conns)
		balancer.Init()
		defer balancer.Close()

		for b.Loop() {

			for i := 0; i < reqsN; i++ {
				c := balancer.NextConn()
				_ = c
			}

		}
	})

	b.Run("atomic", func(b *testing.B) {
		balancer := NewLoadBalancerAtomic(conns)
		for b.Loop() {

			for i := 0; i < reqsN; i++ {
				c := balancer.NextConn()
				_ = c
			}

		}
	})

	b.Run("mutex", func(b *testing.B) {
		balancer := NewLoadBalancerMutex(conns)

		for b.Loop() {

			for i := 0; i < reqsN; i++ {
				c := balancer.NextConn()
				_ = c
			}

		}
	})

}

func makeConnections(n int) []*Target {
	var conns []*Target
	for i := 0; i < n; i++ {
		conns = append(conns, &Target{})
	}

	return conns
}
