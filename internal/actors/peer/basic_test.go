package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/rlog"
	"testing"
)

func TestSendSendsMessageToFollowers(t *testing.T) {
	testData := []struct{
		expectedEntryCount int
		followerLogEntries []rlog.Entry
		peerCount int
	}{
		{
			1,
			[]rlog.Entry{
				rlog.Entry{Command: "doFoo"},
			},
			3,
		},
		{
			1,
			[]rlog.Entry{
				rlog.Entry{
					Command: "not supposed to be here",
				},
				rlog.Entry{
					Command: "not supposed to be here either",
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
			entries := followerLog.Entries(0, 1)
			actualEntryCount := len(entries)
			if actualEntryCount != testDatum.expectedEntryCount {
				t.Fatalf("expected [%d], got [%d]", testDatum.expectedEntryCount, actualEntryCount)
			}

			expectedEntry := leaderLog.Entry(0)
			for _, actualEntry := range entries {
				if expectedEntry.Command != actualEntry.Command {
					t.Fatalf("expected [%s], got [%s]", expectedEntry.Command, actualEntry.Command)
				}
			}
		}
	}
}
