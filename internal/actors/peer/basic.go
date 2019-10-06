package actors

type basicPeer struct {
	id           int
	log          []Entry
	nextIndexMap map[int]int
	peers        []Peer
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
	success := []bool{}

	for _, otherPeer := range p.peers {
		otherPeerId := otherPeer.Id()

		success = append(
			success,
			otherPeer.AddEntries(
				entryInfo{
					entries:   m.entries,
					nextIndex: p.nextIndexMap[otherPeerId],
				},
			),
		)
	}

	if len(success) != (len(m.entries) * len(p.peers)) {
		return false
	}

	for _, otherPeer := range p.peers {
		otherPeerId := otherPeer.Id()
		p.nextIndexMap[otherPeerId] += len(m.entries)
	}

	return true
}
