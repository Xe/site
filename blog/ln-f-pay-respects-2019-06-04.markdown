---
title: "Introducing ln: The Natural Logging Library"
date: 2019-06-04
---

# Introducing ln: The Natural Logging Library

Logging is a very annoyingly complicated topic in programming. Sometimes you log too much and your log servers run out of space or need to only have a week's retention. Sometimes you log too little and are left recreating things from scratch when support tickets come in. Sometimes this means you need to go recreate the problem using your logging infrastructure to get the same error patterns.

Basically it's a mess. A lot of the popular Go libraries around it are [zero-allocation nanosecond scale](https://christine.website/blog/experimental-rilkef-2018-11-30) things that offer a lot of flexibility and speed, but ultimately make this entire endeavor painful more than it needs to be. None of them also seem to offer contextual storage of key->value fields. This means you have to pass a partially constructed logger around instead of the global one Just Doing The Right Thing.

So let's talk about my solution for this called [`ln`](https://github.com/Xe/ln), the [natural logging library](https://en.wikipedia.org/wiki/Natural_logarithm). `ln` is based on the idea of structured logging. `ln` uses key->value pairs like this:

```
// F ields for logging
f := ln.F{
  "azure_diamond": "hunter2",
  "meme_source": "http://bash.org/?244321",
}

ln.Log(ctx, f)
```

and this prints something like:

```
time="2009-11-10T23:00:00Z" azure_diamond=hunter2 meme_source="http://bash.org/?244321"
```

Simple, right?

You can also put a MUTABLE key->value F into a context:

```
func main() {
	ctx := context.Background()
	f := ln.F{
		"azure_diamond": "hunter2",
		"meme_source":   "http://bash.org/?244321",
	}
	ctx = ln.WithF(ctx, f)
	doSomethingElse(ctx)

	ln.Log(ctx)
}

func doSomethingElse(ctx context.Context) {
	ln.WithF(ctx, ln.F{
		"hi": "mom",
	})
}
```

And this [yields](https://play.golang.org/p/0-3-qPA7d6Y):

```
time="2009-11-10T23:00:00Z" azure_diamond=hunter2 meme_source="http://bash.org/?244321" hi=mom
```



---

This blogpost was converted from [this tweetstorm](https://twitter.com/theprincessxena/status/1129917083364597760).
