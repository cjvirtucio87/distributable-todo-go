package manager

import (
	"cjvirtucio87/distributed-todo-go/internal/actors/peer"
	"cjvirtucio87/distributed-todo-go/pkg/config"
	"log"
)

type HttpPeerConfig struct {
	Host   string
	Port   int
	Scheme string
}

type HttpManagerConfig struct {
	Peers []HttpPeerConfig
}

type Manager interface {
	Healthcheck() error
	Start()
	Stop()
}

func NewHttpManager(loader config.Loader) Manager {
	if err := loader.Load(); err != nil {
		log.Fatal(err)
	}

	var c HttpManagerConfig

	if err := loader.Unmarshal(&c); err != nil {
		log.Fatal(err)
	}

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
			),
		)

		id++
	}

	leader := peers[0]

	for _, follower := range peers[1:] {
		leader.AddPeer(follower)
	}

	return &httpManager{
		peers: peers,
	}
}
