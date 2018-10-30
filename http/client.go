package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Client struct {
	httpClient http.Client
}

func NewClient() Client {
	return Client{
		httpClient: http.Client{
			Timeout: time.Millisecond * 500,
		},
	}
}

func (c *Client) SendAppendEntries(peer string, term int) (int, bool, error) {
	req := AppendEntriesRequest{Term:term}
	reqJson, _ := json.Marshal(req)
	res, err := c.httpClient.Post("http://" + peer + "/raft/append-entries", "application/json", bytes.NewBuffer(reqJson))
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

func (c *Client) SendRequestVote(peer string, term int, candidateId string) (int, bool, error) {
	req := RequestVoteRequest{Term:term, CandidateId:candidateId}
	reqJson, _ := json.Marshal(req)
	res, err := c.httpClient.Post("http://" + peer + "/raft/request-vote", "application/json", bytes.NewBuffer(reqJson))
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