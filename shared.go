package main

import (
    "net"
    "net/http"
    "net/http/httputil"
)

const (
    colon   = ":"
    XRealIP = "X-Real-IP"
)

type NoBackend struct{}

func RealIP(req *http.Request) string {
    host, _, _ := net.SplitHostPort(req.RemoteAddr)
    return host
}

type Proxy struct {
    *httputil.ReverseProxy
}

func (p *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    defer func() {
        if err := recover(); err != nil {
            switch err.(type) {
            case NoBackend:
                rw.WriteHeader(503)
                req.Body.Close()
            default:
                panic(err)
            }
        }
    }()
    p.ReverseProxy.ServeHTTP(rw, req)
}
