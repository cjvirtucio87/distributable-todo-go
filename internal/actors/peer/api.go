package actors

import "cjvirtucio87/distributed-todo-go/internal/dto"

type Peer interface {
	AddEntries(e dto.EntryInfo) bool
	AddPeer(peer Peer)
	Entry(idx int) (dto.Entry, bool)
	Followers() []Peer
	Init() error
	PeerCount() int
	LogCount() int
	Id() int
	Send(m dto.Message) bool
}

func NewBasicPeer(id int) Peer {
	return &basicPeer{
		id:           id,
		log:          []dto.Entry{},
		NextIndexMap: map[int]int{},
		peers:        []Peer{},
	}
}

func NewHttpPeer(scheme, host string, port, id int) Peer {
	p := &httpPeer{
		basicPeer: basicPeer{
			id:           id,
			log:          []dto.Entry{},
			NextIndexMap: map[int]int{},
			peers:        []Peer{},
		},
		scheme: scheme,
		host:   host,
		port:   port,
	}

	return p
}
