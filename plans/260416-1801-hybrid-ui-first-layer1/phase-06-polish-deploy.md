# Phase 6: Responsive Polish & Demo Deployment

## Overview
- **Priority:** P1 — Final step before stakeholder demo
- **Status:** Pending
- **Effort:** 1h
- Basic responsive fixes, cross-browser check, deploy to VPS for demo

## Context Links
- PRD: `Root/binhvuong-os-prd.md` Section 11.7 (responsive breakpoints)
- All page templates from Phases 3-5

## Requirements

### Responsive (basic, not full mobile optimization)
Following demo's existing approach (desktop-first):
- **< 640px (mobile)**: Stack columns, full-width cards, hide search bar
- **640-1024px (tablet)**: Reduce grid columns, compact stat cards
- **> 1024px (desktop)**: Full layout as designed

### Priority responsive fixes
1. **Dashboard stat grid**: 4-col → 2-col on tablet, 1-col on mobile
2. **Kanban 5-col**: Horizontal scroll on tablet/mobile
3. **Work logs table**: Horizontal scroll on smaller screens
4. **Inbox 2-col**: Stack on mobile (list on top, triage below)
5. **Companies 3-col grid**: 2-col tablet, 1-col mobile
6. **Navigation tabs**: Horizontal scroll (already has `overflow-x-auto`)

### NOT in scope
- Full mobile-optimized UX (Layer 2+)
- Dark mode
- Touch-specific interactions
- PWA / service worker

## Implementation Steps

### 1. Add responsive classes to grid layouts
```html
<!-- Dashboard stat grid -->
<div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-12 gap-px ...">

<!-- Companies grid -->
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">

<!-- Kanban -->
<div class="flex gap-4 overflow-x-auto pb-4 min-w-0">
  <div class="min-w-[240px] flex-shrink-0"> <!-- each column -->
```

### 2. Table responsive wrapper
```html
<div class="overflow-x-auto -mx-4 px-4">
  <table class="w-full min-w-[800px]">
```

### 3. Hide non-essential elements on mobile
```html
<div class="hidden md:flex ..."> <!-- search bar, date display -->
<div class="hidden md:block ..."> <!-- user name in header -->
```

### 4. Test in browser
- Chrome DevTools responsive mode: 375px, 768px, 1440px
- Verify no horizontal overflow on any page
- Verify text readability on mobile
- Verify tabs scroll horizontally on mobile

### 5. Deploy to VPS via Docker Compose

```yaml
# docker-compose.prod.yml
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "3000:3000"
    restart: unless-stopped
    environment:
      - ENV=production
```

Deploy steps:
```bash
# On VPS
git pull
docker compose -f docker-compose.prod.yml up -d --build
```

Caddy config (if using Caddy as reverse proxy):
```
binhvuong.domain.com {
    reverse_proxy localhost:3000
}
```

## Todo List
- [ ] Add responsive breakpoints to dashboard stat grid
- [ ] Add responsive breakpoints to companies card grid
- [ ] Wrap kanban in horizontal scroll container
- [ ] Wrap tables in horizontal scroll container
- [ ] Stack inbox columns on mobile
- [ ] Hide non-essential header elements on mobile
- [ ] Test all 8 pages at 375px, 768px, 1440px widths
- [ ] Verify no horizontal overflow issues
- [ ] Create `docker-compose.prod.yml`
- [ ] Deploy to VPS + verify accessible via domain

## Success Criteria
- All 8 pages usable (not perfect) on 375px mobile width
- No horizontal page overflow on any screen size
- Desktop layout unchanged from previous phases
- App accessible via `https://binhvuong.domain.com` (or IP:3000 for testing)
- Docker container runs stably

## Risk Assessment
- **Caddy SSL**: First deploy may need DNS propagation. Fallback: use IP:3000 for initial demo.
- **Image optimization**: Stock images may be slow on mobile. Add `loading="lazy"` + consider WebP format.

## Next Steps (after Layer 1 demo)
→ Stakeholder feedback session
→ Plan Layer 2: DB + Auth + Core CRUD (separate plan)
→ Plan Layer 3: Business logic + integrations (separate plan)
