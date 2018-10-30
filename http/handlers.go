package http

import (
	".."
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Handlers struct {
	raft *raft.Raft
}

func NewHandlers(raft *raft.Raft) *Handlers {
	return &Handlers{
		raft: raft,
	}
}

func (h *Handlers) HandleState(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Access-Control-Allow-Origin", "*")
	res := h.raft.State()
	writeJson(res, writer)
}

type AppendEntriesRequest struct {
	Term int `json:"term"`
}

type AppendEntriesResponse struct {
	Term int `json:"term"`
	Success bool `json:"success"`
}

func (h *Handlers) HandleAppendEntries(writer http.ResponseWriter, request *http.Request) {
	appendEntriesRequest := AppendEntriesRequest{}
	readJson(appendEntriesRequest, request)

	success := h.raft.AppendEntries(appendEntriesRequest.Term)

	appendEntriesResponse := AppendEntriesResponse{Term:h.raft.State().CurrentTerm, Success:success}
	writeJson(appendEntriesResponse, writer)
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
	voteRequest := RequestVoteRequest{}
	readJson(voteRequest, request)

	voteGranted := h.raft.RequestVote(voteRequest.Term, voteRequest.CandidateId)

	voteResponse := RequestVoteResponse{VoteGranted:voteGranted, Term:h.raft.State().CurrentTerm}
	writeJson(voteResponse, writer)
}

func readJson(obj interface{}, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(body, obj)
}

func writeJson(obj interface{}, writer http.ResponseWriter) {
	b, _ := json.Marshal(obj)
	writer.Write(b)
	writer.Header().Add("content-type", "application/json")
}