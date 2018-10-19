package balancer

import (
	"errors"
	"fmt"
	"github.com/harrisonturton/submission-control/daemon/client"
	"github.com/harrisonturton/submission-control/daemon/environment"
	"github.com/harrisonturton/submission-control/daemon/remote"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"sync"
)

type Balancer struct {
	Addr         string
	Client       *client.Client
	Logger       *log.Logger
	Remote       *remote.RemoteServer
	Environments map[string]*environment.Environment
}

func NewBalancer(version string, addr string, loggerOut io.Writer) (*Balancer, error) {
	client, err := client.NewClient(version)
	if err != nil {
		return nil, err
	}
	logger := log.New(loggerOut, "", 0)
	return &Balancer{
		Addr:         addr,
		Client:       client,
		Logger:       logger,
		Remote:       remote.NewRemoteServer(client, logger),
		Environments: map[string]*environment.Environment{},
	}, nil
}

// Begin serving the RPC server through remote.RemoteServer
func (balancer *Balancer) ServeRPC(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("RPC server running on " + balancer.Addr)
	rpc.Register(balancer.Remote)
	rpc.HandleHTTP()
	if err := http.ListenAndServe(balancer.Addr, nil); err != nil {
		balancer.Logger.Printf(err.Error())
	}
}

// Add an environment to catch future build * test requests
func (balancer *Balancer) AddEnvironment(baseImage string, commands []string) {
	env := environment.NewEnvironment(baseImage, commands, balancer.Client, balancer.Logger)
	balancer.Environments[baseImage] = env
	balancer.Logger.Printf("Added environment with image " + baseImage)
}

// Launch a container for an environment. This is blocking
func (balancer *Balancer) Run(baseImage string) error {
	if env, ok := balancer.Environments[baseImage]; ok {
		balancer.Logger.Printf("Running environment " + baseImage)
		logs, err := env.RunWithLogs(true, true)
		if err != nil {
			return err
		}
		balancer.Logger.Printf("Successfully ran container with image %s.", baseImage)
		if logs != "" {
			balancer.Logger.Printf("Logs: %s", logs)
		} else {
			balancer.Logger.Printf("No logs.")
		}
		return nil
	}
	return errors.New("Environment " + baseImage + " could not be found.")
}
