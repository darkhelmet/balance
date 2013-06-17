package main

import (
    "flag"
    BA "github.com/darkhelmet/balance/backends"
    // "github.com/hawx/hadfield"
    "io"
    "log"
    "net"
    "net/http"
    "net/http/httputil"
)

var (
    mode     = flag.String("mode", "tcp", "The mode to balance on: tcp|http")
    bind     = flag.String("bind", "", "The address to bind on")
    balance  []string
    backends BA.Backends
)

func init() {
    flag.Parse()

    if *bind == "" {
        log.Fatalln("specify the address to listen on with -bind")
    }

    servers := flag.Args()
    if len(servers) == 0 {
        log.Fatalln("please specify backend servers")
    }
    backends = BA.NewSimpleBackends(servers)
}

func copy(wc io.WriteCloser, r io.Reader) {
    defer wc.Close()
    io.Copy(wc, r)
}

func handleConnection(us net.Conn, backend string) {
    ds, err := net.Dial("tcp", backend)
    if err != nil {
        us.Close()
        log.Printf("failed to dial %s: %s", backend, err)
        return
    }

    go copy(ds, us)
    go copy(us, ds)
}

func tcpBalance() {
    log.Println("using tcp balancing")
    ln, err := net.Listen("tcp", *bind)
    if err != nil {
        log.Fatalf("failed to bind: %s", err)
    }

    log.Printf("listening on %s, balancing %d backends", *bind, backends.Len())

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Printf("failed to accept: %s", err)
            continue
        }
        go handleConnection(conn, backends.Choose())
    }
}

func httpBalance() {
    log.Println("using http balancing")
    proxy := &httputil.ReverseProxy{Director: func(req *http.Request) {
        req.URL.Scheme = "http"
        req.URL.Host = backends.Choose()
    }}
    log.Printf("listening on %s, balancing %d backends", *bind, backends.Len())
    err := http.ListenAndServe(*bind, proxy)
    if err != nil {
        log.Fatalf("failed to bind: %s", err)
    }
}

func main() {
    switch *mode {
    case "tcp":
        tcpBalance()
    case "http":
        httpBalance()
    default:
        log.Printf("invalid balance mode %s", *mode)
    }
}
