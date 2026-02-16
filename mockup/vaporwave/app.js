/* ============================================
   JAMVOTE VAPORWAVE - Shared JavaScript
   ============================================ */

// --- Dark/Light Mode Toggle ---
(function initTheme() {
  const saved = localStorage.getItem('jamvote-theme');
  if (saved === 'dark') {
    document.documentElement.setAttribute('data-theme', 'dark');
  }
})();

function toggleTheme() {
  const html = document.documentElement;
  const isDark = html.getAttribute('data-theme') === 'dark';
  if (isDark) {
    html.removeAttribute('data-theme');
    localStorage.setItem('jamvote-theme', 'light');
  } else {
    html.setAttribute('data-theme', 'dark');
    localStorage.setItem('jamvote-theme', 'dark');
  }
}

// --- Mobile Nav Toggle ---
function toggleNav() {
  const links = document.querySelector('.nav-links');
  if (links) links.classList.toggle('open');
}

// --- Countdown Timer ---
function initCountdown(elementId, targetDate) {
  const el = document.getElementById(elementId);
  if (!el) return;

  function update() {
    const now = new Date().getTime();
    const diff = targetDate - now;

    if (diff <= 0) {
      el.textContent = '';
      const span = document.createElement('span');
      span.className = 'gradient-text';
      span.style.fontSize = '1.2rem';
      span.style.fontWeight = '700';
      span.textContent = "Time's up!";
      el.appendChild(span);
      return;
    }

    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    const mins = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
    const secs = Math.floor((diff % (1000 * 60)) / 1000);

    // Clear existing content
    el.textContent = '';

    const units = [
      { value: days, label: 'Days' },
      { value: hours, label: 'Hours' },
      { value: mins, label: 'Minutes' },
      { value: secs, label: 'Seconds' }
    ];

    units.forEach(function(unit) {
      var item = document.createElement('div');
      item.className = 'countdown-item';

      var num = document.createElement('div');
      num.className = 'countdown-number';
      num.textContent = unit.label === 'Days' ? unit.value : String(unit.value).padStart(2, '0');
      item.appendChild(num);

      var lbl = document.createElement('div');
      lbl.className = 'countdown-label';
      lbl.textContent = unit.label;
      item.appendChild(lbl);

      el.appendChild(item);
    });
  }

  update();
  setInterval(update, 1000);
}

// --- Voting Sliders ---
function initSliders() {
  var sliders = document.querySelectorAll('input[type="range"][data-display]');
  sliders.forEach(function(slider) {
    var display = document.getElementById(slider.dataset.display);
    if (!display) return;

    function updateDisplay() {
      var val = parseFloat(slider.value);
      display.textContent = val.toFixed(1);

      // Update slider track fill via CSS gradient
      var min = parseFloat(slider.min);
      var max = parseFloat(slider.max);
      var pct = ((val - min) / (max - min)) * 100;
      var color = slider.dataset.color || '#b967ff';
      slider.style.background = 'linear-gradient(90deg, ' + color + ' ' + pct + '%, rgba(200,200,200,0.3) ' + pct + '%)';
    }

    slider.addEventListener('input', updateDisplay);
    updateDisplay();
  });

  updateOverallScore();
}

function updateOverallScore() {
  var themeEl = document.getElementById('slider-theme');
  var enjoyEl = document.getElementById('slider-enjoyment');
  var aesthEl = document.getElementById('slider-aesthetics');
  var innovEl = document.getElementById('slider-innovation');
  var bonusEl = document.getElementById('slider-bonus');

  if (!themeEl) return;

  var themeVal = parseFloat(themeEl.value || 0);
  var enjoyVal = parseFloat(enjoyEl.value || 0);
  var aesthVal = parseFloat(aesthEl.value || 0);
  var innovVal = parseFloat(innovEl.value || 0);
  var bonusVal = parseFloat(bonusEl.value || 0);

  var overall = Math.min(5, Math.max(1, (themeVal + enjoyVal + aesthVal + innovVal + bonusVal) / 4.5));

  var overallEl = document.getElementById('overall-score');
  if (overallEl) {
    overallEl.textContent = overall.toFixed(2);
  }
}

// --- Perspective Grid SVG ---
function initGrid() {
  var svg = document.querySelector('.grid-svg');
  if (!svg) return;

  var w = 1200;
  var h = 180;

  // Remove existing children
  while (svg.firstChild) {
    svg.removeChild(svg.firstChild);
  }

  var ns = 'http://www.w3.org/2000/svg';

  // Horizontal lines with perspective
  for (var i = 0; i < 12; i++) {
    var y = h - (Math.pow(i / 11, 1.8) * h * 0.85);
    var opacity = 0.2 + (i / 11) * 0.6;
    var line = document.createElementNS(ns, 'line');
    line.setAttribute('x1', '0');
    line.setAttribute('y1', y);
    line.setAttribute('x2', w);
    line.setAttribute('y2', y);
    line.setAttribute('stroke', 'var(--grid-line)');
    line.setAttribute('stroke-width', '1');
    line.setAttribute('opacity', opacity);
    svg.appendChild(line);
  }

  // Vertical lines converging to vanishing point
  var vanishX = w / 2;
  var vanishY = 10;
  for (var j = 0; j <= 20; j++) {
    var x = (j / 20) * w;
    var op = 0.15 + Math.abs(j - 10) / 10 * 0.35;
    var vline = document.createElementNS(ns, 'line');
    vline.setAttribute('x1', x);
    vline.setAttribute('y1', h);
    vline.setAttribute('x2', vanishX + (x - vanishX) * 0.1);
    vline.setAttribute('y2', vanishY);
    vline.setAttribute('stroke', 'var(--grid-line)');
    vline.setAttribute('stroke-width', '1');
    vline.setAttribute('opacity', op);
    svg.appendChild(vline);
  }
}

// --- Submit Vote (mock) ---
function submitVote() {
  var btn = document.querySelector('.submit-vote-btn');
  if (btn) {
    btn.textContent = 'Submitted!';
    btn.style.background = 'linear-gradient(135deg, var(--mint), var(--baby-blue))';
    btn.style.color = '#1a0a2e';
    setTimeout(function() {
      btn.textContent = 'Submit Vote';
      btn.style.background = '';
      btn.style.color = '';
    }, 2000);
  }
}

// --- Init on DOM ready ---
document.addEventListener('DOMContentLoaded', function() {
  initSliders();
  initGrid();

  // Attach overall recalculation to all voting sliders
  var votingSliders = document.querySelectorAll('.voting-slider');
  votingSliders.forEach(function(s) {
    s.addEventListener('input', updateOverallScore);
  });

  // Highlight active nav link
  var currentPage = window.location.pathname.split('/').pop() || 'index.html';
  document.querySelectorAll('.nav-links a, .sub-nav a').forEach(function(a) {
    var href = a.getAttribute('href');
    if (href === currentPage) {
      a.classList.add('active');
    }
  });
});
