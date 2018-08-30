package main

import (
	"log"
	"net/http"
	"net/http/httputil"

	BA "github.com/darkhelmet/balance/backends"
	"github.com/gonuts/commander"
)

var (
	httpsOptions = struct {
		certFile, keyFile string
	}{}
)

func httpsBalance(bind string, backends BA.Backends) error {
	if httpsOptions.certFile == "" || httpsOptions.keyFile == "" {
		log.Fatalln("specify both -cert and -key")
	}

	log.Println("using https balancing")

	proxy := &Proxy{
		&httputil.ReverseProxy{Director: func(req *http.Request) {
			backend := backends.Choose()
			if backend == nil {
				log.Printf("no backend for client %s", req.RemoteAddr)
				panic(NoBackend{})
			}
			req.URL.Scheme = "http"
			req.Header.Add("X-Forwarded-Proto", "https")
			req.URL.Host = backend.String()
			req.Header.Add(XRealIP, RealIP(req))
		}},
	}
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
	log.Printf("listening on %s, balancing %d backends", bind, backends.Len())
	return http.ListenAndServeTLS(bind, httpsOptions.certFile, httpsOptions.keyFile, proxy)
}

func init() {
	fs := newFlagSet("https")
	fs.StringVar(&httpsOptions.certFile, "cert", "", "the SSL certificate file to use")
	fs.StringVar(&httpsOptions.keyFile, "key", "", "the SSL key file to use")

	cmd.Subcommands = append(cmd.Subcommands, &commander.Command{
		UsageLine: "https [options] <backend> [<more backends>]",
		Short:     "performs https based load balancing",
		Flag:      *fs,
		Run:       balancer(httpsBalance),
	})
}
