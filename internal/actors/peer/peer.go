package actors

type Peer interface {
	AddEntries(e entryInfo) bool
	AddPeer(peer Peer)
	Count() int
	Id() int
	Send(m Message) bool
}

type entryInfo struct {
	entries   []Entry
	nextIndex int
}

type Entry struct {
	command string
}

type Message struct {
	entries []Entry
}

func NewPeer(kind string, id int) Peer {
	var p Peer

	switch kind {
	default:
		p = &basicPeer{
			id:           id,
			log:          []Entry{},
			nextIndexMap: map[int]int{},
			peers:        []Peer{},
		}
		break
	}

	return p
}
