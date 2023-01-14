package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type Event struct {
	ID   string `json:"id"`
	Data any    `json:"data"`
}

type Post struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	PostID int    `json:"post_id"`
}

type PostWithComments map[int]*Post

func main() {
	posts := make(PostWithComments)

	r := chi.NewRouter()

	r.Post("/events", func(w http.ResponseWriter, r *http.Request) {
		var event Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		switch event.ID {
		case "PostCreated":
			data, err := json.Marshal(event.Data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var post Post
			err = json.Unmarshal(data, &post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			posts[post.ID] = &post
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
			if p, ok := posts[comment.PostID]; ok {
				p.Comments = append(p.Comments, comment)
			} else {
				http.Error(w, "post not found", http.StatusNotFound)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	})
	r.Get("/posts", func(w http.ResponseWriter, r *http.Request) {
		res, err := json.MarshalIndent(posts, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res = append(res, '\n')
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})
	log.Println("starting the query service post on 3002")
	err := http.ListenAndServe(":3002", r)
	if err != nil {
		log.Fatalln(err)
	}
}
