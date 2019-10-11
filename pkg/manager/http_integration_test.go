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

	return nil
}

func (l *mockLoader) Load() error {
	return nil
}

func TestIntegrationStart(t *testing.T) {
	m := NewHttpManager(
		&mockLoader{
			peers: []HttpPeerConfig{
				HttpPeerConfig{
					Host:   "localhost",
					Port:   8080,
					Scheme: "http",
				},
			},
		},
	)

	defer m.Stop()

	m.Start()

	if err := m.Healthcheck(); err != nil {
		t.Fatalf(
			"manager failed to start peers; error: %s\n",
			err.Error(),
		)
	}
}
