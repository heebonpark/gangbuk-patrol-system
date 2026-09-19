// 지도위치확인 PWA 서비스워커.
// 오프라인 캐싱이 목적이 아니라(카카오맵 SDK·Supabase는 항상 인터넷이 필요) 설치
// 가능(installable) 조건을 만족시키기 위한 최소 구성 — 요청을 그대로 통과시킨다.
self.addEventListener('install', e => self.skipWaiting());
self.addEventListener('activate', e => self.clients.claim());
self.addEventListener('fetch', e => {
  e.respondWith(fetch(e.request));
});
