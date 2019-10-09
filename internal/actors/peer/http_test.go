package actors

import (
	"fmt"
	"testing"
)

func TestLogCountHttp(t *testing.T) {
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
		err := follower.Init()

		if err != nil {
			t.Fatal(err)
		}

		leader.AddPeer(follower)
	}

	err := leader.Init()

	if err != nil {
		t.Fatal(err)
	}

	expectedLogCount := 0

	for _, follower := range followers {
		actualLogCount := follower.LogCount()

		if expectedLogCount != actualLogCount {
			t.Fatal(fmt.Printf("expected %d, was %d", expectedLogCount, actualLogCount))
		}
	}
}

func TestPeerCountHttp(t *testing.T) {
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
		t.Fatal(err)
	}

	expectedCount := len(followers)
	actualCount := leader.PeerCount()

	if expectedCount != actualCount {
		t.Fatal(fmt.Printf("expected %d, was %d", expectedCount, actualCount))
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
		t.Fatal(err)
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
		t.Fatal(fmt.Printf("expectedSendResult %t, was %t\n", expectedSendResult, actualSendResult))
	}

	expectedPeerLogCount := 1

	for _, p := range leader.Followers() {
		actualPeerLogCount := p.LogCount()

		if expectedPeerLogCount != actualPeerLogCount {
			t.Fatal(fmt.Printf("expectedPeerLogCount %d, was %d\n", expectedPeerLogCount, actualPeerLogCount))
		}

		id := 0

		actualPeerEntry, ok := p.Entry(id)

		if !ok {
			t.Fatal(fmt.Printf("unable to retrieve entry with id %d\n", id))
		} else if expectedEntry != actualPeerEntry {
			t.Fatal(fmt.Printf("expectedEntry %v, was %v", expectedEntry.Command, actualPeerEntry.Command))
		}
	}
}
