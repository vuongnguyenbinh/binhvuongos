# Phase 5: Inbox + Kanban Enhancements

## Overview
- **Priority:** P1
- **Status:** Pending
- **Effort:** 3h
- Inbox: create button, multi-select checkboxes, detail popup with markdown editor
- Kanban: SortableJS drag-drop

## Implementation Steps

### 1. Inbox Enhancements

**"Tạo inbox" button:**
- Add button next to header, opens modal with form
- Form: content (textarea), URL, type dropdown, company dropdown

**Multi-select checkboxes:**
- Add checkbox column to inbox item list
- "Select all" checkbox in header
- When items selected: show action bar (Triage selected, Archive, Delete)

**Detail popup:**
- Click inbox item → `openModal('inbox-detail')`
- Modal shows: source badge, time, URL, content text, attachments
- Large textarea/markdown area for notes
- Triage actions at bottom (same as sidebar but in modal)

### 2. Kanban Drag-Drop

**Add SortableJS:**
In Dockerfile, add to npm install: `sortablejs`
Or use CDN: `<script src="https://cdn.jsdelivr.net/npm/sortablejs@1.15.3/Sortable.min.js"></script>`

**Create `web/static/js/sortable-init.js`:**
```js
document.addEventListener('DOMContentLoaded', function() {
  document.querySelectorAll('[data-kanban-column]').forEach(function(col) {
    new Sortable(col, {
      group: 'kanban',
      animation: 150,
      ghostClass: 'opacity-30',
      dragClass: 'shadow-lg',
      onEnd: function(evt) {
        // In Layer 1 static: just visual feedback
        // Layer 2: hx-post to update status
        var card = evt.item;
        var newColumn = evt.to.dataset.kanbanColumn;
        console.log('Moved to:', newColumn);
      }
    });
  });
});
```

**Update tasks.templ:**
- Add `data-kanban-column="todo"` etc. to each column's card container
- Add `data-card-id="1"` to each card

### 3. Inbox Modal Integration
- Use Modal component from Phase 3
- Content: inbox item info + rich text area (plain textarea for now, Milkdown in Phase 7)

## Files to Create
- `web/static/js/sortable-init.js`

## Files to Modify
- `web/templates/pages/inbox.templ` — create button, checkboxes, modal trigger
- `web/templates/pages/tasks.templ` — data attributes for SortableJS
- `web/templates/layout.templ` — add SortableJS CDN script
- `web/static/js/app.js` — checkbox select-all logic

## Success Criteria
- "Tạo inbox" button opens form modal
- Checkboxes allow multi-select with action bar
- Click inbox item opens detail popup
- Kanban cards draggable between columns
- Drag works on desktop and mobile (touch)
