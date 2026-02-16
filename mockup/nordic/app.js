/* ============================================
   Nordic JamVote - Shared JS
   ============================================ */

// --- Dark/Light Mode Toggle ---
(function () {
  const html = document.documentElement;
  const stored = localStorage.getItem('nordic-theme');
  if (stored) {
    html.setAttribute('data-theme', stored);
  } else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    html.setAttribute('data-theme', 'dark');
  }

  document.addEventListener('DOMContentLoaded', () => {
    const toggle = document.getElementById('theme-toggle');
    if (toggle) {
      toggle.addEventListener('click', () => {
        const current = html.getAttribute('data-theme');
        const next = current === 'dark' ? 'light' : 'dark';
        html.setAttribute('data-theme', next);
        localStorage.setItem('nordic-theme', next);
      });
    }

    // Mobile nav
    const hamburger = document.getElementById('nav-hamburger');
    const navLinks = document.getElementById('nav-links');
    if (hamburger && navLinks) {
      hamburger.addEventListener('click', () => {
        navLinks.classList.toggle('open');
      });
    }
  });
})();

// --- Countdown Timer ---
function initCountdown(targetDate, prefix) {
  function update() {
    const now = new Date().getTime();
    const distance = targetDate - now;

    if (distance < 0) {
      const els = ['days', 'hours', 'minutes', 'seconds'];
      els.forEach(id => {
        const el = document.getElementById(prefix + '-' + id);
        if (el) el.textContent = '0';
      });
      return;
    }

    const days = Math.floor(distance / (1000 * 60 * 60 * 24));
    const hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
    const seconds = Math.floor((distance % (1000 * 60)) / 1000);

    const daysEl = document.getElementById(prefix + '-days');
    const hoursEl = document.getElementById(prefix + '-hours');
    const minutesEl = document.getElementById(prefix + '-minutes');
    const secondsEl = document.getElementById(prefix + '-seconds');

    if (daysEl) daysEl.textContent = days;
    if (hoursEl) hoursEl.textContent = String(hours).padStart(2, '0');
    if (minutesEl) minutesEl.textContent = String(minutes).padStart(2, '0');
    if (secondsEl) secondsEl.textContent = String(seconds).padStart(2, '0');
  }

  update();
  setInterval(update, 1000);
}

// --- Slider Interactivity ---
function initSliders() {
  document.querySelectorAll('.voting-slider').forEach(slider => {
    const display = document.getElementById(slider.dataset.display);
    const updateSlider = () => {
      const val = parseFloat(slider.value);
      display.textContent = val.toFixed(1);
      // Update track fill via CSS gradient
      const min = parseFloat(slider.min);
      const max = parseFloat(slider.max);
      const pct = ((val - min) / (max - min)) * 100;
      slider.style.background = `linear-gradient(to right, var(--slider-fill) 0%, var(--slider-fill) ${pct}%, var(--slider-track) ${pct}%, var(--slider-track) 100%)`;
    };
    slider.addEventListener('input', updateSlider);
    updateSlider();
  });

  // Overall score calculation
  updateOverallScore();
}

function updateOverallScore() {
  const ids = ['theme', 'enjoyment', 'aesthetics', 'innovation', 'bonus'];
  const sliders = ids.map(id => document.getElementById('slider-' + id));
  const overallEl = document.getElementById('overall-score');

  if (!overallEl || sliders.some(s => !s)) return;

  function calc() {
    const theme = parseFloat(sliders[0].value);
    const enjoyment = parseFloat(sliders[1].value);
    const aesthetics = parseFloat(sliders[2].value);
    const innovation = parseFloat(sliders[3].value);
    const bonus = parseFloat(sliders[4].value);
    const overall = Math.min(5, Math.max(1, (theme + enjoyment + aesthetics + innovation + bonus) / 4.5));
    overallEl.textContent = overall.toFixed(2);
  }

  sliders.forEach(s => s.addEventListener('input', calc));
  calc();
}
