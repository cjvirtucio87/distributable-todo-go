package actors

type basicPeer struct {
	id           int
	log          []Entry
	NextIndexMap map[int]int
	peers        []Peer
}

func (p *basicPeer) AddEntries(e EntryInfo) bool {
	if e.NextIndex > len(p.log)+1 {
		return false
	}

	p.log = append(p.log[:e.NextIndex], e.Entries...)

	return true
}

func (p *basicPeer) AddPeer(otherPeer Peer) {
	p.peers = append(p.peers, otherPeer)

	p.NextIndexMap[otherPeer.Id()] = len(p.log)
}

func (p *basicPeer) Entry(idx int) (Entry, bool) {
	var result Entry
	ok := true

	if len(p.log) <= idx {
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

func (p *basicPeer) Send(m Message) bool {
	p.log = append(p.log, m.Entries...)

	var successfulAppendCount int

	for _, otherPeer := range p.peers {
		otherPeerId := otherPeer.Id()
		nextIndex := p.NextIndexMap[otherPeerId]

		successfulAppend := otherPeer.AddEntries(
			EntryInfo{
				Entries:   p.log[nextIndex:],
				NextIndex: nextIndex,
			},
		)

		if !successfulAppend {
			for nextIndex := p.NextIndexMap[otherPeerId]; nextIndex >= 0; nextIndex-- {
				successfulAppend = otherPeer.AddEntries(
					EntryInfo{
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

	if successfulAppendCount != len(p.peers) {
		return false
	}

	for _, otherPeer := range p.peers {
		otherPeerId := otherPeer.Id()
		p.NextIndexMap[otherPeerId] += len(m.Entries)
	}

	return true
}
