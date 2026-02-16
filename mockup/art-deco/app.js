/* ============================================
   JAMVOTE ART DECO — Shared JavaScript
   ============================================ */

// --- Dark/Light Theme Toggle ---
(function () {
    const toggle = document.getElementById('theme-toggle');
    const html = document.documentElement;

    // Restore saved theme
    const saved = localStorage.getItem('jamvote-theme');
    if (saved) {
        html.setAttribute('data-theme', saved);
    }
    updateToggleIcon();

    if (toggle) {
        toggle.addEventListener('click', function () {
            const current = html.getAttribute('data-theme');
            const next = current === 'light' ? 'dark' : 'light';
            html.setAttribute('data-theme', next);
            localStorage.setItem('jamvote-theme', next);
            updateToggleIcon();
        });
    }

    function updateToggleIcon() {
        if (!toggle) return;
        const isLight = html.getAttribute('data-theme') === 'light';
        toggle.querySelector('span').textContent = isLight ? '\u263E' : '\u2600';
    }
})();

// --- Mobile Nav Toggle ---
(function () {
    const btn = document.querySelector('.nav-toggle');
    const links = document.querySelector('.nav-links');
    if (btn && links) {
        btn.addEventListener('click', function () {
            links.classList.toggle('open');
        });
    }
})();

// --- Countdown Timer ---
(function () {
    const el = document.getElementById('countdown');
    if (!el) return;

    // Target: 2 days, 14 hours from now (for demo)
    const target = new Date();
    target.setDate(target.getDate() + 2);
    target.setHours(target.getHours() + 14);

    function update() {
        const now = new Date();
        let diff = Math.max(0, Math.floor((target - now) / 1000));

        const days = Math.floor(diff / 86400);
        diff %= 86400;
        const hours = Math.floor(diff / 3600);
        diff %= 3600;
        const minutes = Math.floor(diff / 60);
        const seconds = diff % 60;

        setDigit('cd-days', days);
        setDigit('cd-hours', hours);
        setDigit('cd-minutes', minutes);
        setDigit('cd-seconds', seconds);
    }

    function setDigit(id, val) {
        const node = document.getElementById(id);
        if (node) node.textContent = String(val).padStart(2, '0');
    }

    update();
    setInterval(update, 1000);
})();

// --- Voting Sliders ---
(function () {
    const sliders = document.querySelectorAll('.vote-slider');
    sliders.forEach(function (slider) {
        const display = document.getElementById(slider.dataset.display);
        if (!display) return;

        function updateDisplay() {
            const val = parseFloat(slider.value);
            display.textContent = val.toFixed(1);
        }

        slider.addEventListener('input', function () {
            updateDisplay();
            computeOverall();
        });
        updateDisplay();
    });

    function computeOverall() {
        const theme = getVal('slider-theme');
        const enjoyment = getVal('slider-enjoyment');
        const aesthetics = getVal('slider-aesthetics');
        const innovation = getVal('slider-innovation');
        const bonus = getVal('slider-bonus');

        if (theme === null) return; // not on voting page

        const overall = Math.min(5, Math.max(1,
            (theme + enjoyment + aesthetics + innovation + bonus) / 4.5
        ));

        const el = document.getElementById('overall-value');
        if (el) el.textContent = overall.toFixed(2);
    }

    function getVal(id) {
        const el = document.getElementById(id);
        return el ? parseFloat(el.value) : null;
    }

    computeOverall();
})();

// --- Animate score bars on results/team pages ---
(function () {
    const bars = document.querySelectorAll('.score-bar-fill');
    if (!bars.length) return;

    const observer = new IntersectionObserver(function (entries) {
        entries.forEach(function (entry) {
            if (entry.isIntersecting) {
                const bar = entry.target;
                bar.style.width = bar.dataset.width;
                observer.unobserve(bar);
            }
        });
    }, { threshold: 0.2 });

    bars.forEach(function (bar) {
        bar.dataset.width = bar.style.width;
        bar.style.width = '0%';
        observer.observe(bar);
    });
})();
