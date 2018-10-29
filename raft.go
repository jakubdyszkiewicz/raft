package raft

import (
	"log"
	"math/rand"
	"time"
)

type State struct {
	CurrentTerm int `json:"currentTerm"`
	LastHeartbeat int64 `json:"lastHeartbeat"`
	Peers []string `json:"peers"`
	Role string `json:"role"`
	NodeId string `json:"nodeId"`
	VotedFor string `json:"votedFor"`
	VotesGranted int `json:"votesGranted"`
}

var state = State{
	Peers: []string{},
	Role:"follower"}

var restartElectionTickerChannel = make(chan(int), 100)

var leaderHeartbeatTicker = time.NewTicker(3 * time.Second)
var restartLeaderHeartbeatsChannel = make(chan(int), 100)

var RequestVoteFunc = func(peer string, term int, candidateId string) (int, bool, error) {
	return -1, false, nil
}
var AppendEntriesFunc = func(peer string, term int) (int, bool, error) {
	return -1, false, nil
}

func Start() {
	startElectionTicker()
	startLeaderHeartbeatsTicker()
}

func startElectionTicker() {
	var randomMillis = rand.Int() % 1000 + 4000
	electionTicker := time.NewTicker(time.Duration(randomMillis) * time.Millisecond)
	go func() {
		for {
			select {
			case <-electionTicker.C:
				HandleElectionTimeout()
			case <-restartElectionTickerChannel:
				go startElectionTicker()
				return
			}
		}
	}()
}

func startLeaderHeartbeatsTicker() {
	leaderHeartbeatTicker = time.NewTicker(3 * time.Second)
	go func() {
		for {
			select {
			case <-leaderHeartbeatTicker.C:
				sendHeartbeats()
			case <-restartLeaderHeartbeatsChannel:
				go startLeaderHeartbeatsTicker()
				return
			}
		}
	}()
}

func HandleElectionTimeout() {
	if isLeader() {
		return
	}
	log.Println("Handle election timeout")
	convertToCandidate()
	startElection()
}

func sendHeartbeats() {
	for _, peer := range state.Peers {
		if !isLeader() {
			return
		}
		log.Printf("Sending heartbeat to %v", peer)
		term, _, err := AppendEntriesFunc(peer, state.CurrentTerm)
		if err != nil {
			log.Printf("Error on append entries to peer %v: %v", peer, err)
		} else {
			updateTermIfNeeded(term)
		}
	}
	updateHeartbeat()
}

func updateHeartbeat() {
	state.LastHeartbeat = time.Now().UnixNano() / int64(time.Millisecond)
}

func isLeader() bool {
	return state.Role == "leader"
}

func startElection() {
	log.Print("Starting election")
	state.CurrentTerm++
	state.VotedFor = state.NodeId
	resetElectionTimer()
	var votesGranted = 1 // voted for myself
	for _, peer := range state.Peers {
		log.Printf("Sending request vote to peer %v", peer)
		term, voteGranted, err := RequestVoteFunc(peer, state.CurrentTerm, state.NodeId)
		log.Printf("Received %v term and vote %v", term, votesGranted)
		if err != nil {
			log.Printf("Error on request vote to peer %v: %v", peer, err)
		} else {
			updateTermIfNeeded(term)
			if !isCandidate() {
				return
			}
			if voteGranted {
				votesGranted++
			}
		}
	}
	state.VotesGranted = votesGranted
	allPeers := len(state.Peers) + 1
	if votesGranted > allPeers / 2 {
		log.Printf("Granted majority of votes %v of %v", votesGranted, allPeers)
		convertToLeader()
	} else {
		log.Printf("Did not granted majority of votes: %v of %v", votesGranted, allPeers)
	}
}
func isCandidate() bool {
	return state.Role == "candidate"
}

func convertToLeader() {
	log.Print("Converting to leader")
	state.Role = "leader"
	sendHeartbeats()
	startLeaderHeartbeatsTicker()
}

func convertToCandidate() {
	log.Print("Converting to candidate")
	state.Role = "candidate"
	resetElectionTimer()
}

func convertToFollower() {
	log.Print("Converting to follower")
	state.Role = "follower"
	state.VotedFor = ""
	resetElectionTimer()
}

func CurrentState() State {
	return state
}

func UpdatePeers(peers []string) {
	state.Peers = peers
}

func UpdateNodeId(nodeId string) {
	state.NodeId = nodeId
}

func RequestVote(term int, candidateId string) bool {
	log.Printf("Requesting vote for candidate %v and term %v", candidateId, term)
	updateTermIfNeeded(term)
	if term < state.CurrentTerm {
		log.Printf("Received lower term %v than current %v", term, state.CurrentTerm)
		return false
	}
	if state.VotedFor != "" {
		log.Printf("Already voted for %v", state.VotedFor)
		return false
	}
	state.VotedFor = candidateId
	log.Printf("Voted granted")
	return true
}
func isFollower() bool {
	return state.Role == "follower"
}

func AppendEntries(term int) bool {
	log.Printf("Appending entry for term %v", term)
	updateTermIfNeeded(term)
	resetElectionTimer()
	if term < state.CurrentTerm {
		log.Printf("Received lower term %v than current %v", term, state.CurrentTerm)
		return false
	}
	updateHeartbeat()
	log.Printf("Entry appended. Heartbeat %v", state.LastHeartbeat)
	return true
}

func resetElectionTimer() {
	log.Println("Reseting election timer")
	restartElectionTickerChannel <- 1
}

func updateTermIfNeeded(term int) {
	if state.CurrentTerm < term {
		log.Printf("Stale term %v, updating to %v", state.CurrentTerm, term)
		state.CurrentTerm = term
		state.VotesGranted = 0
		state.VotedFor = ""
		convertToFollower()
	}
}
