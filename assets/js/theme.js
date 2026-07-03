(function () {
  var root = document.documentElement;
  var storageKey = 'agentrc-theme';
  var stored = null;
  try { stored = localStorage.getItem(storageKey); } catch (e) {}

  function systemTheme() {
    return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  }

  function apply(theme) {
    var next = theme === 'light' || theme === 'dark' ? theme : systemTheme();
    root.setAttribute('data-theme', next);
    var button = document.querySelector('[data-theme-toggle]');
    if (button) {
      button.setAttribute('aria-label', next === 'dark' ? 'Switch to light theme' : 'Switch to dark theme');
      button.setAttribute('title', next === 'dark' ? 'Switch to light theme' : 'Switch to dark theme');
    }
  }

  apply(stored);

  document.addEventListener('DOMContentLoaded', function () {
    apply(stored);
    var button = document.querySelector('[data-theme-toggle]');
    if (!button) return;
    button.addEventListener('click', function () {
      var current = root.getAttribute('data-theme') || systemTheme();
      var next = current === 'dark' ? 'light' : 'dark';
      try { localStorage.setItem(storageKey, next); } catch (e) {}
      apply(next);
    });
  });

  document.addEventListener('click', function (e) {
    var btn = e.target.closest ? e.target.closest('[data-copy-md]') : null;
    if (!btn) return;
    var url = btn.getAttribute('data-md-url');
    if (!url) return;
    var label = btn.querySelector('span');
    var original = label ? label.textContent : '';
    fetch(url)
      .then(function (r) { if (!r.ok) throw new Error('fetch failed'); return r.text(); })
      .then(function (text) {
        if (navigator.clipboard && navigator.clipboard.writeText) return navigator.clipboard.writeText(text);
        throw new Error('clipboard unavailable');
      })
      .then(function () {
        btn.classList.add('is-copied');
        if (label) label.textContent = 'Copied!';
        setTimeout(function () {
          btn.classList.remove('is-copied');
          if (label) label.textContent = original;
        }, 1600);
      })
      .catch(function () { window.open(url, '_blank', 'noopener'); });
  });

  // Copy a literal string (e.g. the install one-liner) from [data-copy].
  document.addEventListener('click', function (e) {
    var btn = e.target.closest ? e.target.closest('[data-copy]') : null;
    if (!btn) return;
    var text = btn.getAttribute('data-copy');
    if (!text || !navigator.clipboard || !navigator.clipboard.writeText) return;
    navigator.clipboard.writeText(text).then(function () {
      btn.classList.add('is-copied');
      setTimeout(function () { btn.classList.remove('is-copied'); }, 1400);
    }).catch(function () {});
  });

  if (window.matchMedia) {
    var mq = window.matchMedia('(prefers-color-scheme: dark)');
    var listener = function () {
      try { stored = localStorage.getItem(storageKey); } catch (e) {}
      if (!stored) apply(null);
    };
    if (mq.addEventListener) mq.addEventListener('change', listener);
    else if (mq.addListener) mq.addListener(listener);
  }
})();
