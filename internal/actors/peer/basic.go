package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/dto"
	"cjvirtucio87/distributed-todo-go/internal/rlog"
	"context"
)

type basicPeer struct {
	id           int
	rlog         rlog.Log
	NextIndexMap map[int]int
	peers        []Peer
}

func (p *basicPeer) AddEntries(e dto.EntryInfo) bool {
	latestIndex := p.LogCount() + 1

	if e.NextIndex > latestIndex {
		return false
	}

	for _, entry := range e.Entries {
		entry.Id = latestIndex
		latestIndex++
	}

	p.rlog.AddEntries(e)

	return true
}

func (p *basicPeer) AddPeer(otherPeer Peer) {
	p.peers = append(p.peers, otherPeer)
}

func (p *basicPeer) Entry(idx int) (dto.Entry, bool) {
	var result dto.Entry
	ok := true

	if p.LogCount() <= idx {
		ok = false
	} else {
		result, ok = p.rlog.Entry(idx)
	}

	return result, ok
}

func (p *basicPeer) Followers() []Peer {
	return p.peers[:]
}

func (p *basicPeer) Init() error {
	for _, otherPeer := range p.peers {
		p.NextIndexMap[otherPeer.Id()] = p.LogCount()
	}

	return nil
}

func (p *basicPeer) LogCount() int {
	return p.rlog.Count()
}

func (p *basicPeer) PeerCount() int {
	return len(p.peers)
}

func (p *basicPeer) Id() int {
	return p.id
}

func (p *basicPeer) Send(m dto.Message) bool {
	p.rlog.AddEntries(
		dto.EntryInfo{
			NextIndex: p.rlog.Count(),
			Entries:   m.Entries,
		},
	)

	var successfulAppendCount int

	for _, otherPeer := range p.peers {
		otherPeerId := otherPeer.Id()
		nextIndex := p.NextIndexMap[otherPeerId]

		entries, ok := p.rlog.Entries(
			nextIndex,
			p.rlog.Count(),
		)

		if !ok {
			return false
		}

		successfulAppend := otherPeer.AddEntries(
			dto.EntryInfo{
				Entries:   entries,
				NextIndex: nextIndex,
			},
		)

		if !successfulAppend {
			for nextIndex := p.NextIndexMap[otherPeerId]; nextIndex >= 0; nextIndex-- {
				successfulAppend = otherPeer.AddEntries(
					dto.EntryInfo{
						Entries:   entries,
						NextIndex: nextIndex,
					},
				)

				if successfulAppend {
					successfulAppendCount++
					break
				}

				p.NextIndexMap[otherPeerId] = nextIndex - 1
			}
		} else {
			successfulAppendCount++
		}
	}

	if successfulAppendCount != p.PeerCount() {
		return false
	}

	for _, otherPeer := range p.peers {
		otherPeerId := otherPeer.Id()
		p.NextIndexMap[otherPeerId] += len(m.Entries)
	}

	return true
}

func (p *basicPeer) Shutdown(ctx context.Context) error {
	p.NextIndexMap = nil

	return nil
}
