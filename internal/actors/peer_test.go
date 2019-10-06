package actors

import "testing"

func TestAddPeer(t *testing.T) {
	leader := NewPeer()

	leader.AddPeer(NewPeer())

	if len(leader.Peers) != 1 {
		t.Fail()
	}
}
