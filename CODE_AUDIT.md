# Code Audit Report

> **Last Updated:** 2026-01-08
> **Purpose:** Identify code duplication, consolidation opportunities, and architectural improvements

This document follows the core principles from `CONTRIBUTING.md` to audit the codebase for issues.

---

## Summary

| Category | Issues Found | Severity |
|----------|-------------|----------|
| ✅ **Good Practices Followed** | 12 | — |
| ⚠️ **Minor Issues** | 6 | Low |
| 🔶 **Consolidation Opportunities** | 4 | Medium |
| ❌ **Duplications** | 3 | Medium |

---

## ✅ Good Practices Already Followed

The codebase already follows many good practices:

### Backend

1. **Centralized response helpers** - `internal/handler/response.go` has `writeJSON`, `writeError` 
2. **Centralized error handling** - `internal/handler/errors.go` has `HandleServiceError`
3. **Centralized constants** - `internal/config/constants.go` has all timeouts and magic numbers
4. **Filesystem abstraction** - All file operations use `filesystem.FS` interface
5. **Path validation** - `internal/pkg/validator/` provides path validation
6. **File utilities** - `internal/pkg/fileutil/` has shared `ToFileInfo` and `DetectMimeType`

### Frontend

7. **Centralized config** - `src/lib/config.ts` has all magic numbers and intervals
8. **Centralized storage** - `src/lib/utils/storage.ts` abstracts localStorage access
9. **Centralized file types** - `src/lib/utils/fileTypes.ts` is the single source of truth for extensions
10. **API client abstraction** - `src/lib/api/client.ts` handles auth, refresh, and errors
11. **Format utilities** - `src/lib/utils/format.ts` has file size, date formatting
12. **Type exports** - Types are re-exported from `src/lib/types/files.ts`

---

## ❌ Duplications Found

### 1. ~~**MIME Type Detection Duplicated** (Backend)~~ ✅ FIXED

**Status:** ✅ Fixed - Duplicate removed from `stream.go`, now uses `fileutil.DetectMimeType`

**What was done:**
- Removed `detectMimeType` method from `internal/handler/stream.go`
- Added import for `internal/pkg/fileutil`
- Updated `Download` and `Preview` handlers to use `fileutil.DetectMimeType`

---

### 2. **Direct localStorage Access in API Files** (Frontend)

**Severity:** Medium

**Storage abstraction exists at:** `src/lib/utils/storage.ts`

**But direct access still in:** `src/lib/api/files.ts` (lines 219, 230)
```typescript
const token = typeof window !== 'undefined' ? localStorage.getItem('accessToken') : null;
```

**Problem:** Storage is centralized in `tokenStorage` but not used everywhere.

**Fix:**
- Import from storage utilities:
```typescript
import { tokenStorage } from '$lib/utils/storage';
const token = tokenStorage.getAccessToken();
```

---

### 3. **Type Definitions in API vs Types Folder** (Frontend)

**Severity:** Low

**`src/lib/api/files.ts` defines:**
- `FileInfo`, `FileList`, `MountPoint`, `DriveStats`, `ListOptions`, `SearchResponse`, etc.

**`src/lib/types/files.ts` defines:**
- `SortField`, `SortDir` (but these are ALSO inline in `ListOptions`)

**Problem:** Types are split between API module and types folder.

**Fix:**
- Move all shared interfaces to `src/lib/types/files.ts`
- Re-export from API for backward compatibility
- Use `SortField` and `SortDir` from types in `ListOptionsState`

---

## 🔶 Consolidation Opportunities

### 1. **API Object Exports Are Redundant** (Frontend)

**Location:** `src/lib/api/files.ts` (lines 199-212)

The file exports both individual functions AND an object:
```typescript
export async function listRoots() { ... }
export async function getDriveStats() { ... }
// ... more functions

// ALSO exports object (redundant)
export const filesApi = {
    listRoots,
    getDriveStats,
    // ...
};
```

**Issue:** Violates DRY - two ways to access the same functions.

**Recommendation:**
- Keep only named exports (more tree-shakeable)
- Remove `filesApi`, `jobsApi`, `authApi` objects

---

### 2. **Stores Still Use Svelte 4 Patterns** (Frontend)

**Locations:**
- `src/lib/stores/auth.ts` - Uses `writable`, `derived`, `get`
- `src/lib/stores/files.ts` - Uses `writable`, `derived`, `get`
- `src/lib/stores/jobs.ts` - Uses `writable`, `derived`, `get`
- `src/lib/stores/settings.ts` - Uses `writable`, `derived`, `get`
- `src/lib/stores/websocket.ts` - Uses `writable`, `derived`, `get`

**Already migrated:**
- `src/lib/stores/clipboard.svelte.ts` - Uses Svelte 5 `$state`

**Issue:** Inconsistent patterns - some stores use old Svelte 4, some use new Svelte 5.

**Recommendation:**
- Migrate all stores to `.svelte.ts` extension
- Use `$state`, `$derived` instead of `writable`, `derived`

---

### 3. **http.Error vs writeError Inconsistency** (Backend)

**`writeError` helper exists at:** `internal/handler/response.go`

**But `http.Error` still used in:** `internal/handler/websocket.go` (lines 100, 107)
```go
http.Error(w, "Missing authentication token", http.StatusUnauthorized)
http.Error(w, "Invalid authentication token", http.StatusUnauthorized)
```

**Issue:** Inconsistent error response format (plain text vs JSON).

**Recommendation:**
- Use `writeError` for JSON consistency, OR
- Document that WebSocket upgrade errors are intentionally plain text

---

### 4. **UploadManager Could Be a Service** (Backend)

**Location:** `internal/handler/stream.go` (lines 161-272)

**Issue:** `UploadManager` is embedded in handler package but has complex logic:
- Session management
- Cleanup goroutine
- Temp file handling

**Recommendation:**
- Consider moving to `internal/service/upload.go`
- Handler would just delegate to service
- Follows handler → service → model pattern

---

## ⚠️ Minor Issues

### 1. **Empty src/lib/index.ts**

**Location:** `src/lib/index.ts`

File exists but only has 326 bytes (minimal exports).

**Recommendation:** Either populate with re-exports or document its purpose.

---

### 2. **SortField/SortDir Inline in Store**

**Location:** `src/lib/stores/files.ts` (lines 107-108)

```typescript
sortBy: 'name' | 'size' | 'modTime' | 'type';
sortDir: 'asc' | 'desc';
```

**Exists in types:** `src/lib/types/files.ts`
```typescript
export type SortField = 'name' | 'size' | 'modTime' | 'type';
export type SortDir = 'asc' | 'desc';
```

**Fix:** Import from types instead of inline definition.

---

### 3. **formatSize Removed but Check for Remnants**

The `format.ts` re-exports from `fileTypes.ts`:
```typescript
export { getFileTypeDescription } from '$lib/utils/fileTypes';
```

**Note:** Good consolidation already done. No issues here.

---

### 4. **CONFIG.auth.tokenRefreshIntervalMs Well-Used**

**Location:** `src/lib/stores/auth.ts` (line 72)
```typescript
}, CONFIG.auth.tokenRefreshIntervalMs);
```

**Note:** Good practice - using centralized config.

---

### 5. **Rate Limit Cleanup Missing StartCleanup Call**

**Location:** `internal/middleware/ratelimit.go`

The `RateLimiter` has `StartCleanup` and `cleanup` methods, but `RateLimit` function creates a new limiter without starting cleanup.

**Current:**
```go
func RateLimit(rps float64) func(http.Handler) http.Handler {
    limiter := NewRateLimiter(rps, int(rps*2))
    // Missing: limiter.StartCleanup(context.Background())
    return limiter.Limit
}
```

**Fix:** Start cleanup or document that cleanup is optional for short-lived limiters.

---

### 6. **FileBrowser.svelte Usage**

**Location:** `src/lib/components/FileBrowser.svelte`

**Issue:** As noted in frontend refactor plan, this component may not be used. Need to verify.

**Recommendation:** Check if used; delete if not.

---

## Where Things Should Live

Following the "Where should this live?" principle:

| What | Current Location | Correct? |
|------|-----------------|----------|
| MIME detection | `pkg/fileutil/` | ✅ Fixed - duplicate removed |
| Upload session management | `handler/stream.go` | 🔶 Could move to service |
| File type constants | `src/lib/utils/fileTypes.ts` | ✅ Correct |
| API URL base | `src/lib/api/client.ts` | ✅ Correct |
| Storage keys | `src/lib/config.ts` | ✅ Correct |
| Backend constants | `internal/config/constants.go` | ✅ Correct |

---

## Action Items

### High Priority

1. [x] ~~**Remove duplicate MIME detection** from `stream.go`, use `fileutil.DetectMimeType`~~ ✅ Done
2. [ ] **Fix localStorage access** in `api/files.ts` to use `tokenStorage`
3. [ ] **Start rate limiter cleanup** or document cleanup strategy

### Medium Priority

4. [ ] **Move types** from `api/files.ts` to `types/files.ts`
5. [ ] **Remove redundant API objects** (`filesApi`, etc.)
6. [ ] **Migrate stores** to Svelte 5 runes (auth, files, jobs, settings, websocket)

### Low Priority

7. [ ] **Consider moving UploadManager** to service layer
8. [ ] **Use consistent error format** in WebSocket handler
9. [ ] **Import SortField/SortDir** from types in stores

---

## Verification Commands

### Backend

```bash
# Check for duplicate implementations
rg "filepath.Ext" --type go backend/

# Verify all handlers use writeError
rg "http\.Error" --type go backend/internal/handler/

# Check response.go usage
rg "writeJSON|writeError|HandleServiceError" --type go backend/internal/handler/
```

### Frontend

```bash
# Check for direct localStorage access
rg "localStorage\." frontend/src/lib/ --type ts

# Check for old Svelte store patterns
rg "writable|derived" frontend/src/lib/stores/ --type ts

# Verify fileTypes.ts is used
rg "from.*fileTypes" frontend/src/ --type ts --type svelte
```
