---
title: "Xeact 0.0.69: A Revolutionary Femtoframework For High Efficiency JavaScript Development"
date: 2021-11-18
tags:
 - javascript
 - framework
 - satire
 - xeact
---

[Writing JavaScript is so lame. All the tools require me to do so much bullshit
to get them to even compile. I shouldn't need to compile JavaScript to
JavaScript in order to deploy stuff to a webpage. How did it get this bad? I
wish there was something easier!](conversation://Cadey/coffee)

Is this you? Have you wished for something like this? Your prayers have been
answered! Keep reading this post to learn more about this revolutionary set of
tools that let you scale up and down as you need to.

<noscript>

[Hey, normally we try to make sure that this blog functions exactly as it would
without JavaScript enabled as it does with JavaScript enabled. This is an
exception. If you want the interactive parts of this post to work, you will need
to allow JavaScript to execute on this page. These examples require <a
href="https://caniuse.com/es6-module">ES6 module</a> support, so you may want to
try the most up to date version of Chrome, Firefox or
Safari.](conversation://Mara/hacker)

</noscript>

At Xeserv we strive to make solid advances forward that allow you to keep your
focus on what matters: your life and more importantly your production
applications. As a part of this mission, we have created a groundbreaking,
revolutionary femtoframework called [Xeact](https://github.com/Xe/Xeact)
([npm](https://www.npmjs.com/package/@xeserv/xeact)).

Xeact is a high performance, developer-efficient and overall ballin'
femtoframework for productive development in Javascript. It will take everything
that is complicated about frontend web development and throw it in the trash.

Don't trust this random webpage? You don't have to! Let's see it in action with
some testimonials that were collected from committed users of this amazing
framework:

```html
<blockquote id="testimonial">Loading...</blockquote>
```

Then you can import it from either
[unpkg](https://unpkg.com/@xeserv/xeact@0.0.69/xeact.js) or a local copy:

```javascript
import { g } from "/static/js/xeact.min.js";

const shuf = (arr) => {
    var rand, temp, i;

    for (i = arr.length - 1; i > 0; i -= 1) {
        rand = Math.floor((i + 1) * Math.random());
        temp = arr[rand];
        arr[rand] = arr[i];
        arr[i] = temp;
    }
    return arr;
};

const testimonials = shuf([
    "It Works™",
    "It shouldn't crash until the heat death of the universe",
    "A necessary addition to your tech stack",
    "Completely revolutionized our deployment cycle",
    "Our engineering team was blown away; versatile, powerful, constantly updated. It's everything we wanted and more",
    "Daunting at first, but works right out of the box, with great documentation and amazing support",
    "With a small footprint and big impact, this is a definite game changer",
    "Something something kill child before killing parent joke",
    "Made with Rust",
    "Reliable, Secure, Good performance, Easy to use UI/UX, High throughput",
    "It doesn't stalk me as much as the competitors!",
    "this software created synergy i didn't think was possible",
    "easy to use and it respects your privacy",
    "this software's cloud backend has only been hacked twice this year",
    "please, they have my wife and kids; i want to see them again",
    "The third best femtoframework available on the market"
]);

let i = 0;

const update = () => {
    i++;
    i = i % testimonials.length;
    g("testimonial").innerText = testimonials[i];
};

update();
setInterval(update, 3000);
```

<script type="module">
import { g } from "/static/js/xeact.min.js";

const shuf = (arr) => {
    var rand, temp, i;

    for (i = arr.length - 1; i > 0; i -= 1) {
        rand = Math.floor((i + 1) * Math.random());
        temp = arr[rand];
        arr[rand] = arr[i];
        arr[i] = temp;
    }
    return arr;
};

const testimonials = shuf([
    "It Works™",
    "It shouldn't crash until the heat death of the universe",
    "A necessary addition to your tech stack",
    "Completely revolutionized our deployment cycle",
    "Our engineering team was blown away; versatile, powerful, constantly updated. It's everything we wanted and more",
    "Daunting at first, but works right out of the box, with great documentation and amazing support",
    "With a small footprint and big impact, this is a definite game changer",
    "Something something kill child before killing parent joke",
    "Made with Rust",
    "Reliable, Secure, Good performance, Easy to use UI/UX, High throughput",
    "It doesn't stalk me as much as the competitors!",
    "this software created synergy i didn't think was possible",
    "easy to use and it respects your privacy",
    "this software's cloud backend has only been hacked twice this year",
    "please, they have my wife and kids; i want to see them again",
    "The third best femtoframework available on the market"
]);

const runTestimonalExample = () => {
  let i = 0;

  const update = () => {
      i++;
      i = i % testimonials.length;
      g("testimonial").innerText = testimonials[i];
  };

  update();
  setInterval(update, 3000);
};

window.runTestimonalExample = runTestimonalExample;
</script>

<blockquote id="testimonial">Loading...</blockquote>

<button onclick="runTestimonalExample()">Run demo</button>

And this is only the tip of the iceberg. Xeact is a poweruser's dream. It makes
it utterly trivial to compose more complicated things out of the same basic
tools that you have for free in the browser. Each of the main functions will be
covered in detail below:

## `g` - gets an element by ID

The most basic function is `g`, `g` lets you get an HTML element out of the DOM
so you can manipulate it, much like we did in the above testimonial example (no,
there is no way to stop it).

Usage is simple:

```html
<div id="g-example"></div>
```

```javascript
import { g } from "/static/js/xeact.min.js";

let elem = g("g-example");
elem.innerText = elem.innerText + "a";
```

<script type="module">
import { g } from "/static/js/xeact.min.js";

const runGExample = () => {
    let elem = g("g-example");
    elem.innerText = elem.innerText + "a";
};

window.runGExample = runGExample;
</script>

<div id="g-example"></div>

<button onclick="runGExample()">Run demo</button>

## `h` - creates an HTML element

Just getting HTML elements by themselves isn't that useful outside of very
simple cases. Sometimes you need to create new HTML elements to munge them back
into the DOM. Let's say you want to create a replica of this Mara conversation
blurb:

[These are some words and I am writing them.](conversation://Mara/hacker)

This actually gets expanded out to HTML that looks like this:

```html
<div class="conversation">
    <div class="conversation-picture conversation-smol">
        <picture>
            <source srcset="https://cdn.xeiaso.net/file/christine-static/stickers/mara/hacker.avif" type="image/avif">
            <source srcset="https://cdn.xeiaso.net/file/christine-static/stickers/mara/hacker.webp" type="image/webp">
            <img src="https://cdn.xeiaso.net/file/christine-static/stickers/mara/hacker.png" alt="Mara is hacker">
        </picture>
    </div>
    <div class="conversation-chat">&lt;<b>Mara</b>&gt; These are some words and I am writing them.</div>
</div>
```

HTML is a kind of a tree internally, so the `h` function lets you build element
trees:

```html
<div id="exampleConversationRoot">Loading...</div>
```

```javascript
import { g, h, x } from "/static/js/xeact.min.js";

const mkConversation = (who, mood, message) =>
    h("div", {className: "conversation"}, [
        h("div", {className: "conversation-picture conversation-smol"}, [
            h("picture", {}, [
                h("source", {type: "image/avif", srcset: `https://cdn.xeiaso.net/file/christine-static/stickers/${who.toLowerCase()}/${mood}.avif`}),
                h("source", {type: "image/webp", srcset: `https://cdn.xeiaso.net/file/christine-static/stickers/${who.toLowerCase()}/${mood}.webp`}),
                h("img", {alt: `${who} is ${mood}`, src: `https://cdn.xeiaso.net/file/christine-static/stickers/${who.toLowerCase()}/${mood}.png`})
            ])
        ]),
        h("div", {className: "conversation-chat"}, [
            h("span", {innerText: "<"}),
            h("b", {innerText: who}),
            h("span", {innerText: "> "}),
            h("span", {innerText: msg})
        ])
    ]);

// clear out #exampleConversationRoot
x(g("exampleConversationRoot"));
g("exampleConversationRoot")
    .append(mkConversation("Mara", "hacker", "These are some words and I am writing them."));
```

<script type="module">
import { g, h, x } from "/static/js/xeact.min.js";

const mkConversation = (who, mood, message) =>
    h("div", {className: "conversation"}, [
        h("div", {className: "conversation-picture conversation-smol"}, [
            h("picture", {}, [
                h("source", {type: "image/avif", srcset: `https://cdn.xeiaso.net/file/christine-static/stickers/${who.toLowerCase()}/${mood}.avif`}),
                h("source", {type: "image/webp", srcset: `https://cdn.xeiaso.net/file/christine-static/stickers/${who.toLowerCase()}/${mood}.webp`}),
                h("img", {alt: `${who} is ${mood}`, src: `https://cdn.xeiaso.net/file/christine-static/stickers/${who.toLowerCase()}/${mood}.png`})
            ])
        ]),
        h("div", {className: "conversation-chat"}, [
            h("span", {innerText: "<"}),
            h("b", {innerText: who}),
            h("span", {innerText: "> "}),
            h("span", {innerText: message})
        ])
    ]);

const runHExample = () => {
    x(g("exampleConversationRoot"));
    g("exampleConversationRoot").append(mkConversation("Mara", "hacker", "These are some words and I am writing them."));

    g("afterHExample").style = "";
}

window.runHExample = runHExample;
</script>

<div id="exampleConversationRoot">Loading...</div>

<button onclick="runHExample()">Run demo</button>

<div id="afterHExample" style="display:none">

[What...the heck? I didn't write those words! What is going on?](conversation://Mara/wat)

[Don't worry, it's all part of the plan](conversation://Cadey/enby)

</div>

You can use this to create whatever you need to create for your webapps. This
was inspired by the [Elm HTML
library](https://package.elm-lang.org/packages/elm/html/latest/). I may end up
making a companion library with common elements as an addon to Xeact that will
be distributed separately.

## `x` - remove all children of an element

This allows you to remove everything from an element so you can recreate it
anew. You should use this right before you add new things to an element to flip
the page on refresh:

```html
<ul id="exampleXRoot">
    <li>stuff</li>
    <li>stuff</li>
    <li>stuff</li>
    <li>stuff</li>
    <li>stuff</li>
</ul>
```

```javascript
import { g, h, x } from "/static/js/xeact.min.js";

const addStuff = () => {
    g("exampleXRoot").append(h("li", {innerText: "stuff"}));
};

const clearStuff = () => {
    x(g("exampleXRoot"));
};
```

<script type="module">
import { g, h, x } from "/static/js/xeact.min.js";

const addStuff = () => {
    g("exampleXRoot").append(h("li", {innerText: "stuff"}));
};

const clearStuff = () => {
    x(g("exampleXRoot"));
};

window.addXStuff = addStuff;
window.clearXStuff = clearStuff;
</script>

<ul id="exampleXRoot">
    <li>stuff</li>
    <li>stuff</li>
    <li>stuff</li>
    <li>stuff</li>
    <li>stuff</li>
</ul>

<button onclick="addXStuff()">Add</button> <button onclick="clearXStuff()">Clear</button>

## `u` - build relative/absolute URLs

A common task in frontend code is to make HTTP requests to the server ~~where
you can write the code in a _real_ language~~. `u` lets you build those URLs
quickly:

```javascript
import { u } from "/static/js/xeact.min.js";

console.log(u("/blog.json"));
```

<script type="module">
import { g, u } from "/static/js/xeact.min.js";

g("uExampleOne").textContent = u("/blog.json");
</script>

The [JSON Feed](https://www.jsonfeed.org/) for this blog can be found here:
<span id="uExampleOne"></span>

You can chain this into `fetch` like so:

```html
<div id="xeblog-root"></div>
```

```javascript
import { g, h, u, x } from "/static/js/xeact.min.js";

const div = (data = {}, children = []) => h("div", data, children);
const ahref = (to, text) => h("a", {href: to, innerText: text});

const h3 = (text, attrs = {}) => {
   attrs["innerText"] = text;
   return h("h3", attrs);
}

(async () => {
    let resp = await fetch(u("/blog.json"));
    if (!resp.ok) {
        throw new Error(`/blog.json status ${resp.status}`);
    }
    let feed = await resp.json();
    
    let root = g("xeblog-root");
    let content = div({}, [
        h3(feed.title),
        h("ul", {}, feed.items.map(item => h("li", {}, [ahref(item.url, item.title)])))
    ]);
    
    x(root);
    root.append(content);
})();
```

<script type="module">
import { g, h, u, x } from "/static/js/xeact.min.js";

const div = (data = {}, children = []) => h("div", data, children);
const ahref = (to, text) => h("a", {href: to, innerText: text});

const h3 = (text, attrs = {}) => {
   attrs["innerText"] = text;
   return h("h3", attrs);
}

const uDemo = async () => {
    let resp = await fetch(u("/blog.json"));
    if (!resp.ok) {
        throw new Error(`/blog.json status ${resp.status}`);
    }
    let feed = await resp.json();
    
    let root = g("xeblog-root");
    let content = div({}, [
        h3(feed.title),
        h("ul", {}, feed.items.map(item => h("li", {}, [ahref(item.url, item.title)])))
    ]);
    
    x(root);
    root.append(content);
};

window.runUExample = uDemo;
</script>

<div id="xeblog-root"></div>

<button onclick="runUExample()">Run</button>

## `r` - run a function once when the page is ready

`r` lets you defer execution of a function until everything in the page has
loaded. This is situationally useful and is difficult to demo.

```javascript
import { r } from "/static/js/xeact.min.js";

console.log("hi, open your dev console to see me");
```

<script type="module">
import { r } from "/static/js/xeact.min.js";

console.log("hi, open your dev console to see me");
</script>

---

What kind of awesome things can you create with Xeact? Use the hashtag `#xeact`
on Twitter and I'll maybe give you a shoutout!

---

EDIT(00:37 M11 19 2021): A prior version of this post made an incorrect
assertion about Facebook's patent grant situation and React. I was apparently
operating under old information. I fucked up.
