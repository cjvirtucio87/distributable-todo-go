package actors

import (
	"fmt"
	"testing"
)

func TestLogCount(t *testing.T) {
	scheme := "http"
	host := "127.0.0.1"
	leaderPort := 8080

	leader := NewHttpPeer(
		scheme,
		host,
		leaderPort,
		0,
	)

	followers := []Peer{}

	for i := 1; i < 3; i++ {
		followers = append(
			followers,
			NewHttpPeer(
				scheme,
				host,
				leaderPort+i,
				i,
			),
		)
	}

	for _, follower := range followers {
		follower.Init()

		leader.AddPeer(follower)
	}

	err := leader.Init()

	if err != nil {
		t.Error(err)
	}

	expectedLogCount := 0

	for _, follower := range followers {
		actualLogCount := follower.LogCount()

		if expectedLogCount != actualLogCount {
			t.Error(fmt.Printf("expected %d, was %d", expectedLogCount, actualLogCount))
		}
	}
}

func TestPeerCount(t *testing.T) {
	scheme := "http"
	host := "127.0.0.1"
	leaderPort := 8080

	leader := NewHttpPeer(
		scheme,
		host,
		leaderPort,
		0,
	)

	followers := []Peer{}

	for i := 1; i < 3; i++ {
		followers = append(
			followers,
			NewHttpPeer(
				scheme,
				host,
				leaderPort+i,
				i,
			),
		)
	}

	for _, follower := range followers {
		leader.AddPeer(follower)
	}

	err := leader.Init()

	if err != nil {
		t.Error(err)
	}

	expectedCount := len(followers)
	actualCount := leader.PeerCount()

	if expectedCount != actualCount {
		t.Error(fmt.Printf("expected %d, was %d", expectedCount, actualCount))
	}
}
