/* ======================================================
   JamVote Steampunk Theme - Shared JavaScript
   ====================================================== */

// --- Dark/Light Mode Toggle ---
(function () {
  const html = document.documentElement;
  const stored = localStorage.getItem("jamvote-theme");

  if (stored) {
    html.setAttribute("data-theme", stored);
  }

  document.addEventListener("DOMContentLoaded", function () {
    const toggle = document.getElementById("theme-toggle");
    if (!toggle) return;

    function updateIcon() {
      const current = html.getAttribute("data-theme");
      toggle.textContent = current === "light" ? "\u2600" : "\u2699";
      toggle.title = current === "light" ? "Switch to dark mode" : "Switch to light mode";
    }

    updateIcon();

    toggle.addEventListener("click", function () {
      const current = html.getAttribute("data-theme");
      const next = current === "light" ? "dark" : "light";
      html.setAttribute("data-theme", next);
      localStorage.setItem("jamvote-theme", next);
      updateIcon();
    });
  });
})();

// --- Mobile Nav Toggle ---
document.addEventListener("DOMContentLoaded", function () {
  const hamburger = document.getElementById("nav-hamburger");
  const navLinks = document.getElementById("nav-links");

  if (hamburger && navLinks) {
    hamburger.addEventListener("click", function () {
      navLinks.classList.toggle("open");
    });

    document.addEventListener("click", function (e) {
      if (!hamburger.contains(e.target) && !navLinks.contains(e.target)) {
        navLinks.classList.remove("open");
      }
    });
  }
});

// --- Interactive Sliders ---
document.addEventListener("DOMContentLoaded", function () {
  var sliders = document.querySelectorAll('input[type="range"][data-display]');

  sliders.forEach(function (slider) {
    var displayId = slider.getAttribute("data-display");
    var display = document.getElementById(displayId);
    if (!display) return;

    function update() {
      display.textContent = parseFloat(slider.value).toFixed(1);
    }

    slider.addEventListener("input", update);
    update();
  });

  // Overall score calculation
  function updateOverall() {
    var theme = parseFloat((document.getElementById("slider-theme") || {}).value || 0);
    var enjoyment = parseFloat((document.getElementById("slider-enjoyment") || {}).value || 0);
    var aesthetics = parseFloat((document.getElementById("slider-aesthetics") || {}).value || 0);
    var innovation = parseFloat((document.getElementById("slider-innovation") || {}).value || 0);
    var bonus = parseFloat((document.getElementById("slider-bonus") || {}).value || 0);

    var overall = Math.min(5, Math.max(1, (theme + enjoyment + aesthetics + innovation + bonus) / 4.5));
    var overallEl = document.getElementById("overall-score");
    if (overallEl) {
      overallEl.textContent = overall.toFixed(2);
    }
  }

  var voteSliders = document.querySelectorAll(".voting-form input[type='range']");
  voteSliders.forEach(function (s) {
    s.addEventListener("input", updateOverall);
  });
  if (voteSliders.length > 0) updateOverall();
});

// --- Countdown Timer ---
document.addEventListener("DOMContentLoaded", function () {
  var countdownEl = document.getElementById("countdown");
  if (!countdownEl) return;

  var targetStr = countdownEl.getAttribute("data-target");
  if (!targetStr) return;
  var target = new Date(targetStr).getTime();

  function updateCountdown() {
    var now = Date.now();
    var diff = target - now;

    if (diff <= 0) {
      countdownEl.textContent = "";
      var span = document.createElement("span");
      span.style.color = "var(--accent-brass)";
      span.style.fontFamily = "Georgia, serif";
      span.style.fontSize = "1.2rem";
      span.textContent = "Time's up!";
      countdownEl.appendChild(span);
      return;
    }

    var days = Math.floor(diff / (1000 * 60 * 60 * 24));
    diff -= days * (1000 * 60 * 60 * 24);
    var hours = Math.floor(diff / (1000 * 60 * 60));
    diff -= hours * (1000 * 60 * 60);
    var minutes = Math.floor(diff / (1000 * 60));
    diff -= minutes * (1000 * 60);
    var seconds = Math.floor(diff / 1000);

    var units = [
      { value: days, label: "Days" },
      { value: hours, label: "Hours" },
      { value: minutes, label: "Minutes" },
      { value: seconds, label: "Seconds" },
    ];

    // Clear existing content
    while (countdownEl.firstChild) {
      countdownEl.removeChild(countdownEl.firstChild);
    }

    units.forEach(function (u) {
      var unitDiv = document.createElement("div");
      unitDiv.className = "countdown-unit";

      var valueSpan = document.createElement("span");
      valueSpan.className = "countdown-value";
      valueSpan.textContent = String(u.value).padStart(2, "0");

      var labelSpan = document.createElement("span");
      labelSpan.className = "countdown-label";
      labelSpan.textContent = u.label;

      unitDiv.appendChild(valueSpan);
      unitDiv.appendChild(labelSpan);
      countdownEl.appendChild(unitDiv);
    });
  }

  updateCountdown();
  setInterval(updateCountdown, 1000);
});

// --- Score Bar Animation ---
document.addEventListener("DOMContentLoaded", function () {
  var bars = document.querySelectorAll(".score-bar-fill");

  var observer = new IntersectionObserver(
    function (entries) {
      entries.forEach(function (entry) {
        if (entry.isIntersecting) {
          var targetWidth = entry.target.getAttribute("data-width");
          if (targetWidth) {
            entry.target.style.width = targetWidth;
          }
          observer.unobserve(entry.target);
        }
      });
    },
    { threshold: 0.1 }
  );

  bars.forEach(function (bar) {
    bar.style.width = "0%";
    observer.observe(bar);
  });
});

// --- Vote Form Submission (mock) ---
document.addEventListener("DOMContentLoaded", function () {
  var voteForm = document.getElementById("vote-form");
  if (!voteForm) return;

  voteForm.addEventListener("submit", function (e) {
    e.preventDefault();
    var btn = voteForm.querySelector('button[type="submit"]');
    var originalText = btn.textContent;
    btn.textContent = "Submitted!";
    btn.disabled = true;
    btn.style.background = "linear-gradient(180deg, #6b8b3d, #4a6b2d)";

    setTimeout(function () {
      btn.textContent = originalText;
      btn.disabled = false;
      btn.style.background = "";
    }, 2000);
  });
});
