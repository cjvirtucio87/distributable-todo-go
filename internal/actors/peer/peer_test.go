package actors

import (
	"fmt"
	"testing"
)

func TestAddPeer(t *testing.T) {
	leader := NewPeer("basic", 0)

	leader.AddPeer(NewPeer("basic", 1))

	expected := 1
	actual := leader.Count()

	if expected != actual {
		t.Error(fmt.Printf("expected %d, was %d", expected, actual))
	}
}

func TestSend(t *testing.T) {
	leader := NewPeer("basic", 0)

	for i := 1; i < 3; i++ {
		leader.AddPeer(NewPeer("basic", i))
	}

	expected := true
	actual := leader.Send(
		Message{
			entries: []Entry{
				Entry{command: "doFoo"},
			},
		},
	)

	if expected != actual {
		t.Error(fmt.Printf("expected %t, was %t", expected, actual))
	}
}
