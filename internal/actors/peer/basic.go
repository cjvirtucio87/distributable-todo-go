package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/dto"
	"cjvirtucio87/distributed-todo-go/internal/rlog"
	"cjvirtucio87/distributed-todo-go/internal/rlogging"
	"context"
	"errors"
	"fmt"
)

type basicPeer struct {
	id           int
	rlog         rlog.Log
	rlogger      rlogging.Logger
	NextIndexMap map[int]int
	peers        []Peer
}

func (p *basicPeer) AddEntries(e dto.EntryInfo) (bool, error) {
	if latestIndex, err := p.LogCount(); err != nil {
		return false, err
	} else if e.NextIndex > latestIndex+1 {
		return false, nil
	} else {
		for _, entry := range e.Entries {
			entry.Id = latestIndex + 1
			latestIndex++
		}

		p.rlog.AddEntries(e)

		return true, nil
	}
}

func (p *basicPeer) AddPeer(otherPeer Peer) {
	p.peers = append(p.peers, otherPeer)
}

func (p *basicPeer) Entry(idx int) (dto.Entry, error) {
	var e dto.Entry
	if result, err := p.LogCount(); err != nil {
		return e, err
	} else if result <= idx {
		return e, errors.New(
			fmt.Sprintf(
				"idx %d exceeds log boundary",
				idx,
			),
		)
	} else {
		return p.rlog.Entry(idx)
	}
}

func (p *basicPeer) Followers() []Peer {
	return p.peers[:]
}

func (p *basicPeer) Init() error {
	for _, otherPeer := range p.peers {
		if result, err := p.LogCount(); err != nil {
			return err
		} else {
			p.NextIndexMap[otherPeer.Id()] = result
		}
	}

	return nil
}

func (p *basicPeer) LogCount() (int, error) {
	return p.rlog.Count(), nil
}

func (p *basicPeer) PeerCount() (int, error) {
	return len(p.peers), nil
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

	if peerCount, err := p.PeerCount(); err != nil {
		return false, err
	} else if successfulAppendCount != peerCount {
		return false, errors.New(
			fmt.Sprintf(
				"only %d out of %d successful append calls",
				successfulAppendCount,
				peerCount,
			),
		)
	} else {
		for _, otherPeer := range p.peers {
			otherPeerId := otherPeer.Id()
			p.NextIndexMap[otherPeerId] += len(m.Entries)
		}

		return true, nil
	}
}

func (p *basicPeer) setLog(log rlog.Log) error {
	p.rlog = log

	return nil
}

func (p *basicPeer) setLogger(rlogger rlogging.Logger) error {
	p.rlogger = rlogger

	return nil
}

func (p *basicPeer) Shutdown(ctx context.Context) error {
	p.NextIndexMap = nil

	return nil
}
