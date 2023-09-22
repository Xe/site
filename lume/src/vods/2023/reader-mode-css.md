---
title: Fixing Xesite in reader mode and RSS readers
date: 2023-01-21
vod:
  path: talks/vod/2023/01-21-reader-mode
tags:
  - css
  - xedn
  - imageProcessing
  - scalability
  - bugFix
---

When you are using reader mode in Firefox, Safari or Google Chrome, the browser rends control of the website's design and renders its own design. This is typically done in order to prevent people's bad design decisions from making webpages unreadable and also to strip away advertisements from content. As a website publisher, I rely on the ability to control the CSS of my blog a lot. This stream covers the research/implementation process for fixing some long-standing issues with the Xesite CSS and making a fix to XeDN so that the site renders acceptably in reader mode.

This stream covers the following topics:

* Understanding complicated CSS rules and creating fixes for issues with them
* Using content distribution networks (CDNs) to help reduce page load time for readers
* Implementing image resizing capabilities into an existing CDN program (XeDN)
* Design with end-users in mind