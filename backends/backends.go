package backends

import (
    "log"
)

type Backends interface {
    Choose() string
    Len() int
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
