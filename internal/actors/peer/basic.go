package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/dto"
	"cjvirtucio87/distributed-todo-go/internal/rlog"
	"context"
	"errors"
	"fmt"
)

type basicPeer struct {
	id           int
	rlog         rlog.Log
	NextIndexMap map[int]int
	peers        []Peer
}

func (p *basicPeer) AddEntries(e dto.EntryInfo) (bool, error) {
	latestIndex := p.LogCount() + 1

	if e.NextIndex > latestIndex {
		return false, nil
	}

	for _, entry := range e.Entries {
		entry.Id = latestIndex
		latestIndex++
	}

	p.rlog.AddEntries(e)

	return true, nil
}

func (p *basicPeer) AddPeer(otherPeer Peer) {
	p.peers = append(p.peers, otherPeer)
}

func (p *basicPeer) Entry(idx int) (dto.Entry, error) {
	var e dto.Entry
	var err error

	if p.LogCount() <= idx {
		err = errors.New(
			fmt.Sprintf(
				"idx %d exceeds log boundary",
				idx,
			),
		)
	} else {
		e, err = p.rlog.Entry(idx)
	}

	return e, err
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

func (p *basicPeer) Send(m dto.Message) (bool, error) {
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
			return false, errors.New("unable to retrieve entries")
		}

		if successfulAppend, err := otherPeer.AddEntries(
			dto.EntryInfo{
				Entries:   entries,
				NextIndex: nextIndex,
			},
		); err != nil {
			return false, err
		} else if !successfulAppend {
			for nextIndex := p.NextIndexMap[otherPeerId]; nextIndex >= 0; nextIndex-- {
				if successfulAppend, err := otherPeer.AddEntries(
					dto.EntryInfo{
						Entries:   entries,
						NextIndex: nextIndex,
					},
				); err != nil {
					return false, err
				} else if successfulAppend {
					successfulAppendCount++
					break
				}

				p.NextIndexMap[otherPeerId] = nextIndex - 1
			}
		} else {
			successfulAppendCount++
		}
	}

	peerCount := p.PeerCount()

	if successfulAppendCount != peerCount {
		return false, errors.New(
			fmt.Sprintf(
				"only %d out of %d successful append calls",
				successfulAppendCount,
				peerCount,
			),
		)
	}

	for _, otherPeer := range p.peers {
		otherPeerId := otherPeer.Id()
		p.NextIndexMap[otherPeerId] += len(m.Entries)
	}

	return true, nil
}

func (p *basicPeer) Shutdown(ctx context.Context) error {
	p.NextIndexMap = nil

	return nil
}
