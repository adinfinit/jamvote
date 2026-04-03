// Theme toggle for Neon Arcade - dark/light mode
(function() {
  var STORAGE_KEY = 'jamvote-theme';

  function getPreferred() {
    var saved = localStorage.getItem(STORAGE_KEY);
    if (saved === 'light' || saved === 'dark') return saved;
    return 'dark'; // default for neon arcade
  }

  function applyTheme(theme, animate) {
    if (animate) {
      document.documentElement.style.transition = 'background-color 0.3s, color 0.3s';
      setTimeout(function() {
        document.documentElement.style.transition = '';
      }, 350);
    }

    if (theme === 'light') {
      document.documentElement.setAttribute('data-theme', 'light');
    } else {
      document.documentElement.removeAttribute('data-theme');
    }
    updateToggleIcon(theme);
  }

  function updateToggleIcon(theme) {
    var btn = document.getElementById('theme-toggle');
    if (!btn) return;
    // Sun for dark mode (click to go light), Moon for light mode (click to go dark)
    btn.textContent = theme === 'light' ? '\u263E' : '\u2600';
    btn.title = theme === 'light' ? 'Switch to dark mode' : 'Switch to light mode';
  }

  function toggle() {
    var current = document.documentElement.getAttribute('data-theme') === 'light' ? 'light' : 'dark';
    var next = current === 'light' ? 'dark' : 'light';
    localStorage.setItem(STORAGE_KEY, next);
    applyTheme(next, true);
  }

  // Apply saved theme immediately (before DOM ready) to prevent flash
  applyTheme(getPreferred(), false);

  // Once DOM is ready, wire up the button
  document.addEventListener('DOMContentLoaded', function() {
    var btn = document.getElementById('theme-toggle');
    if (btn) {
      btn.addEventListener('click', toggle);
      updateToggleIcon(getPreferred());
    }
  });
})();
