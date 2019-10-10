package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/dto"
	"cjvirtucio87/distributed-todo-go/internal/log"
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
			t.Fatalf(err.Error())
		}

		leader.AddPeer(follower)
	}

	err := leader.Init()

	if err != nil {
		t.Fatalf(err.Error())
	}

	expectedLogCount := 0

	for _, follower := range followers {
		actualLogCount := follower.LogCount()

		if expectedLogCount != actualLogCount {
			t.Fatalf("expected %d, was %d", expectedLogCount, actualLogCount)
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
		t.Fatalf(err.Error())
	}

	expectedCount := len(followers)
	actualCount := leader.PeerCount()

	if expectedCount != actualCount {
		t.Fatalf("expected %d, was %d", expectedCount, actualCount)
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
					log: log.NewBasicLog(
						log.WithBackend(
							[]dto.Entry{
								dto.Entry{
									Command: "not supposed to be here",
								},
								dto.Entry{
									Command: "not supposed to be here either",
								},
							},
						),
					),
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
		t.Fatalf(err.Error())
	}

	expectedSendResult := true
	expectedEntries := []dto.Entry{
		dto.Entry{
			Command: "doHip",
		},
		dto.Entry{
			Command: "doBar",
		},
		dto.Entry{
			Command: "doFoo",
		},
	}
	actualSendResult := leader.Send(
		dto.Message{
			Entries: expectedEntries,
		},
	)

	if expectedSendResult != actualSendResult {
		t.Fatalf("expectedSendResult %t, was %t\n", expectedSendResult, actualSendResult)
	}

	expectedPeerLogCount := len(expectedEntries)

	for _, p := range leader.Followers() {
		actualPeerLogCount := p.LogCount()

		if expectedPeerLogCount != actualPeerLogCount {
			t.Fatalf("expectedPeerLogCount %d, was %d\n", expectedPeerLogCount, actualPeerLogCount)
		}

		id := expectedPeerLogCount - 1

		expectedLatestEntry := expectedEntries[id]
		actualPeerEntry, ok := p.Entry(id)

		if !ok {
			t.Fatalf("unable to retrieve entry with id %d\n", id)
		} else if expectedLatestEntry != actualPeerEntry {
			t.Fatalf("expectedLatestEntry %v, was %v", expectedLatestEntry.Command, actualPeerEntry.Command)
		} else if expectedLatestEntry.Id != actualPeerEntry.Id {
			t.Fatalf("expectedLatestEntry %v, was %v", expectedLatestEntry.Id, actualPeerEntry.Id)
		}
	}
}
