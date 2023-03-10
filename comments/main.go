package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

type Comment struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Status string `json:"status"`
	PostID int    `json:"post_id"`
}

type Event struct {
	ID   string `json:"id"`
	Data any    `json:"data"`
}

func main() {
	comments := make(map[int][]Comment)

	r := chi.NewRouter()
	r.Post("/posts/{id:^[0-9]+}/comments", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postIDString := chi.URLParam(r, "id")
		postID, err := strconv.Atoi(postIDString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var comment Comment
		err = json.NewDecoder(r.Body).Decode(&comment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		comment.ID = len(comments[postID]) + 1
		comments[postID] = append(comments[postID], comment)
		comment.Status = "pending"
		comment.PostID = postID
		res, err := json.Marshal(comment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res = append(res, '\n')

		event := Event{
			ID:   "CommentCreated",
			Data: comment,
		}

		eventPayloadInBytes, err := json.Marshal(event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		eventPayload := bytes.NewBuffer(eventPayloadInBytes)
		eventBusRes, _ := http.Post("http://event-bus-srv:3005/events", "application/json", eventPayload)
		defer eventBusRes.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	}))
	r.Get("/posts/{id:^[0-9]+}/comments", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postIDString := chi.URLParam(r, "id")
		postID, err := strconv.Atoi(postIDString)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		res, err := json.Marshal(comments[postID])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res = append(res, '\n')
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}))

	r.Post("/events", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var event Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("Recieved Payload")
		log.Println(event)

		switch event.ID {
		case "CommentModerated":
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

			commentsOfAPost := comments[comment.PostID]
			commentIndex := -1

			for i, c := range commentsOfAPost {
				if c.ID == comment.ID {
					commentIndex = i
					break
				}
			}
			commentsOfAPost[commentIndex] = comment

			event := Event{
				ID:   "CommentUpdated",
				Data: comment,
			}

			eventPayloadInBytes, err := json.Marshal(event)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			eventPayload := bytes.NewBuffer(eventPayloadInBytes)
			eventBusRes, _ := http.Post("http://event-bus-srv:3005/events", "application/json", eventPayload)
			defer eventBusRes.Body.Close()
		}
		w.WriteHeader(http.StatusOK)
	}))

	log.Println("starting the comments server on 3001")

	err := http.ListenAndServe(":3001", r)
	if err != nil {
		log.Fatalln(err)
	}
}
