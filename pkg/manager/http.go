package manager

import (
	"cjvirtucio87/distributed-todo-go/internal/actors/peer"
	"cjvirtucio87/distributed-todo-go/internal/rlogging"
	"context"
	"reflect"
	"time"
)

type httpManager struct {
	logger       rlogging.Logger
	peers        []actors.Peer
	StartTimeout int
	StopTimeout  int
}

func (m *httpManager) Start() {
	m.logger.Infof("Starting peers.\n")

	channels := []chan error{}

	for _, peer := range m.peers {
		peerChannel := make(chan error)

		go func(peerChannel chan error, peer actors.Peer) {
			peerChannel <- peer.Init()
		}(peerChannel, peer)

		channels = append(channels, peerChannel)
	}

	timeoutChannel := make(chan error)

	go func() {
		m.logger.Infof("waiting..")

		for i := 0; i < m.StartTimeout; i++ {
			m.logger.Infof("%d", i+1)

			time.Sleep(1 * time.Second)
		}

		m.logger.Infof("done waiting. no errors")

		timeoutChannel <- nil
	}()

	channels = append(
		channels,
		timeoutChannel,
	)

	selectCases := make([]reflect.SelectCase, len(channels))

	for i, ch := range channels {
		selectCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}

	_, value, _ := reflect.Select(selectCases)

	if valueInterface := value.Interface(); valueInterface != nil {
		m.logger.Errorf(valueInterface.(error).Error())
	}
}

func (m *httpManager) Stop() {
	m.logger.Infof("Stopping peers.\n")
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(m.StopTimeout)*time.Second,
	)

	defer cancel()

	for _, peer := range m.peers {
		if err := peer.Shutdown(ctx); err != nil {
			m.logger.Errorf(err.Error())
		}
	}
}
