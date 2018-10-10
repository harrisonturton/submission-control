package main

import (
	"./config"
	"fmt"
)

const (
	config_path = "config.yaml"
)

func main() {
	fmt.Println("Starting...")
	fmt.Println("Reading config...")
	config, err := config.ReadConfig(config_path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(config)
}
