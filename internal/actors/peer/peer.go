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

func NewPeer(kind string) Peer {
	var p Peer

	switch kind {
	default:
		p = &basicPeer{
			log:   []Entry{},
			peers: []Peer{},
		}
		break
	}

	return p
}
