// +build integration

package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/dto"
	"cjvirtucio87/distributed-todo-go/internal/rlog"
	"cjvirtucio87/distributed-todo-go/internal/rlogging"
	"testing"
)

func TestIntegrationLogCountHttp(t *testing.T) {
	scheme := "http"
	host := "localhost"
	leaderPort := 8080

	rlogger := rlogging.NewZapLogger()

	leader := NewHttpPeer(
		scheme,
		host,
		leaderPort,
		0,
		WithLogger(rlogger),
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
				WithLogger(rlogger),
			),
		)
	}

	for _, follower := range followers {
		if err := follower.Init(); err != nil {
			t.Fatalf(err.Error())
		}

		leader.AddPeer(follower)
	}

	if err := leader.Init(); err != nil {
		t.Fatalf(err.Error())
	}

	expectedLogCount := 0

	for _, follower := range followers {
		if actualLogCount, err := follower.LogCount(); err != nil {
			t.Fatalf("failed to retrieve log count due to error, %s\n", err.Error())
		} else if expectedLogCount != actualLogCount {
			t.Fatalf("expected %d, was %d", expectedLogCount, actualLogCount)
		}
	}
}

func TestIntegrationPeerCountHttp(t *testing.T) {
	scheme := "http"
	host := "localhost"
	leaderPort := 8090

	rlogger := rlogging.NewZapLogger()

	leader := NewHttpPeer(
		scheme,
		host,
		leaderPort,
		0,
		WithLogger(rlogger),
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
				WithLogger(rlogger),
			),
		)
	}

	for _, follower := range followers {
		leader.AddPeer(follower)
	}

	if err := leader.Init(); err != nil {
		t.Fatalf(err.Error())
	}

	expectedCount := len(followers)
	if actualCount, err := leader.PeerCount(); err != nil {
		t.Fatalf("failed to retrieve peer count due to error, %s\n", err.Error())
	} else if expectedCount != actualCount {
		t.Fatalf("expected %d, was %d", expectedCount, actualCount)
	}
}

func TestIntegrationSendHttp(t *testing.T) {
	scheme := "http"
	host := "localhost"
	leaderPort := 8100

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
				WithLog(
					rlog.NewBasicLog(
						rlog.WithBackend(
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
				),
			),
		)
	}

	for _, follower := range followers {
		follower.Init()

		leader.AddPeer(follower)
	}

	if err := leader.Init(); err != nil {
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
	if actualSendResult, err := leader.Send(
		dto.Message{
			Entries: expectedEntries,
		},
	); err != nil {
		t.Fatal(err)
	} else if expectedSendResult != actualSendResult {
		t.Fatalf("expectedSendResult %t, was %t\n", expectedSendResult, actualSendResult)
	}

	expectedPeerLogCount := len(expectedEntries)

	for _, p := range leader.Followers() {
		if actualPeerLogCount, err := p.LogCount(); err != nil {
			t.Fatalf("failed to retrieve log count due to error, %s\n", err.Error())
		} else if expectedPeerLogCount != actualPeerLogCount {
			t.Fatalf("expectedPeerLogCount %d, was %d\n", expectedPeerLogCount, actualPeerLogCount)
		} else {
			id := expectedPeerLogCount - 1

			expectedLatestEntry := expectedEntries[id]

			if actualPeerEntry, err := p.Entry(id); err != nil {
				t.Fatalf(
					"unable to retrieve entry with id %d, due to error, %s\n",
					id,
					err.Error(),
				)
			} else if expectedLatestEntry != actualPeerEntry {
				t.Fatalf("expectedLatestEntry %v, was %v", expectedLatestEntry.Command, actualPeerEntry.Command)
			} else if expectedLatestEntry.Id != actualPeerEntry.Id {
				t.Fatalf("expectedLatestEntry %v, was %v", expectedLatestEntry.Id, actualPeerEntry.Id)
			}
		}
	}
}
