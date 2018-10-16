package main

import (
	"fmt"
	"github.com/harrisonturton/hydra-cli/cli"
	"sync"
)

var wg sync.WaitGroup
var stop = make(chan bool)

func main() {
	wg.Add(1)
	go cli.Run(stop, &wg)
	wg.Wait()
	fmt.Println("Goodbye!")
}
