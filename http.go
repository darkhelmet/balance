package main

import (
    "log"
    "net/http"
    "net/http/httputil"

    BA "github.com/darkhelmet/balance/backends"
    "github.com/gonuts/commander"
)

func httpBalance(bind string, backends BA.Backends) error {
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
    return http.ListenAndServe(bind, proxy)
}

func init() {
    fs := newFlagSet("http")

    cmd.Subcommands = append(cmd.Subcommands, &commander.Command{
        UsageLine: "http [options] [<backends>]",
        Short:     "performs http based load balancing",
        Flag:      *fs,
        Run:       balancer(httpBalance),
    })
}
