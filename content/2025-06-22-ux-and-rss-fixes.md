---
title: "Small UX improvements and RSS feed fixes"
date: 2025-06-22 2:44
slug: "ux-and-rss-fixes"
---

I just closed the [second
issue](https://github.com/RipeFruit08/blog-go/issues/4) on my [blog
project](https://github.com/RipeFruit08/blog-go) which fixed rss issues and
improved the reader experience.

## Background

When I first got my blog up and running I wanted to test the rss feed to ensure
it was being served up properly. I think
[RSS](https://en.wikipedia.org/wiki/RSS) is a really cool standard that allows
for sharing content is a standardized way. Initially, my rss feed was not
fully working. The feed would get consumed and would show all of my blog posts,
but when actually viewing those posts, nothing was rendering on screen besides
the post's title and date. The time was also incorrect

![RSS posts not rendering properly](/static/images/rss-issues.jpeg)

## Causes

Posts not rendering properly in rss was initially unclear to me but I knew why
the time on each post was incorrect - it was because I never actually added time
metadata to my posts. At first, I didn't think this would be a problem until I
wanted to write a second post for the day and had no way to let the blog know
which post came first

## Markdown parsing

My blog posts are written in Markdown, and then my website parses through that
to serve up the HTML for that post. These Markdown files begin with a header
containing metadata associated with the post

```text
# old header format
---
title: "Fixing the RSS feed and learning about Go's Embedding"
date: 2025-06-21
slug: "rss-fixes-go-embedding"
---
```

One of the problems with this header format is that it can only capture the date
associated with a post but no time metadata. To extend functionality, we need to
understand how my application is parsing the header information.

The blog has a notion of a `Post` object defined by the following

```go
// old Post struct
type Post struct {
  Title string
  Slug  string
  Date  time.Time
  HTML  template.HTML
}
```

The pipeline reads each markdown file then runs the contents through a
`parsePost` function which includes the following snippet

```go
func parsePost(data []byte) (Post, error) {
  var post Post

  content := data
  if bytes.HasPrefix(data, []byte("---")) {
    parts := bytes.SplitN(data, []byte("---"), 3)
    if len(parts) >= 3 {
      meta := parts[1]
      content = parts[2]
      if err := yaml.Unmarshal(meta, &post); err != nil {
        // relevant line ^^^
        return post, err
      }
    }
  }

  // additional code omitted...
}
```

`yaml.Unmarshal(meta, &post)` will do the heavy lifting of converting YAML data
into a struct in Go. This can be really convenient and magical if your YAML data
fields match the fields of your struct and if the types in your struct are built
in types. In this case, `meta` is our YAML data and contains the information
from the header in our markdown files, and our struct is the `Post` model
mentioned above. But, how does it know, for example, to parse the date
information from our header (e.g., 2025-06-21) and store it into an instance of
a `Post` struct in the Date field? Under the hood, it is looking through each
type (`time.Time` in this instance) and checking to see if it has a
corresponding implemention for the `yaml.Unmarshal` interface. If you're using
already defined types (e.g., string, time.Time, etc.), then it is very likely
that has already been taken care of for you.

## Extending the date field's functionality in our headers

Since the contents of my blog posts are written by me in markdown, ease
of populating the header metadata should be prioritized. In order for my posts
to also keep track of time, I would like the date in my header to take either
`yyyy-mm-dd` but also `yyyy-mm-dd hh:mm` where `hh:mm` goes off a 24 hour clock
(military time)

```text
# new header format
---
title: "Fixing the RSS feed and learning about Go's Embedding"
date: 2025-06-21 1:11
slug: "rss-fixes-go-embedding"
---
```

Without making any changes to the codebase, the `parsePost` function will break when
it reads our new Date field with a time component. This is because the default
implementation that unmarshals yaml data into `time.Time` does not expect data in this format.

To fix this, we will first update the `Post.Date` field to a custom type

```go
type Post struct {
  Title string
  Slug  string
  Date  CustomTime
  HTML  template.HTML
}
```

We now need to define what the `CustomTime` struct looks like, provide our own
implementation for `yaml.Unmarshal` so it knows how to handle it, and write our
own formatter function to control how we display this field.

```go
type CustomTime struct {
  time.Time
}

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

```

This helped make a noticable improvement on my blog posts when displaying the
date and time associated with the post. See before and after photos below

![Hello world post before updating how time was displayed. Date and time was
displayed as a messy string that was not very human
readable](/static/images/date-fix-before.png)

![Hello world post after updating how time was displayed. Date and time now show
as a nicely formatted date string that is easy to read for
humans](/static/images/date-fix-after.png)

## Fixing the RSS feed

The issues with my RSS feed were the following

- Posts showed, but there was no content rendering
- Following links from the feed did not properly redirect to my website's url

Resolving these issues were relatively straightforward. For the posts not
rendering anything, it turns out that for each item in my feed (e.g, each blog
post from my website), I was leaving out the `Description` field and it appears
that that is the field that is used to render the contents of your post in feed
readers.

For links not going to the right place, all of my hrefs in my feed were not
correctly linked up to my website to and I had to make updates to all of those
references. I took this opportunity to convert all occurrences to my website
into a `baseUrl` variable and to pull that value from a `BASE_URL` environment
variable. This way, if my domain ever changes I just need to make that update in
once place.
