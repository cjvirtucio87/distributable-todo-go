// +build integration

package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/dto"
	"cjvirtucio87/distributed-todo-go/internal/rlogging"
	"reflect"
	"testing"
	"time"
)

type testFactory struct {
	followers []Peer
	leader    Peer
}

func newFactory(t *testing.T) *testFactory {
	rlogger := rlogging.NewZapLogger()

	scheme := "http"
	host := "127.0.0.1"
	leaderPort := 8080
	followers := []Peer{}

	leader := NewHttpPeer(
		scheme,
		host,
		leaderPort,
		0,
		WithLogger(rlogger),
	)

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

	channels := []chan error{}

	for _, follower := range followers {
		followerChannel := make(chan error)

		go func(followerChannel chan error, follower Peer) {
			followerChannel <- follower.Init()
		}(followerChannel, follower)

		leader.AddPeer(follower)

		channels = append(channels, followerChannel)
	}

	leaderChannel := make(chan error)

	go func() {
		leaderChannel <- leader.Init()
	}()

	timeoutChannel := make(chan error)

	go func() {
		t.Log("waiting..")

		for i := 0; i < 5; i++ {
			t.Logf("%d", i+1)

			time.Sleep(1 * time.Second)
		}

		t.Log("done waiting. no errors")

		timeoutChannel <- nil
	}()

	channels = append(
		channels,
		leaderChannel,
		timeoutChannel,
	)

	// inspired by: https://stackoverflow.com/a/19992525/5665947
	selectCases := make([]reflect.SelectCase, len(channels))

	for i, ch := range channels {
		selectCases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}

	_, value, _ := reflect.Select(selectCases)

	if valueInterface := value.Interface(); valueInterface != nil {
		t.Fatalf(valueInterface.(error).Error())
	}

	return &testFactory{
		followers: followers,
		leader:    leader,
	}
}

func TestIntegrationLogCountHttp(t *testing.T) {
	factory := newFactory(t)

	expectedLogCount := 0

	for _, follower := range factory.followers {
		if actualLogCount, err := follower.LogCount(); err != nil {
			t.Fatalf("failed to retrieve log count due to error, %s\n", err.Error())
		} else if expectedLogCount != actualLogCount {
			t.Fatalf("expected %d, was %d", expectedLogCount, actualLogCount)
		}
	}
}

func TestIntegrationPeerCountHttp(t *testing.T) {
	factory := newFactory(t)

	expectedCount := len(factory.followers)

	if actualCount, err := factory.leader.PeerCount(); err != nil {
		t.Fatalf("failed to retrieve peer count due to error, %s\n", err.Error())
	} else if expectedCount != actualCount {
		t.Fatalf("expected %d, was %d", expectedCount, actualCount)
	}
}

func TestIntegrationSendHttp(t *testing.T) {
	factory := newFactory(t)

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
	if actualSendResult, err := factory.leader.Send(
		dto.Message{
			Entries: expectedEntries,
		},
	); err != nil {
		t.Fatal(err)
	} else if expectedSendResult != actualSendResult {
		t.Fatalf("expectedSendResult %t, was %t\n", expectedSendResult, actualSendResult)
	}

	expectedPeerLogCount := len(expectedEntries)

	for _, p := range factory.leader.Followers() {
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
