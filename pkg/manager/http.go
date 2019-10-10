package dtm

import (
	"cjvirtucio87/distributed-todo-go/internal/actors/peer"
	"log"
)

type httpManager struct {
	peers []actors.Peer
}

func (m *httpManager) Start() {
	for _, peer := range m.peers {
		err := peer.Init()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (m *httpManager) Stop() {
}
