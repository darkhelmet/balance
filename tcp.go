package main

import (
	"fmt"
	"io"
	"log"
	"net"

	BA "github.com/darkhelmet/balance/backends"
	"github.com/gonuts/commander"
)

func copy(wc io.WriteCloser, r io.Reader) {
	defer wc.Close()
	io.Copy(wc, r)
}

func handleConnection(us net.Conn, backend BA.Backend) error {
	if backend == nil {
		log.Printf("no backend available for connection from %s", us.RemoteAddr())
		return us.Close()
	}

	ds, err := net.Dial("tcp", backend.String())
	if err != nil {
		us.Close()
		return fmt.Errorf("failed to dial %s: %s", backend, err)
	}

	// FIXME: fetch errors from copies ?
	go copy(ds, us)
	go copy(us, ds)

	return nil
}

func tcpBalance(bind string, backends BA.Backends) error {
	log.Println("using tcp balancing")
	ln, err := net.Listen("tcp", bind)
	if err != nil {
		return fmt.Errorf("failed to bind: %s", err)
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

	return err
}

func init() {
	fs := newFlagSet("tcp")

	cmd.Subcommands = append(cmd.Subcommands, &commander.Command{
		UsageLine: "tcp [options] <backend> [<more backends>]",
		Short:     "performs tcp based load balancing",
		Flag:      *fs,
		Run:       balancer(tcpBalance),
	})
}
