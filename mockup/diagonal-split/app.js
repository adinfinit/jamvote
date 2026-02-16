/* ========================================
   DIAGONAL SPLIT THEME - JamVote JS
   ======================================== */

// --- Dark/Light Theme Toggle ---
(function initTheme() {
  const saved = localStorage.getItem('jamvote-theme');
  if (saved) {
    document.documentElement.setAttribute('data-theme', saved);
  }
})();

function toggleTheme() {
  const html = document.documentElement;
  const current = html.getAttribute('data-theme');
  const next = current === 'light' ? 'dark' : 'light';
  html.setAttribute('data-theme', next);
  localStorage.setItem('jamvote-theme', next);
  updateThemeButton();
}

function updateThemeButton() {
  const btn = document.getElementById('theme-toggle');
  if (!btn) return;
  const theme = document.documentElement.getAttribute('data-theme');
  btn.textContent = theme === 'light' ? 'Dark' : 'Light';
}

// --- Mobile Nav Toggle ---
function toggleNav() {
  const links = document.querySelector('.nav-links');
  if (links) links.classList.toggle('open');
}

// --- Countdown Timer ---
function initCountdown() {
  const el = document.getElementById('countdown');
  if (!el) return;

  const target = el.getAttribute('data-target');
  if (!target) return;

  const targetDate = new Date(target).getTime();

  function buildItem(value, label) {
    const item = document.createElement('div');
    item.className = 'countdown-item';

    const numSpan = document.createElement('span');
    numSpan.className = 'number';
    numSpan.textContent = value;

    const lblSpan = document.createElement('span');
    lblSpan.className = 'label';
    lblSpan.textContent = label;

    item.appendChild(numSpan);
    item.appendChild(lblSpan);
    return item;
  }

  function update() {
    const now = Date.now();
    const diff = targetDate - now;

    // Clear existing children
    while (el.firstChild) {
      el.removeChild(el.firstChild);
    }

    if (diff <= 0) {
      el.appendChild(buildItem('!', 'Ended'));
      return;
    }

    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    const mins = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
    const secs = Math.floor((diff % (1000 * 60)) / 1000);

    el.appendChild(buildItem(String(days), 'Days'));
    el.appendChild(buildItem(String(hours).padStart(2, '0'), 'Hours'));
    el.appendChild(buildItem(String(mins).padStart(2, '0'), 'Mins'));
    el.appendChild(buildItem(String(secs).padStart(2, '0'), 'Secs'));
  }

  update();
  setInterval(update, 1000);
}

// --- Voting Sliders ---
function initSliders() {
  const sliders = document.querySelectorAll('.vote-slider');
  sliders.forEach(function(slider) {
    var valueDisplay = document.getElementById(slider.dataset.display);
    if (valueDisplay) {
      valueDisplay.textContent = parseFloat(slider.value).toFixed(1);
      slider.addEventListener('input', function() {
        valueDisplay.textContent = parseFloat(slider.value).toFixed(1);
        updateOverallScore();
      });
    }
  });
}

function updateOverallScore() {
  var themeEl = document.getElementById('slider-theme');
  var enjoymentEl = document.getElementById('slider-enjoyment');
  var aestheticsEl = document.getElementById('slider-aesthetics');
  var innovationEl = document.getElementById('slider-innovation');
  var bonusEl = document.getElementById('slider-bonus');

  var theme = parseFloat(themeEl ? themeEl.value : 0);
  var enjoyment = parseFloat(enjoymentEl ? enjoymentEl.value : 0);
  var aesthetics = parseFloat(aestheticsEl ? aestheticsEl.value : 0);
  var innovation = parseFloat(innovationEl ? innovationEl.value : 0);
  var bonus = parseFloat(bonusEl ? bonusEl.value : 0);

  var overall = Math.min(5, Math.max(1, (theme + enjoyment + aesthetics + innovation + bonus) / 4.5));

  var overallEl = document.getElementById('overall-score');
  if (overallEl) {
    overallEl.textContent = overall.toFixed(2);
  }
}

// --- Animate Score Bars on Scroll ---
function initScoreBars() {
  var bars = document.querySelectorAll('.score-bar-fill[data-width]');
  if (!bars.length) return;

  var observer = new IntersectionObserver(function(entries) {
    entries.forEach(function(entry) {
      if (entry.isIntersecting) {
        entry.target.style.width = entry.target.getAttribute('data-width');
        observer.unobserve(entry.target);
      }
    });
  }, { threshold: 0.2 });

  bars.forEach(function(bar) {
    bar.style.width = '0%';
    observer.observe(bar);
  });
}

// --- Vote Form Submit ---
function initVoteForm() {
  var form = document.getElementById('vote-form');
  if (!form) return;

  form.addEventListener('submit', function(e) {
    e.preventDefault();
    var btn = form.querySelector('button[type="submit"]');
    if (btn) {
      var inner = btn.querySelector('span');
      if (inner) inner.textContent = 'Submitted!';
      btn.disabled = true;
      btn.style.background = '#10b981';
    }
  });
}

// --- Init ---
document.addEventListener('DOMContentLoaded', function() {
  updateThemeButton();
  initCountdown();
  initSliders();
  initScoreBars();
  initVoteForm();
});
