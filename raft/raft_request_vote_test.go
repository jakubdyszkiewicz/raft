package raft

import "testing"

func Test_ShouldGrantVote(t *testing.T) {
	// given
	var r = NewTestRaft(State{CurrentTerm: 1, Role: "follower"})

	// when voted
	vote := r.RequestVote(1, "candidate")

	// then should grant vote
	if ! vote {
		t.Error("Should grant vote")
	}
	if r.State().VotedFor != "candidate" {
		t.Errorf("Should save voted for to: candidate, was %v", r.State().VotedFor)
	}
}

func Test_ShouldRejectVoteOnStaleTerm(t *testing.T) {
	// given
	var r = NewTestRaft(State{CurrentTerm: 2, Role: "follower"})

	// when voted with stale term
	vote := r.RequestVote(1, "candidate")

	// then should reject vote
	if vote {
		t.Error("Should reject vote because of stale term")
	}
}

func Test_ShouldSwitchToFollowerOnStaleTerm(t *testing.T) {
	// given
	var r = NewTestRaft(State{CurrentTerm: 1, Role: "leader"})

	// when
	r.RequestVote(2, "candidate")

	// then
	if r.State().Role != "follower" {
		t.Errorf("Role incorrect: should be follower, was %v", r.State().Role)
	}

	if r.State().CurrentTerm != 2 {
		t.Errorf("Did not changed CurrentTerm %v. Should be 2", r.State().CurrentTerm)
	}
}

func Test_ShouldRejectVoteWhenAlreadyVoted(t *testing.T) {
	// given
	var r = NewTestRaft(State{CurrentTerm: 1, Role: "follower", VotedFor:"someNode"})

	// when
	vote := r.RequestVote(1, "anotherCandidate")

	// then
	if vote {
		t.Error("Should reject vote because node already voted in this term")
	}
}
