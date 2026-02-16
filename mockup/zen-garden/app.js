/* ============================================
   ZEN GARDEN — JamVote JavaScript
   ============================================ */

// --- Dark/Light Theme Toggle ---
(function initTheme() {
  const saved = localStorage.getItem('zen-theme');
  if (saved === 'dark') {
    document.documentElement.setAttribute('data-theme', 'dark');
  }
  updateToggleLabel();
})();

function toggleTheme() {
  const html = document.documentElement;
  const isDark = html.getAttribute('data-theme') === 'dark';

  if (isDark) {
    html.removeAttribute('data-theme');
    localStorage.setItem('zen-theme', 'light');
  } else {
    html.setAttribute('data-theme', 'dark');
    localStorage.setItem('zen-theme', 'dark');
  }

  updateToggleLabel();
}

function updateToggleLabel() {
  const btn = document.getElementById('themeBtn');
  if (!btn) return;

  const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
  // Small sun/moon symbol rendered as text
  btn.textContent = isDark ? '\u263C' : '\u25CF';
  btn.title = isDark ? 'Switch to light mode' : 'Switch to dark mode';
}

// Ensure toggle label is set after DOM loads
document.addEventListener('DOMContentLoaded', updateToggleLabel);
