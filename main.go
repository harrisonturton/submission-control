package main

import (
	"fmt"
	"github.com/harrisonturton/hydra-daemon/server"
	"sync"
	"time"
)

var stop = make(chan bool)
var wg sync.WaitGroup

// Run server for 10 seconds before stopping (test graceful shutdown)
func main() {
	wg.Add(1)
	go server.NewServer("localhost:3000").Serve(stop, &wg)
	time.Sleep(time.Second * 10)
	close(stop)
	wg.Wait()
	fmt.Println("Finished.")
}
