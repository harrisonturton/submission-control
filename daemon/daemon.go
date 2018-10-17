package daemon

import (
	"errors"
	"fmt"
	"github.com/harrisonturton/submission-control/daemon/container"
	"github.com/harrisonturton/submission-control/daemon/environment"
	"github.com/harrisonturton/submission-control/daemon/remote"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"sync"
)

type Daemon struct {
	Addr         string
	Client       *container.Client
	Logger       *log.Logger
	Remote       *remote.RemoteServer
	Environments map[string]*environment.Environment
}

func NewDaemon(version string, addr string, loggerOut io.Writer) (*Daemon, error) {
	client, err := container.NewClient(version)
	if err != nil {
		return nil, err
	}
	logger := log.New(loggerOut, "", 0)
	return &Daemon{
		Addr:         addr,
		Client:       client,
		Logger:       logger,
		Remote:       remote.NewRemoteServer(client, logger),
		Environments: map[string]*environment.Environment{},
	}, nil
}

// Begin serving the RPC server through remote.RemoteServer
func (daemon *Daemon) ServeRPC(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("RPC server running on " + daemon.Addr)
	rpc.Register(daemon.Remote)
	rpc.HandleHTTP()
	if err := http.ListenAndServe(daemon.Addr, nil); err != nil {
		daemon.Logger.Printf(err.Error())
	}
}

// Add an environment to catch future build * test requests
func (daemon *Daemon) AddEnvironment(baseImage string, commands []string) {
	env := environment.NewEnvironment(baseImage, commands, daemon.Client, daemon.Logger)
	daemon.Environments[baseImage] = env
	log.Printf("Added environment with image " + baseImage)
}

// Launch a container for an environment. This is blocking
func (daemon *Daemon) Run(baseImage string) error {
	if env, ok := daemon.Environments[baseImage]; ok {
		if err := env.Run(); err != nil {
			return err
		}
		daemon.Logger.Printf("Running environment " + baseImage)
		return nil
	}
	return errors.New("Environment " + baseImage + " could not be found.")
}
