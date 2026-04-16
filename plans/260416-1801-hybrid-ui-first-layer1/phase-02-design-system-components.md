# Phase 2: Design System & Shared Components

## Overview
- **Priority:** P0 — shared components used by all page views
- **Status:** Pending
- **Effort:** 3h
- Create reusable templ components + base layout + apply color/image design changes

## Context Links
- Demo HTML: `Root/binhvuong-os-demo.html` lines 1-241 (styles), 244-297 (header/nav)
- Plan: `plan.md` → Design Changes section

## Key Insights

### Color Mapping (amber → red/ember)
| Demo Original | New Value | CSS Variable | Usage |
|---------------|-----------|--------------|-------|
| `#B8741F` (amber2) | `#D94F30` (ember) | `--ember` | Primary accent, notification dots, badges, active pills |
| `#C89B3C` (ochre) | `#E8623A` (flame) | `--flame` | Warning indicators, secondary accent |
| `#FDF4E2` (amber bg) | `#FDE8E4` (ember-light) | — | Pill backgrounds for CAO/CHU Y status |
| `#B8741F` (amber pill text) | `#D94F30` | — | Pill text color |

### Stock Image Sections
Sections currently using solid `bg-forest`:
1. **Dashboard stat box** "Chiến dịch đang chạy" (line 406-431) → landscape + forest overlay 75%
2. **Knowledge featured card** (line 2192-2209) → landscape + forest overlay 75%

Pattern for stock image sections:
```html
<div class="relative overflow-hidden rounded-lg">
  <img src="/static/img/landscape-01.jpg" class="absolute inset-0 w-full h-full object-cover" />
  <div class="absolute inset-0 bg-forest/75 backdrop-blur-sm"></div>
  <div class="relative z-10 p-6 text-white">
    <!-- content -->
  </div>
</div>
```

## Requirements

### Shared Components to Extract
From demo HTML, identify reusable patterns:

| Component | Demo Usage | templ File |
|-----------|------------|------------|
| Layout (html/head/body) | Every page | `layout.templ` |
| Header + Navigation tabs | Lines 247-297 | `partials/header.templ` |
| Footer + sync status | Lines 2326-2341 | `partials/footer.templ` |
| Pill badge | ~30 instances | `components/pill.templ` |
| Progress bar | ~15 instances | `components/progress.templ` |
| Stat hero card | 4 dashboard cards | `components/stat_card.templ` |
| Kanban card | ~15 cards | `components/kanban_card.templ` |
| Checkbox | ~10 instances | `components/check.templ` |
| Hero section (stock img bg) | 2 sections | `components/hero_section.templ` |
| Company row | 5 rows dashboard | `components/company_row.templ` |
| Eyebrow label | ~40 instances | CSS class only (no templ) |
| Spark bar chart | 1 instance | inline (too specific for component) |

## Files to Create

| File | Purpose |
|------|---------|
| `web/templates/layout.templ` | Base HTML layout with head, fonts, CSS, JS |
| `web/templates/partials/header.templ` | Logo + search + notifications + user menu |
| `web/templates/partials/tabs.templ` | Navigation tabs with active state |
| `web/templates/partials/footer.templ` | Footer with sync status |
| `web/templates/components/pill.templ` | Pill badge (color variants) |
| `web/templates/components/progress.templ` | Progress bar (track + fill) |
| `web/templates/components/stat_card.templ` | Stat hero card |
| `web/templates/components/kanban_card.templ` | Kanban card |
| `web/templates/components/check.templ` | Interactive checkbox |
| `web/templates/components/hero_section.templ` | Section with stock image bg + blur overlay |
| `web/templates/components/company_row.templ` | Company list row |
| `web/static/img/landscape-01.jpg` | Mountain/lake landscape (Unsplash) |
| `web/static/img/landscape-02.jpg` | Field/meadow landscape |
| `web/static/img/landscape-03.jpg` | Forest/valley landscape |

## Implementation Steps

### 1. Download stock images
3 landscape images from Unsplash (free license), resize to 1200x600 for optimal load:
- Mountain lake (for campaigns section)
- Rolling hills (for knowledge featured)
- Forest canopy (spare/alternate)

```bash
# Example — actual URLs to be determined at implementation
# Resize to 1200x600, optimize with quality 80
convert input.jpg -resize 1200x600^ -gravity center -extent 1200x600 -quality 80 landscape-01.jpg
```

### 2. Layout template (`layout.templ`)
```go
package templates

templ Layout(title string, activeTab string) {
    <!DOCTYPE html>
    <html lang="vi">
    <head>
        <meta charset="UTF-8"/>
        <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <title>{ title } — Bình V��ơng OS</title>
        <link rel="preconnect" href="https://fonts.googleapis.com"/>
        <link href="https://fonts.googleapis.com/css2?family=Phudu:wght@300;400;500;600;700;800;900&family=Inter:wght@300;400;500;600;700&family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet"/>
        <link rel="stylesheet" href="/static/css/output.css"/>
        <script src="https://unpkg.com/htmx.org@1.9.12"></script>
    </head>
    <body class="min-h-screen paper">
        @Header(activeTab)
        <main class="max-w-[1440px] mx-auto px-8 py-10">
            { children... }
        </main>
        @Footer()
        <script src="/static/js/app.js"></script>
    </body>
    </html>
}
```

### 3. Pill component (`components/pill.templ`)
Variants: forest, ember, rust, sage, muted
```go
package components

type PillVariant string
const (
    PillForest PillVariant = "forest"
    PillEmber  PillVariant = "ember"   // was amber
    PillRust   PillVariant = "rust"
    PillSage   PillVariant = "sage"
    PillMuted  PillVariant = "muted"
)

templ Pill(label string, variant PillVariant, dotColor string) {
    <span class={ "pill", variantClass(variant) }>
        if dotColor != "" {
            <span class="dot" style={ "background:" + dotColor }></span>
        }
        { label }
    </span>
}
```

### 4. Progress bar (`components/progress.templ`)
```go
package components

templ Progress(percent int, height string, gradient string) {
    <div class="progress-track" style={ heightStyle(height) }>
        <div class="progress-fill" style={ widthAndBg(percent, gradient) }></div>
    </div>
}
```

### 5. Hero section with stock image (`components/hero_section.templ`)
```go
package components

templ HeroSection(imageURL string, overlayOpacity string) {
    <div class="relative overflow-hidden">
        <img src={ imageURL } class="absolute inset-0 w-full h-full object-cover" alt=""/>
        <div class={ "absolute inset-0 backdrop-blur-sm", "bg-forest/" + overlayOpacity }></div>
        <div class="relative z-10">
            { children... }
        </div>
    </div>
}
```

### 6. Header with tabs (`partials/header.templ`)
Port lines 247-297 from demo. Key changes:
- Tab links use `href` instead of `onclick="switchView()"` (full page nav for Layer 1)
- Active tab determined by `activeTab` parameter
- Notification badge uses ember color instead of amber2
- `hx-boost="true"` on nav links for HTMX smooth transitions (optional)

### 7. Update CSS variables in `input.css`
```css
:root {
    --bg: #FAF8F3;
    --surface: #FFFFFF;
    --ink: #1A1918;
    --forest: #1F3D2E;
    --ember: #D94F30;     /* was --amber: #B8741F */
    --flame: #E8623A;     /* new, was ochre */
    --muted: #6B665E;
    --hairline: #E8E4DB;
    --cream: #F2EEE4;
}
```

Replace all CSS occurrences:
- `var(--amber)` → `var(--ember)`
- `color: #B8741F` → `color: #D94F30`
- `background: #C89B3C` → `background: #E8623A`
- `.text-amber2` → `.text-ember`
- `.bg-amber2` → `.bg-ember`
- Pill backgrounds: `#FDF4E2` → `#FDE8E4`

## Todo List
- [ ] Download + optimize 3 stock landscape images
- [ ] Create `layout.templ` with full HTML head, fonts, CSS
- [ ] Create `header.templ` with logo, search, notifications, user menu
- [ ] Create `tabs.templ` with active tab logic
- [ ] Create `footer.templ` with sync status bar
- [ ] Create `pill.templ` with color variants (forest/ember/rust/sage/muted)
- [ ] Create `progress.templ` with height/gradient params
- [ ] Create `stat_card.templ` for dashboard stats
- [ ] Create `kanban_card.templ` with priority/status/progress
- [ ] Create `check.templ` for interactive checkboxes
- [ ] Create `hero_section.templ` with stock image + blur overlay
- [ ] Create `company_row.templ` for company list
- [ ] Update `input.css` — replace all amber/ochre refs with ember/flame
- [ ] Update `tailwind.config.js` — verify ember/flame in theme
- [ ] Run `templ generate` + `tailwind build` → verify no errors

## Success Criteria
- All components render correctly in isolation
- Pill colors: ember (đỏ cam) replaces old amber everywhere
- Hero sections show landscape image with blurred forest overlay
- Layout wraps any page content correctly
- Tab navigation highlights active page
- `templ generate` produces valid Go code
- Tailwind output includes all custom classes

## Risk Assessment
- **templ children slot**: `{ children... }` syntax for slot-like composition — verify templ version supports it
- **Stock images file size**: Large images slow page load. Mitigate: compress to <100KB each, lazy load with `loading="lazy"`
- **Tailwind + templ class scanning**: Tailwind needs to scan `.templ` files for class names. Verify `content` config path includes `**/*.templ`

## Next Steps
→ Phase 3: Dashboard + Work Logs views (highest priority pages)
