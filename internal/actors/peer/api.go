package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/rlog"
)

type Peer interface {
	// Add entries attempts to add entries from the leader
	// to its log.
	AddEntries(entries EntryCollection) error
	AddPeer(peer Peer)
	Entry(idx int) *rlog.Entry
	Followers() []Peer
	Id() int
	Init() error
	LogCount() int
	PeerCount() int
	// Send a message to follower peers, adding entries
	// to own log and sending them to followers for them
	// to add to their logs. Entries in follower logs
	// that aren't in the leader's log are discarded.
	// Each follower will only get the entries they
	// lack, based on the 'next index' tracked by
	// the leader. Next index is tracked to ensure
	// that no gaps are left in each follower's log.
	Send(m Message) error
}

type Message struct {
	// The entries to be added to peer logs.
	Entries []*rlog.Entry
}

func NewBasicPeer(id int) Peer {
	return &basicPeer{
		id:           id,
		rlog:         rlog.NewBasicLog(),
		NextIndexMap: make(map[int]int),
		peers:        []Peer{},
	}
}
