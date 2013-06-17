package backends

import (
    "sync"
)

type roundRobin struct {
    backends []string
    n        int
    l        sync.Mutex
}

func NewRoundRobin(specs []string) Backends {
    return &roundRobin{backends: specs}
}

func init() {
    factories["round-robin"] = NewRoundRobin
}

func (b *roundRobin) Choose() string {
    b.l.Lock()
    defer b.l.Unlock()
    idx := b.n % len(b.backends)
    b.n++
    return b.backends[idx]
}

func (b *roundRobin) Len() int {
    return len(b.backends)
}
