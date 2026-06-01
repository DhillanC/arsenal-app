// Arsenal App — theme toggle + PWA service worker registration.
//
// base.html calls toggleTheme() onclick of the theme button. Without this file
// the button is wired to an undefined function and dark/light mode is broken.

(function () {
  const STORAGE_KEY = 'arsenal.theme';
  const root = document.documentElement;

  // Apply persisted preference on load. If none, honor OS preference.
  function applyInitialTheme() {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved === 'dark') {
      root.classList.add('dark');
    } else if (saved === 'light') {
      root.classList.remove('dark');
    } else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
      root.classList.add('dark');
    }
  }

  applyInitialTheme();

  // Exposed globally because base.html uses `onclick="toggleTheme()"`.
  window.toggleTheme = function toggleTheme() {
    const isDark = root.classList.toggle('dark');
    const theme = isDark ? 'dark' : 'light';
    localStorage.setItem(STORAGE_KEY, theme);
    // Sync cookie for server-side rendering on next request
    document.cookie = 'arsenal_theme=' + theme + ';path=/;max-age=31536000';
  };
})();

// PWA service worker registration. The commit promised "Service worker structure";
// this is the actual hookup. Registration is best-effort and silently no-ops if
// the browser doesn't support it (e.g., older Safari, http:// non-localhost).
if ('serviceWorker' in navigator) {
  window.addEventListener('load', function () {
    navigator.serviceWorker.register('/static/js/sw.js').catch(function () {
      // Silent — SW is progressive enhancement, not load-bearing.
    });
  });
}
