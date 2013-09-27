package backends_test

import (
    BA "github.com/darkhelmet/balance/backends"
    . "launchpad.net/gocheck"
    "testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var (
    _   = Suite(&S{})
    a   = "1.2.3.4"
    b   = "2.3.4.5"
)

func (s *S) TestLen(c *C) {
    r := BA.NewRoundRobin([]string{a, b})
    c.Assert(r.Len(), Equals, 2)
}

func (s *S) TestChoose(c *C) {
    r := BA.NewRoundRobin([]string{a, b})
    c.Assert(r.Choose().String(), Equals, a)
    c.Assert(r.Choose().String(), Equals, b)
    c.Assert(r.Choose().String(), Equals, a)
    c.Assert(r.Choose().String(), Equals, b)
}

func (s *S) TestChooseEmpty(c *C) {
    r := BA.NewRoundRobin([]string{})
    c.Assert(r.Choose(), Equals, nil)
}

func (s *S) TestAdd(c *C) {
    r := BA.NewRoundRobin([]string{a})
    c.Assert(r.Choose().String(), Equals, a)
    c.Assert(r.Choose().String(), Equals, a)
    r.Add(b)
    c.Assert(r.Choose().String(), Equals, b)
    c.Assert(r.Choose().String(), Equals, a)
}

func (s *S) TestAddEmpty(c *C) {
    r := BA.NewRoundRobin([]string{})
    c.Assert(r.Len(), Equals, 0)
    r.Add(a)
    c.Assert(r.Len(), Equals, 1)
    c.Assert(r.Choose().String(), Equals, a)
}

func (s *S) TestRemove(c *C) {
    r := BA.NewRoundRobin([]string{a, b})
    c.Assert(r.Len(), Equals, 2)
    c.Assert(r.Choose().String(), Equals, a)
    c.Assert(r.Choose().String(), Equals, b)
    r.Remove(b)
    c.Assert(r.Len(), Equals, 1)
    c.Assert(r.Choose().String(), Equals, a)
    c.Assert(r.Choose().String(), Equals, a)
    r.Remove(a)
    c.Assert(r.Len(), Equals, 0)
}

func (s *S) TestRemoveEmpty(c *C) {
    r := BA.NewRoundRobin([]string{})
    c.Assert(r.Len(), Equals, 0)
    r.Remove(a)
    c.Assert(r.Len(), Equals, 0)
    r.Add(a)
    c.Assert(r.Len(), Equals, 1)
    r.Remove(a)
    c.Assert(r.Len(), Equals, 0)
    r.Remove(a)
}
