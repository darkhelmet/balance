package main

import (
    BA "github.com/darkhelmet/balance/backends"
    "github.com/gonuts/commander"
    "log"
    "net/http"
    "net/http/httputil"
)

func httpBalance(bind string, backends BA.Backends) {
    log.Println("using http balancing")
    proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
        req.URL.Scheme = "http"
        req.URL.Host = backends.Choose()
        req.Header.Add(XRealIP, RealIP(req))
    }}
    log.Printf("listening on %s, balancing %d backends", bind, backends.Len())
    err := http.ListenAndServe(bind, proxy)
    if err != nil {
        log.Fatalf("failed to bind: %s", err)
    }
}

func init() {
    fs := newFlagSet("http")

    cmd.Commands = append(cmd.Commands, &commander.Command{
        UsageLine: "http [options] <backend> [<more backends>]",
        Short:     "performs http based load balancing",
        Flag:      *fs,
        Run:       balancer(httpBalance),
    })
}
