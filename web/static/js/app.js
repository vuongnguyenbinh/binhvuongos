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

// Modal open/close
function openModal(id) {
  var el = document.getElementById(id);
  if (el) { el.classList.remove('hidden'); document.body.style.overflow = 'hidden'; }
}
function closeModal(id) {
  var el = typeof id === 'string' ? document.getElementById(id) : id;
  if (el) { el.classList.add('hidden'); document.body.style.overflow = ''; }
}
// Close modal on backdrop click
document.addEventListener('click', function(e) {
  if (e.target.classList.contains('modal-backdrop')) {
    closeModal(e.target.closest('[data-modal]'));
  }
});
// Close modal on Escape
document.addEventListener('keydown', function(e) {
  if (e.key === 'Escape') {
    var open = document.querySelector('[data-modal]:not(.hidden)');
    if (open) closeModal(open);
  }
});

// Copy to clipboard
function copyText(text, btn) {
  navigator.clipboard.writeText(text).then(function() {
    var orig = btn.textContent;
    btn.textContent = '✓ Đã copy';
    setTimeout(function() { btn.textContent = orig; }, 1500);
  });
}

// Select all checkboxes
function toggleSelectAll(masterCheckbox, containerSelector) {
  var checks = document.querySelectorAll(containerSelector + ' .item-check');
  checks.forEach(function(c) {
    if (masterCheckbox.checked) {
      c.checked = true;
    } else {
      c.checked = false;
    }
  });
  updateBulkActions(containerSelector);
}
function updateBulkActions(containerSelector) {
  var checked = document.querySelectorAll(containerSelector + ' .item-check:checked').length;
  var bar = document.getElementById('bulk-actions');
  if (bar) {
    bar.classList.toggle('hidden', checked === 0);
    var count = bar.querySelector('.bulk-count');
    if (count) count.textContent = checked;
  }
}

// Inbox batch checkboxes
function updateInboxBatch() {
  var checks = document.querySelectorAll('.inbox-check:checked');
  var bar = document.getElementById('inbox-batch-bar');
  var countEl = document.getElementById('inbox-selected-count');
  var idsEl = document.getElementById('inbox-batch-ids');
  if (!bar) return;
  var ids = [];
  checks.forEach(function(c) { ids.push(c.getAttribute('data-id')); });
  bar.classList.toggle('hidden', ids.length === 0);
  if (countEl) countEl.textContent = ids.length;
  if (idsEl) idsEl.value = ids.join(',');
}

// Init SortableJS on kanban columns
document.addEventListener('DOMContentLoaded', function() {
  document.querySelectorAll('[data-kanban-column]').forEach(function(col) {
    if (typeof Sortable !== 'undefined') {
      new Sortable(col, {
        group: 'kanban',
        animation: 150,
        ghostClass: 'opacity-30',
        dragClass: 'shadow-lg',
        onEnd: function(evt) {
          console.log('Moved card to:', evt.to.dataset.kanbanColumn);
        }
      });
    }
  });

  // Fix dark mode icon on load
  var icon = document.getElementById('dark-toggle-icon');
  if (icon) {
    icon.textContent = document.documentElement.classList.contains('dark') ? '☀️' : '🌙';
  }

  // Populate header avatar from /auth/me
  var slot = document.getElementById('avatar-slot');
  if (slot) {
    fetch('/auth/me').then(function(r) { return r.ok ? r.json() : null; }).then(function(me) {
      if (!me) return;
      if (me.avatar_url) {
        slot.innerHTML = '<img src="' + me.avatar_url + '" alt="Avatar" class="w-8 h-8 object-cover rounded-full"/>';
      } else if (me.full_name) {
        var parts = me.full_name.trim().split(/\s+/);
        var init = parts.length === 1 ? parts[0][0] : (parts[0][0] + parts[parts.length-1][0]);
        slot.textContent = init.toUpperCase();
      }
    }).catch(function() {});
  }
});

// Avatar dropdown toggle — used by header button
function toggleAvatarMenu(ev) {
  ev.stopPropagation();
  var menu = document.getElementById('avatar-menu');
  if (!menu) return;
  menu.classList.toggle('hidden');
  if (!menu.classList.contains('hidden')) {
    var closer = function() { menu.classList.add('hidden'); document.removeEventListener('click', closer); };
    setTimeout(function() { document.addEventListener('click', closer); }, 0);
  }
}
