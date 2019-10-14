package manager

import (
	"cjvirtucio87/distributed-todo-go/internal/actors/peer"
	"cjvirtucio87/distributed-todo-go/internal/rlogging"
	"context"
	"time"
)

type httpManager struct {
	logger rlogging.Logger
	peers  []actors.Peer
}

func (m *httpManager) Healthcheck() error {
	for _, peer := range m.peers {
		peer.LogCount()
	}

	return nil
}

func (m *httpManager) Start() {
	m.logger.Infof("Starting peers.\n")
	for _, peer := range m.peers {
		if err := peer.Init(); err != nil {
			m.logger.Errorf(err.Error())
		}
	}
}

func (m *httpManager) Stop() {
	m.logger.Infof("Stopping peers.\n")
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	for _, peer := range m.peers {
		if err := peer.Shutdown(ctx); err != nil {
			m.logger.Errorf(err.Error())
		}
	}
}
