//This is the service worker with the combined offline experience (Offline page + Offline copy of pages)

//Install stage sets up the offline page in the cache and opens a new cache
self.addEventListener('install', function(event) {
  event.waitUntil(preLoad());
});

const cacheName = "cache-2019-11-01";

var preLoad = function(){
  console.log('[PWA Builder] Install Event processing');
  return caches.open(cacheName).then(function(cache) {
    console.log('[PWA Builder] Cached index and offline page during Install');
    return cache.addAll(['/blog/', '/blog', '/', '/contact', '/resume', '/talks', '/gallery']);
  });
};

self.addEventListener('fetch', function(event) {
  if (event.request.cache === 'only-if-cached' && event.request.mode !== 'same-origin') {
    return;
  }
  console.log('[PWA Builder] The service worker is serving the asset.');
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
  return caches.open(cacheName).then(function (cache) {
    return fetch(request).then(function (response) {
      console.log('[PWA Builder] add page to offline: ' + response.url);
      return cache.put(request, response);
    });
  });
};

var returnFromCache = function(request){
  return caches.open(cacheName).then(function (cache) {
    return cache.match(request).then(function (matching) {
     if(!matching || matching.status == 404) {
       return cache.match('offline.html');
     } else {
       return matching;
     }
    });
  });
};
