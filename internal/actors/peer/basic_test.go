package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/rlog"
	"testing"
)

func TestSendSendsMessageToFollowers(t *testing.T) {
	testData := []struct {
		followerLogEntries []*rlog.Entry
		leaderMessage      Message
		peerCount          int
		secondLeaderMessage *Message
	}{
		{
			[]*rlog.Entry{},
			Message{
				Entries: []*rlog.Entry{
					&rlog.Entry{Command: "doFoo"},
				},
			},
			3,
			nil,
		},
		{
			[]*rlog.Entry{
				&rlog.Entry{Command: "doFoo"},
			},
			Message{
				Entries: []*rlog.Entry{
					&rlog.Entry{Command: "doFoo"},
				},
			},
			3,
			nil,
		},
		{
			[]*rlog.Entry{
				&rlog.Entry{
					Command: "not supposed to be here",
				},
				&rlog.Entry{
					Command: "not supposed to be here either",
				},
			},
			Message{
				Entries: []*rlog.Entry{
					&rlog.Entry{Command: "doFoo"},
				},
			},
			3,
			nil,
		},
		{
			[]*rlog.Entry{},
			Message{
				Entries: []*rlog.Entry{
					&rlog.Entry{Command: "doFoo"},
					&rlog.Entry{Command: "doBar"},
					&rlog.Entry{Command: "doBaz"},
				},
			},
			3,
			nil,
		},
		{
			[]*rlog.Entry{},
			Message{
				Entries: []*rlog.Entry{
					&rlog.Entry{Command: "doFoo"},
					&rlog.Entry{Command: "doBar"},
					&rlog.Entry{Command: "doBaz"},
				},
			},
			3,
			&Message{
				Entries: []*rlog.Entry{
					&rlog.Entry{Command: "doHip"},
					&rlog.Entry{Command: "doHop"},
				},
			},
		},
	}

	for _, testDatum := range testData {
		leaderLog := rlog.NewBasicLog()
		leader := &basicPeer{
			id:           0,
			rlog:         leaderLog,
			NextIndexMap: make(map[int]int),
			peers:        []Peer{},
		}

		followerLogs := make([]rlog.Log, testDatum.peerCount)
		for i := 0; i < testDatum.peerCount; i++ {
			if len(testDatum.followerLogEntries) == 0 {
				followerLogs[i] = rlog.NewBasicLog()
			} else {
				followerLogs[i] = rlog.NewBasicLog(
					rlog.WithBackend(testDatum.followerLogEntries),
				)
			}
		}

		followers := make([]Peer, testDatum.peerCount)
		for i, _ := range followers {
			followers[i] = &basicPeer{
				id:           i + 1,
				rlog:         followerLogs[i],
				NextIndexMap: make(map[int]int),
				peers:        []Peer{},
			}

			leader.AddPeer(followers[i])
		}

		err := leader.Init()
		if err != nil {
			t.Fatal(err)
		}

		err = leader.Send(testDatum.leaderMessage)

		if err != nil {
			t.Fatal(err)
		}

		expectedEntries := leaderLog.Entries(0, -1)
		for _, followerLog := range followerLogs {
			actualEntries := followerLog.Entries(0, -1)

			expectedEntryCount := len(expectedEntries)
			actualEntryCount := len(actualEntries)
			if actualEntryCount != expectedEntryCount {
				t.Fatalf("expected [%d], got [%d]", expectedEntryCount, actualEntryCount)
			}

			for i, expectedEntry := range expectedEntries {
				actualEntry := actualEntries[i]
				if expectedEntry.Id != i {
					t.Fatalf("[%v] expected [%d], got [%d]", expectedEntries, expectedEntry.Id, i)
				}

				if expectedEntry.Id != actualEntry.Id {
					t.Fatalf("expected [%d], got [%d]", expectedEntry.Id, actualEntry.Id)
				}

				if expectedEntry.Command != actualEntry.Command {
					t.Fatalf("expected [%s], got [%s]", expectedEntry.Command, actualEntry.Command)
				}
			}
		}

		leader.Commit()

		expectedLastAppliedId := expectedEntries[len(expectedEntries)-1].Id
		actualLastAppliedId := leader.LastAppliedId()
		if expectedLastAppliedId != actualLastAppliedId {
			t.Fatalf("expected [%d], got [%d]", expectedLastAppliedId, actualLastAppliedId)
		}

		for _, follower := range followers {
			actualLastAppliedId = follower.LastAppliedId()
			if expectedLastAppliedId != actualLastAppliedId {
				t.Fatalf("expected [%d], got [%d]", expectedLastAppliedId, actualLastAppliedId)
			}
		}

		if testDatum.secondLeaderMessage != nil {
			err = leader.Send(*testDatum.secondLeaderMessage)

			if err != nil {
				t.Fatal(err)
			}

			expectedEntries := leaderLog.Entries(0, -1)
			for _, followerLog := range followerLogs {
				actualEntries := followerLog.Entries(0, -1)

				expectedEntryCount := len(expectedEntries)
				actualEntryCount := len(actualEntries)
				if actualEntryCount != expectedEntryCount {
					t.Fatalf("expected [%d], got [%d]", expectedEntryCount, actualEntryCount)
				}

				for i, expectedEntry := range expectedEntries {
					actualEntry := actualEntries[i]
					if expectedEntry.Id != i {
						t.Fatalf("[%v] expected [%d], got [%d]", expectedEntries, expectedEntry.Id, i)
					}

					if expectedEntry.Id != actualEntry.Id {
						t.Fatalf("expected [%d], got [%d]", expectedEntry.Id, actualEntry.Id)
					}

					if expectedEntry.Command != actualEntry.Command {
						t.Fatalf("expected [%s], got [%s]", expectedEntry.Command, actualEntry.Command)
					}
				}
			}

			leader.Commit()

			expectedLastAppliedId := expectedEntries[len(expectedEntries)-1].Id
			actualLastAppliedId := leader.LastAppliedId()
			if expectedLastAppliedId != actualLastAppliedId {
				t.Fatalf("expected [%d], got [%d]", expectedLastAppliedId, actualLastAppliedId)
			}

			for _, follower := range followers {
				actualLastAppliedId = follower.LastAppliedId()
				if expectedLastAppliedId != actualLastAppliedId {
					t.Fatalf("expected [%d], got [%d]", expectedLastAppliedId, actualLastAppliedId)
				}
			}
		}
	}
}
