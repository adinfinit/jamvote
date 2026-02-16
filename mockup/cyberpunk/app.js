/* ============================================================
   JAMVOTE CYBERPUNK THEME - app.js
   ============================================================ */

document.addEventListener('DOMContentLoaded', () => {
  initThemeToggle();
  initHamburger();
  initCountdown();
  initSliders();
});

/* --- Dark/Light Theme Toggle --- */
function initThemeToggle() {
  const saved = localStorage.getItem('jv-theme');
  if (saved) {
    document.documentElement.setAttribute('data-theme', saved);
  }

  document.querySelectorAll('.theme-toggle').forEach(btn => {
    updateToggleLabel(btn);
    btn.addEventListener('click', () => {
      const current = document.documentElement.getAttribute('data-theme');
      const next = current === 'light' ? 'dark' : 'light';
      document.documentElement.setAttribute('data-theme', next);
      localStorage.setItem('jv-theme', next);
      document.querySelectorAll('.theme-toggle').forEach(updateToggleLabel);
    });
  });
}

function updateToggleLabel(btn) {
  const theme = document.documentElement.getAttribute('data-theme');
  btn.textContent = theme === 'light' ? '[SYS:DARK]' : '[SYS:LITE]';
}

/* --- Mobile Hamburger --- */
function initHamburger() {
  document.querySelectorAll('.nav-hamburger').forEach(btn => {
    btn.addEventListener('click', () => {
      const links = btn.closest('.nav').querySelector('.nav-links');
      if (links) links.classList.toggle('open');
    });
  });
}

/* --- Countdown Timer --- */
function initCountdown() {
  const el = document.getElementById('countdown');
  if (!el) return;

  const targetStr = el.getAttribute('data-target');
  if (!targetStr) return;
  const target = new Date(targetStr).getTime();

  function update() {
    const now = Date.now();
    let diff = Math.max(0, target - now);

    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    diff -= days * 1000 * 60 * 60 * 24;
    const hours = Math.floor(diff / (1000 * 60 * 60));
    diff -= hours * 1000 * 60 * 60;
    const mins = Math.floor(diff / (1000 * 60));
    diff -= mins * 1000 * 60;
    const secs = Math.floor(diff / 1000);

    setVal('cd-days', days);
    setVal('cd-hours', hours);
    setVal('cd-mins', mins);
    setVal('cd-secs', secs);
  }

  function setVal(id, v) {
    const node = document.getElementById(id);
    if (node) node.textContent = String(v).padStart(2, '0');
  }

  update();
  setInterval(update, 1000);
}

/* --- Interactive Sliders --- */
function initSliders() {
  document.querySelectorAll('.js-slider').forEach(range => {
    const group = range.closest('.slider-group');
    const valueEl = group ? group.querySelector('.slider-value') : null;
    const trackEl = group ? group.querySelector('.slider-track') : null;
    const max = parseFloat(range.max) || 5;
    const step = parseFloat(range.step) || 0.5;

    function refresh() {
      const val = parseFloat(range.value);
      if (valueEl) valueEl.textContent = val.toFixed(1);
      if (trackEl) trackEl.style.setProperty('--slider-pct', (val / max * 100) + '%');
    }

    range.addEventListener('input', refresh);
    refresh();
  });

  // Overall score calculation
  const overallEl = document.getElementById('overall-score');
  if (overallEl) {
    const sliders = document.querySelectorAll('.js-slider');
    function calcOverall() {
      let theme = 3, enjoy = 3, aesth = 3, innov = 3, bonus = 1;
      sliders.forEach(s => {
        const v = parseFloat(s.value);
        switch (s.getAttribute('data-aspect')) {
          case 'theme': theme = v; break;
          case 'enjoyment': enjoy = v; break;
          case 'aesthetics': aesth = v; break;
          case 'innovation': innov = v; break;
          case 'bonus': bonus = v; break;
        }
      });
      const overall = Math.min(5, Math.max(1, (theme + enjoy + aesth + innov + bonus) / 4.5));
      overallEl.textContent = overall.toFixed(2);
    }
    sliders.forEach(s => s.addEventListener('input', calcOverall));
    calcOverall();
  }
}
