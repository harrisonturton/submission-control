package main

import (
	"flag"
	"fmt"
	"github.com/harrisonturton/submission-control/daemon/daemon"
	"os"
	"sync"
)

var port = flag.String("port", "3000", "The port to run the local RPC server on")
var wg sync.WaitGroup

func main() {
	flag.Parse()
	fmt.Println(*port)
	instance, err := daemon.NewDaemon("1.38", "localhost:"+*port, os.Stdout)
	panicError(err)

	wg.Add(2)
	go instance.ServeRPC(&wg)
	go func() {
		defer wg.Done()
		fmt.Println("Woohoo!")
	}()
	wg.Wait()
	fmt.Println("Finishing...")
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
