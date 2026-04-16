# Phase 1: Fix Hero + Dark Mode

## Overview
- **Priority:** P0
- **Status:** Pending
- **Effort:** 4h
- Fix hero section readability, implement Tailwind dark mode across all pages

## Implementation Steps

### 1. Fix Hero Overlay (15min)
In `web/static/css/input.css`:
- `.hero-overlay` background: `rgba(31,61,46,0.78)` → `rgba(31,61,46,0.55)`
- Add text shadow to hero content: `text-shadow: 0 1px 3px rgba(0,0,0,0.3)`
- Ensure `.hero-content` text uses `font-semibold` for better contrast

### 2. Tailwind Dark Mode Config
In `tailwind.config.js`:
```js
module.exports = {
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // existing...
        // Dark palette
        'dark-bg': '#0F1117',
        'dark-surface': '#1A1D26',
        'dark-text': '#E8E6E3',
        'dark-border': '#2D3039',
        'dark-cream': '#1F2229',
        'dark-muted': '#8B8880',
      }
    }
  }
}
```

### 3. CSS Variables for Dark Mode
In `input.css` add:
```css
.dark {
  --bg: #0F1117;
  --surface: #1A1D26;
  --ink: #E8E6E3;
  --forest: #1F3D2E;
  --ember: #D94F30;
  --flame: #E8623A;
  --muted: #8B8880;
  --hairline: #2D3039;
  --cream: #1F2229;
}
```
Update all custom CSS classes to use `var()` where hardcoded.

### 4. Dark Toggle Component
Create `web/templates/components/dark_toggle.templ`:
- Sun/moon icon button in header
- Calls `toggleDarkMode()` JS function

Create `web/static/js/dark-mode.js`:
```js
function toggleDarkMode() {
  document.documentElement.classList.toggle('dark');
  localStorage.setItem('dark', document.documentElement.classList.contains('dark'));
}
// Auto-apply on load
if (localStorage.getItem('dark') === 'true') {
  document.documentElement.classList.add('dark');
}
```

### 5. Add dark: Classes to All Templates
For each template, add dark variants:
- `bg-surface` → `bg-surface dark:bg-dark-surface`
- `bg-ivory` → `bg-ivory dark:bg-dark-bg`
- `text-ink` → `text-ink dark:text-dark-text`
- `border-hairline` → `border-hairline dark:border-dark-border`
- `bg-cream` → `bg-cream dark:bg-dark-cream`
- `text-muted` → `text-muted dark:text-dark-muted`

**Strategy:** Use CSS variable approach for base styles (html/body, .paper, .pill, .kanban-card etc.) so they auto-switch. Only add `dark:` classes for Tailwind utility classes in templates.

### 6. Layout Template Updates
In `layout.templ`:
- Add `<script src="/static/js/dark-mode.js"></script>` in `<head>` (before body renders)
- Add dark toggle button in header next to notification bell
- `<body>` background handled by CSS vars
- `<header>` add `dark:bg-dark-bg/80`

## Files to Modify
- `tailwind.config.js` — add darkMode + dark colors
- `web/static/css/input.css` — dark CSS vars, fix hero overlay
- `web/templates/layout.templ` — dark toggle, dark header
- `web/templates/pages/dashboard.templ` — dark classes
- `web/templates/pages/worklogs.templ` — dark classes
- `web/templates/pages/inbox.templ` — dark classes
- `web/templates/pages/tasks.templ` — dark classes
- `web/templates/pages/content.templ` — dark classes
- `web/templates/pages/companies.templ` — dark classes
- `web/templates/pages/campaigns.templ` — dark classes
- `web/templates/pages/knowledge.templ` — dark classes

## Files to Create
- `web/static/js/dark-mode.js`

## Success Criteria
- Hero text clearly readable on landscape images
- Dark mode toggles correctly with button
- Dark preference persists via localStorage
- All 8 pages render correctly in both light/dark mode
- Accent colors (forest, ember, sage, rust) visible in both modes
