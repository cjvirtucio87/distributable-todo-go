package actors

type Peer interface {
	AddEntries(e EntryInfo) bool
	AddPeer(peer Peer)
	Entry(idx int) (Entry, bool)
	Followers() []Peer
	Init() error
	PeerCount() int
	LogCount() int
	Id() int
	Send(m Message) bool
}

type EntryInfo struct {
	Entries   []Entry
	NextIndex int
}

type Entry struct {
	Id      int
	Command string
}

type Message struct {
	Entries []Entry
}

func NewBasicPeer(id int) Peer {
	return &basicPeer{
		id:           id,
		log:          []Entry{},
		NextIndexMap: map[int]int{},
		peers:        []Peer{},
	}
}

func NewHttpPeer(scheme, host string, port, id int) Peer {
	p := &httpPeer{
		basicPeer: basicPeer{
			id:           id,
			log:          []Entry{},
			NextIndexMap: map[int]int{},
			peers:        []Peer{},
		},
		scheme: scheme,
		host:   host,
		port:   port,
	}

	return p
}
