package backends

import (
    "sync"
)

type simpleBackends struct {
    backends []string
    n        int
    l        sync.Mutex
}

func NewSimpleBackends(ba []string) Backends {
    return &simpleBackends{backends: ba}
}

func (b *simpleBackends) Choose() string {
    b.l.Lock()
    defer b.l.Unlock()
    idx := b.n % len(b.backends)
    b.n++
    return b.backends[idx]
}

func (b *simpleBackends) Len() int {
    return len(b.backends)
}
