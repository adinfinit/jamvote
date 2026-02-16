/* ========================================
   JamVote Synthwave - Shared JS
   ======================================== */

// --- Dark/Light Theme Toggle ---
(function() {
  const html = document.documentElement;
  const stored = localStorage.getItem('jamvote-theme');
  if (stored) {
    html.setAttribute('data-theme', stored);
  }

  document.addEventListener('DOMContentLoaded', function() {
    const btn = document.querySelector('.theme-toggle');
    if (!btn) return;

    btn.addEventListener('click', function() {
      const current = html.getAttribute('data-theme');
      const next = current === 'light' ? 'dark' : 'light';
      html.setAttribute('data-theme', next);
      localStorage.setItem('jamvote-theme', next);
    });
  });
})();

// --- Mobile Nav Toggle ---
document.addEventListener('DOMContentLoaded', function() {
  const hamburger = document.querySelector('.nav-hamburger');
  const links = document.querySelector('.nav-links');
  if (!hamburger || !links) return;

  hamburger.addEventListener('click', function() {
    links.classList.toggle('open');
  });

  links.querySelectorAll('a').forEach(function(link) {
    link.addEventListener('click', function() {
      links.classList.remove('open');
    });
  });
});

// --- Countdown Timer ---
// Uses safe DOM creation instead of innerHTML
function initCountdown(elementId, targetDate) {
  var el = document.getElementById(elementId);
  if (!el) return;

  function createUnit(value, label) {
    var unit = document.createElement('div');
    unit.className = 'countdown-unit';

    var valDiv = document.createElement('div');
    valDiv.className = 'countdown-value';
    valDiv.textContent = String(value).padStart(2, '0');

    var labelDiv = document.createElement('div');
    labelDiv.className = 'countdown-label';
    labelDiv.textContent = label;

    unit.appendChild(valDiv);
    unit.appendChild(labelDiv);
    return unit;
  }

  function update() {
    var now = new Date().getTime();
    var diff = targetDate - now;
    if (diff < 0) diff = 0;

    var days = Math.floor(diff / (1000 * 60 * 60 * 24));
    var hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    var minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
    var seconds = Math.floor((diff % (1000 * 60)) / 1000);

    // Clear and rebuild with safe DOM methods
    while (el.firstChild) { el.removeChild(el.firstChild); }
    el.appendChild(createUnit(days, 'Days'));
    el.appendChild(createUnit(hours, 'Hours'));
    el.appendChild(createUnit(minutes, 'Minutes'));
    el.appendChild(createUnit(seconds, 'Seconds'));

    if (diff > 0) {
      requestAnimationFrame(update);
    }
  }

  update();
}

// --- Voting Sliders ---
function initSliders() {
  document.querySelectorAll('.vote-slider').forEach(function(slider) {
    var scoreEl = slider.closest('.vote-aspect').querySelector('.vote-aspect-score');
    function updateScore() {
      var val = parseFloat(slider.value);
      scoreEl.textContent = val.toFixed(1);
      var min = parseFloat(slider.min);
      var max = parseFloat(slider.max);
      var pct = ((val - min) / (max - min)) * 100;
      slider.style.background = 'linear-gradient(90deg, #ff2a6d 0%, #7b2ff7 ' + pct + '%, rgba(123,47,247,0.2) ' + pct + '%)';
    }
    slider.addEventListener('input', updateScore);
    updateScore();
  });
}

function calculateOverall() {
  var aspects = ['theme', 'enjoyment', 'aesthetics', 'innovation', 'bonus'];
  var total = 0;
  aspects.forEach(function(name) {
    var slider = document.getElementById('slider-' + name);
    if (slider) total += parseFloat(slider.value);
  });
  var overall = Math.min(5, Math.max(1, total / 4.5));
  var el = document.getElementById('overall-score');
  if (el) el.textContent = overall.toFixed(2);
}

document.addEventListener('DOMContentLoaded', function() {
  initSliders();
  document.querySelectorAll('.vote-slider').forEach(function(slider) {
    slider.addEventListener('input', calculateOverall);
  });
  calculateOverall();
});
