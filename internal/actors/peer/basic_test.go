package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/rlog"
	"testing"
)

func TestSendSendsMessageToFollowers(t *testing.T) {
	peerCount := 3
	leaderLog := rlog.NewBasicLog()
	leader := &basicPeer{
		id:           0,
		rlog:         leaderLog,
		NextIndexMap: make(map[int]int),
		peers:        []Peer{},
	}

	followerLogs := make([]rlog.Log, peerCount)
	for i := 0; i < peerCount; i++ {
		followerLogs[i] = rlog.NewBasicLog()
	}

	for i := 0; i < peerCount; i++ {
		leader.AddPeer(
			&basicPeer{
				id:           i + 1,
				rlog:         followerLogs[i],
				NextIndexMap: make(map[int]int),
				peers:        []Peer{},
			},
		)
	}

	err := leader.Init()
	if err != nil {
		t.Fatal(err)
	}

	err = leader.Send(
		Message{
			Entries: []rlog.Entry{
				rlog.Entry{Command: "doFoo"},
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	for _, followerLog := range followerLogs {
		expectedEntryCount := 1
		entries := followerLog.Entries(0, 1)
		actualEntryCount := len(entries)
		if actualEntryCount != expectedEntryCount {
			t.Fatalf("expected [%d], got [%d]", expectedEntryCount, actualEntryCount)
		}

		expectedEntry := leaderLog.Entry(0)
		for _, actualEntry := range entries {
			if expectedEntry.Command != actualEntry.Command {
				t.Fatalf("expected [%s], got [%s]", expectedEntry.Command, actualEntry.Command)
			}
		}
	}
}

func TestSendDiscardsInvalidFollowerLogEntries(t *testing.T) {
	peerCount := 3
	leaderLog := rlog.NewBasicLog()
	leader := &basicPeer{
		id:           0,
		rlog:         leaderLog,
		NextIndexMap: make(map[int]int),
		peers:        []Peer{},
	}

	followerLogs := make([]rlog.Log, peerCount)
	for i := 0; i < peerCount; i++ {
		followerLogs[i] = rlog.NewBasicLog(
			rlog.WithBackend(
				[]rlog.Entry{
					rlog.Entry{
						Command: "not supposed to be here",
					},
					rlog.Entry{
						Command: "not supposed to be here either",
					},
				},
			),
		)
	}

	for i := 0; i < peerCount; i++ {
		leader.AddPeer(
			&basicPeer{
				id: i + 1,
				rlog: followerLogs[i],
				NextIndexMap: map[int]int{},
				peers:        []Peer{},
			},
		)
	}

	err := leader.Init()
	if err != nil {
		t.Fatal(err)
	}

	err = leader.Send(
		Message{
			Entries: []rlog.Entry{
				rlog.Entry{Command: "doFoo"},
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	for _, followerLog := range followerLogs {
		expectedEntryCount := 1
		entries := followerLog.Entries(0, 1)
		actualEntryCount := len(entries)
		if actualEntryCount != expectedEntryCount {
			t.Fatalf("expected [%d], got [%d]", expectedEntryCount, actualEntryCount)
		}

		expectedEntry := leaderLog.Entry(0)
		for _, actualEntry := range entries {
			if expectedEntry.Command != actualEntry.Command {
				t.Fatalf("expected [%s], got [%s]", expectedEntry.Command, actualEntry.Command)
			}
		}
	}
}
