# Phase 7: Milkdown Editor + Deploy

## Overview
- **Priority:** P2
- **Status:** Pending
- **Effort:** 2h
- Integrate Milkdown markdown editor, deploy all changes

## Implementation Steps

### 1. Milkdown Setup
Add to Dockerfile npm install: `@milkdown/core @milkdown/preset-commonmark @milkdown/theme-nord @milkdown/plugin-listener`

Create `web/static/js/milkdown-init.js`:
```js
// Lazy load: only init when editor container exists
document.addEventListener('DOMContentLoaded', async function() {
  var editors = document.querySelectorAll('[data-milkdown]');
  if (!editors.length) return;

  // Dynamic import
  var { Editor, rootCtx } = await import('@milkdown/core');
  var { commonmark } = await import('@milkdown/preset-commonmark');
  var { listener, listenerCtx } = await import('@milkdown/plugin-listener');

  editors.forEach(function(el) {
    var textarea = el.querySelector('textarea[name]');
    Editor.make()
      .config(ctx => {
        ctx.set(rootCtx, el);
        ctx.set(listenerCtx, {
          markdown: [(getMarkdown) => {
            textarea.value = getMarkdown();
          }]
        });
      })
      .use(commonmark)
      .use(listener)
      .create();
  });
});
```

**Alternative (simpler for Layer 1):** Use a simple textarea with markdown preview toggle instead of full Milkdown. Less dependencies, faster build.

### 2. Editor Integration Points
- Inbox detail popup: note field
- Content detail: body field
- Knowledge detail: body field
- Prompt create/edit: prompt body

### 3. Build + Deploy
```bash
# Local
docker build -t binhvuongos .
docker compose up -d

# Server
ssh root@103.97.125.186
cd /opt/binhvuongos
git pull
docker compose up -d --build
```

### 4. Smoke Test All Pages
After deploy, verify:
- All 10 tabs render (8 original + bookmarks + prompts)
- Dark mode toggle works
- Date filters render
- Pagination renders
- Detail pages accessible
- Kanban drag works
- Inbox modal opens
- Markdown editor loads

## Files to Create
- `web/static/js/milkdown-init.js` (or simpler markdown-preview.js)

## Files to Modify
- `Dockerfile` — npm deps
- `web/templates/layout.templ` — script tag for editor
- Various detail templates — add `data-milkdown` containers

## Success Criteria
- Markdown editor loads in popup/detail views
- Output is valid Markdown (paste into Notion/Obsidian preserves format)
- All changes deployed to https://os.binhvuong.vn
- All 10 tabs + detail pages working in production
- Dark mode works end-to-end
