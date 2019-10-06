package actors

type basicPeer struct {
	id           int
	log          []Entry
	nextIndexMap map[int]int
	peers        []Peer
	timeout      int
}

func (p *basicPeer) AddEntries(e entryInfo) bool {
	if e.nextIndex != len(p.log) {
		return false
	}

	p.log = append(p.log, e.entries...)

	return true
}

func (p *basicPeer) AddPeer(otherPeer Peer) {
	p.peers = append(p.peers, otherPeer)

	p.nextIndexMap[otherPeer.Id()] = len(p.log)
}

func (p *basicPeer) Count() int {
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
			for i := 0; i < p.timeout; i++ {
				nextIndex = p.nextIndexMap[otherPeerId]

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
