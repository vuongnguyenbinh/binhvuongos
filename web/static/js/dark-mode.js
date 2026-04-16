// Dark mode: apply before page renders to prevent flash
(function() {
  if (localStorage.getItem('dark') === 'true' ||
      (!localStorage.getItem('dark') && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    document.documentElement.classList.add('dark');
  }
})();

function toggleDarkMode() {
  var isDark = document.documentElement.classList.toggle('dark');
  localStorage.setItem('dark', isDark);
  // Update icon
  var icon = document.getElementById('dark-toggle-icon');
  if (icon) {
    icon.textContent = isDark ? '☀️' : '🌙';
  }
}
