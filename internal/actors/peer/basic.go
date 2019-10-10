package actors

import "cjvirtucio87/distributed-todo-go/internal/dto"

type basicPeer struct {
	id           int
	log          []dto.Entry
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

	p.log = append(p.log[:e.NextIndex], e.Entries...)

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
		result = p.log[idx]

		if &result == nil {
			ok = false
		}
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
	return len(p.log)
}

func (p *basicPeer) PeerCount() int {
	return len(p.peers)
}

func (p *basicPeer) Id() int {
	return p.id
}

func (p *basicPeer) Send(m dto.Message) bool {
	p.log = append(p.log, m.Entries...)

	var successfulAppendCount int

	for _, otherPeer := range p.peers {
		otherPeerId := otherPeer.Id()
		nextIndex := p.NextIndexMap[otherPeerId]

		successfulAppend := otherPeer.AddEntries(
			dto.EntryInfo{
				Entries:   p.log[nextIndex:],
				NextIndex: nextIndex,
			},
		)

		if !successfulAppend {
			for nextIndex := p.NextIndexMap[otherPeerId]; nextIndex >= 0; nextIndex-- {
				successfulAppend = otherPeer.AddEntries(
					dto.EntryInfo{
						Entries:   p.log[nextIndex:],
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
