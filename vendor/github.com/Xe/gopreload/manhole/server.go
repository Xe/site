package main

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/rpc"
)

func init() {
	l, err := net.Listen("tcp", "127.0.0.2:0")
	if err != nil {
		log.Printf("manhole: cannot bind to 127.0.0.2:0: %v", err)
		return
	}

	log.Printf("manhole: Now listening on http://%s", l.Addr())

	rpc.HandleHTTP()
	go http.Serve(l, nil)
}
