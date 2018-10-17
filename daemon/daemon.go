package daemon

import (
	"github.com/harrisonturton/submission-control/daemon/container"
	"github.com/harrisonturton/submission-control/daemon/remote"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"sync"
)

type Daemon struct {
	Addr   string
	Client *container.Client
	Logger *log.Logger
	Remote *remote.RemoteServer
}

func NewDaemon(version string, addr string, loggerOut io.Writer) (*Daemon, error) {
	client, err := container.NewClient(version)
	if err != nil {
		return nil, err
	}
	logger := log.New(loggerOut, "", 0)
	return &Daemon{
		Addr:   addr,
		Client: client,
		Logger: logger,
		Remote: remote.NewRemoteServer(client, logger),
	}, nil
}

// Begin serving the RPC server through remote.RemoteServer
func (daemon *Daemon) ServeRPC(wg *sync.WaitGroup) {
	defer wg.Done()
	rpc.Register(daemon.Remote)
	rpc.HandleHTTP()
	if err := http.ListenAndServe(daemon.Addr, nil); err != nil {
		daemon.Logger.Printf(err.Error())
	}
}
