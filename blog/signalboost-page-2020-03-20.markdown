---
title: "New Site Feature: Signal Boosting"
date: 2020-03-20
tags:
 - signalboost
---

In light of the [COVID-19][covid19] pandemic, people have been losing their
jobs. In normal times, this would be less of an issue, but in the middle of the
pandemic, HR departments have been reluctant to hire people as entire companies
suddenly switch to remote work. I feel utterly powerless during this outbreak. I
can only image what people who have lost their job feel. 

[covid19]: https://www.canada.ca/en/public-health/services/diseases/coronavirus-disease-covid-19.html

I've decided to do what I can to help. I have created a page on my website to
signal boost people who are looking for work. You can find it at
[`/signalboost`](/signalboost). If you want to be added to it, please open a
[GitHub issue](https://github.com/Xe/site/issues/new), [contact me](/contact),
or open a pull request to `signalboost.dhall` in the root.

The schema of this is simple:

```dhall
Person::{ Name = "Nicole Brennan"
        , Tags = [ "python", "go", "rust", "technical-writing" ]
        , GitLink = "https://github.com/Twi"
        , Twitter = "https://twitter.com/TwitterAccountNameHere"
        }
```

This will create a grid entry on the site that looks like this:

<div class="grid">
  <div class="cell -4of12 content"></div>
  <div class="cell -4of12 content">
    <big>Nicole Brennan</big>
    <p>python go rust technical-writing</p>
    <a href="https://github.com/Twi">GitHub</a> - <a href="https://twitter.com/TwitterAccountNameHere">Twitter</a>
  </div>
  <div class="cell -4of12 content"></div>
</div>

I've also changed my footer to point to this page for the forseeable future
instead of linking to my Patreon. Thank you for reading this and please take a
look at the people on [`/signalboost`](/signalboost).

Be well.
