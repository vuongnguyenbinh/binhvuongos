# Phase 2: Mobile Responsive

## Overview
- **Priority:** P0
- **Status:** Pending
- **Effort:** 2h
- Polish responsive layouts for mobile (375px) and tablet (768px)

## Implementation Steps

### 1. Mobile Header
- Hide search bar below `md:`
- Hamburger menu button → toggle drawer
- Tabs → horizontal scroll (already done), reduce gap

### 2. Dashboard Responsive
- Stat grid: `grid-cols-1 sm:grid-cols-2 xl:grid-cols-12` (already done)
- Company list: hide stat columns below `md:`
- Today tasks + inbox: stack below company list on mobile

### 3. Tables (Work Logs, Content)
- Wrap in `overflow-x-auto` (already done)
- Add `min-w-[900px]` to table (already done)
- Consider card view alternative for mobile (stretch goal)

### 4. Kanban (Tasks)
- `overflow-x-auto` with `min-w-[240px]` columns (already done)
- Add horizontal scroll indicator on mobile

### 5. Companies Grid
- `grid-cols-1 md:grid-cols-2 lg:grid-cols-3` (already done)
- Verify card content doesn't overflow

### 6. Knowledge Grid
- `grid-cols-1 md:grid-cols-2 lg:grid-cols-4` (already done)
- Featured card: `md:col-span-2` (already done)

### 7. Forms (Inbox Triage)
- Stack 2-column layout on mobile: `grid-cols-1 xl:grid-cols-12`
- Triage panel: not sticky on mobile, flows below list

### 8. Test All Breakpoints
- Chrome DevTools: 375px (iPhone), 768px (iPad), 1440px (desktop)
- Verify no horizontal overflow on any page

## Files to Modify
- `web/templates/layout.templ` — mobile nav drawer
- All page templates — verify/fix responsive classes

## Success Criteria
- No horizontal page overflow at 375px
- All content readable on mobile
- Tables scroll horizontally
- Navigation usable on mobile
