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
	Id() int
	Init() error
	PeerCount() (int, error)
	LogCount() (int, error)
	Send(m dto.Message) (bool, error)
	setLog(log rlog.Log) error
	setLogger(rlogger rlogging.Logger) error
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

func NewHttpPeer(scheme, host string, port, id int, options ...func(p Peer)) Peer {
	p := &httpPeer{
		basicPeer: basicPeer{
			id:           id,
			rlog:         rlog.NewBasicLog(),
			NextIndexMap: map[int]int{},
			peers:        []Peer{},
		},
		scheme: scheme,
		host:   host,
		port:   port,
	}

	for _, o := range options {
		o(p)
	}

	return p
}

func WithLog(log rlog.Log) func(p Peer) {
	return func(p Peer) {
		p.setLog(log)
	}
}

func WithLogger(rlogger rlogging.Logger) func(p Peer) {
	return func(p Peer) {
		p.setLogger(rlogger)
	}
}
