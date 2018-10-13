package main

import (
	"./container"
	"fmt"
)

func main() {
	client, err := container.NewClient("1.38")
	panicError(err)

	// Create base container
	resp, err := client.CreateContainer("python-submission-base")
	panicError(err)

	// Run it
	if err = client.StartContainer(resp.ID); err != nil {
		fmt.Println("Failed to run container " + resp.ID)
		panicError(err)
	}
}

func panicError(err error) {
	fmt.Println(err.Error())
	panic(err)
}
