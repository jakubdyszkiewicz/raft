package raft

import "testing"

func Test_ShouldStartElectionOnTimeoutAndLooseIt(t *testing.T) {
	// given
	var r = NewTestRaft(State{NodeId: "mynode", CurrentTerm: 1, Role: "follower", Peers:[]string { "anotherpeer" }})

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

func Test_ShouldWinElectionWithMajorityOfVotes(t *testing.T) {
	// given
	var r = NewTestRaft(State{NodeId: "mynode", CurrentTerm: 1, Role: "follower", Peers:[]string { "anotherpeer" }})
	r.AnswerForRequestVote("anotherpeer", 1, true)
	r.AnswerForAppendEntry("anotherpeer", 1, true)

	// when
	r.handleElectionTimeout()

	// then
	if r.State().Role != "leader" {
		t.Errorf("Role should be leader, was %v", r.State().Role)
	}
	if r.State().VotesGranted != 2 {
		t.Errorf("Should have 2 vote, had %v votes", r.State().VotedFor)
	}
}

func Test_ShouldConvertToFollowerWhenReceivedVoteResponseWithNewerTerm(t *testing.T) {
	// given
	var r = NewTestRaft(State{NodeId: "mynode", CurrentTerm: 0, Role: "follower", Peers:[]string { "anotherpeer" }})
	r.AnswerForRequestVote("anotherpeer", 2, true)
	r.AnswerForAppendEntry("anotherpeer", 2, true)

	// when
	r.handleElectionTimeout()

	// then
	if r.State().Role != "follower" {
		t.Errorf("Role should be follower, was %v", r.State().Role)
	}
}