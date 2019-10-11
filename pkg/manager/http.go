package manager

import (
	"cjvirtucio87/distributed-todo-go/internal/actors/peer"
	"context"
	"log"
	"time"
)

type httpManager struct {
	peers []actors.Peer
}

func (m *httpManager) Healthcheck() error {
	for _, peer := range m.peers {
		peer.LogCount()
	}

	return nil
}

func (m *httpManager) Start() {
	log.Printf("Starting peers.\n")
	for _, peer := range m.peers {
		err := peer.Init()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (m *httpManager) Stop() {
	log.Printf("Stopping peers.\n")
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	for _, peer := range m.peers {
		if err := peer.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}
}
