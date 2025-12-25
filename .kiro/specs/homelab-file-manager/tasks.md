# Implementation Plan: Homelab File Manager

## Overview

This plan implements a high-performance file manager with a Go backend and Svelte frontend. Tasks are organized to build incrementally: core infrastructure first, then file operations, streaming, jobs, WebSocket, and finally the frontend.

## Tasks

- [x] 1. Set up project structure and dependencies
  - [x] 1.1 Initialize Go module and install dependencies
    - Create `backend/` directory with `go.mod`
    - Install chi, gorilla/websocket, jwt, afero, viper, zerolog, gopter
    - Create directory structure: `cmd/server/`, `internal/{config,handler,service,middleware,model,websocket,pkg/}`
    - _Requirements: All_

  - [x] 1.2 Initialize Svelte project with SvelteKit
    - SvelteKit scaffold already exists with Tailwind CSS configured
    - _Requirements: All_

  - [x] 1.3 Install frontend application dependencies
    - Install @tanstack/svelte-query for server state management
    - Create directory structure: `src/lib/{components,stores,api,utils}/`
    - _Requirements: All_

  - [x] 1.4 Create backend configuration system
    - Implement `internal/config/config.go` with viper
    - Define ServerConfig struct with mount points, JWT secret, port, chunk sizes
    - Create sample `config.yaml` file
    - _Requirements: 6.1_

- [x] 2. Implement core models and filesystem abstraction
  - [x] 2.1 Create data models
    - Implement `internal/model/file.go` with FileInfo, FileList, ListOptions
    - Implement `internal/model/job.go` with Job, JobType, JobState, JobUpdate
    - Implement `internal/model/config.go` with MountPoint, ServerConfig
    - _Requirements: 1.1, 4.1, 6.1_

  - [x] 2.2 Create filesystem abstraction layer
    - Implement `internal/pkg/filesystem/fs.go` wrapping afero
    - Add methods for common operations (ReadDir, Stat, Open, Create, Remove, Rename)
    - _Requirements: 1.1, 8.1, 8.2, 8.3_

  - [x] 2.3 Implement path validator
    - Implement `internal/pkg/validator/path.go` with SanitizePath function
    - Handle path traversal detection and prevention
    - Validate paths against mount point boundaries
    - _Requirements: 7.2, 6.2_

  - [x] 2.4 Write property test for path traversal prevention

    - **Property 13: Path Traversal Prevention**
    - **Validates: Requirements 7.2**

- [x] 3. Implement authentication system
  - [x] 3.1 Create auth service
    - Implement `internal/service/auth.go` with Login, Refresh, ValidateToken
    - Use golang-jwt/jwt for token generation and validation
    - Support access and refresh token pairs
    - _Requirements: 7.1, 7.4_

  - [x] 3.2 Create auth middleware
    - Implement `internal/middleware/auth.go` with JWTAuth middleware
    - Extract and validate Bearer tokens from Authorization header
    - Add user claims to request context
    - _Requirements: 7.1, 7.5_

  - [x] 3.3 Create auth handler
    - Implement `internal/handler/auth.go` with login, refresh, logout endpoints
    - _Requirements: 7.1, 7.4_

  - [x] 3.4 Write property test for authentication enforcement

    - **Property 12: Authentication Enforcement**
    - **Validates: Requirements 7.1, 7.5**

- [x] 4. Implement security middleware
  - [x] 4.1 Create security headers middleware
    - Implement `internal/middleware/security.go` with SecurityHeaders
    - Add CSP, X-Frame-Options, X-Content-Type-Options, X-XSS-Protection
    - _Requirements: 7.3_

  - [x] 4.2 Create mount point guard middleware
    - Implement MountPointGuard in `internal/middleware/security.go`
    - Validate all paths against configured mount points
    - Enforce read-only restrictions
    - _Requirements: 6.2, 6.3, 6.4_

  - [x] 4.3 Write property tests for security middleware

    - **Property 10: Mount Point Isolation**
    - **Property 11: Read-Only Mount Enforcement**
    - **Property 14: Security Headers Presence**
    - **Validates: Requirements 6.2, 6.3, 6.4, 7.3**

- [x] 5. Checkpoint - Core infrastructure
  - Ensure all tests pass, ask the user if questions arise.

- [x] 6. Implement file service and handlers
  - [x] 6.1 Create file service
    - Implement `internal/service/file.go` with List, GetInfo, CreateDir, Rename, Delete
    - Add pagination support with sorting and filtering
    - Use filesystem abstraction layer
    - _Requirements: 1.1, 1.3, 8.1, 8.2, 8.3_

  - [x] 6.2 Create file handler
    - Implement `internal/handler/file.go` with REST endpoints
    - GET /api/v1/files (list roots), GET /api/v1/files/*path (list/info)
    - POST /api/v1/files/*path (create dir), PUT (rename), DELETE (delete)
    - _Requirements: 1.1, 1.3, 1.5, 8.1, 8.2, 8.3, 8.4, 8.5_

  - [x] 6.3 Write property tests for file operations

    - **Property 1: Directory Listing Metadata Completeness**
    - **Property 2: Pagination Correctness**
    - **Property 3: Non-Existent Path Returns 404**
    - **Property 15: File Rename Correctness**
    - **Property 16: Directory Creation Correctness**
    - **Property 17: File Deletion Correctness**
    - **Validates: Requirements 1.1, 1.3, 1.5, 8.1, 8.2, 8.3**

- [x] 7. Implement search functionality
  - [x] 7.1 Create search service
    - Implement `internal/service/search.go` with recursive directory search
    - Support case-insensitive name matching
    - Return results with full metadata
    - _Requirements: 9.1, 9.2, 9.3_

  - [x] 7.2 Add search endpoint to file handler
    - GET /api/v1/search?path=&q=
    - Validate non-empty query
    - _Requirements: 9.1, 9.2, 9.3_

  - [x] 7.3 Write property test for search correctness

    - **Property 18: Search Result Correctness**
    - **Validates: Requirements 9.1, 9.2**

- [x] 8. Implement streaming upload/download
  - [x] 8.1 Create stream handler for downloads
    - Implement `internal/handler/stream.go` with Download handler
    - Use http.ServeContent for Range header support
    - Set Content-Disposition and MIME type headers
    - _Requirements: 3.1, 3.2, 3.4_

  - [x] 8.2 Create stream handler for chunked uploads
    - Implement Upload handler with chunk assembly
    - Track upload sessions with upload ID
    - Support resumable uploads from last chunk
    - Verify checksum on completion
    - _Requirements: 2.1, 2.3, 2.4, 2.5_

  - [x] 8.3 Write property tests for streaming

    - **Property 4: Upload/Download Round-Trip Integrity**
    - **Property 5: Resumable Upload Correctness**
    - **Property 6: Range Request Correctness**
    - **Validates: Requirements 2.3, 2.5, 3.1, 3.2**

- [x] 9. Checkpoint - File operations complete
  - Ensure all tests pass, ask the user if questions arise.

- [x] 10. Implement WebSocket hub
  - [x] 10.1 Create WebSocket hub
    - Implement `internal/websocket/hub.go` with client registration/unregistration
    - Implement broadcast mechanism for job updates
    - Handle concurrent access with sync.RWMutex
    - _Requirements: 5.1, 5.2, 5.3_

  - [x] 10.2 Create WebSocket client handler
    - Implement `internal/websocket/client.go` for individual connections
    - Handle ping/pong for connection health
    - Parse subscribe/unsubscribe messages
    - _Requirements: 5.1, 5.2, 5.3_

  - [x] 10.3 Create WebSocket handler with auth
    - Implement `internal/handler/websocket.go`
    - Authenticate WebSocket upgrade requests
    - Register clients with hub
    - _Requirements: 5.1_

- [x] 11. Implement job service
  - [x] 11.1 Create job service and executor
    - Implement `internal/service/job.go` with Create, Get, List, Cancel
    - Implement job executor with worker pool
    - Execute copy, move, delete as background jobs
    - Track progress and broadcast updates via WebSocket hub
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6_

  - [x] 11.2 Create job handler
    - Implement `internal/handler/job.go` with REST endpoints
    - GET /api/v1/jobs (list), GET /api/v1/jobs/:id (status)
    - POST /api/v1/jobs (create), DELETE /api/v1/jobs/:id (cancel)
    - _Requirements: 4.1, 4.4, 4.5_

  - [x] 11.3 Write property tests for job service

    - **Property 7: Job Progress Monotonicity**
    - **Property 8: Job Completion Notification**
    - **Property 9: Job Cancellation Cleanup**
    - **Validates: Requirements 4.2, 4.3, 4.5**

- [x] 12. Wire up backend server
  - [x] 12.1 Create main server entry point
    - Implement `cmd/server/main.go`
    - Load configuration with viper
    - Initialize all services and handlers
    - Set up chi router with middleware chain
    - Start HTTP server and WebSocket hub
    - _Requirements: 6.1, 10.4_

  - [x] 12.2 Create router configuration
    - Wire all routes with appropriate middleware
    - Apply auth middleware to protected routes
    - Apply security middleware globally
    - _Requirements: 7.1, 7.3_

- [x] 13. Checkpoint - Backend complete
  - Ensure all tests pass, ask the user if questions arise.

- [x] 14. Implement frontend API layer
  - [x] 14.1 Create HTTP client wrapper
    - Implement `src/lib/api/client.ts` with fetch wrapper
    - Handle JWT token injection and refresh
    - Implement error handling and response parsing
    - _Requirements: 7.1_

  - [x] 14.2 Create file API module
    - Implement `src/lib/api/files.ts` with list, create, rename, delete, search
    - _Requirements: 1.1, 8.1, 8.2, 8.3, 9.1_

  - [x] 14.3 Create job API module
    - Implement `src/lib/api/jobs.ts` with list, get, create, cancel
    - _Requirements: 4.1, 4.4, 4.5_

  - [x] 14.4 Create auth API module
    - Implement `src/lib/api/auth.ts` with login, refresh, logout
    - _Requirements: 7.1_

- [x] 15. Implement frontend stores
  - [x] 15.1 Create auth store
    - Implement `src/lib/stores/auth.ts` with token management
    - Handle login state and token refresh
    - _Requirements: 7.1, 7.4_

  - [x] 15.2 Create files store
    - Implement `src/lib/stores/files.ts` with current path and listing state
    - Use @tanstack/svelte-query for server state
    - _Requirements: 1.1, 1.2_

  - [x] 15.3 Create jobs store
    - Implement `src/lib/stores/jobs.ts` with active jobs tracking
    - _Requirements: 4.4_

  - [x] 15.4 Create WebSocket store
    - Implement `src/lib/stores/websocket.ts` with connection management
    - Handle reconnection with exponential backoff
    - Process job update messages
    - _Requirements: 5.1, 5.4_

- [x] 16. Implement frontend utilities
  - [x] 16.1 Create chunked upload utility
    - Implement `src/lib/utils/upload.ts` with chunk splitting and upload
    - Track progress and support resume
    - Calculate checksums
    - _Requirements: 2.1, 2.2, 2.3, 2.5_

  - [x] 16.2 Create formatting utilities
    - Implement `src/lib/utils/format.ts` with file size and date formatting
    - _Requirements: 1.1_

- [x] 17. Implement frontend components
  - [x] 17.1 Create FileBrowser component
    - Implement `src/lib/components/FileBrowser.svelte` as main container
    - Compose Breadcrumb, SearchBar, FileList, UploadDropzone
    - _Requirements: 1.2, 1.4_

  - [x] 17.2 Create FileList component
    - Implement `src/lib/components/FileList.svelte` with sortable columns
    - Display name, size, modified date, type
    - Handle file/folder click navigation
    - _Requirements: 1.1, 1.2_

  - [x] 17.3 Create Breadcrumb component
    - Implement `src/lib/components/Breadcrumb.svelte` for path navigation
    - _Requirements: 1.4_

  - [x] 17.4 Create UploadDropzone component
    - Implement `src/lib/components/UploadDropzone.svelte` with drag-drop
    - Trigger chunked upload on file drop
    - _Requirements: 2.1, 2.2_

  - [x] 17.5 Create UploadProgress component
    - Implement `src/lib/components/UploadProgress.svelte` showing active uploads
    - Display progress bars and cancel buttons
    - _Requirements: 2.2_

  - [x] 17.6 Create JobMonitor component
    - Implement `src/lib/components/JobMonitor.svelte` showing background jobs
    - Display progress, status, and cancel option
    - _Requirements: 4.2, 4.4, 4.5_

  - [x] 17.7 Create SearchBar component
    - Implement `src/lib/components/SearchBar.svelte` with search input
    - Display loading state during search
    - _Requirements: 9.1, 9.4_

- [x] 18. Implement frontend routes
  - [x] 18.1 Create layout and login page
    - Update `src/routes/+layout.svelte` with auth check
    - Implement `src/routes/login/+page.svelte` with login form
    - _Requirements: 7.1_

  - [x] 18.2 Create browse page
    - Implement `src/routes/browse/[...path]/+page.svelte`
    - Load directory listing based on path parameter
    - Integrate FileBrowser component
    - _Requirements: 1.1, 1.2, 1.4_

  - [x] 18.3 Update main page with redirect
    - Update `src/routes/+page.svelte` to redirect to browse
    - _Requirements: 1.2_

- [x] 19. Final checkpoint - Full integration
  - Ensure all tests pass, ask the user if questions arise.
  - Verify frontend connects to backend
  - Test file browsing, upload, download, and job execution

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties
- Unit tests validate specific examples and edge cases
- Backend uses afero.MemMapFs for testing without real filesystem
- Frontend uses @tanstack/svelte-query for efficient server state management
