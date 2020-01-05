// +build integration

package manager

import (
	"testing"
)

type mockLoader struct {
	peers []HttpPeerConfig
}

func (l *mockLoader) Unmarshal(obj interface{}) error {
	c := obj.(*HttpManagerConfig)

	c.Peers = l.peers

	c.startTimeout = 5

	return nil
}

func (l *mockLoader) Load() error {
	return nil
}

func TestIntegrationStart(t *testing.T) {
	if m, err := NewHttpManager(
		&mockLoader{
			peers: []HttpPeerConfig{
				HttpPeerConfig{
					Host:   "127.0.0.1",
					Port:   8110,
					Scheme: "http",
				},
			},
		},
	); err != nil {
		t.Fatal(err)
	} else {
		m.Start()

		defer m.Stop()
	}
}
