package raft

type RequestVoteAnswer struct {
	term int
	voteGranted bool
}

type AppendEntriesAnswer struct {
	term int
	success bool
}

type TestRaft struct {
	*Raft
	requestVoteAnswers map[string]RequestVoteAnswer
	appendEntriesAnswers map[string]AppendEntriesAnswer
	receivedAppendEntries map[string]int
}

func NewTestRaft(state State) *TestRaft {
	requestVoteAnswers := map[string]RequestVoteAnswer{}
	appendEntriesAnswers := map[string]AppendEntriesAnswer{}
	receivedAppendEntries := map[string]int{}

	reqVoteFun := func(peer string, term int, candidateId string) (int, bool, error) {
		if answer, ok := requestVoteAnswers[peer]; ok {
			return answer.term, answer.voteGranted, nil
		}
		return -1, false, nil
	}
	appendEntriesFun := func(peer string, term int) (int, bool, error) {
		if received, ok :=  receivedAppendEntries[peer]; ok {
			receivedAppendEntries[peer] = received + 1
		} else {
			receivedAppendEntries[peer] = 1
		}
		if answer, ok := appendEntriesAnswers[peer]; ok {
			return answer.term, answer.success, nil
		}
		return -1, false, nil
	}
	raft := NewRaft(reqVoteFun, appendEntriesFun, "testNode", []string{})
	raft.state = state
	return &TestRaft{
		Raft: raft,
		requestVoteAnswers: requestVoteAnswers,
		appendEntriesAnswers: appendEntriesAnswers,
		receivedAppendEntries: receivedAppendEntries,
	}
}

func (t *TestRaft) AnswerForRequestVote(peer string, term int, voteGranted bool) {
	t.requestVoteAnswers[peer] = RequestVoteAnswer{term, voteGranted}
}

func (t *TestRaft) AnswerForAppendEntry(peer string, term int, success bool) {
	t.appendEntriesAnswers[peer] = AppendEntriesAnswer{term, success}
}

func (t *TestRaft) ReceivedAppendEntries(peer string) int {
	if answer, ok := t.receivedAppendEntries[peer]; ok {
		return answer
	} else {
		return 0
	}
}