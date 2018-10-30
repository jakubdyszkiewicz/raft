package http

import (
	"../raft"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Server struct {
	raft *raft.Raft
}

func NewServer(raft *raft.Raft) *Server {
	return &Server{
		raft: raft,
	}
}

func (s *Server) Start(listenAddr string) {
	http.HandleFunc("/raft/state", s.HandleState)
	http.HandleFunc("/raft/request-vote", s.HandleRequestVote)
	http.HandleFunc("/raft/append-entries", s.HandleAppendEntries)
	log.Fatal(http.ListenAndServe(listenAddr, logRequest(http.DefaultServeMux)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func (s *Server) HandleState(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Access-Control-Allow-Origin", "*")
	res := s.raft.State()
	writeJson(res, writer)
}

type AppendEntriesRequest struct {
	Term int `json:"term"`
}

type AppendEntriesResponse struct {
	Term int `json:"term"`
	Success bool `json:"success"`
}

func (s *Server) HandleAppendEntries(writer http.ResponseWriter, request *http.Request) {
	appendEntriesRequest := AppendEntriesRequest{}
	readJson(&appendEntriesRequest, request)

	success := s.raft.AppendEntries(appendEntriesRequest.Term)

	appendEntriesResponse := AppendEntriesResponse{Term:s.raft.State().CurrentTerm, Success:success}
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

func (s *Server) HandleRequestVote(writer http.ResponseWriter, request *http.Request) {
	voteRequest := RequestVoteRequest{}
	readJson(&voteRequest, request)

	voteGranted := s.raft.RequestVote(voteRequest.Term, voteRequest.CandidateId)

	voteResponse := RequestVoteResponse{VoteGranted:voteGranted, Term:s.raft.State().CurrentTerm}
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