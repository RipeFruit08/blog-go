---
title: "Testing Goldmark"
date: 2025-06-21
slug: "testing-goldmark"
---

My blog entries are written in
[Markdown](https://en.wikipedia.org/wiki/Markdown), then converted and served as
HTML using [goldmark](https://github.com/yuin/goldmark), a Markdown parser for
[Go](https://go.dev).

This entry primarily serves to test the functionality of goldmark's Markdown
parser, though, is kind of redundant as the developer of goldmark links to their
own [goldmark playground](https://yuin.github.io/goldmark/playground/) to test
this exact functionality. The below source code comes from a [free Markdown to
HTML converter](https://markdowntohtml.com)

# Sample Markdown

This is some basic, sample markdown.

## Second Heading

- Unordered lists, and:

1. One
1. Two
1. Three

- More

> Blockquote

And **bold**, _italics_, and even \*italics and later **bold\***. Even ~~strikethrough~~. [A link](https://markdowntohtml.com) to somewhere.

And code highlighting:

```js
var foo = "bar";

function baz(s) {
  return foo + ":" + s;
}
```

Or inline code like `var foo = 'bar';`.

Or an image of bears

![bears](http://placebear.com/200/200)

The end ...

👋 It's me again, gonna try rendering an image based off an image hosted on the
blog. This photo is a screenshot of a shiny Totodile I hatched. The screenshot
was taken using [Delta Emulator's](https://faq.deltaemulator.com/) screenshot
feature.

![Shiny Totodile](/static/images/shiny-totodile.jpeg)
