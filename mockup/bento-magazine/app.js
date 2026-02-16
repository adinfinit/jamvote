/* ============================================
   BENTO MAGAZINE — App JS
   ============================================ */

// --- Dark/Light Mode Toggle ---
(function() {
  const html = document.documentElement;
  const stored = localStorage.getItem('bento-theme');
  if (stored) {
    html.setAttribute('data-theme', stored);
  }

  document.addEventListener('DOMContentLoaded', function() {
    const toggle = document.getElementById('theme-toggle');
    if (!toggle) return;

    function updateIcon() {
      const isDark = html.getAttribute('data-theme') === 'dark';
      toggle.textContent = isDark ? '\u2600' : '\u263D';
      toggle.setAttribute('aria-label', isDark ? 'Switch to light mode' : 'Switch to dark mode');
    }

    updateIcon();

    toggle.addEventListener('click', function() {
      const isDark = html.getAttribute('data-theme') === 'dark';
      const next = isDark ? 'light' : 'dark';
      html.setAttribute('data-theme', next);
      localStorage.setItem('bento-theme', next);
      updateIcon();
    });
  });
})();

// --- Mobile Nav Toggle ---
function toggleNav() {
  var links = document.querySelector('.nav-links');
  if (links) links.classList.toggle('open');
}

// --- Countdown Timer ---
function initCountdown(targetDate, elementId) {
  var el = document.getElementById(elementId);
  if (!el) return;

  function padTwo(n) {
    return String(n).padStart(2, '0');
  }

  function buildCountdownDOM(days, hours, minutes, seconds) {
    // Clear existing content
    el.textContent = '';

    var row = document.createElement('div');
    row.className = 'countdown-row';

    var units = [
      { value: days, label: 'Days' },
      { value: hours, label: 'Hours' },
      { value: minutes, label: 'Min' },
      { value: seconds, label: 'Sec' }
    ];

    units.forEach(function(unit, i) {
      if (i > 0) {
        var sep = document.createElement('div');
        sep.className = 'countdown-sep';
        sep.textContent = ':';
        row.appendChild(sep);
      }
      var wrap = document.createElement('div');
      wrap.className = 'countdown-unit';
      var num = document.createElement('div');
      num.className = 'countdown-number';
      num.textContent = padTwo(unit.value);
      var lbl = document.createElement('div');
      lbl.className = 'countdown-label';
      lbl.textContent = unit.label;
      wrap.appendChild(num);
      wrap.appendChild(lbl);
      row.appendChild(wrap);
    });

    el.appendChild(row);
  }

  function buildEndedDOM() {
    el.textContent = '';
    var num = document.createElement('span');
    num.className = 'countdown-number';
    num.textContent = '00';
    var lbl = document.createElement('span');
    lbl.className = 'countdown-label';
    lbl.textContent = 'ENDED';
    el.appendChild(num);
    el.appendChild(lbl);
  }

  function update() {
    var now = new Date().getTime();
    var target = new Date(targetDate).getTime();
    var diff = target - now;

    if (diff <= 0) {
      buildEndedDOM();
      return;
    }

    var days = Math.floor(diff / (1000 * 60 * 60 * 24));
    var hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    var minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
    var seconds = Math.floor((diff % (1000 * 60)) / 1000);

    buildCountdownDOM(days, hours, minutes, seconds);
  }

  update();
  setInterval(update, 1000);
}

// --- Voting Sliders ---
function initSliders() {
  var sliders = document.querySelectorAll('.aspect-slider');
  sliders.forEach(function(slider) {
    var display = document.getElementById(slider.dataset.display);
    if (!display) return;

    function updateDisplay() {
      var val = parseFloat(slider.value);
      display.textContent = val.toFixed(1);
      recalcOverall();
    }

    slider.addEventListener('input', updateDisplay);
    updateDisplay();
  });
}

function recalcOverall() {
  var theme = parseFloat(document.getElementById('val-theme')?.textContent || 0);
  var enjoyment = parseFloat(document.getElementById('val-enjoyment')?.textContent || 0);
  var aesthetics = parseFloat(document.getElementById('val-aesthetics')?.textContent || 0);
  var innovation = parseFloat(document.getElementById('val-innovation')?.textContent || 0);
  var bonus = parseFloat(document.getElementById('val-bonus')?.textContent || 0);

  var overall = (theme + enjoyment + aesthetics + innovation + bonus) / 4.5;
  overall = Math.min(5, Math.max(1, overall));

  var el = document.getElementById('overall-score');
  if (el) el.textContent = overall.toFixed(2);
}

// --- Score Bar Animation ---
function animateScoreBars() {
  var bars = document.querySelectorAll('.score-bar-fill');
  bars.forEach(function(bar) {
    var target = bar.dataset.width;
    if (target) {
      setTimeout(function() {
        bar.style.width = target;
      }, 200);
    }
  });
}

// --- Init on DOM Ready ---
document.addEventListener('DOMContentLoaded', function() {
  initSliders();
  animateScoreBars();

  // Initialize countdown if element exists
  // Voting closes in ~3 days from now
  var countdownEl = document.getElementById('countdown');
  if (countdownEl) {
    var target = new Date();
    target.setDate(target.getDate() + 3);
    target.setHours(18, 0, 0, 0);
    initCountdown(target.toISOString(), 'countdown');
  }
});
