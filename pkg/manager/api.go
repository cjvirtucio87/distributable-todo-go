package manager

import (
	"cjvirtucio87/distributed-todo-go/internal/actors/peer"
	"cjvirtucio87/distributed-todo-go/internal/rlogging"
	"cjvirtucio87/distributed-todo-go/pkg/config"
)

type HttpPeerConfig struct {
	Host   string
	Port   int
	Scheme string
}

type HttpManagerConfig struct {
	logger       rlogging.Logger
	Peers        []HttpPeerConfig
	LivenessWaitTime int
	ShutdownTimeout  int
}

type Manager interface {
	Start()
	Stop()
}

func NewHttpManager(loader config.Loader) (Manager, error) {
	if err := loader.Load(); err != nil {
		return nil, err
	}

	var c HttpManagerConfig

	if err := loader.Unmarshal(&c); err != nil {
		return nil, err
	}

	rlogger := rlogging.NewZapLogger()

	peers := []actors.Peer{}
	id := 0

	for _, httpPeerConfig := range c.Peers {
		peers = append(
			peers,
			actors.NewHttpPeer(
				httpPeerConfig.Scheme,
				httpPeerConfig.Host,
				httpPeerConfig.Port,
				id,
				actors.WithLogger(rlogger),
			),
		)

		id++
	}

	leader := peers[0]

	for _, follower := range peers[1:] {
		leader.AddPeer(follower)
	}

	return &httpManager{
		logger:       rlogger,
		peers:        peers,
		LivenessWaitTime: c.LivenessWaitTime,
		ShutdownTimeout:  c.ShutdownTimeout,
	}, nil
}
