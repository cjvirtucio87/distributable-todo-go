package actors

import "testing"

func TestAddPeer(t *testing.T) {
	leader := NewPeer()

	leader.AddPeer(NewPeer())

	if leader.Count() != 1 {
		t.Fail()
	}
}

func TestSend(t *testing.T) {
	leader := NewPeer()

	for i := 0; i < 3; i++ {
		leader.AddPeer(NewPeer())
	}

	if !leader.Send([]Entry{Entry{command: "doFoo"}}) {
		t.Fail()
	}
}
