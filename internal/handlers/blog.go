package handlers

import (
	"blog-go/internal/blog"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/feeds"
)

func ShowPost(posts []blog.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		for _, p := range posts {
			if p.Slug == slug {
				tmpl := template.Must(template.ParseFiles("internal/assets/templates/post.html"))
				tmpl.Execute(w, p)
				return
			}
		}
		http.NotFound(w, r)
	}
}

func Home(posts []blog.Post) http.HandlerFunc {
	fmt.Printf("Loaded %d posts\n", len(posts))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("internal/assets/templates/index.html"))
		data := struct {
			Posts []blog.Post
		}{
			Posts: posts,
		}
		tmpl.Execute(w, data)
	}
}

func RSS(posts []blog.Post, baseUrl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		feed := &feeds.Feed{
			Title:       "Stephen Kim's Dev Blog",
			Link:        &feeds.Link{Href: baseUrl},
			Description: "My blog posts in RSS",
			Author:      &feeds.Author{Name: "Stephen Kim"},
			Created:     now,
		}

		for _, post := range posts {
			feed.Items = append(feed.Items, &feeds.Item{
				Title:       post.Title,
				Link:        &feeds.Link{Href: baseUrl + "/" + post.Slug},
				Description: string(post.HTML), // You can pull a summary here
				Created:     post.Date.Time,
			})
		}

		rss, err := feed.ToRss()
		if err != nil {
			http.Error(w, "Could not generate RSS", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
		_, err = w.Write([]byte(rss))
		if err != nil {
			log.Println("Error writing RSS response:", err)
		}
	}
}
