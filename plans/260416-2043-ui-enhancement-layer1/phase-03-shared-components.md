# Phase 3: Shared Components (Filter, Pagination, Modal)

## Overview
- **Priority:** P0
- **Status:** Pending
- **Effort:** 3h
- Create reusable templ components for date filter, pagination, modal

## Implementation Steps

### 1. Date Range Filter Component
Create `web/templates/components/date_filter.templ`:
```go
templ DateFilter(fromDate string, toDate string, targetURL string) {
  <div class="flex items-center gap-3 text-xs">
    <span class="eyebrow">Từ ngày</span>
    <input type="date" name="from" value={fromDate}
      class="px-2 py-1.5 border border-hairline rounded bg-surface text-sm dark:bg-dark-surface dark:border-dark-border"
      hx-get={targetURL} hx-target="#results" hx-include="[name='to']"/>
    <span class="eyebrow">đến</span>
    <input type="date" name="to" value={toDate}
      class="px-2 py-1.5 border border-hairline rounded bg-surface text-sm dark:bg-dark-surface dark:border-dark-border"
      hx-get={targetURL} hx-target="#results" hx-include="[name='from']"/>
  </div>
}
```

### 2. Pagination Component
Create `web/templates/components/pagination.templ`:
```go
templ Pagination(currentPage int, totalPages int, baseURL string) {
  <nav class="flex items-center justify-center gap-1 mt-8">
    // Previous
    // Page numbers (show max 7: 1 ... 4 5 6 ... 20)
    // Next
  </nav>
}
```
- Style: numbered pills, active = `bg-forest text-white`, inactive = `bg-surface border`
- HTMX: each page link uses `hx-get` with `?page=N` param

### 3. Modal Component
Create `web/templates/components/modal.templ`:
```go
templ Modal(id string, title string, size string) {
  <div id={id} class="fixed inset-0 z-50 hidden" onclick="closeModal(this)">
    <div class="absolute inset-0 bg-ink/40 dark:bg-black/60 backdrop-blur-sm"></div>
    <div class="relative z-10 mx-auto mt-20 max-w-[size] bg-surface dark:bg-dark-surface rounded-lg shadow-xl border border-hairline dark:border-dark-border">
      <div class="flex items-center justify-between p-6 border-b border-hairline dark:border-dark-border">
        <h3 class="display text-2xl font-bold">{title}</h3>
        <button onclick="closeModal(this.closest('[id]'))" class="text-muted hover:text-ink">✕</button>
      </div>
      <div class="p-6">
        { children... }
      </div>
    </div>
  </div>
}
```

JS in `app.js`:
```js
function openModal(id) { document.getElementById(id).classList.remove('hidden'); }
function closeModal(el) { el.classList.add('hidden'); }
```

### 4. Integrate into Pages
- Add `@DateFilter()` to: worklogs, tasks, content, campaigns, knowledge, inbox
- Add `@Pagination()` to: worklogs, content, campaigns, knowledge
- Modal will be used in Phase 5 (inbox detail)

## Files to Create
- `web/templates/components/date_filter.templ`
- `web/templates/components/pagination.templ`
- `web/templates/components/modal.templ`

## Files to Modify
- `web/static/js/app.js` — modal open/close functions
- `web/templates/pages/worklogs.templ` — add filter + pagination
- `web/templates/pages/content.templ` — add filter + pagination
- `web/templates/pages/campaigns.templ` — add filter
- `web/templates/pages/knowledge.templ` — add filter + pagination
- `web/templates/pages/inbox.templ` — add filter
- `web/templates/pages/tasks.templ` — add filter

## Success Criteria
- Date filter renders on all tabs
- Pagination shows numbered pages (hardcoded 3 pages for demo)
- Modal opens/closes with smooth animation
- All components work in both light/dark mode
