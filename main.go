package main

import (
	"flag"
	"fmt"
	"github.com/harrisonturton/hydra-daemon/remote"
	"net/http"
	"net/rpc"
)

var port = flag.String("port", "3000", "The port to run the local RPC server on")

func main() {
	flag.Parse()
	remote, err := remote.NewRemoteServer("1.38")
	if err != nil {
		panic(err)
	}

	rpc.Register(remote)
	rpc.HandleHTTP()
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		fmt.Println(err.Error())
	}
}
