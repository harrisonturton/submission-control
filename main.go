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

	/*args := Args{5, 10}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	panicError(err)
	fmt.Println(fmt.Sprintf("%d", reply))*/
}

func panicError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
