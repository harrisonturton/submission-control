package main

import (
	"fmt"
	"github.com/harrisonturton/hydra-daemon/container"
	"time"
)

var replicas uint64 = 3
var name = "test"

func main() {
	client, err := container.NewClient("1.38")
	panicErr(err)

	resp, err := client.CreateService(name, "alpine", replicas, []string{"ping", "docker.com"})
	panicErr(err)
	fmt.Println(fmt.Sprintf("ID: %s", resp.ID))

	time.Sleep(time.Second * 10)
	err = client.ScaleService(resp.ID, 7)
	panicErr(err)
	fmt.Println("Increased number of replicas")
}

func panicErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
