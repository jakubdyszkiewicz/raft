package main

import (
	".."
	"../http"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		panic("Pass NodeID as a first arg and peers as a next")
	}

	client := http.NewClient()
	r := raft.NewRaft(client.SendRequestVote, client.SendAppendEntries, args[0], args[1:])
	server := http.NewServer(r)
	server.Start(args[0])
}