package actors

type Peer interface {
	AddEntries(e entryInfo) bool
	AddPeer(peer Peer)
	Entry(idx int) Entry
	Followers() []Peer
	Init() error
	PeerCount() int
	LogCount() int
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

func NewBasicPeer(id int) Peer {
	return &basicPeer{
		id:           id,
		log:          []Entry{},
		nextIndexMap: map[int]int{},
		peers:        []Peer{},
	}
}

func NewHttpPeer(scheme, host, port string, id int) Peer {
	p := &httpPeer{
		basicPeer: basicPeer{
			id:           id,
			log:          []Entry{},
			nextIndexMap: map[int]int{},
			peers:        []Peer{},
		},
		scheme: scheme,
		host:   host,
		port:   port,
	}

	return p
}
