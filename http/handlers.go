package http

import (
	".."
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Handlers struct {
	raft raft.Raft
}

func NewHandlers(raft raft.Raft) Handlers {
	return Handlers{
		raft: raft,
	}
}

func (h *Handlers) HandleState(writer http.ResponseWriter, request *http.Request) {
	res := h.raft.State()
	b, _ := json.Marshal(res)
	writer.Header().Add("Access-Control-Allow-Origin", "*")
	writer.Write(b)
	writer.Header().Add("content-type", "application/json")
}

type AppendEntriesRequest struct {
	Term int `json:"term"`
}

type AppendEntriesResponse struct {
	Term int `json:"term"`
	Success bool `json:"success"`
}

func (h *Handlers) HandleAppendEntries(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	appendEntriesRequest := AppendEntriesRequest{}
	json.Unmarshal(body, &appendEntriesRequest)

	success := h.raft.AppendEntries(appendEntriesRequest.Term)

	appendEntriesResponse := AppendEntriesResponse{Term:h.raft.State().CurrentTerm, Success:success}
	b, _ := json.Marshal(appendEntriesResponse)
	writer.Write(b)
	writer.Header().Add("content-type", "application/json")
}

type RequestVoteRequest struct {
	Term int `json:"term"`
	CandidateId string `json:"candidateId"`
}

type RequestVoteResponse struct {
	Term int `json:"term"`
	VoteGranted bool `json:"voteGranted"`
}

func (h *Handlers) HandleRequestVote(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	voteRequest := RequestVoteRequest{}
	json.Unmarshal(body, &voteRequest)

	voteGranted := h.raft.RequestVote(voteRequest.Term, voteRequest.CandidateId)

	voteResponse := RequestVoteResponse{VoteGranted:voteGranted, Term:h.raft.State().CurrentTerm}
	b, _ := json.Marshal(voteResponse)
	writer.Write(b)
	writer.Header().Add("content-type", "application/json")
}