package handlers

import (
	"errors"
	"fmt"
	"github.com/harrisonturton/hydra-cli/types"
	"net/rpc"
	"strconv"
)

type Command int

type Handler func(argv []string, client *rpc.Client) error

const (
	Stop Command = iota
	Help
	Create
	Remove
	Scale
)

// Convert from the string representation into
// our internal representation
var commandNames = map[string]Command{
	"stop":   Stop,
	"help":   Help,
	"create": Create,
	"remove": Remove,
	"scale":  Scale,
}

// Find the handler for a specific command
var commandHandlers = map[Command]Handler{
	Stop:   stopHandler,
	Help:   helpHandler,
	Create: createHandler,
	Remove: removeHandler,
	Scale:  scaleHandler,
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

// Send an RPC call to creates a new service
func createHandler(argv []string, client *rpc.Client) error {
	if len(argv) < 4 {
		return errors.New("Usage: hydra create BASE_IMAGE REPLICAS [COMMANDS]")
	}
	replicas, err := strconv.ParseUint(argv[3], 10, 64)
	if err != nil {
		return err
	}
	args := types.ServiceCreateSpec{
		BaseImage: argv[2],
		Replicas:  replicas,
		Commands:  argv[1:],
	}
	var reply string
	err = client.Call("RemoteServer.CreateService", args, &reply)
	if err != nil {
		return err
	}
	fmt.Println("Created service: " + reply)
	return nil
}

// Send an RPC call to remove an existing service
func removeHandler(argv []string, client *rpc.Client) error {
	if len(argv) < 3 {
		return errors.New("USAGE: hydra remove SERVICE_ID")
	}
	args := types.ServiceRemoveSpec{
		ServiceID: argv[2],
	}
	var reply int
	return client.Call("RemoteServer.RemoveService", args, &reply)
}

// Send an RPC call to scale a service (i.e. increase or decrease the number of replicas)
func scaleHandler(argv []string, client *rpc.Client) error {
	if len(argv) < 4 {
		return errors.New("Usage: hydra scale SERVICE_ID REPLICAS")
	}
	replicas, err := strconv.ParseUint(argv[3], 10, 64)
	if err != nil {
		return err
	}
	args := types.ServiceScaleSpec{
		ServiceID: argv[2],
		Replicas:  replicas,
	}
	var reply int
	return client.Call("RemoteServer.ScaleService", args, &reply)
}
