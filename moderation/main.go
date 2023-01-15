package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strings"
)

type Event struct {
	ID   string `json:"id"`
	Data any    `json:"data"`
}

type Comment struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	PostID int    `json:"post_id"`
	Status string `json:"status"`
}

var moderatedWord = "orange"

func main() {
	r := chi.NewRouter()

	r.Post("/events", func(w http.ResponseWriter, r *http.Request) {
		var event Event

		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch event.ID {
		case "CommentCreated":
			data, err := json.Marshal(event.Data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var comment Comment
			err = json.Unmarshal(data, &comment)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var event Event
			if strings.Contains(strings.ToLower(comment.Text), moderatedWord) {
				comment.Status = "rejected"
				event = Event{
					ID:   "CommentModerated",
					Data: comment,
				}
			} else {
				comment.Status = "approved"
				event = Event{
					ID:   "CommentModerated",
					Data: comment,
				}
			}
			eventPayloadInBytes, err := json.Marshal(event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			eventPayload := bytes.NewBuffer(eventPayloadInBytes)
			eventBusRes, _ := http.Post("http://localhost:3005/events", "application/json", eventPayload)
			defer eventBusRes.Body.Close()

			w.WriteHeader(http.StatusOK)
		}
	})
	log.Println("starting the moderation server on 3003")
	err := http.ListenAndServe(":3003", r)
	if err != nil {
		log.Fatalln(err)
	}
}
