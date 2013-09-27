package backends

import (
    "container/ring"
    "sync"
)

type RoundRobin struct {
    r   *ring.Ring
    l   sync.RWMutex
}

func NewRoundRobin(strs []string) Backends {
    r := ring.New(len(strs))
    for _, s := range strs {
        r.Value = &backend{s}
        r = r.Next()
    }
    return &RoundRobin{r: r}
}

func init() {
    factories["round-robin"] = NewRoundRobin
}

func (rr *RoundRobin) Len() int {
    rr.l.RLock()
    defer rr.l.RUnlock()
    return rr.r.Len()
}

func (rr *RoundRobin) Choose() Backend {
    rr.l.Lock()
    defer rr.l.Unlock()
    if rr.r == nil {
        return nil
    }
    n := rr.r.Value.(*backend)
    rr.r = rr.r.Next()
    return n
}

func (rr *RoundRobin) Add(s string) {
    rr.l.Lock()
    defer rr.l.Unlock()
    nr := &ring.Ring{Value: &backend{s}}
    if rr.r == nil {
        rr.r = nr
    } else {
        rr.r = rr.r.Link(nr).Next()
    }
}

func (rr *RoundRobin) Remove(s string) {
    rr.l.Lock()
    defer rr.l.Unlock()
    r := rr.r
    if rr.r.Len() == 1 {
        rr.r = ring.New(0)
        return
    }

    for i := rr.r.Len(); i > 0; i-- {
        r = r.Next()
        ba := r.Value.(*backend)
        if s == ba.String() {
            rr.r = r.Unlink(1)
            return
        }
    }
}
