/* ============================================================
   Circular Orbit — JamVote App JS
   ============================================================ */

// --- Dark/Light Theme Toggle ---
(function initTheme() {
  const stored = localStorage.getItem('jamvote-theme');
  if (stored) {
    document.documentElement.setAttribute('data-theme', stored);
  }
})();

function toggleTheme() {
  const html = document.documentElement;
  const current = html.getAttribute('data-theme');
  const next = current === 'light' ? 'dark' : 'light';
  html.setAttribute('data-theme', next);
  localStorage.setItem('jamvote-theme', next);
}

// --- SVG Ring Chart Generator ---
// Note: All values passed to this function are hardcoded in the mockup,
// not derived from user input — safe for innerHTML use in this static demo.
function createRingChart(container, value, max, size, strokeWidth, color, labelHtml) {
  const radius = (size - strokeWidth) / 2;
  const circumference = 2 * Math.PI * radius;
  const progress = (value / max) * circumference;
  const cx = size / 2;
  const cy = size / 2;

  const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
  svg.setAttribute('width', size);
  svg.setAttribute('height', size);
  svg.setAttribute('viewBox', '0 0 ' + size + ' ' + size);
  svg.style.transform = 'rotate(-90deg)';

  // Background ring
  const bgCircle = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
  bgCircle.setAttribute('cx', cx);
  bgCircle.setAttribute('cy', cy);
  bgCircle.setAttribute('r', radius);
  bgCircle.setAttribute('fill', 'none');
  bgCircle.setAttribute('stroke', 'var(--ring-bg)');
  bgCircle.setAttribute('stroke-width', strokeWidth);
  svg.appendChild(bgCircle);

  // Progress ring
  const progressCircle = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
  progressCircle.setAttribute('cx', cx);
  progressCircle.setAttribute('cy', cy);
  progressCircle.setAttribute('r', radius);
  progressCircle.setAttribute('fill', 'none');
  progressCircle.setAttribute('stroke', color);
  progressCircle.setAttribute('stroke-width', strokeWidth);
  progressCircle.setAttribute('stroke-linecap', 'round');
  progressCircle.setAttribute('stroke-dasharray', circumference);
  progressCircle.setAttribute('stroke-dashoffset', circumference);
  progressCircle.style.transition = 'stroke-dashoffset 1s ease-out';
  progressCircle.style.filter = 'drop-shadow(0 0 4px ' + color + ')';
  svg.appendChild(progressCircle);

  while (container.firstChild) container.removeChild(container.firstChild);
  container.style.position = 'relative';
  container.style.display = 'inline-flex';
  container.style.alignItems = 'center';
  container.style.justifyContent = 'center';
  container.appendChild(svg);

  // Label overlay (static mockup content only, not user-generated)
  if (labelHtml !== undefined) {
    const labelEl = document.createElement('div');
    labelEl.className = 'ring-chart-label';
    labelEl.textContent = '';
    // For this static mockup, we build label content safely via DOM
    setRingLabel(labelEl, labelHtml);
    container.appendChild(labelEl);
  }

  // Animate in
  requestAnimationFrame(function() {
    requestAnimationFrame(function() {
      progressCircle.setAttribute('stroke-dashoffset', circumference - progress);
    });
  });

  return { svg: svg, progressCircle: progressCircle, circumference: circumference, radius: radius };
}

// Safe label setter for ring charts (mockup only uses known static values)
function setRingLabel(el, content) {
  if (typeof content === 'object' && content.value !== undefined) {
    el.textContent = '';
    var valSpan = document.createElement('span');
    valSpan.style.fontSize = content.valueFontSize || '1.2rem';
    valSpan.textContent = content.value;
    el.appendChild(valSpan);
    if (content.sub) {
      var subEl = document.createElement('small');
      subEl.textContent = content.sub;
      el.appendChild(subEl);
    }
  } else {
    el.textContent = String(content);
  }
}

// --- Update ring dynamically ---
function updateRing(progressCircle, circumference, value, max) {
  var progress = (value / max) * circumference;
  progressCircle.setAttribute('stroke-dashoffset', circumference - progress);
}

// --- Countdown Timer ---
function startCountdown(targetDate, elementId) {
  var el = document.getElementById(elementId);
  if (!el) return;

  function update() {
    var now = new Date().getTime();
    var distance = targetDate - now;

    if (distance <= 0) {
      el.textContent = '00:00:00';
      return;
    }

    var days = Math.floor(distance / (1000 * 60 * 60 * 24));
    var hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    var minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
    var seconds = Math.floor((distance % (1000 * 60)) / 1000);

    var text = '';
    if (days > 0) text += days + 'd ';
    text += String(hours).padStart(2, '0') + ':' + String(minutes).padStart(2, '0') + ':' + String(seconds).padStart(2, '0');
    el.textContent = text;

    // Also update ring if present
    var ring = document.getElementById(elementId + '-ring');
    if (ring && ring._ringData) {
      var total = ring._ringData.totalMs;
      var remaining = Math.max(0, distance);
      var elapsed = total - remaining;
      updateRing(ring._ringData.progressCircle, ring._ringData.circumference, elapsed, total);
    }
  }

  update();
  setInterval(update, 1000);
}

// --- Position items in orbit on desktop ---
function positionOrbitItems(containerId, centerSelector, items, radius) {
  if (window.innerWidth < 900) return;

  var container = document.getElementById(containerId);
  if (!container) return;
  var center = container.querySelector(centerSelector);
  if (!center) return;

  var containerRect = container.getBoundingClientRect();
  var centerRect = center.getBoundingClientRect();
  var cx = centerRect.left - containerRect.left + centerRect.width / 2;
  var cy = centerRect.top - containerRect.top + centerRect.height / 2;

  var angleStep = (2 * Math.PI) / items.length;
  var startAngle = -Math.PI / 2;

  items.forEach(function(item, i) {
    var angle = startAngle + angleStep * i;
    var x = cx + radius * Math.cos(angle);
    var y = cy + radius * Math.sin(angle);
    item.style.position = 'absolute';
    item.style.left = x + 'px';
    item.style.top = y + 'px';
    item.style.transform = 'translate(-50%, -50%)';
  });
}

// --- Slider value display + overall calculation ---
function initVotingSliders() {
  var sliders = document.querySelectorAll('.vote-slider');
  var overallEl = document.getElementById('overall-score');
  var overallRingContainer = document.getElementById('overall-ring');

  var overallRingData = null;

  if (overallRingContainer) {
    overallRingData = createRingChart(overallRingContainer, 0, 5, 160, 12, '#60a5fa',
      { value: '0.0', valueFontSize: '2rem', sub: 'Overall' });
  }

  function recalcOverall() {
    var theme = 0, enjoyment = 0, aesthetics = 0, innovation = 0, bonus = 0;
    sliders.forEach(function(s) {
      var val = parseFloat(s.value);
      var name = s.getAttribute('data-aspect');
      if (name === 'theme') theme = val;
      if (name === 'enjoyment') enjoyment = val;
      if (name === 'aesthetics') aesthetics = val;
      if (name === 'innovation') innovation = val;
      if (name === 'bonus') bonus = val;
    });

    var overall = Math.min(5, Math.max(1, (theme + enjoyment + aesthetics + innovation + bonus) / 4.5));

    if (overallEl) overallEl.textContent = overall.toFixed(1);
    if (overallRingData && overallRingContainer) {
      updateRing(overallRingData.progressCircle, overallRingData.circumference, overall, 5);
      var labelEl = overallRingContainer.querySelector('.ring-chart-label');
      if (labelEl) {
        setRingLabel(labelEl, { value: overall.toFixed(1), valueFontSize: '2rem', sub: 'Overall' });
      }
    }
  }

  sliders.forEach(function(slider) {
    var display = document.getElementById(slider.id + '-val');
    var miniRing = document.getElementById(slider.id + '-ring');

    var miniRingData = null;
    var max = slider.getAttribute('data-aspect') === 'bonus' ? 2.5 : 5;
    var color = slider.getAttribute('data-color') || '#60a5fa';

    if (miniRing) {
      miniRingData = createRingChart(miniRing, parseFloat(slider.value), max, 50, 5, color,
        { value: parseFloat(slider.value).toFixed(1), valueFontSize: '0.85rem' });
    }

    slider.addEventListener('input', function() {
      var val = parseFloat(slider.value);
      if (display) display.textContent = val.toFixed(1);
      if (miniRingData && miniRing) {
        updateRing(miniRingData.progressCircle, miniRingData.circumference, val, max);
        var label = miniRing.querySelector('.ring-chart-label');
        if (label) {
          setRingLabel(label, { value: val.toFixed(1), valueFontSize: '0.85rem' });
        }
      }
      recalcOverall();
    });
  });

  // Initial calculation
  recalcOverall();
}

// --- Initialize Score Ring Charts on results/team pages ---
function initScoreRings() {
  document.querySelectorAll('[data-ring]').forEach(function(el) {
    var value = parseFloat(el.getAttribute('data-value'));
    var max = parseFloat(el.getAttribute('data-max') || '5');
    var size = parseInt(el.getAttribute('data-size') || '80');
    var stroke = parseInt(el.getAttribute('data-stroke') || '6');
    var color = el.getAttribute('data-color') || '#60a5fa';
    var label = el.getAttribute('data-label') || '';
    var fontSize = size > 100 ? '1.5rem' : '0.75rem';
    createRingChart(el, value, max, size, stroke, color,
      { value: value.toFixed(1), valueFontSize: fontSize, sub: label || undefined });
  });
}

// --- DOM Ready ---
document.addEventListener('DOMContentLoaded', function() {
  // Theme toggle
  var toggleBtn = document.getElementById('theme-toggle');
  if (toggleBtn) toggleBtn.addEventListener('click', toggleTheme);

  // Initialize score rings
  initScoreRings();

  // Initialize voting sliders if present
  if (document.querySelector('.vote-slider')) {
    initVotingSliders();
  }

  // Orbit layouts on desktop
  if (window.innerWidth >= 900) {
    var innerItems = document.querySelectorAll('.orbit-inner-item');
    var outerItems = document.querySelectorAll('.orbit-outer-item');

    if (innerItems.length > 0) {
      positionOrbitItems('orbit-main', '.central-circle', Array.from(innerItems), 200);
    }
    if (outerItems.length > 0) {
      positionOrbitItems('orbit-main', '.central-circle', Array.from(outerItems), 320);
    }
  }

  // Countdown
  var countdownEl = document.getElementById('countdown-text');
  if (countdownEl) {
    var target = new Date(countdownEl.getAttribute('data-target'));
    startCountdown(target.getTime(), 'countdown-text');
  }
});
