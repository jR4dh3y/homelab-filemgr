# Configuration Guide

The Homelab File Manager is configured through a YAML configuration file and environment variables.

## Configuration File

The backend reads configuration from `config.yaml`. The default location is:
- `./config.yaml` (current directory)
- `/app/config.yaml` (Docker container)

Override with the `-config` flag:
```bash
./server -config /path/to/config.yaml
```

## Full Configuration Example

```yaml
# Server settings
port: 8080
host: "0.0.0.0"

# JWT authentication secret
# IMPORTANT: Change this in production!
# Generate with: openssl rand -base64 32
jwt_secret: "change-me-in-production-use-a-long-random-string"

# Upload settings
max_upload_mb: 10240  # Maximum upload size (10GB)
chunk_size_mb: 5      # Chunk size for uploads (5MB)

# Mount points - directories accessible through the file manager
mount_points:
  - name: "media"
    path: "/data/media"
    read_only: false

  - name: "documents"
    path: "/home/user/documents"
    read_only: false

  - name: "backups"
    path: "/mnt/backups"
    read_only: true
```

## Configuration Options

### Server Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `port` | int | 8080 | HTTP server port |
| `host` | string | "0.0.0.0" | Bind address |
| `jwt_secret` | string | (required) | Secret for JWT signing |

### Upload Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `max_upload_mb` | int | 10240 | Maximum upload size in MB |
| `chunk_size_mb` | int | 5 | Chunk size for uploads in MB |

### Security Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `users` | map[string]string | (optional) | Username to password mapping |
| `rate_limit_rps` | float | 10.0 | Auth endpoint rate limit (requests per second per IP) |
| `allowed_origins` | string[] | [] | WebSocket/CORS allowed origins (empty = allow all) |

**Example security configuration:**

```yaml
# config.yaml
users:
  admin: "secure-password-here"
  user2: "another-password"

rate_limit_rps: 10

allowed_origins:
  - "http://localhost:3000"
  - "https://myapp.example.com"
  - "*.internal.lan"  # Wildcard subdomain support
```

### Mount Points

Each mount point has:

| Option | Type | Required | Description |
|--------|------|----------|-------------|
| `name` | string | Yes | Display name and URL path prefix |
| `path` | string | Yes | Absolute filesystem path |
| `read_only` | bool | No | If true, write operations are blocked |
| `auto_discover` | bool | No | If true, auto-discover subdirectory mount points |

## Environment Variables

Environment variables override config file values. All variables use the `FM_` prefix:

| Variable | Config Key | Description |
|----------|------------|-------------|
| `FM_JWT_SECRET` | jwt_secret | JWT signing secret |
| `FM_PORT` | port | HTTP server port |
| `FM_HOST` | host | Bind address |
| `FM_RATE_LIMIT_RPS` | rate_limit_rps | Rate limit for auth endpoints |
| `FM_ALLOWED_ORIGINS` | allowed_origins | Comma-separated allowed origins |
| `FM_USERS_<username>` | users.<username> | User password (e.g., `FM_USERS_admin=password`) |
| `CONFIG_PATH` | - | Path to config file |

**Example environment setup:**

```bash
export FM_JWT_SECRET="your-secure-random-string"
export FM_USERS_admin="secure-password"
export FM_RATE_LIMIT_RPS="10"
export FM_ALLOWED_ORIGINS="http://localhost:3000,https://myapp.example.com"
```

## Mount Point Configuration

### Basic Setup

```yaml
mount_points:
  - name: "media"
    path: "/data/media"
    read_only: false
```

This creates a mount point accessible at `/api/v1/files/media/*`.

### Read-Only Mounts

For directories that should not be modified:

```yaml
mount_points:
  - name: "backups"
    path: "/mnt/backups"
    read_only: true
```

Write operations (create, rename, delete, upload) will return 403 Forbidden.

### Multiple Mounts

```yaml
mount_points:
  - name: "media"
    path: "/data/media"
    read_only: false

  - name: "documents"
    path: "/home/user/documents"
    read_only: false

  - name: "photos"
    path: "/data/photos"
    read_only: false

  - name: "backups"
    path: "/mnt/backups"
    read_only: true

  - name: "nas"
    path: "/mnt/nas"
    read_only: false
```

### Network Mounts

Network shares (NFS, SMB) work the same as local directories:

```yaml
mount_points:
  - name: "nas"
    path: "/mnt/nas"  # NFS mount point
    read_only: false

  - name: "share"
    path: "/mnt/smb/share"  # SMB mount point
    read_only: false
```

Ensure the mount is available before starting the server.

## Docker Configuration

When using Docker, paths in `config.yaml` are container paths. Map host directories in `docker-compose.yml`:

**config.yaml:**
```yaml
mount_points:
  - name: "media"
    path: "/data/media"  # Container path
```

**docker-compose.yml:**
```yaml
volumes:
  - /host/path/to/media:/data/media  # Host:Container
```

## Security Considerations

### JWT Secret

Always use a strong, random JWT secret in production:

```bash
# Generate a secure secret
openssl rand -base64 32
```

Never commit secrets to version control. Use environment variables:

```bash
export JWT_SECRET="your-secure-random-string"
```

### Mount Point Security

1. **Principle of Least Privilege**: Only mount directories that need to be accessible
2. **Read-Only When Possible**: Use `read_only: true` for backup directories
3. **Avoid System Directories**: Never mount `/`, `/etc`, `/var`, etc.
4. **Path Validation**: The server validates all paths to prevent traversal attacks

### Recommended Mount Structure

```yaml
mount_points:
  # User data - read/write
  - name: "media"
    path: "/data/media"
    read_only: false

  - name: "documents"
    path: "/data/documents"
    read_only: false

  # Backups - read only
  - name: "backups"
    path: "/backups"
    read_only: true

  # Shared network storage
  - name: "nas"
    path: "/mnt/nas"
    read_only: false
```

## Validation

The server validates configuration on startup:

1. **Mount point paths must exist** (warning if not)
2. **Mount point names must be unique**
3. **JWT secret must be set**

Check logs for configuration issues:

```bash
docker-compose logs backend | grep -i config
```

## Hot Reload

Configuration changes require a server restart:

```bash
docker-compose restart backend
```

## Troubleshooting

### Mount point not accessible

```bash
# Check if path exists
ls -la /data/media

# Check permissions
stat /data/media

# Check container can access
docker exec filemanager-backend ls -la /data/media
```

### Configuration not loading

```bash
# Check config file syntax
cat backend/config.yaml | python -c "import yaml,sys; yaml.safe_load(sys.stdin)"

# Check environment variables
docker exec filemanager-backend env | grep -E "(JWT|CONFIG)"
```

### Permission denied on write

1. Check `read_only` setting in config
2. Check filesystem permissions
3. Check container user has write access
