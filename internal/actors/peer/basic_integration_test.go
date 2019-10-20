// +build integration

package actors

import (
	"cjvirtucio87/distributed-todo-go/internal/dto"
	"cjvirtucio87/distributed-todo-go/internal/rlog"
	"testing"
)

func TestIntegrationAddPeer(t *testing.T) {
	leader := NewBasicPeer(0)

	leader.AddPeer(NewBasicPeer(1))

	expectedCount := 1
	if actualCount, err := leader.PeerCount(); err != nil {
		t.Fatalf("failed to retrieve peer count due to error, %s\n", err.Error())
	} else if expectedCount != actualCount {
		t.Fatalf("expectedCount %d, was %d", expectedCount, actualCount)
	}
}

func TestIntegrationSend(t *testing.T) {
	leader := NewBasicPeer(0)

	for i := 1; i < 3; i++ {
		leader.AddPeer(NewBasicPeer(i))
	}

	expectedSendResult := true
	if actualSendResult, err := leader.Send(
		dto.Message{
			Entries: []dto.Entry{
				dto.Entry{Command: "doFoo"},
			},
		},
	); err != nil {
		t.Fatal(err)
	} else if expectedSendResult != actualSendResult {
		t.Fatalf("expectedSendResult %t, was %t", expectedSendResult, actualSendResult)
	}

	leader = NewBasicPeer(0)
	for i := 1; i < 3; i++ {
		leader.AddPeer(
			&basicPeer{
				id: i,
				rlog: rlog.NewBasicLog(
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
				NextIndexMap: map[int]int{},
				peers:        []Peer{},
			},
		)
	}

	expectedEntry := dto.Entry{Command: "doFoo"}
	if actualSendResult, err := leader.Send(
		dto.Message{
			Entries: []dto.Entry{
				expectedEntry,
			},
		},
	); err != nil {
		t.Fatal(err)
	} else if expectedSendResult != actualSendResult {
		t.Fatalf("expectedSendResult %t, was %t", expectedSendResult, actualSendResult)
	}

	expectedPeerLogCount := 1
	id := 0

	for _, p := range leader.Followers() {
		if actualPeerLogCount, err := p.LogCount(); err != nil {
			t.Fatalf("unable to retrieve log count for entry with id %d due to error, %s\n", id, err.Error())
		} else if expectedPeerLogCount != actualPeerLogCount {
			t.Fatalf("expectedPeerLogCount %d, was %d\n", expectedPeerLogCount, actualPeerLogCount)
		} else if actualPeerEntry, err := p.Entry(id); err != nil {
			t.Fatalf("unable to retrieve entry with id %d due to error, %s\n", id, err.Error())
		} else if expectedEntry != actualPeerEntry {
			t.Fatalf("expectedEntry %v, was %v", expectedEntry.Command, actualPeerEntry.Command)
		}
	}
}
