package actors

type basicPeer struct {
	id           int
	log          []Entry
	nextIndexMap map[int]int
	peers        []Peer
}

func (p *basicPeer) AddEntries(e entryInfo) bool {
	if e.nextIndex > len(p.log)+1 {
		return false
	}

	p.log = append(p.log[:e.nextIndex], e.entries...)

	return true
}

func (p *basicPeer) AddPeer(otherPeer Peer) {
	p.peers = append(p.peers, otherPeer)

	p.nextIndexMap[otherPeer.Id()] = len(p.log)
}

func (p *basicPeer) Entry(idx int) Entry {
	return p.log[idx]
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
	p.log = append(p.log, m.entries...)

	var successfulAppendCount int

	for _, otherPeer := range p.peers {
		otherPeerId := otherPeer.Id()
		nextIndex := p.nextIndexMap[otherPeerId]

		successfulAppend := otherPeer.AddEntries(
			entryInfo{
				entries:   p.log[nextIndex:],
				nextIndex: nextIndex,
			},
		)

		if !successfulAppend {
			for nextIndex := p.nextIndexMap[otherPeerId]; nextIndex >= 0; nextIndex-- {
				successfulAppend = otherPeer.AddEntries(
					entryInfo{
						entries:   p.log[nextIndex:],
						nextIndex: nextIndex,
					},
				)

				if successfulAppend {
					successfulAppendCount++
					break
				}

				p.nextIndexMap[otherPeerId] = nextIndex - 1
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
		p.nextIndexMap[otherPeerId] += len(m.entries)
	}

	return true
}
