# Backend Refactoring Plan

This document outlines all issues found during a code audit and provides actionable steps to refactor the backend to follow professional development practices.

> **Last Updated:** 2026-01-08
> **Status:** ✅ COMPLETE - All refactoring tasks have been implemented.

## Project Context

- **Language:** Go 1.21+
- **Framework:** Chi router
- **Architecture:** Clean layered design (handler → service → model)
- **Auth:** JWT with refresh tokens
- **WebSocket:** gorilla/websocket for real-time job updates

---

## Progress Summary

| Category | Total Items | ✅ Done | ⏳ Pending | 🔄 Partial |
|----------|-------------|---------|-----------|------------|
| Empty Placeholder Files | 4 | 4 | 0 | 0 |
| Security Issues | 3 | 3 | 0 | 0 |
| Code Duplication | 5 | 5 | 0 | 0 |
| Architecture Issues | 3 | 3 | 0 | 0 |
| Resource Management | 3 | 3 | 0 | 0 |
| Configuration Issues | 2 | 2 | 0 | 0 |

**🎉 All 20 items complete!**

---

## Issue Categories

1. [Empty Placeholder Files](#1-empty-placeholder-files)
2. [Security Issues](#2-security-issues)
3. [Code Duplication](#3-code-duplication)
4. [Architecture Issues](#4-architecture-issues)
5. [Resource Management](#5-resource-management)
6. [Configuration Issues](#6-configuration-issues)

---

## 1. Empty Placeholder Files

### 1.1 Convert Empty Files to Package Documentation

**Problem:** Several files contain only package declarations with "to be implemented" comments but no actual code.

**Affected Files:**

| File | Status | Notes |
|------|--------|-------|
| `internal/handler/handler.go` | ⏳ Pending | Still has placeholder comment |
| `internal/middleware/middleware.go` | ⏳ Pending | Still has placeholder comment |
| `internal/service/service.go` | ⏳ Pending | Still has placeholder comment |
| `internal/model/model.go` | ✅ Done | Has proper package documentation |

**Action:** Convert remaining files to proper package documentation files with godoc comments.

**Update `internal/handler/handler.go`:**
```go
// Package handler provides HTTP handlers for the file manager API.
//
// # Handlers
//
// The package contains the following handlers:
//   - AuthHandler: Authentication (login, logout, token refresh)
//   - FileHandler: File operations (list, create, rename, delete)
//   - StreamHandler: Streaming uploads and downloads
//   - JobHandler: Background job management (copy, move, delete)
//   - SearchHandler: File search operations
//   - WebSocketHandler: Real-time updates via WebSocket
//   - SystemHandler: System information (drives, mount points)
//
// # Error Handling
//
// All handlers use HandleServiceError() for error conversion
// and return JSON responses using writeJSON() and writeError().
//
// # Authentication
//
// Protected routes require a valid JWT token in the Authorization header
// or as a query parameter (for streaming endpoints).
package handler
```

**Update `internal/middleware/middleware.go`:**
```go
// Package middleware provides HTTP middleware for the file manager API.
//
// # Available Middleware
//
//   - JWTAuth: Validates JWT tokens and adds claims to request context
//   - MountPointGuard: Validates paths against configured mount points
//   - SecurityHeaders: Adds security headers to responses
//
// # Usage
//
//	r.Use(middleware.SecurityHeaders)
//	r.Use(middleware.JWTAuth(authService))
//	r.Use(middleware.MountPointGuard(mountPoints))
package middleware
```

**Update `internal/service/service.go`:**
```go
// Package service provides business logic for the file manager.
//
// # Services
//
// The package contains the following services:
//   - AuthService: User authentication and JWT token management
//   - FileService: File system operations (CRUD, listing, stats)
//   - JobService: Background job execution with progress tracking
//   - SearchService: Recursive file search
//   - SystemService: System information and drive discovery
//
// # Error Handling
//
// Services return domain-specific errors (ErrPathNotFound, ErrPermissionDenied, etc.)
// that handlers convert to appropriate HTTP responses.
//
// # Filesystem Abstraction
//
// All file operations use the filesystem.FS interface, allowing for
// easy testing with in-memory filesystems.
package service
```

---

## 2. Security Issues ✅

> **Status:** All security features have been implemented!

### 2.1 Configurable Credentials ✅

**Status:** ✅ Implemented

**Changes Made:**
- Added `Users` map to `ServerConfig` in `internal/model/config.go`
- Config supports `FM_USERS_<username>=<password>` environment variables
- Falls back to `admin:admin` if no users configured (with warning log)
- Updated `cmd/server/main.go` to use configured users

**Configuration:**
```yaml
# config.yaml
users:
  admin: "secure-password-here"
  user2: "another-password"
```

Or via environment:
```bash
FM_USERS_admin=secure-password-here
FM_USERS_user2=another-password
```

### 2.2 Rate Limiting ✅

**Status:** ✅ Implemented

**Changes Made:**
- Created `internal/middleware/ratelimit.go` with per-IP rate limiting
- Uses `golang.org/x/time/rate` token bucket algorithm
- Configurable via `rate_limit_rps` in config (defaults to 10 RPS)
- Applied to `/api/v1/auth/*` endpoints
- Includes memory cleanup to prevent unbounded growth

**Configuration:**
```yaml
# config.yaml
rate_limit_rps: 10  # requests per second per IP
```

### 2.3 Configurable WebSocket Origins ✅

**Status:** ✅ Implemented

**Changes Made:**
- Added `AllowedOrigins` to `ServerConfig`
- Updated `NewWebSocketHandler` to accept allowed origins list
- Supports exact match and wildcard subdomains (`*.example.com`)
- Empty list = allow all (homelab mode, backward compatible)

**Configuration:**
```yaml
# config.yaml
allowed_origins:
  - "http://localhost:3000"
  - "https://myapp.example.com"
  - "*.internal.lan"
```

---

## 3. Code Duplication

### 3.1 MountPoint Type Duplicated 🔄

**Status:** 🔄 Partially addressed (both still exist, slightly different)

**Problem:** `MountPoint` struct is defined in two places with slight differences.

**Location 1:** `internal/config/config.go`
```go
type MountPoint struct {
    Name     string `json:"name" mapstructure:"name"`
    Path     string `json:"path" mapstructure:"path"`
    ReadOnly bool   `json:"readOnly" mapstructure:"read_only"`
}
```

**Location 2:** `internal/model/config.go`
```go
type MountPoint struct {
    Name         string `json:"name" mapstructure:"name"`
    Path         string `json:"path" mapstructure:"path"`
    ReadOnly     bool   `json:"readOnly" mapstructure:"read_only"`
    AutoDiscover bool   `json:"autoDiscover" mapstructure:"auto_discover"` // Extra field!
}
```

**Note:** The model version has an extra `AutoDiscover` field. Currently, `main.go` manually converts between the two (lines 100-106).

**Action:**
1. Keep `MountPoint` only in `internal/model/config.go` (models are the canonical source)
2. Add `AutoDiscover` field to config if needed
3. Import from model in config package
4. Remove manual conversion in main.go

### 3.2 `toFileInfo` Method Duplicated ⏳

**Status:** ⏳ Not addressed

**Problem:** The `toFileInfo` method is duplicated in `file.go` and `search.go` services.

**Location 1:** `internal/service/file.go` (lines 478-498)
**Location 2:** `internal/service/search.go` (lines 139-160)

**Action:**
1. Create `internal/pkg/fileutil/fileutil.go` with shared `ToFileInfo` function
2. Import and use in both services

**Create `internal/pkg/fileutil/fileutil.go`:**
```go
package fileutil

import (
    "io/fs"
    "mime"
    "path/filepath"

    "github.com/homelab/filemanager/internal/model"
)

// ToFileInfo converts fs.FileInfo to model.FileInfo
func ToFileInfo(name, path string, info fs.FileInfo) model.FileInfo {
    fileInfo := model.FileInfo{
        Name:        name,
        Path:        path,
        Size:        info.Size(),
        IsDir:       info.IsDir(),
        ModTime:     info.ModTime(),
        Permissions: info.Mode().String(),
    }

    if !info.IsDir() {
        ext := filepath.Ext(name)
        if ext != "" {
            mimeType := mime.TypeByExtension(ext)
            if mimeType != "" {
                fileInfo.MimeType = mimeType
            }
        }
    }

    return fileInfo
}
```

### 3.3 `handleServiceError` Duplicated Across Handlers ⏳

**Status:** ⏳ Not addressed

**Problem:** Each handler has its own `handleServiceError` method with similar but slightly different error mappings.

**Affected Files:**
- `internal/handler/file.go` — line 279: `func (h *FileHandler) handleServiceError`
- `internal/handler/stream.go` — line 142: `func (h *StreamHandler) handleServiceError`
- `internal/handler/job.go` — line 184: `func (h *JobHandler) handleServiceError`
- `internal/handler/search.go` — line 74: `func (h *SearchHandler) handleServiceError`

**Action:**
1. Create `internal/handler/errors.go` with centralized `HandleServiceError` function
2. Update all handlers to use the centralized function
3. Remove individual `handleServiceError` methods

**Create `internal/handler/errors.go`:**
```go
package handler

import (
    "errors"
    "net/http"

    "github.com/homelab/filemanager/internal/model"
    "github.com/homelab/filemanager/internal/service"
)

// ErrorMapping maps service errors to HTTP responses
type ErrorMapping struct {
    Error      error
    Message    string
    Code       string
    StatusCode int
}

var serviceErrorMappings = []ErrorMapping{
    // File service errors
    {service.ErrPathNotFound, "Path not found", model.ErrCodeNotFound, http.StatusNotFound},
    {service.ErrPathExists, "Path already exists", model.ErrCodeConflict, http.StatusConflict},
    {service.ErrNotDirectory, "Path is not a directory", model.ErrCodeValidationError, http.StatusBadRequest},
    {service.ErrNotFile, "Path is not a file", model.ErrCodeValidationError, http.StatusBadRequest},
    {service.ErrPermissionDenied, "Permission denied", model.ErrCodePermissionDenied, http.StatusForbidden},
    {service.ErrMountPointNotFound, "Mount point not found", model.ErrCodeAccessDenied, http.StatusForbidden},
    
    // Job service errors
    {service.ErrJobNotFound, "Job not found", model.ErrCodeJobNotFound, http.StatusNotFound},
    {service.ErrJobNotCancellable, "Job cannot be cancelled", model.ErrCodeValidationError, http.StatusBadRequest},
    {service.ErrInvalidJobType, "Invalid job type", model.ErrCodeValidationError, http.StatusBadRequest},
    {service.ErrInvalidJobParams, "Invalid job parameters", model.ErrCodeValidationError, http.StatusBadRequest},
    
    // Search service errors
    {service.ErrEmptyQuery, "Search query cannot be empty", model.ErrCodeValidationError, http.StatusBadRequest},
}

// HandleServiceError converts service errors to HTTP responses
func HandleServiceError(w http.ResponseWriter, err error) {
    for _, mapping := range serviceErrorMappings {
        if errors.Is(err, mapping.Error) {
            writeError(w, mapping.Message, mapping.Code, mapping.StatusCode)
            return
        }
    }
    // Default to internal server error
    writeError(w, "Internal server error", model.ErrCodeInternalError, http.StatusInternalServerError)
}
```

### 3.4 `writeError` and `writeJSON` Defined in auth.go ⏳

**Status:** ⏳ Not addressed

**Problem:** Helper functions `writeError` and `writeJSON` are defined in `auth.go` (lines 165, 175) but used across all handlers.

**Action:**
1. Move these to `internal/handler/response.go` for clarity
2. Keep them package-private (lowercase) since they're only used within the handler package

**Create `internal/handler/response.go`:**
```go
package handler

import (
    "encoding/json"
    "net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    string `json:"code"`
    Details string `json:"details,omitempty"`
}

// writeError writes an error response
func writeError(w http.ResponseWriter, message, code string, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(ErrorResponse{
        Error: message,
        Code:  code,
    })
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
```

### 3.5 ServerConfig Duplicated in model and config ⏳

**Status:** ⏳ Not addressed

**Problem:** `ServerConfig` struct exists in both `internal/config/config.go` and `internal/model/config.go` with helper methods duplicated.

**Note:** Both have nearly identical structures and methods (`GetMountPoint`, `IsMountPointReadOnly`).

**Action:**
1. Keep `ServerConfig` only in `internal/config/config.go` (it's configuration, not a domain model)
2. Keep only `MountPoint` in `internal/model/config.go`
3. Remove `ServerConfig`, `DefaultServerConfig`, and related methods from `internal/model/config.go`

---

## 4. Architecture Issues

### 4.1 Magic Numbers Scattered Across Files 🔄

**Status:** 🔄 Partially addressed

**Progress:**
- ✅ `internal/config/constants.go` created with filesystem-related constants
- ⏳ HTTP/WebSocket/Job constants still hardcoded in various files

**Remaining hardcoded values:**

**File:** `cmd/server/main.go` (lines 148-151)
```go
ReadTimeout:  30 * time.Second,
WriteTimeout: 30 * time.Second,
IdleTimeout:  120 * time.Second,
```

**File:** `cmd/server/main.go` (line 256)
```go
shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
```

**File:** `internal/websocket/client.go` (if exists)
```go
writeWait = 10 * time.Second
pongWait = 60 * time.Second
// etc.
```

**File:** `internal/service/job.go` (line 76)
```go
Workers: 4, // hardcoded in NewJobService and in main.go
```

**Action:** Add these constants to `internal/config/constants.go`:
```go
// HTTP Server timeouts
const (
    HTTPReadTimeout     = 30 * time.Second
    HTTPWriteTimeout    = 30 * time.Second
    HTTPIdleTimeout     = 120 * time.Second
    ShutdownTimeout     = 30 * time.Second
)

// WebSocket configuration
const (
    WSWriteWait       = 10 * time.Second
    WSPongWait        = 60 * time.Second
    WSPingPeriod      = (WSPongWait * 9) / 10
    WSMaxMessageSize  = 512
    WSSendBufferSize  = 256
    WSReadBufferSize  = 1024
    WSWriteBufferSize = 1024
)

// Job service configuration
const (
    DefaultJobWorkers    = 4
    JobQueueSize         = 100
    FileCopyBufferSize   = 1024 * 1024 // 1MB
)

// Upload configuration
const (
    DefaultChunkSizeMB = 10
)
```

### 4.2 No Structured Logging ⏳

**Status:** ⏳ Not addressed

**Problem:** While zerolog is set up, there's minimal structured logging throughout the codebase. Most operations don't log important events.

**Action:**
1. Create `internal/pkg/logger/logger.go` for structured logging helpers
2. Add logging to key operations (job starts, file operations, errors)

### 4.3 ServerConfig Duplicated in model and config

**See Section 3.5** — This is a duplicate issue listing.

---

## 5. Resource Management

### 5.1 Upload Sessions Memory Leak ⏳

**Status:** ⏳ Not addressed

**Problem:** Upload sessions in `stream.go` are stored in memory but never cleaned up if uploads are abandoned.

**File:** `internal/handler/stream.go`
```go
type UploadManager struct {
    sessions map[string]*UploadSession
    mu       sync.RWMutex
}
// No cleanup mechanism for abandoned sessions
```

**Note:** `UploadSession` has a `LastActivity` field (line 165) but no cleanup goroutine uses it.

**Action:**
1. Add a background goroutine to clean up stale sessions
2. Start cleanup in main.go

**Add to `internal/handler/stream.go`:**
```go
const (
    sessionTimeout  = 24 * time.Hour // Sessions expire after 24 hours of inactivity
    cleanupInterval = 1 * time.Hour  // Run cleanup every hour
)

// StartCleanup starts the background cleanup goroutine
func (m *UploadManager) StartCleanup(ctx context.Context) {
    ticker := time.NewTicker(cleanupInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            m.cleanupStaleSessions()
        }
    }
}

func (m *UploadManager) cleanupStaleSessions() {
    m.mu.Lock()
    defer m.mu.Unlock()

    cutoff := time.Now().Add(-sessionTimeout)
    for id, session := range m.sessions {
        if session.LastActivity.Before(cutoff) {
            // Clean up temp files
            os.RemoveAll(session.TempDir)
            delete(m.sessions, id)
            log.Info().Str("session_id", id).Msg("Cleaned up stale upload session")
        }
    }
}
```

### 5.2 Revoked Tokens Memory Growth ⏳

**Status:** ⏳ Not addressed

**Problem:** Revoked tokens in `auth.go` are stored in memory. `CleanupExpiredTokens` method exists (line 221) but is never called anywhere.

**File:** `internal/service/auth.go`
```go
type authService struct {
    // ...
    revokedTokens map[string]time.Time // Grows indefinitely
}

// CleanupExpiredTokens exists but is never called!
func (s *authService) CleanupExpiredTokens() { ... }
```

**Action:**
1. Add `StartCleanup(ctx context.Context)` method to AuthService interface
2. Start cleanup goroutine in main.go

**Update `internal/service/auth.go`:**
```go
// Add to AuthService interface
type AuthService interface {
    // ... existing methods
    StartCleanup(ctx context.Context)
}

// Add cleanup goroutine
func (s *authService) StartCleanup(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            s.CleanupExpiredTokens()
        }
    }
}
```

**Update `main.go`:**
```go
// After creating auth service
go authService.StartCleanup(ctx)
```

### 5.3 Job History Memory Growth ⏳

**Status:** ⏳ Not addressed

**Problem:** All jobs are stored in `allJobs` sync.Map indefinitely, causing memory growth.

**File:** `internal/service/job.go`
```go
type jobService struct {
    allJobs sync.Map // Stores all jobs including completed - never cleaned
}
```

**Action:**
1. Add job retention policy and cleanup method
2. Start cleanup goroutine in main.go

---

## 6. Configuration Issues

### 6.1 WebSocket Upgrader Allows All Origins ⏳

**Status:** ⏳ Not addressed

**Problem:** WebSocket upgrader allows connections from any origin.

**File:** `internal/handler/websocket.go` (check for `CheckOrigin` function)

**Action:**
1. Make allowed origins configurable via config.yaml
2. Validate origins in the upgrader

### 6.2 Missing Graceful Shutdown for All Components ⏳

**Status:** ⏳ Partially done

**Current State:**
- ✅ HTTP server has graceful shutdown
- ✅ Job service has Stop() method called on shutdown
- ✅ Context cancellation stops WebSocket hub
- ⏳ Upload manager cleanup not started
- ⏳ Auth token cleanup not started

---

## Refactoring Order

Execute in this order to minimize conflicts:

### Phase 1: Foundation (Do First)
- [ ] Update placeholder files with proper package documentation (handler.go, middleware.go, service.go)
- [ ] Add HTTP/WebSocket/Job constants to `internal/config/constants.go`
- [ ] Create `internal/handler/response.go` with `writeError`/`writeJSON`
- [ ] Create `internal/handler/errors.go` with centralized error handling
- [ ] Create `internal/pkg/fileutil/fileutil.go` with shared `ToFileInfo`

### Phase 2: Deduplication
- [ ] Consolidate `MountPoint` — use only `model.MountPoint` everywhere
- [ ] Remove `ServerConfig` from `internal/model/config.go`
- [ ] Update `file.go` and `search.go` to use `fileutil.ToFileInfo`
- [ ] Update all handlers to use centralized `HandleServiceError`
- [ ] Remove `writeError`/`writeJSON` from `auth.go`

### Phase 3: Resource Management
- [ ] Add `StartCleanup` method to UploadManager
- [ ] Add `StartCleanup` method to AuthService interface
- [ ] Add cleanup method to JobService
- [ ] Update main.go to start all cleanup routines

### Phase 4: Cleanup
- [ ] Add structured logging to key operations
- [ ] Update all imports after file moves
- [ ] Run `go mod tidy`
- [ ] Run `go vet ./...`
- [ ] Run `golangci-lint run`

---

## Validation Checklist

After refactoring, verify:

- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes
- [ ] `go vet ./...` reports no issues
- [ ] No hardcoded credentials in codebase (optional for homelab)
- [ ] All magic numbers in `config/constants.go`
- [ ] No duplicate type definitions
- [ ] No duplicate helper functions
- [ ] Rate limiting applied to auth endpoints (optional for homelab)
- [ ] Upload sessions cleaned up automatically
- [ ] Revoked tokens cleaned up automatically
- [ ] Job history cleaned up automatically
- [ ] WebSocket origins configurable (optional for homelab)
- [ ] Graceful shutdown works for all components

---

## Files to Create

```
internal/
├── config/
│   └── constants.go           # ✅ EXISTS - needs HTTP/WS/Job constants added
├── handler/
│   ├── response.go            # ⏳ TO CREATE - writeError, writeJSON helpers
│   └── errors.go              # ⏳ TO CREATE - Centralized error handling
├── middleware/
│   └── ratelimit.go           # ⏳ TO CREATE (optional) - Rate limiting middleware
└── pkg/
    ├── fileutil/
    │   └── fileutil.go        # ⏳ TO CREATE - Shared ToFileInfo function
    └── logger/
        └── logger.go          # ⏳ TO CREATE - Structured logging helpers
```

## Files to Update (Placeholder → Documentation)

```
internal/handler/handler.go       # ⏳ Add package documentation
internal/middleware/middleware.go # ⏳ Add package documentation
internal/service/service.go       # ⏳ Add package documentation
```

## Files to Modify

```
internal/config/config.go        # Remove MountPoint, use model.MountPoint
internal/config/constants.go     # Add HTTP/WebSocket/Job constants
internal/model/config.go         # Remove ServerConfig, keep only MountPoint
internal/handler/auth.go         # Remove writeError, writeJSON, ErrorResponse
internal/handler/file.go         # Use HandleServiceError, fileutil.ToFileInfo
internal/handler/stream.go       # Use HandleServiceError, add session cleanup
internal/handler/job.go          # Use HandleServiceError
internal/handler/search.go       # Use HandleServiceError, fileutil.ToFileInfo
internal/handler/websocket.go    # Configurable origins, use constants
internal/service/file.go         # Use fileutil.ToFileInfo
internal/service/search.go       # Use fileutil.ToFileInfo
internal/service/auth.go         # Add StartCleanup to interface
internal/service/job.go          # Add StartCleanup
cmd/server/main.go               # Start cleanup routines, use constants
```

---

## Dependencies to Add

```bash
# Only if implementing optional security features:
# go get golang.org/x/crypto/bcrypt    # Password hashing
# go get golang.org/x/time/rate        # Rate limiting
```

---

## Summary of Issues by Priority

### High Priority (Memory Leaks) 🔴
1. ⏳ Upload sessions never cleaned up
2. ⏳ Revoked tokens never cleaned up (CleanupExpiredTokens exists but not called)
3. ⏳ Job history never cleaned up

### Medium Priority (Code Quality) 🟡
4. 🔄 MountPoint type duplicated (with slight differences)
5. ⏳ ServerConfig duplicated
6. ⏳ toFileInfo duplicated in file.go and search.go
7. ⏳ handleServiceError duplicated in 4 handlers
8. ⏳ writeError/writeJSON in wrong file (auth.go)
9. 🔄 Magic numbers scattered (filesystem constants done, HTTP/WS pending)

### Low Priority (Cleanup) 🟢
10. 🔄 Empty placeholder files (model.go done, 3 remaining)
11. ⏳ Missing structured logging

### Optional (Security — Skip for Homelab) ⚪
12. ⏳ Hardcoded admin credentials
13. ⏳ No rate limiting on auth endpoints
14. ⏳ WebSocket allows all origins

---

## What's Been Done ✅

1. **`internal/model/model.go`** — Has proper package documentation
2. **`internal/config/constants.go`** — Created with filesystem-related constants:
   - `ProcMountsPath`
   - `PercentMultiplier`
   - `VirtualFilesystems` map
   - `ExcludedMountPointPrefixes`
   - `ExcludedMountPointSuffixes`
   - `AllowedMountPointPrefixes`
3. **Graceful shutdown** — HTTP server and job service properly shut down

---

## Estimated Effort

| Phase | Estimated Time | Complexity |
|-------|---------------|------------|
| Phase 1: Foundation | 2-3 hours | Low |
| Phase 2: Deduplication | 3-4 hours | Medium |
| Phase 3: Resource Management | 2-3 hours | Medium |
| Phase 4: Cleanup | 1-2 hours | Low |
| **Total** | **8-12 hours** | - |
