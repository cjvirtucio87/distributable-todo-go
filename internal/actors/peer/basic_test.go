package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/rlog"
	"testing"
)

func TestSendSendsMessageToFollowers(t *testing.T) {
	testData := []struct{
		followerLogEntries []rlog.Entry
		leaderMessage Message
		peerCount int
	}{
		{
			[]rlog.Entry{},
			Message{
				Entries: []rlog.Entry{
					rlog.Entry{Command: "doFoo"},
				},
			},
			3,
		},
		{
			[]rlog.Entry{
				rlog.Entry{Command: "doFoo"},
			},
			Message{
				Entries: []rlog.Entry{
					rlog.Entry{Command: "doFoo"},
				},
			},
			3,
		},
		{
			[]rlog.Entry{
				rlog.Entry{
					Command: "not supposed to be here",
				},
				rlog.Entry{
					Command: "not supposed to be here either",
				},
			},
			Message{
				Entries: []rlog.Entry{
					rlog.Entry{Command: "doFoo"},
				},
			},
			3,
		},
		{
			[]rlog.Entry{},
			Message{
				Entries: []rlog.Entry{
					rlog.Entry{Command: "doFoo"},
					rlog.Entry{Command: "doBar"},
					rlog.Entry{Command: "doBaz"},
				},
			},
			3,
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

		for i := 0; i < testDatum.peerCount; i++ {
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

		err = leader.Send(testDatum.leaderMessage)

		if err != nil {
			t.Fatal(err)
		}

		for _, followerLog := range followerLogs {
			expectedEntries := leaderLog.Entries(0, -1)
			actualEntries := followerLog.Entries(0, -1)

			expectedEntryCount := len(expectedEntries)
			actualEntryCount := len(actualEntries)
			if actualEntryCount != expectedEntryCount {
				t.Fatalf("expected [%d], got [%d]", expectedEntryCount, actualEntryCount)
			}

			for i, expectedEntry := range expectedEntries {
				actualEntry := actualEntries[i]
				if expectedEntry.Id != actualEntry.Id {
					t.Fatalf("expected [%d], got [%d]", expectedEntry.Id, actualEntry.Id)
				}

				if expectedEntry.Command != actualEntry.Command {
					t.Fatalf("expected [%s], got [%s]", expectedEntry.Command, actualEntry.Command)
				}
			}
		}
	}
}
