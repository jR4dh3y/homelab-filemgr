# Architecture Overview

This document describes the system architecture of the Homelab File Manager.

## System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Browser (Svelte SPA)                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐ │
│  │ File     │  │ Upload   │  │ Job      │  │ WebSocket        │ │
│  │ Browser  │  │ Manager  │  │ Monitor  │  │ Client           │ │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────────┬─────────┘ │
└───────┼─────────────┼─────────────┼─────────────────┼───────────┘
        │ REST        │ REST        │ REST            │ WS
        ▼             ▼             ▼                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Go Backend API Server                       │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    HTTP Router (Chi)                        ││
│  │  /api/v1/files/*  /api/v1/jobs/*  /api/v1/ws  /api/v1/auth  ││
│  └─────────────────────────────────────────────────────────────┘│
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐ │
│  │ File     │  │ Stream   │  │ Job      │  │ WebSocket        │ │
│  │ Handler  │  │ Handler  │  │ Handler  │  │ Hub              │ │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────────┬─────────┘ │
│       │             │             │                 │           │
│  ┌────▼─────────────▼─────────────▼─────────────────▼─────────┐ │
│  │                     Service Layer                          │ │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐            │ │
│  │  │ FileService│  │ JobService │  │ AuthService│            │ │
│  │  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘            │ │
│  └────────┼───────────────┼───────────────┼───────────────────┘ │
│           │               │               │                     │
│  ┌────────▼───────────────▼───────────────▼───────────────────┐ │
│  │                   Infrastructure Layer                      │ │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐            │ │
│  │  │ Filesystem │  │ JobQueue   │  │ Config     │            │ │
│  │  │ (afero)    │  │ (in-memory)│  │ (viper)    │            │ │
│  │  └────────────┘  └────────────┘  └────────────┘            │ │
│  └─────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Filesystem (Local + Mounts)                   │
│  /data/media    /mnt/nas    /home/user/documents                │
└─────────────────────────────────────────────────────────────────┘
```

## Backend Architecture

### Package Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Entry point, server setup
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loading (viper)
│   ├── handler/
│   │   ├── auth.go              # Authentication endpoints
│   │   ├── file.go              # File operations endpoints
│   │   ├── job.go               # Job management endpoints
│   │   ├── search.go            # Search endpoint
│   │   ├── stream.go            # Upload/download streaming
│   │   └── websocket.go         # WebSocket handler
│   ├── middleware/
│   │   ├── auth.go              # JWT validation
│   │   └── security.go          # Security headers, mount guard
│   ├── model/
│   │   ├── config.go            # Configuration models
│   │   ├── error.go             # Error types
│   │   ├── file.go              # File/directory models
│   │   └── job.go               # Job models and states
│   ├── service/
│   │   ├── auth.go              # JWT token management
│   │   ├── file.go              # File operations logic
│   │   ├── job.go               # Job execution and tracking
│   │   └── search.go            # File search logic
│   ├── websocket/
│   │   ├── client.go            # Individual client handling
│   │   └── hub.go               # Connection management
│   └── pkg/
│       ├── filesystem/
│       │   └── fs.go            # Filesystem abstraction
│       └── validator/
│           └── path.go          # Path validation
├── go.mod
└── go.sum
```

### Layer Responsibilities

#### Handlers (HTTP Layer)

- Parse HTTP requests
- Validate input
- Call services
- Format responses
- Handle errors

#### Services (Business Logic)

- Implement business rules
- Coordinate operations
- Manage transactions
- Emit events

#### Infrastructure

- Filesystem access (via afero)
- Configuration management
- Job queue management

### Key Components

#### FileService

Handles all file operations:
- Directory listing with pagination
- File metadata retrieval
- Create, rename, delete operations
- Mount point validation

#### JobService

Manages background operations:
- Worker pool for concurrent execution
- Progress tracking
- Cancellation support
- WebSocket notifications

#### StreamHandler

Handles large file transfers:
- Chunked uploads with resume support
- Range request downloads
- Checksum verification

#### WebSocket Hub

Real-time communication:
- Client connection management
- Job update broadcasting
- Ping/pong health checks

## Frontend Architecture

### Package Structure

```
frontend/
├── src/
│   ├── lib/
│   │   ├── api/
│   │   │   ├── auth.ts          # Auth API calls
│   │   │   ├── client.ts        # HTTP client wrapper
│   │   │   ├── files.ts         # File API calls
│   │   │   └── jobs.ts          # Job API calls
│   │   ├── components/
│   │   │   ├── Breadcrumb.svelte
│   │   │   ├── FileBrowser.svelte
│   │   │   ├── FileList.svelte
│   │   │   ├── JobMonitor.svelte
│   │   │   ├── SearchBar.svelte
│   │   │   ├── UploadDropzone.svelte
│   │   │   └── UploadProgress.svelte
│   │   ├── stores/
│   │   │   ├── auth.ts          # Authentication state
│   │   │   ├── files.ts         # File listing state
│   │   │   ├── jobs.ts          # Job tracking state
│   │   │   └── websocket.ts     # WebSocket connection
│   │   └── utils/
│   │       ├── format.ts        # Formatting utilities
│   │       └── upload.ts        # Chunked upload logic
│   └── routes/
│       ├── +layout.svelte       # App layout
│       ├── +page.svelte         # Home redirect
│       ├── login/
│       │   └── +page.svelte     # Login page
│       └── browse/
│           └── [...path]/
│               └── +page.svelte # File browser
├── svelte.config.js
└── vite.config.ts
```

### State Management

- **TanStack Query**: Server state (file listings, jobs)
- **Svelte Stores**: Client state (auth, WebSocket)
- **Runes ($state, $derived)**: Component state

### Component Hierarchy

```
+layout.svelte
├── Login Page
│   └── Login Form
└── Browse Page
    └── FileBrowser
        ├── Breadcrumb
        ├── SearchBar
        ├── FileList
        ├── UploadDropzone
        ├── UploadProgress
        └── JobMonitor
```

## Data Flow

### File Listing

```
User navigates to /browse/media/movies
         │
         ▼
┌─────────────────┐
│ SvelteKit Route │
│ [...path]       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ TanStack Query  │
│ useQuery()      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ API Client      │
│ GET /files/...  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Backend Handler │
│ FileHandler     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ FileService     │
│ List()          │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Filesystem      │
│ ReadDir()       │
└─────────────────┘
```

### Background Job

```
User initiates copy operation
         │
         ▼
┌─────────────────┐
│ POST /jobs      │
│ {type: "copy"}  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ JobService      │
│ Create()        │
└────────┬────────┘
         │
         ├──────────────────────┐
         ▼                      ▼
┌─────────────────┐    ┌─────────────────┐
│ Return job ID   │    │ Worker Pool     │
│ to client       │    │ Execute job     │
└─────────────────┘    └────────┬────────┘
                                │
                       ┌────────┴────────┐
                       ▼                 ▼
              ┌─────────────┐   ┌─────────────┐
              │ Copy files  │   │ Broadcast   │
              │ with        │   │ progress    │
              │ progress    │   │ via WS      │
              └─────────────┘   └─────────────┘
```

### Chunked Upload

```
┌──────────┐                    ┌──────────┐
│  Browser │                    │  Backend │
└────┬─────┘                    └────┬─────┘
     │                               │
     │ POST chunk 0                  │
     │──────────────────────────────>│
     │                               │ Write to temp
     │ 200 OK                        │
     │<──────────────────────────────│
     │                               │
     │ POST chunk 1                  │
     │──────────────────────────────>│
     │                               │
     │ ... repeat ...                │
     │                               │
     │ POST final chunk + checksum   │
     │──────────────────────────────>│
     │                               │ Verify checksum
     │                               │ Move to final path
     │ 201 Created                   │
     │<──────────────────────────────│
```

## Security Architecture

### Authentication Flow

```
┌─────────┐     ┌─────────┐     ┌─────────┐
│ Client  │     │ Backend │     │ JWT     │
└────┬────┘     └────┬────┘     └────┬────┘
     │               │               │
     │ POST /login   │               │
     │──────────────>│               │
     │               │ Validate      │
     │               │ credentials   │
     │               │               │
     │               │ Generate      │
     │               │──────────────>│
     │               │               │
     │               │ Access +      │
     │               │ Refresh token │
     │               │<──────────────│
     │ Tokens        │               │
     │<──────────────│               │
     │               │               │
     │ GET /files    │               │
     │ + Bearer token│               │
     │──────────────>│               │
     │               │ Validate      │
     │               │──────────────>│
     │               │               │
     │               │ Claims        │
     │               │<──────────────│
     │ Response      │               │
     │<──────────────│               │
```

### Security Layers

1. **JWT Authentication**: All API requests require valid token
2. **Path Validation**: Prevents directory traversal attacks
3. **Mount Point Guard**: Restricts access to configured directories
4. **Read-Only Enforcement**: Blocks writes on read-only mounts
5. **Security Headers**: CSP, X-Frame-Options, etc.

## Scalability Considerations

### Current Limitations

- Single server deployment
- In-memory job queue (not persistent)
- No clustering support

### Future Improvements

- Redis for job queue persistence
- Horizontal scaling with load balancer
- Distributed file operations
- Database for metadata caching

## Technology Choices

### Why Go for Backend?

- Excellent concurrency (goroutines)
- Fast compilation and execution
- Strong standard library
- Easy deployment (single binary)

### Why Svelte for Frontend?

- Minimal bundle size
- Reactive by default
- No virtual DOM overhead
- Great developer experience

### Why Chi Router?

- Lightweight and fast
- Middleware support
- Compatible with net/http
- Good documentation

### Why Afero?

- Filesystem abstraction
- Easy testing with in-memory FS
- Supports multiple backends
