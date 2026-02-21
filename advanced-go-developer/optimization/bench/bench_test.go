package bench

import (
	"testing"
)

func BenchmarkLoadBalancers(b *testing.B) {
	n := 50
	targets := makeTargets(n)

	b.Run("chan", func(b *testing.B) {
		balancer := NewLoadBalancerChan(targets)
		balancer.Init()
		defer balancer.Close()

		for b.Loop() {
			t := balancer.NextConn()
			_ = t
		}

	})

	b.Run("atomic", func(b *testing.B) {
		balancer := NewLoadBalancerAtomic(targets)

		for b.Loop() {
			t := balancer.NextConn()
			_ = t
		}

	})

	b.Run("mutex", func(b *testing.B) {
		balancer := NewLoadBalancerMutex(targets)

		for b.Loop() {
			t := balancer.NextConn()
			_ = t
		}
	})

}

func makeTargets(n int) []*Target {
	var targets []*Target
	for i := 0; i < n; i++ {
		targets = append(targets, &Target{})
	}

	return targets
}
