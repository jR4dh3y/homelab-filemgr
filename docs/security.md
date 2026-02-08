# Security Guide

This document covers security features, best practices, and hardening recommendations for the Homelab File Manager.

## Security Features

### Authentication

The application uses JWT (JSON Web Tokens) for authentication:

- **Access Tokens**: Short-lived (1 hour default)
- **Refresh Tokens**: Longer-lived for session continuity
- **Token Validation**: Every API request validates the token

### Path Traversal Prevention

All file paths are validated to prevent directory traversal attacks:

```go
// Blocked patterns:
// - ../
// - ..\\
// - %2e%2e%2f (URL encoded)
// - %2e%2e/ (mixed encoding)
```

The validator:
1. Decodes URL-encoded characters
2. Cleans the path (removes `.` and `..`)
3. Verifies the result stays within mount boundaries

### Mount Point Isolation

Files are only accessible within configured mount points:

- Requests outside mount points return 403 Forbidden
- Each mount point can be read-only or read-write
- Mount names are validated against configuration

### Security Headers

All responses include security headers:

| Header | Value | Purpose |
|--------|-------|---------|
| X-Content-Type-Options | nosniff | Prevent MIME sniffing |
| X-Frame-Options | DENY | Prevent clickjacking |
| X-XSS-Protection | 1; mode=block | XSS filter |
| Content-Security-Policy | default-src 'self' | Restrict resource loading |
| Referrer-Policy | strict-origin-when-cross-origin | Control referrer info |

## Configuration Security

### JWT Secret

**Critical**: Always use a strong, random JWT secret in production.

Generate a secure secret:
```bash
openssl rand -base64 32
```

Never:
- Use default secrets in production
- Commit secrets to version control
- Share secrets in logs or error messages

### Environment Variables

Store sensitive configuration in environment variables:

```bash
export JWT_SECRET="your-secure-random-string"
```

Or use a `.env` file (not committed to git):
```
JWT_SECRET=your-secure-random-string
```

### Mount Point Security

Follow the principle of least privilege:

```yaml
mount_points:
  # Only mount what's needed
  - name: "media"
    path: "/data/media"
    read_only: false

  # Use read-only for sensitive data
  - name: "backups"
    path: "/backups"
    read_only: true
```

**Never mount:**
- Root filesystem (`/`)
- System directories (`/etc`, `/var`, `/usr`)
- Home directories with sensitive files
- Docker socket or system sockets

## Docker Security

### Container Privileges

The backend container uses specific capabilities instead of full root:

```yaml
cap_add:
  - DAC_READ_SEARCH   # Read files regardless of permissions
  - CHOWN             # Change file ownership
  - FOWNER            # Bypass permission checks
```

This is more secure than `privileged: true`.

### Security Options

```yaml
security_opt:
  - no-new-privileges:true  # Prevent privilege escalation
```

### Read-Only Filesystem

The frontend container runs with a read-only filesystem:

```yaml
read_only: true
tmpfs:
  - /tmp
```

### Non-Root User

The frontend runs as a non-root user:

```yaml
user: "1001:1001"
```

### Resource Limits

Prevent resource exhaustion:

```yaml
deploy:
  resources:
    limits:
      memory: 512M
```

## Network Security

### HTTPS

Always use HTTPS in production:

1. Obtain SSL certificates (Let's Encrypt recommended)
2. Configure nginx with TLS
3. Enable HSTS

```nginx
ssl_protocols TLSv1.2 TLSv1.3;
ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
add_header Strict-Transport-Security "max-age=63072000" always;
```

### Firewall

Restrict access to necessary ports only:

```bash
# Allow only HTTPS
ufw allow 443/tcp

# Block direct backend access from outside
ufw deny 8080/tcp
```

### Network Isolation

Use Docker networks to isolate services:

```yaml
networks:
  filemanager-net:
    driver: bridge
    internal: false  # Set to true if no external access needed
```

## Authentication Hardening

### Configurable Credentials ✅

Credentials are now configurable via config file or environment variables.

**Via config.yaml:**
```yaml
users:
  admin: "your-secure-password"
  user2: "another-password"
```

**Via environment variables:**
```bash
FM_USERS_admin=your-secure-password
FM_USERS_user2=another-password
```

If no users are configured, the system falls back to `admin:admin` with a warning log.

### Session Management

- Access tokens expire after 15 minutes (configurable)
- Refresh tokens allow session continuity (7 days default)
- Logout invalidates tokens
- Revoked tokens are automatically cleaned up

### Rate Limiting ✅

Rate limiting is now built-in for authentication endpoints:

**Configuration:**
```yaml
# config.yaml
rate_limit_rps: 10  # requests per second per IP
```

**Via environment:**
```bash
FM_RATE_LIMIT_RPS=10
```

Features:
- Per-IP rate limiting using token bucket algorithm
- Supports proxy headers (X-Forwarded-For, X-Real-IP)
- Returns HTTP 429 Too Many Requests when exceeded
- Memory-efficient with automatic cleanup

You can also add additional rate limiting via nginx:

```nginx
limit_req_zone $binary_remote_addr zone=login:10m rate=5r/m;

location /api/v1/auth/login {
    limit_req zone=login burst=3 nodelay;
    proxy_pass http://backend;
}
```

## File Upload Security

### Size Limits

Configure maximum upload size:

```yaml
max_upload_mb: 10240  # 10GB
```

### Checksum Verification

Uploads are verified with SHA256 checksums to ensure integrity.

### Temporary File Handling

- Chunks are stored in a temporary directory
- Incomplete uploads are cleaned up
- Final files are moved atomically

## Logging and Monitoring

### Access Logs

Enable access logging in nginx:

```nginx
access_log /var/log/nginx/access.log main;
```

### Security Events

Monitor for:
- Failed login attempts
- Path traversal attempts (403 errors)
- Unusual file access patterns

### Log Rotation

Configure log rotation to prevent disk exhaustion:

```bash
# /etc/logrotate.d/filemanager
/var/log/nginx/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 www-data adm
}
```

## Security Checklist

### Before Deployment

- [ ] Change JWT secret to a secure random value
- [ ] Change default admin credentials
- [ ] Configure mount points with minimal access
- [ ] Set up HTTPS with valid certificates
- [ ] Configure firewall rules
- [ ] Review Docker security settings

### Regular Maintenance

- [ ] Update dependencies regularly
- [ ] Review access logs for anomalies
- [ ] Rotate secrets periodically
- [ ] Test backup and recovery procedures
- [ ] Audit mount point configurations

### Incident Response

If you suspect a security breach:

1. **Isolate**: Disconnect the server from the network
2. **Preserve**: Save logs before they rotate
3. **Investigate**: Review access logs and file changes
4. **Remediate**: Patch vulnerabilities, rotate secrets
5. **Report**: Document the incident

## Reporting Security Issues

If you discover a security vulnerability:

1. **Do not** open a public issue
2. Email security concerns to [security@example.com]
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

We will respond within 48 hours and work with you to address the issue.
