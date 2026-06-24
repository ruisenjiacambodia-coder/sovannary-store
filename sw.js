// ============================================================
// Sovannary Store - Service Worker
// Version: 1.0.0
// Strategy: Network-first for API, Cache-first for assets
// ============================================================

const CACHE_NAME = 'sovannary-v1.0.0';
const STATIC_ASSETS = [
  '/',
  '/index.html',
  '/dashboard.html',
  '/manifest.json',
  'https://cdn.tailwindcss.com',
  'https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css',
  'https://fonts.googleapis.com/css2?family=Battambang:wght@300;400;700;900&family=Playfair+Display:wght@400;600;700;900&family=Inter:wght@300;400;500;600;700&display=swap',
];

const API_CACHE = 'sovannary-api-v1';
const MAX_API_CACHE_AGE = 5 * 60 * 1000; // 5 minutes

// ============================================================
// Install Event
// ============================================================
self.addEventListener('install', (event) => {
  console.log('[SW] Installing...');
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      console.log('[SW] Caching static assets');
      return cache.addAll(STATIC_ASSETS).catch(err => {
        console.warn('[SW] Some assets failed to cache:', err);
      });
    })
  );
  self.skipWaiting();
});

// ============================================================
// Activate Event - Clean old caches
// ============================================================
self.addEventListener('activate', (event) => {
  console.log('[SW] Activating...');
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames
          .filter((name) => name !== CACHE_NAME && name !== API_CACHE)
          .map((name) => {
            console.log('[SW] Deleting old cache:', name);
            return caches.delete(name);
          })      );
    })
  );
  self.clients.claim();
});

// ============================================================
// Fetch Event - Routing strategy
// ============================================================
self.addEventListener('fetch', (event) => {
  const { request } = event;
  const url = new URL(request.url);

  // Skip non-GET requests
  if (request.method !== 'GET') return;

  // API requests: Network-first with cache fallback
  if (url.pathname.startsWith('/api/')) {
    event.respondWith(networkFirstWithCache(request, API_CACHE));
    return;
  }

  // Image requests: Cache-first with network fallback
  if (request.destination === 'image' || url.pathname.match(/\.(png|jpg|jpeg|gif|webp|svg)$/i)) {
    event.respondWith(cacheFirstWithNetwork(request));
    return;
  }

  // Static assets: Cache-first
  if (isStaticAsset(url.pathname)) {
    event.respondWith(cacheFirstWithNetwork(request, CACHE_NAME));
    return;
  }

  // HTML/Navigation: Network-first with offline fallback
  if (request.destination === 'document' || request.headers.get('accept')?.includes('text/html')) {
    event.respondWith(networkFirstWithCache(request, CACHE_NAME, '/index.html'));
    return;
  }

  // Default: Network with cache fallback
  event.respondWith(networkFirstWithCache(request));
});

// ============================================================
// Strategy: Network First with Cache Fallback
// ============================================================
async function networkFirstWithCache(request, cacheName = CACHE_NAME, fallbackUrl = null) {
  try {
    const networkResponse = await fetch(request);    if (networkResponse.ok) {
      const cache = await caches.open(cacheName);
      // Clone and add timestamp for API responses
      const responseToCache = networkResponse.clone();
      cache.put(request, responseToCache);
    }
    return networkResponse;
  } catch (err) {
    console.log('[SW] Network failed, trying cache for:', request.url);
    const cache = await caches.open(cacheName);
    const cachedResponse = await cache.match(request);
    if (cachedResponse) {
      return cachedResponse;
    }
    // Try fallback
    if (fallbackUrl) {
      const fallback = await cache.match(fallbackUrl);
      if (fallback) return fallback;
    }
    // Return offline page
    return new Response(
      JSON.stringify({ error: 'Offline', message: 'Please check your connection' }),
      { status: 503, headers: { 'Content-Type': 'application/json' } }
    );
  }
}

// ============================================================
// Strategy: Cache First with Network Fallback
// ============================================================
async function cacheFirstWithNetwork(request, cacheName = CACHE_NAME) {
  const cache = await caches.open(cacheName);
  const cachedResponse = await cache.match(request);
  if (cachedResponse) {
    // Stale-while-revalidate for images
    fetch(request).then((networkResponse) => {
      if (networkResponse.ok) {
        cache.put(request, networkResponse);
      }
    }).catch(() => {});
    return cachedResponse;
  }
  try {
    const networkResponse = await fetch(request);
    if (networkResponse.ok) {
      cache.put(request, networkResponse.clone());
    }
    return networkResponse;
  } catch (err) {
    return new Response('', { status: 404 });  }
}

// ============================================================
// Background Sync Event
// ============================================================
self.addEventListener('sync', (event) => {
  if (event.tag === 'sync-pending') {
    console.log('[SW] Background sync triggered');
    event.waitUntil(doBackgroundSync());
  }
});

async function doBackgroundSync() {
  try {
    // Notify all clients to trigger sync
    const clients = await self.clients.matchAll();
    clients.forEach((client) => {
      client.postMessage({ type: 'SYNC_TRIGGER' });
    });
  } catch (err) {
    console.error('[SW] Background sync failed:', err);
    throw err;
  }
}

// ============================================================
// Push Notifications (optional)
// ============================================================
self.addEventListener('push', (event) => {
  if (!event.data) return;
  const data = event.data.json();
  event.waitUntil(
    self.registration.showNotification(data.title || 'Sovannary Store', {
      body: data.body || 'You have a new notification',
      icon: '/manifest.json',
      badge: '/manifest.json',
      vibrate: [200, 100, 200],
    })
  );
});

// ============================================================
// Helpers
// ============================================================
function isStaticAsset(pathname) {
  return STATIC_ASSETS.some(asset => pathname === asset || pathname.endsWith(asset));
}