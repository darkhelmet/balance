package main

import (
    BA "github.com/darkhelmet/balance/backends"
    "github.com/gonuts/commander"
    "log"
    "net/http"
    "net/http/httputil"
)

var (
    httpsOptions = struct {
        certFile, keyFile string
    }{}
)

func httpsBalance(bind string, backends BA.Backends) {
    if httpsOptions.certFile == "" || httpsOptions.keyFile == "" {
        log.Fatalln("specify both -cert and -key")
    }

    log.Println("using https balancing")
    proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
        req.URL.Scheme = "http"
        req.Header.Add("X-Forwarded-Proto", "https")
        req.URL.Host = backends.Choose()
        req.Header.Add(XRealIP, RealIP(req))
    }}
    log.Printf("listening on %s, balancing %d backends", bind, backends.Len())
    err := http.ListenAndServeTLS(bind, httpsOptions.certFile, httpsOptions.keyFile, proxy)
    if err != nil {
        log.Fatalf("failed to bind: %s", err)
    }
}

func init() {
    fs := newFlagSet("https")
    fs.StringVar(&httpsOptions.certFile, "cert", "", "the SSL certificate file to use")
    fs.StringVar(&httpsOptions.keyFile, "key", "", "the SSL key file to use")

    cmd.Commands = append(cmd.Commands, &commander.Command{
        UsageLine: "https [options] <backend> [<more backends>]",
        Short:     "performs https based load balancing",
        Flag:      *fs,
        Run:       balancer(httpsBalance),
    })
}
