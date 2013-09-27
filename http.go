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
    proxy := &Proxy{
        &httputil.ReverseProxy{Director: func(req *http.Request) {
            backend := backends.Choose()
            if backend == nil {
                log.Printf("no backend for client %s", req.RemoteAddr)
                panic(NoBackend{})
            }
            req.URL.Scheme = "http"
            req.URL.Host = backend.String()
            req.Header.Add(XRealIP, RealIP(req))
        }},
    }
    log.Printf("listening on %s, balancing %d backends", bind, backends.Len())
    err := http.ListenAndServe(bind, proxy)
    if err != nil {
        log.Fatalf("failed to bind: %s", err)
    }
}

func init() {
    fs := newFlagSet("http")

    cmd.Commands = append(cmd.Commands, &commander.Command{
        UsageLine: "http [options] [<backends>]",
        Short:     "performs http based load balancing",
        Flag:      *fs,
        Run:       balancer(httpBalance),
    })
}
