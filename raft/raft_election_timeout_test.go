package raft

import "testing"

func Test_ShouldStartElectionOnTimeoutAndLooseIt(t *testing.T) {
	// given
	var r = Raft{
		state:State{NodeId: "mynode", CurrentTerm: 1, Role: "follower", Peers:[]string { "anotherpeer" }},
		restartElectionTickerChannel: make(chan int, 100),
		requestVoteFunc: func(peer string, term int, candidateId string) (int, bool, error) {
			return -1, false, nil
		},
	}

	// when
	r.handleElectionTimeout()

	// then
	if r.State().CurrentTerm != 2 {
		t.Errorf("Should bump term to 2, was %v", r.State().CurrentTerm)
	}
	if r.State().VotedFor != "mynode" {
		t.Errorf("Should vote for itself, but voted for %v", r.State().VotedFor)
	}
	if r.State().Role != "candidate" {
		t.Errorf("Role should be candidate, was %v", r.State().Role)
	}
	if r.State().VotesGranted != 1 {
		t.Errorf("Should have only 1 vote (voted for itself), had %v votes", r.State().VotedFor)
	}
}