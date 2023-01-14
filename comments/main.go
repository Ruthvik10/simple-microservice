package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

type Comment struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
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
		res, err := json.Marshal(comment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res = append(res, '\n')
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

	log.Println("starting the comments server on 3001")

	err := http.ListenAndServe(":3001", r)
	if err != nil {
		log.Fatalln(err)
	}
}
