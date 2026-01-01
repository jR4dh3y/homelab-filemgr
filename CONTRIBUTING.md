# Contributing Guidelines

This document defines coding standards and patterns for the Homelab File Manager project. Follow these guidelines when implementing features or making changes.

---

## Core Principles

### 1. Don't Reinvent the Wheel

Before writing new code, **check if it already exists**:

```
Backend:
- File utilities      → internal/pkg/fileutil/
- Path validation     → internal/pkg/validator/
- Filesystem ops      → internal/pkg/filesystem/
- Error handling      → internal/handler/errors.go
- Response helpers    → internal/handler/response.go
- Constants           → internal/config/constants.go

Frontend:
- File type detection → src/lib/utils/fileTypes.ts
- Formatting          → src/lib/utils/format.ts
- API calls           → src/lib/api/
- Storage access      → src/lib/utils/storage.ts
- Config values       → src/lib/config.ts
- UI components       → src/lib/components/ui/
```

**If you need something that might be reusable, check these locations first.**

### 2. No Monolithic Code

Break down large functions and components:

```go
// BAD - 200 line function doing everything
func (h *Handler) DoEverything(w http.ResponseWriter, r *http.Request) {
    // parse request
    // validate input
    // check permissions
    // do business logic
    // format response
    // handle errors
    // ... 200 lines later
}

// GOOD - small focused functions
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    req, err := h.parseCreateRequest(r)
    if err != nil {
        writeError(w, err.Error(), model.ErrCodeValidationError, http.StatusBadRequest)
        return
    }
    
    result, err := h.service.Create(r.Context(), req)
    if err != nil {
        HandleServiceError(w, err)
        return
    }
    
    writeJSON(w, result, http.StatusCreated)
}
```

```svelte
<!-- BAD - 500 line component -->
<script>
  // 200 lines of logic
</script>
<!-- 300 lines of markup -->

<!-- GOOD - composed from smaller components -->
<script>
  import { FileList, Toolbar, Sidebar } from '$lib/components';
</script>

<div class="layout">
  <Sidebar />
  <main>
    <Toolbar />
    <FileList />
  </main>
</div>
```

### 3. Single Source of Truth

Never duplicate definitions:

```typescript
// BAD - same extensions in multiple files
// FileList.svelte
const imageExts = ['jpg', 'png', 'gif'];
// FilePreview.svelte  
const imageExts = ['jpg', 'png', 'gif'];
// format.ts
const imageExts = ['jpg', 'png', 'gif'];

// GOOD - one place, import everywhere
// fileTypes.ts
export const FILE_EXTENSIONS = {
  image: ['jpg', 'jpeg', 'png', 'gif', 'webp'],
  // ...
};

// Other files
import { FILE_EXTENSIONS } from '$lib/utils/fileTypes';
```

### 4. Keep Functions Small and Focused

Each function should do ONE thing:

```go
// BAD - function does too much
func ProcessFile(path string) error {
    // validate path
    // read file
    // parse content
    // transform data
    // write output
    // send notification
}

// GOOD - single responsibility
func ValidatePath(path string) error { ... }
func ReadFile(path string) ([]byte, error) { ... }
func ParseContent(data []byte) (*Content, error) { ... }
func TransformData(c *Content) (*Output, error) { ... }
```

### 5. Reuse Before You Create

When you need functionality:

1. **Search the codebase** — it might already exist
2. **Check if existing code can be extended** — add to existing utils
3. **Only then create new** — and put it in the right place for others to reuse

```typescript
// Need to format a file size?
// DON'T create a new function
// DO check src/lib/utils/format.ts — formatFileSize() exists!

// Need to detect file type?
// DON'T create extension arrays
// DO use src/lib/utils/fileTypes.ts
```

### 6. Extract Shared Logic

If you write similar code twice, extract it:

```go
// You wrote this in handler A:
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(status)
json.NewEncoder(w).Encode(data)

// And again in handler B... STOP!
// Extract to a shared helper:
func writeJSON(w http.ResponseWriter, data interface{}, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
```

### 7. Consistent Patterns

Follow existing patterns in the codebase:

- Handlers follow the same structure
- Services follow the same interface pattern
- Components follow the same props pattern
- Error handling is consistent everywhere

**When in doubt, look at how similar things are done elsewhere in the code.**

---

## Project Structure

```
backend/                    # Go API server
├── cmd/server/            # Application entry point
├── internal/
│   ├── config/            # Configuration loading and constants
│   ├── handler/           # HTTP handlers (thin layer, delegates to services)
│   ├── middleware/        # HTTP middleware (auth, security, guards)
│   ├── model/             # Data models and types
│   ├── pkg/               # Shared utilities
│   │   ├── filesystem/    # Filesystem abstraction (uses afero)
│   │   ├── fileutil/      # File helper functions
│   │   └── validator/     # Path validation and sanitization
│   ├── service/           # Business logic layer
│   └── websocket/         # WebSocket hub and client management

frontend/                   # SvelteKit application
├── src/
│   ├── lib/
│   │   ├── api/           # API client functions
│   │   ├── components/    # Svelte components
│   │   │   ├── ui/        # Reusable base components
│   │   │   └── preview/   # File preview components
│   │   ├── stores/        # Svelte stores (state management)
│   │   ├── types/         # TypeScript type definitions
│   │   └── utils/         # Utility functions
│   └── routes/            # SvelteKit routes/pages
```

---

## Backend (Go) Guidelines

### Architecture Pattern

Follow the **Handler → Service → Model** pattern:

```
HTTP Request → Handler → Service → Model/Filesystem
                  ↓
            HTTP Response
```

- **Handlers**: Parse requests, validate input, call services, format responses
- **Services**: Business logic, filesystem operations, error handling
- **Models**: Data structures, no logic

### Creating a New Handler

```go
// internal/handler/example.go
package handler

type ExampleHandler struct {
    exampleService service.ExampleService
}

func NewExampleHandler(svc service.ExampleService) *ExampleHandler {
    return &ExampleHandler{exampleService: svc}
}

func (h *ExampleHandler) RegisterRoutes(r chi.Router) {
    r.Get("/", h.List)
    r.Post("/", h.Create)
    r.Get("/{id}", h.Get)
}

func (h *ExampleHandler) List(w http.ResponseWriter, r *http.Request) {
    result, err := h.exampleService.List(r.Context())
    if err != nil {
        HandleServiceError(w, err)  // Use centralized error handler
        return
    }
    writeJSON(w, result, http.StatusOK)  // Use centralized response helper
}
```

### Creating a New Service

```go
// internal/service/example.go
package service

// Define errors at package level
var (
    ErrExampleNotFound = errors.New("example not found")
    ErrInvalidExample  = errors.New("invalid example")
)

// Define interface first
type ExampleService interface {
    List(ctx context.Context) ([]model.Example, error)
    Get(ctx context.Context, id string) (*model.Example, error)
    Create(ctx context.Context, params model.ExampleParams) (*model.Example, error)
}

// Implement with unexported struct
type exampleService struct {
    fs          filesystem.FS
    mountPoints []model.MountPoint
}

type ExampleServiceConfig struct {
    MountPoints []model.MountPoint
}

func NewExampleService(fsys filesystem.FS, cfg ExampleServiceConfig) ExampleService {
    return &exampleService{
        fs:          fsys,
        mountPoints: cfg.MountPoints,
    }
}
```

### Error Handling

1. Services return domain errors (e.g., `ErrPathNotFound`)
2. Handlers convert to HTTP responses using `HandleServiceError()`
3. Add new error mappings to `internal/handler/errors.go`

```go
// Adding a new error mapping
var serviceErrorMappings = []ErrorMapping{
    // ... existing mappings
    {service.ErrExampleNotFound, "Example not found", model.ErrCodeNotFound, http.StatusNotFound},
}
```

### Constants and Configuration

Put magic numbers in `internal/config/constants.go`:

```go
// BAD - magic number in code
buf := make([]byte, 1024*1024)

// GOOD - use named constant
buf := make([]byte, config.FileCopyBufferSize)
```

### Filesystem Operations

Always use the `filesystem.FS` interface, never `os` package directly:

```go
// BAD
file, err := os.Open(path)

// GOOD
file, err := s.fs.Open(path)
```

### Path Validation

Always validate paths against mount points:

```go
mount, fsPath, err := validator.ValidatePathAgainstMounts(path, s.mountPoints)
if err != nil {
    return nil, ErrMountPointNotFound
}
```

---

## Frontend (SvelteKit) Guidelines

### Component Structure

Use Svelte 5 runes syntax:

```svelte
<script lang="ts">
  import type { Snippet } from 'svelte';
  
  interface Props {
    title: string;
    disabled?: boolean;
    children: Snippet;
    onclick?: () => void;
  }
  
  let { title, disabled = false, children, onclick }: Props = $props();
  
  // Reactive state
  let count = $state(0);
  
  // Derived values
  let doubled = $derived(count * 2);
</script>

<button {disabled} {onclick}>
  {@render children()}
</button>
```

### Styling

Use Tailwind CSS classes, NOT `<style>` blocks:

```svelte
<!-- BAD -->
<button class="my-btn">Click</button>
<style>
  .my-btn { background: blue; padding: 8px; }
</style>

<!-- GOOD -->
<button class="bg-blue-500 px-4 py-2 rounded hover:bg-blue-600">
  Click
</button>
```

Use design tokens from `app.css`:

```svelte
<!-- Use semantic color names -->
<div class="bg-surface-primary text-text-primary border-border-primary">
```

### State Management

Use Svelte 5 runes for stores (`.svelte.ts` files):

```typescript
// src/lib/stores/example.svelte.ts
interface ExampleState {
  items: Item[];
  loading: boolean;
  error: string | null;
}

let state = $state<ExampleState>({
  items: [],
  loading: false,
  error: null,
});

export const exampleStore = {
  get items() { return state.items; },
  get loading() { return state.loading; },
  get error() { return state.error; },
  
  async load() {
    state.loading = true;
    state.error = null;
    try {
      state.items = await api.getItems();
    } catch (e) {
      state.error = e.message;
    } finally {
      state.loading = false;
    }
  },
};
```

### API Calls

Use the API client in `src/lib/api/`:

```typescript
// src/lib/api/example.ts
import { apiClient } from './client';

export async function getExamples(): Promise<Example[]> {
  return apiClient.get('/examples');
}

export async function createExample(data: CreateExampleRequest): Promise<Example> {
  return apiClient.post('/examples', data);
}
```

### File Type Utilities

All file type logic goes in `src/lib/utils/fileTypes.ts`:

```typescript
// DON'T duplicate extension arrays in components
// DO import from fileTypes.ts
import { getFileCategory, getPreviewType, FILE_EXTENSIONS } from '$lib/utils/fileTypes';
```

### Configuration

Put magic numbers in `src/lib/config.ts`:

```typescript
// BAD
const REFRESH_INTERVAL = 14 * 60 * 1000;

// GOOD
import { CONFIG } from '$lib/config';
const interval = CONFIG.auth.tokenRefreshIntervalMs;
```

---

## Shared Patterns

### Naming Conventions

| Type | Backend (Go) | Frontend (TS) |
|------|--------------|---------------|
| Files | `snake_case.go` | `camelCase.ts` or `PascalCase.svelte` |
| Types/Structs | `PascalCase` | `PascalCase` |
| Functions | `PascalCase` (exported), `camelCase` (private) | `camelCase` |
| Constants | `PascalCase` or `SCREAMING_SNAKE` | `SCREAMING_SNAKE` |
| Variables | `camelCase` | `camelCase` |

### API Response Format

All API responses follow this structure:

```json
// Success
{
  "data": { ... }
}

// Error
{
  "error": "Human readable message",
  "code": "ERROR_CODE",
  "details": "Optional details"
}
```

### Error Codes

Use constants from `internal/model/error.go`:

```go
const (
    ErrCodeNotFound         = "NOT_FOUND"
    ErrCodeAccessDenied     = "ACCESS_DENIED"
    ErrCodeValidationError  = "VALIDATION_ERROR"
    ErrCodeInternalError    = "INTERNAL_ERROR"
    // ... etc
)
```

---

## Do's and Don'ts

### DO

- ✅ **Search before creating** — check if functionality exists
- ✅ **Reuse existing utilities** — don't duplicate code
- ✅ **Keep functions small** — single responsibility
- ✅ **Extract shared logic** — if you write it twice, make it reusable
- ✅ **Follow existing patterns** — consistency matters
- ✅ Use the filesystem abstraction (`filesystem.FS`)
- ✅ Validate all paths against mount points
- ✅ Use centralized error handling
- ✅ Use Tailwind classes for styling
- ✅ Use Svelte 5 runes (`$state`, `$derived`, `$props`)
- ✅ Put constants in config files
- ✅ Add context cancellation support for long operations
- ✅ Use interfaces for services (enables testing)

### DON'T

- ❌ **Create monolithic functions** — break them down
- ❌ **Duplicate code** — extract and reuse
- ❌ **Reinvent existing utilities** — check first
- ❌ **Copy-paste code** — if you're copying, you should be extracting
- ❌ Use `os` package directly for file operations
- ❌ Hardcode paths or magic numbers
- ❌ Add `<style>` blocks to Svelte components
- ❌ Duplicate type definitions across files
- ❌ Use legacy Svelte 4 stores (`writable`, `derived`)
- ❌ Access `localStorage` directly (use `storage.ts`)
- ❌ Duplicate file extension arrays (use `fileTypes.ts`)

---

## Adding a New Feature Checklist

### Backend

- [ ] Create model in `internal/model/`
- [ ] Create service interface and implementation in `internal/service/`
- [ ] Create handler in `internal/handler/`
- [ ] Register routes in `cmd/server/main.go`
- [ ] Add error mappings to `internal/handler/errors.go`
- [ ] Add any new constants to `internal/config/constants.go`
- [ ] Write tests

### Frontend

- [ ] Add types to `src/lib/types/`
- [ ] Add API functions to `src/lib/api/`
- [ ] Create store if needed in `src/lib/stores/`
- [ ] Create components in `src/lib/components/`
- [ ] Add route if needed in `src/routes/`
- [ ] Use Tailwind for all styling
- [ ] Add any new constants to `src/lib/config.ts`


---

## Before You Start Coding

### Ask Yourself:

1. **Does this already exist?**
   - Search the codebase for similar functionality
   - Check the utility folders listed in "Core Principles"

2. **Can I extend something existing?**
   - Maybe a utility just needs one more function
   - Maybe a component just needs one more prop

3. **Where should this live?**
   - Is it reusable? → Put in `utils/` or `pkg/`
   - Is it specific to one feature? → Keep it local
   - Is it a constant? → Put in config

4. **Am I duplicating anything?**
   - If you're copying code, stop and extract it
   - If you're defining the same type twice, use the existing one

5. **Is this function doing too much?**
   - Can you describe it in one sentence without "and"?
   - If not, break it down

---

## Quick Reference: Where Things Go

| What you need | Backend location | Frontend location |
|---------------|------------------|-------------------|
| Constants/magic numbers | `internal/config/constants.go` | `src/lib/config.ts` |
| Type definitions | `internal/model/` | `src/lib/types/` |
| API response helpers | `internal/handler/response.go` | — |
| Error handling | `internal/handler/errors.go` | — |
| File type detection | — | `src/lib/utils/fileTypes.ts` |
| Formatting (size, date) | — | `src/lib/utils/format.ts` |
| Path validation | `internal/pkg/validator/` | — |
| Filesystem operations | `internal/pkg/filesystem/` | — |
| Reusable UI components | — | `src/lib/components/ui/` |
| API client functions | — | `src/lib/api/` |
| State management | — | `src/lib/stores/` |
| localStorage access | — | `src/lib/utils/storage.ts` |
