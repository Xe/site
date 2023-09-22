---
title: Let it Snow
date: 2018-12-17
for: the lols
tags:
 - fluff
---

I have very terribly added snow to this website for the holidays. See [the CSS](https://github.com/Xe/site/blob/00d930e46939fff5700735bf97a62beaa674eb02/css/snow.css) for how I did this, it's really low-tech. Feel free to steal this trick, it is low-effort for maximum niceness. I have the `background-color` of the `snowframe` class identical to the `background-color` of the main page. This and `opacity: 1.0` seems to be the ticket.

Happy holidays, all.

---

<link rel="stylesheet" href="/css/snow.css" />

More detailed usage:

```
<html>
  <head>
    <link rel="stylesheet" href="/css/snow.css" />
  </head>
  
  <body class="snow">
    <div class="container">
      <div class="snowframe">
        <!-- The rest of your page here -->
      </div>
    </div>
  </body>
</html>
```

Then you should have content not being occluded by snow.
