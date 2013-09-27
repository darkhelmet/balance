package main

import (
    BA "github.com/darkhelmet/balance/backends"
    "github.com/gonuts/commander"
    "io"
    "log"
    "net"
)

func copy(wc io.WriteCloser, r io.Reader) {
    defer wc.Close()
    io.Copy(wc, r)
}

func handleConnection(us net.Conn, backend BA.Backend) {
    if backend == nil {
        log.Printf("no backend available for connection from %s", us.RemoteAddr())
        us.Close()
        return
    }

    ds, err := net.Dial("tcp", backend.String())
    if err != nil {
        us.Close()
        log.Printf("failed to dial %s: %s", backend, err)
        return
    }

    go copy(ds, us)
    go copy(us, ds)
}

func tcpBalance(bind string, backends BA.Backends) {
    log.Println("using tcp balancing")
    ln, err := net.Listen("tcp", bind)
    if err != nil {
        log.Fatalf("failed to bind: %s", err)
    }

    log.Printf("listening on %s, balancing %d backends", bind, backends.Len())

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Printf("failed to accept: %s", err)
            continue
        }
        go handleConnection(conn, backends.Choose())
    }
}

func init() {
    fs := newFlagSet("tcp")

    cmd.Commands = append(cmd.Commands, &commander.Command{
        UsageLine: "tcp [options] <backend> [<more backends>]",
        Short:     "performs tcp based load balancing",
        Flag:      *fs,
        Run:       balancer(tcpBalance),
    })
}
