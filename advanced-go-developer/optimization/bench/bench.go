package bench

import (
    "sync"
    "sync/atomic"
)

type Target struct{}

type LoadBalancer interface {
    NextTarget() *Target
}

type LoadBalancerChan struct {
    targets []*Target
    ch    chan *Target
    stop  chan struct{}
}

func NewLoadBalancerChan(targets []*Target) *LoadBalancerChan {
    return &LoadBalancerChan{targets: targets, ch: make(chan *Target), stop: make(chan struct{})}
}

func (b *LoadBalancerChan) Init() {
    go b.worker()
}

func (b *LoadBalancerChan) Close() {
    b.stop <- struct{}{}
}

func (b *LoadBalancerChan) worker() {
    for i := 0; ; {
        select {
        case b.ch <- b.targets[i]:
            i++
            if i == len(b.targets) {
                i = 0
            }

        case <-b.stop:
            return
        }
    }
}

func (b *LoadBalancerChan) NextConn() *Target {
    return <-b.ch
}

type LoadBalancerAtomic struct {
    targets   []*Target
    counter uint32
}

func NewLoadBalancerAtomic(targets []*Target) *LoadBalancerAtomic {
    return &LoadBalancerAtomic{targets: targets}
}

func (b *LoadBalancerAtomic) NextConn() *Target {
    i := atomic.AddUint32(&b.counter, 1) % uint32(len(b.targets))
    return b.targets[i]
}

type LoadBalancerMutex struct {
    targets   []*Target
    counter int
    mu      sync.Mutex
}

func NewLoadBalancerMutex(conns []*Target) *LoadBalancerMutex {
    return &LoadBalancerMutex{targets: conns}
}

func (b *LoadBalancerMutex) NextConn() *Target {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.counter = (b.counter + 1) % len(b.targets)
    return b.targets[b.counter]
}
