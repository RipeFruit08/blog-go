---
title: "Logging Visitors"
date: 2025-06-17
slug: "logging-visitors"
---

Today, I was able to link up this blog with a postgres database to log
information on visitors. When any endpoint is hit on my website, it will log the
type of method made (GET, POST, etc.), the path (/rss.xml), the UserAgent string,
the platform (macOS, iOS, Windows, etc.) the duration it took for the request to
complete, and a timestamp of when the event occurred. I was able to confirm on
my computer and phone that it was hitting my server. I also have my rss feed
saved on [https://netnewswire.com](NetNewsWire) and was able to confirm that it
had recently hit the /rss.xml to check for updates. Pretty cool!
