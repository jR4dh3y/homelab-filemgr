<h1 align="center">Homelab File Manager</h1>

<p align="center">
  <strong>A modern, self-hosted file manager for your homelab</strong>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#documentation">Documentation</a> •
  <a href="#tech-stack">Tech Stack</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/SvelteKit-FF3E00?style=flat-square&logo=svelte&logoColor=white" alt="SvelteKit">
  <img src="https://img.shields.io/badge/TypeScript-3178C6?style=flat-square&logo=typescript&logoColor=white" alt="TypeScript">
  <img src="https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white" alt="Docker">
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License">
</p>



## Features

| Feature | Description |
|---------|-------------|
| **Multi-Mount Browsing** | Browse files across multiple mount points with a clean, intuitive interface |
| **Chunked Uploads** | Upload large files with resumable transfers and progress tracking |
| **Real-time Updates** | WebSocket-powered live updates for file changes and job progress |
| **Background Jobs** | Copy, move, and delete operations run asynchronously with monitoring |
| **Fast Search** | Recursive directory scanning with name filtering |
| **Secure by Default** | JWT auth, path traversal protection, configurable credentials, rate limiting |
| **Configurable Access** | Per-user credentials, read-only mounts, WebSocket origin restrictions |

---

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/yourusername/homelab-filemanager.git
cd homelab-filemanager

# Copy environment file and configure
cp .env.example .env

# Start the services
docker compose up -d
```

Access the web interface at `http://localhost:3000`

### Manual Setup

See [Development Guide](docs/development.md) for detailed instructions.

---


## Documentation

| Document | Description |
|----------|-------------|
| [API Reference](docs/api.md) | REST API endpoints and WebSocket events |
| [Architecture](docs/architecture.md) | System design and component overview |
| [Configuration](docs/configuration.md) | Environment variables and config options |
| [Development](docs/development.md) | Local setup and development workflow |
| [Docker](docs/docker.md) | Container deployment and orchestration |
| [Security](docs/security.md) | Authentication, authorization, and best practices |

---

## Tech Stack

**Backend**
- Go with Chi router
- JWT authentication
- WebSocket support (Gorilla)
- Afero filesystem abstraction

**Frontend**
- SvelteKit
- TypeScript
- Tailwind CSS

**Infrastructure**
- Docker & Docker Compose
- Nginx reverse proxy
- Multi-stage builds

---

## Project Structure

```
homelab-filemanager/
├── backend/                 # Go backend service
│   ├── cmd/server/          # Application entrypoint
│   ├── internal/
│   │   ├── config/          # Configuration loading
│   │   ├── handler/         # HTTP handlers
│   │   ├── middleware/      # Auth & security middleware
│   │   ├── model/           # Data models
│   │   ├── pkg/             # Shared utilities
│   │   ├── service/         # Business logic
│   │   └── websocket/       # WebSocket hub & clients
│   └── Dockerfile
├── frontend/                # SvelteKit frontend
│   ├── src/
│   │   ├── lib/
│   │   │   ├── api/         # API client modules
│   │   │   ├── components/  # Svelte components
│   │   │   ├── stores/      # Svelte stores
│   │   │   └── utils/       # Helper functions
│   │   └── routes/          # SvelteKit routes
│   └── Dockerfile
├── nginx/                   # Reverse proxy config
├── docs/                    # Documentation
└── docker-compose.yml       # Container orchestration
```

---

## Contributing

Contributions are welcome! Please read the [Development Guide](docs/development.md) before submitting a PR.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
