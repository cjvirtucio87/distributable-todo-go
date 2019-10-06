package actors

type Peer interface {
	AddEntries(entries []Entry) bool
	AddPeer(peer Peer)
	Send(entries []Entry) bool
	Count() int
}

type Entry struct {
	command string
}

type basicPeer struct {
	log   []Entry
	peers []Peer
}

func (p *basicPeer) AddEntries(entries []Entry) bool {
	p.log = append(p.log, entries...)

	return true
}

func (p *basicPeer) AddPeer(otherPeer Peer) {
	p.peers = append(p.peers, otherPeer)
}

func (p *basicPeer) Count() int {
	return len(p.peers)
}

func (p *basicPeer) Send(entries []Entry) bool {
	success := []bool{}

	for _, otherPeer := range p.peers {
		success = append(success, otherPeer.AddEntries(entries))
	}

	return len(success) == (len(entries) * len(p.peers))
}

func NewPeer() Peer {
	p := &basicPeer{
		log:   []Entry{},
		peers: []Peer{},
	}

	return p
}
