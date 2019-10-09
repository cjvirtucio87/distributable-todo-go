package actors

import (
	"fmt"
	"testing"
)

func TestAddPeer(t *testing.T) {
	leader := NewBasicPeer(0)

	leader.AddPeer(NewBasicPeer(1))

	expectedCount := 1
	actualCount := leader.PeerCount()

	if expectedCount != actualCount {
		t.Fatal(fmt.Printf("expectedCount %d, was %d", expectedCount, actualCount))
	}
}

func TestSend(t *testing.T) {
	leader := NewBasicPeer(0)

	for i := 1; i < 3; i++ {
		leader.AddPeer(NewBasicPeer(i))
	}

	expectedSendResult := true
	actualSendResult := leader.Send(
		Message{
			Entries: []Entry{
				Entry{Command: "doFoo"},
			},
		},
	)

	if expectedSendResult != actualSendResult {
		t.Fatal(fmt.Printf("expectedSendResult %t, was %t", expectedSendResult, actualSendResult))
	}

	leader = NewBasicPeer(0)
	for i := 1; i < 3; i++ {
		leader.AddPeer(
			&basicPeer{
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
		)
	}

	expectedEntry := Entry{Command: "doFoo"}
	actualSendResult = leader.Send(
		Message{
			Entries: []Entry{
				expectedEntry,
			},
		},
	)

	if expectedSendResult != actualSendResult {
		t.Fatal(fmt.Printf("expectedSendResult %t, was %t", expectedSendResult, actualSendResult))
	}

	expectedPeerLogCount := 1

	for _, p := range leader.Followers() {
		actualPeerLogCount := p.LogCount()

		if expectedPeerLogCount != actualPeerLogCount {
			t.Fatal(fmt.Printf("expectedPeerLogCount %d, was %d", expectedPeerLogCount, actualPeerLogCount))
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
