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
