package main

import (
	"fmt"
	"github.com/harrisonturton/hydra-cli/handlers"
	"os"
)

func main() {
	err := handlers.RunCommand(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}
