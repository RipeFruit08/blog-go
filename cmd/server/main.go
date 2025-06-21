package main

import (
	"blog-go/content"
	"blog-go/internal/blog"
	"blog-go/internal/handlers"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	connStr := os.Getenv("DB_CONN_STRING")
	if connStr == "" {
		log.Fatal("DB_CONN_STRING not set")
	}

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

	// db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/your_db?sslmode=disable")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	r := chi.NewRouter()

	// Use built-in middleware or your own custom one
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(loggingAndDBMiddleware(db))

	// Serve static files under the "/static/*" route
	r.Route("/static", func(r chi.Router) {
		fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
		r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fs.ServeHTTP(w, r)
		}))
	})

	// Pass posts to handlers
	r.Get("/", handlers.Home(posts))
	r.Get("/rss.xml", handlers.RSS(posts))
	r.Get("/{slug}", handlers.ShowPost(posts))

	log.Println("Server listening on :3010")
	http.ListenAndServe(":3010", r)
}

func loggingAndDBMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			method := r.Method
			path := r.URL.Path
			userAgent := r.UserAgent()
			ip := getClientIP(r)
			platform := detectPlatform(userAgent)

			next.ServeHTTP(w, r)

			duration := time.Since(start).Milliseconds()

			// Log to stdout
			log.Printf("Request: %s %s from %s [%s] (%s) in %dms", method, path, ip, platform, userAgent, duration)

			// Insert into DB
			go func() {
				_, err := db.Exec(`
				INSERT INTO request_logs (method, path, ip_address, user_agent, platform, duration_ms)
				VALUES ($1, $2, $3, $4, $5, $6)`,
					method, path, ip, userAgent, platform, duration,
				)
				if err != nil {
					log.Printf("Error logging request to DB: %v", err)
				}
			}()
		})
	}
}

func getClientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.Split(fwd, ",")[0]
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func detectPlatform(ua string) string {
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "iphone"), strings.Contains(ua, "ipad"):
		return "iOS"
	case strings.Contains(ua, "android"):
		return "Android"
	case strings.Contains(ua, "mac os"):
		return "macOS"
	case strings.Contains(ua, "windows"):
		return "Windows"
	case strings.Contains(ua, "linux"):
		return "Linux"
	default:
		return "Unknown"
	}
}
