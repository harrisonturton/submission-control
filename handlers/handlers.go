package handlers

import (
	"errors"
	"fmt"
	"net/rpc"
	"strconv"
)

type Command int

type Handler func(argv []string, client *rpc.Client) error

const (
	Stop Command = iota
	Help
	Scale
)

// Convert from the string representation into
// our internal representation
var commandNames = map[string]Command{
	"stop":  Stop,
	"help":  Help,
	"scale": Scale,
}

// Find the handler for a specific command
var commandHandlers = map[Command]Handler{
	Stop:  stopHandler,
	Help:  helpHandler,
	Scale: scaleHandler,
}

// Attempt to run a command
func RunCommand(input []string, client *rpc.Client) error {
	if len(input) <= 1 {
		return errors.New("Usage: hydra COMMAND [OPTIONS]")
	}
	name, ok := commandNames[input[1]]
	if !ok {
		return errors.New(fmt.Sprintf("%s: command not found", input[0]))
	}
	handler, ok := commandHandlers[name]
	if !ok {
		return errors.New(fmt.Sprintf("%s: command handler not found", input[0]))
	}
	return handler(input, client)
}

// Handler for the "stop" command
func stopHandler(argv []string, client *rpc.Client) error {
	_, err := fmt.Println("Exiting!")
	return err
}

// Handler for the "help" command
func helpHandler(argv []string, client *rpc.Client) error {
	fmt.Println("Usage: hydra COMMAND [OPTIONS]")
	fmt.Println("\thelp shows the usage for possible commands")
	fmt.Println("\tstop will stop the hydra daemon")
	return nil
}

// Handler for the "scale" command
func scaleHandler(argv []string, client *rpc.Client) error {
	if len(argv) < 3 {
		fmt.Println(argv)
		return errors.New("Usage: hydra scale SERVICE_ID REPLICAS")
	}
	replicas, err := strconv.Atoi(argv[2])
	if err != nil {
		return err
	}
	fmt.Println(argv[1])
	fmt.Println(argv[2])
	return nil
}
