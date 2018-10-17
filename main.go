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
	instance, err := daemon.NewDaemon("1.38", "localhost:"+*port, os.Stdout)
	panicError(err)

	wg.Add(2)
	go instance.ServeRPC(&wg)
	go func() {
		defer wg.Done()
		fmt.Println("Woohoo!")
		instance.AddEnvironment("hello-world", []string{})
		if err := instance.Run("hello-world"); err != nil {
			fmt.Println("Error running hello-world")
			fmt.Println(err.Error())
		}
	}()
	wg.Wait()
	fmt.Println("Finishing...")
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
