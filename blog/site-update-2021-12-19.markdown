---
title: "Site Updates: Better Contrast Ratio and Using Xeact"
date: 2021-12-19
---

Happy holidays all! As the year rolls to a close I wanted to take a moment to
let you know of some improvements that I've done over the year to try and make
the reading experience as best as it can be.

## Better Contrast Ratio

Over the years I have gotten reports that my site is hard to read. I've been
taking them seriously but I have never really been sure what to do about them. I
think I have found the core of the problem and I have changed the site's
contrast ratio to hopefully have more contrast between the text of the website
and the background.

Here are some comparisons to before and after the change:

<div style="background-color:#000000;padding:0.5em">
  <div style="background-color:#282828;color:#fbf1c7;padding:1em">
    This is what my site looked like before I made the changes (dark mode).
  </div>
  
  <div style="padding:0.25em"></div>

  <div style="background-color:#1d2021;color:#f9f5d7;padding:1em">
    This is what my site looks like after I made the changes (dark mode).
  </div>
  
  <div style="padding:0.25em"></div>

  <div style="color:#282828;background-color:#fbf1c7;padding:1em">
    This is what my site looked like before I made the changes (light mode).
  </div>
  
  <div style="padding:0.25em"></div>

  <div style="color:#1d2021;background-color:#f9f5d7;padding:1em">
    This is what my site looks like after I made the changes (light mode).
  </div>
</div>

Hopefully this should improve the contrast ratio a lot more. I've always wanted
this website to look a lot like [my emacs
config](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+from+2021-12-19+12-06-39.png),
but these changes should hopefully reach a balance of readability and stylistic
choices to get across the vision I have for my website.

## Using Xeact

<noscript>

Sorry, but you may want to scroll past this section. At the time of writing I
don't currently have a good fallback set up for people that don't have
JavaScript enabled on their browser. If you have ideas, please [email
me](mailto:me@xeiaso.net) and let me know them.

</noscript>

I want to use [Xeact](https://xeiaso.net/blog/xeact-0.0.69-2021-11-18)
more in my website. I am trying to hit a balance of avoiding structural
JavaScript while also allowing me to experiment with new and interesting ways of
doing things. To this end I have created a custom HTML element that allows me to
embed the little conversation fragments that I feel makes this blog unique. They
are currently done by [a horrific hack I made to my markdown
parser](https://github.com/Xe/site/blob/540ae4a3a9735d3f55ebceb1d271e472cd7f950e/src/app/markdown.rs#L38-L65).
This also makes it hard for me to put links in the little conversation
fragments, so when you see me do that, that's because I've written something of
this form:

```markdown
[This is a test <a href="https://zombo.com">with a link</a>.](conversation://Mara/hacker)
```

[This is a test <a href="https://zombo.com">with a link</a>.](conversation://Mara/hacker)

This has become an unworkable mess. With the [html
component](https://github.com/Xe/site/blob/540ae4a3a9735d3f55ebceb1d271e472cd7f950e/static/js/conversation.js)
I made recently, I can instead write things that look like this:

```markdown
<xeblog-conv name="Mara" mood="hacker">
This is a test [with a link](https://zombo.com).
</xeblog-conv>
```

<xeblog-conv name="Mara" mood="hacker">

This is a test [with a link](https://zombo.com).

</xeblog-conv>

This lets me use Xeact in a way that allows me to enhance the experience of
using my website as well as my experience in writing it. I don't really like how
this adds JavaScript into the mix for rendering the page. I have tried to avoid
it, but this is getting unworkable for me. When I get things to the point that I
feel more comfortable with making it the default option (this may involve some
custom hacking to the CSS to make it degrade more gracefully, I don't know what
I'm doing here), I will use it more as my main way to write these asides. Until
then, I can deal with the link syntax form. I may end up writing something to
munge things into place.

This also depends on you having modern JavaScript support. I can't change
anything about that without introducing pain and suffering to my development
workflow.

<xeblog-conv name="Cadey" mood="coffee">

This also breaks people reading the RSS feed, but it's already been [pretty darn
broken](https://github.com/Xe/site/issues/388) already. Worst case this makes
people a bit confused with the RSS feed, but until I can find a good workaround
I think I can tolerate this weirdness.

I am not good with CSS.

</xeblog-conv>

## Server-side Syntax Highlighting

You may have noticed that code syntax highlighting has a notably different color
scheme to it compared to the rest of the blog as of late. This is not on
accident. I used to use [Prism](https://prismjs.com/) to do this on the client
side. This worked great, but Prism is a huge (near 1 MB) download. I really want
to avoid wasting bandwidth, so I added this at the markdown rendering step with
[syntect](https://docs.rs/syntect/latest/syntect/). The theme I use
(base16-mocha) is not perfectly aligned with my editor, but this is as close as
I can get with their default themes until I have the energy to port gruvbox into
this.

Hopefully this should work in the RSS feed too.

---

Have a good rest of the year and stay safe! I'm gonna try and take a load off
over the holidays, so I may end up posting less frequently. 2021 has been a lot.

[Please get vaccinated. I want to be able to travel for conventions without it
being so scary all the time.](conversation://Cadey/enby)
