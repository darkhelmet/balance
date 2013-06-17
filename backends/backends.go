package backends

type Backends interface {
    Choose() string
    Len() int
}
