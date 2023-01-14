package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Event struct {
	ID   string `json:"id"`
	Data any    `json:"data"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var event Event
			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			eventPayloadInBytes, err := json.Marshal(event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			payload := bytes.NewBuffer(eventPayloadInBytes)
			postRes, _ := http.Post("http://localhost:3000/events", "application/json", payload)

			defer postRes.Body.Close()
			//commentRes, _ := http.Post("http.//localhost:3001/events", "application/json", payload)
			//
			//defer commentRes.Body.Close()
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
