package handlers

import (
	"errors"
	"fmt"
)

type Command int

type Handler func(argv []string) error

const (
	Stop Command = iota
	Help
)

// Convert from the string representation into
// our internal representation
var commandNames = map[string]Command{
	"stop": Stop,
	"help": Help,
}

// Find the handler for a specific command
var commandHandlers = map[Command]Handler{
	Stop: stopHandler,
	Help: helpHandler,
}

// Attempt to run a command
func RunCommand(input []string) error {
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
	return handler(input)
}

// Handler for the "stop" command
func stopHandler(argv []string) error {
	_, err := fmt.Println("Exiting!")
	return err
}

// Handler for the "help" command
func helpHandler(argv []string) error {
	fmt.Println("Usage: hydra COMMAND [OPTIONS]")
	fmt.Println("\thelp shows the usage for possible commands")
	fmt.Println("\tstop will stop the hydra daemon")
	return nil
}
