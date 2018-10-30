package main

import (
	".."
	raftHttp "../http"
	"log"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		panic("Pass NodeID as a first arg and peers as a next")
	}

	client := raftHttp.NewClient()
	r := raft.NewRaft(client.SendRequestVote, client.SendAppendEntries, args[0], args[1:])
	handlers := raftHttp.NewHandlers(r)

	http.HandleFunc("/raft/state", handlers.HandleState)
	http.HandleFunc("/raft/request-vote", handlers.HandleRequestVote)
	http.HandleFunc("/raft/append-entries", handlers.HandleAppendEntries)
	log.Fatal(http.ListenAndServe(args[0], logRequest(http.DefaultServeMux)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}