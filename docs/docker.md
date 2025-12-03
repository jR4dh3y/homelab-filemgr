# Docker Deployment Guide

This guide explains how to deploy the Homelab File Manager using Docker.

## Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- At least 1GB RAM available

## Quick Start

1. **Copy environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Configure your settings in `.env`:**
   ```bash
   # Generate a secure JWT secret
   JWT_SECRET=$(openssl rand -base64 32)
   
   # Set your mount point paths
   MEDIA_PATH=/path/to/your/media
   DOCUMENTS_PATH=/path/to/your/documents
   BACKUPS_PATH=/path/to/your/backups
   ```

3. **Start the services:**
   ```bash
   docker-compose up -d
   ```

4. **Access the application:**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## Default Credentials

- Username: `admin`
- Password: `admin`

⚠️ **Change these in production!** See [Security Guide](security.md) for instructions.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | (required) | Secret key for JWT tokens |
| `BACKEND_PORT` | 8080 | Backend API port |
| `FRONTEND_PORT` | 3000 | Frontend UI port |
| `ORIGIN` | http://localhost:3000 | Origin URL for CORS |
| `MEDIA_PATH` | /data/media | Host path to media files |
| `DOCUMENTS_PATH` | /home/user/documents | Host path to documents |
| `BACKUPS_PATH` | /mnt/backups | Host path to backups |

## Docker Compose Files

### Development (`docker-compose.yml`)

Basic setup with backend and frontend:

```bash
docker-compose up -d
```

### Production (`docker-compose.prod.yml`)

Includes nginx reverse proxy with HTTPS support:

```bash
docker-compose -f docker-compose.prod.yml up -d
```

## Volume Mounts

Map your host directories to container paths in `docker-compose.yml`:

```yaml
volumes:
  # Host path : Container path
  - /your/media:/data/media
  - /your/documents:/home/user/documents
  - /your/backups:/mnt/backups:ro  # Read-only
```

The container paths must match what's configured in `backend/config.yaml`.

## HTTPS Setup

### Using Let's Encrypt

1. **Obtain certificates:**
   ```bash
   certbot certonly --standalone -d files.yourdomain.com
   ```

2. **Copy certificates:**
   ```bash
   mkdir -p nginx/certs
   cp /etc/letsencrypt/live/files.yourdomain.com/fullchain.pem nginx/certs/
   cp /etc/letsencrypt/live/files.yourdomain.com/privkey.pem nginx/certs/
   ```

3. **Enable HTTPS in nginx config:**
   
   Edit `nginx/nginx.conf` and uncomment the HTTPS server block.

4. **Update environment:**
   ```bash
   ORIGIN=https://files.yourdomain.com
   ```

5. **Restart services:**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

### Using Custom Certificates

Place your certificates in `nginx/certs/`:
- `fullchain.pem` - Certificate chain
- `privkey.pem` - Private key

## File Permissions

The backend container needs appropriate permissions to access mounted directories.

### Option 1: Run as Root (Default)

The container runs as root by default to access files with various ownership. This is the simplest option for homelab use.

### Option 2: Specific UID/GID

If your files have consistent ownership:

```yaml
backend:
  user: "1000:1000"  # Match your file ownership
```

### Option 3: ACLs

Use filesystem ACLs for fine-grained control:

```bash
# Grant access to container user
setfacl -R -m u:1000:rwx /data/media
setfacl -R -d -m u:1000:rwx /data/media
```

## Resource Limits

Default limits in docker-compose.yml:

| Service | Memory Limit | Memory Reserved |
|---------|--------------|-----------------|
| Backend | 512MB | 128MB |
| Frontend | 256MB | 64MB |

Adjust based on your server capacity:

```yaml
deploy:
  resources:
    limits:
      memory: 1G
    reservations:
      memory: 256M
```

## Health Checks

Both services include health checks:

```bash
# View health status
docker-compose ps

# Check specific service
docker inspect --format='{{.State.Health.Status}}' filemanager-backend
```

## Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend

# Last 100 lines
docker-compose logs --tail=100 backend
```

## Updating

```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## Backup

Important data to backup:
- `.env` - Environment configuration
- `backend/config.yaml` - Mount point configuration
- `nginx/certs/` - SSL certificates (if using HTTPS)

## Troubleshooting

### Container won't start

```bash
# Check logs
docker-compose logs backend

# Check container status
docker-compose ps
```

### Permission denied errors

```bash
# Verify mount point exists and is accessible
ls -la /data/media

# Check container user
docker exec filemanager-backend id
```

### WebSocket connection fails

```bash
# Check nginx logs
docker-compose logs nginx

# Verify WebSocket endpoint
curl -i -N -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  http://localhost/api/v1/ws
```

### Large file uploads fail

Increase nginx client_max_body_size in `nginx/nginx.conf`:

```nginx
client_max_body_size 20G;
```

### Out of memory

Increase resource limits or add swap:

```bash
# Add 2GB swap
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```
