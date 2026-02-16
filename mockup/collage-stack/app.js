/* ============================================
   COLLAGE STACK — JamVote App JavaScript
   ============================================ */

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
  const next = current === 'dark' ? 'light' : 'dark';
  html.setAttribute('data-theme', next);
  localStorage.setItem('jamvote-theme', next);
  updateThemeIcon();
}

function updateThemeIcon() {
  const btn = document.querySelector('.theme-toggle');
  if (!btn) return;
  const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
  btn.textContent = isDark ? '\u2600' : '\u263E';
  btn.title = isDark ? 'Switch to light mode' : 'Switch to dark mode';
}

// --- Mobile Nav Toggle ---
function toggleNav() {
  const links = document.querySelector('.nav-links');
  if (links) links.classList.toggle('open');
}

// --- Countdown Timer ---
function initCountdown(elementId, targetDate) {
  var el = document.getElementById(elementId);
  if (!el) return;

  function update() {
    var now = new Date().getTime();
    var target = new Date(targetDate).getTime();
    var diff = target - now;

    if (diff <= 0) {
      el.textContent = '';
      var span = document.createElement('span');
      span.style.color = 'var(--tomato)';
      span.style.fontWeight = 'bold';
      span.textContent = "Time's up!";
      el.appendChild(span);
      return;
    }

    var days = Math.floor(diff / (1000 * 60 * 60 * 24));
    var hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    var minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
    var seconds = Math.floor((diff % (1000 * 60)) / 1000);

    // Clear and rebuild using DOM methods
    el.textContent = '';

    var units = [
      { value: String(days), label: 'Days' },
      { value: String(hours).padStart(2, '0'), label: 'Hrs' },
      { value: String(minutes).padStart(2, '0'), label: 'Min' },
      { value: String(seconds).padStart(2, '0'), label: 'Sec' }
    ];

    units.forEach(function(u) {
      var unitDiv = document.createElement('div');
      unitDiv.className = 'countdown-unit';

      var numDiv = document.createElement('div');
      numDiv.className = 'countdown-number';
      numDiv.textContent = u.value;

      var labelDiv = document.createElement('div');
      labelDiv.className = 'countdown-label';
      labelDiv.textContent = u.label;

      unitDiv.appendChild(numDiv);
      unitDiv.appendChild(labelDiv);
      el.appendChild(unitDiv);
    });
  }

  update();
  setInterval(update, 1000);
}

// --- Voting Sliders ---
function initSliders() {
  var sliders = document.querySelectorAll('.aspect-slider');
  sliders.forEach(function(slider) {
    var container = slider.closest('.slider-container');
    var valueDisplay = container ? container.querySelector('.slider-value') : null;
    if (!valueDisplay) return;

    function updateValue() {
      var val = parseFloat(slider.value);
      valueDisplay.textContent = val.toFixed(1);
      updateOverallScore();
    }

    slider.addEventListener('input', updateValue);
    updateValue();
  });
}

function updateOverallScore() {
  var themeEl = document.getElementById('slider-theme');
  var enjoymentEl = document.getElementById('slider-enjoyment');
  var aestheticsEl = document.getElementById('slider-aesthetics');
  var innovationEl = document.getElementById('slider-innovation');
  var bonusEl = document.getElementById('slider-bonus');

  var theme = themeEl ? parseFloat(themeEl.value) : 0;
  var enjoyment = enjoymentEl ? parseFloat(enjoymentEl.value) : 0;
  var aesthetics = aestheticsEl ? parseFloat(aestheticsEl.value) : 0;
  var innovation = innovationEl ? parseFloat(innovationEl.value) : 0;
  var bonus = bonusEl ? parseFloat(bonusEl.value) : 0;

  var overall = Math.min(5, Math.max(1, (theme + enjoyment + aesthetics + innovation + bonus) / 4.5));

  var scoreEl = document.getElementById('overall-score-value');
  if (scoreEl) {
    scoreEl.textContent = overall.toFixed(2);
  }
}

// --- Random rotation helper for dynamic content ---
function applyScatterRotations() {
  var rotations = [-2, 1.5, -1, 2.5, -0.5, 1, -3, 0.5];
  document.querySelectorAll('.scatter-rotate').forEach(function(el, i) {
    var rot = rotations[i % rotations.length];
    el.style.transform = 'rotate(' + rot + 'deg)';
  });
}

// --- Initialize on DOM ready ---
document.addEventListener('DOMContentLoaded', function () {
  updateThemeIcon();
  initSliders();
  applyScatterRotations();

  // Mark active nav link
  var currentPage = window.location.pathname.split('/').pop() || 'index.html';
  document.querySelectorAll('.nav-links a, .subnav a').forEach(function(link) {
    var href = link.getAttribute('href');
    if (href === currentPage) {
      link.classList.add('active');
    }
  });
});
