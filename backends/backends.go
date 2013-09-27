package backends

import (
    "log"
)

type backend struct {
    hostname string
}

func (b *backend) String() string {
    return b.hostname
}

type Backend interface {
    String() string
}

type Backends interface {
    Choose() Backend
    Len() int
    Add(string)
    Remove(string)
}

type Factory func([]string) Backends

var factories = make(map[string]Factory)

func Build(algorithm string, specs []string) Backends {
    factory, found := factories[algorithm]
    if !found {
        log.Fatalf("balance algorithm %s not supported", algorithm)
    }
    return factory(specs)
}
