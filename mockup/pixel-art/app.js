/* ========================================
   JamVote Pixel Art - Shared JS
   ======================================== */

// --- Dark/Light Theme Toggle ---
(function() {
  const html = document.documentElement;
  const stored = localStorage.getItem('jamvote-theme');
  if (stored) {
    html.setAttribute('data-theme', stored);
  }
  // default is dark (no attribute needed, but be explicit)
  if (!stored) {
    html.setAttribute('data-theme', 'dark');
  }
})();

function toggleTheme() {
  const html = document.documentElement;
  const current = html.getAttribute('data-theme') || 'dark';
  const next = current === 'dark' ? 'light' : 'dark';
  html.setAttribute('data-theme', next);
  localStorage.setItem('jamvote-theme', next);
  updateThemeButton();
}

function updateThemeButton() {
  const btn = document.getElementById('theme-btn');
  if (!btn) return;
  const current = document.documentElement.getAttribute('data-theme') || 'dark';
  btn.textContent = current === 'dark' ? '[DAY]' : '[NITE]';
}

document.addEventListener('DOMContentLoaded', updateThemeButton);

// --- Mobile Nav Toggle ---
function toggleNav() {
  const links = document.getElementById('nav-links');
  if (links) links.classList.toggle('open');
}

// --- Countdown Timer ---
function initCountdown(elementId, targetDate) {
  const el = document.getElementById(elementId);
  if (!el) return;

  function update() {
    const now = new Date().getTime();
    const target = new Date(targetDate).getTime();
    let diff = target - now;

    if (diff <= 0) {
      el.querySelector('.cd-days').textContent = '00';
      el.querySelector('.cd-hrs').textContent = '00';
      el.querySelector('.cd-min').textContent = '00';
      el.querySelector('.cd-sec').textContent = '00';
      return;
    }

    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    diff -= days * 1000 * 60 * 60 * 24;
    const hrs = Math.floor(diff / (1000 * 60 * 60));
    diff -= hrs * 1000 * 60 * 60;
    const min = Math.floor(diff / (1000 * 60));
    diff -= min * 1000 * 60;
    const sec = Math.floor(diff / 1000);

    el.querySelector('.cd-days').textContent = String(days).padStart(2, '0');
    el.querySelector('.cd-hrs').textContent = String(hrs).padStart(2, '0');
    el.querySelector('.cd-min').textContent = String(min).padStart(2, '0');
    el.querySelector('.cd-sec').textContent = String(sec).padStart(2, '0');
  }

  update();
  setInterval(update, 1000);
}

// --- Interactive Score Slider ---
function initScoreSlider(name) {
  const slider = document.getElementById('slider-' + name);
  const display = document.getElementById('display-' + name);
  const segments = document.querySelectorAll('#bar-' + name + ' .seg');

  if (!slider || !display) return;

  function updateBar(val) {
    display.textContent = parseFloat(val).toFixed(1);
    const max = parseFloat(slider.max);
    const step = max / segments.length;

    segments.forEach(function(seg, i) {
      if ((i + 1) * step <= parseFloat(val) + step * 0.01) {
        seg.classList.add('filled');
      } else {
        seg.classList.remove('filled');
      }
    });
  }

  slider.addEventListener('input', function() {
    updateBar(this.value);
    updateOverall();
  });

  // clickable segments
  segments.forEach(function(seg, i) {
    seg.addEventListener('click', function() {
      const max = parseFloat(slider.max);
      const step = max / segments.length;
      const val = (i + 1) * step;
      slider.value = val;
      updateBar(val);
      updateOverall();
    });
  });

  updateBar(slider.value);
}

function updateOverall() {
  const aspects = ['theme', 'enjoyment', 'aesthetics', 'innovation', 'bonus'];
  let sum = 0;
  aspects.forEach(function(name) {
    const slider = document.getElementById('slider-' + name);
    if (slider) sum += parseFloat(slider.value);
  });
  const overall = Math.min(5, Math.max(1, sum / 4.5));
  const el = document.getElementById('overall-score');
  if (el) el.textContent = overall.toFixed(2);
}

// --- Init all voting sliders ---
function initVotingPage() {
  ['theme', 'enjoyment', 'aesthetics', 'innovation', 'bonus'].forEach(initScoreSlider);
  updateOverall();
}

// --- Fake form submit ---
function showPixelAlert(msg) {
  alert(msg);
}
