package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Event struct {
	ID   string `json:"id"`
	Data any    `json:"data"`
}

func main() {
	mux := http.NewServeMux()
	numberOfServices := 4
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var event Event
			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			log.Println("Recieved Payload: ", event)
			eventPayloadInBytes, err := json.Marshal(event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var wg sync.WaitGroup
			wg.Add(numberOfServices)

			go func(wg *sync.WaitGroup) {
				postPayload := bytes.NewBuffer(eventPayloadInBytes)
				postRes, _ := http.Post("http://localhost:3000/events", "application/json", postPayload)
				defer wg.Done()
				defer postRes.Body.Close()
			}(&wg)
			go func(wg *sync.WaitGroup) {
				commentPayload := bytes.NewBuffer(eventPayloadInBytes)
				commentRes, _ := http.Post("http://localhost:3001/events", "application/json", commentPayload)
				defer wg.Done()
				defer commentRes.Body.Close()
			}(&wg)
			go func(wg *sync.WaitGroup) {
				queryPayload := bytes.NewBuffer(eventPayloadInBytes)
				queryRes, _ := http.Post("http://localhost:3002/events", "application/json", queryPayload)
				defer wg.Done()
				defer queryRes.Body.Close()
			}(&wg)

			go func(wg *sync.WaitGroup) {
				moderationPayload := bytes.NewBuffer(eventPayloadInBytes)
				moderationRes, _ := http.Post("http://localhost:3003/events", "application/json", moderationPayload)
				defer wg.Done()
				defer moderationRes.Body.Close()
			}(&wg)

			w.WriteHeader(http.StatusOK)

		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	log.Println("starting the event bus server on port 3005")
	err := http.ListenAndServe(":3005", mux)
	log.Fatalln(err)
}
