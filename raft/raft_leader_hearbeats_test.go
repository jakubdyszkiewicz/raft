package raft

import "testing"

func Test_ShouldSendHeartbeatsWhenNodeIsLeader(t *testing.T) {
	// given
	var r = NewTestRaft(State{CurrentTerm: 1, Role: "leader", Peers: []string{"p1", "p2"}})
	r.AnswerForAppendEntry("p1", 1, true)
	r.AnswerForAppendEntry("p2", 1, true)

	// when
	r.sendHeartbeats()

	// then
	if r.ReceivedAppendEntries("p1") != 1 {
		t.Errorf("p1 should receive 1 heartbeat, was %v", r.ReceivedAppendEntries("p1"))
	}
	if r.ReceivedAppendEntries("p2") != 1 {
		t.Errorf("p2 should receive 1 heartbeat, was %v", r.ReceivedAppendEntries("p2"))
	}
	if r.State().LastHeartbeat == 0 {
		t.Errorf("Did not update last heartbeat")
	}
	if r.State().Role != "leader" {
		t.Errorf("Should be leader, was %v", r.State().Role)
	}
}

func Test_ShouldSwitchToFollowerWhenReceivedAnswerWithHigherTermOnHeartbeat(t *testing.T) {
	// given
	var r = NewTestRaft(State{CurrentTerm: 0, Role: "leader", Peers: []string{"p1", "p2"}})
	r.AnswerForAppendEntry("p1", 2, false)

	// when
	r.sendHeartbeats()

	// then
	if r.State().Role != "follower" {
		t.Errorf("Should be follower, was %v", r.State().Role)
	}
}