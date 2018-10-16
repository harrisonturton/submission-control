package main

import (
	"fmt"
	"github.com/harrisonturton/hydra-cli/handlers"
	"net/rpc"
	"os"
)

func main() {
	client, err := rpc.DialHTTP("tcp", ":3000")
	panicError(err)
	err = handlers.RunCommand(os.Args, client)
	panicError(err)
}

func panicError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
