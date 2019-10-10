package dtm

type HttpPeerConfig struct {
	Host   string
	Port   string
	Scheme string
}

type HttpManagerConfig struct {
	HttpPeers []HttpPeerConfig
}

type Manager interface {
	Start()
	Stop()
}

func NewHttpManager() {
}
