package actors

import (
	"fmt"
	"testing"
)

func TestAddPeer(t *testing.T) {
	leader := NewPeer("basic")

	leader.AddPeer(NewPeer("basic"))

	expected := 1
	actual := leader.Count()

	if expected != actual {
		t.Error(fmt.Printf("expected %d, was %d", expected, actual))
	}
}

func TestSend(t *testing.T) {
	leader := NewPeer("basic")

	for i := 0; i < 3; i++ {
		leader.AddPeer(NewPeer("basic"))
	}

	expected := true
	actual := leader.Send([]Entry{Entry{command: "doFoo"}})

	if expected != actual {
		t.Error(fmt.Printf("expected %t, was %t", expected, actual))
	}
}
