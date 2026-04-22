# Phase 4 — Header Avatar Dropdown

**Effort:** 20m

## Files modify

- `web/templates/layout.templ` — header avatar dropdown
- `web/static/js/app.js` — toggle helper (optional, plain onclick works)

## Template

Trong `Header()` templ, thay label user bằng:

```html
<div class="relative" onclick="toggleAvatarMenu(event)">
    if user.AvatarURL != "" {
        <img src={ user.AvatarURL } class="w-8 h-8 rounded-full cursor-pointer"/>
    } else {
        <div class="w-8 h-8 rounded-full bg-forest text-white flex items-center justify-center cursor-pointer text-sm font-semibold">
            { initials(user.FullName) }
        </div>
    }
    <div id="avatar-menu" class="hidden absolute right-0 top-10 bg-surface border border-hairline rounded shadow-lg py-2 min-w-[150px] z-50">
        <a href="/profile" class="block px-4 py-2 text-sm hover:bg-cream/50">Profile</a>
        <form method="POST" action="/auth/logout" class="block">
            <button type="submit" class="w-full text-left px-4 py-2 text-sm hover:bg-cream/50 text-rust">Logout</button>
        </form>
    </div>
</div>
```

## JS helper

```js
// app.js — append
function toggleAvatarMenu(ev) {
    ev.stopPropagation();
    const menu = document.getElementById('avatar-menu');
    menu.classList.toggle('hidden');
    // click anywhere else → close
    const onDocClick = () => { menu.classList.add('hidden'); document.removeEventListener('click', onDocClick); };
    if (!menu.classList.contains('hidden')) {
        setTimeout(() => document.addEventListener('click', onDocClick), 0);
    }
}
```

## initials helper

Trong layout.templ hoặc view_helpers.go:
```go
func Initials(name string) string {
    parts := strings.Fields(name)
    if len(parts) == 0 { return "?" }
    if len(parts) == 1 { return strings.ToUpper(parts[0][:1]) }
    return strings.ToUpper(string(parts[0][0]) + string(parts[len(parts)-1][0]))
}
```

## Layout needs user.AvatarURL

Check: Layout() currently receives user from context? If not, thread it via templ param or use `ctx` pattern. Most templ-based Go projects pass user down — check existing `layout.templ` signature.

## Todo
- [ ] Check layout.templ signature — add User parameter if missing
- [ ] Initials helper
- [ ] Dropdown HTML
- [ ] JS toggle
- [ ] Bump `/static/js/app.js?v=` in layout for cache bust

## Success criteria
- Avatar visible in header (img nếu có, initials fallback)
- Click → dropdown toggle
- Click outside → dropdown đóng
- Click Profile → /profile
- Click Logout → form POST /auth/logout → /login
