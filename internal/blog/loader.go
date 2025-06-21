// In /internal/blog/loader.go, implement a loader that:

// * Walks the /content/ folder

// * Parses frontmatter

// * Converts Markdown to HTML

// * Stores as a list of []Post
package blog

import (
	"blog-go/content"
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"html/template"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)


type Post struct {
	Title string
	Slug  string
	Date  CustomTime
	HTML  template.HTML // <-- this tells Go it's safe to render as HTML
}

type CustomTime struct {
	time.Time
}

// method on CustomTime pointer receiver
// Implement yaml.Unmarshaler interface
func (ct *CustomTime) UnmarshalYAML(value *yaml.Node) error {
	const fullLayout = "2006-01-02 15:04"
	const dateOnlyLayout = "2006-01-02"

	// Load the Eastern Time location
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return fmt.Errorf("failed to load timezone: %w", err)
	}

	// Try parsing full date+time
	t, err := time.ParseInLocation(fullLayout, value.Value, loc)
	if err != nil {
		// Try parsing date only, defaulting time to midnight in Eastern
		t, err = time.ParseInLocation(dateOnlyLayout, value.Value, loc)
		if err != nil {
			return fmt.Errorf("invalid date format: %s", value.Value)
		}
	}

	ct.Time = t
	return nil
}

// method on CustomTime value receiver
func (ct CustomTime) FormatReadable() string {
	return ct.Format("January 2, 2006 at 3:04 PM")
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
