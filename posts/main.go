package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
)

type Post struct {
	Title string `json:"title"`
	ID    int    `json:"id"`
}

type Event struct {
	ID   string `json:"id"`
	Data any    `json:"data"`
}

func main() {

	posts := make(map[int]Post)

	r := chi.NewRouter()

	r.Post("/posts", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var post Post
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		post.ID = len(posts) + 1
		posts[post.ID] = post

		res, err := json.Marshal(post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		res = append(res, '\n')

		event := Event{
			ID:   "PostCreated",
			Data: post,
		}

		eventPayloadInBytes, err := json.Marshal(event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		eventPayload := bytes.NewBuffer(eventPayloadInBytes)
		eventBusRes, _ := http.Post("http://localhost:3005/events", "application/json", eventPayload)
		defer eventBusRes.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	}))

	r.Get("/posts", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := json.Marshal(posts)
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
		res, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("Recieved Payload")
		log.Println(string(res))
		w.WriteHeader(http.StatusOK)
	}))

	log.Println("starting the posts server on 3000")

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatalln(err)
	}
}
