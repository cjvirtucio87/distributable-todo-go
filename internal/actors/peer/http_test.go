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

func TestSendHttp(t *testing.T) {
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
			&httpPeer{
				scheme: scheme,
				host:   host,
				port:   leaderPort + i,
				basicPeer: basicPeer{
					id: i,
					log: []Entry{
						Entry{
							Command: "not supposed to be here",
						},
						Entry{
							Command: "not supposed to be here either",
						},
					},
					NextIndexMap: map[int]int{},
					peers:        []Peer{},
				},
			},
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

	expectedSendResult := true
	expectedEntry := Entry{Command: "doFoo"}
	actualSendResult := leader.Send(
		Message{
			Entries: []Entry{
				expectedEntry,
			},
		},
	)

	if expectedSendResult != actualSendResult {
		t.Error(fmt.Printf("expectedSendResult %t, was %t\n", expectedSendResult, actualSendResult))
	}

	expectedPeerLogCount := 1

	for _, p := range leader.Followers() {
		actualPeerLogCount := p.LogCount()

		if expectedPeerLogCount != actualPeerLogCount {
			t.Error(fmt.Printf("expectedPeerLogCount %d, was %d\n", expectedPeerLogCount, actualPeerLogCount))
		}

		actualPeerEntry := p.Entry(0)

		if expectedEntry != actualPeerEntry {
			t.Error(fmt.Printf("expectedEntry %v, was %v\n", expectedEntry.Command, actualPeerEntry.Command))
		}
	}
}
