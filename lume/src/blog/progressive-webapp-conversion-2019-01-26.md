---
title: How To Make a Progressive Web App Out Of Your Existing Website
date: 2019-01-26
also_for: "Heroku's blog (TODO: put link here)"
thanks: Nathanial, Andrew Konoff
---

[Progressive web apps](https://developer.mozilla.org/en-US/docs/Web/Apps/Progressive) enable websites to trade some flexibility to function more like native apps, without all the overhead of app store approvals and tons of platform-specific native code. Progressive web apps allow users to install them to their home screen and launch them into their own pseudo-app frame. However, that frame is locked down and restricted, and only allows access to pages that are subpaths of the scope of the progressive web app. They also have to be served over HTTPS. Updates to these can be deployed without needing to wait for app store approval.

The core of progressive web apps are [service workers](https://developers.google.com/web/fundamentals/primers/service-workers/), which are effectively client-side Javascript daemons. Service workers can listen for a few kinds of events and react to them. One of the most commonly supported events is the [fetch event](https://jakearchibald.github.io/isserviceworkerready/#fetch-event); this can be used to cache web content offline as explained below.

There are a large number of web apps that fit just fine within these [rules and restrictions](https://developer.mozilla.org/en-US/docs/Web/Apps/Progressive/App_structure), however there could potentially be compatibility issues with existing code. Instead of waiting for Apple or Google to approve and push out app updates, service worker (and by extension progressive web app) updates will be fetched [following standard HTTP caching rules](https://stackoverflow.com/questions/38843970/service-worker-javascript-update-frequency-every-24-hours). Plus, you get to use plenty of native APIs, including geolocation, camera, and sensor APIs that only native mobile apps used to be able to take advantage of.

In this post, we’ll show you how to convert your existing website into a progressive web app. It’s fairly simple, only really requiring the following steps:

* Creating an app manifest
* Adding it to your base HTML template
* Creating the [service worker](https://developers.google.com/web/fundamentals/primers/service-workers/)
    * Serving the service worker on the root of the scope you used in the manifest
* Adding a `<script>` block to your base HTML template to load the service worker
* Deploying
* Using Your Progressive Web App

If you want a more guided version of this post, the folks at [https://pwabuilder.com](https://pwabuilder.com) have created an [online interface](https://www.pwabuilder.com/generate) for doing most of the below steps automatically.

## Creating an app manifest

An [app manifest](https://developer.mozilla.org/en-US/docs/Web/Manifest) is a combination of the following information:

* The canonical name of the website
* A short version of that name (for icons)
* The theme color of the website for OS integration
* The background color of the website for OS integration
* The URL scope that the progressive web app is limited to
* The start URL that new instances of the progressive web app will implicitly load
* A human-readable description
* Orientation restrictions (it is unwise to change this from `"any"` without a hard technical limit)
* Any icons for your website to be used on the home screen (see the above manifest generator for autogenerating icons)

This information will be used as the OS-level metadata for your progressive web app when it is installed.

Here is an example web app manifest [from my portfolio site](https://github.com/Xe/site/blob/d7b817b22db9e10dbdfe97082ec2330e13cfff18/static/manifest.json).

```json
{
    "name": "Xeblog",
    "short_name": "Xeblog",
    "theme_color": "#ffcbe4",
    "background_color": "#fa99ca",
    "display": "standalone",
    "scope": "/",
    "start_url": "/",
    "description": "Blog and Resume for Xe Iaso",
    "orientation": "any",
    "icons": [
        {
            "src": "/static/img/avatar.png",
            "sizes": "1024x1024"
        }
    ]
}
```

If you just want to create a manifest quickly, check out [this](https://app-manifest.firebaseapp.com) online wizard.

## Add Manifest to Your Base HTML Template

I suggest adding the HTML link for the manifest to the most base HTML template you can, or in the case of a purely client side web app its main `index.html` file, as it needs to be as visible by the client trying to install the app. Adding this is [simple](https://developer.mozilla.org/en-US/docs/Web/Apps/Progressive/Installable_PWAs), assuming you are hosting this manifest on [/static/manifest.json](https://xeiaso.net/static/manifest.json) – simply add it to the `<head>` section:

```html
<link rel="manifest" href="/static/manifest.json">
```

## Create offline.html as an alias to index.html

By default the service worker code below will render `/offline.html` instead of any resource it can't fetch while offline. Create a file at `<your-scope>/offline.html` to give your user a more helpful error message, explaining that this data isn't cached and the user is offline.

If you are adapting a single-page web app, you might want to make `offline.html` a symbolic link to your `index.html` file and have the offline 404 handler be done inside there. If users can't get back out of the offline page, it can potentially confuse or strand users at a fairly useless looking and feeling "offline" screen; this obviates a lot of the point of progressive web apps in the first place. Be sure to have some kind of "back" button on all error pages.


To set up a symbolic link if you are adapting a single-page web app, just enter this in your console:

```
$ ln -s index.html offline.html
```

Now we can create and add the service worker.

## Creating The Service Worker

When service workers are used with the [fetch event](https://developer.mozilla.org/en-US/docs/Web/API/FetchEvent), you can set up caching of assets and pages as the user browses. This makes content available offline and loads it significantly faster. We are just going to focus on the offline caching features of service workers today instead of automated background sync, [because iOS doesn't support background sync yet](https://jakearchibald.github.io/isserviceworkerready/).

At a high level, consider what assets and pages you want users of your website to always be able to access some copy of (even if it goes out of date). These pages will additionally be cached for every user to that website with a browser that supports service workers. I suggest implicitly caching at least the following:

* Any CSS, Javascript or image files core to the operations of your website that your starting route does not load
* Contact information for the person, company or service running the progressive web app
* Any other pages or information you might find useful for users of your website

For example, I have the following precached for [my portfolio site](https://xeiaso.net):

* My homepage (implicitly includes all of the CSS on the site) `/`
* My blog index `/blog/`
* My contact information `/contact`
* My resume `/resume`
* The offline information page `/offline.html`

And this translates into the following service worker code:

```javascript
self.addEventListener("install", function(event) {
  event.waitUntil(preLoad());
});

var preLoad = function(){
  console.log("Installing web app");
  return caches.open("offline").then(function(cache) {
    console.log("caching index and important routes");
    return cache.addAll(["/blog/", "/blog", "/", "/contact", "/resume", "/offline.html"]);
  });
};

self.addEventListener("fetch", function(event) {
  event.respondWith(checkResponse(event.request).catch(function() {
    return returnFromCache(event.request);
  }));
  event.waitUntil(addToCache(event.request));
});

var checkResponse = function(request){
  return new Promise(function(fulfill, reject) {
    fetch(request).then(function(response){
      if(response.status !== 404) {
        fulfill(response);
      } else {
        reject();
      }
    }, reject);
  });
};

var addToCache = function(request){
  return caches.open("offline").then(function (cache) {
    return fetch(request).then(function (response) {
      console.log(response.url + " was cached");
      return cache.put(request, response);
    });
  });
};

var returnFromCache = function(request){
  return caches.open("offline").then(function (cache) {
    return cache.match(request).then(function (matching) {
     if(!matching || matching.status == 404) {
       return cache.match("offline.html");
     } else {
       return matching;
     }
    });
  });
};
```

You host the above at `<your-scope>/sw.js`. This file must be served from the same level as the scope. There is no way around this, unfortunately.

## Load the Service Worker

To load the service worker, we just add the following to your base HTML template at the end of your `<body>` tag:

```html
<script>
 if (!navigator.serviceWorker.controller) {
     navigator.serviceWorker.register("/sw.js").then(function(reg) {
         console.log("Service worker has been registered for scope: " + reg.scope);
     });
 }
</script>
```

And then deploy these changes – you should see your service worker posting logs in your browser’s console. If you are testing this from a phone, see platform-specific instructions [here for iOS+Safari](https://www.dummies.com/web-design-development/how-to-use-developer-tools-in-safari-on-ios/) and [here for Chrome+Android](https://developers.google.com/web/tools/chrome-devtools/remote-debugging/?hl=en).

## Deploying

Deploying your web app is going to be specific to how your app is developed. If you don't have a place to put it already, [Heroku](https://heroku.com) offers a nice and simple way to host progressive web apps. Using [the static buildpack](https://github.com/heroku/heroku-buildpack-static) is the fastest way to deploy a static application already built to Javascript and HTML. You can look at [my fork of GraphvizOnline](https://github.com/Xe/GraphvizOnline) for an example of a Heroku-compatible progressive web app.

## Using Your Progressive Web App

For iOS Safari, go to the webpage you want to add as an app, then click the share button (you may have to tap the bottom of the screen to get the share button to show up on an iPhone). Scroll the bottom part of the share sheet over to "Add to Home Screen.” The resulting dialog will let you name and change the URL starting page of the progressive web app before it gets added to the home screen. Users can then launch, manage and delete it like any other app, with no effect on any other apps on the device.

For Android with Chrome, tap on the hamburger menu in the upper right hand corner of the browser window and then tap "Add to Home screen.” This may prompt you for confirmation, then it will put the icon on your homescreen and you can launch, multitask or delete it like any other app. Unlike iOS, you cannot edit the starting URL or name of a progressive web app with Android.

After all of these steps, you will have a progressive web app. Any page or asset that the users of that progressive web app (or any browser that supports service workers) loads will seamlessly be cached for future offline access. It will be exciting to see how service workers develop in the future. I'm personally excited the most for [background sync](https://developers.google.com/web/updates/2015/12/background-sync) – I feel it could enable some fascinatingly robust experiences.

---

Also posted on the [Heroku Engineering Blog](https://blog.heroku.com/how-to-make-progressive-web-app).
