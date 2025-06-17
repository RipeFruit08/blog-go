// In /internal/blog/loader.go, implement a loader that:

// * Walks the /content/ folder

// * Parses frontmatter

// * Converts Markdown to HTML

// * Stores as a list of []Post
package blog

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
	"blog-go/content"

	"gopkg.in/yaml.v3"
	"github.com/yuin/goldmark"
	"html/template"
)


type Post struct {
	Title string
	Slug  string
	Date  time.Time
	HTML  template.HTML // <-- this tells Go it's safe to render as HTML
}

func LoadPosts() ([]Post, error) {
	var posts []Post

	err := fs.WalkDir(content.Content, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}

		data, err := content.Content.ReadFile(path)
		if err != nil {
			return err
		}

		post, err := parsePost(data)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		// Default slug to filename if not set
		if post.Slug == "" {
			post.Slug = strings.TrimSuffix(filepath.Base(path), ".md")
		}

		posts = append(posts, post)
		return nil
	})

	return posts, err
}

func parsePost(data []byte) (Post, error) {
	var post Post

	content := data
	if bytes.HasPrefix(data, []byte("---")) {
		parts := bytes.SplitN(data, []byte("---"), 3)
		if len(parts) >= 3 {
			meta := parts[1]
			content = parts[2]
			if err := yaml.Unmarshal(meta, &post); err != nil {
				return post, err
			}
		}
	}

	md := goldmark.New()
	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		return post, err
	}

	post.HTML = template.HTML(buf.String())
	return post, nil
}
