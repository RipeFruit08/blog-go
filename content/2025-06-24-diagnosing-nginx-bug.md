---
title: "Diagnosing Nginx Bug"
date: 2025-06-24 00:24
slug: "diagnosing-nginx-bug"
---

The other day while I was out eating lunch, I went to my website to ensure it
was up and running. This website is the first public thing I've put out in the
world and I wanted to make sure it was available for whoever might be interested
in visiting my little corner of the internet. I was surprised to see something I
had not recognized before.

![A screenshot of the homepage of stephenkim.net but instead of my normal home
page it is showing a boilerplate page titled "Apache2 Ubuntu Default
Page"](/static/images/default-page.jpeg)

## Thoughts?

I had no clue what could be causing this and my immediate thought was that I was
being hacked 😅. It doesn't help I recently [logged visitors](/logging-visitors)
and get random requests from servers around the world (though, I think they're
mostly web crawlers and datacenters). Since I was out I couldn't get to a
computer to diagnose the issue. I have ssh access from my phone but working on
such a tiny screen is cumbersone. I was able to verify that I can still ssh into
my machine and my webserver was up and running - that assuaged my concerns of an
attack. Oddly enough, when I shut down my webserver I was still getting served
that "Apache2 Ubuntu Default Page" page which _should have_ made the issue
clear, but it wasn't obvious to me in the moment.

## What happened?

So what happened? I know
[Apache](https://en.wikipedia.org/wiki/Apache_HTTP_Server) is a webserver but I
was pretty sure that I wasn't even running Apache on my server. I confirmed this
suspicion by running `systemctl status apache2` and confirming that the service
was indeed not active.

When I got to my computer to diagnose the problem I hit my server again and
noticed that everything was running just fine. I checked it again on my phone
and it was also working fine. Strange, what was different?

I don't recall what compelled me to do this but I thought I would try
disconnecting my phone from wifi and hitting my webserver on cellular data. And
that caused me encounter this problem once more.

So, it turns out that this was a configuration problem. My
[Nginx](https://en.wikipedia.org/wiki/Nginx) server block for my domain was
listening for [IPv4](https://en.wikipedia.org/wiki/IPv4) connections _but not_
for [IPv6](https://en.wikipedia.org/wiki/IPv6) connections. I needed to update
my server block with the following statement

```nginx
server{
  listen 80;      # listens for IPv4 connections
  listen [::]:80; # listens for IPv6 connections <--- this is what I was missing

  # other config stuff below
}
```

Don't forget to test the config and reload nginx

```zsh
# test the configuration
sudo nginx -t

# reload Nginx
sudo systemctl reload nginx
```

## Why was I seeing the Apache2 Ubuntu Default Page?

My understanding of how Nginx functions when processing incoming requests is like this:

1. It will try to match the request with one of the Nginx config's server block.
1. If it does not find a matching server block it will fallback to a "default"
   server block to handle all other cases. You will most likely find it at
   `/etc/nginx/sites-available/default` but your `nginx.conf` may specify it
   elsewhere

In my case, the default server block was serving up an index file located at
`/var/www/html`. The contents of that file was a static page that was the
"Apache2 Default Ubuntu Page"

So, since the config for my `stephenkim.net` domain was not configured to listen
to IPv6 connections, it would not match that server block and thus would not
serve my go webserver pointing to this blog. Instead, it would fallback to my
default server block which will staticly serve the elusive Apache2 Ubuntu
Default Page. I should probably update that to something less confusing...

After looking into my home router settings, it turns out that I have IPv6 turned
off entirely so that would explain why I did not see this problem until I got on
cellular, where it was attempting to make an IPv6 connection.
