package main

import (
	"fmt"
	"log"
	"io/fs"
	"net/http"
	"blog-go/content"
	"blog-go/internal/blog"
	"blog-go/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	fmt.Println("Hello World!")
	entries, err := fs.ReadDir(content.Content, ".")
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		fmt.Println("Embedded:", e.Name())
	}


	posts, err := blog.LoadPosts()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	// Pass posts to handlers
	r.Get("/", handlers.Home(posts))
	r.Get("/rss.xml", handlers.RSS(posts))
	r.Get("/{slug}", handlers.ShowPost(posts))

	log.Println("Server listening on :3010")
	http.ListenAndServe(":3010", r)
}
