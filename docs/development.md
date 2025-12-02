# Development Guide

This guide covers setting up a local development environment for the Homelab File Manager.

## Prerequisites

- Go 1.23+
- Node.js 22+ (or Bun)
- Git

## Quick Start

### Clone Repository

```bash
git clone https://github.com/yourusername/homelab-file-manager.git
cd homelab-file-manager
```

### Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Create test directories
mkdir -p /tmp/filemanager/media
mkdir -p /tmp/filemanager/documents

# Update config.yaml with test paths
cat > config.yaml << EOF
port: 8080
host: "0.0.0.0"
jwt_secret: "dev-secret-change-in-production"
max_upload_mb: 1024
chunk_size_mb: 5

mount_points:
  - name: "media"
    path: "/tmp/filemanager/media"
    read_only: false
  - name: "documents"
    path: "/tmp/filemanager/documents"
    read_only: false
EOF

# Run the server
go run ./cmd/server
```

### Frontend Setup

```bash
cd frontend

# Install dependencies
npm install
# or
bun install

# Start development server
npm run dev
# or
bun dev
```

Access the application at http://localhost:5173

## Project Structure

```
homelab-file-manager/
├── backend/                 # Go backend
│   ├── cmd/server/         # Entry point
│   ├── internal/           # Internal packages
│   │   ├── config/         # Configuration
│   │   ├── handler/        # HTTP handlers
│   │   ├── middleware/     # Middleware
│   │   ├── model/          # Data models
│   │   ├── service/        # Business logic
│   │   ├── websocket/      # WebSocket hub
│   │   └── pkg/            # Shared utilities
│   ├── config.yaml         # Configuration file
│   ├── go.mod
│   └── go.sum
├── frontend/               # Svelte frontend
│   ├── src/
│   │   ├── lib/           # Components, stores, utils
│   │   └── routes/        # SvelteKit routes
│   ├── package.json
│   └── svelte.config.js
├── docs/                   # Documentation
├── nginx/                  # Nginx configuration
├── docker-compose.yml
└── README.md
```

## Running Tests

### Backend Tests

```bash
cd backend

# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test -v ./internal/service/...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Property-Based Tests

The backend includes property-based tests using gopter:

```bash
# Run property tests (may take longer)
go test -v -timeout 300s ./internal/service/...
```

### Frontend Tests

```bash
cd frontend

# Type checking
npm run check

# Linting
npm run lint

# Format code
npm run format
```

## Code Style

### Go

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Run `go vet` before committing

```bash
# Format code
gofmt -w .

# Run vet
go vet ./...
```

### TypeScript/Svelte

- Use Prettier for formatting
- Follow ESLint rules

```bash
# Format
npm run format

# Lint
npm run lint
```

## API Development

### Adding a New Endpoint

1. **Define the model** in `internal/model/`:
   ```go
   type NewFeature struct {
       ID   string `json:"id"`
       Name string `json:"name"`
   }
   ```

2. **Create the service** in `internal/service/`:
   ```go
   type NewFeatureService interface {
       Get(ctx context.Context, id string) (*model.NewFeature, error)
   }
   ```

3. **Create the handler** in `internal/handler/`:
   ```go
   func (h *NewFeatureHandler) Get(w http.ResponseWriter, r *http.Request) {
       // Implementation
   }
   ```

4. **Register routes** in `cmd/server/main.go`:
   ```go
   r.Route("/newfeature", func(r chi.Router) {
       newFeatureHandler.RegisterRoutes(r)
   })
   ```

### Adding Frontend API

1. **Add API function** in `src/lib/api/`:
   ```typescript
   export async function getNewFeature(id: string): Promise<NewFeature> {
       return api.get<NewFeature>(`/newfeature/${id}`);
   }
   ```

2. **Create store** if needed in `src/lib/stores/`:
   ```typescript
   export const newFeatureStore = writable<NewFeature | null>(null);
   ```

3. **Use in component**:
   ```svelte
   <script lang="ts">
       import { getNewFeature } from '$lib/api/newfeature';
       
       let feature = $state<NewFeature | null>(null);
       
       async function load() {
           feature = await getNewFeature('123');
       }
   </script>
   ```

## Debugging

### Backend Debugging

Using VS Code:

```json
// .vscode/launch.json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Backend",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/backend/cmd/server",
            "cwd": "${workspaceFolder}/backend"
        }
    ]
}
```

Using Delve:

```bash
cd backend
dlv debug ./cmd/server
```

### Frontend Debugging

Use browser DevTools or VS Code debugger with the Svelte extension.

### Logging

Backend uses zerolog:

```go
import "github.com/rs/zerolog/log"

log.Info().Str("path", path).Msg("Processing request")
log.Error().Err(err).Msg("Operation failed")
```

## Common Tasks

### Update Dependencies

```bash
# Backend
cd backend
go get -u ./...
go mod tidy

# Frontend
cd frontend
npm update
```

### Generate Mocks (if needed)

```bash
go install github.com/golang/mock/mockgen@latest
mockgen -source=internal/service/file.go -destination=internal/service/mock_file.go
```

### Build for Production

```bash
# Backend
cd backend
CGO_ENABLED=0 go build -ldflags="-w -s" -o server ./cmd/server

# Frontend
cd frontend
npm run build
```

## Troubleshooting

### Port already in use

```bash
# Find process using port
lsof -i :8080
# or on Windows
netstat -ano | findstr :8080

# Kill process
kill -9 <PID>
```

### Go module issues

```bash
go clean -modcache
go mod download
```

### Node module issues

```bash
rm -rf node_modules
rm package-lock.json
npm install
```

### CORS errors

Ensure the backend is running and the frontend is configured to use the correct API URL.

For development, the Vite proxy handles CORS. Check `vite.config.ts`.

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Run tests: `go test ./...` and `npm run check`
5. Commit: `git commit -m "Add my feature"`
6. Push: `git push origin feature/my-feature`
7. Create a Pull Request

### Commit Messages

Follow conventional commits:

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation
- `refactor:` Code refactoring
- `test:` Adding tests
- `chore:` Maintenance
