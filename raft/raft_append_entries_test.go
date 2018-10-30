package raft

import "testing"

func TestShouldSwitchToFollowerIfStaleTermOnAppendEntries(t *testing.T) {
	var r = Raft{
		state:State{CurrentTerm: 1, Role: "leader"},
		restartElectionTickerChannel: make(chan int, 100),
	}

	r.AppendEntries(2)

	if r.State().Role != "follower" {
		t.Errorf("Role incorrect: should be follower was %v", r.State().Role)
	}

	if r.State().CurrentTerm != 2 {
		t.Errorf("Did not changed CurrentTerm %v. Should be 2", r.State().CurrentTerm)
	}

}