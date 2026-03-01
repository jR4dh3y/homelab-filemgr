# File Manager - Feature Roadmap & TODO

> **Last Updated:** 2026-01-08

## ✅ Recently Completed

- [x] **Context Menu (Right-click)** - Copy, Cut, Paste, Rename, Delete, Download, Properties, Open With
- [x] **Backend Security Refactoring** - All security features implemented:
  - [x] Configurable user credentials (via config/env vars)
  - [x] Rate limiting middleware for auth endpoints
  - [x] Configurable WebSocket origins
  - [x] Cross-platform build support (Linux/Windows)

---

## High Priority

- [ ] **Keyboard Shortcuts** - Ctrl+C/V/X/H(toggle hidden files), Delete, F2 (rename), Enter (open), Backspace (go up)
- [ ] **Drag & Drop Upload** - Drop files anywhere to upload to current folder
- [ ] **Search Bar** - Quick search in current folder or global search
- [ ] **Breadcrumb Quick Navigation** - Click any folder in path to jump there

## Medium Priority

- [ ] **Multi-select Actions** - Bulk delete, download as zip, move/copy multiple files
- [ ] **File/Folder Creation** - New folder button, new text file
- [ ] **Rename Inline** - Click on filename to rename in place
- [ ] **Sort Persistence** - Remember sort preferences per folder
- [ ] **Quick Preview on Hover** - Thumbnail preview for images on hover

## Nice to Have

- [ ] **Favorites/Bookmarks** - Pin frequently accessed folders
- [ ] **Recent Files** - Show recently opened/modified files
- [ ] **Dual Pane View** - Split view for easy file moving
- [ ] **File Info Panel** - Side panel showing file details, EXIF data for images
- [ ] **Dark/Light Theme Toggle** - User preference for theme
- [ ] **Grid View** - Icon/thumbnail grid view for images folder
- [ ] **Clipboard History** - Track recent copy operations
- [ ] **Progress Indicator** - Show upload/download progress in status bar

---

## Documentation Discrepancies (To Fix)

These are inconsistencies found between documentation and actual code:

- [ ] **Config loader users parsing** - The config loader (`internal/config/config.go`) doesn't currently parse `FM_USERS_*` environment variables. Need to implement user loading from env vars with prefix pattern matching.
- [ ] **Allowed origins env parsing** - The config loader doesn't parse `FM_ALLOWED_ORIGINS` as comma-separated list. Need to implement string splitting.
- [ ] **Go version mismatch** - `docs/development.md` says Go 1.23+, but `go.mod` was upgraded to Go 1.24.0. Should align documentation.
- [ ] **Viper env prefix** - Need to configure Viper to use `FM_` prefix for all environment variables as documented.

---

## Frontend Refactoring (from frontend/REFACTOR_PLAN.md)

See `frontend/REFACTOR_PLAN.md` for full details. Summary of remaining work:

### Phase 1: Foundation
- [ ] Create `src/lib/config.ts` with centralized constants
- [ ] Create `src/lib/utils/storage.ts` for localStorage access
- [ ] Consolidate `fileTypes.ts` as single source of truth
- [ ] Remove duplicate `formatSize` from `fileTypes.ts`
- [ ] Delete unused files (`.gitkeep`, empty `index.ts`)

### Phase 2: Design System
- [ ] Add design tokens to `layout.css`
- [ ] Create `src/lib/components/ui/` base components
- [ ] Create `Button.svelte`, `Spinner.svelte`, `Card.svelte`

### Phase 3: Component Migration
- [ ] Convert all CSS `<style>` blocks to Tailwind classes
- [ ] Start with small components, then medium, then large
- [ ] Update routes last

### Phase 4: Store Migration
- [ ] Convert stores to Svelte 5 runes (`.svelte.ts`)
- [ ] Update all imports

---

## Backend Technical Debt

These are minor issues identified during refactoring:

- [ ] **Markdownlint warnings** - Several markdown files have formatting issues (blanks around lists/fences)
- [ ] **Test coverage** - Add more unit tests for new middleware (rate limiting, origin checking)
- [ ] **bcrypt password hashing** - Currently passwords are stored in plain text config. For production, implement bcrypt hashing.
- [ ] **Config validation** - Add validation for new config fields (users, rate_limit_rps, allowed_origins)

---

## Notes

### Backend Refactoring Complete ✅

All tasks from `backend/REFACTOR_PLAN.md` are complete:
- Empty placeholder files → Added package documentation
- Code duplication → Consolidated into shared utilities
- Architecture issues → Centralized constants, error handling, response helpers
- Resource management → Added cleanup for auth tokens, upload sessions, job history
- Security issues → Configurable credentials, rate limiting, WebSocket origins

---

## Integration Drift Audit (2026-03-01)

Context: A backend/frontend linkage audit was done after stream/upload fixes were shipped.  
Goal: catch places where backend features exist but frontend behavior is missing, stale, or not wired correctly.

### High Priority

- [ ] **Upload ID mismatch between store and uploader**
  - Context: `uploadStore` generates `uploadId`, but `uploadFile()` generates a different one internally.
  - Impact: progress tracking/resume wiring can drift and become unreliable.
  - References:
    - `frontend/src/lib/stores/upload.svelte.ts` (queue/progress lifecycle)
    - `frontend/src/lib/utils/upload.ts` (`uploadFile()` internal `generateUploadId()`)
  - Fix direction: pass `uploadId` from store into uploader and remove duplicate ID generation.

- [ ] **Resumable upload feature exists but is not integrated in UI flow**
  - Context: `resumeUpload()` and `/stream/upload/status` are implemented, but normal upload path only calls `uploadFile()`.
  - Impact: interrupted uploads restart instead of resuming.
  - References:
    - `frontend/src/lib/utils/upload.ts` (`resumeUpload`, `getUploadStatus`)
    - `frontend/src/lib/stores/upload.svelte.ts` (only `uploadFile()` used)
  - Fix direction: detect prior upload session and route through `resumeUpload()` when available.

- [ ] **WebSocket backend capability not fully consumed by frontend**
  - Context: backend can batch multiple JSON messages in one WS frame (newline-delimited), frontend assumes one JSON object per frame.
  - Impact: dropped/parse-failed real-time updates under burst conditions.
  - References:
    - `backend/internal/websocket/client.go` (`WritePump` writes newline-delimited messages)
    - `frontend/src/lib/stores/websocket.ts` (`JSON.parse(event.data)` single parse path)
  - Fix direction: split incoming payload by newline and parse each JSON message safely.

- [ ] **WebSocket store not actively wired into page lifecycle**
  - Context: connection helpers exist in `websocketStore`, but app-level usage of `connect()`/subscriptions is not clearly wired from routes.
  - Impact: real-time job updates may silently rely on polling only.
  - References:
    - `frontend/src/lib/stores/websocket.ts`
    - route/layout integration points
  - Fix direction: connect/disconnect in authenticated app lifecycle and subscribe job IDs when jobs are active.

### Medium Priority

- [ ] **Streaming API docs are stale vs actual response contract**
  - Context: docs still describe old upload response shape (`received`, `size`, `checksum`) while handler returns `receivedChunks`, `totalChunks`, `complete`, `path`.
  - Impact: client implementers use wrong schema.
  - References:
    - `docs/api.md` (stream upload section)
    - `backend/internal/handler/stream.go` (`UploadResponse`)
  - Fix direction: update docs to current endpoint behavior and response payloads.

- [ ] **Stream property tests use non-production route prefix**
  - Context: tests exercise `/api/v1/upload/*`, while production mounts stream routes under `/api/v1/stream/*`.
  - Impact: route-level regressions can slip through tests.
  - References:
    - `backend/internal/handler/stream_property_test.go`
    - `backend/cmd/server/main.go` route mounting
  - Fix direction: align test request paths with production route prefix.
