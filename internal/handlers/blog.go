package handlers

import (
	"net/http"
	"fmt"
	"blog-go/internal/blog"
	"github.com/go-chi/chi/v5"
	"html/template"
	"time"
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

func RSS(posts []blog.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		feed := &feeds.Feed{
			Title:       "My Go Blog",
			Link:        &feeds.Link{Href: "https://your-domain.com"},
			Description: "My blog posts in RSS",
			Author:      &feeds.Author{Name: "Your Name"},
			Created:     now,
		}

		for _, post := range posts {
			feed.Items = append(feed.Items, &feeds.Item{
				Title:       post.Title,
				Link:        &feeds.Link{Href: "https://your-domain.com/" + post.Slug},
				Description: "", // You can pull a summary here
				Created:     post.Date,
			})
		}

		rss, err := feed.ToRss()
		if err != nil {
			http.Error(w, "Could not generate RSS", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/rss+xml")
		w.Write([]byte(rss))
	}
}
