(function() {
    var STORAGE_KEY = 'jamvote-theme';

    function getPreferred() {
        var saved = localStorage.getItem(STORAGE_KEY);
        if (saved === 'light' || saved === 'dark') return saved;
        return 'dark';
    }

    function apply(theme) {
        if (theme === 'light') {
            document.documentElement.setAttribute('data-theme', 'light');
        } else {
            document.documentElement.removeAttribute('data-theme');
        }
        updateButtons(theme);
    }

    function updateButtons(theme) {
        var btns = document.querySelectorAll('.theme-toggle');
        for (var i = 0; i < btns.length; i++) {
            btns[i].textContent = theme === 'light' ? '[CRT]' : '[PAPER]';
            btns[i].title = theme === 'light'
                ? 'Switch to CRT mode'
                : 'Switch to paper printout mode';
        }
    }

    function toggle() {
        var current = getPreferred();
        var next = current === 'dark' ? 'light' : 'dark';
        localStorage.setItem(STORAGE_KEY, next);
        apply(next);
    }

    // Apply saved theme immediately
    apply(getPreferred());

    // Expose toggle for onclick
    window.toggleTheme = toggle;
})();
