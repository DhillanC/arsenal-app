// Arsenal App — minimal service worker.
//
// Pass-through fetch handler with no caching. This makes the PWA install
// prompt eligible (a registered SW is one of the PWA criteria) without
// committing to a caching strategy before it's been designed.
//
// When offline support is actually scoped, replace this with a versioned
// cache (e.g., workbox or a hand-rolled cache-first/network-first strategy).

self.addEventListener('install', function (event) {
  self.skipWaiting();
});

self.addEventListener('activate', function (event) {
  event.waitUntil(self.clients.claim());
});

self.addEventListener('fetch', function (event) {
  // No-op: let the request go to the network as usual.
});
