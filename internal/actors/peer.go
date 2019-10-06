package actors

type entry struct {
	command string
}

type peer struct {
	log   []entry
	Peers []*peer
}

func (p *peer) AddPeer(otherPeer *peer) {
	p.Peers = append(p.Peers, otherPeer)
}

func NewPeer() *peer {
	p := &peer{
		log:   []entry{},
		Peers: []*peer{},
	}

	return p
}
