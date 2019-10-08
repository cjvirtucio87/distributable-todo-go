package actors

import (
	"fmt"
	"testing"
)

func TestPeerCount(t *testing.T) {
	scheme := "http"
	host := "127.0.0.1"
	leaderPort := "8080"

	leader := NewHttpPeer(
		scheme,
		host,
		leaderPort,
		0,
	)

	follower := NewHttpPeer(
		scheme,
		host,
		"8081",
		1,
	)

	leader.AddPeer(follower)

	err := leader.Init()

	if err != nil {
		t.Error(err)
	}

	expectedCount := 1
	actualCount := leader.PeerCount()

	if expectedCount != actualCount {
		t.Error(fmt.Printf("expected %d, was %d", expectedCount, actualCount))
	}
}
