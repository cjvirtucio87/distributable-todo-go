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
	actualCount := leader.PeerCount()

	if expectedCount != actualCount {
		t.Fatalf("expectedCount %d, was %d", expectedCount, actualCount)
	}
}

func TestIntegrationSend(t *testing.T) {
	leader := NewBasicPeer(0)

	for i := 1; i < 3; i++ {
		leader.AddPeer(NewBasicPeer(i))
	}

	expectedSendResult := true
	actualSendResult := leader.Send(
		dto.Message{
			Entries: []dto.Entry{
				dto.Entry{Command: "doFoo"},
			},
		},
	)

	if expectedSendResult != actualSendResult {
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
	actualSendResult = leader.Send(
		dto.Message{
			Entries: []dto.Entry{
				expectedEntry,
			},
		},
	)

	if expectedSendResult != actualSendResult {
		t.Fatalf("expectedSendResult %t, was %t", expectedSendResult, actualSendResult)
	}

	expectedPeerLogCount := 1

	for _, p := range leader.Followers() {
		actualPeerLogCount := p.LogCount()

		if expectedPeerLogCount != actualPeerLogCount {
			t.Fatalf("expectedPeerLogCount %d, was %d", expectedPeerLogCount, actualPeerLogCount)
		}

		id := 0

		var actualPeerEntry dto.Entry
		var err error

		if actualPeerEntry, err = p.Entry(id); err != nil {
			t.Fatalf("unable to retrieve entry with id %d\n", id)
		} else if expectedEntry != actualPeerEntry {
			t.Fatalf("expectedEntry %v, was %v", expectedEntry.Command, actualPeerEntry.Command)
		}
	}
}
