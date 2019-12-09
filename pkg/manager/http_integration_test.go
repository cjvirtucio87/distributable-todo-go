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
