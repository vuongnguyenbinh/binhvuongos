// Checkbox toggle
document.addEventListener('click', function(e) {
  var check = e.target.closest('.check');
  if (!check) return;
  e.stopPropagation();
  check.classList.toggle('checked');
  if (check.classList.contains('checked')) {
    check.innerHTML = '<svg width="8" height="8" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="3"><path d="M20 6L9 17l-5-5"/></svg>';
  } else {
    check.innerHTML = '';
  }
});
