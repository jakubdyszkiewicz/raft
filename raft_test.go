package raft

import "testing"

func TestShouldSwitchToFollowerIfStaleTermOnAppendEntries(t *testing.T) {
	state = State{CurrentTerm: 1, Role: "leader"}

	AppendEntries(2)

	if state.Role != "follower" {
		t.Errorf("Role incorrect: should be follower was %v", state.Role)
	}

	if state.CurrentTerm != 2 {
		t.Errorf("Did not changed CurrentTerm %v. Should be 2", state.CurrentTerm)
	}
}

func TestShouldSwitchToFollowerIfStaleTermOnRequestVote(t *testing.T) {
	state = State{CurrentTerm: 1, Role: "leader"}

	RequestVote(2, 1)

	if state.Role != "follower" {
		t.Errorf("Role incorrect: should be follower was %v", state.Role)
	}

	if state.CurrentTerm != 2 {
		t.Errorf("Did not changed CurrentTerm %v. Should be 2", state.CurrentTerm)
	}
}