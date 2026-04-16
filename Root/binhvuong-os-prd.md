# PRD — Solo Expert Bình Vương OS

> **Product Requirements Document**
> Version: 1.0
> Ngày: 16/04/2026
> Stakeholder: Bình Vương (Solo Founder)
> Status: Draft — Ready for Development

---

## Mục Lục

1. [Executive Summary](#1-executive-summary)
2. [Bối Cảnh & Vấn Đề](#2-bối-cảnh--vấn-đề)
3. [Mục Tiêu & Phạm Vi](#3-mục-tiêu--phạm-vi)
4. [Personas & Use Cases](#4-personas--use-cases)
5. [Kiến Trúc Hệ Thống](#5-kiến-trúc-hệ-thống)
6. [Functional Requirements](#6-functional-requirements)
7. [Non-Functional Requirements](#7-non-functional-requirements)
8. [User Flows](#8-user-flows)
9. [Data Model](#9-data-model)
10. [API Specifications](#10-api-specifications)
11. [UI/UX Requirements](#11-uiux-requirements)
12. [Notion Sync Specifications](#12-notion-sync-specifications)
13. [Security & Permissions](#13-security--permissions)
14. [Roadmap & Milestones](#14-roadmap--milestones)
15. [Success Metrics](#15-success-metrics)
16. [Risks & Mitigations](#16-risks--mitigations)
17. [Out of Scope](#17-out-of-scope)

---

## 1. Executive Summary

### 1.1 Tầm nhìn sản phẩm

Xây dựng hệ điều hành nội bộ (internal OS) cho một Solo Founder đang đồng thời quản lý doanh nghiệp riêng, vai trò co-founder và mentor cho dưới 10 công ty client. Hệ thống tập trung hoá việc assign task, theo dõi tiến độ, quản lý sản xuất nội dung, và đo lường output hàng ngày của 20 nhân sự/CTV qua một giao diện duy nhất có thương hiệu riêng.

### 1.2 Giá trị cốt lõi

- **Cho Solo Founder:** Một cockpit điều khiển nhìn thấy tình hình tất cả dự án, phát hiện vấn đề sớm, không còn phân tán trên 10+ công cụ
- **Cho nhân sự/CTV:** Một nơi duy nhất xem task, submit báo cáo, không phải học 5 công cụ cho 5 công ty khác nhau
- **Cho kinh doanh:** Tự động hoá báo cáo, giảm thời gian quản trị từ 2–3 giờ/ngày xuống 30–45 phút/ngày

### 1.3 Stack kỹ thuật

- **Backend:** Go (Fiber hoặc Chi) + PostgreSQL 16
- **Frontend:** templ + HTMX (server-rendered)
- **Deployment:** Docker Compose trên VPS Linux hiện có
- **Reverse proxy:** Caddy (auto HTTPS)
- **Integration:** Notion API (sync 1 chiều mỗi giờ)
- **Capture:** Telegram Bot (webhook)

### 1.4 Timeline

MVP: **8 tuần**. Full production-ready: **12 tuần**.

---

## 2. Bối Cảnh & Vấn Đề

### 2.1 Hiện trạng

Bình Vương hiện quản lý đa dự án với công cụ phân tán:

- **Notion Plus** — dùng cho cá nhân và 1 số company
- **Google Docs** — 80% content CTV viết nằm rải rác
- **Google Sheets** — nhân sự báo cáo output hàng ngày (backlink, content, ads)
- **Telegram** — capture cá nhân, brief CTV
- **Zalo** — giao tiếp với nhân sự client
- **Email** — báo cáo client, brief dự án

**Vấn đề chính:**

1. **Phân tán data** — không có nguồn sự thật duy nhất
2. **Rate limit Notion API** — không thể dùng Notion làm backend cho 20 user đồng thời (giới hạn 3 req/s)
3. **Không có dashboard tổng hợp** — phải vào từng tool để thấy số liệu
4. **Tracking output thủ công** — mỗi cuối tuần phải tự cộng số từ nhiều sheet
5. **Phân quyền rối** — mỗi client có tool riêng, mỗi nhân sự dùng 2–3 tool

### 2.2 Giải pháp đề xuất

Xây dựng web app tự host với:
- **PostgreSQL làm source of truth** — không còn giới hạn rate limit
- **UI có brand riêng** — professional khi làm việc với client
- **Sync 1 chiều sang Notion** — giữ Notion như kho tham khảo cá nhân cho Solo Founder
- **Tích hợp Telegram Bot** — capture nhanh không rời khỏi thói quen hiện tại
- **Giữ Google Sheets** — nhân sự vẫn dùng sheet chi tiết, app chỉ lưu số tổng + link

---

## 3. Mục Tiêu & Phạm Vi

### 3.1 Business Goals

| # | Mục tiêu | Đo lường | Timeline |
|---|---|---|---|
| G1 | Giảm thời gian quản trị hàng ngày | Từ 2–3h → < 45 phút | Sau 3 tháng go-live |
| G2 | Tăng độ chính xác báo cáo output | 100% data có verification | Sau 1 tháng |
| G3 | Xoá bỏ phụ thuộc công cụ rời rạc | Chỉ dùng 1 app + Notion (read-only) | Sau 2 tháng |
| G4 | Nâng chuyên nghiệp với client | Brand riêng trên portal | Từ MVP |

### 3.2 In Scope (MVP)

✅ **Module 1:** Authentication & User Management
✅ **Module 2:** Company Management (dưới 10 công ty)
✅ **Module 3:** Inbox (capture từ Telegram + manual)
✅ **Module 4:** Tasks (thay Projects bằng group_name)
✅ **Module 5:** Content Pipeline (gộp pipeline + library)
✅ **Module 6:** Work Logs (module tracking output mới)
✅ **Module 7:** Campaigns
✅ **Module 8:** Knowledge Base (gộp knowledge + raw material)
✅ **Module 9:** Dashboard (đa chiều cho Owner)
✅ **Module 10:** Notion Sync Worker (1h/lần, 1 chiều)

### 3.3 Out of Scope (MVP) — xem Section 17

❌ Mobile app native (dùng responsive web)
❌ HR module (chấm công, lương)
❌ Accounting/Finance
❌ CRM cho khách hàng cuối
❌ Chat nội bộ (dùng Zalo)
❌ Sync 2 chiều với Notion
❌ AI auto-classification cho Inbox

---

## 4. Personas & Use Cases

### 4.1 Personas

#### Persona 1: Bình Vương (Owner)
- **Vai trò:** Solo Founder, quản lý + mentor dưới 10 công ty
- **Đặc điểm:** Guồng làm việc nhanh, multitask cao, cần bird's-eye view
- **Thiết bị chính:** Desktop (làm việc sâu) + iPhone (capture di động)
- **Pain points:** Context switching, mất thời gian tổng hợp báo cáo
- **Mục tiêu dùng app:** Điều hành, review, ra quyết định

#### Persona 2: Staff Cốt Lõi (2–3 người)
- **Vai trò:** Cánh tay phải của Bình Vương, tham gia nhiều công ty
- **Đặc điểm:** Có kỹ năng đa dạng, làm việc với nhiều CTV
- **Thiết bị chính:** Desktop + mobile
- **Pain points:** Theo dõi task đa dự án, review công việc CTV
- **Mục tiêu dùng app:** Quản lý đội nhỏ, delegate, review

#### Persona 3: CTV (5–10 người)
- **Vai trò:** Nhận việc theo dự án (content, SEO, design)
- **Đặc điểm:** Làm part-time, có thể tham gia 1–2 công ty
- **Thiết bị chính:** Desktop
- **Pain points:** Nhớ deadline, nộp bài, theo dõi phí
- **Mục tiêu dùng app:** Nhận task, submit bài, log output hàng ngày

#### Persona 4: Nhân sự Client (5–10 người)
- **Vai trò:** Nhân viên của công ty khách, được Bình Vương mentor
- **Đặc điểm:** Không quen tool mới, ưa dùng Zalo + Google Sheets
- **Thiết bị chính:** Desktop + mobile
- **Pain points:** Học tool mới khó, sợ bị phức tạp hoá
- **Mục tiêu dùng app:** Nhận hướng dẫn, nộp báo cáo cuối ngày đơn giản

### 4.2 Core Use Cases

| UC | Persona | Mô tả |
|---|---|---|
| UC-01 | Owner | Ghi nhanh ý tưởng/link vào Inbox qua Telegram |
| UC-02 | Owner | Triage Inbox mỗi sáng (10 phút) |
| UC-03 | Owner | Assign task cho CTV/Staff |
| UC-04 | Owner | Review work_logs cuối ngày (approve/reject) |
| UC-05 | Owner | Xem dashboard tổng quan — Health từng công ty |
| UC-06 | Owner | Tạo campaign + set target theo loại work |
| UC-07 | Staff | Quản lý task của CTV mà mình phụ trách |
| UC-08 | Staff | Review content draft của CTV |
| UC-09 | CTV | Xem task được giao hôm nay |
| UC-10 | CTV | Submit content sau khi viết xong |
| UC-11 | CTV/Staff | Log output cuối ngày (form 30s) |
| UC-12 | Client Staff | Nhận task + bài giảng từ Bình Vương |
| UC-13 | Client Staff | Submit work log cuối ngày với link sheet |
| UC-14 | All | Truy cập Knowledge Base xem bài giảng/SOP |

---

## 5. Kiến Trúc Hệ Thống

### 5.1 High-level Architecture

```
┌───────────────────────────────────────────────────────┐
│   CLIENTS                                              │
│   ├── Browser (20 users, desktop + mobile responsive) │
│   └── Telegram (Owner capture only)                   │
└────────────────────────┬──────────────────────────────┘
                         │ HTTPS
                         ▼
┌───────────────────────────────────────────────────────┐
│   VPS (Linux + Docker)                                 │
│                                                        │
│   ┌─────────────────────────────────────────────┐    │
│   │  Caddy (reverse proxy, auto SSL)            │    │
│   └──────────────────┬──────────────────────────┘    │
│                      │                                │
│   ┌──────────────────▼──────────────────────────┐    │
│   │  Go App Server (Fiber/Chi)                  │    │
│   │  ├── HTTP Handlers                          │    │
│   │  ├── Auth Middleware (JWT)                  │    │
│   │  ├── Business Logic Services                │    │
│   │  ├── templ Renderer                         │    │
│   │  └── HTMX endpoints                         │    │
│   └──────────┬────────────────────────┬─────────┘    │
│              │                        │               │
│   ┌──────────▼──────────┐  ┌─────────▼─────────┐    │
│   │  PostgreSQL 16      │  │  Sync Worker      │    │
│   │  (Source of Truth)  │  │  (Go, cron 1h)    │    │
│   └─────────────────────┘  └─────────┬─────────┘    │
│                                      │               │
│   ┌──────────────────────────────────▼─────────┐    │
│   │  Telegram Bot Handler (Go webhook)         │    │
│   └────────────────────────────────────────────┘    │
└───────────────────────────┬───────────────────────────┘
                            │
                            ▼
┌───────────────────────────────────────────────────────┐
│   EXTERNAL SERVICES                                    │
│   ├── Notion API (sync 1 chiều)                       │
│   ├── Telegram Bot API (webhook in)                   │
│   └── SMTP (email notifications, optional)            │
└───────────────────────────────────────────────────────┘
```

### 5.2 Component Responsibilities

| Component | Trách nhiệm |
|---|---|
| **Caddy** | Terminate SSL, route to Go app, gzip, security headers |
| **Go App Server** | HTTP handlers, business logic, render UI |
| **PostgreSQL** | Source of truth, ACID transactions |
| **Sync Worker** | Background job — đẩy thay đổi sang Notion mỗi giờ |
| **Telegram Bot Handler** | Webhook endpoint nhận message, lưu vào Inbox |

### 5.3 Data Flow Examples

**Flow 1: Owner capture qua Telegram**
```
Owner gửi link TikTok vào Bot
  → Telegram Bot API POST webhook → Go app
  → App verify telegram_id khớp user `owner`
  → Tạo inbox_items record (source='telegram', status='raw')
  → Reply Telegram: "✅ Đã lưu vào Inbox"
```

**Flow 2: CTV submit work log**
```
CTV mở app → /work-log/new
  → Chọn company, work_type, điền quantity, paste link sheet
  → Submit → POST /api/work-logs
  → App tạo work_logs record (status='submitted')
  → HTMX swap thông báo "Đã gửi, chờ duyệt"
```

**Flow 3: Sync sang Notion**
```
Cron 1h chạy → Sync Worker
  → Query: SELECT * FROM work_logs WHERE sync_status='pending' AND status='approved'
  → Rate limit 2 req/s
  → Với mỗi record: POST Notion API → tạo page
  → Update record: sync_status='synced', notion_page_id='...'
  → Log vào notion_sync_log
```

---

## 6. Functional Requirements

### 6.1 Module 1: Authentication & User Management

#### FR-AUTH-01: Đăng nhập email/password
**Priority:** P0 (Must have)

- User nhập email + password
- App verify với `users` table (bcrypt hash)
- Tạo JWT token (httpOnly cookie, 7 ngày)
- Redirect to dashboard theo role

**Acceptance criteria:**
- [ ] Sai password → hiển thị lỗi, không tiết lộ email có tồn tại không
- [ ] Sau 5 lần sai → lock 15 phút
- [ ] Remember me checkbox → extend token lên 30 ngày

#### FR-AUTH-02: Owner tạo user mới
**Priority:** P0

- Chỉ `owner` có quyền truy cập `/admin/users/new`
- Form: email, full_name, role, password tạm, companies được gán
- System gửi email chào mừng với link đổi password (optional ở MVP)

**Acceptance criteria:**
- [ ] Tạo user thành công → user có thể login ngay
- [ ] Email trùng → báo lỗi
- [ ] Role không thuộc enum → báo lỗi

#### FR-AUTH-03: User đổi password
**Priority:** P1

- `/profile/security` — form đổi password
- Yêu cầu nhập password cũ để verify
- Password mới ≥ 8 ký tự

#### FR-AUTH-04: Forgot password
**Priority:** P2 (Post-MVP)

- Flow gửi email reset link (cần SMTP setup)

---

### 6.2 Module 2: Company Management

#### FR-CO-01: Owner tạo công ty
**Priority:** P0

- Form: name, short_code, industry, my_role, scope, contact info
- Tự động tạo slug từ name
- Sau khi tạo → redirect to company detail page

**Acceptance criteria:**
- [ ] short_code unique
- [ ] Có thể upload logo (tuỳ chọn)

#### FR-CO-02: Danh sách công ty với health indicator
**Priority:** P0

- `/companies` — list all companies với filter theo status
- Mỗi row hiển thị: name, logo, industry, my_role, scope (chips), health (🟢🟡🔴), task_count active
- Click vào → company detail page

#### FR-CO-03: Company detail page
**Priority:** P0

- Layout: header (info) + tabs (Tasks | Content | Work Logs | Campaigns | Objectives | Knowledge | People)
- Mỗi tab là một filtered view của module tương ứng

#### FR-CO-04: Owner assign user vào company
**Priority:** P0

- Tab "People" → button "Add person"
- Chọn user từ dropdown + role_in_company + quyền (can_view, can_edit, can_approve)
- Tạo record trong `user_company_assignments`

#### FR-CO-05: Update health
**Priority:** P1

- Owner có thể đổi health 🟢🟡🔴 bằng click nhanh trên dashboard
- Ghi nhận vào `companies.health` + timestamp vào updated_at

---

### 6.3 Module 3: Inbox

#### FR-INBOX-01: Nhận tin nhắn từ Telegram Bot
**Priority:** P0

- Endpoint webhook: `POST /api/telegram/webhook`
- Verify: `from.id` phải khớp với một `users.telegram_id`
- Parse message:
  - Text đơn thuần → lưu vào `content`
  - Có URL → lưu vào `url`, auto-detect source (tiktok/facebook/web)
  - Có photo/document → upload vào storage, lưu URL vào `attachments`
- Tạo `inbox_items` với status='raw', source='telegram'
- Reply Telegram: `"✅ Đã lưu vào Inbox (ID: xxx)"`

**Acceptance criteria:**
- [ ] Telegram_id không hợp lệ → reply "Bạn chưa được cấp quyền"
- [ ] File > 20MB → reply "File quá lớn, vui lòng upload qua web"
- [ ] Response Telegram trong < 3s

#### FR-INBOX-02: Danh sách Inbox
**Priority:** P0

- `/inbox` — list items filter mặc định `status='raw'`
- Filter: status, source, có company/không
- Sort: created_at DESC
- Mỗi row: content preview, source badge, attachments icon, created_at relative, button "Triage"

#### FR-INBOX-03: Triage 1 item
**Priority:** P0

- Modal triage với options:
  - Loại: idea | task | content | knowledge | trash
  - Destination: tasks | content | knowledge | archive | delete
  - Company (optional)
  - Button "Convert" → tạo record trong bảng đích, update inbox_items.converted_to_*
- Status `raw` → `done`

**Acceptance criteria:**
- [ ] Chuyển sang tasks → pre-fill title từ inbox.content
- [ ] Archive/Delete → không tạo record mới, chỉ update status
- [ ] Có attachments → carry over sang record đích

#### FR-INBOX-04: Auto-archive sau 7 ngày
**Priority:** P1

- Cron daily 2am: `UPDATE inbox_items SET status='archived' WHERE status='raw' AND created_at < NOW() - INTERVAL '7 days'`

#### FR-INBOX-05: Manual add
**Priority:** P0

- Button "+ New" trên /inbox
- Form ngắn: content (required), url, type, company (optional)

---

### 6.4 Module 4: Tasks

#### FR-TASK-01: Tạo task
**Priority:** P0

- Form: title, description, company, assignee, status, priority, due_date, category, group_name, campaign
- Sau tạo → notification cho assignee (in-app + optional Telegram)

**Acceptance criteria:**
- [ ] Title và assignee bắt buộc
- [ ] Due_date không được < today (warning, không block)

#### FR-TASK-02: Kanban board
**Priority:** P0

- View Kanban cột theo status
- Drag-drop để đổi status
- Filter: company, assignee, priority, group_name

#### FR-TASK-03: My Tasks view
**Priority:** P0

- Dashboard personal: "Việc hôm nay" — filter `assignee_id=me AND due_date <= today AND status != done`
- Section "Quá hạn" màu đỏ
- Section "Sắp đến hạn" (3 ngày tới)

#### FR-TASK-04: Group by project view
**Priority:** P1

- Group by `group_name` — thay thế Projects DB
- Hiển thị progress bar: `done / total`

#### FR-TASK-05: Update status nhanh
**Priority:** P0

- Quick action: click status badge → dropdown để đổi
- Khi chuyển sang `done` → tự set `completed_at = NOW()`

#### FR-TASK-06: Calendar view
**Priority:** P1

- View calendar với task theo `due_date`
- Filter: company, assignee

---

### 6.5 Module 5: Content Pipeline

#### FR-CONT-01: Tạo content
**Priority:** P0

- Form: title, company, author, content_type, platforms, status (default 'idea'), topics, publish_date
- Source_file_url (Google Docs link) là trường chính cho text

#### FR-CONT-02: Pipeline board
**Priority:** P0

- Kanban theo status: Idea → Drafting → Review → Revise → Approved → Published
- Ẩn trạng thái `killed` mặc định
- Filter: company, author, content_type

#### FR-CONT-03: Review & feedback
**Priority:** P0

- Trạng thái `review` → Owner/Staff reviewer để lại `review_notes`
- Button "Approve" → status → `approved`
- Button "Request revision" → status → `revise` + notes bắt buộc

#### FR-CONT-04: Publish & track performance
**Priority:** P0

- Trạng thái `approved` → button "Mark as Published"
- Form điền: `published_url`, `publish_date`
- Sau publish → hiện form nhập Reach/Engagement (nhân sự tự điền sau)

#### FR-CONT-05: Published view (Content Library)
**Priority:** P0

- Gallery view: filter `status='published'`
- Sort by engagement_rate, reach, publish_date
- Click card → detail page

#### FR-CONT-06: Performance tracking
**Priority:** P1

- Bảng so sánh content tháng này — cột: title, platform, reach, engagement, rate
- Export CSV

---

### 6.6 Module 6: Work Logs (Module tracking output chính)

#### FR-WL-01: Submit work log
**Priority:** P0

- User vào `/work-logs/new`
- Form gọn:
  - **Ngày làm việc** (default = today, cho phép lùi 2 ngày)
  - **Công ty** (dropdown filter theo user_company_assignments)
  - **Campaign** (optional, filter theo company)
  - **Loại công việc** (từ work_types)
  - **Số lượng** (decimal, required)
  - **Link Google Sheet** (URL)
  - **Link ảnh/drive** (URL)
  - **Ảnh upload** (tối đa 5 ảnh)
  - **Ghi chú** (textarea)
- Button "Submit" → tạo record status='submitted'

**Acceptance criteria:**
- [ ] Không cho submit cùng user+date+work_type+campaign (unique constraint)
- [ ] File upload > 5MB → resize client-side hoặc báo lỗi
- [ ] Sau submit → redirect /work-logs với message success

#### FR-WL-02: Danh sách work logs cá nhân
**Priority:** P0

- `/work-logs` — list logs của user hiện tại
- Filter: month, company, work_type, status
- Cột: date, company, work_type, quantity+unit, status badge, notes preview

#### FR-WL-03: Owner/Staff review work logs
**Priority:** P0

- `/admin/work-logs/review` — filter mặc định `status='submitted'`
- Mỗi row có action:
  - ✅ Approve → status='approved'
  - ❌ Reject → status='rejected' + bắt buộc admin_notes
  - 🔄 Need fix → status='needs_fix' + admin_notes
- Batch approve: chọn nhiều + click "Approve selected"

**Acceptance criteria:**
- [ ] Chỉ user có `can_approve=true` trong user_company_assignments mới thấy button
- [ ] Reviewed log → lưu `reviewed_by` + `reviewed_at`

#### FR-WL-04: Dashboard output hàng tháng
**Priority:** P0

- Hiển thị tổng theo từng work_type tháng hiện tại
- Cards: 🔗 120 link · 📝 45 bài · 📢 8 campaign ...
- Click card → drill down theo company/người

#### FR-WL-05: Campaign progress view
**Priority:** P0

- Tab "Campaigns" trong Company Detail
- Mỗi campaign hiển thị:
  - Progress bar theo work_type: `187/300 (62%)` cho backlink
  - Days elapsed / total days
  - Status badge
- Indicator 🟢🟡🔴 theo % progress

#### FR-WL-06: User performance report
**Priority:** P1

- `/reports/performance` (chỉ Owner/Staff)
- Table: user, work_type, daily_avg, month_total, active_days
- Filter: month, company, work_type

#### FR-WL-07: Export CSV
**Priority:** P2

- Download toàn bộ work_logs theo filter
- Cột đầy đủ để phân tích ngoài

---

### 6.7 Module 7: Campaigns

#### FR-CAMP-01: Tạo campaign
**Priority:** P0

- Form: name, company, campaign_type, start_date, end_date, target_json (UI kéo thả chọn work_types + target)
- Budget (optional)

**UX for target_json:**
```
[+] Add target:
   Work type: [Backlink ▼]
   Target:    [300]
   Unit:      link (auto từ work_type)
   
[+] Add target:
   Work type: [Content ▼]
   Target:    [50]
   Unit:      bài
```

Lưu dưới dạng JSONB: `{"backlink": 300, "content": 50}`

#### FR-CAMP-02: Campaign detail dashboard
**Priority:** P0

- Header: name, status, date range
- Section "Progress": mỗi work_type 1 progress bar
- Section "Recent work logs": work_logs thuộc campaign này
- Section "Tasks": tasks gắn campaign
- Section "Content": content gắn campaign

#### FR-CAMP-03: List campaigns
**Priority:** P0

- Filter: company, status, campaign_type
- Card view với progress bar tổng

---

### 6.8 Module 8: Knowledge Base

#### FR-KB-01: Tạo knowledge item
**Priority:** P0

- Form: title, category, topics, scope, source_url, attachments
- Nếu scope='company_specific' → bắt buộc chọn company
- Nếu muốn write trong app → trường `body` dạng markdown editor

#### FR-KB-02: Browse knowledge
**Priority:** P0

- `/knowledge` — filter theo category, topic, quality_rating
- Default view: grid cards
- Click → detail page (render markdown nếu có body, link ngoài nếu source_url)

#### FR-KB-03: Full-text search
**Priority:** P1

- Search bar trên /knowledge
- Query GIN index với tsvector + unaccent (tiếng Việt)
- Highlight từ khoá trong kết quả

#### FR-KB-04: Share với company
**Priority:** P0

- Owner edit knowledge → multi-select "Visible to companies"
- Nhân sự company đó vào Company Detail tab Knowledge → thấy tài liệu được share

#### FR-KB-05: Rating (cho raw material)
**Priority:** P1

- User đánh giá 1-3 sao cho nguyên liệu đã dùng
- Filter "3 sao" để tìm nguyên liệu tốt nhất

---

### 6.9 Module 9: Dashboard

#### FR-DASH-01: Owner dashboard
**Priority:** P0

Layout:

```
┌────────────────────────────────────────────────┐
│ 📊 Tổng quan tháng 4/2026                       │
├────────────────┬────────────────┬──────────────┤
│ Work output    │ Content        │ Tasks        │
│ 🔗 320 link    │ 📝 42 bài      │ ✅ 58/72     │
│ 📢 15 camps    │ 📤 18 đã đăng  │ 🔴 3 overdue│
├────────────────┴────────────────┴──────────────┤
│ 🏢 Công ty (health status)                      │
│ ABC ✅  XYZ ⚠️  EDU ✅  FIN 🔴                   │
├─────────────────────────────────────────────────┤
│ 📥 Inbox (12 raw)    👀 Cần duyệt (5 work_logs)│
│                      👀 Content review (3)      │
├─────────────────────────────────────────────────┤
│ 📅 Việc hôm nay (8 tasks)                       │
│ 🔴 Quá hạn (2 tasks)                            │
└─────────────────────────────────────────────────┘
```

#### FR-DASH-02: Staff dashboard
**Priority:** P0

- Việc hôm nay (của staff)
- Việc mình manage (assigned bởi staff cho CTV)
- Content cần review
- Quick links to companies mình phụ trách

#### FR-DASH-03: CTV dashboard
**Priority:** P0

- Việc hôm nay
- Content đang viết
- Form quick submit work log
- Thông báo feedback từ review

#### FR-DASH-04: Client Staff dashboard
**Priority:** P0

- Giống CTV dashboard
- Thêm section Knowledge Base được share
- Không thấy work_logs/tasks của người khác

---

### 6.10 Module 10: Notion Sync

> Xem Section 12 — Notion Sync Specifications

---

## 7. Non-Functional Requirements

### 7.1 Performance

| Metric | Target | Notes |
|---|---|---|
| Dashboard load time | < 500ms | 95th percentile |
| Submit work_log | < 200ms | 95th percentile |
| Search knowledge | < 300ms | Với 1,000 records |
| Telegram webhook response | < 3s | Yêu cầu của Telegram |
| Notion sync 100 records | < 90s | Do rate limit 2 req/s |

### 7.2 Scalability

- Thiết kế cho 20 users, 10 companies ở MVP
- Có thể scale lên 50 users, 30 companies mà không cần đổi kiến trúc
- Database index cho các query hot path

### 7.3 Reliability

- Uptime target: **99.5%** (cho phép ~3.5h downtime/tháng)
- Backup PostgreSQL hàng ngày, giữ 7 ngày gần nhất
- Graceful degradation: Nếu Notion API down → app vẫn hoạt động, sync retry sau

### 7.4 Security

Xem Section 13.

### 7.5 Accessibility

- Desktop-first UI, responsive cho mobile
- Font size base 14–16px (dễ đọc)
- Contrast ratio đạt WCAG AA
- Keyboard navigation cho mọi action chính

### 7.6 Localization

- **MVP:** Tiếng Việt 100%
- UI terms, error messages, email templates đều tiếng Việt
- Date format: DD/MM/YYYY
- Number format: 1.234,56 (EU/VN style)
- Timezone: Asia/Ho_Chi_Minh (UTC+7)

### 7.7 Browser Support

- Chrome, Edge, Safari (2 bản mới nhất)
- Firefox (bản mới nhất)
- Không support IE11

---

## 8. User Flows

### 8.1 Flow: Owner capture ý tưởng qua Telegram

```
┌────────────────────────────────────────┐
│ 1. Owner lướt TikTok thấy bài hay       │
└─────────────┬──────────────────────────┘
              ▼
┌────────────────────────────────────────┐
│ 2. Copy link + mở Telegram Bot         │
│    "Bình Vương Capture Bot"            │
└─────────────┬──────────────────────────┘
              ▼
┌────────────────────────────────────────┐
│ 3. Paste link + gõ "idea content SEO"  │
│    + Send                              │
└─────────────┬──────────────────────────┘
              ▼
┌────────────────────────────────────────┐
│ 4. Bot reply:                          │
│    "✅ Lưu vào Inbox #a1b2"             │
└─────────────┬──────────────────────────┘
              ▼
┌────────────────────────────────────────┐
│ 5. Sáng hôm sau mở app → Inbox         │
│    thấy item chờ triage                │
└────────────────────────────────────────┘
```

### 8.2 Flow: CTV submit work log cuối ngày

```
┌─────────────────────────────────────────┐
│ 1. 18:00 — CTV hoàn thành công việc      │
└─────────────┬───────────────────────────┘
              ▼
┌─────────────────────────────────────────┐
│ 2. Mở app → Dashboard → Quick form      │
│    hoặc click "+ Work Log"              │
└─────────────┬───────────────────────────┘
              ▼
┌─────────────────────────────────────────┐
│ 3. Chọn Company: ABC Education           │
│    Chọn Campaign: SEO Q2                 │
│    Chọn Loại: Backlink                   │
│    Số lượng: 25                          │
│    Link sheet: docs.google.com/...       │
│    Ảnh: [2 ảnh screenshot]               │
│    Note: "Hôm nay 2 site từ chối"       │
└─────────────┬───────────────────────────┘
              ▼
┌─────────────────────────────────────────┐
│ 4. Submit                               │
│    → Record tạo với status='submitted'   │
│    → Quay về dashboard, thấy log mới    │
└─────────────┬───────────────────────────┘
              ▼
┌─────────────────────────────────────────┐
│ 5. Sáng hôm sau mở app → thấy           │
│    status đã chuyển 'approved' (Owner   │
│    đã review)                           │
└─────────────────────────────────────────┘
```

### 8.3 Flow: Owner review work logs sáng

```
┌────────────────────────────────────────┐
│ 1. 8:00 sáng — Owner mở Dashboard      │
└─────────────┬──────────────────────────┘
              ▼
┌────────────────────────────────────────┐
│ 2. Thấy notification "5 work logs      │
│    chờ duyệt"                          │
└─────────────┬──────────────────────────┘
              ▼
┌────────────────────────────────────────┐
│ 3. Click → /admin/work-logs/review     │
│    List 5 logs mới nhất                │
└─────────────┬──────────────────────────┘
              ▼
┌────────────────────────────────────────┐
│ 4. Với mỗi log:                        │
│    - Xem quantity, link sheet, ảnh     │
│    - Click ✅ Approve (hoặc batch)      │
└─────────────┬──────────────────────────┘
              ▼
┌────────────────────────────────────────┐
│ 5. Log status → 'approved'             │
│    → sẽ sync sang Notion lúc X:00      │
└────────────────────────────────────────┘
```

### 8.4 Flow: Owner setup Campaign mới

```
┌────────────────────────────────────────┐
│ 1. Company ABC bắt đầu dự án SEO Q2    │
└─────────────┬──────────────────────────┘
              ▼
┌────────────────────────────────────────┐
│ 2. Owner → Companies → ABC → Campaigns │
│    → Button "+ New Campaign"           │
└─────────────┬──────────────────────────┘
              ▼
┌────────────────────────────────────────┐
│ 3. Form:                               │
│    - Name: "SEO Q2 2026"               │
│    - Type: seo                         │
│    - Dates: 01/04 → 30/06              │
│    - Targets:                          │
│      🔗 Backlink: 300 link             │
│      📝 Content: 50 bài                │
│    - Budget: 30,000,000 VND            │
└─────────────┬──────────────────────────┘
              ▼
┌────────────────────────────────────────┐
│ 4. Save → Campaign detail page          │
│    Progress bar: 0/300, 0/50            │
└─────────────┬──────────────────────────┘
              ▼
┌────────────────────────────────────────┐
│ 5. Owner thông báo CTV:                │
│    "Tháng này đóng góp cho Campaign     │
│    SEO Q2 ABC"                         │
│    CTV khi submit work log → chọn       │
│    campaign này → tự cộng vào progress  │
└────────────────────────────────────────┘
```

---

## 9. Data Model

> Xem file `binhvuong-schema-postgres.md` cho schema chi tiết 12 bảng.

### 9.1 Tóm tắt

| Bảng | Rows dự kiến/năm | Purpose |
|---|---|---|
| users | 20–30 | Tài khoản |
| companies | 10–15 | Công ty |
| user_company_assignments | 30–50 | Phân quyền |
| inbox_items | 3,000+ | Ghi nhanh |
| tasks | 5,000+ | Công việc |
| content | 1,500+ | Nội dung |
| campaigns | 50–100 | Chiến dịch |
| work_types | 6–10 | Config |
| **work_logs** | **15,000+** | **Output hàng ngày** |
| knowledge_items | 500–1,000 | Tài liệu |
| objectives | 50–100 | OKR |
| notion_sync_log | 50,000+ | Debug (auto-cleanup) |

### 9.2 Quan hệ trọng tâm

```
users ◄─── user_company_assignments ───► companies
                                              │
                                    ┌─────────┼─────────┐
                                    ▼         ▼         ▼
                                  tasks    content    campaigns
                                              ▲         │
                                              │         ▼
                                              └────  work_logs ──► work_types
```

---

## 10. API Specifications

### 10.1 Conventions

- Base URL: `https://binhvuong.domain.com/api/v1`
- Auth: JWT via httpOnly cookie
- Content-Type: `application/json`
- Response format:

```json
{
  "success": true,
  "data": { ... },
  "error": null,
  "meta": {
    "page": 1,
    "total": 42
  }
}
```

Error:
```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Field 'quantity' must be greater than 0"
  }
}
```

### 10.2 Core Endpoints

#### Authentication
| Method | Endpoint | Purpose |
|---|---|---|
| POST | `/auth/login` | Login |
| POST | `/auth/logout` | Logout |
| POST | `/auth/refresh` | Refresh token |
| GET | `/auth/me` | Get current user |

#### Companies
| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/companies` | List companies (filtered by permissions) |
| POST | `/companies` | Create (owner only) |
| GET | `/companies/:id` | Get detail |
| PATCH | `/companies/:id` | Update |
| DELETE | `/companies/:id` | Soft delete |

#### Inbox
| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/inbox?status=raw` | List inbox items |
| POST | `/inbox` | Manual add |
| POST | `/inbox/:id/triage` | Triage item |
| DELETE | `/inbox/:id` | Delete |
| POST | `/telegram/webhook` | Telegram Bot webhook (no auth, HMAC verify) |

#### Tasks
| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/tasks` | List (filtered) |
| POST | `/tasks` | Create |
| GET | `/tasks/:id` | Detail |
| PATCH | `/tasks/:id` | Update |
| PATCH | `/tasks/:id/status` | Quick status change |
| DELETE | `/tasks/:id` | Soft delete |

#### Content
| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/content` | List |
| POST | `/content` | Create |
| GET | `/content/:id` | Detail |
| PATCH | `/content/:id` | Update |
| POST | `/content/:id/review` | Submit review (approve/revise) |
| POST | `/content/:id/publish` | Mark as published |
| PATCH | `/content/:id/metrics` | Update reach/engagement |

#### Work Logs
| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/work-logs` | List mine |
| POST | `/work-logs` | Submit new |
| GET | `/work-logs/:id` | Detail |
| PATCH | `/work-logs/:id` | Update (only if submitted) |
| DELETE | `/work-logs/:id` | Cancel (only if submitted) |
| POST | `/work-logs/:id/approve` | Approve (permission required) |
| POST | `/work-logs/:id/reject` | Reject |
| POST | `/work-logs/:id/need-fix` | Request fix |
| POST | `/work-logs/batch-approve` | Batch approve |
| GET | `/work-logs/summary` | Aggregate stats |

#### Campaigns
| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/campaigns` | List |
| POST | `/campaigns` | Create |
| GET | `/campaigns/:id` | Detail with progress |
| PATCH | `/campaigns/:id` | Update |
| GET | `/campaigns/:id/progress` | Get progress data |

#### Knowledge Base
| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/knowledge` | List (filtered) |
| POST | `/knowledge` | Create |
| GET | `/knowledge/:id` | Detail |
| PATCH | `/knowledge/:id` | Update |
| GET | `/knowledge/search?q=...` | Full-text search |

#### Dashboard
| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/dashboard/owner` | Owner dashboard data |
| GET | `/dashboard/staff` | Staff dashboard |
| GET | `/dashboard/user` | Personal dashboard |

---

## 11. UI/UX Requirements

### 11.1 Design Principles

1. **Tiếng Việt hoàn toàn** — mọi label, button, message
2. **Minimal clicks** — tối đa 2 clicks để đến bất kỳ action quan trọng
3. **Consistent layout** — header + sidebar + main content
4. **Visual hierarchy** — quan trọng nhất ở trên + to hơn
5. **Status by color** — 🟢 ok, 🟡 attention, 🔴 urgent, dùng nhất quán
6. **Dark mode** — P2, không bắt buộc ở MVP

### 11.2 Layout Structure

```
┌────────────────────────────────────────────────┐
│ Header: Logo │ Search │ Notif │ User Menu       │
├──────┬─────────────────────────────────────────┤
│      │                                          │
│ Side │  Main Content Area                       │
│ bar  │                                          │
│      │  Breadcrumb                              │
│ • 🏠 │  Page Title                              │
│ • 📥 │  ─────────────────────────────           │
│ • ✅ │                                          │
│ • 📝 │  Content...                              │
│ • 📚 │                                          │
│ • 🏢 │                                          │
│      │                                          │
└──────┴─────────────────────────────────────────┘
```

**Sidebar items theo role:**

| Role | Sidebar items |
|---|---|
| Owner | Dashboard · Inbox · Tasks · Content · Work Logs · Campaigns · Knowledge · Companies · People · Reports |
| Staff | Dashboard · Inbox · Tasks · Content · Work Logs · Campaigns · Knowledge · Companies |
| CTV | Dashboard · My Tasks · Content · Work Logs · Knowledge |
| Client Staff | Dashboard · My Tasks · Work Logs · Knowledge (filtered) |

### 11.3 Color System

```css
--primary: #4F46E5;     /* Indigo — brand color */
--success: #10B981;     /* Green */
--warning: #F59E0B;     /* Amber */
--danger:  #EF4444;     /* Red */
--info:    #3B82F6;     /* Blue */
--gray-50: #F9FAFB;
--gray-900: #111827;
```

### 11.4 Typography

- **Display font:** Be Vietnam Pro (tiếng Việt đẹp, miễn phí Google Fonts)
- **Body font:** Be Vietnam Pro
- **Monospace:** JetBrains Mono (cho code, UUIDs)
- **Base size:** 14px
- **Line height:** 1.5

### 11.5 Forms — Rules

1. Label trên input (không bên cạnh)
2. Required field có dấu `*` đỏ
3. Error message ngay dưới field, màu đỏ
4. Help text màu xám nhỏ dưới label
5. Submit button luôn bên phải (hoặc full width mobile)
6. Cancel bên trái (outline style)
7. Auto-save draft cho form dài (work log, content)

### 11.6 Components sử dụng

- **shadcn/ui adapt cho templ** — hoặc viết component lại
- HTMX cho interactivity:
  - `hx-post` cho form submit
  - `hx-get` cho load dynamic content
  - `hx-swap="outerHTML"` cho replace
  - Indicator loading state

### 11.7 Responsive Breakpoints

| Breakpoint | Target | Note |
|---|---|---|
| < 640px | Mobile | Sidebar → drawer, full-width forms |
| 640 – 1024px | Tablet | Sidebar thu gọn icon |
| > 1024px | Desktop | Full sidebar mở |

---

## 12. Notion Sync Specifications

### 12.1 Scope

- **Sync direction:** 1 chiều (Postgres → Notion)
- **Frequency:** Mỗi 1 giờ (cron)
- **Tables to sync:** companies, users, tasks, content, work_logs, campaigns, knowledge_items, objectives
- **Tables NOT sync:** inbox_items (staging), notion_sync_log (debug), user_company_assignments (internal)

### 12.2 Sync Worker Logic

```python
# Pseudocode
def sync_to_notion():
    rate_limiter = RateLimiter(2, 1)  # 2 req/sec
    
    for table in SYNC_TABLES:
        records = db.query(f"""
            SELECT * FROM {table}
            WHERE sync_status IN ('pending', 'error')
               OR updated_at > synced_at
            LIMIT 100
        """)
        
        for record in records:
            rate_limiter.wait()
            try:
                if record.notion_page_id is None:
                    page = notion.create_page(mapper(record))
                    db.update(table, record.id, {
                        'notion_page_id': page.id,
                        'synced_at': now(),
                        'sync_status': 'synced'
                    })
                else:
                    notion.update_page(record.notion_page_id, mapper(record))
                    db.update(table, record.id, {
                        'synced_at': now(),
                        'sync_status': 'synced'
                    })
                log_sync(table, record.id, 'success')
            except RateLimitError:
                rate_limiter.backoff()
            except Exception as e:
                db.update(table, record.id, {
                    'sync_status': 'error',
                    'sync_error': str(e)
                })
                log_sync(table, record.id, 'error', str(e))
```

### 12.3 Mapper Specifications

Mỗi bảng cần 1 mapper function `postgresRow → notionProperties`.

**Ví dụ mapper cho `companies`:**

```go
func MapCompanyToNotion(c Company) map[string]interface{} {
    return map[string]interface{}{
        "Tên công ty": map[string]interface{}{
            "title": []map[string]interface{}{{"text": map[string]interface{}{"content": c.Name}}},
        },
        "Mã": richText(c.ShortCode),
        "Ngành": selectOption(mapIndustry(c.Industry)),
        "Vai trò": selectOption(mapRole(c.MyRole)),
        "Trạng thái": selectOption(mapStatus(c.Status)),
        "Phạm vi": multiSelect(c.Scope),
        "Tình trạng": selectOption(mapHealth(c.Health)),
        "Liên hệ chính": richText(fmt.Sprintf("%s - %s", c.ContactName, c.ContactZalo)),
        "Ngày bắt đầu": dateField(c.StartDate),
        "Ghi chú": richText(c.Description),
    }
}
```

### 12.4 Setup Notion

**Owner làm 1 lần:**

1. Tạo Notion Integration: notion.so/my-integrations → New Integration
2. Copy secret key → `.env` của app Go
3. Tạo các database trong Notion với property schema khớp mapper
4. Share mỗi database với Integration (click Share → add connection)
5. Copy database_id → config của app

### 12.5 Rate Limit Handling

API rate limit là 3 requests/giây trung bình (2,700/15 phút). Khi vượt quá sẽ bị HTTP 429.

**App dùng 2 req/s (an toàn):**
- Max ~120 records/phút
- Max ~7,200 records/giờ
- Với 100 records cập nhật/giờ → chỉ dùng ~1.5% capacity

**Handling 429:**
- Parse header `Retry-After`
- Sleep + retry (max 3 lần)
- Nếu vẫn lỗi → mark `sync_status='error'`, chờ lần sync sau

### 12.6 Error Recovery

| Error Type | Action |
|---|---|
| Network timeout | Retry 3 lần với exponential backoff |
| 429 Rate limit | Wait Retry-After header + retry |
| 400 Validation | Log, mark error, skip đến lần sync sau |
| 401 Auth | Alert Owner qua Telegram, stop worker |
| 404 Database not found | Log error, skip record |
| 500 Notion server | Retry 3 lần, skip nếu vẫn lỗi |

---

## 13. Security & Permissions

### 13.1 Authentication

- **Password:** Bcrypt với cost=12
- **Session:** JWT trong httpOnly cookie, Secure flag, SameSite=Lax
- **Token lifetime:** 7 ngày (30 nếu remember me)
- **Rate limit login:** 5 attempts / 15 phút per IP

### 13.2 Authorization Matrix

| Resource | Owner | Staff | CTV | Client Staff |
|---|---|---|---|---|
| Create company | ✅ | ❌ | ❌ | ❌ |
| Update company | ✅ | ✅ (assigned) | ❌ | ❌ |
| Create user | ✅ | ❌ | ❌ | ❌ |
| Create task | ✅ | ✅ (assigned cty) | ❌ | ❌ |
| View task | ✅ all | ✅ (assigned cty) | ✅ (mine only) | ✅ (mine only) |
| Approve work_log | ✅ all | ✅ (can_approve=true) | ❌ | ❌ |
| Submit work_log | ✅ | ✅ | ✅ | ✅ |
| View Knowledge | ✅ all | ✅ (shared+scope) | ✅ (shared+scope) | ✅ (shared+scope) |
| Create Knowledge | ✅ | ✅ | ❌ | ❌ |
| View reports | ✅ | ✅ | ❌ | ❌ |

### 13.3 Data Protection

- Database connection: TLS required
- Backups: Encrypted at rest (GPG)
- Passwords: Never logged, never returned in API
- PII: Không có PII nhạy cảm (credit card, ID number)
- File uploads: Scan size, validate MIME type, rename to UUID

### 13.4 OWASP Top 10 Checklist

- [x] **Injection:** Dùng parameterized queries (pgx/sqlc)
- [x] **Broken Auth:** Bcrypt + JWT + rate limit
- [x] **Sensitive Data:** HTTPS + no secrets in logs
- [x] **XXE:** Không parse XML từ user
- [x] **Broken Access Control:** Middleware check permission mọi endpoint
- [x] **Security Misconfig:** Caddy security headers
- [x] **XSS:** templ auto-escape + CSP header
- [x] **Insecure Deserialization:** Validate JSON schema
- [x] **Known Vulnerabilities:** `go mod tidy` + dependabot
- [x] **Insufficient Logging:** Structured logs với slog

### 13.5 Telegram Webhook Security

- Verify HMAC signature với secret token
- Whitelist IP của Telegram (optional)
- Verify `from.id` khớp với `users.telegram_id`

---

## 14. Roadmap & Milestones

### 14.1 Timeline Overview

```
Week 1-2: Foundation  ─── MVP Core (Week 3-6) ─── Polish (Week 7-8) ─── Post-MVP (Week 9-12)
```

### 14.2 Tuần 1–2: Foundation

**Goal:** Setup hạ tầng + auth + users

- [ ] Setup Docker Compose stack (Postgres + Go + Caddy)
- [ ] Tạo migrations đầu tiên (users, companies, user_company_assignments)
- [ ] Implement auth (login, logout, JWT middleware)
- [ ] Tạo template base UI (header, sidebar, layout với templ)
- [ ] CRUD Companies + Users (Owner only)
- [ ] Deploy preview lên VPS

**Deliverable:** Owner đăng nhập được, tạo công ty, tạo user, phân quyền.

### 14.3 Tuần 3–4: Inbox + Tasks

**Goal:** Core flow capture + assign

- [ ] Inbox module (CRUD + triage flow)
- [ ] Telegram bot integration
- [ ] Tasks module (CRUD + Kanban view)
- [ ] My Tasks dashboard
- [ ] Notifications cơ bản (in-app)

**Deliverable:** Owner có thể capture qua Telegram, triage, tạo task, CTV thấy task của mình.

### 14.4 Tuần 5: Content

**Goal:** Content pipeline hoàn chỉnh

- [ ] Content CRUD
- [ ] Pipeline Kanban
- [ ] Review & feedback flow
- [ ] Publish & metrics tracking
- [ ] Content Library view

**Deliverable:** CTV submit content, Owner review, publish, track engagement.

### 14.5 Tuần 6: Work Logs (MODULE QUAN TRỌNG NHẤT)

**Goal:** Tracking output hàng ngày

- [ ] work_types seed + admin UI
- [ ] Campaigns CRUD
- [ ] Work logs submit form (mobile-friendly)
- [ ] Review flow (approve/reject)
- [ ] Dashboard output monthly
- [ ] Campaign progress view

**Deliverable:** 20 nhân sự submit log cuối ngày, Owner review, dashboard hiện tổng.

### 14.6 Tuần 7: Knowledge + Dashboard

**Goal:** Hoàn thiện trải nghiệm

- [ ] Knowledge Base CRUD
- [ ] Full-text search
- [ ] Share với company flow
- [ ] Owner dashboard hoàn chỉnh
- [ ] Role-based dashboards

**Deliverable:** Knowledge accessible, dashboard đầy đủ số liệu.

### 14.7 Tuần 8: Notion Sync + Polish

**Goal:** Production-ready

- [ ] Notion sync worker
- [ ] Mappers cho tất cả bảng
- [ ] Error handling + retry logic
- [ ] Performance tuning (index, query)
- [ ] Bug fixes từ internal testing
- [ ] Documentation cho user

**Deliverable:** MVP go-live cho 5 users đầu tiên (soft launch).

### 14.8 Tuần 9–12: Post-MVP

**Goal:** Hoàn thiện dựa trên feedback

- [ ] Performance tracking + export CSV
- [ ] Forgot password flow
- [ ] Email notifications
- [ ] Mobile responsive polish
- [ ] Advanced reports
- [ ] Onboarding flow cho user mới
- [ ] Full rollout 20 users

---

## 15. Success Metrics

### 15.1 Adoption Metrics

| Metric | Target (3 tháng) | Cách đo |
|---|---|---|
| Daily active users | 15/20 (75%) | Login trong ngày |
| Work log submission rate | 90% weekdays | Logs / expected logs |
| Inbox triage time | < 24h median | Time from create to process |
| Task completion rate | 80% ontime | Done before due_date / total |

### 15.2 Business Metrics

| Metric | Target | Note |
|---|---|---|
| Thời gian quản trị hàng ngày của Owner | Giảm 60% | Từ 2–3h → 45 phút |
| Accuracy của báo cáo output | 95%+ | Verified reports / total |
| Số lượng công cụ dùng song song | Giảm từ 6 → 2 | App + Notion |

### 15.3 Technical Metrics

| Metric | Target |
|---|---|
| Uptime | > 99.5% |
| Dashboard load time (p95) | < 500ms |
| Notion sync success rate | > 95% |
| Zero critical bugs in production | ✓ |

---

## 16. Risks & Mitigations

| # | Risk | Probability | Impact | Mitigation |
|---|---|---|---|---|
| R1 | Go learning curve cho bạn | Medium | High | Hoặc thuê dev Go 2–3 tháng, hoặc dùng framework dễ hơn (Node.js Fastify) |
| R2 | Notion API thay đổi breaking | Low | Medium | Version lock, monitor changelog, có fallback không sync |
| R3 | User refuse dùng app mới | High | High | Onboarding kỹ, UX cực đơn giản cho form submit, giữ Google Sheets song song |
| R4 | VPS bị hỏng/mất data | Low | Critical | Daily backup + offsite backup (S3/R2), disaster recovery plan |
| R5 | Scope creep | High | Medium | Strict MVP scope, parking lot cho ý tưởng mới |
| R6 | Security breach | Low | Critical | Security checklist, pen test trước production, rotate secrets |
| R7 | Performance degrade khi data lớn | Medium | Medium | Index đúng từ đầu, query profiling, pagination everywhere |
| R8 | Sync Notion lỗi → mất data | Low | Low | Postgres là source of truth, Notion chỉ là mirror |

---

## 17. Out of Scope

### 17.1 Tính năng KHÔNG làm ở MVP

- ❌ **Mobile app native** — dùng responsive web
- ❌ **Chat nội bộ** — giữ Zalo/Telegram
- ❌ **Video call** — dùng Google Meet/Zoom
- ❌ **File storage như Google Drive** — chỉ upload nhỏ vào app, file lớn link ngoài
- ❌ **Accounting/Invoice** — dùng tool riêng
- ❌ **CRM khách hàng cuối** — không phải scope
- ❌ **Email marketing** — không phải scope
- ❌ **Timesheet chấm công** — không phải scope
- ❌ **AI features** (summarize, auto-classify) — post-MVP
- ❌ **Sync 2 chiều với Notion** — phức tạp, rủi ro conflict
- ❌ **Public API cho external** — chỉ internal use
- ❌ **White-label cho client** — mỗi client có subdomain riêng, không làm
- ❌ **Multi-tenant** — single workspace duy nhất
- ❌ **Real-time collaboration** (like Google Docs) — không cần
- ❌ **Audit log chi tiết** — chỉ log sync, không log mọi action

### 17.2 Có thể cân nhắc Post-MVP (Phase 2)

- Telegram bot 2 chiều (nhận lệnh phức tạp hơn)
- AI auto-classify Inbox
- Google Sheets API để verify work_log quantity
- Analytics nâng cao + custom reports
- API cho Zapier/Make.com
- Dark mode
- Multi-language (English)
- Export PDF reports cho client
- Gantt chart cho campaigns
- Time tracking trên task (optional)

---

## Phụ Lục

### A. Glossary

| Thuật ngữ | Định nghĩa |
|---|---|
| Solo Founder | Người sáng lập một mình, không có co-founder cùng vai trò |
| CTV | Cộng tác viên, làm part-time |
| Client Staff | Nhân viên của công ty khách mà Owner mentor |
| Work Log | Báo cáo output hàng ngày của 1 nhân sự |
| Campaign | Nhóm công việc có mục tiêu và thời hạn |
| Triage | Quá trình phân loại Inbox item vào các module |
| Source of Truth | Nguồn dữ liệu chính xác duy nhất |
| Sync 1 chiều | Đồng bộ một hướng, không có conflict resolution |

### B. Tài liệu liên quan

- `binhvuong-schema-postgres.md` — Database schema chi tiết
- `notion-schema-toi-gian.md` — Schema Notion workspace (bản mirror)
- Docker Compose template (TBD)
- Go project skeleton (TBD)

### C. Changelog

| Version | Date | Author | Changes |
|---|---|---|---|
| 1.0 | 16/04/2026 | AI Assistant | Initial draft |

---

**END OF DOCUMENT**

*Tài liệu này là bản thiết kế hoàn chỉnh. Trước khi bắt đầu development, cần review kỹ với team kỹ thuật và xác nhận lại scope MVP.*
