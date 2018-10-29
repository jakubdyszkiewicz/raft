package main

import (
	".."
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func handleState(writer http.ResponseWriter, request *http.Request) {
	res := raft.CurrentState()
	b, _ := json.Marshal(res)
	writer.Header().Add("Access-Control-Allow-Origin", "*")
	writer.Write(b)
}

type AppendEntriesRequest struct {
	Term int `json:"term"`
}

type AppendEntriesResponse struct {
	Term int `json:"term"`
	Success bool `json:"success"`
}

func handleAppendEntries(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	appendEntriesRequest := AppendEntriesRequest{}
	json.Unmarshal(body, &appendEntriesRequest)

	success := raft.AppendEntries(appendEntriesRequest.Term)

	appendEntriesResponse := AppendEntriesResponse{Term:raft.CurrentState().CurrentTerm, Success:success}
	b, _ := json.Marshal(appendEntriesResponse)
	writer.Write(b)
}

type RequestVoteRequest struct {
	Term int `json:"term"`
	CandidateId string `json:"candidateId"`
}

type RequestVoteResponse struct {
	Term int `json:"term"`
	VoteGranted bool `json:"voteGranted"`
}

func handleRequestVote(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	voteRequest := RequestVoteRequest{}
	json.Unmarshal(body, &voteRequest)

	voteGranted := raft.RequestVote(voteRequest.Term, voteRequest.CandidateId)

	voteResponse := RequestVoteResponse{VoteGranted:voteGranted, Term:raft.CurrentState().CurrentTerm}
	b, _ := json.Marshal(voteResponse)
	writer.Write(b)
}

func sendAppendEntries(peer string, term int) (int, bool, error) {
	req := AppendEntriesRequest{Term:term}
	reqJson, _ := json.Marshal(req)
	res, err := http.Post("http://" + peer + "/raft/append-entries", "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		return -1, false, err
	}

	if res.StatusCode != 200 {
		log.Printf("Received status code %v", res.StatusCode)
		return -1, false, nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1, false, err
	}
	var appendEntriesResponse = AppendEntriesResponse{}
	json.Unmarshal(body, &appendEntriesResponse)

	return appendEntriesResponse.Term, appendEntriesResponse.Success, nil
}

func sendRequestVote(peer string, term int, candidateId string) (int, bool, error) {
	req := RequestVoteRequest{Term:term, CandidateId:candidateId}
	reqJson, _ := json.Marshal(req)
	res, err := http.Post("http://" + peer + "/raft/request-vote", "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		return -1, false, err
	}

	if res.StatusCode != 200 {
		log.Printf("Received status code %v", res.StatusCode)
		return -1, false, nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1, false, err
	}

	var requestVoteResponse = RequestVoteResponse{}
	json.Unmarshal(body, &requestVoteResponse)

	return requestVoteResponse.Term, requestVoteResponse.VoteGranted, nil
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		panic("Pass NodeID as a first arg and peers as a next")
	}

	raft.UpdateNodeId(args[0])
	raft.UpdatePeers(args[1:])

	raft.Start()
	raft.AppendEntriesFunc = sendAppendEntries
	raft.RequestVoteFunc = sendRequestVote
	http.HandleFunc("/raft/state", handleState)
	http.HandleFunc("/raft/request-vote", handleRequestVote)
	http.HandleFunc("/raft/append-entries", handleAppendEntries)
	log.Fatal(http.ListenAndServe(args[0], logRequest(http.DefaultServeMux)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}