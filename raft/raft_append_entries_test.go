package raft

import "testing"

func Test_ShouldSwitchToFollowerIfStaleTermOnAppendEntries(t *testing.T) {
	// given
	var r = NewTestRaft(State{CurrentTerm: 1, Role: "leader"})

	// when
	r.AppendEntries(2)

	// then
	if r.State().Role != "follower" {
		t.Errorf("Role incorrect: should be follower was %v", r.State().Role)
	}

	if r.State().CurrentTerm != 2 {
		t.Errorf("Did not changed CurrentTerm %v. Should be 2", r.State().CurrentTerm)
	}

}