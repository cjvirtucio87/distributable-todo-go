package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/rlog"
	"fmt"
)

type basicPeer struct {
	id           int
	rlog         rlog.Log
	NextIndexMap map[int]int
	peers        []Peer
}

// Entries to be appended to the log.
type EntryCollection struct {
	Entries   []rlog.Entry
	// Starting point for the entries. Pre-existing
	// entries in the follower's log beginning from
	// this index will be discarded in favor of the new
	// entries from the leader.
	NextIndex int
}


func (p *basicPeer) AddEntries(entries EntryCollection) error {
	count := p.LogCount()

	if entries.NextIndex > count {
		return fmt.Errorf(
			"NextIndex [%d] exceeds count [%d]",
			entries.NextIndex,
			count,
		)
	}

	id := entries.NextIndex
	for _, entry := range entries.Entries {
		entry.Id = id
		id++
	}

	return p.rlog.AddEntries(entries.NextIndex, entries.Entries)
}

func (p *basicPeer) AddPeer(otherPeer Peer) {
	p.peers = append(p.peers, otherPeer)
}

func (p *basicPeer) Entry(idx int) *rlog.Entry {
	return p.rlog.Entry(idx)
}

func (p *basicPeer) Followers() []Peer {
	return p.peers[:]
}

func (p *basicPeer) Id() int {
	return p.id
}

func (p *basicPeer) Init() error {
	for _, otherPeer := range p.peers {
		result := p.LogCount()
		p.NextIndexMap[otherPeer.Id()] = result
	}

	return nil
}

func (p *basicPeer) LogCount() int {
	return p.rlog.Count()
}

func (p *basicPeer) PeerCount() int {
	return len(p.peers)
}

func (p *basicPeer) Send(m Message) error {
	err := p.rlog.AddEntries(p.rlog.Count(), m.Entries)

	if err != nil {
		return err
	}

	var successfulAppendCount int
	for _, otherPeer := range p.peers {
		otherPeerId := otherPeer.Id()
		// Keep trying to AddEntries() until it succeeds
		for nextIndex := p.NextIndexMap[otherPeerId]; nextIndex >= 0; nextIndex-- {
			entries := p.rlog.Entries(
				nextIndex,
				p.LogCount(),
			)

			p.NextIndexMap[otherPeerId] = nextIndex

			err := otherPeer.AddEntries(
				EntryCollection {
					Entries:   entries,
					NextIndex: nextIndex,
				},
			)

			if err == nil {
				successfulAppendCount++
				break
			}
		}
	}

	peerCount := p.PeerCount()
	if successfulAppendCount != peerCount {
		return fmt.Errorf(
			"only %d out of %d successful append calls",
			successfulAppendCount,
			peerCount,
		)
	}

	for _, otherPeer := range p.peers {
		otherPeerId := otherPeer.Id()
		p.NextIndexMap[otherPeerId] += len(m.Entries)
	}

	return nil
}
