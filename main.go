package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type event struct {
	UserID  string `json:"userId"`
	Payload string
}

type eventResponse struct {
	Msg string `json:"msg"`
}

func parseRequest(req *http.Request) (ev event) {
	// No validation on request method or body yet
	decoder := json.NewDecoder(req.Body)
	_ = decoder.Decode(&ev)
	log.Printf("%+v\n", ev)
	return
}

func sendResponse(w http.ResponseWriter, msg string) {
	var resp = eventResponse{
		Msg: msg,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func receiveEvent(w http.ResponseWriter, req *http.Request) {
	parseRequest(req)
	sendResponse(w, "Event received")
}

func main() {
	http.HandleFunc("/receive_event", receiveEvent)

	port := "8090"
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Unable to start http service : %s", err)
		return
	}
}
