---
title: I need help with AVIF files on iOS 16 Safari
date: 2022-09-13
tags:
 - avif
 - ios16
 - apple
 - ios
---

<xeblog-hero ai="Waifu Diffusion v1.2" file="happy-fireworks-pagodas" prompt="a beautiful landscape, studio ghibli, hatsune miku, fuji-san, cherry blossoms, pagoda, bubble tea, happy fireworks, ornate composition, matte painting"></xeblog-hero>

Hey there. I'm having some problems with AVIF files on iOS 16. If you're reading
this from an iPhone running iOS 16 you may see the problem.

I recently started using a new program called
[uploud](https://github.com/Xe/x/blob/master/cmd/uploud/main.go) to automate
converting and uploading things like my hero images and slides to [Backblaze
B2](https://www.backblaze.com/b2/cloud-storage.html), my storage provider.
However, images that I make with
[`github.com/Kagami/go-avif`](https://github.com/Kagami/go-avif) won't load on
iOS 16 Safari. Here are the encoding settings I am using:

```go
&avif.Options{
	Threads: runtime.GOMAXPROCS(0),
	Speed:   0,
	Quality: 48,
}
```

This creates images that can be viewed in Chrome:

![](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+2022-09-12+213003.png)

But not on iOS 16:

![](https://cdn.xeiaso.net/file/christine-static/blog/IMG_1499.png)

What am I doing wrong? I get no useful debug information from the browser
inspector on Safari. It just says the image is _wrong_ but won't tell me _why_.

Please [contact me](https://xeiaso.net/contact) if you have any ideas.
