package raft

import "testing"

func Test_RequestVote_ShouldGrantVote(t *testing.T) {
	// given
	var r = Raft{
		state:State{CurrentTerm: 1, Role: "follower"},
		restartElectionTickerChannel: make(chan int, 100),
	}

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

func Test_RequestVote_RejectVoteOnStaleTerm(t *testing.T) {
	// given
	var r = Raft{
		state:State{CurrentTerm: 2, Role: "follower"},
		restartElectionTickerChannel: make(chan int, 100),
	}

	// when voted with stale term
	vote := r.RequestVote(1, "candidate")

	// then should reject vote
	if vote {
		t.Error("Should reject vote because of stale term")
	}
}

func Test_RequestVote_SwitchToFollowerOnStaleTerm(t *testing.T) {
	// given
	var r = Raft{
		state:State{CurrentTerm: 1, Role: "leader"},
		restartElectionTickerChannel: make(chan int, 100),
	}

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

func Test_RequestVote_RejectVoteWhenAlreadyVoted(t *testing.T) {
	// given
	var r = Raft{
		state:State{CurrentTerm: 1, Role: "follower", VotedFor:"someNode"},
		restartElectionTickerChannel: make(chan int, 100),
	}

	// when
	vote := r.RequestVote(1, "anotherCandidate")

	// then
	if vote {
		t.Error("Should reject vote because node already voted in this term")
	}
}
