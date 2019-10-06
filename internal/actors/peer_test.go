package actors

import (
	"fmt"
	"testing"
)

func TestAddPeer(t *testing.T) {
	leader := NewPeer()

	leader.AddPeer(NewPeer())

	expected := 1
	actual := leader.Count()

	if expected != actual {
		t.Error(fmt.Printf("expected %d, was %d", expected, actual))
	}
}

func TestSend(t *testing.T) {
	leader := NewPeer()

	for i := 0; i < 3; i++ {
		leader.AddPeer(NewPeer())
	}

	expected := true
	actual := leader.Send([]Entry{Entry{command: "doFoo"}})

	if expected != actual {
		t.Error(fmt.Printf("expected %t, was %t", expected, actual))
	}
}
