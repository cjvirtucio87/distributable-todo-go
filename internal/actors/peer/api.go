package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/dto"
	"cjvirtucio87/distributed-todo-go/internal/rlog"
	"cjvirtucio87/distributed-todo-go/internal/rlogging"
	"context"
)

type Peer interface {
	AddEntries(e dto.EntryInfo) (bool, error)
	AddPeer(peer Peer)
	Entry(idx int) (dto.Entry, error)
	Followers() []Peer
	Init() error
	PeerCount() int
	LogCount() int
	Id() int
	Send(m dto.Message) (bool, error)
	Shutdown(ctx context.Context) error
}

func NewBasicPeer(id int) Peer {
	return &basicPeer{
		id:           id,
		rlog:         rlog.NewBasicLog(),
		NextIndexMap: map[int]int{},
		peers:        []Peer{},
	}
}

func NewHttpPeer(scheme, host string, port, id int) Peer {
	p := &httpPeer{
		basicPeer: basicPeer{
			id:           id,
			rlog:         rlog.NewBasicLog(),
			NextIndexMap: map[int]int{},
			peers:        []Peer{},
		},
		logger: rlogging.NewZapLogger(),
		scheme: scheme,
		host:   host,
		port:   port,
	}

	return p
}
