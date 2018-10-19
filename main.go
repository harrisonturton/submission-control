package main

import (
	"flag"
	"fmt"
	"github.com/harrisonturton/submission-control/daemon/server"
	"os"
	"sync"
	"time"
)

var port = flag.String("port", "3000", "The port to run the local RPC server on")
var wg sync.WaitGroup
var done chan bool

func main() {
	flag.Parse()
	instance, err := server.New(
		"1.38", "localhost:"+*port, []string{"hello-world"}, os.Stdout)
	panicError(err)

	done := make(chan bool)
	wg.Add(2)
	go instance.Serve(done, &wg)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 2)
		for i, _ := range []int{0, 0, 0, 0, 0, 0} {
			fmt.Printf("\rStopping in %d", 5-i)
			time.Sleep(time.Second)
		}
		fmt.Printf("\n")
		close(done)
	}()
	wg.Wait()
	fmt.Println("Finished.")
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
