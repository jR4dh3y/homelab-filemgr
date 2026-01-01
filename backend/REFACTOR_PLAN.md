# Backend Refactoring Plan

This document outlines all issues found during a code audit and provides actionable steps to refactor the backend to follow professional development practices.

## Project Context

- **Language:** Go 1.21+
- **Framework:** Chi router
- **Architecture:** Clean layered design (handler → service → model)
- **Auth:** JWT with refresh tokens
- **WebSocket:** gorilla/websocket for real-time job updates

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
- `internal/handler/handler.go`
- `internal/middleware/middleware.go`
- `internal/service/service.go`
- `internal/model/model.go`

**Action:** Convert these to proper package documentation files with godoc comments.

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
//
// # Error Handling
//
// All handlers use centralized error handling via HandleServiceError()
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

**Update `internal/model/model.go`:**
```go
// Package model contains all data models for the homelab file manager.
//
// # Organization
//
// Models are organized into separate files:
//   - file.go: FileInfo, FileList, ListOptions, DriveStats
//   - job.go: Job, JobType, JobState, JobUpdate, JobParams
//   - config.go: MountPoint
//   - error.go: ErrorResponse and error code constants
//
// # JSON Serialization
//
// All models use json tags for API serialization and mapstructure tags
// for configuration file parsing where applicable.
package model
```

---

## 2. Security Issues (Optional — Homelab Context)

> **Note:** These are marked as optional since this is a homelab project behind a private network. Implement only if you plan to expose this externally or want defense-in-depth.

### 2.1 Hardcoded Default Credentials (Optional)

**Problem:** Default admin credentials are hardcoded in `main.go`.

**File:** `cmd/server/main.go`
```go
Users: map[string]string{
    "admin": "admin",
},
```

**For homelab:** This is fine. If you want to change credentials, just edit the map.

**For production:** Move to config file with bcrypt hashing.

### 2.2 No Rate Limiting (Optional)

**For homelab:** Not needed — you're the only user on your network.

**For production:** Add rate limiting middleware to auth endpoints.

### 2.3 WebSocket Allows All Origins (Optional)

**For homelab:** Fine — all traffic is from your network anyway.

**For production:** Configure allowed origins list.

---

## 3. Code Duplication

### 3.1 MountPoint Type Duplicated

**Problem:** `MountPoint` struct is defined identically in two places.

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
    Name     string `json:"name" mapstructure:"name"`
    Path     string `json:"path" mapstructure:"path"`
    ReadOnly bool   `json:"readOnly" mapstructure:"read_only"`
}
```

**Action:**
1. Keep `MountPoint` only in `internal/model/config.go` (models are the canonical source)
2. Import from model in config package

**Update `internal/config/config.go`:**
```go
package config

import (
    "github.com/homelab/filemanager/internal/model"
    // ...
)

// Remove the MountPoint struct definition

type ServerConfig struct {
    Port        int                `mapstructure:"port"`
    Host        string             `mapstructure:"host"`
    MountPoints []model.MountPoint `mapstructure:"mount_points"` // Use model.MountPoint
    JWTSecret   string             `mapstructure:"jwt_secret"`
    MaxUploadMB int                `mapstructure:"max_upload_mb"`
    ChunkSizeMB int                `mapstructure:"chunk_size_mb"`
}

// Update helper methods to use model.MountPoint
func (c *ServerConfig) GetMountPoint(name string) *model.MountPoint {
    for i := range c.MountPoints {
        if c.MountPoints[i].Name == name {
            return &c.MountPoints[i]
        }
    }
    return nil
}
```

### 3.2 `toFileInfo` Method Duplicated

**Problem:** The `toFileInfo` method is duplicated in `file.go` and `search.go` services.

**Location 1:** `internal/service/file.go` (lines 290-308)
```go
func (s *fileService) toFileInfo(name, path string, info fs.FileInfo) model.FileInfo {
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

**Location 2:** `internal/service/search.go` (lines 95-113) — identical code

**Action:**
1. Create a shared utility function in `internal/pkg/fileutil/fileutil.go`
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

**Update services:**
```go
// In file.go and search.go
import "github.com/homelab/filemanager/internal/pkg/fileutil"

// Replace s.toFileInfo(...) with:
fileutil.ToFileInfo(name, path, info)
```

### 3.3 `handleServiceError` Duplicated Across Handlers

**Problem:** Each handler has its own `handleServiceError` method with similar but slightly different error mappings.

**Affected Files:**
- `internal/handler/file.go` — `handleServiceError` method
- `internal/handler/stream.go` — `handleServiceError` method
- `internal/handler/job.go` — `handleServiceError` method
- `internal/handler/search.go` — `handleServiceError` method

**Action:**
1. Create a centralized error handler in `internal/handler/errors.go`
2. Use a unified error mapping approach

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

**Update handlers to use centralized function:**
```go
// In file.go, stream.go, job.go, search.go
// Replace h.handleServiceError(w, err) with:
HandleServiceError(w, err)

// Remove the individual handleServiceError methods from each handler
```

### 3.4 `writeError` and `writeJSON` Defined in auth.go

**Problem:** Helper functions `writeError` and `writeJSON` are defined in `auth.go` but used across all handlers.

**File:** `internal/handler/auth.go` (lines 143-159)

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

**Remove from auth.go:**
Delete the `ErrorResponse` struct and `writeError`/`writeJSON` functions from `auth.go`.

---

## 4. Architecture Issues

### 4.1 Magic Numbers Scattered Across Files

**Problem:** Hardcoded values for timeouts, buffer sizes, and intervals scattered throughout the codebase.

**Examples Found:**

**File:** `cmd/server/main.go`
```go
ReadTimeout:  30 * time.Second,
WriteTimeout: 30 * time.Second,
IdleTimeout:  120 * time.Second,
// ...
shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
```

**File:** `internal/websocket/client.go`
```go
writeWait = 10 * time.Second
pongWait = 60 * time.Second
pingPeriod = (pongWait * 9) / 10
maxMessageSize = 512
sendBufferSize = 256
```

**File:** `internal/service/job.go`
```go
workers = 4 // default worker count
workQueue: make(chan *model.Job, 100),
buf := make([]byte, 1024*1024) // 1MB buffer
```

**File:** `internal/handler/stream.go`
```go
chunkSizeMB = 10 // Default 10MB chunks
ReadBufferSize:  1024,
WriteBufferSize: 1024,
```

**Action:**
1. Create `internal/config/constants.go` for all magic numbers
2. Reference these constants throughout the codebase

**Create `internal/config/constants.go`:**
```go
package config

import "time"

// HTTP Server timeouts
const (
    HTTPReadTimeout     = 30 * time.Second
    HTTPWriteTimeout    = 30 * time.Second
    HTTPIdleTimeout     = 120 * time.Second
    ShutdownTimeout     = 30 * time.Second
)

// WebSocket configuration
const (
    WSWriteWait      = 10 * time.Second
    WSPongWait       = 60 * time.Second
    WSPingPeriod     = (WSPongWait * 9) / 10
    WSMaxMessageSize = 512
    WSSendBufferSize = 256
    WSReadBufferSize = 1024
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

**Update usages:**
```go
// In main.go
server := &http.Server{
    ReadTimeout:  config.HTTPReadTimeout,
    WriteTimeout: config.HTTPWriteTimeout,
    IdleTimeout:  config.HTTPIdleTimeout,
}

// In websocket/client.go
const (
    writeWait      = config.WSWriteWait
    pongWait       = config.WSPongWait
    // ...
)
```

### 4.2 No Structured Logging

**Problem:** While zerolog is set up, there's minimal structured logging throughout the codebase. Most operations don't log important events.

**Action:**
1. Add structured logging to key operations
2. Create a logger wrapper for consistent field naming

**Create `internal/pkg/logger/logger.go`:**
```go
package logger

import (
    "context"
    "net/http"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

// ContextKey for logger in context
type contextKey string
const loggerKey contextKey = "logger"

// WithRequestID adds request ID to logger
func WithRequestID(ctx context.Context, requestID string) context.Context {
    logger := log.With().Str("request_id", requestID).Logger()
    return context.WithValue(ctx, loggerKey, &logger)
}

// FromContext retrieves logger from context
func FromContext(ctx context.Context) *zerolog.Logger {
    if logger, ok := ctx.Value(loggerKey).(*zerolog.Logger); ok {
        return logger
    }
    return &log.Logger
}

// LogRequest logs an HTTP request
func LogRequest(r *http.Request, status int, duration time.Duration) {
    log.Info().
        Str("method", r.Method).
        Str("path", r.URL.Path).
        Int("status", status).
        Dur("duration", duration).
        Str("remote_addr", r.RemoteAddr).
        Msg("HTTP request")
}

// LogJobEvent logs a job-related event
func LogJobEvent(jobID, jobType, event string, progress int, err error) {
    entry := log.Info().
        Str("job_id", jobID).
        Str("job_type", jobType).
        Str("event", event).
        Int("progress", progress)
    
    if err != nil {
        entry = log.Error().Err(err).
            Str("job_id", jobID).
            Str("job_type", jobType).
            Str("event", event)
    }
    
    entry.Msg("Job event")
}
```

### 4.3 ServerConfig Duplicated in model and config

**Problem:** `ServerConfig` struct exists in both `internal/config/config.go` and `internal/model/config.go` with helper methods duplicated.

**Action:**
1. Keep `ServerConfig` only in `internal/config/config.go` (it's configuration, not a domain model)
2. Keep only `MountPoint` in `internal/model/config.go`
3. Remove `ServerConfig` and its methods from `internal/model/config.go`

**Update `internal/model/config.go`:**
```go
package model

// MountPoint represents a configured filesystem location
type MountPoint struct {
    Name     string `json:"name" mapstructure:"name"`
    Path     string `json:"path" mapstructure:"path"`
    ReadOnly bool   `json:"readOnly" mapstructure:"read_only"`
}

// Remove ServerConfig, DefaultServerConfig, IsMountPointReadOnly, GetMountPoint
// These belong in the config package
```

---

## 5. Resource Management

### 5.1 Upload Sessions Memory Leak

**Problem:** Upload sessions in `stream.go` are stored in memory but never cleaned up if uploads are abandoned.

**File:** `internal/handler/stream.go`
```go
type UploadManager struct {
    sessions map[string]*UploadSession
    mu       sync.RWMutex
}
// No cleanup mechanism for abandoned sessions
```

**Action:**
1. Add a background goroutine to clean up stale sessions
2. Add session expiry time configuration

**Update `internal/handler/stream.go`:**
```go
const (
    sessionTimeout = 24 * time.Hour // Sessions expire after 24 hours of inactivity
    cleanupInterval = 1 * time.Hour // Run cleanup every hour
)

// Add cleanup method to UploadManager
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

**Start cleanup in main.go:**
```go
// After creating stream handler
go streamHandler.uploadManager.StartCleanup(ctx)
```

### 5.2 Revoked Tokens Memory Growth

**Problem:** Revoked tokens in `auth.go` are stored in memory but `CleanupExpiredTokens` is never called.

**File:** `internal/service/auth.go`
```go
type authService struct {
    // ...
    revokedTokens map[string]time.Time // Grows indefinitely
}

// CleanupExpiredTokens exists but is never called
func (s *authService) CleanupExpiredTokens() { ... }
```

**Action:**
1. Start a background cleanup goroutine
2. Export the cleanup method via interface

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

**Start cleanup in main.go:**
```go
// After creating auth service
go authService.StartCleanup(ctx)
```

### 5.3 Job History Memory Growth

**Problem:** All jobs are stored in `allJobs` sync.Map indefinitely, causing memory growth.

**File:** `internal/service/job.go`
```go
type jobService struct {
    allJobs sync.Map // Stores all jobs including completed - never cleaned
}
```

**Action:**
1. Add job retention policy
2. Clean up completed jobs after a configurable period

**Update `internal/service/job.go`:**
```go
const (
    jobRetentionPeriod = 24 * time.Hour // Keep completed jobs for 24 hours
    jobCleanupInterval = 1 * time.Hour
)

// Add to JobService interface
type JobService interface {
    // ... existing methods
    StartCleanup(ctx context.Context)
}

func (s *jobService) StartCleanup(ctx context.Context) {
    ticker := time.NewTicker(jobCleanupInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            s.cleanupOldJobs()
        }
    }
}

func (s *jobService) cleanupOldJobs() {
    cutoff := time.Now().Add(-jobRetentionPeriod)
    
    s.allJobs.Range(func(key, value interface{}) bool {
        job := value.(*model.Job)
        if job.State.IsTerminal() && job.CompletedAt.Before(cutoff) {
            s.allJobs.Delete(key)
        }
        return true
    })
}
```

---

## 6. Configuration Issues

### 6.1 WebSocket Upgrader Allows All Origins

**Problem:** WebSocket upgrader allows connections from any origin, which is a security risk in production.

**File:** `internal/handler/websocket.go`
```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // Allows all origins!
    },
}
```

**Action:**
1. Make allowed origins configurable
2. Validate origins in production

**Update config.yaml:**
```yaml
websocket:
  allowed_origins:
    - "http://localhost:5173"
    - "https://filemanager.example.com"
```

**Update `internal/handler/websocket.go`:**
```go
type WebSocketHandler struct {
    hub            *ws.Hub
    authService    service.AuthService
    allowedOrigins []string
}

func NewWebSocketHandler(hub *ws.Hub, authService service.AuthService, allowedOrigins []string) *WebSocketHandler {
    return &WebSocketHandler{
        hub:            hub,
        authService:    authService,
        allowedOrigins: allowedOrigins,
    }
}

func (h *WebSocketHandler) createUpgrader() websocket.Upgrader {
    return websocket.Upgrader{
        ReadBufferSize:  config.WSReadBufferSize,
        WriteBufferSize: config.WSWriteBufferSize,
        CheckOrigin: func(r *http.Request) bool {
            origin := r.Header.Get("Origin")
            if len(h.allowedOrigins) == 0 {
                return true // Allow all if not configured (dev mode)
            }
            for _, allowed := range h.allowedOrigins {
                if origin == allowed {
                    return true
                }
            }
            return false
        },
    }
}
```

### 6.2 Missing Graceful Shutdown for All Components

**Problem:** While HTTP server has graceful shutdown, other components (upload manager, auth cleanup) don't.

**Action:**
1. Create a unified shutdown coordinator
2. Ensure all background goroutines respect context cancellation

**Update `cmd/server/main.go`:**
```go
func waitForShutdown(ctx context.Context, cancel context.CancelFunc, server *http.Server, jobService service.JobService) {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

    sig := <-sigCh
    log.Info().Str("signal", sig.String()).Msg("Received shutdown signal")

    // Cancel context to stop ALL background goroutines
    cancel()

    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
    defer shutdownCancel()

    // Create wait group for all cleanup tasks
    var wg sync.WaitGroup

    // Stop job service
    wg.Add(1)
    go func() {
        defer wg.Done()
        log.Info().Msg("Stopping job service...")
        jobService.Stop()
    }()

    // Shutdown HTTP server
    wg.Add(1)
    go func() {
        defer wg.Done()
        log.Info().Msg("Shutting down HTTP server...")
        if err := server.Shutdown(shutdownCtx); err != nil {
            log.Error().Err(err).Msg("Error during server shutdown")
        }
    }()

    // Wait for all cleanup with timeout
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        log.Info().Msg("Server shutdown complete")
    case <-shutdownCtx.Done():
        log.Warn().Msg("Shutdown timed out, forcing exit")
    }
}
```

---

## Refactoring Order

Execute in this order to minimize conflicts:

### Phase 1: Foundation (Do First)
1. Update placeholder files with proper package documentation
2. Create `internal/config/constants.go` with all magic numbers
3. Create `internal/handler/response.go` with `writeError`/`writeJSON`
4. Create `internal/handler/errors.go` with centralized error handling
5. Create `internal/pkg/fileutil/fileutil.go` with shared `ToFileInfo`

### Phase 2: Deduplication
1. Remove `MountPoint` from `internal/config/config.go` — use `model.MountPoint`
2. Remove `ServerConfig` from `internal/model/config.go`
3. Remove `toFileInfo` from `file.go` and `search.go` — use `fileutil.ToFileInfo`
4. Remove `handleServiceError` from individual handlers — use `HandleServiceError`
5. Remove `writeError`/`writeJSON` from `auth.go`

### Phase 3: Resource Management
1. Add upload session cleanup goroutine
2. Add revoked token cleanup goroutine
3. Add job history cleanup goroutine
4. Update main.go to start all cleanup routines

### Phase 4: Cleanup
1. Add structured logging to key operations
2. Update all imports after file moves
3. Run `go mod tidy`
4. Run `go vet ./...`
5. Run `golangci-lint run`

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
│   └── constants.go           # Centralized magic numbers
├── handler/
│   ├── response.go            # writeError, writeJSON helpers
│   └── errors.go              # Centralized error handling
├── middleware/
│   └── ratelimit.go           # Rate limiting middleware
└── pkg/
    ├── fileutil/
    │   └── fileutil.go        # Shared ToFileInfo function
    └── logger/
        └── logger.go          # Structured logging helpers
```

## Files to Update (Placeholder → Documentation)

```
internal/handler/handler.go       # Add package documentation
internal/middleware/middleware.go # Add package documentation
internal/service/service.go       # Add package documentation
internal/model/model.go           # Add package documentation
```

## Files to Modify

```
internal/config/config.go        # Remove MountPoint, use model.MountPoint
internal/model/config.go         # Remove ServerConfig, keep only MountPoint
internal/handler/auth.go         # Remove writeError, writeJSON, ErrorResponse
internal/handler/file.go         # Use HandleServiceError, fileutil.ToFileInfo
internal/handler/stream.go       # Use HandleServiceError, add session cleanup
internal/handler/job.go          # Use HandleServiceError
internal/handler/search.go       # Use HandleServiceError, fileutil.ToFileInfo
internal/handler/websocket.go    # Configurable origins, use constants
internal/service/file.go         # Use fileutil.ToFileInfo
internal/service/search.go       # Use fileutil.ToFileInfo
internal/service/auth.go         # Add bcrypt, StartCleanup
internal/service/job.go          # Add StartCleanup
internal/websocket/client.go     # Use constants from config
cmd/server/main.go               # Remove hardcoded creds, start cleanup routines
config.yaml                      # Add users, websocket.allowed_origins
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

### High (Memory Leaks)
1. Upload sessions never cleaned up
2. Revoked tokens never cleaned up
3. Job history never cleaned up

### Medium (Code Quality)
4. MountPoint type duplicated
5. ServerConfig duplicated
6. toFileInfo duplicated
7. handleServiceError duplicated
8. writeError/writeJSON in wrong file
9. Magic numbers scattered

### Low (Cleanup)
10. Empty placeholder files
11. Missing structured logging

### Optional (Security — Skip for Homelab)
12. Hardcoded admin credentials
13. No rate limiting on auth endpoints
14. WebSocket allows all origins
